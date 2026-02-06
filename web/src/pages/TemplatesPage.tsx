import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { fetchNoteTypes } from '#/lib/api'
import { TemplateEditor } from '#/components/TemplateEditor'
import type { NoteType } from '#/lib/api'

export function TemplatesPage() {
  const [selectedNoteType, setSelectedNoteType] = useState<NoteType | null>(null)
  const [showEditor, setShowEditor] = useState(false)

  const { data: noteTypes, isLoading, error } = useQuery({
    queryKey: ['note-types'],
    queryFn: fetchNoteTypes,
  })

  const handleEditTemplate = (noteType: NoteType) => {
    setSelectedNoteType(noteType)
    setShowEditor(true)
  }

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto">
        <div className="text-gray-600">Loading note types...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-5xl mx-auto">
        <div className="text-red-600">
          Error: {error instanceof Error ? error.message : 'Failed to load note types'}
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Card Templates</h1>
        <p className="text-gray-600">
          Manage templates for your note types. Templates control how cards are generated and displayed.
        </p>
      </div>

      <div className="grid gap-6">
        {noteTypes && noteTypes.map((noteType) => (
          <div key={noteType.name} className="bg-white rounded-lg shadow border border-gray-200">
            <div className="p-6">
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-gray-900">{noteType.name}</h3>
                  <p className="text-sm text-gray-600 mt-1">
                    {noteType.fields.length} field{noteType.fields.length !== 1 ? 's' : ''}, {' '}
                    {noteType.templates.length} template{noteType.templates.length !== 1 ? 's' : ''}
                  </p>
                </div>
                <button
                  onClick={() => handleEditTemplate(noteType)}
                  className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 text-sm font-medium"
                >
                  Edit Template
                </button>
              </div>

              {/* Fields */}
              <div className="mb-4">
                <h4 className="text-sm font-medium text-gray-700 mb-2">Fields:</h4>
                <div className="flex flex-wrap gap-2">
                  {noteType.fields.map((field) => (
                    <span
                      key={field}
                      className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm"
                    >
                      {field}
                    </span>
                  ))}
                </div>
              </div>

              {/* Templates */}
              <div>
                <h4 className="text-sm font-medium text-gray-700 mb-2">Card Templates:</h4>
                <div className="space-y-2">
                  {noteType.templates.map((template) => (
                    <div
                      key={template.name}
                      className="p-3 bg-gray-50 rounded-md border border-gray-200"
                    >
                      <div className="flex items-center justify-between">
                        <span className="font-medium text-gray-900">{template.name}</span>
                        {template.isCloze && (
                          <span className="px-2 py-1 bg-purple-100 text-purple-700 text-xs rounded">
                            Cloze
                          </span>
                        )}
                      </div>
                      {template.ifFieldNonEmpty && (
                        <p className="text-xs text-gray-500 mt-1">
                          Conditional: Generated only if "{template.ifFieldNonEmpty}" is not empty
                        </p>
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {!noteTypes || noteTypes.length === 0 && (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <p className="text-gray-500">No note types found.</p>
        </div>
      )}

      {/* Template Editor Modal */}
      {showEditor && selectedNoteType && (
        <TemplateEditor
          noteType={selectedNoteType}
          onClose={() => {
            setShowEditor(false)
            setSelectedNoteType(null)
          }}
        />
      )}
    </div>
  )
}
