import { Link } from "react-router";
import { useQuery } from "@tanstack/react-query";
import { useAppRepository } from "#/lib/app-repository";

function StatCard({
  label,
  value,
  detail,
}: {
  label: string;
  value: string | number;
  detail: string;
}) {
  return (
    <div className="rounded-3xl border border-[var(--app-line)] bg-[var(--app-card)] p-5 shadow-sm">
      <p className="text-xs uppercase tracking-[0.24em] text-[var(--app-muted)]">
        {label}
      </p>
      <p className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)]">
        {value}
      </p>
      <p className="mt-2 text-sm text-[var(--app-text-soft)]">{detail}</p>
    </div>
  );
}

function formatLastStudiedAt(value?: string) {
  if (!value) {
    return "No completed study sessions yet.";
  }
  return `Last studied ${new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
  }).format(new Date(value))}.`;
}

export function HomePage() {
  const repository = useAppRepository();
  const dashboardQuery = useQuery({
    queryKey: ["dashboard"],
    queryFn: () => repository.fetchDashboard(),
  });

  const plan = dashboardQuery.data?.plan?.toUpperCase() ?? "FREE";
  const noteUsage = dashboardQuery.data?.usage.notes ?? 0;
  const noteLimit = dashboardQuery.data?.limits.maxNotes ?? 0;
  const cardUsage = dashboardQuery.data?.usage.cardsTotal ?? 0;
  const cardLimit = dashboardQuery.data?.limits.maxCardsTotal ?? 0;
  const totalDecks = dashboardQuery.data?.totalDecks ?? 0;
  const totalNotes = dashboardQuery.data?.totalNotes ?? 0;
  const dueToday = dashboardQuery.data?.dueToday ?? 0;
  const studyAnalytics = dashboardQuery.data?.studyAnalytics;

  return (
    <div className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.4fr)_minmax(0,0.9fr)]">
        <div className="rounded-[2rem] bg-[var(--app-card-strong)] px-6 py-8 text-[var(--app-text)] shadow-sm md:px-8 md:py-10">
          <p className="text-xs uppercase tracking-[0.3em] text-[var(--app-accent)]">
            Workspace overview
          </p>
          <h2 className="mt-4 max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl">
            Keep deck building, review, and template tuning in one browser-first
            workspace.
          </h2>
          <p className="mt-4 max-w-2xl text-sm leading-7 text-[var(--app-text-soft)] md:text-base">
            Use the dashboard to watch plan usage, jump into notes management,
            and keep new cards moving into study without losing structure.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              to="/notes/add"
              className="inline-flex items-center rounded-2xl bg-[var(--app-accent)] px-4 py-2.5 text-sm font-medium text-[var(--app-accent-ink)] transition hover:brightness-105"
            >
              Add note
            </Link>
            <Link
              to="/notes/view"
              className="inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)]"
            >
              Open notes
            </Link>
            <Link
              to="/templates"
              className="inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)]"
            >
              Open templates
            </Link>
            <Link
              to="/stats"
              className="inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)]"
            >
              Open analytics
            </Link>
          </div>
        </div>

        <div className="rounded-[2rem] border border-[var(--app-line)] bg-[var(--app-card)] p-6 shadow-sm">
          <p className="text-xs uppercase tracking-[0.28em] text-[var(--app-muted)]">
            Plan usage
          </p>
          <div className="mt-5 space-y-5">
            <div>
              <div className="flex items-center justify-between text-sm text-[var(--app-text-soft)]">
                <span>Notes</span>
                <span>
                  {noteUsage} / {noteLimit}
                </span>
              </div>
              <div className="mt-2 h-2 rounded-full bg-[var(--app-muted-surface)]">
                <div
                  className="h-2 rounded-full bg-[var(--app-accent)]"
                  style={{
                    width: `${noteLimit > 0 ? Math.min((noteUsage / noteLimit) * 100, 100) : 0}%`,
                  }}
                />
              </div>
            </div>
            <div>
              <div className="flex items-center justify-between text-sm text-[var(--app-text-soft)]">
                <span>Total cards</span>
                <span>
                  {cardUsage} / {cardLimit}
                </span>
              </div>
              <div className="mt-2 h-2 rounded-full bg-[var(--app-muted-surface)]">
                <div
                  className="h-2 rounded-full bg-[var(--app-accent-strong)]"
                  style={{
                    width: `${cardLimit > 0 ? Math.min((cardUsage / cardLimit) * 100, 100) : 0}%`,
                  }}
                />
              </div>
            </div>
            <div className="rounded-2xl bg-[var(--app-muted-surface)] p-4">
              <p className="text-xs uppercase tracking-[0.2em] text-[var(--app-muted)]">
                Current plan
              </p>
              <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                {plan}
              </p>
              <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                Study Groups creation unlocks on Team. Marketplace publishing
                arrives in a later tranche for eligible plans.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Total Decks"
          value={totalDecks}
          detail="Active decks in this workspace."
        />
        <StatCard
          label="Total Notes"
          value={totalNotes}
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

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Current Streak"
          value={studyAnalytics?.currentStreak ?? 0}
          detail={formatLastStudiedAt(studyAnalytics?.lastStudiedAt)}
        />
        <StatCard
          label="Sessions (7d)"
          value={studyAnalytics?.sessions7d ?? 0}
          detail="Completed review sessions from the last week."
        />
        <StatCard
          label="Cards Reviewed (7d)"
          value={studyAnalytics?.cardsReviewed7d ?? 0}
          detail="Answered cards captured through persisted study sessions."
        />
        <StatCard
          label="Minutes Studied (7d)"
          value={studyAnalytics?.minutesStudied7d ?? 0}
          detail="Time accumulated from completed sessions this week."
        />
      </section>

      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
        <div className="rounded-[2rem] border border-[var(--app-line)] bg-[var(--app-card)] p-6 shadow-sm">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Recent note activity
              </h3>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Jump back into the notes you created or edited most recently.
              </p>
            </div>
            <Link
              to="/notes/view"
              className="text-sm font-medium text-[var(--app-accent)] hover:brightness-110"
            >
              View all
            </Link>
          </div>

          {dashboardQuery.isLoading && (
            <p className="mt-6 text-sm text-[var(--app-text-soft)]">
              Loading recent notes...
            </p>
          )}
          {!dashboardQuery.isLoading &&
            (dashboardQuery.data?.recentNotes.length ?? 0) === 0 && (
              <p className="mt-6 text-sm text-[var(--app-text-soft)]">
                No notes yet. Start with your first note and the dashboard will
                reflect it here.
              </p>
            )}
          {(dashboardQuery.data?.recentNotes?.length ?? 0) > 0 && (
            <ul className="mt-6 space-y-3">
              {dashboardQuery.data?.recentNotes.map((note) => (
                <li
                  key={note.id}
                  className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-2">
                        <span className="rounded-full bg-[var(--app-card)] px-2.5 py-1 text-xs font-medium text-[var(--app-text)] ring-1 ring-[var(--app-line)]">
                          {note.typeId}
                        </span>
                        {note.deckName && (
                          <span className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">
                            {note.deckName}
                          </span>
                        )}
                      </div>
                      <p className="mt-2 truncate text-sm font-medium text-[var(--app-text)]">
                        {note.fieldPreview || "Untitled note"}
                      </p>
                      {note.tags.length > 0 && (
                        <div className="mt-2 flex flex-wrap gap-1">
                          {note.tags.slice(0, 3).map((tag) => (
                            <span
                              key={tag}
                              className="rounded-full bg-[var(--app-card)] px-2 py-1 text-xs text-[var(--app-text-soft)] ring-1 ring-[var(--app-line)]"
                            >
                              #{tag}
                            </span>
                          ))}
                        </div>
                      )}
                    </div>
                    <Link
                      to={`/notes/view${note.deckId ? `?deckId=${note.deckId}&` : "?"}noteId=${note.id}`}
                      className="shrink-0 text-sm font-medium text-[var(--app-accent)] hover:brightness-110"
                    >
                      Edit
                    </Link>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </div>

        <div className="rounded-[2rem] border border-dashed border-[var(--app-line-strong)] bg-[var(--app-card)] p-6 shadow-sm">
          <p className="text-xs uppercase tracking-[0.28em] text-[var(--app-muted)]">
            Study Groups
          </p>
          <h3 className="mt-4 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
            Deck-linked collaboration is planned next.
          </h3>
          <p className="mt-3 text-sm leading-7 text-[var(--app-text-soft)]">
            Team and Enterprise workspaces will be able to create study groups,
            invite members, and manage deck-centric group participation from one
            place.
          </p>
          <Link
            to="/study-groups"
            className="mt-6 inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-muted-surface)]"
          >
            Open Study Groups
          </Link>
        </div>
      </section>
    </div>
  );
}
