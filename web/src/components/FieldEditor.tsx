import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { addField, renameField, removeField, reorderFields } from '#/lib/api'
import type { NoteType } from '#/lib/api'

interface FieldEditorProps {
  noteType: NoteType
  onClose: () => void
}

export function FieldEditor({ noteType, onClose }: FieldEditorProps) {
  const queryClient = useQueryClient()
  const [fields, setFields] = useState<string[]>(noteType.fields)
  const [newFieldName, setNewFieldName] = useState('')
  const [editingField, setEditingField] = useState<string | null>(null)
  const [editFieldName, setEditFieldName] = useState('')
  const [error, setError] = useState<string | null>(null)

  const invalidateNoteTypes = () => {
    queryClient.invalidateQueries({ queryKey: ['note-types'] })
    queryClient.invalidateQueries({ queryKey: ['note-type', noteType.name] })
  }

  const addFieldMutation = useMutation({
    mutationFn: (fieldName: string) => addField(noteType.name, { fieldName }),
    onSuccess: (data) => {
      setFields(data.fields)
      setNewFieldName('')
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const renameFieldMutation = useMutation({
    mutationFn: ({ oldName, newName }: { oldName: string; newName: string }) =>
      renameField(noteType.name, { oldName, newName }),
    onSuccess: (data) => {
      setFields(data.fields)
      setEditingField(null)
      setEditFieldName('')
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const removeFieldMutation = useMutation({
    mutationFn: (fieldName: string) => removeField(noteType.name, { fieldName }),
    onSuccess: (data) => {
      setFields(data.fields)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const reorderFieldsMutation = useMutation({
    mutationFn: (newFields: string[]) => reorderFields(noteType.name, { fields: newFields }),
    onSuccess: (data) => {
      setFields(data.fields)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const handleAddField = (e: React.FormEvent) => {
    e.preventDefault()
    if (!newFieldName.trim()) return
    addFieldMutation.mutate(newFieldName.trim())
  }

  const handleRenameField = (oldName: string) => {
    if (!editFieldName.trim() || editFieldName === oldName) {
      setEditingField(null)
      return
    }
    renameFieldMutation.mutate({ oldName, newName: editFieldName.trim() })
  }

  const handleRemoveField = (fieldName: string) => {
    if (fields.length <= 1) {
      setError('Cannot remove the last field')
      return
    }
    if (confirm(`Are you sure you want to remove the field "${fieldName}"?`)) {
      removeFieldMutation.mutate(fieldName)
    }
  }

  const moveField = (index: number, direction: 'up' | 'down') => {
    const newIndex = direction === 'up' ? index - 1 : index + 1
    if (newIndex < 0 || newIndex >= fields.length) return

    const newFields = [...fields]
    const [removed] = newFields.splice(index, 1)
    newFields.splice(newIndex, 0, removed)
    setFields(newFields)
    reorderFieldsMutation.mutate(newFields)
  }

  const isPending = addFieldMutation.isPending || renameFieldMutation.isPending ||
    removeFieldMutation.isPending || reorderFieldsMutation.isPending

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div className="flex items-center justify-between p-4 border-b">
          <h2 className="text-lg font-semibold text-gray-900">
            Edit Fields: {noteType.name}
          </h2>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="p-4">
          {/* Error message */}
          {error && (
            <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
              {error}
            </div>
          )}

          {/* Field list */}
          <div className="space-y-2 mb-4">
            {fields.map((field, index) => (
              <div
                key={field}
                className="flex items-center gap-2 p-2 bg-gray-50 rounded-md"
              >
                {editingField === field ? (
                  <input
                    type="text"
                    value={editFieldName}
                    onChange={(e) => setEditFieldName(e.target.value)}
                    onBlur={() => handleRenameField(field)}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') handleRenameField(field)
                      if (e.key === 'Escape') setEditingField(null)
                    }}
                    className="flex-1 px-2 py-1 border border-blue-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                    autoFocus
                  />
                ) : (
                  <span
                    className="flex-1 cursor-pointer hover:text-blue-600"
                    onClick={() => {
                      setEditingField(field)
                      setEditFieldName(field)
                    }}
                    title="Click to rename"
                  >
                    {field}
                  </span>
                )}

                {/* Move up/down buttons */}
                <button
                  onClick={() => moveField(index, 'up')}
                  disabled={index === 0 || isPending}
                  className="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
                  title="Move up"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                  </svg>
                </button>
                <button
                  onClick={() => moveField(index, 'down')}
                  disabled={index === fields.length - 1 || isPending}
                  className="p-1 text-gray-400 hover:text-gray-600 disabled:opacity-30 disabled:cursor-not-allowed"
                  title="Move down"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>

                {/* Delete button */}
                <button
                  onClick={() => handleRemoveField(field)}
                  disabled={fields.length <= 1 || isPending}
                  className="p-1 text-red-400 hover:text-red-600 disabled:opacity-30 disabled:cursor-not-allowed"
                  title="Remove field"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                </button>
              </div>
            ))}
          </div>

          {/* Add new field */}
          <form onSubmit={handleAddField} className="flex gap-2">
            <input
              id="new-field-name"
              type="text"
              value={newFieldName}
              onChange={(e) => setNewFieldName(e.target.value)}
              placeholder="New field name..."
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={isPending}
            />
            <button
              type="submit"
              disabled={!newFieldName.trim() || isPending}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Add
            </button>
          </form>

          <p className="mt-2 text-xs text-gray-500">
            Click a field name to rename it. Reserved names: Tags, Type, Deck, Card, FrontSide
          </p>
        </div>

        <div className="flex justify-end gap-2 p-4 border-t bg-gray-50">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
