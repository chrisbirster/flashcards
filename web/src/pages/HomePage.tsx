import { Link } from "react-router";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import { EmptyState, PageContainer, PageSection, SurfaceCard } from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";
import type { Deck, NoteListItem, StudySessionSummary } from "#/lib/api";

function formatPlanLabel(plan: "guest" | "free" | "pro" | "team" | "enterprise") {
  return plan.charAt(0).toUpperCase() + plan.slice(1);
}

function formatMinutes(value: number) {
  if (value <= 0) {
    return "0m";
  }
  if (value < 60) {
    return `${value}m`;
  }
  const hours = Math.floor(value / 60);
  const minutes = value % 60;
  return minutes > 0 ? `${hours}h ${minutes}m` : `${hours}h`;
}

function formatDateTime(value?: string) {
  if (!value) {
    return "No study sessions yet";
  }
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return "No study sessions yet";
  }
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  }).format(date);
}

function formatUsage(used: number, limit: number) {
  return `${used.toLocaleString()} / ${limit.toLocaleString()}`;
}

function usagePercent(used: number, limit: number) {
  if (limit <= 0) {
    return 0;
  }
  return Math.min(100, Math.max(0, Math.round((used / limit) * 100)));
}

function sessionLabel(session: StudySessionSummary) {
  if (session.mode === "focus") {
    if (session.protocol === "deep-focus") return "Deep focus";
    if (session.protocol === "custom") return "Custom focus";
    return "Pomodoro";
  }
  return session.deckName || "Study session";
}

function notePreview(note: NoteListItem) {
  return note.fieldPreview || Object.values(note.fieldVals ?? {})[0] || "Untitled note";
}

function OverviewStat({
  label,
  value,
  detail,
}: {
  label: string;
  value: string | number;
  detail: string;
}) {
  return (
    <div className="rounded-[1.5rem] border border-[var(--app-stat-line)] bg-[var(--app-stat-surface)] p-5 shadow-sm">
      <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">{label}</p>
      <p className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)]">{value}</p>
      <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">{detail}</p>
    </div>
  );
}

function UsageMeter({
  label,
  used,
  limit,
}: {
  label: string;
  used: number;
  limit: number;
}) {
  const percent = usagePercent(used, limit);
  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between gap-3">
        <p className="text-sm font-medium text-[var(--app-text)]">{label}</p>
        <p className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">
          {formatUsage(used, limit)}
        </p>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-[var(--app-muted-surface)]">
        <div
          className="h-full rounded-full bg-[var(--app-accent)] transition-[width]"
          style={{ width: `${percent}%` }}
        />
      </div>
    </div>
  );
}

function DueDeckRow({
  deck,
  recommended,
}: {
  deck: Deck;
  recommended?: boolean;
}) {
  return (
    <div className="rounded-[1.4rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4">
      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <p className="truncate text-base font-semibold text-[var(--app-text)]">{deck.name}</p>
            {recommended ? (
              <span className="rounded-full bg-[var(--app-accent)] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.16em] text-[var(--app-accent-ink)]">
                Recommended
              </span>
            ) : null}
          </div>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Priority {deck.priorityOrder + 1} • {deck.noteCount} note{deck.noteCount === 1 ? "" : "s"} •{" "}
            {deck.cardCount} card{deck.cardCount === 1 ? "" : "s"}
          </p>
          <div className="mt-3 flex flex-wrap gap-2">
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text)]">
              {deck.dueToday} due today
            </span>
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text)]">
              backlog {deck.dueReviewBacklog}
            </span>
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text)]">
              {deck.newCardsPerDay} new / {deck.reviewsPerDay} reviews
            </span>
          </div>
          {deck.newCardsPaused ? (
            <p className="mt-3 text-sm leading-6 text-[var(--app-warning-text)]">
              New cards paused until backlog drops below your review cap.
            </p>
          ) : null}
        </div>
        <div className="flex flex-wrap gap-3">
          <Link
            to={`/study/${deck.id}`}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
          >
            Study deck
          </Link>
          <Link
            to="/decks"
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
          >
            Manage
          </Link>
        </div>
      </div>
    </div>
  );
}

function QuickAction({
  to,
  title,
  description,
}: {
  to: string;
  title: string;
  description: string;
}) {
  return (
    <Link
      to={to}
      className="group rounded-[1.4rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4 transition hover:border-[var(--app-line-strong)] hover:bg-[var(--app-card-strong)]"
    >
      <p className="text-base font-semibold text-[var(--app-text)]">{title}</p>
      <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">{description}</p>
      <p className="mt-4 text-xs font-semibold uppercase tracking-[0.18em] text-[var(--app-accent)]">
        Open
      </p>
    </Link>
  );
}

function RecentNoteRow({ note }: { note: NoteListItem }) {
  return (
    <Link
      to={`/notes/view?noteId=${note.id}`}
      className="block rounded-[1.3rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4 transition hover:border-[var(--app-line-strong)] hover:bg-[var(--app-card-strong)]"
    >
      <div className="flex flex-wrap items-center gap-2">
        <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.16em] text-[var(--app-muted)]">
          {note.deckName || "No deck"}
        </span>
        <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.16em] text-[var(--app-muted)]">
          {note.typeId}
        </span>
      </div>
      <p className="mt-3 line-clamp-2 text-base font-semibold text-[var(--app-text)]">{notePreview(note)}</p>
      <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
        {note.cardCount} card{note.cardCount === 1 ? "" : "s"} • updated {formatDateTime(note.modifiedAt)}
      </p>
    </Link>
  );
}

export function HomePage() {
  const repository = useAppRepository();

  const dashboardQuery = useQuery({
    queryKey: ["dashboard"],
    queryFn: () => repository.fetchDashboard(),
  });
  const decksQuery = useQuery({
    queryKey: ["decks"],
    queryFn: () => repository.fetchDecks(),
  });
  const sessionQuery = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const isLoading =
    dashboardQuery.isLoading || decksQuery.isLoading || sessionQuery.isLoading;
  const error =
    (dashboardQuery.error as Error | null) ??
    (decksQuery.error as Error | null) ??
    (sessionQuery.error as Error | null);

  const dashboard = dashboardQuery.data;
  const workspaceName = sessionQuery.data?.workspace?.name || "Workspace";

  const sortedDecks = useMemo(
    () =>
      [...(decksQuery.data ?? [])].sort((left, right) => {
        if (left.priorityOrder !== right.priorityOrder) {
          return left.priorityOrder - right.priorityOrder;
        }
        if (right.dueToday !== left.dueToday) {
          return right.dueToday - left.dueToday;
        }
        return left.name.localeCompare(right.name);
      }),
    [decksQuery.data],
  );

  const dueDecks = sortedDecks.filter((deck) => deck.dueToday > 0);
  const recommendedDeck = dueDecks[0] ?? sortedDecks[0];
  const recentNotes = (dashboard?.recentNotes ?? []).slice(0, 5);
  const recentSessions = (dashboard?.studyAnalytics.recentSessions ?? []).slice(0, 3);

  if (isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading workspace overview...
        </PageSection>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {error.message || "Failed to load the dashboard."}
        </PageSection>
      </PageContainer>
    );
  }

  if (!dashboard) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="No dashboard data yet"
          description="Create a deck or add a few notes and Vutadex will start surfacing your study queue, usage, and momentum here."
          action={
            <Link
              to="/notes/add"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Add your first note
            </Link>
          }
        />
      </PageContainer>
    );
  }

  return (
    <PageContainer className="space-y-6">
      <PageSection className="overflow-hidden bg-[var(--app-card-strong)] p-6 sm:p-8">
        <div className="grid gap-6 xl:grid-cols-[minmax(0,1.2fr)_minmax(18rem,0.8fr)]">
          <div>
            <p className="text-[11px] uppercase tracking-[0.28em] text-[var(--app-accent)]">Overview</p>
            <h1 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--app-text)] sm:text-4xl">
              {workspaceName}
            </h1>
            <p className="mt-4 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)] sm:text-base">
              A calmer view of what matters today: what is due, where your study time is going,
              and the next deck worth opening.
            </p>
          </div>

          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-1">
            <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card)]">
              <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Today</p>
              <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
                {dashboard.dueToday > 0
                  ? `${dashboard.dueToday} cards waiting across ${dueDecks.length || 1} deck${dueDecks.length === 1 ? "" : "s"}`
                  : "You are caught up for now"}
              </p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Review cards come first. New cards pause automatically when a deck backlog runs past its review cap.
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card)]">
              <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Plan</p>
              <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
                {formatPlanLabel(dashboard.plan)}
              </p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                {dashboard.features.studyGroups
                  ? "Collaboration and published installs are available in this workspace."
                  : "Upgrade when you want collaboration and publishing features unlocked."}
              </p>
            </SurfaceCard>
          </div>
        </div>
      </PageSection>

      <div className="grid gap-4 sm:grid-cols-2 xl:grid-cols-4">
        <OverviewStat
          label="Total decks"
          value={dashboard.totalDecks}
          detail="Organized by deck priority so the home queue stays predictable."
        />
        <OverviewStat
          label="Total notes"
          value={dashboard.totalNotes}
          detail="Your note count tracks the source material feeding card creation."
        />
        <OverviewStat
          label="Due today"
          value={dashboard.dueToday}
          detail="Review and relearning cards stay visible even if you skip a day."
        />
        <OverviewStat
          label="Current streak"
          value={dashboard.studyAnalytics.currentStreak}
          detail={`Last study activity ${formatDateTime(dashboard.studyAnalytics.lastStudiedAt)}.`}
        />
      </div>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)]">
        <PageSection className="p-5 sm:p-6">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Study queue</p>
              <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
                Start with the highest-priority due deck.
              </h2>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Deck priority now controls the order here. Each deck still respects its own new-card and review caps.
              </p>
            </div>
            <Link
              to="/decks"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
            >
              Reorder decks
            </Link>
          </div>

          <div className="mt-5 space-y-4">
            {recommendedDeck ? (
              <DueDeckRow deck={recommendedDeck} recommended />
            ) : (
              <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card-strong)]">
                <p className="text-lg font-semibold text-[var(--app-text)]">Nothing is due right now.</p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  Use this window to add notes, tune deck caps, or start a focus block before the next review wave lands.
                </p>
                <div className="mt-4 flex flex-wrap gap-3">
                  <Link
                    to="/notes/add"
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
                  >
                    Add note
                  </Link>
                  <Link
                    to="/focus"
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-5 text-sm font-semibold text-[var(--app-text)]"
                  >
                    Start focus block
                  </Link>
                </div>
              </SurfaceCard>
            )}

            {dueDecks.length > 1 ? (
              <div className="space-y-3">
                <div className="flex items-center justify-between gap-3">
                  <h3 className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--app-muted)]">
                    Next in line
                  </h3>
                  <span className="text-sm text-[var(--app-text-soft)]">
                    {dueDecks.length - 1} more deck{dueDecks.length - 1 === 1 ? "" : "s"} with work ready
                  </span>
                </div>
                {dueDecks.slice(1, 5).map((deck) => (
                  <DueDeckRow key={deck.id} deck={deck} />
                ))}
              </div>
            ) : null}
          </div>
        </PageSection>

        <div className="grid gap-4">
          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Plan & usage</p>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
              Capacity at a glance.
            </h2>
            <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
              This workspace is on {formatPlanLabel(dashboard.plan)}. Study Groups and Marketplace installs still keep private review history per user.
            </p>
            <div className="mt-5 space-y-4">
              <UsageMeter label="Decks" used={dashboard.usage.decks} limit={dashboard.limits.maxDecks} />
              <UsageMeter label="Notes" used={dashboard.usage.notes} limit={dashboard.limits.maxNotes} />
              <UsageMeter
                label="Cards"
                used={dashboard.usage.cardsTotal}
                limit={dashboard.limits.maxCardsTotal}
              />
            </div>
            <div className="mt-5 flex flex-wrap gap-3">
              <Link
                to="/settings"
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
              >
                Manage plan
              </Link>
              {sessionQuery.data?.workspace?.organizationId ? (
                <Link
                  to="/team"
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
                >
                  Open team
                </Link>
              ) : null}
            </div>
          </PageSection>

          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Weekly momentum</p>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
              Stay steady instead of spiky.
            </h2>
            <div className="mt-5 grid gap-3 sm:grid-cols-2">
              <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">Sessions</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                  {dashboard.studyAnalytics.sessions7d}
                </p>
              </SurfaceCard>
              <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">Cards reviewed</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                  {dashboard.studyAnalytics.cardsReviewed7d}
                </p>
              </SurfaceCard>
              <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">Minutes</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                  {formatMinutes(dashboard.studyAnalytics.minutesStudied7d)}
                </p>
              </SurfaceCard>
              <SurfaceCard className="border-[var(--app-line-strong)] bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">Focus blocks</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                  {dashboard.studyAnalytics.focusSessions7d}
                </p>
              </SurfaceCard>
            </div>

            {recentSessions.length > 0 ? (
              <div className="mt-5 space-y-3">
                <h3 className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--app-muted)]">
                  Recent sessions
                </h3>
                {recentSessions.map((session) => (
                  <div
                    key={session.id}
                    className="rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-card)] px-4 py-3"
                  >
                    <div className="flex items-center justify-between gap-3">
                      <p className="text-sm font-semibold text-[var(--app-text)]">{sessionLabel(session)}</p>
                      <span className="text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">
                        {session.status}
                      </span>
                    </div>
                    <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                      {session.cardsReviewed} cards • {formatMinutes(session.minutesStudied)} •{" "}
                      {formatDateTime(session.endedAt || session.updatedAt)}
                    </p>
                  </div>
                ))}
              </div>
            ) : null}
          </PageSection>
        </div>
      </div>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,0.8fr)_minmax(0,1.2fr)]">
        <div className="grid gap-4">
          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Quick actions</p>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
              Keep the next move obvious.
            </h2>
            <div className="mt-5 grid gap-3 sm:grid-cols-2">
              <QuickAction
                to="/notes/add"
                title="Add note"
                description="Capture source material and turn it into cards without leaving the app shell."
              />
              <QuickAction
                to="/decks"
                title="Open decks"
                description="Adjust caps, reorder priorities, or jump into a specific deck."
              />
              <QuickAction
                to="/focus"
                title="Start focus block"
                description="Use a Pomodoro or deep-focus session when you want time-boxed work."
              />
              <QuickAction
                to="/stats"
                title="Review stats"
                description="Check streaks, answer mix, and weekly study volume."
              />
            </div>
          </PageSection>

          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Collaboration</p>
            <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
              Study Groups stay private by design.
            </h2>
            <p className="mt-3 text-sm leading-7 text-[var(--app-text-soft)]">
              Members install workspace-local copies, keep private review history, and opt into source updates when they are ready.
            </p>
            <div className="mt-5 flex flex-wrap gap-3">
              <Link
                to="/study-groups"
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
              >
                Open Study Groups
              </Link>
              {dashboard.features.studyGroups ? null : (
                <Link
                  to="/settings"
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
                >
                  View plan options
                </Link>
              )}
            </div>
          </PageSection>
        </div>

        <PageSection className="p-5 sm:p-6">
          <div className="flex flex-col gap-2 sm:flex-row sm:items-end sm:justify-between">
            <div>
              <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Recent notes</p>
              <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
                Recent source material, not just card counts.
              </h2>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                The most recent notes stay close so editing and refinement is only one tap away.
              </p>
            </div>
            <Link
              to="/notes/view"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
            >
              View all notes
            </Link>
          </div>

          <div className="mt-5 space-y-3">
            {recentNotes.length > 0 ? (
              recentNotes.map((note) => <RecentNoteRow key={note.id} note={note} />)
            ) : (
              <EmptyState
                title="No recent notes yet"
                description="Add a note and Vutadex will keep the newest material here for quick editing and review."
                action={
                  <Link
                    to="/notes/add"
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
                  >
                    Add note
                  </Link>
                }
              />
            )}
          </div>
        </PageSection>
      </div>
    </PageContainer>
  );
}
