import { useEffect } from "react";
import { Link, useNavigate, useSearchParams } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  EmptyState,
  PageContainer,
  PageSection,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";

export function BillingCompletePage() {
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const checkoutState = searchParams.get("checkout");
  const checkoutSessionId = searchParams.get("checkout_session_id")?.trim() ?? "";

  const sessionQuery = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const syncMutation = useMutation({
    mutationFn: (sessionId: string) => repository.syncBillingCheckoutSession(sessionId),
    onSuccess: async (payload) => {
      if (payload.session) {
        queryClient.setQueryData(["auth-session"], payload.session);
      }
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
        queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
        queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
        queryClient.invalidateQueries({ queryKey: ["organization"] }),
        queryClient.invalidateQueries({ queryKey: ["study-groups"] }),
      ]);

      const nextSession = payload.session;
      if (nextSession?.user?.onboarding) {
        return;
      }

      const fromOnboarding = sessionQuery.data?.user?.onboarding ?? false;
      navigate(fromOnboarding ? "/" : "/settings?billing=success", {
        replace: true,
      });
    },
  });

  const retrySync = () => {
    if (!checkoutSessionId || syncMutation.isPending) return;
    syncMutation.mutate(checkoutSessionId);
  };

  useEffect(() => {
    if (checkoutState === "success" && checkoutSessionId && syncMutation.isIdle) {
      syncMutation.mutate(checkoutSessionId);
    }
  }, [checkoutSessionId, checkoutState, syncMutation]);

  if (sessionQuery.isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading billing status...
        </PageSection>
      </PageContainer>
    );
  }

  if (!sessionQuery.data?.authenticated) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="Sign in required"
          description="Sign in again to finish syncing the billing change for this workspace."
          action={
            <Link
              to="/login"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Go to login
            </Link>
          }
        />
      </PageContainer>
    );
  }

  if (checkoutState === "cancelled") {
    const returnPath = sessionQuery.data.user?.onboarding ? "/onboarding/plan" : "/settings";
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="Checkout canceled"
          description="No billing change was applied. You can return to plan management whenever you're ready."
          action={
            <Link
              to={returnPath}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Return to plan management
            </Link>
          }
        />
      </PageContainer>
    );
  }

  if (checkoutState !== "success" || !checkoutSessionId) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="Billing return link is incomplete"
          description="We couldn't find a Stripe checkout session to sync. Head back to settings and try again."
          action={
            <Link
              to="/settings"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Open settings
            </Link>
          }
        />
      </PageContainer>
    );
  }

  const hasSucceeded = syncMutation.isSuccess && syncMutation.data.completed;
  const needsRetry = syncMutation.isSuccess && !syncMutation.data.completed;

  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-6 sm:p-8">
        <p className="text-[11px] uppercase tracking-[0.28em] text-[var(--app-accent)]">
          Billing
        </p>
        <h1 className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)] sm:text-4xl">
          {hasSucceeded ? "Billing updated" : "Syncing your billing change"}
        </h1>
        <p className="mt-4 max-w-2xl text-sm leading-7 text-[var(--app-text-soft)]">
          {hasSucceeded
            ? "Your workspace subscription is now in sync with Stripe. We’re refreshing the app state so the right entitlements show up everywhere."
            : "Stripe has returned control to Vutadex. We just need to sync the checkout result back into your workspace before we continue."}
        </p>

        <div className="mt-6 flex flex-wrap gap-3">
          {!hasSucceeded ? (
            <button
              type="button"
              onClick={retrySync}
              disabled={syncMutation.isPending}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
            >
              {syncMutation.isPending ? "Syncing..." : "Sync billing now"}
            </button>
          ) : null}
          <Link
            to={sessionQuery.data.user?.onboarding ? "/onboarding/plan" : "/settings"}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-5 text-sm font-medium text-[var(--app-text)]"
          >
            Back to app
          </Link>
        </div>

        {needsRetry ? (
          <p className="mt-4 text-sm text-[var(--app-text-soft)]">
            Stripe checkout has not fully completed yet. Give it a moment and
            run the sync again.
          </p>
        ) : null}

        {syncMutation.isError ? (
          <p className="mt-4 text-sm text-[var(--app-danger-text)]">
            {syncMutation.error instanceof Error
              ? syncMutation.error.message
              : "Failed to sync the billing checkout."}
          </p>
        ) : null}
      </PageSection>
    </PageContainer>
  );
}
