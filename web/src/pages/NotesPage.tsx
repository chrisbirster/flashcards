import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Note, NoteType } from '#/lib/api'
import { useAppRepository } from '#/lib/app-repository'
import { ActionBar, FormActions } from '#/components/action-bar'
import { FieldRow } from '#/components/field-row'
import { EmptyState, PageContainer, PageSection } from '#/components/page-layout'
import { FullscreenSheet, Sheet } from '#/components/sheet'

function formatTimestamp(value: string) {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return new Intl.DateTimeFormat(undefined, {
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  }).format(date)
}

function tagsToInput(tags: string[]) {
  return tags.join(', ')
}

function parseTags(value: string) {
  return value
    .split(/[,\s]+/)
    .map((tag) => tag.trim())
    .filter(Boolean)
}

function buildFieldValues(noteType: NoteType | undefined, source: Record<string, string>) {
  if (!noteType) return source
  const shaped: Record<string, string> = {}
  noteType.fields.forEach((field) => {
    shaped[field] = source[field] || ''
  })
  return shaped
}

export function NotesPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const [searchParams, setSearchParams] = useSearchParams()
  const [mobileFiltersOpen, setMobileFiltersOpen] = useState(false)

  const deckId = Number(searchParams.get('deckId') || '') || undefined
  const noteId = Number(searchParams.get('noteId') || '') || undefined
  const search = searchParams.get('q') || ''
  const typeId = searchParams.get('typeId') || ''
  const tag = searchParams.get('tag') || ''
  const cursor = searchParams.get('cursor') || undefined

  const noteTypesQuery = useQuery({
    queryKey: ['note-types'],
    queryFn: () => repository.fetchNoteTypes(),
  })
  const decksQuery = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })
  const notesQuery = useQuery({
    queryKey: ['notes', { deckId, search, typeId, tag, cursor }],
    queryFn: () =>
      repository.fetchNotes({
        deckId,
        q: search || undefined,
        typeId: typeId || undefined,
        tag: tag || undefined,
        cursor,
        limit: 25,
      }),
  })
  const selectedNoteQuery = useQuery({
    queryKey: ['note', noteId],
    queryFn: () => repository.fetchNote(noteId!),
    enabled: Boolean(noteId),
  })

  const noteTypeMap = useMemo(() => {
    const map = new Map<string, NoteType>()
    for (const noteType of noteTypesQuery.data ?? []) {
      map.set(noteType.name, noteType)
    }
    return map
  }, [noteTypesQuery.data])

  const selectedNote = selectedNoteQuery.data
  const [draftTypeId, setDraftTypeId] = useState('')
  const [draftDeckId, setDraftDeckId] = useState<number>(deckId || 0)
  const [draftFieldVals, setDraftFieldVals] = useState<Record<string, string>>({})
  const [draftTags, setDraftTags] = useState('')
  const [editorMessage, setEditorMessage] = useState<string | null>(null)
  const [editorError, setEditorError] = useState<string | null>(null)

  useEffect(() => {
    if (!selectedNote) return
    setDraftTypeId(selectedNote.typeId)
    setDraftDeckId(selectedNote.deckId || deckId || 0)
    setDraftFieldVals(buildFieldValues(noteTypeMap.get(selectedNote.typeId), selectedNote.fieldVals))
    setDraftTags(tagsToInput(selectedNote.tags))
    setEditorMessage(null)
    setEditorError(null)
  }, [selectedNote, noteTypeMap, deckId])

  const updateSearchParams = (updates: Record<string, string | number | undefined | null>) => {
    const next = new URLSearchParams(searchParams)
    Object.entries(updates).forEach(([key, value]) => {
      if (value === undefined || value === null || value === '') {
        next.delete(key)
      } else {
        next.set(key, String(value))
      }
    })
    setSearchParams(next)
  }

  const updateNoteMutation = useMutation({
    mutationFn: (payload: {id: number; note: Note}) =>
      repository.updateNote(payload.id, {
        typeId: draftTypeId,
        deckId: draftDeckId,
        fieldVals: draftFieldVals,
        tags: parseTags(draftTags),
      }),
    onSuccess: () => {
      setEditorMessage('Note updated.')
      setEditorError(null)
      queryClient.invalidateQueries({ queryKey: ['notes'] })
      queryClient.invalidateQueries({ queryKey: ['note', noteId] })
      queryClient.invalidateQueries({ queryKey: ['deck-notes'] })
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats'] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
    },
    onError: (error: Error) => {
      setEditorError(error.message)
      setEditorMessage(null)
    },
  })

  const deleteNoteMutation = useMutation({
    mutationFn: (targetNoteId: number) => repository.deleteNote(targetNoteId),
    onSuccess: (_, deletedNoteId) => {
      const nextParams = new URLSearchParams(searchParams)
      if (noteId === deletedNoteId) {
        nextParams.delete('noteId')
        setSearchParams(nextParams)
      }
      queryClient.invalidateQueries({ queryKey: ['notes'] })
      queryClient.invalidateQueries({ queryKey: ['note'] })
      queryClient.invalidateQueries({ queryKey: ['deck-notes'] })
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats'] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
      setEditorMessage('Note deleted.')
      setEditorError(null)
    },
    onError: (error: Error) => {
      setEditorError(error.message)
      setEditorMessage(null)
    },
  })

  const handleDelete = (targetNoteId: number) => {
    if (!window.confirm('Delete this note and all cards generated from it?')) {
      return
    }
    deleteNoteMutation.mutate(targetNoteId)
  }

  const currentNoteType = noteTypeMap.get(draftTypeId)
  const editorFields = currentNoteType?.fields ?? Object.keys(draftFieldVals)
  const parsedDraftTags = parseTags(draftTags)

  const renderFilters = () => (
    <div className="grid gap-4 lg:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)_minmax(0,0.8fr)_minmax(0,0.8fr)]">
      <FieldRow label="Search notes">
        <input
          type="text"
          value={search}
          onChange={(event) => updateSearchParams({ q: event.target.value, cursor: undefined })}
          placeholder="Search note content, tags, or type"
          className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
        />
      </FieldRow>
      <FieldRow label="Deck">
        <select
          value={deckId ?? ''}
          onChange={(event) => updateSearchParams({ deckId: event.target.value || undefined, cursor: undefined, noteId: undefined })}
          className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
        >
          <option value="">All decks</option>
          {(decksQuery.data ?? []).map((deck) => (
            <option key={deck.id} value={deck.id}>
              {deck.name}
            </option>
          ))}
        </select>
      </FieldRow>
      <FieldRow label="Note type">
        <select
          value={typeId}
          onChange={(event) => updateSearchParams({ typeId: event.target.value || undefined, cursor: undefined, noteId: undefined })}
          className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
        >
          <option value="">All note types</option>
          {(noteTypesQuery.data ?? []).map((noteType) => (
            <option key={noteType.name} value={noteType.name}>
              {noteType.name}
            </option>
          ))}
        </select>
      </FieldRow>
      <FieldRow label="Tag">
        <input
          type="text"
          value={tag}
          onChange={(event) => updateSearchParams({ tag: event.target.value || undefined, cursor: undefined, noteId: undefined })}
          placeholder="Filter by tag"
          className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
        />
      </FieldRow>
    </div>
  )

  const renderEditorForm = (mode: 'desktop' | 'mobile') => {
    if (!noteId) {
      return (
        <div className="p-6">
          <EmptyState
            title="Select a note to edit"
            description="Pick a note from the list to update fields, move it to another deck, or delete it."
          />
        </div>
      )
    }

    if (selectedNoteQuery.isLoading) {
      return <div className="p-6 text-sm text-[var(--app-text-soft)]">Loading note details...</div>
    }

    if (!selectedNote) {
      return <div className="p-6 text-sm text-[var(--app-danger-text)]">Failed to load note details.</div>
    }

    return (
      <form
        onSubmit={(event) => {
          event.preventDefault()
          updateNoteMutation.mutate({ id: selectedNote.id, note: selectedNote })
        }}
        className="flex h-full min-h-0 flex-col"
      >
        <div className="space-y-5 p-5">
          {editorMessage ? (
            <div className="rounded-2xl border border-[var(--app-success-line)] bg-[var(--app-success-surface)] px-4 py-3 text-sm text-[var(--app-success-text)]">
              {editorMessage}
            </div>
          ) : null}
          {editorError ? (
            <div className="rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 py-3 text-sm text-[var(--app-danger-text)]">
              {editorError}
            </div>
          ) : null}

          <div className="grid gap-4 md:grid-cols-2">
            <FieldRow label="Note type">
              <select
                value={draftTypeId}
                onChange={(event) => {
                  const nextTypeId = event.target.value
                  setDraftTypeId(nextTypeId)
                  setDraftFieldVals(buildFieldValues(noteTypeMap.get(nextTypeId), draftFieldVals))
                }}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              >
                {(noteTypesQuery.data ?? []).map((noteType) => (
                  <option key={noteType.name} value={noteType.name}>
                    {noteType.name}
                  </option>
                ))}
              </select>
            </FieldRow>

            <FieldRow label="Deck">
              <select
                value={draftDeckId}
                onChange={(event) => setDraftDeckId(Number(event.target.value))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              >
                {(decksQuery.data ?? []).map((deck) => (
                  <option key={deck.id} value={deck.id}>
                    {deck.name}
                  </option>
                ))}
              </select>
            </FieldRow>
          </div>

          <div className="space-y-4">
            {editorFields.map((field) => (
              <FieldRow key={field} label={field}>
                <textarea
                  value={draftFieldVals[field] || ''}
                  onChange={(event) =>
                    setDraftFieldVals((current) => ({
                      ...current,
                      [field]: event.target.value,
                    }))
                  }
                  rows={mode === 'mobile' ? 5 : 4}
                  className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 font-mono text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                />
              </FieldRow>
            ))}
          </div>

          <FieldRow label="Tags" hint="Comma or space separated">
            <div className="space-y-3">
              {parsedDraftTags.length > 0 ? (
                <div className="flex flex-wrap gap-2">
                  {parsedDraftTags.map((entry) => (
                    <span
                      key={entry}
                      className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]"
                    >
                      #{entry}
                    </span>
                  ))}
                </div>
              ) : null}
              <input
                type="text"
                value={draftTags}
                onChange={(event) => setDraftTags(event.target.value)}
                placeholder="comma or space separated tags"
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              />
            </div>
          </FieldRow>

          {mode === 'desktop' ? (
            <div className="flex flex-wrap items-center justify-between gap-3 border-t border-[var(--app-line)] pt-4">
              <p className="text-sm text-[var(--app-text-soft)]">
                Created {formatTimestamp(selectedNote.createdAt)} • Updated {formatTimestamp(selectedNote.modifiedAt)}
              </p>
              <FormActions>
                <button
                  type="button"
                  onClick={() => handleDelete(selectedNote.id)}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 text-sm font-medium text-[var(--app-danger-text)]"
                >
                  Delete note
                </button>
                <button
                  type="submit"
                  disabled={updateNoteMutation.isPending || draftDeckId === 0 || !draftTypeId}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {updateNoteMutation.isPending ? 'Saving...' : 'Save changes'}
                </button>
              </FormActions>
            </div>
          ) : null}
        </div>

        {mode === 'mobile' ? (
          <ActionBar>
            <div className="space-y-2">
              <p className="text-xs text-[var(--app-muted)]">
                Created {formatTimestamp(selectedNote.createdAt)} • Updated {formatTimestamp(selectedNote.modifiedAt)}
              </p>
              <FormActions>
                <button
                  type="button"
                  onClick={() => handleDelete(selectedNote.id)}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 text-sm font-medium text-[var(--app-danger-text)]"
                >
                  Delete note
                </button>
                <button
                  type="submit"
                  disabled={updateNoteMutation.isPending || draftDeckId === 0 || !draftTypeId}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {updateNoteMutation.isPending ? 'Saving...' : 'Save changes'}
                </button>
              </FormActions>
            </div>
          </ActionBar>
        ) : null}
      </form>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <PageSection className="hidden p-5 md:p-6 xl:block">
        {renderFilters()}
      </PageSection>

      <PageSection className="p-4 xl:hidden">
        <div className="flex items-center justify-between gap-3">
          <div>
            <p className="text-[11px] uppercase tracking-[0.22em] text-[var(--app-muted)]">Filters</p>
            <p className="text-sm text-[var(--app-text-soft)]">
              {notesQuery.data?.total ?? 0} matching note{(notesQuery.data?.total ?? 0) === 1 ? '' : 's'}
            </p>
          </div>
          <button
            type="button"
            onClick={() => setMobileFiltersOpen(true)}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Filters
          </button>
        </div>
      </PageSection>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)]">
        <PageSection className="overflow-hidden">
          <div className="flex items-center justify-between border-b border-[var(--app-line)] px-5 py-4">
            <div>
              <h2 className="text-lg font-semibold text-[var(--app-text)]">Notes</h2>
              <p className="text-sm text-[var(--app-text-soft)]">{notesQuery.data?.total ?? 0} matching notes</p>
            </div>
            <Link
              to={deckId ? `/notes/add?deckId=${deckId}` : '/notes/add'}
              className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
            >
              Add note
            </Link>
          </div>

          {notesQuery.isLoading ? <p className="p-5 text-sm text-[var(--app-text-soft)]">Loading notes...</p> : null}
          {!notesQuery.isLoading && (notesQuery.data?.notes.length ?? 0) === 0 && (
            <div className="p-5">
              <EmptyState
                title="No notes match the current filters"
                description="Try a wider search, switch decks, or open the Add Note flow to create a new note."
              />
            </div>
          )}
          {(notesQuery.data?.notes.length ?? 0) > 0 && (
            <ul className="space-y-3 p-4">
              {notesQuery.data?.notes.map((note) => (
                <li key={note.id}>
                  <button
                    type="button"
                    onClick={() => updateSearchParams({ noteId: note.id, deckId: note.deckId ?? deckId })}
                    className={`w-full rounded-[1.5rem] border px-4 py-4 text-left transition-colors ${
                      note.id === noteId
                        ? 'border-[var(--app-accent)] bg-[var(--app-card-strong)] text-[var(--app-text)]'
                        : 'border-[var(--app-line)] bg-[var(--app-muted-surface)] hover:border-[var(--app-line-strong)] hover:bg-[var(--app-card)]'
                    }`}
                  >
                    <div className="flex items-start justify-between gap-4">
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${note.id === noteId ? 'bg-[var(--app-accent)] text-[var(--app-accent-ink)]' : 'bg-[var(--app-card)] text-[var(--app-text)] ring-1 ring-[var(--app-line)]'}`}>
                            {note.typeId}
                          </span>
                          {note.deckName && (
                            <span className={`text-xs uppercase tracking-[0.18em] ${note.id === noteId ? 'text-[var(--app-muted)]' : 'text-[var(--app-muted)]'}`}>
                              {note.deckName}
                            </span>
                          )}
                        </div>
                        <p className="mt-3 truncate text-sm font-medium text-[var(--app-text)]">{note.fieldPreview || 'Untitled note'}</p>
                        <p className="mt-2 text-xs text-[var(--app-text-soft)]">
                          {note.cardCount} card{note.cardCount === 1 ? '' : 's'} • Updated {formatTimestamp(note.modifiedAt)}
                        </p>
                        {note.tags.length > 0 && (
                          <div className="mt-3 flex flex-wrap gap-1">
                            {note.tags.slice(0, 4).map((entry) => (
                              <span
                                key={entry}
                                className={`rounded-full px-2 py-1 text-xs ${note.id === noteId ? 'bg-[var(--app-muted-surface)] text-[var(--app-text-soft)]' : 'bg-[var(--app-card)] text-[var(--app-text-soft)] ring-1 ring-[var(--app-line)]'}`}
                              >
                                #{entry}
                              </span>
                            ))}
                          </div>
                        )}
                      </div>
                      <div className="flex shrink-0 flex-col items-end gap-2">
                        <span className="text-xs text-[var(--app-muted)]">
                          {formatTimestamp(note.createdAt)}
                        </span>
                        <button
                          type="button"
                          onClick={(event) => {
                            event.stopPropagation()
                            handleDelete(note.id)
                          }}
                          className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 text-xs font-medium text-[var(--app-text-soft)] hover:border-[var(--app-danger-line)] hover:text-[var(--app-danger-text)]"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          )}

          {(notesQuery.data?.nextCursor || notesQuery.data?.prevCursor) && (
            <div className="flex items-center justify-between border-t border-[var(--app-line)] px-5 py-4">
              <button
                type="button"
                onClick={() => updateSearchParams({ cursor: notesQuery.data?.prevCursor })}
                disabled={!notesQuery.data?.prevCursor}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)] disabled:cursor-not-allowed disabled:opacity-40"
              >
                Previous
              </button>
              <button
                type="button"
                onClick={() => updateSearchParams({ cursor: notesQuery.data?.nextCursor })}
                disabled={!notesQuery.data?.nextCursor}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)] disabled:cursor-not-allowed disabled:opacity-40"
              >
                Next
              </button>
            </div>
          )}
        </PageSection>

        <PageSection className="hidden overflow-hidden xl:block">
          <div className="flex items-center justify-between border-b border-[var(--app-line)] px-5 py-4">
            <div>
              <h2 className="text-lg font-semibold text-[var(--app-text)]">Editor</h2>
              <p className="text-sm text-[var(--app-text-soft)]">
                {selectedNote ? `Editing note #${selectedNote.id}` : 'Select a note to start editing'}
              </p>
            </div>
          </div>
          {renderEditorForm('desktop')}
        </PageSection>
      </div>

      <Sheet open={mobileFiltersOpen} onClose={() => setMobileFiltersOpen(false)} title="Filters">
        <div className="space-y-4">
          {renderFilters()}
          <button
            type="button"
            onClick={() => setMobileFiltersOpen(false)}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
          >
            Apply filters
          </button>
        </div>
      </Sheet>

      <FullscreenSheet
        open={Boolean(noteId)}
        onClose={() => updateSearchParams({ noteId: undefined })}
        title={selectedNote ? `Note #${selectedNote.id}` : 'Note editor'}
      >
        {renderEditorForm('mobile')}
      </FullscreenSheet>
    </PageContainer>
  )
}
