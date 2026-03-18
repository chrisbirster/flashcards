import { useQuery } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'
import type { RecentDeckNoteSummary } from '#/lib/api'

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

function RecentDeckNoteItem({ note }: {note: RecentDeckNoteSummary}) {
  return (
    <li className="border border-gray-200 rounded-lg p-3">
      <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2">
        <div className="min-w-0">
          <div className="flex items-center gap-2 flex-wrap">
            <span className="inline-flex items-center rounded-full bg-blue-50 px-2 py-0.5 text-xs font-medium text-blue-700">
              {note.noteType}
            </span>
            {note.cardCountInDeck > 1 && (
              <span className="text-xs text-gray-500">
                {note.cardCountInDeck} cards
              </span>
            )}
          </div>
          <p className="mt-2 text-sm text-gray-900 break-words">
            {note.fieldPreview || 'No preview available'}
          </p>
          {note.tags.length > 0 && (
            <div className="mt-2 flex flex-wrap gap-1">
              {note.tags.slice(0, 4).map((tag) => (
                <span
                  key={tag}
                  className="inline-flex items-center rounded-full bg-gray-100 px-2 py-0.5 text-xs text-gray-600"
                >
                  #{tag}
                </span>
              ))}
            </div>
          )}
        </div>
        <div className="text-xs text-gray-500 shrink-0">
          {formatTimestamp(note.createdAt)}
        </div>
      </div>
    </li>
  )
}

export function RecentDeckNotesPanel({ deckId }: {deckId?: number }) {
  const repository = useAppRepository()
  const enabled = typeof deckId === 'number' && deckId > 0
  const { data, isLoading, isError } = useQuery({
    queryKey: ['deck-notes', deckId],
    queryFn: () => repository.fetchDeckNotes(deckId!, 20),
    enabled,
  })

  if (!enabled) {
    return null
  }

  return (
    <section className="bg-white rounded-lg shadow p-4 sm:p-6" data-testid="recent-deck-notes">
      <div className="flex items-center justify-between gap-2 mb-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Recent Notes in This Deck</h2>
          <p className="text-sm text-gray-500">
            Newest first so you can track deck-building progress while adding notes.
          </p>
        </div>
        {data?.notes?.length ? (
          <span className="text-xs text-gray-500">{data.notes.length} shown</span>
        ) : null}
      </div>

      {isLoading && <p className="text-sm text-gray-500">Loading recent notes...</p>}
      {isError && <p className="text-sm text-red-600">Failed to load recent notes.</p>}
      {!isLoading && !isError && data?.notes.length === 0 && (
        <p className="text-sm text-gray-500">No notes in this deck yet.</p>
      )}
      {!isLoading && !isError && data?.notes.length ? (
        <ul className="space-y-3">
          {data.notes.map((note) => (
            <RecentDeckNoteItem key={note.noteId} note={note} />
          ))}
        </ul>
      ) : null}
    </section>
  )
}
