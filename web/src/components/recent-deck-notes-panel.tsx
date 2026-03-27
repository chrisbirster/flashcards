import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { useAppRepository } from '#/lib/app-repository'
import type { RecentDeckNoteSummary } from '#/lib/api'
import { ConfirmSheet } from '#/components/sheet'

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

function EditIcon() {
  return (
    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M16.862 5.487a2.25 2.25 0 113.182 3.182L8.25 20.463 3 21l.537-5.25 13.325-10.263z" />
    </svg>
  )
}

function TrashIcon() {
  return (
    <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M4.5 7.5h15m-10.5 0V6a1.5 1.5 0 011.5-1.5h3A1.5 1.5 0 0115 6v1.5m-7.5 0V18a1.5 1.5 0 001.5 1.5h6A1.5 1.5 0 0016.5 18V7.5" />
    </svg>
  )
}

function RecentDeckNoteItem({
  note,
  onEdit,
  onDelete,
  deleting,
}: {
  note: RecentDeckNoteSummary
  onEdit: () => void
  onDelete: () => void
  deleting: boolean
}) {
  return (
    <li className="group relative rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4">
      <div className="flex flex-col gap-2 sm:flex-row sm:items-start sm:justify-between">
        <div className="min-w-0 pr-0 sm:pr-28">
          <div className="flex flex-wrap items-center gap-2">
            <span className="inline-flex items-center rounded-full bg-[var(--app-accent)] px-2.5 py-1 text-xs font-semibold text-[var(--app-accent-ink)]">
              {note.noteType}
            </span>
            {note.cardCountInDeck > 1 && (
              <span className="text-xs text-[var(--app-muted)]">{note.cardCountInDeck} cards</span>
            )}
          </div>
          <p className="mt-3 break-words text-sm leading-6 text-[var(--app-text)]">
            {note.fieldPreview || 'No preview available'}
          </p>
          {note.tags.length > 0 && (
            <div className="mt-3 flex flex-wrap gap-2">
              {note.tags.slice(0, 4).map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-2.5 py-1 text-xs text-[var(--app-text-soft)]"
                >
                  #{tag}
                </span>
              ))}
            </div>
          )}
        </div>

        <div className="shrink-0 text-xs text-[var(--app-muted)]">{formatTimestamp(note.createdAt)}</div>
      </div>

      <div className="mt-3 flex items-center justify-end gap-2 sm:absolute sm:bottom-3 sm:right-3 sm:mt-0 sm:opacity-0 sm:transition-opacity sm:group-hover:opacity-100 sm:group-focus-within:opacity-100">
        <button
          type="button"
          onClick={onEdit}
          className="inline-flex min-h-11 items-center gap-1 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-3 text-xs font-medium text-[var(--app-text)] hover:border-[var(--app-accent)]"
          aria-label={`Edit note ${note.noteId}`}
        >
          <EditIcon />
          Edit
        </button>
        <button
          type="button"
          onClick={onDelete}
          disabled={deleting}
          className="inline-flex min-h-11 items-center gap-1 rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-3 text-xs font-medium text-[var(--app-danger-text)] disabled:cursor-not-allowed disabled:opacity-50"
          aria-label={`Delete note ${note.noteId}`}
        >
          <TrashIcon />
          {deleting ? 'Deleting...' : 'Delete'}
        </button>
      </div>
    </li>
  )
}

export function RecentDeckNotesPanel({ deckId }: {deckId?: number }) {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const [pendingDeleteNoteId, setPendingDeleteNoteId] = useState<number | null>(null)
  const enabled = typeof deckId === 'number' && deckId > 0
  const { data, isLoading, isError } = useQuery({
    queryKey: ['deck-notes', deckId],
    queryFn: () => repository.fetchDeckNotes(deckId!, 10),
    enabled,
  })
  const notes = Array.isArray(data?.notes) ? data.notes : []

  const deleteMutation = useMutation({
    mutationFn: (noteId: number) => repository.deleteNote(noteId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['deck-notes'] })
      queryClient.invalidateQueries({ queryKey: ['notes'] })
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats'] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
    },
  })

  if (!enabled) {
    return null
  }

  return (
    <section className="rounded-[1.75rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4 shadow-sm sm:p-6" data-testid="recent-deck-notes">
      <div className="mb-4 flex items-center justify-between gap-2">
        <div>
          <h2 className="text-lg font-semibold text-[var(--app-text)]">Recent Notes in This Deck</h2>
          <p className="text-sm text-[var(--app-text-soft)]">
            Newest first so you can track deck-building progress while adding notes.
          </p>
        </div>
        <button
          type="button"
          onClick={() => navigate(`/notes/view?deckId=${deckId}`)}
          className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)]"
        >
          View all
        </button>
      </div>

      {isLoading && <p className="text-sm text-[var(--app-text-soft)]">Loading recent notes...</p>}
      {isError && <p className="text-sm text-[var(--app-danger-text)]">Failed to load recent notes.</p>}
      {!isLoading && !isError && notes.length === 0 && (
        <p className="text-sm text-[var(--app-text-soft)]">No notes in this deck yet.</p>
      )}
      {!isLoading && !isError && notes.length > 0 ? (
        <ul className="space-y-3">
          {notes.map((note) => (
            <RecentDeckNoteItem
              key={note.noteId}
              note={note}
              deleting={deleteMutation.isPending && deleteMutation.variables === note.noteId}
              onEdit={() => navigate(`/notes/view?deckId=${deckId}&noteId=${note.noteId}`)}
              onDelete={() => setPendingDeleteNoteId(note.noteId)}
            />
          ))}
        </ul>
      ) : null}
      <ConfirmSheet
        open={pendingDeleteNoteId !== null}
        onClose={() => setPendingDeleteNoteId(null)}
        title="Delete note"
        description="Delete this note and every card generated from it."
        confirmLabel={deleteMutation.isPending ? 'Deleting...' : 'Delete note'}
        onConfirm={() => {
          if (pendingDeleteNoteId !== null) {
            deleteMutation.mutate(pendingDeleteNoteId)
          }
        }}
        destructive
      />
    </section>
  )
}
