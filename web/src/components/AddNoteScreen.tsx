import { useState, useEffect, useCallback, useRef, type CSSProperties } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchNoteTypes, fetchDecks, createNote, checkDuplicate } from '#/lib/api'
import type { NoteBrief, FieldOptions } from '#/lib/api'
import { FieldEditor } from './FieldEditor'
import { TemplateEditor } from './TemplateEditor'
import DOMPurify from 'dompurify';
import { TemplateFieldPreview } from './template-field-preview'
import { ErrorMessage, SuccessMessage } from './message'
import { EditFieldIcon } from './edit-field-icon'
import { IconButton } from './edit-field-icon-button'
import { ShowTemplateIcon } from './show-template-icon'


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
        return `<span class="text-blue-600 font-semibold">${answer}</span>`
      } else {
        return `<span class="bg-blue-100 text-blue-800 px-1 rounded">[${hint || '...'}]</span>`
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
  const queryClient = useQueryClient()
  const [selectedNoteType, setSelectedNoteType] = useState<string>('')
  const [selectedDeckId, setSelectedDeckId] = useState<number>(deckId || 0)
  const [fieldValues, setFieldValues] = useState<Record<string, string>>({})
  const [tags, setTags] = useState<string>('')
  const [duplicates, setDuplicates] = useState<NoteBrief[]>([])
  const [isCheckingDuplicate, setIsCheckingDuplicate] = useState(false)
  const [showDuplicateWarning, setShowDuplicateWarning] = useState(false)
  const [activeField, setActiveField] = useState<string | null>(null)
  const [showFieldEditor, setShowFieldEditor] = useState(false)
  const [showTemplateEditor, setShowTemplateEditor] = useState(false)
  const textareaRefs = useRef<Record<string, HTMLTextAreaElement | null>>({})

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
  }, [activeField, fieldValues])

  const { data: noteTypes, isLoading: loadingNoteTypes } = useQuery({
    queryKey: ['note-types'],
    queryFn: fetchNoteTypes,
  })

  const { data: decks, isLoading: loadingDecks } = useQuery({
    queryKey: ['decks'],
    queryFn: fetchDecks,
  })

  // Set default note type when loaded
  useEffect(() => {
    if (noteTypes && noteTypes.length > 0 && !selectedNoteType) {
      setSelectedNoteType(noteTypes[0].name)
    }
  }, [noteTypes, selectedNoteType])

  // Set default deck when loaded
  useEffect(() => {
    if (decks && decks.length > 0 && !selectedDeckId) {
      setSelectedDeckId(decks[0].id)
    }
  }, [decks, selectedDeckId])

  // Get current note type
  const currentNoteType = noteTypes?.find(nt => nt.name === selectedNoteType)

  // Reset field values when note type changes
  useEffect(() => {
    if (currentNoteType) {
      const newFieldValues: Record<string, string> = {}
      currentNoteType.fields.forEach(field => {
        newFieldValues[field] = fieldValues[field] || ''
      })
      setFieldValues(newFieldValues)
    }
  }, [currentNoteType?.name])

  // Debounced duplicate check for the first field
  const checkForDuplicates = useCallback(async (fieldName: string, value: string) => {
    if (!value.trim() || !selectedNoteType) {
      setDuplicates([])
      setShowDuplicateWarning(false)
      return
    }

    setIsCheckingDuplicate(true)
    try {
      const result = await checkDuplicate({
        typeId: selectedNoteType,
        fieldName,
        value,
        deckId: selectedDeckId || undefined,
      })
      setDuplicates(result.duplicates || [])
      setShowDuplicateWarning(result.isDuplicate)
    } catch {
      // Silently fail duplicate check - not critical
      setDuplicates([])
      setShowDuplicateWarning(false)
    } finally {
      setIsCheckingDuplicate(false)
    }
  }, [selectedNoteType, selectedDeckId])

  // Debounce duplicate check when first field changes
  useEffect(() => {
    if (!currentNoteType) return
    const firstField = currentNoteType.fields[0]
    if (!firstField) return

    const firstFieldValue = fieldValues[firstField]
    const timeoutId = setTimeout(() => {
      checkForDuplicates(firstField, firstFieldValue || '')
    }, 500)

    return () => clearTimeout(timeoutId)
  }, [fieldValues, currentNoteType, checkForDuplicates])

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
    mutationFn: createNote,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats', selectedDeckId] })
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

  // Check if required fields have content
  const hasRequiredContent = () => {
    if (!currentNoteType) return false
    // At least one field should have content
    return currentNoteType.fields.some(field => fieldValues[field]?.trim())
  }

  if (loadingNoteTypes || loadingDecks) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-3xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold text-gray-900">Add Note</h1>
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
          >
            Close
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Note Type and Deck Selectors */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Note Type
                </label>
                <div className="flex gap-2">
                  <select
                    value={selectedNoteType}
                    onChange={(e) => setSelectedNoteType(e.target.value)}
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    {noteTypes?.map((nt) => (
                      <option key={nt.name} value={nt.name}>
                        {nt.name}
                      </option>
                    ))}
                  </select>
                  <IconButton
                    title="Edit Fields"
                    testId="edit-fields-button"
                    handleClick={() => setShowFieldEditor(true)}
                    icon={<EditFieldIcon />}
                  />
                  <IconButton
                    title="Edit Templates"
                    testId="edit-templates-button"
                    handleClick={() => setShowTemplateEditor(true)}
                    icon={<ShowTemplateIcon />}
                  />
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Deck
                </label>
                <select
                  value={selectedDeckId}
                  onChange={(e) => setSelectedDeckId(Number(e.target.value))}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  {decks?.map((deck) => (
                    <option key={deck.id} value={deck.id}>
                      {deck.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
          </div>

          {/* Field Inputs */}
          <div className="bg-white rounded-lg shadow p-6">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold text-gray-900">Fields</h2>
              {/* Cloze toolbar - only shown for Cloze note type */}
              {isClozeType && (
                <div className="flex items-center gap-2">
                  <button
                    type="button"
                    onClick={insertCloze}
                    disabled={!activeField}
                    className="px-3 py-1.5 bg-blue-100 text-blue-700 text-sm font-medium rounded-md hover:bg-blue-200 disabled:bg-gray-100 disabled:text-gray-400 disabled:cursor-not-allowed flex items-center gap-1"
                    title="Add Cloze (Ctrl+Shift+C)"
                    data-testid="add-cloze-button"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    [...] Cloze
                  </button>
                  <span className="text-xs text-gray-400">Ctrl+Shift+C</span>
                </div>
              )}
            </div>
            <div className="space-y-4">
              {currentNoteType?.fields.map((field) => (
                <div key={field}>
                  {(() => {
                    const options = currentNoteType?.fieldOptions?.[field]
                    const isHtmlEditor = options?.htmlEditor || false
                    const fieldStyle = buildFieldEditorStyle(options)

                    return (
                      <>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    {field}
                  </label>
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
                    rows={3}
                    className={`w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 resize-y ${
                      isHtmlEditor ? '' : 'font-mono'
                    }`}
                    style={fieldStyle}
                    dir={options?.rtl ? 'rtl' : 'ltr'}
                  />
                  {isHtmlEditor && (
                    <p className="mt-1 text-xs text-gray-500">HTML editor default enabled for this field.</p>
                  )}
                  {isClozeType && field === 'Text' && (
                    <p className="mt-1 text-xs text-gray-500">
                      Select text and click "[...] Cloze" or press Ctrl+Shift+C to create cloze deletions
                    </p>
                  )}
                      </>
                    )
                  })()}
                </div>
              ))}
            </div>
          </div>

          {/* Duplicate Warning */}
          {showDuplicateWarning && duplicates.length > 0 && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4" data-testid="duplicate-warning">
              <div className="flex items-start">
                <svg className="w-5 h-5 text-yellow-600 mt-0.5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div className="flex-1">
                  <h3 className="text-sm font-medium text-yellow-800">
                    Possible duplicate found
                  </h3>
                  <div className="mt-2 text-sm text-yellow-700">
                    <p>A similar note already exists:</p>
                    <ul className="mt-1 list-disc list-inside">
                      {duplicates.slice(0, 3).map((dup) => {
                        const firstFieldKey = Object.keys(dup.fieldVals)[0]
                        const preview = dup.fieldVals[firstFieldKey] || ''
                        return (
                          <li key={dup.id} className="truncate">
                            {preview.substring(0, 50)}{preview.length > 50 ? '...' : ''}
                          </li>
                        )
                      })}
                    </ul>
                    {duplicates.length > 3 && (
                      <p className="mt-1 text-xs">...and {duplicates.length - 3} more</p>
                    )}
                  </div>
                  <p className="mt-2 text-xs text-yellow-600">
                    You can still add this note if it's intentional.
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Checking indicator */}
          {isCheckingDuplicate && (
            <div className="text-sm text-gray-500 flex items-center gap-2">
              <svg className="animate-spin h-4 w-4 text-gray-400" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Checking for duplicates...
            </div>
          )}

          {/* Tags */}
          <div className="bg-white rounded-lg shadow p-6">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Tags
            </label>
            <input
              type="text"
              value={tags}
              onChange={(e) => setTags(e.target.value)}
              placeholder="Enter tags (comma or space separated)..."
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <p className="mt-1 text-sm text-gray-500">
              Separate multiple tags with commas or spaces
            </p>
          </div>

          {/* Actions */}
          <div className="flex gap-4">
            <button
              type="submit"
              disabled={createNoteMutation.isPending || !hasRequiredContent()}
              className="flex-1 px-6 py-3 bg-blue-600 text-white font-medium rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              {createNoteMutation.isPending ? 'Adding...' : 'Add Note'}
            </button>
            <button
              type="button"
              onClick={onClose}
              className="px-6 py-3 text-gray-700 bg-gray-100 font-medium rounded-lg hover:bg-gray-200"
            >
              Cancel
            </button>
          </div>

          {/* Error message */}
          {createNoteMutation.isError && <ErrorMessage />}
          {/* Success message */}
          {createNoteMutation.isSuccess && <SuccessMessage />}
        </form>

        {/* Field Editor Modal */}
        {showFieldEditor && currentNoteType && (
          <FieldEditor
            noteType={currentNoteType}
            onClose={() => setShowFieldEditor(false)}
          />
        )}

        {/* Template Editor Modal */}
        {showTemplateEditor && currentNoteType && (
          <TemplateEditor
            noteType={currentNoteType}
            onClose={() => setShowTemplateEditor(false)}
          />
        )}

        {/* Preview section */}
        {currentNoteType && hasRequiredContent() && (
          <div className="mt-6 bg-white rounded-lg shadow p-6">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Preview</h2>
            <div className="space-y-4">
              {isClozeType ? (
                // Cloze preview - show one card per cloze ordinal
                (() => {
                  const textField = fieldValues['Text'] || ''
                  const ordinals = extractClozeOrdinals(textField)

                  if (ordinals.length === 0) {
                    return (
                      <div className="text-sm text-gray-500 italic">
                        No cloze deletions found. Use {"{{c1::text}}"} syntax or click the Cloze button.
                      </div>
                    )
                  }

                  return ordinals.map((ordinal) => {
                    const frontHtml = renderClozePreview(textField, ordinal, false)
                    const backHtml = renderClozePreview(textField, ordinal, true)

                    return (
                      <div key={ordinal} className="border rounded-lg p-4">
                        <div className="text-sm text-gray-500 mb-2">
                          Card {ordinal}: Cloze {ordinal}
                        </div>
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <div className="text-xs text-gray-400 mb-1">Front (hidden)</div>
                            <div
                              className="p-2 bg-gray-50 rounded text-sm whitespace-pre-wrap"
                              dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(frontHtml) }}
                            />
                          </div>
                          <div>
                            <div className="text-xs text-gray-400 mb-1">Back (revealed)</div>
                            <div
                              className="p-2 bg-gray-50 rounded text-sm whitespace-pre-wrap"
                              dangerouslySetInnerHTML={{
                                __html:
                                  DOMPurify.sanitize(backHtml)
                              }}
                            />
                          </div>
                        </div>
                      </div>
                    )
                  })
                })()
              ) : (
                // Regular note type preview
                currentNoteType.templates.map((template, idx) => {
                  // Simple preview - replace field placeholders
                  let front = template.qFmt
                  let back = template.aFmt
                  Object.entries(fieldValues).forEach(([field, value]) => {
                    const regex = new RegExp(`\\{\\{${field}\\}\\}`, 'g')
                    front = front.replace(regex, value || `[${field}]`)
                    back = back.replace(regex, value || `[${field}]`)
                  })

                  return (
                    <div key={idx} className="border rounded-lg p-4">
                      <div className="text-sm text-gray-500 mb-2">
                        Card {idx + 1}: {template.name}
                      </div>
                      <div className="grid grid-cols-2 gap-4">
                        <TemplateFieldPreview previewContent={front} label="Front" />
                        <TemplateFieldPreview previewContent={back} label="Back" />
                      </div>
                    </div>
                  )
                })
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
