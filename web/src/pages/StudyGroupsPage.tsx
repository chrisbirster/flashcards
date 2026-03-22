import { type ReactNode, useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useNavigate, useParams, useSearchParams } from 'react-router'
import { useAppRepository } from '#/lib/app-repository'
import { EmptyState, PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'
import { Sheet } from '#/components/sheet'

function formatDateTime(value?: string) {
  if (!value) return 'Unknown'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? 'Unknown' : date.toLocaleString()
}

function SectionHeader({
  eyebrow,
  title,
  description,
  action,
}: {
  eyebrow: string
  title: string
  description: string
  action?: ReactNode
}) {
  return (
    <PageSection className="p-5 sm:p-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">{eyebrow}</p>
          <h1 className="mt-3 text-3xl font-semibold tracking-tight text-[var(--app-text)]">{title}</h1>
          <p className="mt-3 max-w-3xl text-sm leading-7 text-[var(--app-text-soft)]">{description}</p>
        </div>
        {action}
      </div>
    </PageSection>
  )
}

export function StudyGroupsPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const [createOpen, setCreateOpen] = useState(false)
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [primaryDeckId, setPrimaryDeckId] = useState<number | ''>('')

  const groupsQuery = useQuery({
    queryKey: ['study-groups'],
    queryFn: () => repository.fetchStudyGroups(),
  })
  const decksQuery = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })

  const createMutation = useMutation({
    mutationFn: () =>
      repository.createStudyGroup({
        name,
        description,
        primaryDeckId: Number(primaryDeckId),
        visibility: 'private',
        joinPolicy: 'invite',
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['study-groups'] })
      setCreateOpen(false)
      setName('')
      setDescription('')
      setPrimaryDeckId('')
    },
  })

  const canCreate = entitlementsQuery.data?.features.studyGroups ?? false

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Study Groups"
        title="Canonical source decks with private member installs."
        description="Owners and admins publish explicit source versions. Members install personal copies, keep their own review history, and opt into updates when they are ready."
        action={
          canCreate ? (
            <button
              type="button"
              onClick={() => setCreateOpen(true)}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Create group
            </button>
          ) : null
        }
      />

      <div className="grid gap-4 md:grid-cols-3">
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Model</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">Source deck + personal installs</p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            This phase avoids live shared review queues. Members study their own installed copies and updates stay opt-in.
          </p>
        </SurfaceCard>
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Permissions</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">
            {canCreate ? 'Team / Enterprise enabled' : 'Upgrade required'}
          </p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Creation and management are gated to Team and Enterprise workspaces. Invited members can still join and install.
          </p>
        </SurfaceCard>
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Current status</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">
            {groupsQuery.data?.length ?? 0} group{groupsQuery.data?.length === 1 ? '' : 's'}
          </p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Pending invites, active groups, install status, and update availability now come from the real Phase 2 APIs.
          </p>
        </SurfaceCard>
      </div>

      {groupsQuery.isLoading ? (
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading study groups...</PageSection>
      ) : groupsQuery.isError ? (
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {groupsQuery.error instanceof Error ? groupsQuery.error.message : 'Failed to load study groups.'}
        </PageSection>
      ) : groupsQuery.data && groupsQuery.data.length > 0 ? (
        <div className="grid gap-4">
          {groupsQuery.data.map((group) => (
            <PageSection key={group.id} className="p-5 sm:p-6">
              <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
                <div className="space-y-3">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
                      {group.role || 'member'}
                    </span>
                    <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-text-soft)]">
                      {group.membershipStatus}
                    </span>
                    {group.updateAvailable ? (
                      <span className="rounded-full bg-[var(--app-accent)] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.16em] text-[var(--app-accent-ink)]">
                        Update available
                      </span>
                    ) : null}
                  </div>
                  <div>
                    <h2 className="text-xl font-semibold text-[var(--app-text)]">{group.name}</h2>
                    <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">{group.description || 'No description yet.'}</p>
                  </div>
                  <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
                    <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                      <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Source deck</p>
                      <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">{group.sourceDeckName || 'Unknown deck'}</p>
                    </SurfaceCard>
                    <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                      <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Latest version</p>
                      <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">
                        {group.latestVersionNumber > 0 ? `v${group.latestVersionNumber}` : 'Not published'}
                      </p>
                    </SurfaceCard>
                    <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                      <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Members</p>
                      <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">{group.memberCount}</p>
                    </SurfaceCard>
                    <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                      <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Your install</p>
                      <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">
                        {group.currentUserInstall
                          ? `v${group.currentUserInstall.sourceVersionNumber} • ${group.currentUserInstall.syncState}`
                          : 'Not installed'}
                      </p>
                    </SurfaceCard>
                  </div>
                </div>
                <Link
                  to={`/study-groups/${group.id}`}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
                >
                  Open group
                </Link>
              </div>
            </PageSection>
          ))}
        </div>
      ) : (
        <EmptyState
          title="No study groups yet"
          description={
            canCreate
              ? 'Create the first group from a source deck, publish a version, and invite members to install their own copies.'
              : 'Study group creation is reserved for Team and Enterprise workspaces. Invited members can still join once a group exists.'
          }
          action={
            canCreate ? (
              <button
                type="button"
                onClick={() => setCreateOpen(true)}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
              >
                Create group
              </button>
            ) : undefined
          }
        />
      )}

      <Sheet open={createOpen} onClose={() => setCreateOpen(false)} title="Create study group">
        <div className="space-y-4">
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Name</span>
            <input
              value={name}
              onChange={(event) => setName(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Med School Cohort"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Description</span>
            <textarea
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              rows={4}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Canonical source deck for the cohort. Members install personal copies and opt into updates."
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Source deck</span>
            <select
              value={primaryDeckId}
              onChange={(event) => setPrimaryDeckId(event.target.value ? Number(event.target.value) : '')}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            >
              <option value="">Select a deck</option>
              {(decksQuery.data ?? []).map((deck) => (
                <option key={deck.id} value={deck.id}>
                  {deck.name}
                </option>
              ))}
            </select>
          </label>
          {createMutation.isError ? (
            <p className="text-sm text-[var(--app-danger-text)]">
              {createMutation.error instanceof Error ? createMutation.error.message : 'Failed to create group.'}
            </p>
          ) : null}
          <button
            type="button"
            onClick={() => createMutation.mutate()}
            disabled={createMutation.isPending || !name.trim() || !primaryDeckId}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {createMutation.isPending ? 'Creating...' : 'Create group'}
          </button>
        </div>
      </Sheet>
    </PageContainer>
  )
}

export function StudyGroupDetailPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const params = useParams()
  const groupId = params.groupId ?? ''
  const [inviteOpen, setInviteOpen] = useState(false)
  const [publishOpen, setPublishOpen] = useState(false)
  const [inviteEmail, setInviteEmail] = useState('')
  const [inviteRole, setInviteRole] = useState<'admin' | 'member'>('member')
  const [changeSummary, setChangeSummary] = useState('')
  const [workspaceId, setWorkspaceId] = useState('')

  const detailQuery = useQuery({
    queryKey: ['study-group', groupId],
    queryFn: () => repository.fetchStudyGroup(groupId),
    enabled: groupId.length > 0,
  })

  const detail = detailQuery.data
  const activeMembers = useMemo(
    () => (detail?.members ?? []).filter((member) => member.status === 'active'),
    [detail?.members]
  )

  const refreshGroup = () => {
    queryClient.invalidateQueries({ queryKey: ['study-groups'] })
    queryClient.invalidateQueries({ queryKey: ['study-group', groupId] })
  }

  const inviteMutation = useMutation({
    mutationFn: () => repository.inviteStudyGroupMember(groupId, {email: inviteEmail, role: inviteRole}),
    onSuccess: () => {
      refreshGroup()
      setInviteEmail('')
      setInviteRole('member')
      setInviteOpen(false)
    },
  })

  const publishMutation = useMutation({
    mutationFn: () => repository.publishStudyGroupVersion(groupId, {changeSummary}),
    onSuccess: () => {
      refreshGroup()
      setChangeSummary('')
      setPublishOpen(false)
    },
  })

  const installMutation = useMutation({
    mutationFn: (destinationWorkspaceId: string) => repository.installStudyGroupDeck(groupId, {destinationWorkspaceId}),
    onSuccess: refreshGroup,
  })

  const updateInstallMutation = useMutation({
    mutationFn: ({installId, destinationWorkspaceId}: {installId: string; destinationWorkspaceId?: string}) =>
      repository.updateStudyGroupInstall(groupId, installId, {destinationWorkspaceId}),
    onSuccess: refreshGroup,
  })

  const removeInstallMutation = useMutation({
    mutationFn: (installId: string) => repository.removeStudyGroupInstall(groupId, installId),
    onSuccess: refreshGroup,
  })

  const memberMutation = useMutation({
    mutationFn: ({memberId, role}: {memberId: string; role: 'admin' | 'member'}) =>
      repository.updateStudyGroupMember(groupId, memberId, {role}),
    onSuccess: refreshGroup,
  })

  const removeMemberMutation = useMutation({
    mutationFn: (memberId: string) => repository.deleteStudyGroupMember(groupId, memberId),
    onSuccess: refreshGroup,
  })

  if (detailQuery.isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading study group...</PageSection>
      </PageContainer>
    )
  }

  if (detailQuery.isError || !detail) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {detailQuery.error instanceof Error ? detailQuery.error.message : 'Failed to load study group.'}
        </PageSection>
      </PageContainer>
    )
  }

  const availableWorkspaces = detail.availableWorkspaces.length > 0 ? detail.availableWorkspaces : []
  const selectedWorkspaceId = workspaceId || availableWorkspaces[0]?.id || ''

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Study Group"
        title={detail.group.name}
        description={detail.group.description || 'Canonical source deck with explicit versions and private member installs.'}
        action={
          <div className="flex flex-wrap gap-2">
            {detail.canInvite ? (
              <button
                type="button"
                onClick={() => setInviteOpen(true)}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-semibold text-[var(--app-text)]"
              >
                Invite
              </button>
            ) : null}
            {detail.canPublishVersion ? (
              <button
                type="button"
                onClick={() => setPublishOpen(true)}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
              >
                Publish update
              </button>
            ) : null}
          </div>
        }
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1.25fr)_minmax(0,0.9fr)]">
        <div className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Source deck</p>
              <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">{detail.sourceDeckName}</p>
              <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                Members never study this deck directly. It is the canonical publish source for versioned installs.
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Membership</p>
              <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
                {detail.role} • {detail.membershipStatus}
              </p>
              <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                Owners and admins can invite members and publish source versions. Member installs keep private review history.
              </p>
            </SurfaceCard>
          </div>

          <PageSection className="p-5 sm:p-6">
            <div className="flex flex-col gap-4">
              <div>
                <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Latest source version</p>
                <p className="mt-3 text-2xl font-semibold text-[var(--app-text)]">
                  {detail.latestVersion ? `v${detail.latestVersion.versionNumber}` : 'Not published yet'}
                </p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  {detail.latestVersion
                    ? `${detail.latestVersion.noteCount} notes • ${detail.latestVersion.cardCount} cards • ${detail.latestVersion.changeSummary || 'No change summary'}`
                    : 'Publish the first version to make this source deck installable.'}
                </p>
              </div>
              {detail.currentUserInstall ? (
                <div className="rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
                  <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Your install</p>
                  <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
                    v{detail.currentUserInstall.sourceVersionNumber} • {detail.currentUserInstall.syncState}
                  </p>
                  <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                    {detail.currentUserInstall.installedDeckName || 'Installed deck'} lives in your workspace. Updating creates a fresh new copy and keeps the old one intact.
                  </p>
                  <div className="mt-4 flex flex-col gap-3 sm:flex-row">
                    {detail.updateAvailable ? (
                      <button
                        type="button"
                        onClick={() =>
                          updateInstallMutation.mutate({
                            installId: detail.currentUserInstall!.id,
                            destinationWorkspaceId: selectedWorkspaceId || detail.currentUserInstall!.destinationWorkspaceId,
                          })
                        }
                        disabled={updateInstallMutation.isPending}
                        className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
                      >
                        {updateInstallMutation.isPending ? 'Updating...' : 'Install latest version'}
                      </button>
                    ) : null}
                    <button
                      type="button"
                      onClick={() => removeInstallMutation.mutate(detail.currentUserInstall!.id)}
                      disabled={removeInstallMutation.isPending}
                      className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-semibold text-[var(--app-text-soft)] disabled:opacity-60"
                    >
                      {removeInstallMutation.isPending ? 'Removing...' : 'Remove install'}
                    </button>
                  </div>
                </div>
              ) : (
                <div className="rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
                  <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Install</p>
                  <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">No personal copy yet</p>
                  <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                    Install the latest published version into your workspace. Studying it will not change anyone else’s due queue.
                  </p>
                  <div className="mt-4 flex flex-col gap-3">
                    <select
                      value={selectedWorkspaceId}
                      onChange={(event) => setWorkspaceId(event.target.value)}
                      className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                    >
                      {availableWorkspaces.map((workspace) => (
                        <option key={workspace.id} value={workspace.id}>
                          {workspace.name}
                        </option>
                      ))}
                    </select>
                    <button
                      type="button"
                      onClick={() => installMutation.mutate(selectedWorkspaceId)}
                      disabled={installMutation.isPending || !selectedWorkspaceId || !detail.latestVersion}
                      className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
                    >
                      {installMutation.isPending ? 'Installing...' : 'Install latest version'}
                    </button>
                  </div>
                </div>
              )}
              {(installMutation.error || updateInstallMutation.error || removeInstallMutation.error) ? (
                <p className="text-sm text-[var(--app-danger-text)]">
                  {(installMutation.error || updateInstallMutation.error || removeInstallMutation.error) instanceof Error
                    ? ((installMutation.error || updateInstallMutation.error || removeInstallMutation.error) as Error).message
                    : 'Study group install action failed.'}
                </p>
              ) : null}
            </div>
          </PageSection>

          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Published versions</p>
            <div className="mt-4 space-y-3">
              {detail.versions.length > 0 ? (
                detail.versions.map((version) => (
                  <SurfaceCard key={version.id} className="border-none bg-[var(--app-card-strong)] p-4">
                    <div className="flex items-center justify-between gap-3">
                      <div>
                        <p className="text-sm font-semibold text-[var(--app-text)]">v{version.versionNumber}</p>
                        <p className="mt-1 text-xs text-[var(--app-muted)]">{formatDateTime(version.createdAt)}</p>
                      </div>
                      <span className="text-sm text-[var(--app-text-soft)]">
                        {version.noteCount} notes • {version.cardCount} cards
                      </span>
                    </div>
                    <p className="mt-3 text-sm leading-6 text-[var(--app-text-soft)]">
                      {version.changeSummary || 'No change summary provided.'}
                    </p>
                  </SurfaceCard>
                ))
              ) : (
                <EmptyState
                  title="No published versions"
                  description="Publish the first source version to make this group installable."
                />
              )}
            </div>
          </PageSection>
        </div>

        <div className="space-y-4">
          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Dashboard</p>
            <div className="mt-4 grid gap-3 sm:grid-cols-2">
              <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Members</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">{detail.dashboard.memberCount}</p>
              </SurfaceCard>
              <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Active 7d</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">{detail.dashboard.activeMembers7d}</p>
              </SurfaceCard>
              <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Reviews 7d</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">{detail.dashboard.reviews7d}</p>
              </SurfaceCard>
              <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Latest adoption</p>
                <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">{detail.dashboard.latestVersionAdoption}</p>
              </SurfaceCard>
            </div>
            <div className="mt-4 space-y-3">
              {(detail.dashboard.leaderboard ?? []).map((entry) => (
                <SurfaceCard key={`${entry.email}-${entry.reviews7d}`} className="border-none bg-[var(--app-card-strong)] p-4">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="text-sm font-semibold text-[var(--app-text)]">{entry.displayName || entry.email}</p>
                      <p className="mt-1 text-xs text-[var(--app-muted)]">{entry.email}</p>
                    </div>
                    <span className="text-sm font-semibold text-[var(--app-accent)]">{entry.reviews7d} reviews</span>
                  </div>
                </SurfaceCard>
              ))}
            </div>
          </PageSection>

          <PageSection className="p-5 sm:p-6">
            <div className="flex items-center justify-between gap-3">
              <div>
                <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Members</p>
                <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                  Active members install personal copies. Owners and admins manage invites and source publishing.
                </p>
              </div>
              <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
                {activeMembers.length} active
              </span>
            </div>
            <div className="mt-4 space-y-3">
              {detail.members.map((member) => (
                <SurfaceCard key={member.id} className="border-none bg-[var(--app-card-strong)] p-4">
                  <div className="flex flex-col gap-3">
                    <div className="flex items-start justify-between gap-3">
                      <div>
                        <p className="text-sm font-semibold text-[var(--app-text)]">{member.email}</p>
                        <p className="mt-1 text-xs text-[var(--app-muted)]">
                          {member.role} • {member.status}
                        </p>
                      </div>
                      {detail.canEdit && member.role !== 'owner' ? (
                        <div className="flex flex-wrap gap-2">
                          <button
                            type="button"
                            onClick={() => memberMutation.mutate({memberId: member.id, role: member.role === 'admin' ? 'member' : 'admin'})}
                            className="inline-flex min-h-10 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 text-xs font-semibold text-[var(--app-text-soft)]"
                          >
                            {member.role === 'admin' ? 'Make member' : 'Make admin'}
                          </button>
                          <button
                            type="button"
                            onClick={() => removeMemberMutation.mutate(member.id)}
                            className="inline-flex min-h-10 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 text-xs font-semibold text-[var(--app-danger-text)]"
                          >
                            Remove
                          </button>
                        </div>
                      ) : null}
                    </div>
                    {member.inviteExpiresAt ? (
                      <p className="text-xs text-[var(--app-muted)]">Invite expires {formatDateTime(member.inviteExpiresAt)}</p>
                    ) : null}
                  </div>
                </SurfaceCard>
              ))}
            </div>
          </PageSection>

          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Recent activity</p>
            <div className="mt-4 space-y-3">
              {detail.recentEvents.length > 0 ? (
                detail.recentEvents.map((event) => (
                  <SurfaceCard key={event.id} className="border-none bg-[var(--app-card-strong)] p-4">
                    <p className="text-sm font-semibold text-[var(--app-text)]">{event.eventType.replaceAll('_', ' ')}</p>
                    <p className="mt-2 text-xs text-[var(--app-muted)]">{formatDateTime(event.createdAt)}</p>
                  </SurfaceCard>
                ))
              ) : (
                <p className="text-sm text-[var(--app-text-soft)]">No activity yet.</p>
              )}
            </div>
          </PageSection>
        </div>
      </div>

      <Sheet open={inviteOpen} onClose={() => setInviteOpen(false)} title="Invite member">
        <div className="space-y-4">
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Email</span>
            <input
              value={inviteEmail}
              onChange={(event) => setInviteEmail(event.target.value)}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="student@example.com"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Role</span>
            <select
              value={inviteRole}
              onChange={(event) => setInviteRole(event.target.value as 'admin' | 'member')}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
            >
              <option value="member">Member</option>
              <option value="admin">Admin</option>
            </select>
          </label>
          {inviteMutation.isError ? (
            <p className="text-sm text-[var(--app-danger-text)]">
              {inviteMutation.error instanceof Error ? inviteMutation.error.message : 'Failed to invite member.'}
            </p>
          ) : null}
          <button
            type="button"
            onClick={() => inviteMutation.mutate()}
            disabled={inviteMutation.isPending || !inviteEmail.trim()}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {inviteMutation.isPending ? 'Inviting...' : 'Send invite'}
          </button>
        </div>
      </Sheet>

      <Sheet open={publishOpen} onClose={() => setPublishOpen(false)} title="Publish source version">
        <div className="space-y-4">
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Change summary</span>
            <textarea
              value={changeSummary}
              onChange={(event) => setChangeSummary(event.target.value)}
              rows={4}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Added new nephrology cards and corrected two anatomy answers."
            />
          </label>
          {publishMutation.isError ? (
            <p className="text-sm text-[var(--app-danger-text)]">
              {publishMutation.error instanceof Error ? publishMutation.error.message : 'Failed to publish version.'}
            </p>
          ) : null}
          <button
            type="button"
            onClick={() => publishMutation.mutate()}
            disabled={publishMutation.isPending}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {publishMutation.isPending ? 'Publishing...' : 'Publish version'}
          </button>
        </div>
      </Sheet>
    </PageContainer>
  )
}

export function StudyGroupJoinPage() {
  const repository = useAppRepository()
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const token = searchParams.get('token') ?? ''
  const [installLatest, setInstallLatest] = useState(true)
  const sessionQuery = useQuery({
    queryKey: ['auth-session'],
    queryFn: () => repository.fetchSession(),
  })
  const joinMutation = useMutation({
    mutationFn: () =>
      repository.joinStudyGroup({
        token,
        destinationWorkspaceId: sessionQuery.data?.workspace?.id ?? '',
        installLatest,
      }),
    onSuccess: (detail) => {
      navigate(`/study-groups/${detail.group.id}`)
    },
  })

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Join"
        title="Accept study group invite"
        description="Joining activates your membership. Installing the latest version creates a personal copy in your current workspace and leaves your review history private."
      />

      {!token ? (
        <EmptyState
          title="Invite token missing"
          description="Open the join link from your invite email so Vutadex can connect the request to the correct study group."
        />
      ) : (
        <PageSection className="mx-auto max-w-2xl p-5 sm:p-6">
          <div className="space-y-4">
            <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
              <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Destination workspace</p>
              <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
                {sessionQuery.data?.workspace?.name || 'Current workspace'}
              </p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                This phase installs into the active workspace from your session. You can update or remove the personal copy later from the group detail page.
              </p>
            </SurfaceCard>

            <label className="flex items-start gap-3 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
              <input
                type="checkbox"
                checked={installLatest}
                onChange={(event) => setInstallLatest(event.target.checked)}
                className="mt-1 h-4 w-4 rounded border-[var(--app-line-strong)]"
              />
              <span className="text-sm leading-6 text-[var(--app-text-soft)]">
                Install the latest published version immediately after joining.
              </span>
            </label>

            {joinMutation.isError ? (
              <p className="text-sm text-[var(--app-danger-text)]">
                {joinMutation.error instanceof Error ? joinMutation.error.message : 'Failed to join study group.'}
              </p>
            ) : null}

            <button
              type="button"
              onClick={() => joinMutation.mutate()}
              disabled={joinMutation.isPending || sessionQuery.isLoading}
              className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
            >
              {joinMutation.isPending ? 'Joining...' : 'Accept invite'}
            </button>
          </div>
        </PageSection>
      )}
    </PageContainer>
  )
}
