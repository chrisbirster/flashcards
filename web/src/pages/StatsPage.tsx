import { Link } from "react-router";
import { useQuery } from "@tanstack/react-query";
import {
  EmptyState,
  PageContainer,
  PageSection,
  StatCard,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";

function formatShortDate(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
  }).format(new Date(value));
}

function formatDayLabel(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    weekday: "short",
  }).format(new Date(value));
}

function formatRelativeDateTime(value: string) {
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit",
  }).format(new Date(value));
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

function formatSessionModeLabel(mode: string, protocol?: string) {
  if (mode !== "focus") {
    return mode;
  }
  switch (protocol) {
    case "deep-focus":
      return "deep focus";
    case "custom":
      return "custom focus";
    default:
      return "pomodoro";
  }
}

function formatSessionTitle(session: {
  deckName?: string;
  mode: string;
  protocol?: string;
}) {
  if (session.mode === "focus") {
    return `${formatSessionModeLabel(session.mode, session.protocol)} block`;
  }
  return session.deckName || "Workspace study session";
}

function ActivityBars({
  dailyActivity,
}: {
  dailyActivity: Array<{
    date: string;
    sessions: number;
    cardsReviewed: number;
    minutesStudied: number;
  }>;
}) {
  const maxCards = Math.max(
    ...dailyActivity.map((entry) => entry.cardsReviewed),
    1,
  );

  return (
    <div className="space-y-4">
      <div className="flex items-end gap-2">
        {dailyActivity.map((entry) => {
          const height =
            entry.cardsReviewed > 0
              ? Math.max((entry.cardsReviewed / maxCards) * 100, 14)
              : 8;
          return (
            <div
              key={entry.date}
              className="flex min-w-0 flex-1 flex-col items-center gap-2"
            >
              <div className="flex h-40 w-full items-end rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-2">
                <div
                  className="w-full rounded-xl bg-[var(--app-accent)] shadow-[0_10px_30px_rgba(112,214,108,0.22)]"
                  style={{ height: `${height}%` }}
                  title={`${formatShortDate(entry.date)}: ${entry.cardsReviewed} cards, ${entry.sessions} sessions`}
                />
              </div>
              <div className="text-center">
                <p className="text-xs font-medium text-[var(--app-text)]">
                  {entry.cardsReviewed}
                </p>
                <p className="text-[11px] uppercase tracking-[0.16em] text-[var(--app-muted)]">
                  {formatDayLabel(entry.date)}
                </p>
              </div>
            </div>
          );
        })}
      </div>
      <p className="text-sm leading-6 text-[var(--app-text-soft)]">
        Bars show reviewed cards per day for the last week. Session totals and
        time studied stay visible in the summary cards above.
      </p>
    </div>
  );
}

function AnswerBreakdownCard({
  answerBreakdown,
}: {
  answerBreakdown: {
    again: number;
    hard: number;
    good: number;
    easy: number;
  };
}) {
  const segments = [
    {
      label: "Again",
      value: answerBreakdown.again,
      tone: "bg-rose-500/12 border-rose-500/30",
    },
    {
      label: "Hard",
      value: answerBreakdown.hard,
      tone: "bg-orange-500/12 border-orange-500/30",
    },
    {
      label: "Good",
      value: answerBreakdown.good,
      tone: "bg-emerald-500/12 border-emerald-500/30",
    },
    {
      label: "Easy",
      value: answerBreakdown.easy,
      tone: "bg-sky-500/12 border-sky-500/30",
    },
  ];
  const total = segments.reduce((sum, segment) => sum + segment.value, 0);

  return (
    <div className="space-y-3">
      {segments.map((segment) => {
        const percent =
          total > 0 ? Math.round((segment.value / total) * 100) : 0;
        return (
          <div
            key={segment.label}
            className="rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-3"
          >
            <div className="flex items-center justify-between gap-3">
              <div>
                <p className="text-sm font-semibold text-[var(--app-text)]">
                  {segment.label}
                </p>
                <p className="mt-1 text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">
                  {segment.value} answers
                </p>
              </div>
              <span
                className={`rounded-full border px-3 py-1 text-xs font-semibold text-[var(--app-text)] ${segment.tone}`}
              >
                {percent}%
              </span>
            </div>
          </div>
        );
      })}
    </div>
  );
}

export function StatsPage() {
  const repository = useAppRepository();
  const analyticsQuery = useQuery({
    queryKey: ["study-analytics"],
    queryFn: () => repository.fetchStudyAnalyticsOverview(),
  });
  const decksQuery = useQuery({
    queryKey: ["decks"],
    queryFn: () => repository.fetchDecks(),
  });

  const analytics = analyticsQuery.data;
  const topDecks = [...(decksQuery.data ?? [])]
    .filter(
      (deck) =>
        deck.analytics.cardsReviewed7d > 0 || deck.analytics.sessions7d > 0,
    )
    .sort((left, right) => {
      if (right.analytics.cardsReviewed7d !== left.analytics.cardsReviewed7d) {
        return right.analytics.cardsReviewed7d - left.analytics.cardsReviewed7d;
      }
      return right.analytics.sessions7d - left.analytics.sessions7d;
    })
    .slice(0, 5);

  const isLoading = analyticsQuery.isLoading || decksQuery.isLoading;
  const error =
    (analyticsQuery.error as Error | null) ??
    (decksQuery.error as Error | null);

  if (isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading study analytics...
        </PageSection>
      </PageContainer>
    );
  }

  if (error) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {error.message || "Failed to load analytics."}
        </PageSection>
      </PageContainer>
    );
  }

  if (!analytics) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="No analytics yet"
          description="Complete a few review or focus sessions and Vutadex will start surfacing streaks, answer mix, and deck trends here."
          action={
            <Link
              to="/decks"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Open decks
            </Link>
          }
        />
      </PageContainer>
    );
  }

  return (
    <PageContainer className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.3fr)_minmax(0,0.9fr)]">
        <div className="rounded-[2rem] bg-[var(--app-card-strong)] px-6 py-8 text-[var(--app-text)] shadow-sm md:px-8 md:py-10">
          <p className="text-xs uppercase tracking-[0.3em] text-[var(--app-accent)]">
            Study analytics
          </p>
          <h2 className="mt-4 max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl">
            Watch your weekly rhythm instead of guessing how steady your review
            habit feels.
          </h2>
          <p className="mt-4 max-w-2xl text-sm leading-7 text-[var(--app-text-soft)] md:text-base">
            This page rolls up persisted study sessions into daily activity,
            answer mix, and deck momentum so we can see where study is
            consistent and where it needs help.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              to="/decks"
              className="inline-flex items-center rounded-2xl bg-[var(--app-accent)] px-4 py-2.5 text-sm font-medium text-[var(--app-accent-ink)] transition hover:brightness-105"
            >
              Open decks
            </Link>
            <Link
              to="/study-groups"
              className="inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)]"
            >
              Group progress
            </Link>
          </div>
        </div>

        <PageSection className="p-6">
          <p className="text-xs uppercase tracking-[0.24em] text-[var(--app-muted)]">
            Last completed session
          </p>
          {analytics.recentSessions.length > 0 ? (
            <div className="mt-5 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-5">
              <p className="text-lg font-semibold text-[var(--app-text)]">
                {formatSessionTitle(analytics.recentSessions[0])}
              </p>
              <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                {formatRelativeDateTime(
                  analytics.recentSessions[0].endedAt ||
                    analytics.recentSessions[0].updatedAt,
                )}
              </p>
              <div className="mt-4 flex flex-wrap gap-2">
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                  {analytics.recentSessions[0].mode === "focus"
                    ? `${analytics.recentSessions[0].targetMinutes || 0}m target`
                    : `${analytics.recentSessions[0].cardsReviewed} cards reviewed`}
                </span>
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                  {formatMinutes(analytics.recentSessions[0].minutesStudied)}
                </span>
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                  {formatSessionModeLabel(
                    analytics.recentSessions[0].mode,
                    analytics.recentSessions[0].protocol,
                  )}
                </span>
              </div>
            </div>
          ) : (
            <p className="mt-4 text-sm text-[var(--app-text-soft)]">
              No completed study sessions yet.
            </p>
          )}
          <p className="mt-5 text-sm leading-6 text-[var(--app-text-soft)]">
            Current streak advances on days with completed review work or focus
            blocks, so the signal reflects both memory reps and deep-work rhythm.
          </p>
        </PageSection>
      </section>

      <section className="grid gap-4 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-6">
        <StatCard
          label="Current Streak"
          value={analytics.currentStreak}
          detail={
            analytics.lastStudiedAt
              ? `Last studied ${formatShortDate(analytics.lastStudiedAt)}.`
              : "No completed sessions yet."
          }
        />
        <StatCard
          label="Sessions (7d)"
          value={analytics.sessions7d}
          detail="Completed review sessions captured this week."
        />
        <StatCard
          label="Cards Reviewed (7d)"
          value={analytics.cardsReviewed7d}
          detail="Reviewed cards across all persisted sessions."
        />
        <StatCard
          label="Minutes Studied (7d)"
          value={analytics.minutesStudied7d}
          detail="Estimated time spent inside completed review sessions."
        />
        <StatCard
          label="Focus Blocks (7d)"
          value={analytics.focusSessions7d}
          detail="Completed pomodoro and focus sessions this week."
        />
        <StatCard
          label="Focus Minutes (7d)"
          value={analytics.focusMinutes7d}
          detail="Estimated time spent inside completed focus blocks."
        />
      </section>

      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)]">
        <PageSection className="p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Daily activity
              </h3>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Seven-day view of reviewed cards and session rhythm.
              </p>
            </div>
            <span className="rounded-full bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
              UTC calendar days
            </span>
          </div>
          <div className="mt-6">
            <ActivityBars dailyActivity={analytics.dailyActivity} />
          </div>
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Answer mix
          </h3>
          <p className="mt-1 text-sm text-[var(--app-text-soft)]">
            Review outcomes across the same seven-day window.
          </p>
          <div className="mt-6">
            <AnswerBreakdownCard answerBreakdown={analytics.answerBreakdown} />
          </div>
        </PageSection>
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <PageSection className="p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Recent sessions
              </h3>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                The last few completed or active study runs captured by Vutadex.
              </p>
            </div>
          </div>

          {analytics.recentSessions.length === 0 ? (
            <p className="mt-6 text-sm text-[var(--app-text-soft)]">
              No completed sessions yet.
            </p>
          ) : (
            <ul className="mt-6 space-y-3">
              {analytics.recentSessions.map((session) => (
                <li
                  key={session.id}
                  className="rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4"
                >
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div className="min-w-0">
                      <p className="text-sm font-semibold text-[var(--app-text)]">
                        {formatSessionTitle(session)}
                      </p>
                      <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                        {formatRelativeDateTime(
                          session.endedAt || session.updatedAt,
                        )}
                      </p>
                    </div>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {formatSessionModeLabel(session.mode, session.protocol)}
                    </span>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    {session.mode === "focus" ? (
                      <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                        {session.targetMinutes || 0}m target / {session.breakMinutes || 0}m break
                      </span>
                    ) : (
                      <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                        {session.cardsReviewed} cards
                      </span>
                    )}
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {formatMinutes(session.minutesStudied)}
                    </span>
                    {session.mode === "focus" ? null : (
                      <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                        {session.goodCount} good / {session.againCount} again
                      </span>
                    )}
                  </div>
                </li>
              ))}
            </ul>
          )}
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Active decks
              </h3>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Ranked by reviewed cards from the last week.
              </p>
            </div>
            <Link
              to="/decks"
              className="text-sm font-medium text-[var(--app-accent)] hover:brightness-110"
            >
              Open decks
            </Link>
          </div>

          {topDecks.length === 0 ? (
            <p className="mt-6 text-sm text-[var(--app-text-soft)]">
              Deck-level analytics will appear after you complete a few study
              sessions.
            </p>
          ) : (
            <ul className="mt-6 space-y-3">
              {topDecks.map((deck) => (
                <li
                  key={deck.id}
                  className="rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div className="min-w-0">
                      <p className="truncate text-sm font-semibold text-[var(--app-text)]">
                        {deck.name}
                      </p>
                      <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                        {deck.analytics.cardsReviewed7d} cards across{" "}
                        {deck.analytics.sessions7d} sessions
                      </p>
                    </div>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {deck.dueToday} due
                    </span>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {formatMinutes(deck.analytics.minutesStudied7d)}
                    </span>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {deck.analytics.averageCardsPerSession7d.toFixed(1)} cards
                      / session
                    </span>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      Good {deck.analytics.goodCount7d}
                    </span>
                  </div>
                </li>
              ))}
            </ul>
          )}
        </PageSection>
      </section>
    </PageContainer>
  );
}
