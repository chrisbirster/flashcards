import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Note, NoteType } from '#/lib/api'
import { useAppRepository } from '#/lib/app-repository'

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
  const showingEditorOnMobile = Boolean(noteId)

  return (
    <div className="space-y-6">
      <section className="rounded-[2rem] border border-slate-200 bg-white p-5 shadow-sm md:p-6">
        <div className="flex flex-col gap-4 lg:flex-row lg:items-end">
          <div className="flex-1">
            <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Search notes</label>
            <input
              type="text"
              value={search}
              onChange={(event) => updateSearchParams({ q: event.target.value, cursor: undefined })}
              placeholder="Search note content, tags, or type"
              className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
            />
          </div>
          <div className="grid gap-4 sm:grid-cols-3">
            <div>
              <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Deck</label>
              <select
                value={deckId ?? ''}
                onChange={(event) => updateSearchParams({ deckId: event.target.value || undefined, cursor: undefined, noteId: undefined })}
                className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
              >
                <option value="">All decks</option>
                {(decksQuery.data ?? []).map((deck) => (
                  <option key={deck.id} value={deck.id}>
                    {deck.name}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Note type</label>
              <select
                value={typeId}
                onChange={(event) => updateSearchParams({ typeId: event.target.value || undefined, cursor: undefined, noteId: undefined })}
                className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
              >
                <option value="">All note types</option>
                {(noteTypesQuery.data ?? []).map((noteType) => (
                  <option key={noteType.name} value={noteType.name}>
                    {noteType.name}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Tag</label>
              <input
                type="text"
                value={tag}
                onChange={(event) => updateSearchParams({ tag: event.target.value || undefined, cursor: undefined, noteId: undefined })}
                placeholder="Filter by tag"
                className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
              />
            </div>
          </div>
        </div>
      </section>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)]">
        <section className={`${showingEditorOnMobile ? 'hidden xl:block' : 'block'} rounded-[2rem] border border-slate-200 bg-white shadow-sm`}>
          <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <div>
              <h2 className="text-lg font-semibold text-slate-950">Notes</h2>
              <p className="text-sm text-slate-500">{notesQuery.data?.total ?? 0} matching notes</p>
            </div>
            <Link
              to={deckId ? `/notes/add?deckId=${deckId}` : '/notes/add'}
              className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 hover:border-slate-400 hover:bg-stone-50"
            >
              Add note
            </Link>
          </div>

          {notesQuery.isLoading && <p className="p-5 text-sm text-slate-500">Loading notes...</p>}
          {!notesQuery.isLoading && (notesQuery.data?.notes.length ?? 0) === 0 && (
            <p className="p-5 text-sm text-slate-500">No notes match the current filters.</p>
          )}
          {(notesQuery.data?.notes.length ?? 0) > 0 && (
            <ul className="divide-y divide-slate-200">
              {notesQuery.data?.notes.map((note) => (
                <li key={note.id} className="px-5 py-4">
                  <button
                    type="button"
                    onClick={() => updateSearchParams({ noteId: note.id, deckId: note.deckId ?? deckId })}
                    className={`w-full rounded-2xl border px-4 py-4 text-left transition-colors ${
                      note.id === noteId
                        ? 'border-slate-950 bg-slate-950 text-white'
                        : 'border-slate-200 bg-stone-50 hover:border-slate-300 hover:bg-white'
                    }`}
                  >
                    <div className="flex items-start justify-between gap-4">
                      <div className="min-w-0">
                        <div className="flex flex-wrap items-center gap-2">
                          <span className={`rounded-full px-2.5 py-1 text-xs font-medium ${note.id === noteId ? 'bg-white/15 text-white' : 'bg-white text-slate-700 ring-1 ring-slate-200'}`}>
                            {note.typeId}
                          </span>
                          {note.deckName && (
                            <span className={`text-xs uppercase tracking-[0.18em] ${note.id === noteId ? 'text-white/70' : 'text-slate-400'}`}>
                              {note.deckName}
                            </span>
                          )}
                        </div>
                        <p className="mt-3 truncate text-sm font-medium">{note.fieldPreview || 'Untitled note'}</p>
                        <p className={`mt-2 text-xs ${note.id === noteId ? 'text-white/70' : 'text-slate-500'}`}>
                          {note.cardCount} card{note.cardCount === 1 ? '' : 's'} • Updated {formatTimestamp(note.modifiedAt)}
                        </p>
                        {note.tags.length > 0 && (
                          <div className="mt-3 flex flex-wrap gap-1">
                            {note.tags.slice(0, 4).map((entry) => (
                              <span
                                key={entry}
                                className={`rounded-full px-2 py-1 text-xs ${note.id === noteId ? 'bg-white/15 text-white' : 'bg-white text-slate-500 ring-1 ring-slate-200'}`}
                              >
                                #{entry}
                              </span>
                            ))}
                          </div>
                        )}
                      </div>
                      <div className="flex shrink-0 items-center gap-2">
                        <span className={`text-xs ${note.id === noteId ? 'text-white/70' : 'text-slate-400'}`}>
                          {formatTimestamp(note.createdAt)}
                        </span>
                        <button
                          type="button"
                          onClick={(event) => {
                            event.stopPropagation()
                            handleDelete(note.id)
                          }}
                          className={`rounded-xl px-3 py-2 text-xs font-medium ${note.id === noteId ? 'bg-white text-slate-950' : 'border border-slate-300 text-slate-600 hover:border-red-300 hover:text-red-600'}`}
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
            <div className="flex items-center justify-between border-t border-slate-200 px-5 py-4">
              <button
                type="button"
                onClick={() => updateSearchParams({ cursor: notesQuery.data?.prevCursor })}
                disabled={!notesQuery.data?.prevCursor}
                className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-40"
              >
                Previous
              </button>
              <button
                type="button"
                onClick={() => updateSearchParams({ cursor: notesQuery.data?.nextCursor })}
                disabled={!notesQuery.data?.nextCursor}
                className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 disabled:cursor-not-allowed disabled:opacity-40"
              >
                Next
              </button>
            </div>
          )}
        </section>

        <section className={`${showingEditorOnMobile ? 'block' : 'hidden xl:block'} rounded-[2rem] border border-slate-200 bg-white shadow-sm`}>
          <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
            <div>
              <h2 className="text-lg font-semibold text-slate-950">Editor</h2>
              <p className="text-sm text-slate-500">
                {selectedNote ? `Editing note #${selectedNote.id}` : 'Select a note to start editing'}
              </p>
            </div>
            {showingEditorOnMobile && (
              <button
                type="button"
                onClick={() => updateSearchParams({ noteId: undefined })}
                className="rounded-xl border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 xl:hidden"
              >
                Back to list
              </button>
            )}
          </div>

          {!noteId && (
            <div className="p-8 text-sm text-slate-500">
              Pick a note from the list to update fields, move it to another deck, or delete it.
            </div>
          )}

          {noteId && selectedNoteQuery.isLoading && (
            <div className="p-8 text-sm text-slate-500">Loading note details...</div>
          )}

          {noteId && selectedNote && (
            <form
              onSubmit={(event) => {
                event.preventDefault()
                updateNoteMutation.mutate({ id: selectedNote.id, note: selectedNote })
              }}
              className="space-y-5 p-5"
            >
              {editorMessage && (
                <div className="rounded-2xl border border-emerald-200 bg-emerald-50 px-4 py-3 text-sm text-emerald-800">
                  {editorMessage}
                </div>
              )}
              {editorError && (
                <div className="rounded-2xl border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
                  {editorError}
                </div>
              )}

              <div className="grid gap-4 md:grid-cols-2">
                <div>
                  <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Note type</label>
                  <select
                    value={draftTypeId}
                    onChange={(event) => {
                      const nextTypeId = event.target.value
                      setDraftTypeId(nextTypeId)
                      setDraftFieldVals(buildFieldValues(noteTypeMap.get(nextTypeId), draftFieldVals))
                    }}
                    className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
                  >
                    {(noteTypesQuery.data ?? []).map((noteType) => (
                      <option key={noteType.name} value={noteType.name}>
                        {noteType.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="text-xs font-medium uppercase tracking-[0.18em] text-slate-400">Deck</label>
                  <select
                    value={draftDeckId}
                    onChange={(event) => setDraftDeckId(Number(event.target.value))}
                    className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
                  >
                    {(decksQuery.data ?? []).map((deck) => (
                      <option key={deck.id} value={deck.id}>
                        {deck.name}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="space-y-4">
                {editorFields.map((field) => (
                  <div key={field}>
                    <label className="text-sm font-medium text-slate-700">{field}</label>
                    <textarea
                      value={draftFieldVals[field] || ''}
                      onChange={(event) =>
                        setDraftFieldVals((current) => ({
                          ...current,
                          [field]: event.target.value,
                        }))
                      }
                      rows={4}
                      className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 font-mono text-sm focus:border-slate-500 focus:outline-none"
                    />
                  </div>
                ))}
              </div>

              <div>
                <label className="text-sm font-medium text-slate-700">Tags</label>
                <input
                  type="text"
                  value={draftTags}
                  onChange={(event) => setDraftTags(event.target.value)}
                  placeholder="comma or space separated tags"
                  className="mt-2 w-full rounded-2xl border border-slate-300 px-4 py-3 text-sm focus:border-slate-500 focus:outline-none"
                />
              </div>

              <div className="flex flex-wrap items-center justify-between gap-3 border-t border-slate-200 pt-4">
                <p className="text-sm text-slate-500">
                  Created {formatTimestamp(selectedNote.createdAt)} • Updated {formatTimestamp(selectedNote.modifiedAt)}
                </p>
                <div className="flex flex-wrap gap-3">
                  <button
                    type="button"
                    onClick={() => handleDelete(selectedNote.id)}
                    className="rounded-2xl border border-red-200 px-4 py-2.5 text-sm font-medium text-red-700 hover:bg-red-50"
                  >
                    Delete note
                  </button>
                  <button
                    type="submit"
                    disabled={updateNoteMutation.isPending || draftDeckId === 0 || !draftTypeId}
                    className="rounded-2xl bg-slate-950 px-4 py-2.5 text-sm font-medium text-white hover:bg-slate-800 disabled:cursor-not-allowed disabled:bg-slate-300"
                  >
                    {updateNoteMutation.isPending ? 'Saving...' : 'Save changes'}
                  </button>
                </div>
              </div>
            </form>
          )}
        </section>
      </div>
    </div>
  )
}
