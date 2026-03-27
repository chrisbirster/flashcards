import { Link } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ThemeToggle } from "#/components/theme-toggle";
import {
  EmptyState,
  PageContainer,
  PageSection,
  SurfaceCard,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";
import type { UpdateWorkspacePlanRequest } from "#/lib/api";

function SettingsRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="flex flex-col gap-1 rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4">
      <p className="text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
        {label}
      </p>
      <p className="text-sm font-medium text-[var(--app-text)]">{value}</p>
    </div>
  );
}

function formatDate(value?: string) {
  if (!value) return "Not set";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "Not set" : date.toLocaleDateString();
}

const planOptions: Array<{
  plan: UpdateWorkspacePlanRequest["plan"];
  label: string;
  description: string;
  price: string;
}> = [
  {
    plan: "free",
    label: "Free",
    description: "Personal starter limits and solo workflow.",
    price: "$0",
  },
  {
    plan: "pro",
    label: "Pro",
    description: "Higher limits and AI-assisted solo study.",
    price: "$12/mo",
  },
  {
    plan: "team",
    label: "Team",
    description: "Team roles, member installs, and shared admin surfaces.",
    price: "$8/user/mo",
  },
  {
    plan: "enterprise",
    label: "Enterprise",
    description: "Enterprise controls, limits, and support posture.",
    price: "Custom",
  },
];

export function SettingsPage() {
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const sessionQuery = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const logoutMutation = useMutation({
    mutationFn: () => repository.logout(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["auth-session"] });
      queryClient.invalidateQueries({ queryKey: ["entitlements"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });

  const updatePlanMutation = useMutation({
    mutationFn: async (plan: UpdateWorkspacePlanRequest["plan"]) => {
      const session = sessionQuery.data!;
      const currentPlan = session.entitlements.plan;
      const provider = session.subscription?.provider ?? "";

      if (plan === "enterprise") {
        window.location.assign(
          "mailto:sales@vutadex.com?subject=Vutadex%20Enterprise",
        );
        return { mode: "external" as const };
      }

      if (provider === "manual") {
        const nextSession = await repository.updateWorkspacePlan(
          session.workspace!.id,
          { plan },
        );
        return { mode: "session" as const, session: nextSession };
      }

      if (currentPlan === "free" || currentPlan === "guest") {
        const response = await repository.startBillingCheckout({ plan });
        return { mode: "checkout" as const, response };
      }

      const response = await repository.openBillingPortal({ plan });
      return { mode: "portal" as const, response };
    },
    onSuccess: async (result) => {
      if (result.mode === "session") {
        queryClient.setQueryData(["auth-session"], result.session);
        await Promise.all([
          queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
          queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
          queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
        ]);
        return;
      }

      if (result.mode === "checkout") {
        if (result.response.completed && result.response.session) {
          queryClient.setQueryData(["auth-session"], result.response.session);
          await Promise.all([
            queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
            queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
            queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
          ]);
          return;
        }
        if (result.response.checkoutUrl) {
          window.location.assign(result.response.checkoutUrl);
          return;
        }
        throw new Error("Billing checkout did not return a redirect URL.");
      }

      if (result.mode === "portal") {
        window.location.assign(result.response.url);
      }
    },
  });

  const billingPortalMutation = useMutation({
    mutationFn: () => repository.openBillingPortal(),
    onSuccess: (result) => {
      window.location.assign(result.url);
    },
  });

  if (sessionQuery.isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading settings...
        </PageSection>
      </PageContainer>
    );
  }

  if (!sessionQuery.data?.authenticated || !sessionQuery.data.user) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="Settings unavailable"
          description="Sign in to manage your account, workspace plan, and team context."
        />
      </PageContainer>
    );
  }

  const session = sessionQuery.data!;
  const user = session.user!;
  const userLabel = user.displayName || user.email;
  const userInitial = userLabel.trim().charAt(0).toUpperCase() || "U";
  const currentPlan = session.entitlements.plan;
  const subscription = session.subscription;
  const organizationRole = session.organizationMember?.role;
  const canManagePlan = session.workspace?.organizationId
    ? organizationRole === "admin" || organizationRole === "owner"
    : session.workspace?.ownerUserId === user.id;
  const canOpenBillingPortal =
    canManagePlan &&
    subscription?.provider === "stripe" &&
    Boolean(subscription?.providerCustomerId);

  return (
    <PageContainer className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)]">
        <div className="rounded-[2rem] bg-[var(--app-card-strong)] px-6 py-8 text-[var(--app-text)] shadow-sm md:px-8 md:py-10">
          <p className="text-xs uppercase tracking-[0.3em] text-[var(--app-accent)]">
            User settings
          </p>
          <h2 className="mt-4 text-4xl font-semibold tracking-tight md:text-5xl">
            Manage your profile, plan, and workspace context.
          </h2>
          <p className="mt-4 max-w-2xl text-sm leading-7 text-[var(--app-text-soft)] md:text-base">
            This is the account surface for the current signed-in user. Plan
            changes live here, and team-backed workspaces link directly into
            member and role management from this page.
          </p>
        </div>

        <PageSection className="p-6">
          <div className="flex items-center gap-4">
            <span className="flex h-16 w-16 items-center justify-center rounded-full bg-[var(--app-accent)] text-xl font-semibold text-[var(--app-accent-ink)]">
              {userInitial}
            </span>
            <div className="min-w-0">
              <p className="truncate text-lg font-semibold text-[var(--app-text)]">
                {userLabel}
              </p>
              <p className="mt-1 text-sm text-[var(--app-text-soft)]">
                {user.email}
              </p>
              <p className="mt-2 text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">
                {currentPlan.toUpperCase()} plan
              </p>
            </div>
          </div>
        </PageSection>
      </section>

      <section className="grid gap-6 lg:grid-cols-2">
        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Profile
          </h3>
          <div className="mt-5 grid gap-3">
            <SettingsRow label="Display name" value={user.displayName} />
            <SettingsRow label="Email" value={user.email} />
            <SettingsRow
              label="Current plan"
              value={currentPlan.toUpperCase()}
            />
          </div>
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Workspace
          </h3>
          <div className="mt-5 grid gap-3">
            <SettingsRow
              label="Workspace name"
              value={session.workspace?.name || "Unknown workspace"}
            />
            {session.organization ? (
              <SettingsRow label="Team" value={session.organization.name} />
            ) : (
              <SettingsRow label="Workspace type" value="Personal" />
            )}
            {session.organizationMember ? (
              <SettingsRow
                label="Your team role"
                value={session.organizationMember.role}
              />
            ) : null}
          </div>
        </PageSection>
      </section>

      <PageSection className="p-5 sm:p-6">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
              Plan management
            </h3>
            <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
              Choose the plan that fits the current workspace. Team-backed
              workspaces can only be changed by team admins and owners.
            </p>
          </div>
          {!canManagePlan ? (
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.16em] text-[var(--app-muted)]">
              Read only
            </span>
          ) : null}
        </div>

        <div className="mt-5 grid gap-3 lg:grid-cols-4">
          {planOptions.map((option) => {
            const isCurrent = currentPlan === option.plan;
            const isPending =
              updatePlanMutation.isPending &&
              updatePlanMutation.variables === option.plan;

            return (
              <SurfaceCard
                key={option.plan}
                className={[
                  "flex h-full flex-col border-none bg-[var(--app-card-strong)] p-4",
                  isCurrent ? "ring-1 ring-[var(--app-accent)]" : "",
                ].join(" ")}
              >
                <p className="text-sm font-semibold text-[var(--app-text)]">
                  {option.label}
                </p>
                <p className="mt-2 text-2xl font-semibold tracking-tight text-[var(--app-text)]">
                  {option.price}
                </p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  {option.description}
                </p>
                {option.plan === "team" ? (
                  <p className="mt-2 text-xs text-[var(--app-muted)]">
                    3-seat minimum. Seats are billed for active team members.
                  </p>
                ) : null}
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
                      onClick={() => updatePlanMutation.mutate(option.plan)}
                      disabled={
                        isCurrent ||
                        !canManagePlan ||
                        updatePlanMutation.isPending ||
                        billingPortalMutation.isPending
                      }
                      className={[
                        "inline-flex min-h-11 w-full items-center justify-center rounded-2xl px-4 text-center text-sm font-semibold transition disabled:opacity-60",
                        isCurrent
                          ? "border border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text)]"
                          : "bg-[var(--app-accent)] text-[var(--app-accent-ink)] hover:brightness-105",
                      ].join(" ")}
                    >
                      {isPending
                        ? "Updating..."
                        : isCurrent
                          ? "Current plan"
                          : currentPlan === "free" || currentPlan === "guest"
                            ? option.plan === "team"
                              ? "Start Team checkout"
                              : `Start ${option.label}`
                            : `Manage ${option.label} in billing`}
                    </button>
                  )}
                </div>
              </SurfaceCard>
            );
          })}
        </div>

        <div className="mt-5 grid gap-3 md:grid-cols-3">
          <SettingsRow
            label="Billing provider"
            value={subscription?.provider?.toUpperCase() || "LOCAL"}
          />
          <SettingsRow
            label="Current period end"
            value={formatDate(subscription?.currentPeriodEnd)}
          />
          <SettingsRow
            label="Subscription status"
            value={
              subscription?.cancelAtPeriodEnd
                ? `Scheduled to end on ${formatDate(subscription.currentPeriodEnd)}`
                : subscription?.status || "No paid subscription"
            }
          />
        </div>

        {canOpenBillingPortal ? (
          <div className="mt-5">
            <button
              type="button"
              onClick={() => billingPortalMutation.mutate()}
              disabled={billingPortalMutation.isPending || updatePlanMutation.isPending}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-5 text-sm font-medium text-[var(--app-text)] transition hover:border-[var(--app-accent)] hover:bg-[var(--app-muted-surface)] disabled:opacity-60"
            >
              {billingPortalMutation.isPending
                ? "Opening billing..."
                : "Manage billing, invoices, and payment methods"}
            </button>
          </div>
        ) : null}

        {!canManagePlan ? (
          <p className="mt-4 text-sm text-[var(--app-muted)]">
            {session.workspace?.organizationId
              ? "Only team admins and owners can change the current team plan."
              : "Only the workspace owner can change the current plan."}
          </p>
        ) : null}

        {(updatePlanMutation.isError || billingPortalMutation.isError) ? (
          <p className="mt-4 text-sm text-[var(--app-danger-text)]">
            {(
              updatePlanMutation.error || billingPortalMutation.error
            ) instanceof Error
              ? (
                  updatePlanMutation.error || billingPortalMutation.error
                )!.message
              : "Failed to update plan."}
          </p>
        ) : null}
      </PageSection>

      {session.organization ? (
        <PageSection className="p-5 sm:p-6">
          <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
            <div>
              <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
                Team
              </h3>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Open the team view to manage members, roles, and team-level
                plan controls for this workspace.
              </p>
            </div>
            <Link
              to="/team"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Open team
            </Link>
          </div>
        </PageSection>
      ) : null}

      <section className="grid gap-6 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Appearance
          </h3>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Toggle between the Vutadex dark and light themes.
          </p>
          <div className="mt-5">
            <ThemeToggle />
          </div>
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <h3 className="text-xl font-semibold tracking-tight text-[var(--app-text)]">
            Account actions
          </h3>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Use sign out if you want to switch accounts or leave the current
            workspace session on this device.
          </p>
          <button
            type="button"
            onClick={() => logoutMutation.mutate()}
            disabled={logoutMutation.isPending}
            className="mt-5 inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-5 text-sm font-medium text-[var(--app-text)] transition hover:border-[var(--app-accent)] hover:bg-[var(--app-muted-surface)] disabled:opacity-60"
          >
            {logoutMutation.isPending ? "Signing out..." : "Sign out"}
          </button>
        </PageSection>
      </section>
    </PageContainer>
  );
}
