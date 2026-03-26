import { afterEach, describe, expect, it, vi } from "vitest";
import { cleanup, render, screen } from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router";
import { StatsPage } from "#/pages/StatsPage";
import {
  AppRepositoryProvider,
  type AppRepository,
} from "#/lib/app-repository";

afterEach(() => {
  cleanup();
});

function renderStatsPage(repository: AppRepository) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  render(
    <MemoryRouter>
      <QueryClientProvider client={queryClient}>
        <AppRepositoryProvider repository={repository}>
          <StatsPage />
        </AppRepositoryProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

describe("StatsPage", () => {
  it("renders analytics overview, recent sessions, and active decks", async () => {
    const repository = {
      fetchStudyAnalyticsOverview: vi.fn().mockResolvedValue({
        sessions7d: 3,
        cardsReviewed7d: 47,
        minutesStudied7d: 62,
        focusSessions7d: 2,
        focusMinutes7d: 40,
        currentStreak: 4,
        lastStudiedAt: "2026-03-22T15:00:00.000Z",
        answerBreakdown: {
          again: 6,
          hard: 7,
          good: 24,
          easy: 10,
        },
        dailyActivity: [
          {
            date: "2026-03-16",
            sessions: 0,
            cardsReviewed: 0,
            minutesStudied: 0,
          },
          {
            date: "2026-03-17",
            sessions: 1,
            cardsReviewed: 8,
            minutesStudied: 12,
          },
          {
            date: "2026-03-18",
            sessions: 0,
            cardsReviewed: 0,
            minutesStudied: 0,
          },
          {
            date: "2026-03-19",
            sessions: 1,
            cardsReviewed: 14,
            minutesStudied: 20,
          },
          {
            date: "2026-03-20",
            sessions: 0,
            cardsReviewed: 0,
            minutesStudied: 0,
          },
          {
            date: "2026-03-21",
            sessions: 1,
            cardsReviewed: 25,
            minutesStudied: 30,
          },
          {
            date: "2026-03-22",
            sessions: 0,
            cardsReviewed: 0,
            minutesStudied: 0,
          },
        ],
        recentSessions: [
          {
            id: "sts_1",
            deckId: 1,
            deckName: "Biology",
            mode: "review",
            status: "completed",
            cardsReviewed: 25,
            minutesStudied: 30,
            againCount: 3,
            hardCount: 4,
            goodCount: 12,
            easyCount: 6,
            startedAt: "2026-03-21T14:00:00.000Z",
            endedAt: "2026-03-21T14:30:00.000Z",
            updatedAt: "2026-03-21T14:30:00.000Z",
          },
        ],
      }),
      fetchDecks: vi.fn().mockResolvedValue([
        {
          id: 1,
          name: "Biology",
            cardIds: [1, 2],
            dueToday: 6,
            dueReviewBacklog: 3,
            newCardsPerDay: 20,
            reviewsPerDay: 200,
            priorityOrder: 1,
            newCardsPaused: false,
            noteCount: 10,
            cardCount: 20,
          canDelete: false,
          deleteBlockedReason: "Deck is not empty.",
          analytics: {
            sessions7d: 2,
            cardsReviewed7d: 34,
            minutesStudied7d: 40,
            averageCardsPerSession7d: 17,
            againCount7d: 4,
            hardCount7d: 3,
            goodCount7d: 20,
            easyCount7d: 7,
            lastStudiedAt: "2026-03-21T14:30:00.000Z",
          },
        },
        {
          id: 2,
          name: "Chemistry",
            cardIds: [3],
            dueToday: 2,
            dueReviewBacklog: 1,
            newCardsPerDay: 10,
            reviewsPerDay: 100,
            priorityOrder: 2,
            newCardsPaused: false,
            noteCount: 3,
            cardCount: 5,
          canDelete: false,
          deleteBlockedReason: "Deck is not empty.",
          analytics: {
            sessions7d: 1,
            cardsReviewed7d: 13,
            minutesStudied7d: 22,
            averageCardsPerSession7d: 13,
            againCount7d: 2,
            hardCount7d: 4,
            goodCount7d: 5,
            easyCount7d: 2,
            lastStudiedAt: "2026-03-19T13:00:00.000Z",
          },
        },
      ]),
    } as unknown as AppRepository;

    renderStatsPage(repository);

    expect(await screen.findByText("Study analytics")).toBeInTheDocument();
    expect(await screen.findByText("Current Streak")).toBeInTheDocument();
    expect(await screen.findByText("Focus Blocks (7d)")).toBeInTheDocument();
    expect(await screen.findByText("Daily activity")).toBeInTheDocument();
    expect(await screen.findByText("Recent sessions")).toBeInTheDocument();
    expect(await screen.findAllByText("Biology")).toHaveLength(3);
    expect(
      await screen.findByText(/34 cards across 2 sessions/i),
    ).toBeInTheDocument();
    expect(
      await screen.findByText(/13 cards across 1 sessions/i),
    ).toBeInTheDocument();
  });
});
