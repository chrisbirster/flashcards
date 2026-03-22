import { useEffect, useMemo, useState, type ReactNode } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useParams, useSearchParams } from 'react-router'
import { useAppRepository } from '#/lib/app-repository'
import { EmptyState, PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'
import { Sheet } from '#/components/sheet'
import type {
  CreateMarketplaceListingRequest,
  MarketplaceListingSummary,
  UpdateMarketplaceListingRequest,
} from '#/lib/api'

function formatDateTime(value?: string) {
  if (!value) return 'Unknown'
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? 'Unknown' : date.toLocaleString()
}

function formatPrice(mode: 'free' | 'premium', cents: number, currency: string) {
  if (mode === 'free') return 'Free'
  return new Intl.NumberFormat(undefined, {
    style: 'currency',
    currency: currency || 'USD',
    maximumFractionDigits: 0,
  }).format((cents || 0) / 100)
}

function parseTags(value: string) {
  return value
    .split(',')
    .map((tag) => tag.trim())
    .filter(Boolean)
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

function ListingCard({
  listing,
  action,
}: {
  listing: MarketplaceListingSummary
  action?: ReactNode
}) {
  return (
    <PageSection className="p-5 sm:p-6">
      <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
        <div className="space-y-3">
          <div className="flex flex-wrap items-center gap-2">
            {listing.category ? (
              <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">
                {listing.category}
              </span>
            ) : null}
            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-[11px] uppercase tracking-[0.18em] text-[var(--app-text-soft)]">
              {listing.status}
            </span>
            <span className="rounded-full bg-[var(--app-accent)] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.16em] text-[var(--app-accent-ink)]">
              {formatPrice(listing.priceMode, listing.priceCents, listing.currency)}
            </span>
            {listing.updateAvailable ? (
              <span className="rounded-full border border-[var(--app-accent)] bg-[var(--app-card)] px-3 py-1 text-[11px] font-semibold uppercase tracking-[0.16em] text-[var(--app-accent)]">
                Update available
              </span>
            ) : null}
          </div>
          <div>
            <h2 className="text-xl font-semibold text-[var(--app-text)]">{listing.title}</h2>
            <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
              {listing.summary || listing.description || 'No summary yet.'}
            </p>
          </div>
          <div className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
            <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Creator</p>
              <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">
                {listing.creatorDisplayName || listing.creatorEmail || 'Unknown'}
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Source deck</p>
              <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">{listing.sourceDeckName || 'Unknown deck'}</p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Latest version</p>
              <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">
                {listing.latestVersionNumber > 0 ? `v${listing.latestVersionNumber}` : 'Draft'}
              </p>
            </SurfaceCard>
            <SurfaceCard className="border-none bg-[var(--app-card-strong)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Installs</p>
              <p className="mt-2 text-sm font-semibold text-[var(--app-text)]">{listing.installCount}</p>
            </SurfaceCard>
          </div>
          {listing.tags.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {listing.tags.map((tag) => (
                <span
                  key={tag}
                  className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-xs text-[var(--app-text-soft)]"
                >
                  {tag}
                </span>
              ))}
            </div>
          ) : null}
          {listing.currentUserInstall ? (
            <p className="text-sm text-[var(--app-text-soft)]">
              Your install: v{listing.currentUserInstall.sourceVersionNumber}
            </p>
          ) : listing.currentUserLicense ? (
            <p className="text-sm text-[var(--app-text-soft)]">
              Purchased at v{listing.currentUserLicense.grantedVersionNumber}. You can install it into any workspace you own.
            </p>
          ) : null}
        </div>
        {action}
      </div>
    </PageSection>
  )
}

export function MarketplacePage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const [searchParams] = useSearchParams()
  const checkoutState = searchParams.get('checkout')
  const checkoutSessionId = searchParams.get('checkout_session_id') ?? ''
  const publishedQuery = useQuery({
    queryKey: ['marketplace-listings', 'public'],
    queryFn: () => repository.fetchMarketplaceListings(),
  })
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })
  const checkoutSyncQuery = useQuery({
    queryKey: ['marketplace-checkout-sync', checkoutSessionId],
    queryFn: () => repository.syncMarketplaceCheckoutSession(checkoutSessionId),
    enabled: checkoutState === 'success' && checkoutSessionId.length > 0,
  })

  useEffect(() => {
    if (!checkoutSyncQuery.data?.completed) return
    queryClient.invalidateQueries({queryKey: ['marketplace-listings']})
  }, [checkoutSyncQuery.data?.completed, queryClient])

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Marketplace"
        title="Install expert decks without sharing your review history."
        description="Marketplace installs reuse the Study Group model: source versions are published explicitly, installs create personal workspace copies, and premium checkout now layers on top of the same install foundation."
        action={
          entitlementsQuery.data?.features.marketplacePublish ? (
            <Link
              to="/marketplace/publish"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Publish listing
            </Link>
          ) : null
        }
      />

      {checkoutState === 'cancelled' ? (
        <PageSection className="border-[var(--app-line)] bg-[var(--app-card-strong)] p-5 text-sm text-[var(--app-text-soft)]">
          Checkout was cancelled. Your marketplace install options are unchanged.
        </PageSection>
      ) : null}

      {checkoutSyncQuery.isLoading ? (
        <PageSection className="border-[var(--app-line)] bg-[var(--app-card-strong)] p-5 text-sm text-[var(--app-text-soft)]">
          Finalizing your marketplace purchase...
        </PageSection>
      ) : null}

      {checkoutSyncQuery.isSuccess && checkoutSyncQuery.data.completed ? (
        <PageSection className="border-[var(--app-accent)] bg-[var(--app-card-strong)] p-5 text-sm text-[var(--app-text)]">
          Purchase confirmed. Your premium marketplace license is active and the deck can now be installed into your workspaces.
        </PageSection>
      ) : null}

      {checkoutSyncQuery.isError ? (
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {checkoutSyncQuery.error instanceof Error ? checkoutSyncQuery.error.message : 'Failed to confirm marketplace checkout.'}
        </PageSection>
      ) : null}

      <div className="grid gap-4 md:grid-cols-3">
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Install model</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">Personal copies, not shared queues</p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Listings publish explicit source versions. Installs create workspace-local copies so review history stays private per user.
          </p>
        </SurfaceCard>
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Publishing</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">
            {entitlementsQuery.data?.features.marketplacePublish ? 'Pro and above enabled' : 'Upgrade required'}
          </p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            This phase supports creator listings, explicit version publishing, and premium licensing. Creator payout setup is required before premium decks can go live.
          </p>
        </SurfaceCard>
        <SurfaceCard className="border-none bg-[var(--app-card-strong)]">
          <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Catalog</p>
          <p className="mt-3 text-base font-semibold text-[var(--app-text)]">
            {publishedQuery.data?.length ?? 0} listing{publishedQuery.data?.length === 1 ? '' : 's'}
          </p>
          <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
            Browse published decks now. Free decks install directly, and premium decks unlock after checkout grants a personal license.
          </p>
        </SurfaceCard>
      </div>

      {publishedQuery.isLoading ? (
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading marketplace...</PageSection>
      ) : publishedQuery.isError ? (
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {publishedQuery.error instanceof Error ? publishedQuery.error.message : 'Failed to load marketplace listings.'}
        </PageSection>
      ) : publishedQuery.data && publishedQuery.data.length > 0 ? (
        <div className="grid gap-4">
          {publishedQuery.data.map((listing) => (
            <ListingCard
              key={listing.id}
              listing={listing}
              action={
                <Link
                  to={`/marketplace/${listing.slug}`}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
                >
                  View listing
                </Link>
              }
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title="No published listings yet"
          description="Publish the first source deck to start the marketplace catalog. Free installs are ready now, and premium listings can sell once creator payout setup is complete."
          action={
            entitlementsQuery.data?.features.marketplacePublish ? (
              <Link
                to="/marketplace/publish"
                className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
              >
                Publish listing
              </Link>
            ) : undefined
          }
        />
      )}
    </PageContainer>
  )
}

export function MarketplaceDetailPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const params = useParams()
  const slug = params.slug ?? ''
  const [workspaceId, setWorkspaceId] = useState('')

  const detailQuery = useQuery({
    queryKey: ['marketplace-listing', slug],
    queryFn: () => repository.fetchMarketplaceListing(slug),
    enabled: slug.length > 0,
  })

  const detail = detailQuery.data
  const availableWorkspaces = detail?.availableWorkspaces ?? []
  const selectedWorkspaceId = workspaceId || availableWorkspaces[0]?.id || ''

  const refreshListing = () => {
    queryClient.invalidateQueries({ queryKey: ['marketplace-listings'] })
    queryClient.invalidateQueries({ queryKey: ['marketplace-listing', slug] })
    queryClient.invalidateQueries({ queryKey: ['decks'] })
  }

  const installMutation = useMutation({
    mutationFn: () => repository.installMarketplaceListing(slug, {destinationWorkspaceId: selectedWorkspaceId}),
    onSuccess: refreshListing,
  })
  const checkoutMutation = useMutation({
    mutationFn: () => repository.checkoutMarketplaceListing(slug),
    onSuccess: (response) => {
      if (response.checkoutUrl) {
        window.location.assign(response.checkoutUrl)
        return
      }
      refreshListing()
    },
  })
  const updateInstallMutation = useMutation({
    mutationFn: ({installId, destinationWorkspaceId}: {installId: string; destinationWorkspaceId?: string}) =>
      repository.updateMarketplaceInstall(slug, installId, {destinationWorkspaceId}),
    onSuccess: refreshListing,
  })
  const removeInstallMutation = useMutation({
    mutationFn: (installId: string) => repository.removeMarketplaceInstall(slug, installId),
    onSuccess: refreshListing,
  })

  if (detailQuery.isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading listing...</PageSection>
      </PageContainer>
    )
  }

  if (detailQuery.isError || !detail) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {detailQuery.error instanceof Error ? detailQuery.error.message : 'Failed to load listing.'}
        </PageSection>
      </PageContainer>
    )
  }

  const listing = detail.listing
  const hasPremiumAccess = listing.priceMode === 'free' || listing.canEdit || Boolean(detail.currentUserLicense)
  const installLabel = listing.priceMode === 'premium'
    ? detail.currentUserLicense
      ? 'Install purchased deck'
      : 'Install creator copy'
    : 'Install free deck'

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Marketplace Listing"
        title={listing.title}
        description={listing.description || listing.summary || 'No listing description yet.'}
        action={
          listing.canEdit ? (
            <Link
              to="/marketplace/publish"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-5 text-sm font-semibold text-[var(--app-text)]"
            >
              Manage listing
            </Link>
          ) : null
        }
      />

      <div className="grid gap-4 lg:grid-cols-[minmax(0,1.25fr)_minmax(0,0.9fr)]">
        <div className="space-y-4">
          <ListingCard listing={listing} />

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
                <EmptyState title="Not published yet" description="This listing needs a published version before users can install it." />
              )}
            </div>
          </PageSection>
        </div>

        <div className="space-y-4">
          <PageSection className="p-5 sm:p-6">
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Install</p>
            {detail.currentUserInstall ? (
              <div className="mt-4 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
                <p className="text-lg font-semibold text-[var(--app-text)]">
                  v{detail.currentUserInstall.sourceVersionNumber} installed
                </p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  {detail.currentUserInstall.installedDeckName || 'Installed deck'} lives in your workspace as a personal copy.
                </p>
                <div className="mt-4 flex flex-col gap-3">
                  {detail.updateAvailable ? (
                    <button
                      type="button"
                      onClick={() =>
                        updateInstallMutation.mutate({
                          installId: detail.currentUserInstall!.id,
                          destinationWorkspaceId: selectedWorkspaceId || detail.currentUserInstall!.workspaceId,
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
            ) : !hasPremiumAccess ? (
              <div className="mt-4 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
                <p className="text-lg font-semibold text-[var(--app-text)]">Premium deck purchase required</p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  Buying this deck grants a personal license first. After checkout, you can install versioned workspace copies whenever you need them.
                </p>
                <div className="mt-4 flex flex-col gap-3">
                  <button
                    type="button"
                    onClick={() => checkoutMutation.mutate()}
                    disabled={checkoutMutation.isPending}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
                  >
                    {checkoutMutation.isPending ? 'Processing...' : `Buy now for ${formatPrice(listing.priceMode, listing.priceCents, listing.currency)}`}
                  </button>
                </div>
              </div>
            ) : (
              <div className="mt-4 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-4">
                <p className="text-lg font-semibold text-[var(--app-text)]">
                  {detail.currentUserLicense ? 'Purchase complete' : 'No install yet'}
                </p>
                <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
                  Installing creates a workspace-local copy with source attribution and published version metadata.
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
                    onClick={() => installMutation.mutate()}
                    disabled={installMutation.isPending || !selectedWorkspaceId || !detail.latestVersion}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
                  >
                    {installMutation.isPending ? 'Installing...' : installLabel}
                  </button>
                </div>
              </div>
            )}
            {(installMutation.error || updateInstallMutation.error || removeInstallMutation.error || checkoutMutation.error) ? (
              <p className="mt-4 text-sm text-[var(--app-danger-text)]">
                {(installMutation.error || updateInstallMutation.error || removeInstallMutation.error || checkoutMutation.error) instanceof Error
                  ? ((installMutation.error || updateInstallMutation.error || removeInstallMutation.error || checkoutMutation.error) as Error).message
                  : 'Marketplace install action failed.'}
              </p>
            ) : null}
          </PageSection>
        </div>
      </div>
    </PageContainer>
  )
}

type PublishFormState = {
  ref?: string
  deckId: number | ''
  title: string
  slug: string
  summary: string
  description: string
  category: string
  tags: string
  coverImageUrl: string
  priceMode: 'free' | 'premium'
  priceCents: string
}

function emptyPublishForm(): PublishFormState {
  return {
    ref: undefined,
    deckId: '',
    title: '',
    slug: '',
    summary: '',
    description: '',
    category: '',
    tags: '',
    coverImageUrl: '',
    priceMode: 'free',
    priceCents: '',
  }
}

function listingToForm(listing: MarketplaceListingSummary): PublishFormState {
  return {
    ref: listing.id,
    deckId: listing.sourceDeckId,
    title: listing.title,
    slug: listing.slug,
    summary: listing.summary,
    description: listing.description,
    category: listing.category,
    tags: listing.tags.join(', '),
    coverImageUrl: listing.coverImageUrl,
    priceMode: listing.priceMode,
    priceCents: listing.priceCents > 0 ? String(listing.priceCents) : '',
  }
}

export function MarketplacePublishPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const [searchParams, setSearchParams] = useSearchParams()
  const [editorOpen, setEditorOpen] = useState(false)
  const [publishOpen, setPublishOpen] = useState(false)
  const [activeListing, setActiveListing] = useState<MarketplaceListingSummary | null>(null)
  const [changeSummary, setChangeSummary] = useState('')
  const [form, setForm] = useState<PublishFormState>(emptyPublishForm)

  const listingsQuery = useQuery({
    queryKey: ['marketplace-listings', 'mine'],
    queryFn: () => repository.fetchMarketplaceListings('mine'),
  })
  const decksQuery = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })
  const entitlementsQuery = useQuery({
    queryKey: ['entitlements'],
    queryFn: () => repository.fetchEntitlements(),
  })
  const creatorStatusQuery = useQuery({
    queryKey: ['marketplace-creator-account'],
    queryFn: () => repository.fetchMarketplaceCreatorAccountStatus(),
  })

  const canPublish = entitlementsQuery.data?.features.marketplacePublish ?? false

  useEffect(() => {
    const editRef = searchParams.get('edit')
    if (!editRef || !listingsQuery.data) return
    const listing = listingsQuery.data.find((entry) => entry.id === editRef)
    if (!listing) return
    setActiveListing(listing)
    setForm(listingToForm(listing))
    setEditorOpen(true)
  }, [listingsQuery.data, searchParams])

  const refreshListings = () => {
    queryClient.invalidateQueries({ queryKey: ['marketplace-listings'] })
    queryClient.invalidateQueries({ queryKey: ['marketplace-creator-account'] })
  }

  const creatorOnboardingMutation = useMutation({
    mutationFn: () => repository.startMarketplaceCreatorAccount(),
    onSuccess: (response) => {
      queryClient.invalidateQueries({ queryKey: ['marketplace-creator-account'] })
      if (response.account?.onboardingUrl) {
        window.location.assign(response.account.onboardingUrl)
        return
      }
      if (response.account?.dashboardUrl) {
        window.location.assign(response.account.dashboardUrl)
      }
    },
  })

  const createMutation = useMutation({
    mutationFn: (payload: CreateMarketplaceListingRequest) => repository.createMarketplaceListing(payload),
    onSuccess: () => {
      refreshListings()
      setEditorOpen(false)
      setForm(emptyPublishForm())
      setSearchParams({})
    },
  })

  const updateMutation = useMutation({
    mutationFn: ({ref, payload}: {ref: string; payload: UpdateMarketplaceListingRequest}) => repository.updateMarketplaceListing(ref, payload),
    onSuccess: () => {
      refreshListings()
      setEditorOpen(false)
      setActiveListing(null)
      setForm(emptyPublishForm())
      setSearchParams({})
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (ref: string) => repository.deleteMarketplaceListing(ref),
    onSuccess: () => {
      refreshListings()
      setActiveListing(null)
      setEditorOpen(false)
      setSearchParams({})
    },
  })

  const publishMutation = useMutation({
    mutationFn: ({ref, payload}: {ref: string; payload: {changeSummary?: string}}) => repository.publishMarketplaceListing(ref, payload),
    onSuccess: () => {
      refreshListings()
      setPublishOpen(false)
      setActiveListing(null)
      setChangeSummary('')
    },
  })

  const submitPayload = useMemo(() => {
    const base = {
      deckId: Number(form.deckId),
      title: form.title.trim(),
      slug: form.slug.trim(),
      summary: form.summary.trim(),
      description: form.description.trim(),
      category: form.category.trim(),
      tags: parseTags(form.tags),
      coverImageUrl: form.coverImageUrl.trim(),
      priceMode: form.priceMode,
      priceCents: Number(form.priceCents || '0'),
      currency: 'USD' as const,
    }
    return base
  }, [form])

  if (!canPublish) {
    return (
      <PageContainer className="space-y-4">
        <SectionHeader
          eyebrow="Publish"
          title="Marketplace publishing requires Pro or above."
          description="This phase ships creator listing management, but only Pro, Team, and Enterprise workspaces can create and publish listings."
        />
        <EmptyState
          title="Upgrade required"
          description="Browse the marketplace from the catalog page, or upgrade this workspace to publish free and premium listing metadata."
          action={
            <Link
              to="/marketplace"
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Browse marketplace
            </Link>
          }
        />
      </PageContainer>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <SectionHeader
        eyebrow="Publish"
        title="Publish versioned source decks to the marketplace."
        description="Listings stay separate from installs. Publishing creates explicit source versions, and premium decks now depend on creator payout setup before buyers can check out."
        action={
          <button
            type="button"
            onClick={() => {
              setActiveListing(null)
              setForm(emptyPublishForm())
              setEditorOpen(true)
            }}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
          >
            New listing
          </button>
        }
      />

      <PageSection className="p-5 sm:p-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div>
            <p className="text-[11px] uppercase tracking-[0.2em] text-[var(--app-muted)]">Creator payouts</p>
            <p className="mt-3 text-lg font-semibold text-[var(--app-text)]">
              {creatorStatusQuery.data?.canSellPremium ? 'Premium checkout enabled' : 'Finish creator setup to sell premium decks'}
            </p>
            <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">
              Premium marketplace listings need a creator payout account. If Stripe is configured, this button opens live Connect onboarding; local development still auto-enables the flow when Stripe is absent.
            </p>
            {creatorStatusQuery.data?.account ? (
              <div className="mt-4 flex flex-wrap gap-2">
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-xs text-[var(--app-text-soft)]">
                  Provider: {creatorStatusQuery.data.account.provider}
                </span>
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-xs text-[var(--app-text-soft)]">
                  Charges {creatorStatusQuery.data.account.chargesEnabled ? 'enabled' : 'pending'}
                </span>
                <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card-strong)] px-3 py-1 text-xs text-[var(--app-text-soft)]">
                  Payouts {creatorStatusQuery.data.account.payoutsEnabled ? 'enabled' : 'pending'}
                </span>
              </div>
            ) : null}
          </div>
          <button
            type="button"
            onClick={() => creatorOnboardingMutation.mutate()}
            disabled={creatorOnboardingMutation.isPending}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {creatorOnboardingMutation.isPending
              ? 'Setting up...'
              : creatorStatusQuery.data?.canSellPremium
                ? 'Open creator dashboard'
                : 'Start creator setup'}
          </button>
        </div>
        {creatorStatusQuery.isError || creatorOnboardingMutation.error ? (
          <p className="mt-4 text-sm text-[var(--app-danger-text)]">
            {(creatorOnboardingMutation.error || creatorStatusQuery.error) instanceof Error
              ? ((creatorOnboardingMutation.error || creatorStatusQuery.error) as Error).message
              : 'Failed to load creator payout status.'}
          </p>
        ) : null}
      </PageSection>

      {listingsQuery.isLoading ? (
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading your listings...</PageSection>
      ) : listingsQuery.isError ? (
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          {listingsQuery.error instanceof Error ? listingsQuery.error.message : 'Failed to load marketplace listings.'}
        </PageSection>
      ) : listingsQuery.data && listingsQuery.data.length > 0 ? (
        <div className="grid gap-4">
          {listingsQuery.data.map((listing) => (
            <ListingCard
              key={listing.id}
              listing={listing}
              action={
                <div className="flex flex-wrap gap-2">
                  <Link
                    to={`/marketplace/${listing.slug}`}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-semibold text-[var(--app-text)]"
                  >
                    View
                  </Link>
                  <button
                    type="button"
                    onClick={() => {
                      setActiveListing(listing)
                      setForm(listingToForm(listing))
                      setEditorOpen(true)
                    }}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-semibold text-[var(--app-text)]"
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setActiveListing(listing)
                      setChangeSummary('')
                      setPublishOpen(true)
                    }}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
                  >
                    {listing.latestVersionNumber > 0 ? 'Publish update' : 'Publish'}
                  </button>
                </div>
              }
            />
          ))}
        </div>
      ) : (
        <EmptyState
          title="No marketplace listings yet"
          description="Create the first listing from a source deck, publish an explicit version, and let free users install personal copies."
          action={
            <button
              type="button"
              onClick={() => setEditorOpen(true)}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              New listing
            </button>
          }
        />
      )}

      <Sheet open={editorOpen} onClose={() => {
        setEditorOpen(false)
        setActiveListing(null)
        setSearchParams({})
      }} title={activeListing ? 'Edit listing' : 'Create listing'}>
        <div className="space-y-4">
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Source deck</span>
            <select
              value={form.deckId}
              onChange={(event) => setForm((current) => ({...current, deckId: event.target.value ? Number(event.target.value) : ''}))}
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
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Title</span>
            <input
              value={form.title}
              onChange={(event) => setForm((current) => ({...current, title: event.target.value}))}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="USMLE Step 1 Foundations"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Slug</span>
            <input
              value={form.slug}
              onChange={(event) => setForm((current) => ({...current, slug: event.target.value}))}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="auto-generated if left blank"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Summary</span>
            <textarea
              value={form.summary}
              onChange={(event) => setForm((current) => ({...current, summary: event.target.value}))}
              rows={3}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Short marketplace teaser"
            />
          </label>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Description</span>
            <textarea
              value={form.description}
              onChange={(event) => setForm((current) => ({...current, description: event.target.value}))}
              rows={5}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Long-form description for the listing detail page"
            />
          </label>
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--app-text)]">Category</span>
              <input
                value={form.category}
                onChange={(event) => setForm((current) => ({...current, category: event.target.value}))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                placeholder="Medicine"
              />
            </label>
            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--app-text)]">Tags</span>
              <input
                value={form.tags}
                onChange={(event) => setForm((current) => ({...current, tags: event.target.value}))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                placeholder="anki, exam, med school"
              />
            </label>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--app-text)]">Price mode</span>
              <select
                value={form.priceMode}
                onChange={(event) => setForm((current) => ({...current, priceMode: event.target.value as 'free' | 'premium'}))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              >
                <option value="free">Free</option>
                <option value="premium">Premium</option>
              </select>
            </label>
            <label className="block space-y-2">
              <span className="text-sm font-medium text-[var(--app-text)]">Price cents</span>
              <input
                value={form.priceCents}
                onChange={(event) => setForm((current) => ({...current, priceCents: event.target.value.replace(/[^\d]/g, '')}))}
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
                placeholder="0"
              />
            </label>
          </div>
          {form.priceMode === 'premium' ? (
            <p className="text-sm leading-6 text-[var(--app-text-soft)]">
              Premium listings can be drafted anytime, but publishing them now requires creator payout setup to be complete.
            </p>
          ) : null}
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Cover image URL</span>
            <input
              value={form.coverImageUrl}
              onChange={(event) => setForm((current) => ({...current, coverImageUrl: event.target.value}))}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="https://..."
            />
          </label>
          {(createMutation.error || updateMutation.error || deleteMutation.error) ? (
            <p className="text-sm text-[var(--app-danger-text)]">
              {(createMutation.error || updateMutation.error || deleteMutation.error) instanceof Error
                ? ((createMutation.error || updateMutation.error || deleteMutation.error) as Error).message
                : 'Marketplace listing action failed.'}
            </p>
          ) : null}
          <div className="flex flex-col gap-3">
            <button
              type="button"
              onClick={() => {
                if (activeListing?.id) {
                  updateMutation.mutate({ref: activeListing.id, payload: submitPayload})
                  return
                }
                createMutation.mutate(submitPayload)
              }}
              disabled={
                createMutation.isPending ||
                updateMutation.isPending ||
                !submitPayload.deckId ||
                !submitPayload.title
              }
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
            >
              {createMutation.isPending || updateMutation.isPending
                ? 'Saving...'
                : activeListing
                  ? 'Save changes'
                  : 'Create listing'}
            </button>
            {activeListing ? (
              <button
                type="button"
                onClick={() => deleteMutation.mutate(activeListing.id)}
                disabled={deleteMutation.isPending}
                className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-semibold text-[var(--app-danger-text)] disabled:opacity-60"
              >
                {deleteMutation.isPending ? 'Deleting...' : 'Delete listing'}
              </button>
            ) : null}
          </div>
        </div>
      </Sheet>

      <Sheet open={publishOpen} onClose={() => setPublishOpen(false)} title="Publish marketplace version">
        <div className="space-y-4">
          <p className="text-sm leading-6 text-[var(--app-text-soft)]">
            Publishing creates an explicit source version for installs. Users who already installed an older version can then opt into a fresh-copy update.
          </p>
          <label className="block space-y-2">
            <span className="text-sm font-medium text-[var(--app-text)]">Change summary</span>
            <textarea
              value={changeSummary}
              onChange={(event) => setChangeSummary(event.target.value)}
              rows={4}
              className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
              placeholder="Added 120 anatomy cards and revised explanations."
            />
          </label>
          {publishMutation.isError ? (
            <p className="text-sm text-[var(--app-danger-text)]">
              {publishMutation.error instanceof Error ? publishMutation.error.message : 'Failed to publish listing.'}
            </p>
          ) : null}
          <button
            type="button"
            onClick={() => {
              if (!activeListing) return
              publishMutation.mutate({ref: activeListing.id, payload: {changeSummary}})
            }}
            disabled={publishMutation.isPending || !activeListing}
            className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {publishMutation.isPending ? 'Publishing...' : 'Publish version'}
          </button>
        </div>
      </Sheet>
    </PageContainer>
  )
}
