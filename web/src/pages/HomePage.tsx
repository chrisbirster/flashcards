import { Link } from 'react-router'
import { useQueries, useQuery } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'

function StatCard({
  label,
  value,
  detail,
}: {
  label: string
  value: string | number
  detail: string
}) {
  return (
    <div className="rounded-3xl border border-slate-200 bg-white p-5 shadow-sm">
      <p className="text-xs uppercase tracking-[0.24em] text-slate-400">{label}</p>
      <p className="mt-4 text-3xl font-semibold tracking-tight text-slate-950">{value}</p>
      <p className="mt-2 text-sm text-slate-500">{detail}</p>
    </div>
  )
}

export function HomePage() {
  const repository = useAppRepository()
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })
  const decksQuery = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })
  const recentNotesQuery = useQuery({
    queryKey: ['notes', 'home-recent'],
    queryFn: () => repository.fetchNotes({ limit: 5 }),
  })

  const deckStatsQueries = useQueries({
    queries: (decksQuery.data ?? []).map((deck) => ({
      queryKey: ['deck-stats', deck.id],
      queryFn: () => repository.fetchDeckStats(deck.id),
    })),
  })

  const dueToday = deckStatsQueries.reduce((sum, query) => sum + (query.data?.dueToday ?? 0), 0)
  const plan = entitlementsQuery.data?.plan?.toUpperCase() ?? 'FREE'
  const noteUsage = entitlementsQuery.data?.usage.notes ?? 0
  const noteLimit = entitlementsQuery.data?.limits.maxNotes ?? 0
  const cardUsage = entitlementsQuery.data?.usage.cardsTotal ?? 0
  const cardLimit = entitlementsQuery.data?.limits.maxCardsTotal ?? 0

  return (
    <div className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,0.9fr)]">
        <div className="rounded-[2rem] bg-slate-950 px-6 py-8 text-white shadow-sm md:px-8 md:py-10">
          <p className="text-xs uppercase tracking-[0.3em] text-amber-300">Workspace overview</p>
          <h2 className="mt-4 max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl">
            Keep deck building, review, and template tuning in one browser-first workspace.
          </h2>
          <p className="mt-4 max-w-2xl text-sm leading-7 text-slate-300 md:text-base">
            Use the dashboard to watch plan usage, jump into notes management, and keep new cards moving into study without losing structure.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              to="/notes/add"
              className="inline-flex items-center rounded-2xl bg-white px-4 py-2.5 text-sm font-medium text-slate-950 hover:bg-slate-100"
            >
              Add note
            </Link>
            <Link
              to="/notes/view"
              className="inline-flex items-center rounded-2xl border border-slate-700 px-4 py-2.5 text-sm font-medium text-white hover:border-slate-500"
            >
              Open notes
            </Link>
            <Link
              to="/templates"
              className="inline-flex items-center rounded-2xl border border-slate-700 px-4 py-2.5 text-sm font-medium text-white hover:border-slate-500"
            >
              Open templates
            </Link>
          </div>
        </div>

        <div className="rounded-[2rem] border border-slate-200 bg-white p-6 shadow-sm">
          <p className="text-xs uppercase tracking-[0.28em] text-slate-400">Plan usage</p>
          <div className="mt-5 space-y-5">
            <div>
              <div className="flex items-center justify-between text-sm text-slate-600">
                <span>Notes</span>
                <span>{noteUsage} / {noteLimit}</span>
              </div>
              <div className="mt-2 h-2 rounded-full bg-slate-100">
                <div
                  className="h-2 rounded-full bg-slate-950"
                  style={{ width: `${noteLimit > 0 ? Math.min((noteUsage / noteLimit) * 100, 100) : 0}%` }}
                />
              </div>
            </div>
            <div>
              <div className="flex items-center justify-between text-sm text-slate-600">
                <span>Total cards</span>
                <span>{cardUsage} / {cardLimit}</span>
              </div>
              <div className="mt-2 h-2 rounded-full bg-slate-100">
                <div
                  className="h-2 rounded-full bg-amber-400"
                  style={{ width: `${cardLimit > 0 ? Math.min((cardUsage / cardLimit) * 100, 100) : 0}%` }}
                />
              </div>
            </div>
            <div className="rounded-2xl bg-stone-50 p-4">
              <p className="text-xs uppercase tracking-[0.2em] text-slate-400">Current plan</p>
              <p className="mt-2 text-2xl font-semibold text-slate-950">{plan}</p>
              <p className="mt-2 text-sm text-slate-500">
                Study Groups creation unlocks on Team. Marketplace publishing arrives in a later tranche for eligible plans.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Total Decks"
          value={entitlementsQuery.data?.usage.decks ?? decksQuery.data?.length ?? 0}
          detail="Active decks in this workspace."
        />
        <StatCard
          label="Total Notes"
          value={noteUsage}
          detail="Structured notes powering all generated cards."
        />
        <StatCard
          label="Due Today"
          value={dueToday}
          detail="Cards currently scheduled for review."
        />
        <StatCard
          label="Plan"
          value={plan}
          detail="Workspace entitlements and limits."
        />
      </section>

      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
        <div className="rounded-[2rem] border border-slate-200 bg-white p-6 shadow-sm">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-slate-950">Recent note activity</h3>
              <p className="mt-1 text-sm text-slate-500">Jump back into the notes you created or edited most recently.</p>
            </div>
            <Link to="/notes/view" className="text-sm font-medium text-slate-950 hover:text-slate-700">
              View all
            </Link>
          </div>

          {recentNotesQuery.isLoading && <p className="mt-6 text-sm text-slate-500">Loading recent notes...</p>}
          {!recentNotesQuery.isLoading && (recentNotesQuery.data?.notes.length ?? 0) === 0 && (
            <p className="mt-6 text-sm text-slate-500">No notes yet. Start with your first note and the dashboard will reflect it here.</p>
          )}
          {(recentNotesQuery.data?.notes?.length ?? 0) > 0 && (
            <ul className="mt-6 space-y-3">
              {recentNotesQuery.data?.notes.map((note) => (
                <li key={note.id} className="rounded-2xl border border-slate-200 bg-stone-50 p-4">
                  <div className="flex items-start justify-between gap-4">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="rounded-full bg-white px-2.5 py-1 text-xs font-medium text-slate-700 ring-1 ring-slate-200">
                          {note.typeId}
                        </span>
                        {note.deckName && (
                          <span className="text-xs uppercase tracking-[0.16em] text-slate-400">{note.deckName}</span>
                        )}
                      </div>
                      <p className="mt-2 truncate text-sm font-medium text-slate-900">{note.fieldPreview || 'Untitled note'}</p>
                      {note.tags.length > 0 && (
                        <div className="mt-2 flex flex-wrap gap-1">
                          {note.tags.slice(0, 3).map((tag) => (
                            <span key={tag} className="rounded-full bg-white px-2 py-1 text-xs text-slate-500 ring-1 ring-slate-200">
                              #{tag}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                    <Link
                      to={`/notes/view${note.deckId ? `?deckId=${note.deckId}&` : '?'}noteId=${note.id}`}
                      className="shrink-0 text-sm font-medium text-slate-700 hover:text-slate-950"
                    >
                      Edit
                    </Link>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="rounded-[2rem] border border-dashed border-slate-300 bg-white p-6 shadow-sm">
          <p className="text-xs uppercase tracking-[0.28em] text-slate-400">Study Groups</p>
          <h3 className="mt-4 text-2xl font-semibold tracking-tight text-slate-950">Deck-linked collaboration is planned next.</h3>
          <p className="mt-3 text-sm leading-7 text-slate-500">
            Team and Enterprise workspaces will be able to create study groups, invite members, and manage deck-centric group participation from one place.
          </p>
          <Link
            to="/study-groups"
            className="mt-6 inline-flex items-center rounded-2xl border border-slate-300 px-4 py-2.5 text-sm font-medium text-slate-800 hover:border-slate-400 hover:bg-stone-50"
          >
            Open Study Groups
          </Link>
        </div>
      </section>
    </div>
  )
}
