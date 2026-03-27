import { useState, useEffect, useCallback, useRef, type CSSProperties } from 'react'
import { useNavigate } from 'react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { DeckNotesResponse, NoteBrief, FieldOptions, RecentDeckNoteSummary } from '#/lib/api'
import { AddNoteFormProvider } from './add-note-form-provider'
import { useAddNoteFormContext } from './add-note-form-context'
import DOMPurify from 'dompurify'
import { TemplateFieldPreview } from './template-field-preview'
import { ErrorMessage, SuccessMessage } from './message'
import { EditFieldIcon } from './edit-field-icon'
import { ShowTemplateIcon } from './show-template-icon'
import { RecentDeckNotesPanel } from './recent-deck-notes-panel'
import { AICardSuggestionPanel } from './ai-suggestion-panel'
import { useAppRepository } from '#/lib/app-repository'
import { ActionBar, FormActions } from '#/components/action-bar'
import { FieldRow } from '#/components/field-row'
import { EmptyState, PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'
import { Sheet } from '#/components/sheet'
import { useMediaQuery } from '#/lib/use-media-query'


// Helper to find the next cloze number in text
function getNextClozeNumber(text: string): number {
  const matches = text.match(/\{\{c(\d+)::/g) || []
  const numbers = matches.map(m => parseInt(m.match(/\d+/)?.[0] || '0'))
  return numbers.length > 0 ? Math.max(...numbers) + 1 : 1
}

// Helper to render cloze preview (show [...] for hidden, show text for revealed)
function renderClozePreview(text: string, targetOrdinal: number, reveal: boolean): string {
  // Pattern: {{c1::answer}} or {{c1::answer::hint}}
  return text.replace(/\{\{c(\d+)::([^}]*?)(?:::([^}]*?))?\}\}/g, (_, num, answer, hint) => {
    const clozeNum = parseInt(num)
    if (clozeNum === targetOrdinal) {
      if (reveal) {
        return `<span style="color: var(--app-accent); font-weight: 600;">${answer}</span>`
      } else {
        return `<span style="background: var(--app-muted-surface); color: var(--app-accent); padding: 0 0.35rem; border-radius: 0.4rem;">[${hint || '...'}]</span>`
      }
    }
    return answer
  })
}

// Get all cloze ordinals from text
function extractClozeOrdinals(text: string): number[] {
  const ordinals = new Set<number>()
  const regex = /\{\{c(\d+)::/g
  let match
  while ((match = regex.exec(text)) !== null) {
    ordinals.add(parseInt(match[1]))
  }
  return Array.from(ordinals).sort((a, b) => a - b)
}

interface AddNoteScreenProps {
  deckId?: number
  onClose: () => void
  onSuccess?: () => void
}

function buildAISourceText(fieldNames: string[], values: Record<string, string>) {
  return fieldNames
    .map((field) => {
      const value = values[field]?.trim()
      if (!value) return ''
      return `${field}: ${value}`
    })
    .filter(Boolean)
    .join('\n')
}

function buildFieldEditorStyle(options?: FieldOptions): CSSProperties {
  const style: CSSProperties = {}

  if (!options) {
    return style
  }

  if (options.font) {
    style.fontFamily = options.font
  }
  if (options.fontSize) {
    style.fontSize = options.fontSize
  }
  if (options.rtl) {
    style.direction = 'rtl'
  }

  return style
}

export function AddNoteScreen({ deckId, onClose, onSuccess }: AddNoteScreenProps) {
  return (
    <AddNoteFormProvider deckId={deckId}>
      <AddNoteScreenContent onClose={onClose} onSuccess={onSuccess} />
    </AddNoteFormProvider>
  )
}

function AddNoteScreenContent({ onClose, onSuccess }: Omit<AddNoteScreenProps, 'deckId'>) {
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const repository = useAppRepository()
  const {
    selectedNoteType,
    setSelectedNoteType,
    selectedDeckId,
    setSelectedDeckId,
    fieldValues,
    setFieldValues,
    tags,
    setTags,
    activeField,
    setActiveField,
    textareaRefs,
  } = useAddNoteFormContext()
  const [duplicates, setDuplicates] = useState<NoteBrief[]>([])
  const [isCheckingDuplicate, setIsCheckingDuplicate] = useState(false)
  const [showDuplicateWarning, setShowDuplicateWarning] = useState(false)
  const [aiOpen, setAiOpen] = useState(false)
  const isDesktopAI = useMediaQuery('(min-width: 768px)')
  const lastDuplicateCheckKeyRef = useRef('')
  const duplicateCheckRequestIdRef = useRef(0)

  // Check if current note type is Cloze
  const isClozeType = selectedNoteType === 'Cloze'

  // Insert cloze deletion at cursor position in active textarea
  const insertCloze = useCallback(() => {
    if (!activeField) return
    const textarea = textareaRefs.current[activeField]
    if (!textarea) return

    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const text = fieldValues[activeField] || ''
    const selectedText = text.substring(start, end)

    // Find next cloze number based on all text fields
    const allText = Object.values(fieldValues).join(' ')
    const nextNum = getNextClozeNumber(allText)

    // Wrap selected text or insert placeholder
    const clozeText = selectedText ? `{{c${nextNum}::${selectedText}}}` : `{{c${nextNum}::}}`
    const newText = text.substring(0, start) + clozeText + text.substring(end)

    setFieldValues(prev => ({ ...prev, [activeField]: newText }))

    // Move cursor inside the cloze if no text was selected
    setTimeout(() => {
      if (!selectedText && textarea) {
        const cursorPos = start + `{{c${nextNum}::`.length
        textarea.setSelectionRange(cursorPos, cursorPos)
        textarea.focus()
      }
    }, 0)
  }, [activeField, fieldValues, setFieldValues, textareaRefs])

  const { data: noteTypes, isLoading: loadingNoteTypes } = useQuery({
    queryKey: ['note-types'],
    queryFn: () => repository.fetchNoteTypes(),
  })

  const { data: decks, isLoading: loadingDecks } = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })

  // Set default note type when loaded
  useEffect(() => {
    if (noteTypes && noteTypes.length > 0 && !selectedNoteType) {
      setSelectedNoteType(noteTypes[0].name)
    }
  }, [noteTypes, selectedNoteType, setSelectedNoteType])

  // Set default deck when loaded
  useEffect(() => {
    if (decks && decks.length > 0 && !selectedDeckId) {
      setSelectedDeckId(decks[0].id)
    }
  }, [decks, selectedDeckId, setSelectedDeckId])

  // Get current note type
  const currentNoteType = noteTypes?.find(nt => nt.name === selectedNoteType)
  const aiInitialSource = currentNoteType ? buildAISourceText(currentNoteType.fields, fieldValues) : ''
  const primaryFieldName = currentNoteType?.fields[0] || ''
  const primaryFieldValue = primaryFieldName ? fieldValues[primaryFieldName] || '' : ''

  // Reset field values when note type changes
  useEffect(() => {
    if (currentNoteType) {
      setFieldValues((prev) => {
        const nextFieldValues: Record<string, string> = {}
        currentNoteType.fields.forEach((field) => {
          nextFieldValues[field] = prev[field] || ''
        })

        const hasSameShape = Object.keys(prev).length === currentNoteType.fields.length
        const hasSameValues = currentNoteType.fields.every((field) => prev[field] === nextFieldValues[field])
        if (hasSameShape && hasSameValues) {
          return prev
        }

        return nextFieldValues
      })
    }
  }, [currentNoteType, setFieldValues])

  // Debounced duplicate check for the first field
  const checkForDuplicates = useCallback(async (fieldName: string, value: string) => {
    const normalizedValue = value.trim()

    if (!normalizedValue || !selectedNoteType) {
      lastDuplicateCheckKeyRef.current = ''
      duplicateCheckRequestIdRef.current += 1
      setDuplicates([])
      setShowDuplicateWarning(false)
      setIsCheckingDuplicate(false)
      return
    }

    const requestKey = `${selectedNoteType}:${selectedDeckId || 0}:${fieldName}:${normalizedValue}`
    if (lastDuplicateCheckKeyRef.current === requestKey) {
      return
    }

    lastDuplicateCheckKeyRef.current = requestKey
    const requestId = ++duplicateCheckRequestIdRef.current
    setIsCheckingDuplicate(true)
    try {
      const result = await repository.checkDuplicate({
        typeId: selectedNoteType,
        fieldName,
        value: normalizedValue,
        deckId: selectedDeckId || undefined,
      })
      if (requestId !== duplicateCheckRequestIdRef.current) {
        return
      }
      setDuplicates(result.duplicates || [])
      setShowDuplicateWarning(result.isDuplicate)
    } catch {
      if (requestId !== duplicateCheckRequestIdRef.current) {
        return
      }
      // Silently fail duplicate check - not critical
      setDuplicates([])
      setShowDuplicateWarning(false)
    } finally {
      if (requestId === duplicateCheckRequestIdRef.current) {
        setIsCheckingDuplicate(false)
      }
    }
  }, [repository, selectedNoteType, selectedDeckId])

  // Debounce duplicate check when first field changes
  useEffect(() => {
    if (!primaryFieldName) return
    const timeoutId = setTimeout(() => {
      checkForDuplicates(primaryFieldName, primaryFieldValue)
    }, 500)

    return () => clearTimeout(timeoutId)
  }, [primaryFieldName, primaryFieldValue, checkForDuplicates])

  // Keyboard shortcut for cloze (Ctrl/Cmd+Shift+C)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'c') {
        if (isClozeType && activeField) {
          e.preventDefault()
          insertCloze()
        }
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [isClozeType, activeField, insertCloze])

  const createNoteMutation = useMutation({
    mutationFn: (req: Parameters<typeof repository.createNote>[0]) => repository.createNote(req),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats', selectedDeckId] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
      const primaryFieldName = currentNoteType?.fields[0]
      const fieldPreview = (primaryFieldName && data.note.fieldVals[primaryFieldName]) ||
        Object.values(data.note.fieldVals).find((value) => value.trim()) ||
        ''
      const recentNote: RecentDeckNoteSummary = {
        noteId: data.note.id,
        noteType: data.note.typeId,
        createdAt: data.note.createdAt,
        modifiedAt: data.note.modifiedAt,
        tags: data.note.tags,
        fieldPreview,
        cardCountInDeck: data.cards.filter((card) => card.deckId === selectedDeckId).length || data.cards.length,
      }
      queryClient.setQueryData<DeckNotesResponse>(['deck-notes', selectedDeckId], (existing) => ({
        notes: [recentNote, ...(existing?.notes || []).filter((note) => note.noteId !== recentNote.noteId)].slice(0, 10),
      }))
      // Clear fields for next note
      if (currentNoteType) {
        const newFieldValues: Record<string, string> = {}
        currentNoteType.fields.forEach(field => {
          newFieldValues[field] = ''
        })
        setFieldValues(newFieldValues)
      }
      setTags('')
      setDuplicates([])
      setShowDuplicateWarning(false)
      onSuccess?.()
    },
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!selectedNoteType || !selectedDeckId) return

    // Parse tags (comma or space separated)
    const tagList = tags
      .split(/[,\s]+/)
      .map(t => t.trim())
      .filter(t => t.length > 0)

    createNoteMutation.mutate({
      typeId: selectedNoteType,
      deckId: selectedDeckId,
      fieldVals: fieldValues,
      tags: tagList,
    })
  }

  const handleFieldChange = (field: string, value: string) => {
    setFieldValues(prev => ({ ...prev, [field]: value }))
  }

  const openFieldEditorRoute = () => {
    if (!selectedNoteType) return
    navigate(`note-types/${encodeURIComponent(selectedNoteType)}/fields`)
  }

  const openTemplateEditorRoute = () => {
    if (!selectedNoteType) return
    navigate(`note-types/${encodeURIComponent(selectedNoteType)}/templates`)
  }

  // Check if required fields have content
  const hasRequiredContent = () => {
    if (!currentNoteType) return false
    // At least one field should have content
    return currentNoteType.fields.some(field => fieldValues[field]?.trim())
  }

  const applyAISuggestion = (suggestion: {fieldVals: Record<string, string>}) => {
    if (!currentNoteType) return
    const nextFieldValues: Record<string, string> = {}
    currentNoteType.fields.forEach((field) => {
      nextFieldValues[field] = suggestion.fieldVals[field] || ''
    })
    setFieldValues(nextFieldValues)
    setActiveField(currentNoteType.fields[0] || '')
    lastDuplicateCheckKeyRef.current = ''
    setShowDuplicateWarning(false)
    setDuplicates([])
    setAiOpen(false)
  }

  if (loadingNoteTypes || loadingDecks) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-text-soft)]">
          Loading add-note workspace...
        </PageSection>
      </PageContainer>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-4 sm:p-5">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Create</p>
            <h1 className="mt-2 text-2xl font-semibold tracking-tight text-[var(--app-text)]">Add Note</h1>
            <p className="mt-2 max-w-2xl text-sm leading-6 text-[var(--app-text-soft)]">
              Build notes in a single vertical flow, preview the cards they generate, and keep recent work close by.
            </p>
          </div>
          <button
            onClick={onClose}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Close
          </button>
        </div>
      </PageSection>

      <form onSubmit={handleSubmit} className="space-y-4">
        <PageSection className="p-4 sm:p-5">
          <div className="grid gap-4">
            <FieldRow label="Note type">
              <div className="space-y-3">
                <select
                  value={selectedNoteType}
                  onChange={(e) => setSelectedNoteType(e.target.value)}
                  className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                >
                  {noteTypes?.map((nt) => (
                    <option key={nt.name} value={nt.name}>
                      {nt.name}
                    </option>
                  ))}
                </select>
                <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                  <button
                    type="button"
                    onClick={openFieldEditorRoute}
                    className="inline-flex min-h-11 items-center justify-center gap-2 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
                    data-testid="edit-fields-button"
                  >
                    <EditFieldIcon />
                    Edit fields
                  </button>
                  <button
                    type="button"
                    onClick={openTemplateEditorRoute}
                    className="inline-flex min-h-11 items-center justify-center gap-2 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
                    data-testid="edit-templates-button"
                  >
                    <ShowTemplateIcon />
                    Edit templates
                  </button>
                </div>
                <button
                  type="button"
                  onClick={() => setAiOpen((current) => !current)}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
                >
                  AI suggestions
                </button>
              </div>
            </FieldRow>

            <FieldRow label="Deck">
              <select
                value={selectedDeckId}
                onChange={(e) => setSelectedDeckId(Number(e.target.value))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              >
                {decks?.map((deck) => (
                  <option key={deck.id} value={deck.id}>
                    {deck.name}
                  </option>
                ))}
              </select>
            </FieldRow>
          </div>
        </PageSection>

        {aiOpen && isDesktopAI ? (
          <PageSection className="p-4 sm:p-5">
            <AICardSuggestionPanel
              open={aiOpen}
              noteType={currentNoteType}
              initialSourceText={aiInitialSource}
              existingFieldVals={fieldValues}
              onApplySuggestion={applyAISuggestion}
            />
          </PageSection>
        ) : null}

        <PageSection className="p-4 sm:p-5">
          <div className="flex flex-col gap-3 border-b border-[var(--app-line)] pb-4">
            <div>
              <h2 className="text-lg font-semibold text-[var(--app-text)]">Fields</h2>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Mobile uses a single-column editor so each field stays readable while typing.
              </p>
              {isCheckingDuplicate ? (
                <p className="mt-2 text-xs font-medium uppercase tracking-[0.16em] text-[var(--app-muted)]">
                  Checking duplicates…
                </p>
              ) : null}
            </div>
            {isClozeType && (
              <div className="flex flex-wrap items-center gap-2">
                <button
                  type="button"
                  onClick={insertCloze}
                  disabled={!activeField}
                  className="inline-flex min-h-11 items-center gap-2 rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-55"
                  title="Add Cloze (Ctrl+Shift+C)"
                  data-testid="add-cloze-button"
                >
                  <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                  </svg>
                  Insert cloze
                </button>
                <span className="text-xs text-[var(--app-muted)]">Shortcut: Ctrl+Shift+C</span>
              </div>
            )}
          </div>

          <div className="mt-5 space-y-5">
            {currentNoteType?.fields.map((field) => {
              const options = currentNoteType?.fieldOptions?.[field]
              const isHtmlEditor = options?.htmlEditor || false
              const fieldStyle = buildFieldEditorStyle(options)

              return (
                <FieldRow
                  key={field}
                  label={field}
                  hint={
                    isClozeType && field === 'Text'
                      ? 'Select text and insert cloze markers directly inside the active field.'
                      : isHtmlEditor
                        ? 'HTML editor default enabled for this field.'
                        : undefined
                  }
                >
                  <textarea
                    ref={(el) => { textareaRefs.current[field] = el }}
                    value={fieldValues[field] || ''}
                    onChange={(e) => handleFieldChange(field, e.target.value)}
                    onFocus={() => setActiveField(field)}
                    placeholder={isClozeType && field === 'Text'
                      ? 'Enter text with {{c1::cloze}} deletions...'
                      : isHtmlEditor
                        ? `Enter ${field.toLowerCase()} (HTML)...`
                        : `Enter ${field.toLowerCase()}...`}
                    rows={4}
                    className={`w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)] ${isHtmlEditor ? '' : 'font-mono'}`}
                    style={fieldStyle}
                    dir={options?.rtl ? 'rtl' : 'ltr'}
                  />
                </FieldRow>
              )
            })}
          </div>
        </PageSection>

        {showDuplicateWarning && duplicates.length > 0 && (
          <PageSection className="border-[var(--app-warning-line)] bg-[var(--app-warning-surface)] p-4 sm:p-5" data-testid="duplicate-warning">
            <div className="flex items-start gap-3">
              <svg className="mt-0.5 h-5 w-5 text-[var(--app-warning-text)]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div className="flex-1">
                <h3 className="text-sm font-semibold text-[var(--app-warning-text)]">Possible duplicate found</h3>
                <div className="mt-3 space-y-2 text-sm text-[var(--app-warning-text)]/90">
                  <p>A similar note already exists:</p>
                  <ul className="space-y-1">
                    {duplicates.slice(0, 3).map((dup) => {
                      const firstFieldKey = Object.keys(dup.fieldVals)[0]
                      const preview = dup.fieldVals[firstFieldKey] || ''
                      return (
                        <li key={dup.id} className="rounded-xl border border-[var(--app-warning-line)]/80 bg-black/5 px-3 py-2">
                          {preview.substring(0, 80)}{preview.length > 80 ? '...' : ''}
                        </li>
                      )
                    })}
                  </ul>
                  {duplicates.length > 3 && (
                    <p className="text-xs text-[var(--app-warning-text)]/80">...and {duplicates.length - 3} more</p>
                  )}
                </div>
              </div>
            </div>
          </PageSection>
        )}

        <PageSection className="p-4 sm:p-5">
          <FieldRow label="Tags" hint="Separate multiple tags with commas or spaces.">
            <input
              type="text"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="Enter tags..."
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            />
          </FieldRow>
        </PageSection>

        <PageSection className="overflow-hidden">
          {createNoteMutation.isError ? (
            <div className="border-b border-[var(--app-line)] px-4 py-4 sm:px-5">
              <ErrorMessage message={createNoteMutation.error instanceof Error ? createNoteMutation.error.message : undefined} />
            </div>
          ) : null}
          {createNoteMutation.isSuccess ? (
            <div className="border-b border-[var(--app-line)] px-4 py-4 sm:px-5">
              <SuccessMessage />
            </div>
          ) : null}
          <ActionBar>
            <FormActions>
              <button
                type="button"
                onClick={onClose}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={createNoteMutation.isPending || !hasRequiredContent()}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60"
              >
                {createNoteMutation.isPending ? 'Adding...' : 'Add Note'}
              </button>
            </FormActions>
          </ActionBar>
        </PageSection>
      </form>

      {currentNoteType && hasRequiredContent() && (
        <PageSection className="p-4 sm:p-5">
          <div className="flex items-start justify-between gap-4 border-b border-[var(--app-line)] pb-4">
            <div>
              <h2 className="text-lg font-semibold text-[var(--app-text)]">Preview</h2>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Review the generated card faces before saving the note.
              </p>
            </div>
          </div>
          <div className="mt-5 space-y-4">
            {isClozeType ? (
              (() => {
                const textField = fieldValues['Text'] || ''
                const ordinals = extractClozeOrdinals(textField)

                if (ordinals.length === 0) {
                  return (
                    <EmptyState
                      title="No cloze deletions yet"
                      description='Use the cloze button or enter {{c1::text}} markers to generate preview cards.'
                    />
                  )
                }

                return ordinals.map((ordinal) => {
                  const frontHtml = renderClozePreview(textField, ordinal, false)
                  const backHtml = renderClozePreview(textField, ordinal, true)

                  return (
                    <SurfaceCard key={ordinal} className="space-y-4">
                      <div className="text-sm font-medium text-[var(--app-text-soft)]">Card {ordinal}: Cloze {ordinal}</div>
                      <div className="grid gap-4 md:grid-cols-2">
                        <div className="space-y-2">
                          <div className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Front (hidden)</div>
                          <div
                            className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card-strong)] p-3 text-sm text-[var(--app-text)]"
                            dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(frontHtml) }}
                          />
                        </div>
                        <div className="space-y-2">
                          <div className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Back (revealed)</div>
                          <div
                            className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card-strong)] p-3 text-sm text-[var(--app-text)]"
                            dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(backHtml) }}
                          />
                        </div>
                      </div>
                    </SurfaceCard>
                  )
                })
              })()
            ) : (
              currentNoteType.templates.map((template, idx) => {
                let front = template.qFmt
                let back = template.aFmt
                Object.entries(fieldValues).forEach(([field, value]) => {
                  const regex = new RegExp(`\\{\\{${field}\\}\\}`, 'g')
                  front = front.replace(regex, value || `[${field}]`)
                  back = back.replace(regex, value || `[${field}]`)
                })

                return (
                  <SurfaceCard key={idx} className="space-y-4">
                    <div className="text-sm font-medium text-[var(--app-text-soft)]">Card {idx + 1}: {template.name}</div>
                    <div className="grid gap-4 md:grid-cols-2">
                      <TemplateFieldPreview previewContent={front} label="Front" />
                      <TemplateFieldPreview previewContent={back} label="Back" />
                    </div>
                  </SurfaceCard>
                )
              })
            )}
          </div>
        </PageSection>
      )}

      <Sheet open={aiOpen && !isDesktopAI} onClose={() => setAiOpen(false)} title="AI card suggestions">
        <AICardSuggestionPanel
          open={aiOpen && !isDesktopAI}
          noteType={currentNoteType}
          initialSourceText={aiInitialSource}
          existingFieldVals={fieldValues}
          onApplySuggestion={applyAISuggestion}
        />
      </Sheet>

      <RecentDeckNotesPanel deckId={selectedDeckId} />
    </PageContainer>
  )
}
