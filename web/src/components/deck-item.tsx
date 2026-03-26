import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Deck } from "#/lib/api";
import { useAppRepository } from "#/lib/app-repository";
import { ConfirmSheet, Sheet } from "#/components/sheet";

function StatPill({ label, value }: { label: string; value: string | number }) {
  return (
    <span className="inline-flex items-center gap-1 rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
      <span className="text-[var(--app-text)]">{value}</span>
      <span>{label}</span>
    </span>
  );
}

function formatLastStudiedAt(value?: string) {
  if (!value) {
    return "Not studied yet.";
  }
  return `Last studied ${new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
  }).format(new Date(value))}.`;
}

export function DeckItem({ deck }: { deck: Deck }) {
  const navigate = useNavigate();
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const [isEditing, setIsEditing] = useState(false);
  const [draftName, setDraftName] = useState(deck.name);
  const [actionError, setActionError] = useState<string | null>(null);
  const [actionsOpen, setActionsOpen] = useState(false);
  const [confirmDeleteOpen, setConfirmDeleteOpen] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(false);
  const [newCardsPerDay, setNewCardsPerDay] = useState(String(deck.newCardsPerDay));
  const [reviewsPerDay, setReviewsPerDay] = useState(String(deck.reviewsPerDay));
  const [priorityOrder, setPriorityOrder] = useState(String(deck.priorityOrder));

  useEffect(() => {
    setDraftName(deck.name);
    setNewCardsPerDay(String(deck.newCardsPerDay));
    setReviewsPerDay(String(deck.reviewsPerDay));
    setPriorityOrder(String(deck.priorityOrder));
  }, [deck.id, deck.name, deck.newCardsPerDay, deck.priorityOrder, deck.reviewsPerDay]);

  const renameMutation = useMutation({
    mutationFn: (name: string) => repository.updateDeck(deck.id, { name }),
    onSuccess: () => {
      setIsEditing(false);
      setActionError(null);
      queryClient.invalidateQueries({ queryKey: ["decks"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
    onError: (error: Error) => setActionError(error.message),
  });

  const deleteMutation = useMutation({
    mutationFn: () => repository.deleteDeck(deck.id),
    onSuccess: () => {
      setActionError(null);
      queryClient.invalidateQueries({ queryKey: ["decks"] });
      queryClient.invalidateQueries({ queryKey: ["entitlements"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
    onError: (error: Error) => setActionError(error.message),
  });

  const settingsMutation = useMutation({
    mutationFn: (payload: {
      newCardsPerDay: number;
      reviewsPerDay: number;
      priorityOrder: number;
    }) => repository.updateDeck(deck.id, payload),
    onSuccess: () => {
      setActionError(null);
      setSettingsOpen(false);
      queryClient.invalidateQueries({ queryKey: ["decks"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
    onError: (error: Error) => setActionError(error.message),
  });

  const canStudy = deck.dueToday > 0;

  const runDelete = () => {
    deleteMutation.mutate();
  };

  const submitSettings = (event: React.FormEvent) => {
    event.preventDefault();
    const parsedNew = Number(newCardsPerDay);
    const parsedReviews = Number(reviewsPerDay);
    const parsedPriority = Number(priorityOrder);
    if (
      !Number.isFinite(parsedNew) ||
      !Number.isFinite(parsedReviews) ||
      !Number.isFinite(parsedPriority)
    ) {
      setActionError("Workload settings must be valid numbers.");
      return;
    }
    settingsMutation.mutate({
      newCardsPerDay: Math.max(0, Math.floor(parsedNew)),
      reviewsPerDay: Math.max(0, Math.floor(parsedReviews)),
      priorityOrder: Math.max(1, Math.floor(parsedPriority)),
    });
  };

  return (
    <>
      <li className="rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4 shadow-sm sm:p-5">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            {isEditing ? (
              <form
                onSubmit={(event) => {
                  event.preventDefault();
                  if (!draftName.trim()) return;
                  renameMutation.mutate(draftName.trim());
                }}
                className="space-y-3"
              >
                <input
                  type="text"
                  value={draftName}
                  onChange={(event) => setDraftName(event.target.value)}
                  className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                />
                <div className="flex flex-col gap-2 sm:flex-row">
                  <button
                    type="submit"
                    disabled={renameMutation.isPending || !draftName.trim()}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
                  >
                    {renameMutation.isPending ? "Saving..." : "Save"}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setIsEditing(false);
                      setDraftName(deck.name);
                    }}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            ) : (
              <>
                <div className="flex flex-wrap items-center gap-3">
                  <h3 className="truncate text-lg font-semibold text-[var(--app-text)]">
                    {deck.name}
                  </h3>
                  {deck.dueToday > 0 ? (
                    <span className="rounded-full bg-[var(--app-success-surface)] px-2.5 py-1 text-xs font-semibold text-[var(--app-success-text)] ring-1 ring-[var(--app-success-line)]">
                      {deck.dueToday} due today
                    </span>
                  ) : null}
                </div>

                <div className="mt-3 flex flex-wrap gap-2">
                  <StatPill label="notes" value={deck.noteCount} />
                  <StatPill label="cards" value={deck.cardCount} />
                  <StatPill label="due" value={deck.dueToday} />
                  <StatPill label="review backlog" value={deck.dueReviewBacklog} />
                  <StatPill label="new/day" value={deck.newCardsPerDay} />
                  <StatPill label="reviews/day" value={deck.reviewsPerDay} />
                  <StatPill label="priority" value={deck.priorityOrder} />
                  <StatPill
                    label="sessions (7d)"
                    value={deck.analytics.sessions7d}
                  />
                  <StatPill
                    label="reviewed (7d)"
                    value={deck.analytics.cardsReviewed7d}
                  />
                </div>

                <p className="mt-3 text-sm leading-6 text-[var(--app-text-soft)]">
                  {deck.canDelete
                    ? "Empty decks can be deleted directly."
                    : deck.deleteBlockedReason ||
                      "Delete is disabled until this deck is empty."}
                </p>
                {deck.newCardsPaused ? (
                  <p className="mt-2 text-sm leading-6 text-[var(--app-warning-text)]">
                    New cards paused until review backlog drops below your review cap.
                  </p>
                ) : (
                  <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                    New cards stay available while review backlog is at or below your review cap.
                  </p>
                )}
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  {formatLastStudiedAt(deck.analytics.lastStudiedAt)}
                </p>

                {actionError ? (
                  <p className="mt-3 rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 py-3 text-sm text-[var(--app-danger-text)]">
                    {actionError}
                  </p>
                ) : null}
              </>
            )}
          </div>

          {!isEditing ? (
            <button
              type="button"
              onClick={() => setActionsOpen(true)}
              className="inline-flex min-h-11 min-w-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text-soft)] md:hidden"
              aria-label={`Open actions for ${deck.name}`}
            >
              <svg
                className="h-5 w-5"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.8"
                  d="M12 6h.01M12 12h.01M12 18h.01"
                />
              </svg>
            </button>
          ) : null}
        </div>

        {!isEditing ? (
          <div className="mt-4 hidden flex-wrap gap-3 md:flex">
            <button
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60"
              disabled={!canStudy}
              onClick={() => navigate(`/study/${deck.id}`)}
            >
              Study
            </button>
            <button
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
              onClick={() => navigate(`/notes/add?deckId=${deck.id}`)}
            >
              Add Note
            </button>
            <button
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
              onClick={() => {
                setActionError(null);
                setSettingsOpen(true);
              }}
            >
              Workload
            </button>
            <button
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
              onClick={() => {
                setActionError(null);
                setIsEditing(true);
              }}
            >
              Rename
            </button>
            <button
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 text-sm font-medium text-[var(--app-danger-text)] disabled:cursor-not-allowed disabled:opacity-45"
              disabled={!deck.canDelete || deleteMutation.isPending}
              onClick={() => {
                if (!window.confirm(`Delete the deck "${deck.name}"?`)) return;
                runDelete();
              }}
            >
              {deleteMutation.isPending ? "Deleting..." : "Delete"}
            </button>
          </div>
        ) : null}
      </li>

      <Sheet
        open={actionsOpen}
        onClose={() => setActionsOpen(false)}
        title={deck.name}
      >
        <div className="space-y-3">
          <button
            type="button"
            disabled={!canStudy}
            onClick={() => {
              navigate(`/study/${deck.id}`);
              setActionsOpen(false);
            }}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            Study
          </button>
          <button
            type="button"
            onClick={() => {
              navigate(`/notes/add?deckId=${deck.id}`);
              setActionsOpen(false);
            }}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Add note
          </button>
          <button
            type="button"
            onClick={() => {
              setActionError(null);
              setSettingsOpen(true);
              setActionsOpen(false);
            }}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Workload
          </button>
          <button
            type="button"
            onClick={() => {
              setActionError(null);
              setIsEditing(true);
              setActionsOpen(false);
            }}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
          >
            Rename
          </button>
          <button
            type="button"
            disabled={!deck.canDelete || deleteMutation.isPending}
            onClick={() => {
              setActionsOpen(false);
              setConfirmDeleteOpen(true);
            }}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 text-sm font-medium text-[var(--app-danger-text)] disabled:opacity-45"
          >
            {deleteMutation.isPending ? "Deleting..." : "Delete"}
          </button>
          {!deck.canDelete && deck.deleteBlockedReason ? (
            <p className="text-sm leading-6 text-[var(--app-text-soft)]">
              {deck.deleteBlockedReason}
            </p>
          ) : null}
        </div>
      </Sheet>

      <Sheet
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        title={`Deck workload: ${deck.name}`}
      >
        <form onSubmit={submitSettings} className="space-y-4">
          <p className="text-sm leading-6 text-[var(--app-text-soft)]">
            Tune how many new and review cards this deck can surface each day,
            and decide where it sits in your manual deck order.
          </p>
          <label className="block space-y-2">
            <span className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--app-muted)]">
              New cards / day
            </span>
            <input
              type="number"
              min={0}
              value={newCardsPerDay}
              onChange={(event) => setNewCardsPerDay(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--app-muted)]">
              Reviews / day
            </span>
            <input
              type="number"
              min={0}
              value={reviewsPerDay}
              onChange={(event) => setReviewsPerDay(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--app-muted)]">
              Priority
            </span>
            <input
              type="number"
              min={1}
              value={priorityOrder}
              onChange={(event) => setPriorityOrder(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            />
          </label>
          <div className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4 text-sm text-[var(--app-text-soft)]">
            {deck.newCardsPaused
              ? "New cards are currently paused for this deck because review backlog is above the review cap."
              : "New cards stay available while review backlog remains at or below the review cap."}
          </div>
          <div className="flex flex-col gap-3 sm:flex-row">
            <button
              type="submit"
              disabled={settingsMutation.isPending}
              className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
            >
              {settingsMutation.isPending ? "Saving..." : "Save workload"}
            </button>
            <button
              type="button"
              onClick={() => setSettingsOpen(false)}
              className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
            >
              Close
            </button>
          </div>
        </form>
      </Sheet>

      <ConfirmSheet
        open={confirmDeleteOpen}
        onClose={() => setConfirmDeleteOpen(false)}
        title={`Delete ${deck.name}?`}
        description="This only works for empty decks in the current tranche. Move or delete cards first if the action is disabled."
        confirmLabel="Delete deck"
        onConfirm={runDelete}
        destructive
      />
    </>
  );
}
