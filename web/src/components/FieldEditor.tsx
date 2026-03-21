import { useState } from 'react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
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
import { useAppRepository } from '#/lib/app-repository'

interface FieldEditorProps {
  noteType: NoteType
  onClose: () => void
}

export function FieldEditor({ noteType, onClose }: FieldEditorProps) {
  const queryClient = useQueryClient()
  const repository = useAppRepository()
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
    mutationFn: (fieldName: string) => repository.addField(noteType.name, { fieldName }),
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
      repository.renameField(noteType.name, { oldName, newName }),
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
    mutationFn: (fieldName: string) => repository.removeField(noteType.name, { fieldName }),
    onSuccess: (data) => {
      setFields(data.fields)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const reorderFieldsMutation = useMutation({
    mutationFn: (newFields: string[]) => repository.reorderFields(noteType.name, { fields: newFields }),
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
    mutationFn: (fieldIndex: number) => repository.setSortField(noteType.name, { fieldIndex }),
    onSuccess: (data) => {
      setSortFieldIndexState(data.sortFieldIndex)
      setError(null)
      invalidateNoteTypes()
    },
    onError: (err: Error) => setError(err.message),
  })

  const setFieldOptionsMutation = useMutation({
    mutationFn: ({ fieldName, options }: { fieldName: string; options: FieldOptions }) =>
      repository.setFieldOptions(noteType.name, fieldName, options),
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
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/60 p-2 sm:items-center sm:p-0">
      <div className="flex h-[100dvh] max-h-[100dvh] w-full max-w-2xl flex-col rounded-t-[1.75rem] border border-[var(--app-line)] bg-[var(--app-panel)] shadow-xl sm:mx-4 sm:h-auto sm:max-h-[90vh] sm:rounded-[1.75rem]">
        <EditFieldPanelHeader noteTypeName={noteType.name} onClose={onClose} />
        <div className="p-3 sm:p-4 overflow-auto">
          {/* Error message */}
          {error && <ErrorMessage message={error} />}

          {/* Field list */}
          <div className="mb-4">
            <FieldEditorHeader sortField={fields[sortFieldIndex]} />
            <div className="space-y-2">
              {fields.map((field, index) => (
                <div key={field} className="space-y-2">
                  <div className="flex flex-wrap items-center gap-2 rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-3">
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
                        className="flex-1 rounded-xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 py-2 text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                        autoFocus
                      />
                    ) : (
                      <span
                        className="flex-1 min-w-[10rem] cursor-pointer text-[var(--app-text)] hover:text-[var(--app-accent)]"
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
                        className="rounded-xl px-2 py-1 text-xs text-[var(--app-text-soft)] hover:text-[var(--app-accent)]"
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
                        fieldValue={fieldOptions[field]?.font}
                        handleChange={(e) => handleUpdateFieldOptions(field, {
                          ...fieldOptions[field],
                          font: e.target.value || undefined,
                        })}
                        isPending={isPending}
                        datatestid={`field-font-${field}`}
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
                      {/* HTML editor default */}
                      <div>
                        <label className="mb-1 block text-xs text-[var(--app-text-soft)]">Editor</label>
                        <label className="flex items-center gap-2 cursor-pointer">
                          <input
                            type="checkbox"
                            checked={fieldOptions[field]?.htmlEditor || false}
                            onChange={() => handleUpdateFieldOptions(field, {
                              ...fieldOptions[field],
                              htmlEditor: !fieldOptions[field]?.htmlEditor,
                            })}
                            className="h-4 w-4 rounded border-[var(--app-line-strong)] text-[var(--app-accent)]"
                            disabled={isPending}
                            data-testid={`field-html-${field}`}
                          />
                          <span className="text-xs text-[var(--app-text)]">HTML by default</span>
                        </label>
                      </div>
                    </ EditFieldPopup>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* Add new field */}
          <form onSubmit={handleAddField} className="flex flex-col sm:flex-row gap-2">
            <input
              id="new-field-name"
              type="text"
              value={newFieldName}
              onChange={(e) => setNewFieldName(e.target.value)}
              placeholder="New field name..."
              className="flex-1 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 py-2 text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              disabled={isPending}
            />
            <button
              type="submit"
              disabled={!newFieldName.trim() || isPending}
              className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60 sm:w-auto"
              data-testid="field-editor-add-button"
            >
              Add
            </button>
          </form>

          <p className="mt-2 text-xs text-[var(--app-text-soft)]">
            Click a field name to rename it. Reserved names: Tags, Type, Deck, Card, FrontSide
          </p>
        </div>

        <div className="flex justify-end gap-2 border-t border-[var(--app-line)] bg-[color:var(--app-header)]/95 p-3 backdrop-blur sm:p-4">
          <button
            onClick={onClose}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)] sm:w-auto"
            data-testid="close-field-editor-footer"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}
