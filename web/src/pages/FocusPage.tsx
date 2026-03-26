import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Link } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  EmptyState,
  PageContainer,
  PageSection,
  StatCard,
  SurfaceCard,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";

type FocusProtocol = "pomodoro" | "deep-focus" | "custom";

type FocusPreset = {
  id: FocusProtocol;
  label: string;
  description: string;
  targetMinutes: number;
  breakMinutes: number;
};

const FOCUS_PRESETS: FocusPreset[] = [
  {
    id: "pomodoro",
    label: "Pomodoro",
    description: "25 minutes of focus with a short 5 minute reset after.",
    targetMinutes: 25,
    breakMinutes: 5,
  },
  {
    id: "deep-focus",
    label: "Deep focus",
    description: "50 minutes for heavier work with a 10 minute recovery break.",
    targetMinutes: 50,
    breakMinutes: 10,
  },
  {
    id: "custom",
    label: "Custom",
    description: "Start from a 30 minute block and tune it for the day.",
    targetMinutes: 30,
    breakMinutes: 10,
  },
];

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

function formatCountdown(totalSeconds: number) {
  const safe = Math.max(0, totalSeconds);
  const minutes = Math.floor(safe / 60)
    .toString()
    .padStart(2, "0");
  const seconds = (safe % 60).toString().padStart(2, "0");
  return `${minutes}:${seconds}`;
}

export function FocusPage() {
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const completionInFlightRef = useRef(false);

  const analyticsQuery = useQuery({
    queryKey: ["study-analytics"],
    queryFn: () => repository.fetchStudyAnalyticsOverview(),
  });

  const [customTargetMinutes, setCustomTargetMinutes] = useState("30");
  const [customBreakMinutes, setCustomBreakMinutes] = useState("10");
  const [selectedPresetId, setSelectedPresetId] =
    useState<FocusProtocol>("pomodoro");
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [activePreset, setActivePreset] = useState<FocusPreset | null>(null);
  const [timerPhase, setTimerPhase] = useState<
    "idle" | "running" | "paused" | "completed"
  >("idle");
  const [remainingSeconds, setRemainingSeconds] = useState(0);
  const [deadlineAt, setDeadlineAt] = useState<number | null>(null);
  const [lastOutcome, setLastOutcome] = useState<
    "completed" | "abandoned" | null
  >(null);

  const createFocusSessionMutation = useMutation({
    mutationFn: ({
      protocol,
      targetMinutes,
      breakMinutes,
    }: {
      protocol: FocusProtocol;
      targetMinutes: number;
      breakMinutes: number;
    }) =>
      repository.createStudySession({
        mode: "focus",
        protocol,
        targetMinutes,
        breakMinutes,
      }),
  });

  const updateStudySessionMutation = useMutation({
    mutationFn: ({
      id,
      status,
    }: {
      id: string;
      status: "completed" | "abandoned";
    }) =>
      repository.updateStudySession(id, {
        status,
        endedAt: new Date().toISOString(),
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["study-analytics"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
      queryClient.invalidateQueries({ queryKey: ["decks"] });
    },
  });

  const selectedPreset = useMemo(() => {
    if (selectedPresetId !== "custom") {
      return FOCUS_PRESETS.find((preset) => preset.id === selectedPresetId)!;
    }

    const parsedTarget = Number.parseInt(customTargetMinutes, 10);
    const parsedBreak = Number.parseInt(customBreakMinutes, 10);
    return {
      id: "custom",
      label: "Custom",
      description: "A flexible focus block you can tune whenever the day shifts.",
      targetMinutes:
        Number.isFinite(parsedTarget) && parsedTarget > 0 ? parsedTarget : 30,
      breakMinutes:
        Number.isFinite(parsedBreak) && parsedBreak >= 0 ? parsedBreak : 10,
    } satisfies FocusPreset;
  }, [customBreakMinutes, customTargetMinutes, selectedPresetId]);

  const resetFocusState = useCallback(() => {
    setActiveSessionId(null);
    setActivePreset(null);
    setTimerPhase("idle");
    setRemainingSeconds(0);
    setDeadlineAt(null);
    completionInFlightRef.current = false;
  }, []);

  const completeFocusSession = useCallback(async () => {
    if (!activeSessionId || completionInFlightRef.current) {
      return;
    }
    completionInFlightRef.current = true;
    try {
      await updateStudySessionMutation.mutateAsync({
        id: activeSessionId,
        status: "completed",
      });
      setTimerPhase("completed");
      setDeadlineAt(null);
      setRemainingSeconds(0);
      setLastOutcome("completed");
    } finally {
      completionInFlightRef.current = false;
    }
  }, [activeSessionId, updateStudySessionMutation]);

  useEffect(() => {
    if (timerPhase !== "running" || deadlineAt == null) {
      return;
    }

    const tick = () => {
      const nextRemainingSeconds = Math.max(
        0,
        Math.ceil((deadlineAt - Date.now()) / 1000),
      );
      setRemainingSeconds(nextRemainingSeconds);
      if (nextRemainingSeconds === 0) {
        void completeFocusSession();
      }
    };

    tick();
    const intervalId = window.setInterval(tick, 1000);
    return () => window.clearInterval(intervalId);
  }, [completeFocusSession, deadlineAt, timerPhase]);

  const startFocusPreset = async () => {
    if (createFocusSessionMutation.isPending || updateStudySessionMutation.isPending) {
      return;
    }

    const created = await createFocusSessionMutation.mutateAsync({
      protocol: selectedPreset.id,
      targetMinutes: selectedPreset.targetMinutes,
      breakMinutes: selectedPreset.breakMinutes,
    });

    setLastOutcome(null);
    setActiveSessionId(created.id);
    setActivePreset(selectedPreset);
    setTimerPhase("running");
    setRemainingSeconds(selectedPreset.targetMinutes * 60);
    setDeadlineAt(Date.now() + selectedPreset.targetMinutes * 60 * 1000);
  };

  const pauseFocusTimer = () => {
    if (timerPhase !== "running" || deadlineAt == null) {
      return;
    }
    setRemainingSeconds(Math.max(0, Math.ceil((deadlineAt - Date.now()) / 1000)));
    setDeadlineAt(null);
    setTimerPhase("paused");
  };

  const resumeFocusTimer = () => {
    if (timerPhase !== "paused" || remainingSeconds <= 0) {
      return;
    }
    setDeadlineAt(Date.now() + remainingSeconds * 1000);
    setTimerPhase("running");
  };

  const abandonFocusSession = async () => {
    if (timerPhase === "completed") {
      resetFocusState();
      return;
    }
    if (!activeSessionId || updateStudySessionMutation.isPending) {
      resetFocusState();
      setLastOutcome("abandoned");
      return;
    }

    await updateStudySessionMutation.mutateAsync({
      id: activeSessionId,
      status: "abandoned",
    });
    resetFocusState();
    setLastOutcome("abandoned");
  };

  const analytics = analyticsQuery.data;

  return (
    <PageContainer className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.15fr)_minmax(0,0.85fr)]">
        <div className="rounded-[2rem] bg-[var(--app-card-strong)] px-6 py-8 text-[var(--app-text)] shadow-sm md:px-8 md:py-10">
          <p className="text-xs uppercase tracking-[0.3em] text-[var(--app-accent)]">
            Focus sessions
          </p>
          <h2 className="mt-4 max-w-3xl text-4xl font-semibold tracking-tight md:text-5xl">
            Give yourself a clean study block before the day gets noisy.
          </h2>
          <p className="mt-4 max-w-2xl text-sm leading-7 text-[var(--app-text-soft)] md:text-base">
            Pomodoro-style sessions now live alongside review analytics, so we can
            track deep work separately from card volume without changing FSRS.
          </p>
          <div className="mt-8 flex flex-wrap gap-3">
            <Link
              to="/stats"
              className="inline-flex items-center rounded-2xl bg-[var(--app-accent)] px-4 py-2.5 text-sm font-medium text-[var(--app-accent-ink)] transition hover:brightness-105"
            >
              Open stats
            </Link>
            <Link
              to="/decks"
              className="inline-flex items-center rounded-2xl border border-[var(--app-line-strong)] px-4 py-2.5 text-sm font-medium text-[var(--app-text)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)]"
            >
              Study decks
            </Link>
          </div>
        </div>

        <PageSection className="p-6">
          <p className="text-xs uppercase tracking-[0.24em] text-[var(--app-muted)]">
            This week
          </p>
          <div className="mt-5 grid gap-4 sm:grid-cols-2">
            <StatCard
              label="Focus blocks (7d)"
              value={analytics?.focusSessions7d ?? 0}
              detail="Completed timed focus sessions tracked this week."
            />
            <StatCard
              label="Focus minutes (7d)"
              value={formatMinutes(analytics?.focusMinutes7d ?? 0)}
              detail="Time spent inside completed focus blocks."
            />
          </div>
          <p className="mt-5 text-sm leading-6 text-[var(--app-text-soft)]">
            Focus sessions count separately from card-review sessions, so we can
            see whether the issue is workload, attention, or both.
          </p>
        </PageSection>
      </section>

      <section className="grid gap-6 xl:grid-cols-[minmax(0,1.1fr)_minmax(0,0.9fr)]">
        <PageSection className="p-5 sm:p-6">
          <div className="flex items-center justify-between gap-3">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Pick a focus block
              </h3>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                Start with a preset, then adjust the custom block when your day needs
                a different rhythm.
              </p>
            </div>
          </div>

          <div className="mt-6 grid gap-4 md:grid-cols-3">
            {FOCUS_PRESETS.map((preset) => {
              const isSelected = selectedPresetId === preset.id;
              return (
                <button
                  key={preset.id}
                  type="button"
                  onClick={() => setSelectedPresetId(preset.id)}
                  className={[
                    "rounded-[1.5rem] border p-4 text-left transition",
                    isSelected
                      ? "border-[var(--app-accent)] bg-[var(--app-card-strong)]"
                      : "border-[var(--app-line)] bg-[var(--app-muted-surface)] hover:border-[var(--app-line-strong)]",
                  ].join(" ")}
                >
                  <p className="text-sm font-semibold text-[var(--app-text)]">
                    {preset.label}
                  </p>
                  <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                    {preset.description}
                  </p>
                  <div className="mt-4 flex flex-wrap gap-2">
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {preset.targetMinutes}m focus
                    </span>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]">
                      {preset.breakMinutes}m break
                    </span>
                  </div>
                </button>
              );
            })}
          </div>

          <SurfaceCard className="mt-4 bg-[var(--app-muted-surface)] p-4">
            <div className="flex flex-col gap-4 sm:flex-row">
              <label className="flex-1 space-y-2">
                <span className="text-xs uppercase tracking-[0.2em] text-[var(--app-muted)]">
                  Custom focus minutes
                </span>
                <input
                  type="number"
                  min={1}
                  value={customTargetMinutes}
                  onChange={(event) => {
                    setSelectedPresetId("custom");
                    setCustomTargetMinutes(event.target.value);
                  }}
                  className="min-h-11 w-full rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] px-4 text-sm text-[var(--app-text)] outline-none transition focus:border-[var(--app-accent)]"
                />
              </label>
              <label className="flex-1 space-y-2">
                <span className="text-xs uppercase tracking-[0.2em] text-[var(--app-muted)]">
                  Custom break minutes
                </span>
                <input
                  type="number"
                  min={0}
                  value={customBreakMinutes}
                  onChange={(event) => {
                    setSelectedPresetId("custom");
                    setCustomBreakMinutes(event.target.value);
                  }}
                  className="min-h-11 w-full rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] px-4 text-sm text-[var(--app-text)] outline-none transition focus:border-[var(--app-accent)]"
                />
              </label>
            </div>
          </SurfaceCard>

          <div className="mt-4 flex flex-wrap gap-3">
            <button
              type="button"
              onClick={() => void startFocusPreset()}
              disabled={timerPhase === "running" || createFocusSessionMutation.isPending}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
            >
              {timerPhase === "running" ? "Focus block running" : `Start ${selectedPreset.label}`}
            </button>
            <span className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] px-4 text-sm text-[var(--app-text-soft)]">
              {selectedPreset.targetMinutes}m focus + {selectedPreset.breakMinutes}m break
            </span>
          </div>
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Current timer
          </h3>
          <p className="mt-1 text-sm text-[var(--app-text-soft)]">
            Use timed focus when review volume is fine but attention is the real bottleneck.
          </p>

          {activePreset ? (
            <div className="mt-6 space-y-4">
              <div className="rounded-[1.75rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] px-5 py-6 text-center">
                <p className="text-xs uppercase tracking-[0.26em] text-[var(--app-accent)]">
                  {activePreset.label}
                </p>
                <p className="mt-4 text-5xl font-semibold tracking-tight text-[var(--app-text)] md:text-6xl">
                  {formatCountdown(remainingSeconds)}
                </p>
                <p className="mt-3 text-sm text-[var(--app-text-soft)]">
                  {timerPhase === "completed"
                    ? `Focus block complete. Take a ${activePreset.breakMinutes} minute break before the next push.`
                    : timerPhase === "paused"
                      ? "Paused. Resume when you are ready to keep going."
                      : `${activePreset.targetMinutes} minute focus block with a ${activePreset.breakMinutes} minute break recommendation after.`}
                </p>
              </div>

              <div className="grid gap-3 sm:grid-cols-2">
                {timerPhase === "running" ? (
                  <button
                    type="button"
                    onClick={pauseFocusTimer}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
                  >
                    Pause timer
                  </button>
                ) : null}
                {timerPhase === "paused" ? (
                  <button
                    type="button"
                    onClick={resumeFocusTimer}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
                  >
                    Resume timer
                  </button>
                ) : null}
                {timerPhase !== "completed" ? (
                  <button
                    type="button"
                    onClick={() => void completeFocusSession()}
                    disabled={updateStudySessionMutation.isPending}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)] disabled:cursor-not-allowed disabled:opacity-60"
                  >
                    Complete focus block
                  </button>
                ) : null}
                <button
                  type="button"
                  onClick={() => void abandonFocusSession()}
                  disabled={updateStudySessionMutation.isPending}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-4 text-sm font-medium text-[var(--app-text-soft)] disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {timerPhase === "completed" ? "Clear timer" : "Abandon session"}
                </button>
              </div>
            </div>
          ) : (
            <div className="mt-6">
              <EmptyState
                title="No active focus block"
                description="Start a timed focus session from the presets on the left and we will capture it alongside your study analytics."
              />
            </div>
          )}

          {lastOutcome === "abandoned" ? (
            <p className="mt-4 text-sm text-[var(--app-text-soft)]">
              Focus session marked as abandoned. Start another block whenever you are ready.
            </p>
          ) : null}
        </PageSection>
      </section>
    </PageContainer>
  );
}
