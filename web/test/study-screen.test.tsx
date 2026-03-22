import { afterEach, describe, expect, it, vi } from "vitest";
import {
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { StudyScreen } from "#/components/StudyScreen";
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
    status: "active" | "completed" | "abandoned";
    cardsReviewed: number;
    againCount: number;
    hardCount: number;
    goodCount: number;
    easyCount: number;
  }> = {},
) {
  return {
    id: overrides.id ?? "sts_1",
    userId: "usr_1",
    workspaceId: "ws_1",
    deckId: 1,
    mode: "review",
    status: overrides.status ?? "active",
    startedAt: "2026-03-22T12:00:00.000Z",
    endedAt:
      overrides.status && overrides.status !== "active"
        ? "2026-03-22T12:05:00.000Z"
        : undefined,
    cardsReviewed: overrides.cardsReviewed ?? 0,
    againCount: overrides.againCount ?? 0,
    hardCount: overrides.hardCount ?? 0,
    goodCount: overrides.goodCount ?? 0,
    easyCount: overrides.easyCount ?? 0,
    createdAt: "2026-03-22T12:00:00.000Z",
    updatedAt: "2026-03-22T12:05:00.000Z",
  };
}

function renderStudyScreen(repository: AppRepository, onExit = vi.fn()) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  });

  render(
    <QueryClientProvider client={queryClient}>
      <AppRepositoryProvider repository={repository}>
        <StudyScreen deckId={1} deckName="Biology" onExit={onExit} />
      </AppRepositoryProvider>
    </QueryClientProvider>,
  );

  return { onExit };
}

describe("StudyScreen", () => {
  it("creates and completes a study session when the last card is answered", async () => {
    const repository = {
      fetchDueCards: vi.fn().mockResolvedValue([
        {
          id: 10,
          noteId: 20,
          deckId: 1,
          templateName: "Card 1",
          ordinal: 0,
          front: "<p>Question</p>",
          back: "<p>Answer</p>",
          flag: 0,
          marked: false,
          suspended: false,
        },
      ]),
      createStudySession: vi.fn().mockResolvedValue(buildStudySession()),
      answerCard: vi.fn().mockResolvedValue({
        id: 10,
        noteId: 20,
        deckId: 1,
        templateName: "Card 1",
        ordinal: 0,
        front: "<p>Question</p>",
        back: "<p>Answer</p>",
        flag: 0,
        marked: false,
        suspended: false,
      }),
      updateStudySession: vi
        .fn()
        .mockImplementation(
          (_id: string, req: Record<string, number | string | undefined>) =>
            Promise.resolve(
              buildStudySession({
                status:
                  (req.status as
                    | "active"
                    | "completed"
                    | "abandoned"
                    | undefined) ?? "active",
                cardsReviewed: Number(req.cardsReviewed ?? 0),
                againCount: Number(req.againCount ?? 0),
                hardCount: Number(req.hardCount ?? 0),
                goodCount: Number(req.goodCount ?? 0),
                easyCount: Number(req.easyCount ?? 0),
              }),
            ),
        ),
      updateCard: vi.fn(),
    } as unknown as AppRepository;

    const { onExit } = renderStudyScreen(repository);

    await waitFor(() =>
      expect(repository.createStudySession).toHaveBeenCalled(),
    );

    fireEvent.click(await screen.findByRole("button", { name: "Show Answer" }));
    fireEvent.click(screen.getByRole("button", { name: /Good/i }));

    await waitFor(() =>
      expect(repository.answerCard).toHaveBeenCalledWith(
        10,
        expect.objectContaining({ rating: 3 }),
      ),
    );
    await waitFor(() =>
      expect(repository.updateStudySession).toHaveBeenCalledWith(
        "sts_1",
        expect.objectContaining({
          status: "completed",
          cardsReviewed: 1,
          goodCount: 1,
        }),
      ),
    );
    await waitFor(() => expect(onExit).toHaveBeenCalled());
  });

  it("abandons an active study session when the user exits early", async () => {
    const repository = {
      fetchDueCards: vi.fn().mockResolvedValue([
        {
          id: 10,
          noteId: 20,
          deckId: 1,
          templateName: "Card 1",
          ordinal: 0,
          front: "<p>Question</p>",
          back: "<p>Answer</p>",
          flag: 0,
          marked: false,
          suspended: false,
        },
      ]),
      createStudySession: vi.fn().mockResolvedValue(buildStudySession()),
      answerCard: vi.fn(),
      updateStudySession: vi
        .fn()
        .mockResolvedValue(buildStudySession({ status: "abandoned" })),
      updateCard: vi.fn(),
    } as unknown as AppRepository;

    const { onExit } = renderStudyScreen(repository);

    await waitFor(() =>
      expect(repository.createStudySession).toHaveBeenCalled(),
    );

    fireEvent.click(screen.getByRole("button", { name: "Exit" }));

    await waitFor(() =>
      expect(repository.updateStudySession).toHaveBeenCalledWith(
        "sts_1",
        expect.objectContaining({
          status: "abandoned",
          cardsReviewed: 0,
        }),
      ),
    );
    await waitFor(() => expect(onExit).toHaveBeenCalled());
  });
});
