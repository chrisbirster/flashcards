import { useState } from 'react'
import { useMutation, useQueryClient, useQuery } from '@tanstack/react-query'
import type { NoteType, CardTemplate } from '#/lib/api'
import DOMPurify from 'dompurify'
import { useAppRepository } from '#/lib/app-repository'

interface TemplateEditorProps {
  noteType: NoteType
  onClose: () => void
}

type TabType = 'front' | 'back' | 'styling'

function buildInitialSampleFieldVals(fields: string[]): Record<string, string> {
  const initialVals: Record<string, string> = {}
  fields.forEach((field, i) => {
    initialVals[field] = `[${field} sample ${i + 1}]`
  })
  return initialVals
}

export function TemplateEditor({ noteType, onClose }: TemplateEditorProps) {
  const queryClient = useQueryClient()
  const repository = useAppRepository()
  const [selectedTemplate, setSelectedTemplate] = useState<CardTemplate>(noteType.templates[0])
  const [activeTab, setActiveTab] = useState<TabType>('front')
  const [qFmt, setQFmt] = useState(selectedTemplate?.qFmt || '')
  const [aFmt, setAFmt] = useState(selectedTemplate?.aFmt || '')
  const [styling, setStyling] = useState(selectedTemplate?.styling || '')
  const [ifFieldNonEmpty, setIfFieldNonEmpty] = useState(selectedTemplate?.ifFieldNonEmpty || '')
  const [deckOverride, setDeckOverride] = useState(selectedTemplate?.deckOverride || '')
  const [sampleFieldVals, setSampleFieldVals] = useState<Record<string, string>>(() =>
    buildInitialSampleFieldVals(noteType.fields),
  )
  const [error, setError] = useState<string | null>(null)
  const [hasChanges, setHasChanges] = useState(false)

  // Fetch decks for deck override dropdown
  const { data: decks = [] } = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })

  const applyTemplate = (template: CardTemplate) => {
    setSelectedTemplate(template)
    setQFmt(template.qFmt)
    setAFmt(template.aFmt)
    setStyling(template.styling || '')
    setIfFieldNonEmpty(template.ifFieldNonEmpty || '')
    setDeckOverride(template.deckOverride || '')
    setHasChanges(false)
    setError(null)
  }

  // Validate cloze templates contain {{cloze:Field}} pattern
  const validateClozeTemplate = (frontTemplate: string): string | null => {
    if (!selectedTemplate?.isCloze) return null
    
    const clozePattern = /\{\{cloze:\w+\}\}/
    if (!clozePattern.test(frontTemplate)) {
      return 'Cloze templates must contain at least one {{cloze:FieldName}} tag'
    }
    return null
  }

  const clozeValidationError = validateClozeTemplate(qFmt)

  const handleAutoFixClozeTemplate = () => {
    if (!selectedTemplate?.isCloze) return

    const fallbackField = noteType.fields.includes('Text') ? 'Text' : noteType.fields[0]
    if (!fallbackField) return

    const clozeToken = `{{cloze:${fallbackField}}}`
    const nextTemplate = qFmt.trim() ? `${qFmt}\n${clozeToken}` : clozeToken

    setQFmt(nextTemplate)
    setActiveTab('front')
    setHasChanges(true)
    setError(null)
  }

  const invalidateNoteTypes = () => {
    queryClient.invalidateQueries({ queryKey: ['note-types'] })
    queryClient.invalidateQueries({ queryKey: ['note-type', noteType.name] })
  }

  const syncNoteTypeCache = (updatedTemplates: CardTemplate[]) => {
    queryClient.setQueryData<NoteType[]>(['note-types'], (existing) => {
      if (!existing) return existing
      return existing.map((candidate) =>
        candidate.name === noteType.name
          ? { ...candidate, templates: updatedTemplates }
          : candidate,
      )
    })

    queryClient.setQueryData<NoteType>(['note-type', noteType.name], (existing) => {
      if (!existing) return existing
      return { ...existing, templates: updatedTemplates }
    })
  }

  const handleTemplatesResponse = (templates: CardTemplate[], preferredTemplateName?: string) => {
    syncNoteTypeCache(templates)
    const nextTemplate =
      templates.find((template) => template.name === preferredTemplateName) ??
      templates.find((template) => template.name === selectedTemplate.name) ??
      templates[0]

    if (nextTemplate) {
      applyTemplate(nextTemplate)
    }
    invalidateNoteTypes()
  }

  const updateTemplateMutation = useMutation({
    mutationFn: () =>
      repository.updateTemplate(noteType.name, selectedTemplate.name, {
        qFmt,
        aFmt,
        styling,
        ifFieldNonEmpty: ifFieldNonEmpty || undefined,
        deckOverride: deckOverride || undefined,
      }),
    onSuccess: (data) => {
      handleTemplatesResponse(data.templates)
    },
    onError: (err: Error) => setError(err.message),
  })

  const createTemplateMutation = useMutation({
    mutationFn: (req: {name: string; sourceTemplateName?: string}) => repository.createTemplate(noteType.name, req),
    onSuccess: (data, variables) => {
      handleTemplatesResponse(data.templates, variables.name)
      setError(null)
    },
    onError: (err: Error) => setError(err.message),
  })

  const renameTemplateMutation = useMutation({
    mutationFn: (nextName: string) =>
      repository.updateTemplate(noteType.name, selectedTemplate.name, { name: nextName }),
    onSuccess: (data, nextName) => {
      handleTemplatesResponse(data.templates, nextName)
      setError(null)
    },
    onError: (err: Error) => setError(err.message),
  })

  const deleteTemplateMutation = useMutation({
    mutationFn: () => repository.deleteTemplate(noteType.name, selectedTemplate.name),
    onSuccess: (data) => {
      handleTemplatesResponse(data.templates)
      setError(null)
    },
    onError: (err: Error) => setError(err.message),
  })

  const handleSave = () => {
    if (clozeValidationError) {
      setError(clozeValidationError)
      return
    }
    updateTemplateMutation.mutate()
  }

  const handleTemplateChange = (templateName: string) => {
    const template = noteType.templates.find(t => t.name === templateName)
    if (template) {
      applyTemplate(template)
    }
  }

  const handleFieldChange = (tab: TabType, value: string) => {
    setHasChanges(true)
    if (tab === 'front') setQFmt(value)
    else if (tab === 'back') setAFmt(value)
    else setStyling(value)
  }

  const handleSampleFieldChange = (field: string, value: string) => {
    setSampleFieldVals(prev => ({ ...prev, [field]: value }))
  }

  const handleCreateTemplate = () => {
    const nextName = window.prompt('Template name')
    if (!nextName?.trim()) return
    createTemplateMutation.mutate({ name: nextName.trim() })
  }

  const handleDuplicateTemplate = () => {
    const suggestedName = `${selectedTemplate.name} Copy`
    const nextName = window.prompt('Duplicate template as', suggestedName)
    if (!nextName?.trim()) return
    createTemplateMutation.mutate({
      name: nextName.trim(),
      sourceTemplateName: selectedTemplate.name,
    })
  }

  const handleRenameTemplate = () => {
    const nextName = window.prompt('Rename template', selectedTemplate.name)
    if (!nextName?.trim() || nextName.trim() === selectedTemplate.name) return
    renameTemplateMutation.mutate(nextName.trim())
  }

  const handleDeleteTemplate = () => {
    if (!window.confirm(`Delete the template "${selectedTemplate.name}"?`)) return
    deleteTemplateMutation.mutate()
  }

  // Render template with sample values (non-recursive version)
  const renderPreview = (template: string, fieldVals: Record<string, string>, isBack: boolean = false): string => {
    let result = template
    for (const [field, value] of Object.entries(fieldVals)) {
      const regex = new RegExp(`\\{\\{${field}\\}\\}`, 'g')
      result = result.replace(regex, value || `[${field}]`)
    }
    // Handle special tokens like {{FrontSide}} (only on back template, no recursion)
    if (isBack) {
      // Replace {{FrontSide}} with rendered front template (one level only)
      let frontRendered = qFmt
      for (const [field, value] of Object.entries(fieldVals)) {
        const regex = new RegExp(`\\{\\{${field}\\}\\}`, 'g')
        frontRendered = frontRendered.replace(regex, value || `[${field}]`)
      }
      result = result.replace(/\{\{FrontSide\}\}/g, frontRendered)
    }
    // Handle cloze:Field tokens
    result = result.replace(/\{\{cloze:(\w+)\}\}/g, (_, fieldName) => {
      return fieldVals[fieldName] || `[${fieldName}]`
    })
    return result
  }

  const currentValue = activeTab === 'front' ? qFmt : activeTab === 'back' ? aFmt : styling
  const frontPreviewHtml = DOMPurify.sanitize(renderPreview(qFmt, sampleFieldVals))
  const backPreviewHtml = DOMPurify.sanitize(renderPreview(aFmt, sampleFieldVals, true))
  const tabLabels: Record<TabType, string> = {
    front: 'Front Template',
    back: 'Back Template',
    styling: 'Styling (CSS)',
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-end sm:items-center justify-center z-50 p-2 sm:p-0">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-0 sm:mx-4 h-[95vh] sm:h-auto max-h-[95vh] sm:max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-start sm:items-center justify-between gap-3 p-3 sm:p-4 border-b">
          <h2 className="text-base sm:text-lg font-semibold text-gray-900">
            Edit Templates: {noteType.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600 shrink-0"
            data-testid="close-template-editor"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="flex-1 overflow-auto p-3 sm:p-4">
          {/* Error message */}
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
              {error}
            </div>
          )}

          {/* Cloze validation warning */}
          {clozeValidationError && !error && (
            <div className="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-md text-yellow-800 text-sm">
              <div className="flex items-start gap-2">
                <svg className="w-5 h-5 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
                </svg>
                <div className="flex-1">
                  <div>{clozeValidationError}</div>
                  <button
                    type="button"
                    onClick={handleAutoFixClozeTemplate}
                    className="mt-2 px-2 py-1 text-xs font-medium bg-yellow-100 text-yellow-900 rounded hover:bg-yellow-200"
                    data-testid="auto-fix-cloze-template"
                  >
                    Auto-fix
                  </button>
                </div>
              </div>
            </div>
          )}

          <div className="mb-4 space-y-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Template
              </label>
              <select
                value={selectedTemplate?.name}
                onChange={(e) => handleTemplateChange(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                data-testid="template-selector"
              >
                {noteType.templates.map((t) => (
                  <option key={t.name} value={t.name}>
                    {t.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                onClick={handleCreateTemplate}
                disabled={createTemplateMutation.isPending}
                className="rounded-md border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60"
              >
                New template
              </button>
              <button
                type="button"
                onClick={handleDuplicateTemplate}
                disabled={createTemplateMutation.isPending}
                className="rounded-md border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60"
              >
                Duplicate
              </button>
              <button
                type="button"
                onClick={handleRenameTemplate}
                disabled={renameTemplateMutation.isPending}
                className="rounded-md border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-60"
              >
                Rename
              </button>
              <button
                type="button"
                onClick={handleDeleteTemplate}
                disabled={noteType.templates.length <= 1 || deleteTemplateMutation.isPending}
                className="rounded-md border border-red-200 px-3 py-2 text-sm font-medium text-red-700 hover:bg-red-50 disabled:opacity-60"
              >
                Delete
              </button>
            </div>
          </div>

          {/* Conditional generation */}
          {!selectedTemplate?.isCloze && (
            <div className="mb-4 p-3 bg-blue-50 border border-blue-200 rounded-md">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Conditional Generation
              </label>
              <div className="flex flex-col sm:flex-row sm:items-center gap-2">
                <span className="text-sm text-gray-600">Only generate card if</span>
                <select
                  value={ifFieldNonEmpty}
                  onChange={(e) => {
                    setIfFieldNonEmpty(e.target.value)
                    setHasChanges(true)
                  }}
                  className="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  data-testid="if-field-non-empty"
                >
                  <option value="">(always generate)</option>
                  {noteType.fields.map((field) => (
                    <option key={field} value={field}>
                      {field}
                    </option>
                  ))}
                </select>
                <span className="text-sm text-gray-600">is not empty</span>
              </div>
              <p className="mt-2 text-xs text-gray-500">
                If set, this template will only create a card when the selected field has content.
              </p>
            </div>
          )}

          {/* Deck override */}
          <div className="mb-4 p-3 bg-purple-50 border border-purple-200 rounded-md">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Deck Override
            </label>
            <select
              value={deckOverride}
              onChange={(e) => {
                setDeckOverride(e.target.value)
                setHasChanges(true)
              }}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              data-testid="deck-override"
            >
              <option value="">(use note's deck)</option>
              {decks.map((deck) => (
                <option key={deck.id} value={deck.name}>
                  {deck.name}
                </option>
              ))}
            </select>
            <p className="mt-2 text-xs text-gray-500">
              If set, cards from this template will be placed in the specified deck instead of the note's deck.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Left side: Editor */}
            <div>
              {/* Tabs */}
              <div className="flex border-b mb-4 overflow-x-auto">
                {(['front', 'back', 'styling'] as TabType[]).map((tab) => (
                  <button
                    key={tab}
                    onClick={() => setActiveTab(tab)}
                    className={`px-3 sm:px-4 py-2 text-sm font-medium border-b-2 -mb-px whitespace-nowrap ${
                      activeTab === tab
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700'
                    }`}
                    data-testid={`tab-${tab}`}
                  >
                    {tabLabels[tab]}
                  </button>
                ))}
              </div>

              {/* Editor textarea */}
              <textarea
                value={currentValue}
                onChange={(e) => handleFieldChange(activeTab, e.target.value)}
                placeholder={
                  activeTab === 'styling'
                    ? '.card {\n  font-family: sans-serif;\n  font-size: 20px;\n}'
                    : `Enter ${activeTab} template HTML...`
                }
                rows={10}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 resize-y font-mono text-sm"
                data-testid={`editor-${activeTab}`}
              />

              <p className="mt-2 text-xs text-gray-500">
                {activeTab === 'styling' ? (
                  <>CSS styling applied to both front and back of cards.</>
                ) : selectedTemplate?.isCloze && activeTab === 'front' ? (
                  <>
                    Use {'{{cloze:FieldName}}'} to mark text for cloze deletion. Regular field tags like {'{{FieldName}}'} also work.
                  </>
                ) : (
                  <>
                    Use {'{{FieldName}}'} to insert field values.
                    {activeTab === 'back' && ' Use {{FrontSide}} to include the front template.'}
                  </>
                )}
              </p>
            </div>

            {/* Right side: Preview */}
            <div>
              <h3 className="text-sm font-medium text-gray-700 mb-2">Live Preview</h3>

              {/* Sample field inputs */}
              <div className="mb-4 p-3 bg-gray-50 rounded-md">
                <h4 className="text-xs font-medium text-gray-600 mb-2">Sample Values</h4>
                <div className="space-y-2">
                  {noteType.fields.map((field) => (
                    <div key={field} className="flex items-center gap-2">
                      <label className="text-xs text-gray-500 w-16 sm:w-20 truncate" title={field}>
                        {field}:
                      </label>
                      <input
                        type="text"
                        value={sampleFieldVals[field] || ''}
                        onChange={(e) => handleSampleFieldChange(field, e.target.value)}
                        className="flex-1 px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
                        data-testid={`sample-${field}`}
                      />
                    </div>
                  ))}
                </div>
              </div>

              {/* Preview pane */}
              <div className="border rounded-md overflow-hidden">
                <div className="bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 border-b">
                  {activeTab === 'front' ? 'Front Preview' : activeTab === 'back' ? 'Back Preview' : 'Styling Preview'}
                </div>
                <div className="p-4 min-h-[200px] bg-white">
                  {styling && (
                    <style>{styling}</style>
                  )}
                  {activeTab === 'styling' ? (
                    <div className="space-y-4" data-testid="styling-preview">
                      <p className="text-xs text-gray-500">
                        Previewing the current front and back card output with this CSS applied.
                      </p>
                      <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
                        <div className="overflow-hidden rounded-md border border-gray-200">
                          <div className="border-b bg-gray-50 px-3 py-2 text-xs font-medium uppercase tracking-[0.2em] text-gray-500">
                            Front card
                          </div>
                          <div className="bg-slate-50 p-4">
                            <div
                              className="card"
                              dangerouslySetInnerHTML={{ __html: frontPreviewHtml }}
                              data-testid="preview-content-front"
                            />
                          </div>
                        </div>
                        <div className="overflow-hidden rounded-md border border-gray-200">
                          <div className="border-b bg-gray-50 px-3 py-2 text-xs font-medium uppercase tracking-[0.2em] text-gray-500">
                            Back card
                          </div>
                          <div className="bg-slate-50 p-4">
                            <div
                              className="card"
                              dangerouslySetInnerHTML={{ __html: backPreviewHtml }}
                              data-testid="preview-content-back"
                            />
                          </div>
                        </div>
                      </div>
                      {!styling.trim() && (
                        <p className="text-xs text-gray-500">(No custom styling yet. Add CSS on the left to change the card preview.)</p>
                      )}
                    </div>
                  ) : (
                    <div
                      className="card"
                      dangerouslySetInnerHTML={{
                        __html: activeTab === 'front' ? frontPreviewHtml : backPreviewHtml,
                      }}
                      data-testid="preview-content"
                    />
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 p-3 sm:p-4 border-t bg-gray-50">
          <div className="text-sm text-gray-500 min-h-5">
            {hasChanges && <span className="text-amber-600">Unsaved changes</span>}
          </div>
          <div className="flex w-full sm:w-auto flex-col-reverse sm:flex-row gap-2">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 w-full sm:w-auto"
            >
              Close
            </button>
            <button
              onClick={handleSave}
              disabled={!hasChanges || updateTemplateMutation.isPending || !!clozeValidationError}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed w-full sm:w-auto"
              data-testid="save-template"
            >
              {updateTemplateMutation.isPending ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
