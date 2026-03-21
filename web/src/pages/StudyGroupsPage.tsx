import { useQuery } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'
import { PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'

const upcomingFeatures = [
  'Deck-linked study groups with owner, admin, and member roles.',
  'Invite and remove members with explicit membership status.',
  'Shared deck dashboards, weekly leaderboards, and announcements.',
]

const futureWorkflows = [
  'Attach a primary deck and keep collaboration scoped to that study surface.',
  'Track member participation, due-card totals, and weekly activity from one dashboard.',
  'Open the feature to Team and Enterprise workspaces first, with invitees able to join.',
]

export function StudyGroupsPage() {
  const repository = useAppRepository()
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })

  const canCreate = entitlementsQuery.data?.features.studyGroups ?? false

  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-5 sm:p-6">
        <p className="text-[11px] uppercase tracking-[0.26em] text-[var(--app-muted)]">Coming soon</p>
        <h1 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--app-text)]">Study Groups</h1>
        <p className="mt-4 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)]">
          Study Groups will give each shared deck a collaborative layer for invites, participation, and accountability without overloading basic deck CRUD.
        </p>
      </PageSection>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,0.9fr)]">
        <PageSection className="p-5 sm:p-6">
          <div className="grid gap-4 md:grid-cols-2">
            <SurfaceCard className="space-y-3 border-none bg-[var(--app-card-strong)]">
              <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Next tranche</p>
              <ul className="space-y-3 text-sm leading-6 text-[var(--app-text-soft)]">
                {upcomingFeatures.map((feature) => (
                  <li key={feature} className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] px-4 py-3">
                    {feature}
                  </li>
                ))}
              </ul>
            </SurfaceCard>

            <SurfaceCard className="space-y-3 border-none bg-[var(--app-card-strong)]">
              <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">How it will work</p>
              <ul className="space-y-3 text-sm leading-6 text-[var(--app-text-soft)]">
                {futureWorkflows.map((workflow) => (
                  <li key={workflow} className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] px-4 py-3">
                    {workflow}
                  </li>
                ))}
              </ul>
            </SurfaceCard>
          </div>
        </PageSection>

        <PageSection className="p-5 sm:p-6">
          <div className="rounded-[1.75rem] border border-[var(--app-line)] bg-[linear-gradient(180deg,rgba(112,214,108,0.10),rgba(112,214,108,0.03))] p-5">
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Plan access</p>
            <p className="mt-3 text-2xl font-semibold text-[var(--app-text)]">
              {canCreate ? 'Eligible to create later' : 'Upgrade required'}
            </p>
            <p className="mt-3 text-sm leading-6 text-[var(--app-text-soft)]">
              {canCreate
                ? 'This workspace can create study groups as soon as the feature is shipped.'
                : 'Study group creation is reserved for Team and Enterprise workspaces.'}
            </p>
          </div>

          <div className="mt-4 grid gap-3">
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-sm font-semibold text-[var(--app-text)]">Shared decks</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Group-owned decks stay collaborative while each member keeps their own review state and due queue.
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-sm font-semibold text-[var(--app-text)]">Member management</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                Owners and admins will be able to invite, approve, and remove people without touching core deck ownership.
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-sm font-semibold text-[var(--app-text)]">Team-first rollout</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                The placeholder stays intentionally simple until the per-user review-state split is ready.
              </p>
            </SurfaceCard>
          </div>
        </PageSection>
      </div>
    </PageContainer>
  )
}
