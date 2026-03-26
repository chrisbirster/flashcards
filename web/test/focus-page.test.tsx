import { afterEach, describe, expect, it, vi } from "vitest";
import {
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { MemoryRouter } from "react-router";
import { FocusPage } from "#/pages/FocusPage";
import {
  AppRepositoryProvider,
  type AppRepository,
} from "#/lib/app-repository";

afterEach(() => {
  cleanup();
});

function buildStudySession(
  overrides: Partial<{
    id: string;
    mode: string;
    protocol: "pomodoro" | "deep-focus" | "custom";
    targetMinutes: number;
    breakMinutes: number;
    status: "active" | "completed" | "abandoned";
  }> = {},
) {
  return {
    id: overrides.id ?? "sts_focus_1",
    userId: "usr_1",
    workspaceId: "ws_1",
    mode: overrides.mode ?? "focus",
    protocol: overrides.protocol ?? "pomodoro",
    targetMinutes: overrides.targetMinutes ?? 25,
    breakMinutes: overrides.breakMinutes ?? 5,
    status: overrides.status ?? "active",
    startedAt: "2026-03-25T12:00:00.000Z",
    endedAt:
      overrides.status && overrides.status !== "active"
        ? "2026-03-25T12:25:00.000Z"
        : undefined,
    cardsReviewed: 0,
    againCount: 0,
    hardCount: 0,
    goodCount: 0,
    easyCount: 0,
    createdAt: "2026-03-25T12:00:00.000Z",
    updatedAt: "2026-03-25T12:25:00.000Z",
  };
}

function renderFocusPage(repository: AppRepository) {
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
          <FocusPage />
        </AppRepositoryProvider>
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

describe("FocusPage", () => {
  it("starts and completes a pomodoro focus session", async () => {
    const repository = {
      fetchStudyAnalyticsOverview: vi.fn().mockResolvedValue({
        sessions7d: 0,
        cardsReviewed7d: 0,
        minutesStudied7d: 0,
        focusSessions7d: 0,
        focusMinutes7d: 0,
        currentStreak: 0,
        answerBreakdown: { again: 0, hard: 0, good: 0, easy: 0 },
        dailyActivity: [],
        recentSessions: [],
      }),
      createStudySession: vi.fn().mockResolvedValue(buildStudySession()),
      updateStudySession: vi
        .fn()
        .mockResolvedValue(buildStudySession({ status: "completed" })),
    } as unknown as AppRepository;

    renderFocusPage(repository);

    fireEvent.click(await screen.findByRole("button", { name: /Start Pomodoro/i }));

    await waitFor(() =>
      expect(repository.createStudySession).toHaveBeenCalledWith({
        mode: "focus",
        protocol: "pomodoro",
        targetMinutes: 25,
        breakMinutes: 5,
      }),
    );

    fireEvent.click(screen.getByRole("button", { name: /Complete focus block/i }));

    await waitFor(() =>
      expect(repository.updateStudySession).toHaveBeenCalledWith(
        "sts_focus_1",
        expect.objectContaining({ status: "completed" }),
      ),
    );

    expect(
      await screen.findByText(/Focus block complete\. Take a 5 minute break/i),
    ).toBeInTheDocument();
  });
});
