import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Navigate, useNavigate } from "react-router";
import {
  PageContainer,
  PageSection,
  SurfaceCard,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";
import type { UpdateWorkspacePlanRequest } from "#/lib/api";

const planOptions: Array<{
  plan: UpdateWorkspacePlanRequest["plan"];
  title: string;
  price: string;
  description: string;
  bullets: string[];
}> = [
  {
    plan: "free",
    title: "Free",
    price: "$0",
    description: "Tight by design for personal testing and simple study flows.",
    bullets: [
      "Small deck and note limits",
      "Solo workflow",
      "Great for getting started",
    ],
  },
  {
    plan: "pro",
    title: "Pro",
    price: "$12/mo",
    description: "For serious solo learners who want more room and AI support.",
    bullets: [
      "Higher limits",
      "AI-assisted note-to-card generation",
      "Marketplace publishing",
    ],
  },
  {
    plan: "team",
    title: "Team",
    price: "$8/user/mo",
    description:
      "Adds a real team workspace, member roles, and group workflows with a 3-seat minimum.",
    bullets: [
      "Creates a team-backed workspace",
      "Team roles and member management",
      "Study groups and shared publishing flows",
      "3 billed seats minimum at launch",
    ],
  },
  {
    plan: "enterprise",
    title: "Enterprise",
    price: "Custom",
    description: "For larger orgs that need admin control and custom agreements.",
    bullets: [
      "Enterprise limits",
      "Admin-friendly controls",
      "Custom support and rollout",
    ],
  },
];

export function OnboardingPlanPage() {
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const sessionQuery = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const selectPlanMutation = useMutation({
    mutationFn: async (plan: UpdateWorkspacePlanRequest["plan"]) => {
      if (plan === "enterprise") {
        throw new Error(
          "Enterprise onboarding is handled manually. Contact sales@vutadex.com to continue.",
        );
      }
      if (plan === "free") {
        const session = await repository.completeOnboardingPlan({ plan });
        return { mode: "session" as const, session };
      }
      const response = await repository.startBillingCheckout({ plan });
      return { mode: "billing" as const, response };
    },
    onSuccess: async (result) => {
      if (result.mode === "session") {
        queryClient.setQueryData(["auth-session"], result.session);
        await Promise.all([
          queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
          queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
          queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
          queryClient.invalidateQueries({ queryKey: ["decks"] }),
        ]);
        navigate("/", { replace: true });
        return;
      }

      if (result.response.completed && result.response.session) {
        queryClient.setQueryData(["auth-session"], result.response.session);
        await Promise.all([
          queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
          queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
          queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
          queryClient.invalidateQueries({ queryKey: ["decks"] }),
        ]);
        navigate("/", { replace: true });
        return;
      }

      if (result.response.checkoutUrl) {
        window.location.assign(result.response.checkoutUrl);
        return;
      }

      throw new Error("Billing checkout did not return a redirect URL.");
    },
  });

  if (sessionQuery.isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading onboarding...
        </PageSection>
      </PageContainer>
    );
  }

  if (!sessionQuery.data?.authenticated || !sessionQuery.data.user) {
    return <Navigate to="/login" replace />;
  }

  if (!sessionQuery.data.user.onboarding) {
    return <Navigate to="/" replace />;
  }

  const currentPlan = sessionQuery.data.entitlements.plan;

  return (
    <PageContainer className="space-y-6">
      <PageSection className="px-6 py-8 md:px-8 md:py-10">
        <p className="text-[11px] uppercase tracking-[0.28em] text-[var(--app-accent)]">
          Welcome to Vutadex
        </p>
        <h1 className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)] sm:text-4xl">
          Choose the plan that matches how you want to study.
        </h1>
        <p className="mt-4 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)] sm:text-base">
          We use this once right after sign-in so the current workspace starts
          with the right limits and collaboration features. You can change the
          plan later from User Settings at any time.
        </p>
      </PageSection>

      <div className="grid gap-4 xl:grid-cols-4">
        {planOptions.map((option) => {
          const isCurrent = currentPlan === option.plan;
          const isPending =
            selectPlanMutation.isPending &&
            selectPlanMutation.variables === option.plan;

          return (
            <SurfaceCard
              key={option.plan}
              className={[
                "flex h-full flex-col",
                isCurrent ? "border-[var(--app-accent)] bg-[var(--app-card-strong)]" : "",
              ].join(" ")}
            >
              <div className="flex items-start justify-between gap-3">
                <div>
                  <p className="text-sm font-semibold text-[var(--app-text)]">
                    {option.title}
                  </p>
                  <p className="mt-2 text-3xl font-semibold tracking-tight text-[var(--app-text)]">
                    {option.price}
                  </p>
                </div>
                {isCurrent ? (
                  <span className="rounded-full bg-[var(--app-accent)] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.16em] text-[var(--app-accent-ink)]">
                    Current
                  </span>
                ) : null}
              </div>

              <p className="mt-4 text-sm leading-6 text-[var(--app-text-soft)]">
                {option.description}
              </p>

              <ul className="mt-5 space-y-2 text-sm text-[var(--app-text-soft)]">
                {option.bullets.map((bullet) => (
                  <li key={bullet} className="flex gap-2">
                    <span className="mt-1 text-[var(--app-accent)]">•</span>
                    <span>{bullet}</span>
                  </li>
                ))}
              </ul>

              <div className="mt-auto pt-6">
                {option.plan === "enterprise" ? (
                  <a
                    href="mailto:sales@vutadex.com?subject=Vutadex%20Enterprise"
                    className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-center text-sm font-semibold text-[var(--app-text)] transition hover:border-[var(--app-accent)]"
                  >
                    Contact sales
                  </a>
                ) : (
                  <button
                    type="button"
                    onClick={() => selectPlanMutation.mutate(option.plan)}
                    disabled={selectPlanMutation.isPending}
                    className={[
                      "inline-flex min-h-11 w-full items-center justify-center rounded-2xl px-4 text-center text-sm font-semibold transition disabled:opacity-60",
                      isCurrent
                        ? "border border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text)]"
                        : "bg-[var(--app-accent)] text-[var(--app-accent-ink)] hover:brightness-105",
                    ].join(" ")}
                  >
                    {isPending
                      ? "Saving..."
                      : isCurrent
                        ? "Continue with this plan"
                        : option.plan === "free"
                          ? "Choose Free"
                          : `Start ${option.title}`}
                  </button>
                )}
              </div>
            </SurfaceCard>
          );
        })}
      </div>

      {selectPlanMutation.isError ? (
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {selectPlanMutation.error instanceof Error
            ? selectPlanMutation.error.message
            : "Failed to save your plan selection."}
        </PageSection>
      ) : null}
    </PageContainer>
  );
}
