import { Link } from 'react-router'
import { EmptyState, PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'

function StudyGroupsPlaceholder({
  title,
  description,
}: {
  title: string
  description: string
}) {
  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-5 sm:p-6">
        <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Study Groups</p>
        <h1 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--app-text)]">{title}</h1>
        <p className="mt-3 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)]">{description}</p>
      </PageSection>

      <div className="grid gap-4 md:grid-cols-3">
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Launch timing</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">Coming June 2026</p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            We&apos;re holding the public rollout until the collaboration and admin experience is polished enough for real teams.
          </p>
        </SurfaceCard>

        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">What it will include</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">Source decks + personal installs</p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Members will install their own study copies, keep private review history, and opt into published updates when they are ready.
          </p>
        </SurfaceCard>

        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">In the meantime</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">Use decks, templates, and marketplace</p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            The current launch path is solo and team-ready study with deck management, plan controls, and versioned installs through marketplace.
          </p>
        </SurfaceCard>
      </div>

      <EmptyState
        title="Study Groups are coming soon"
        description="We&apos;ve hidden this area for now so early users don&apos;t run into a half-finished collaboration workflow. We&apos;re targeting June 2026 for the first public release."
        action={
          <div className="flex flex-wrap items-center justify-center gap-3">
            <Link
              to="/decks"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Open decks
            </Link>
            <Link
              to="/marketplace"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-5 text-sm font-medium text-[var(--app-text)]"
            >
              Browse marketplace
            </Link>
          </div>
        }
      />
    </PageContainer>
  )
}

export function StudyGroupsPage() {
  return (
    <StudyGroupsPlaceholder
      title="Collaborative study spaces are on the way."
      description="Study Groups are planned for the June 2026 release. We&apos;re keeping the collaboration model private for now while we finish team roles, publishing flows, and the install/update experience."
    />
  )
}

export function StudyGroupDetailPage() {
  return (
    <StudyGroupsPlaceholder
      title="This Study Group view is not public yet."
      description="Direct Study Group routes are temporarily disabled while we prepare the June 2026 rollout. For now, use decks and marketplace installs as the supported study workflow."
    />
  )
}
