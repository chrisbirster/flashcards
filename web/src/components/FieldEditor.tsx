import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { addField, renameField, removeField, reorderFields, setSortField, setFieldOptions } from '#/lib/api'
import type { NoteType, FieldOptions } from '#/lib/api'
import { ErrorMessage } from './message'
import { FieldEditorHeader } from './field-editor-header'
import { MoveDownIconButton, MoveUpIconButton } from './move-icon-button'
import { DeleteButton } from './delete-icon-button'
import { FieldOptionsIconButton } from './field-options-icon-button'
import { RtlOptionField } from './rtl-option-field'
import { FontSizeOptionField, FontTypeOptionField } from './font-option-field'
import { SortFieldInfoChip, RTLFieldInfoChip } from './field-info-chips'
import { EditFieldPanelHeader } from './edit-field-panel-header'
import { EditFieldPopup } from './edit-field-popup'

interface FieldEditorProps {
  noteType: NoteType
  onClose: () => void
}

export function FieldEditor({ noteType, onClose }: FieldEditorProps) {
  const queryClient = useQueryClient()
  const [fields, setFields] = useState<string[]>(noteType.fields)
  const [sortFieldIndex, setSortFieldIndexState] = useState<number>(noteType.sortFieldIndex || 0)
  const [fieldOptions, setFieldOptionsState] = useState<Record<string, FieldOptions>>(noteType.fieldOptions || {})
  const [newFieldName, setNewFieldName] = useState('')
  const [editingField, setEditingField] = useState<string | null>(null)
  const [editFieldName, setEditFieldName] = useState('')
  const [editingOptionsField, setEditingOptionsField] = useState<string | null>(null)
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

  const setSortFieldMutation = useMutation({
    mutationFn: (fieldIndex: number) => setSortField(noteType.name, { fieldIndex }),
    onSuccess: (data) => {
      setSortFieldIndexState(data.sortFieldIndex)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const setFieldOptionsMutation = useMutation({
    mutationFn: ({ fieldName, options }: { fieldName: string; options: FieldOptions }) =>
      setFieldOptions(noteType.name, fieldName, options),
    onSuccess: (data) => {
      setFieldOptionsState(data.fieldOptions)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const handleSetSortField = (index: number) => {
    setSortFieldMutation.mutate(index)
  }

  const handleUpdateFieldOptions = (fieldName: string, options: FieldOptions) => {
    setFieldOptionsMutation.mutate({ fieldName, options })
  }

  const isPending = addFieldMutation.isPending || renameFieldMutation.isPending ||
    removeFieldMutation.isPending || reorderFieldsMutation.isPending || setSortFieldMutation.isPending ||
    setFieldOptionsMutation.isPending

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <EditFieldPanelHeader noteTypeName={noteType.name} onClose={onClose} />
        <div className="p-4">
          {/* Error message */}
          {error && <ErrorMessage message={error} />}

          {/* Field list */}
          <div className="mb-4">
            <FieldEditorHeader sortField={fields[sortFieldIndex]} />
            <div className="space-y-2">
              {fields.map((field, index) => (
                <div key={field} className="space-y-2">
                  <div className="flex items-center gap-2 p-2 bg-gray-50 rounded-md">
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
                        <SortFieldInfoChip
                          booleanIndicator={sortFieldIndex === index}
                          datatestid={`sort-chip-${field}`}
                        />
                        <RTLFieldInfoChip
                          booleanIndicator={fieldOptions[field]?.rtl ?? false}
                          datatestid={`rtl-chip-${field}`}
                        />
                      </span>
                    )}

                    {/* Set as sort field */}
                    {sortFieldIndex !== index && (
                      <button
                        onClick={() => handleSetSortField(index)}
                        className="text-xs text-gray-500 hover:text-blue-600 px-2 py-1"
                        title="Set as sort field"
                        disabled={isPending}
                      >
                        Set Sort
                      </button>
                    )}

                    {/* Field options button */}

                    <FieldOptionsIconButton
                      handleClick={() => setEditingOptionsField(editingOptionsField === field ? null : field)}
                      isPending={false}
                      isEditing={editingOptionsField === field}
                      data-testid={`field-options-${field}`}
                    />

                    {/* Move up/down buttons */}
                    <MoveUpIconButton
                      handleClick={() => moveField(index, 'up')}
                      disabled={index === 0 || isPending}
                    />
                    <MoveDownIconButton
                      handleClick={() => moveField(index, 'down')}
                      disabled={index === fields.length - 1 || isPending}
                    />

                    {/* Delete button */}
                    <DeleteButton
                      onDelete={() => handleRemoveField(field)}
                      disabled={fields.length <= 1 || isPending}
                    />
                  </div>

                  {/* Field options panel */}
                  {editingOptionsField === field && (
                    <EditFieldPopup field={field}>
                      {/* Font */}
                      <FontTypeOptionField
                        fieldValue={fieldOptions[field]?.fontSize}
                        handleChange={(e) => handleUpdateFieldOptions(field, {
                          ...fieldOptions[field],
                          fontSize: parseInt(e.target.value) || undefined,
                        })}
                        isPending={isPending}
                        datatestid={`field-type-${field}`}
                      />
                      {/* Font Size */}
                      <FontSizeOptionField
                        fieldValue={fieldOptions[field]?.fontSize}
                        handleChange={(e) => handleUpdateFieldOptions(field, {
                          ...fieldOptions[field],
                          fontSize: parseInt(e.target.value) || undefined,
                        })}
                        isPending={isPending}
                        datatestid={`field-size-${field}`}
                      />
                      {/* RTL */}
                      <RtlOptionField
                        isChecked={fieldOptions[field]?.rtl || false}
                        handleChange={() => handleUpdateFieldOptions(field, {
                          ...fieldOptions[field],
                          rtl: !fieldOptions[field]?.rtl,
                        })}
                        datatestid={`field-rtl-${field}`}
                        isPending={isPending}
                      />
                    </ EditFieldPopup>
                  )}
                </div>
              ))}
            </div>
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
