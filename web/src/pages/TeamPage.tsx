import { useEffect, useMemo, useState } from "react";
import { Link, Navigate, useNavigate } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ConfirmSheet, Sheet } from "#/components/sheet";
import {
  EmptyState,
  PageContainer,
  PageSection,
  SurfaceCard,
} from "#/components/page-layout";
import { useAppRepository } from "#/lib/app-repository";
import type {
  OrganizationMember,
  UpdateWorkspacePlanRequest,
} from "#/lib/api";

const planOptions: Array<{
  plan: UpdateWorkspacePlanRequest["plan"];
  label: string;
  detail: string;
}> = [
  { plan: "free", label: "Free", detail: "Solo starter limits" },
  { plan: "pro", label: "Pro", detail: "Larger solo workflow" },
  { plan: "team", label: "Team", detail: "Roles and team workspace" },
  { plan: "enterprise", label: "Enterprise", detail: "Enterprise controls" },
];

const memberRoleOptions: Array<OrganizationMember["role"]> = [
  "read",
  "edit",
  "admin",
  "owner",
];

function formatDateTime(value?: string) {
  if (!value) return "Unknown";
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? "Unknown" : date.toLocaleString();
}

export function TeamPage() {
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [inviteOpen, setInviteOpen] = useState(false);
  const [deleteOpen, setDeleteOpen] = useState(false);
  const [teamName, setTeamName] = useState("");
  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteRole, setInviteRole] =
    useState<Exclude<OrganizationMember["role"], "owner">>("read");

  const sessionQuery = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const orgId =
    sessionQuery.data?.organization?.id ??
    sessionQuery.data?.workspace?.organizationId;

  const teamQuery = useQuery({
    queryKey: ["organization", orgId],
    queryFn: () => repository.fetchOrganization(orgId!),
    enabled: Boolean(orgId),
  });

  const team = teamQuery.data;

  useEffect(() => {
    if (team?.organization.name) {
      setTeamName(team.organization.name);
    }
  }, [team?.organization.name]);

  const refresh = async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: ["auth-session"] }),
      queryClient.invalidateQueries({ queryKey: ["organization", orgId] }),
      queryClient.invalidateQueries({ queryKey: ["dashboard"] }),
      queryClient.invalidateQueries({ queryKey: ["entitlements"] }),
      queryClient.invalidateQueries({ queryKey: ["study-groups"] }),
    ]);
  };

  const updateTeamMutation = useMutation({
    mutationFn: () =>
      repository.updateOrganization(orgId!, {
        name: teamName,
        slug: team?.organization.slug,
      }),
    onSuccess: async () => {
      await refresh();
    },
  });

  const inviteMutation = useMutation({
    mutationFn: () =>
      repository.addOrganizationMember(orgId!, {
        email: inviteEmail,
        role: inviteRole,
      }),
    onSuccess: async () => {
      setInviteEmail("");
      setInviteRole("read");
      setInviteOpen(false);
      await refresh();
    },
  });

  const updateMemberMutation = useMutation({
    mutationFn: ({
      memberId,
      role,
    }: {
      memberId: string;
      role: OrganizationMember["role"];
    }) => repository.updateOrganizationMember(orgId!, memberId, { role }),
    onSuccess: async () => {
      await refresh();
    },
  });

  const removeMemberMutation = useMutation({
    mutationFn: (memberId: string) =>
      repository.deleteOrganizationMember(orgId!, memberId),
    onSuccess: async () => {
      await refresh();
    },
  });

  const updatePlanMutation = useMutation({
    mutationFn: (plan: UpdateWorkspacePlanRequest["plan"]) =>
      repository.updateWorkspacePlan(sessionQuery.data!.workspace!.id, { plan }),
    onSuccess: async (session) => {
      queryClient.setQueryData(["auth-session"], session);
      await refresh();
    },
  });

  const deleteTeamMutation = useMutation({
    mutationFn: () => repository.deleteOrganization(orgId!),
    onSuccess: async () => {
      await refresh();
      navigate("/settings", { replace: true });
    },
  });

  const currentRole = team?.membership.role;
  const canManageMembers = team?.canManageMembers ?? false;
  const canManagePlan = team?.canManagePlan ?? false;
  const canEditMetadata = currentRole === "admin" || currentRole === "owner";
  const canDeleteTeam = currentRole === "owner";

  const currentPlan =
    sessionQuery.data?.entitlements.plan ??
    team?.subscription?.plan ??
    "free";

  const activeMembers = useMemo(
    () => (team?.members ?? []).filter((member) => member.status === "active"),
    [team?.members],
  );

  if (sessionQuery.isLoading || (orgId && teamQuery.isLoading)) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">
          Loading team...
        </PageSection>
      </PageContainer>
    );
  }

  if (!sessionQuery.data?.authenticated) {
    return <Navigate to="/login" replace />;
  }

  if (!orgId) {
    return (
      <PageContainer className="space-y-4">
        <EmptyState
          title="No team attached to this workspace"
          description="Upgrade the current workspace to Team from User Settings when you want member roles, team billing, and group administration."
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

  if (teamQuery.isError || !team) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {teamQuery.error instanceof Error
            ? teamQuery.error.message
            : "Failed to load team."}
        </PageSection>
      </PageContainer>
    );
  }

  return (
    <PageContainer className="space-y-6">
      <section className="grid gap-6 lg:grid-cols-[minmax(0,1.2fr)_minmax(0,0.8fr)]">
        <PageSection className="px-6 py-8 md:px-8 md:py-10">
          <p className="text-[11px] uppercase tracking-[0.28em] text-[var(--app-accent)]">
            Team
          </p>
          <h1 className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)] sm:text-4xl">
            {team.organization.name}
          </h1>
          <p className="mt-4 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)] sm:text-base">
            This is the admin surface for the current team-backed workspace.
            Team roles gate content editing, member management, and plan
            control without exposing internal workspace storage details.
          </p>
          <div className="mt-6 flex flex-wrap gap-2">
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
              {currentRole}
            </span>
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
              {activeMembers.length} active members
            </span>
          </div>
        </PageSection>

        <PageSection className="p-6">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">
            Current workspace
          </p>
          <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
            {team.workspace?.name || "Team workspace"}
          </p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Plan changes here affect the active team workspace. Study Groups and
            Marketplace installs still create workspace-local copies with
            private review history.
          </p>
        </PageSection>
      </section>

      <section className="grid gap-6 xl:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
        <PageSection className="p-5 sm:p-6">
          <div className="flex items-start justify-between gap-3">
            <div>
              <h2 className="text-xl font-semibold text-[var(--app-text)]">
                Team settings
              </h2>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Team admins and owners can update team metadata. Owners can also
                delete the team.
              </p>
            </div>
          </div>

          <div className="mt-5 space-y-3">
            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--app-text)]">
                Team name
              </span>
              <input
                value={teamName}
                onChange={(event) => setTeamName(event.target.value)}
                disabled={!canEditMetadata || updateTeamMutation.isPending}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)] disabled:opacity-60"
              />
            </label>
            <div className="rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-muted-surface)] p-4">
              <p className="text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
                Created
              </p>
              <p className="mt-2 text-sm font-medium text-[var(--app-text)]">
                {formatDateTime(team.organization.createdAt)}
              </p>
            </div>
          </div>

          <div className="mt-5 flex flex-wrap gap-3">
            <button
              type="button"
              onClick={() => updateTeamMutation.mutate()}
              disabled={
                !canEditMetadata ||
                updateTeamMutation.isPending ||
                !teamName.trim() ||
                teamName.trim() === team.organization.name
              }
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
            >
              {updateTeamMutation.isPending ? "Saving..." : "Save changes"}
            </button>
            {canDeleteTeam ? (
              <button
                type="button"
                onClick={() => setDeleteOpen(true)}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-5 text-sm font-semibold text-[var(--app-danger-text)]"
              >
                Delete team
              </button>
            ) : null}
          </div>

          {updateTeamMutation.isError ? (
            <p className="mt-4 text-sm text-[var(--app-danger-text)]">
              {updateTeamMutation.error instanceof Error
                ? updateTeamMutation.error.message
                : "Failed to update team settings."}
            </p>
          ) : null}
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <h2 className="text-xl font-semibold text-[var(--app-text)]">
            Team plan
          </h2>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Owners and admins can change the current workspace plan here. Other
            roles can view the current billing level but cannot change it.
          </p>
          <div className="mt-5 grid gap-3 sm:grid-cols-2">
            {planOptions.map((option) => {
              const isCurrent = currentPlan === option.plan;
              const isPending =
                updatePlanMutation.isPending &&
                updatePlanMutation.variables === option.plan;

              return (
                <SurfaceCard
                  key={option.plan}
                  className={[
                    "border-[var(--app-line)] bg-[var(--app-card-strong)] p-4",
                    isCurrent ? "border-[var(--app-accent)]" : "",
                  ].join(" ")}
                >
                  <p className="text-sm font-semibold text-[var(--app-text)]">
                    {option.label}
                  </p>
                  <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                    {option.detail}
                  </p>
                  <button
                    type="button"
                    onClick={() => updatePlanMutation.mutate(option.plan)}
                    disabled={isCurrent || !canManagePlan || updatePlanMutation.isPending}
                    className={[
                      "mt-4 inline-flex min-h-11 w-full items-center justify-center rounded-2xl px-4 text-sm font-semibold transition disabled:opacity-60",
                      isCurrent
                        ? "border border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text)]"
                        : "bg-[var(--app-accent)] text-[var(--app-accent-ink)] hover:brightness-105",
                    ].join(" ")}
                  >
                    {isPending
                      ? "Updating..."
                      : isCurrent
                        ? "Current plan"
                        : `Switch to ${option.label}`}
                  </button>
                </SurfaceCard>
              );
            })}
          </div>
          {!canManagePlan ? (
            <p className="mt-4 text-sm text-[var(--app-muted)]">
              Plan changes are limited to team admins and owners.
            </p>
          ) : null}
          {updatePlanMutation.isError ? (
            <p className="mt-4 text-sm text-[var(--app-danger-text)]">
              {updatePlanMutation.error instanceof Error
                ? updatePlanMutation.error.message
                : "Failed to update the team plan."}
            </p>
          ) : null}
        </PageSection>
      </section>

      <PageSection className="p-5 sm:p-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <h2 className="text-xl font-semibold text-[var(--app-text)]">
              Members
            </h2>
            <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
              Roles are read, edit, admin, and owner. Only admins and owners
              can add or remove members.
            </p>
          </div>
          {canManageMembers ? (
            <button
              type="button"
              onClick={() => setInviteOpen(true)}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Invite member
            </button>
          ) : null}
        </div>

        <div className="mt-5 space-y-3">
          {team.members.map((member) => {
            const canMutateMember =
              canManageMembers &&
              member.id !== team.membership.id &&
              (currentRole === "owner" || member.role !== "owner");
            const roleChoices =
              currentRole === "owner"
                ? memberRoleOptions
                : memberRoleOptions.filter((role) => role !== "owner");

            return (
              <SurfaceCard
                key={member.id}
                className="border-none bg-[var(--app-card-strong)] p-4"
              >
                <div className="flex flex-col gap-4">
                  <div className="flex flex-wrap items-start justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-[var(--app-text)]">
                        {member.email}
                      </p>
                      <p className="mt-1 text-xs uppercase tracking-[0.16em] text-[var(--app-muted)]">
                        {member.role} • {member.status}
                      </p>
                      {member.inviteExpiresAt ? (
                        <p className="mt-2 text-xs text-[var(--app-muted)]">
                          Invite expires {formatDateTime(member.inviteExpiresAt)}
                        </p>
                      ) : null}
                    </div>
                    {member.id === team.membership.id ? (
                      <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-[11px] uppercase tracking-[0.16em] text-[var(--app-muted)]">
                        You
                      </span>
                    ) : null}
                  </div>

                  {canMutateMember ? (
                    <div className="flex flex-col gap-3 sm:flex-row">
                      <select
                        value={member.role}
                        onChange={(event) =>
                          updateMemberMutation.mutate({
                            memberId: member.id,
                            role: event.target.value as OrganizationMember["role"],
                          })
                        }
                        disabled={updateMemberMutation.isPending}
                        className="min-h-11 flex-1 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)] disabled:opacity-60"
                      >
                        {roleChoices.map((role) => (
                          <option key={role} value={role}>
                            {role}
                          </option>
                        ))}
                      </select>
                      <button
                        type="button"
                        onClick={() => removeMemberMutation.mutate(member.id)}
                        disabled={removeMemberMutation.isPending}
                        className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 text-sm font-semibold text-[var(--app-danger-text)] disabled:opacity-60"
                      >
                        Remove
                      </button>
                    </div>
                  ) : null}
                </div>
              </SurfaceCard>
            );
          })}
        </div>

        {(inviteMutation.isError ||
          updateMemberMutation.isError ||
          removeMemberMutation.isError) ? (
          <p className="mt-4 text-sm text-[var(--app-danger-text)]">
            {(inviteMutation.error ||
              updateMemberMutation.error ||
              removeMemberMutation.error) instanceof Error
              ? (
                  inviteMutation.error ||
                  updateMemberMutation.error ||
                  removeMemberMutation.error
                )!.message
              : "Failed to update team membership."}
          </p>
        ) : null}
      </PageSection>

      <Sheet open={inviteOpen} onClose={() => setInviteOpen(false)} title="Invite team member">
        <div className="space-y-4">
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">
              Email
            </span>
            <input
              value={inviteEmail}
              onChange={(event) => setInviteEmail(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="teammate@example.com"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">
              Role
            </span>
            <select
              value={inviteRole}
              onChange={(event) =>
                setInviteRole(
                  event.target.value as Exclude<OrganizationMember["role"], "owner">,
                )
              }
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            >
              <option value="read">Read</option>
              <option value="edit">Edit</option>
              <option value="admin">Admin</option>
            </select>
          </label>
          <button
            type="button"
            onClick={() => inviteMutation.mutate()}
            disabled={inviteMutation.isPending || !inviteEmail.trim()}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {inviteMutation.isPending ? "Inviting..." : "Send invite"}
          </button>
        </div>
      </Sheet>

      <ConfirmSheet
        open={deleteOpen}
        onClose={() => setDeleteOpen(false)}
        title="Delete team"
        description="This removes the team record from the current workspace context and marks memberships removed. Use this only when you are sure the team should be dismantled."
        confirmLabel={deleteTeamMutation.isPending ? "Deleting..." : "Delete team"}
        destructive
        onConfirm={() => deleteTeamMutation.mutate()}
      />
    </PageContainer>
  );
}
