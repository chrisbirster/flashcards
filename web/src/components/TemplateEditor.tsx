import { useState, useEffect } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { updateTemplate } from '#/lib/api'
import type { NoteType, CardTemplate } from '#/lib/api'
import DOMPurify from 'dompurify'

interface TemplateEditorProps {
  noteType: NoteType
  onClose: () => void
}

type TabType = 'front' | 'back' | 'styling'

export function TemplateEditor({ noteType, onClose }: TemplateEditorProps) {
  const queryClient = useQueryClient()
  const [selectedTemplate, setSelectedTemplate] = useState<CardTemplate>(noteType.templates[0])
  const [activeTab, setActiveTab] = useState<TabType>('front')
  const [qFmt, setQFmt] = useState(selectedTemplate?.qFmt || '')
  const [aFmt, setAFmt] = useState(selectedTemplate?.aFmt || '')
  const [styling, setStyling] = useState(selectedTemplate?.styling || '')
  const [sampleFieldVals, setSampleFieldVals] = useState<Record<string, string>>({})
  const [error, setError] = useState<string | null>(null)
  const [hasChanges, setHasChanges] = useState(false)

  // Initialize sample field values
  useEffect(() => {
    const initialVals: Record<string, string> = {}
    noteType.fields.forEach((field, i) => {
      initialVals[field] = `[${field} sample ${i + 1}]`
    })
    setSampleFieldVals(initialVals)
  }, [noteType.fields])

  // Update local state when template changes
  useEffect(() => {
    if (selectedTemplate) {
      setQFmt(selectedTemplate.qFmt)
      setAFmt(selectedTemplate.aFmt)
      setStyling(selectedTemplate.styling || '')
      setHasChanges(false)
    }
  }, [selectedTemplate])

  const invalidateNoteTypes = () => {
    queryClient.invalidateQueries({ queryKey: ['note-types'] })
    queryClient.invalidateQueries({ queryKey: ['note-type', noteType.name] })
  }

  const updateTemplateMutation = useMutation({
    mutationFn: () =>
      updateTemplate(noteType.name, selectedTemplate.name, {
        qFmt,
        aFmt,
        styling,
      }),
    onSuccess: (data) => {
      // Update selected template with new data
      const updated = data.templates.find(t => t.name === selectedTemplate.name)
      if (updated) {
        setSelectedTemplate(updated)
      }
      setHasChanges(false)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const handleSave = () => {
    updateTemplateMutation.mutate()
  }

  const handleTemplateChange = (templateName: string) => {
    const template = noteType.templates.find(t => t.name === templateName)
    if (template) {
      setSelectedTemplate(template)
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
  const tabLabels: Record<TabType, string> = {
    front: 'Front Template',
    back: 'Back Template',
    styling: 'Styling (CSS)',
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b">
          <h2 className="text-lg font-semibold text-gray-900">
            Edit Templates: {noteType.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
            data-testid="close-template-editor"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="flex-1 overflow-auto p-4">
          {/* Error message */}
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
              {error}
            </div>
          )}

          {/* Template selector (for note types with multiple templates) */}
          {noteType.templates.length > 1 && (
            <div className="mb-4">
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
          )}

          <div className="grid grid-cols-2 gap-4">
            {/* Left side: Editor */}
            <div>
              {/* Tabs */}
              <div className="flex border-b mb-4">
                {(['front', 'back', 'styling'] as TabType[]).map((tab) => (
                  <button
                    key={tab}
                    onClick={() => setActiveTab(tab)}
                    className={`px-4 py-2 text-sm font-medium border-b-2 -mb-px ${
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
                      <label className="text-xs text-gray-500 w-20 truncate" title={field}>
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
                  {activeTab === 'styling' ? (
                    <div className="text-xs font-mono text-gray-600 whitespace-pre-wrap">
                      {styling || '(No custom styling)'}
                    </div>
                  ) : (
                    <>
                      {styling && (
                        <style>{styling}</style>
                      )}
                      <div
                        className="card"
                        dangerouslySetInnerHTML={{
                          __html: DOMPurify.sanitize(
                            renderPreview(activeTab === 'front' ? qFmt : aFmt, sampleFieldVals, activeTab === 'back')
                          ),
                        }}
                        data-testid="preview-content"
                      />
                    </>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex justify-between items-center gap-2 p-4 border-t bg-gray-50">
          <div className="text-sm text-gray-500">
            {hasChanges && <span className="text-amber-600">Unsaved changes</span>}
          </div>
          <div className="flex gap-2">
            <button
              onClick={onClose}
              className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
            >
              Close
            </button>
            <button
              onClick={handleSave}
              disabled={!hasChanges || updateTemplateMutation.isPending}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
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
