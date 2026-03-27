const API_BASE = (import.meta.env.VITE_API_BASE ?? "/api").replace(/\/$/, "");

export class APIError extends Error {
  status: number;
  code?: string;

  constructor(message: string, status: number, code?: string) {
    super(message);
    this.name = "APIError";
    this.status = status;
    this.code = code;
  }
}

async function requestJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    credentials: "include",
    ...init,
  });

  if (!res.ok) {
    if (
      res.status === 401 &&
      typeof window !== "undefined" &&
      !window.location.pathname.startsWith("/login")
    ) {
      window.location.assign("/login");
    }
    const contentType = res.headers.get("content-type") || "";
    if (contentType.includes("application/json")) {
      const payload = (await res.json().catch(() => null)) as {
        message?: string;
        code?: string;
      } | null;
      throw new APIError(
        payload?.message || "Request failed",
        res.status,
        payload?.code,
      );
    }
    const text = await res.text();
    throw new APIError(text || "Request failed", res.status);
  }

  if (res.status === 204) {
    return undefined as T;
  }

  return res.json();
}

export interface Deck {
  id: number;
  name: string;
  parentId?: number;
  cardIds: number[];
  dueToday: number;
  dueReviewBacklog: number;
  newCardsPerDay: number;
  reviewsPerDay: number;
  priorityOrder: number;
  newCardsPaused: boolean;
  noteCount: number;
  cardCount: number;
  canDelete: boolean;
  deleteBlockedReason?: string;
  analytics: DeckStudyAnalytics;
}

export interface DeckStats {
  deckId: number;
  newCards: number;
  learning: number;
  review: number;
  relearning: number;
  suspended: number;
  buried: number;
  totalCards: number;
  dueToday: number;
  dueReviewBacklog: number;
}

export interface Note {
  id: number;
  typeId: string;
  fieldVals: Record<string, string>;
  tags: string[];
  createdAt: string;
  modifiedAt: string;
  deckId?: number;
  cardCount?: number;
}

export interface Card {
  id: number;
  noteId: number;
  deckId: number;
  templateName: string;
  ordinal: number;
  front: string;
  back: string;
  flag: number;
  marked: boolean;
  suspended: boolean;
}

export interface CardTemplate {
  name: string;
  qFmt: string;
  aFmt: string;
  styling: string;
  ifFieldNonEmpty?: string;
  isCloze: boolean;
  deckOverride?: string;
  browserQFmt?: string;
  browserAFmt?: string;
}

export interface FieldOptions {
  font?: string;
  fontSize?: number;
  rtl?: boolean;
  htmlEditor?: boolean;
}

export interface NoteType {
  name: string;
  fields: string[];
  templates: CardTemplate[];
  sortFieldIndex: number;
  fieldOptions?: Record<string, FieldOptions>;
}

export interface CreateDeckRequest {
  name: string;
}

export type ImportSource = "auto" | "native" | "anki" | "quizlet";

export interface ImportFileRequest {
  file: File;
  source?: ImportSource;
  deckName?: string;
  noteType?: string;
  format?: string;
}

export interface ImportNotesResponse {
  imported: number;
  skipped: number;
  source: string;
  format: string;
  decksCreated?: string[];
  errors?: string[];
}

export interface CreateNoteRequest {
  typeId: string;
  deckId: number;
  fieldVals: Record<string, string>;
  tags?: string[];
}

export interface UpdateNoteRequest {
  typeId: string;
  deckId: number;
  fieldVals: Record<string, string>;
  tags?: string[];
}

export interface UpdateDeckRequest {
  name?: string;
  newCardsPerDay?: number;
  reviewsPerDay?: number;
  priorityOrder?: number;
}

export interface CreateTemplateRequest {
  name: string;
  sourceTemplateName?: string;
}

export interface AnswerCardRequest {
  rating: number; // 1=Again, 2=Hard, 3=Good, 4=Easy
  timeTakenMs?: number; // Time spent on the card in milliseconds
}

export interface CheckDuplicateRequest {
  typeId: string;
  fieldName: string;
  value: string;
  deckId?: number;
}

export interface NoteBrief {
  id: number;
  typeId: string;
  fieldVals: Record<string, string>;
  deckId?: number;
}

export interface DuplicateResult {
  isDuplicate: boolean;
  duplicates?: NoteBrief[];
}

export interface RecentDeckNoteSummary {
  noteId: number;
  noteType: string;
  createdAt: string;
  modifiedAt: string;
  tags: string[];
  fieldPreview: string;
  cardCountInDeck: number;
}

export interface DeckNotesResponse {
  notes: RecentDeckNoteSummary[];
}

export interface NoteListItem {
  id: number;
  typeId: string;
  fieldVals: Record<string, string>;
  fieldPreview: string;
  tags: string[];
  createdAt: string;
  modifiedAt: string;
  deckId?: number;
  deckName?: string;
  cardCount: number;
}

export interface ListNotesResponse {
  notes: NoteListItem[];
  total: number;
  nextCursor?: string;
  prevCursor?: string;
}

export interface ListNotesParams {
  deckId?: number;
  q?: string;
  typeId?: string;
  tag?: string;
  limit?: number;
  cursor?: string;
}

export interface PlanLimits {
  maxDecks: number;
  maxNotes: number;
  maxCardsTotal: number;
  maxSharedDecks: number;
  maxSyncDevices: number;
  maxWorkspaces: number;
}

export interface EntitlementUsage {
  decks: number;
  notes: number;
  cardsTotal: number;
  sharedDecks: number;
  syncDevices: number;
  workspaces: number;
}

export interface EntitlementFeatures {
  googleLogin: boolean;
  accountBacked: boolean;
  sync: boolean;
  shareDecks: boolean;
  organizations: boolean;
  studyGroups: boolean;
  marketplacePublish: boolean;
  enterprise: boolean;
}

export interface Entitlements {
  plan: "guest" | "free" | "pro" | "team" | "enterprise";
  limits: PlanLimits;
  usage: EntitlementUsage;
  features: EntitlementFeatures;
}

export interface DashboardResponse {
  totalDecks: number;
  totalNotes: number;
  dueToday: number;
  plan: Entitlements["plan"];
  usage: EntitlementUsage;
  limits: PlanLimits;
  features: EntitlementFeatures;
  studyAnalytics: StudyAnalyticsOverview;
  recentNotes: NoteListItem[];
}

export interface GenerateAICardSuggestionsRequest {
  sourceText: string;
  noteType: string;
  existingFieldVals?: Record<string, string>;
  maxSuggestions?: number;
}

export interface AICardSuggestion {
  title: string;
  rationale: string;
  fieldVals: Record<string, string>;
}

export interface AICardSuggestionsResponse {
  suggestions: AICardSuggestion[];
  provider: string;
  model?: string;
}

export interface StudyAnalyticsOverview {
  sessions7d: number;
  cardsReviewed7d: number;
  minutesStudied7d: number;
  focusSessions7d: number;
  focusMinutes7d: number;
  currentStreak: number;
  lastStudiedAt?: string;
  answerBreakdown: StudyAnswerBreakdown;
  dailyActivity: StudyAnalyticsDay[];
  recentSessions: StudySessionSummary[];
}

export interface DeckStudyAnalytics {
  sessions7d: number;
  cardsReviewed7d: number;
  minutesStudied7d: number;
  averageCardsPerSession7d: number;
  againCount7d: number;
  hardCount7d: number;
  goodCount7d: number;
  easyCount7d: number;
  lastStudiedAt?: string;
}

export interface StudyAnswerBreakdown {
  again: number;
  hard: number;
  good: number;
  easy: number;
}

export interface StudyAnalyticsDay {
  date: string;
  sessions: number;
  cardsReviewed: number;
  minutesStudied: number;
}

export interface StudySessionSummary {
  id: string;
  deckId?: number;
  deckName?: string;
  mode: string;
  protocol?: string;
  targetMinutes?: number;
  breakMinutes?: number;
  status: "active" | "completed" | "abandoned";
  cardsReviewed: number;
  minutesStudied: number;
  againCount: number;
  hardCount: number;
  goodCount: number;
  easyCount: number;
  startedAt: string;
  endedAt?: string;
  updatedAt: string;
}

export interface StudySession {
  id: string;
  userId: string;
  workspaceId: string;
  deckId?: number;
  mode: string;
  protocol?: string;
  targetMinutes?: number;
  breakMinutes?: number;
  status: "active" | "completed" | "abandoned";
  startedAt: string;
  endedAt?: string;
  cardsReviewed: number;
  againCount: number;
  hardCount: number;
  goodCount: number;
  easyCount: number;
  createdAt: string;
  updatedAt: string;
}

export interface CreateStudySessionRequest {
  deckId?: number;
  mode?: string;
  protocol?: "pomodoro" | "deep-focus" | "custom";
  targetMinutes?: number;
  breakMinutes?: number;
}

export interface UpdateStudySessionRequest {
  status?: "active" | "completed" | "abandoned";
  cardsReviewed?: number;
  againCount?: number;
  hardCount?: number;
  goodCount?: number;
  easyCount?: number;
  endedAt?: string;
}

export interface AccountUser {
  id: string;
  email: string;
  displayName: string;
  avatarUrl?: string;
  onboarding: boolean;
}

export interface WorkspaceSession {
  id: string;
  name: string;
  slug: string;
  collectionId: string;
  ownerUserId?: string;
  organizationId?: string;
}

export interface Organization {
  id: string;
  name: string;
  slug: string;
  createdAt: string;
  updatedAt: string;
}

export interface OrganizationMember {
  id: string;
  organizationId: string;
  userId?: string;
  email: string;
  role: "read" | "edit" | "admin" | "owner";
  status: "invited" | "active" | "removed" | "declined" | "expired";
  inviteToken?: string;
  inviteExpiresAt?: string;
  joinedAt?: string;
  removedAt?: string;
  createdAt: string;
}

export interface Subscription {
  id: string;
  workspaceId?: string;
  organizationId?: string;
  plan: Entitlements["plan"];
  status: string;
  provider?: string;
  providerCustomerId?: string;
  providerSubscriptionId?: string;
  providerSubscriptionItemId?: string;
  providerCheckoutSessionId?: string;
  scheduledPlan?: Entitlements["plan"];
  currentPeriodEnd?: string;
  cancelAtPeriodEnd?: boolean;
  billedQuantity?: number;
  createdAt: string;
  updatedAt: string;
}

export interface OrganizationDetail {
  organization: Organization;
  workspace?: WorkspaceSession;
  subscription?: Subscription;
  membership: OrganizationMember;
  members: OrganizationMember[];
  canManagePlan: boolean;
  canManageMembers: boolean;
  canEdit: boolean;
}

export interface AuthSessionResponse {
  authenticated: boolean;
  googleAuthConfigured: boolean;
  otpAuthEnabled: boolean;
  user?: AccountUser;
  workspace?: WorkspaceSession;
  organization?: Organization;
  organizationMember?: OrganizationMember;
  subscription?: Subscription;
  entitlements: Entitlements;
}

export interface BillingCheckoutResponse {
  provider: string;
  plan: UpdateWorkspacePlanRequest["plan"];
  checkoutUrl?: string;
  completed: boolean;
  message?: string;
  session?: AuthSessionResponse;
  subscription?: Subscription;
}

export interface BillingPortalResponse {
  provider: string;
  url: string;
  message?: string;
}

export interface BillingCheckoutSyncResponse {
  provider: string;
  completed: boolean;
  message?: string;
  session?: AuthSessionResponse;
  subscription?: Subscription;
}

export interface OTPRequestResponse {
  ok: boolean;
  expiresAt: string;
  retryAfterSeconds: number;
  delivery?: "email" | "dev-inline";
  devCode?: string;
}

export interface StudyGroupMember {
  id: string;
  studyGroupId: string;
  userId?: string;
  email: string;
  role: "owner" | "admin" | "edit" | "read";
  status: "invited" | "active" | "removed" | "declined" | "expired";
  inviteToken?: string;
  inviteExpiresAt?: string;
  joinedAt?: string;
  removedAt?: string;
  createdAt: string;
}

export interface StudyGroupVersion {
  id: string;
  studyGroupId: string;
  versionNumber: number;
  sourceDeckId: number;
  publishedByUserId: string;
  changeSummary: string;
  noteCount: number;
  cardCount: number;
  createdAt: string;
}

export interface StudyGroupInstall {
  id: string;
  studyGroupId: string;
  studyGroupMemberId: string;
  destinationWorkspaceId: string;
  installedDeckId: number;
  installedDeckName?: string;
  sourceVersionNumber: number;
  status: "active" | "superseded" | "removed";
  syncState: "clean" | "forked";
  supersededByInstallId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface StudyGroupLeaderboardEntry {
  userId?: string;
  email: string;
  displayName?: string;
  sessions7d: number;
  minutes7d: number;
  reviews7d: number;
}

export interface StudyGroupDashboard {
  memberCount: number;
  activeMembers7d: number;
  activeInstalls: number;
  reviewsToday: number;
  reviews7d: number;
  sessions7d: number;
  minutesStudied7d: number;
  latestVersionNumber: number;
  latestVersionAdoption: number;
  latestVersionAdoptionPercent: number;
  dailyActivity: StudyAnalyticsDay[];
  leaderboard: StudyGroupLeaderboardEntry[];
}

export interface StudyGroupEvent {
  id: string;
  studyGroupId: string;
  actorUserId?: string;
  eventType: string;
  payload: string;
  createdAt: string;
}

export interface StudyGroupSummary {
  id: string;
  name: string;
  description: string;
  sourceDeckId: number;
  sourceDeckName: string;
  role: string;
  membershipStatus: string;
  latestVersionNumber: number;
  memberCount: number;
  activeMembers7d: number;
  updateAvailable: boolean;
  currentUserInstall?: StudyGroupInstall;
}

export interface StudyGroupDetail {
  group: {
    id: string;
    workspaceId: string;
    primaryDeckId: number;
    name: string;
    description: string;
    visibility: string;
    joinPolicy: string;
    createdByUserId: string;
    createdAt: string;
    updatedAt: string;
  };
  role: string;
  membershipStatus: string;
  sourceDeckName: string;
  latestVersion?: StudyGroupVersion;
  versions: StudyGroupVersion[];
  members: StudyGroupMember[];
  currentUserInstall?: StudyGroupInstall;
  updateAvailable: boolean;
  canEdit: boolean;
  canManageMembers: boolean;
  canInvite: boolean;
  canPublishVersion: boolean;
  dashboard: StudyGroupDashboard;
  recentEvents: StudyGroupEvent[];
  availableWorkspaces: WorkspaceSession[];
}

export interface CreateStudyGroupRequest {
  name: string;
  description?: string;
  primaryDeckId: number;
  visibility?: string;
  joinPolicy?: string;
}

export interface UpdateStudyGroupRequest {
  name?: string;
  description?: string;
  visibility?: string;
  joinPolicy?: string;
}

export interface InviteStudyGroupMemberRequest {
  email: string;
  role: "admin" | "edit" | "read";
}

export interface UpdateStudyGroupMemberRequest {
  role?: "admin" | "edit" | "read" | "owner";
  status?: "invited" | "active" | "removed" | "declined" | "expired";
}

export interface JoinStudyGroupRequest {
  token: string;
  destinationWorkspaceId: string;
  installLatest: boolean;
}

export interface PublishStudyGroupVersionRequest {
  changeSummary?: string;
}

export interface InstallStudyGroupDeckRequest {
  destinationWorkspaceId: string;
}

export interface UpdateStudyGroupInstallRequest {
  destinationWorkspaceId?: string;
}

export interface CreateOrganizationRequest {
  name: string;
  slug?: string;
}

export interface AddOrganizationMemberRequest {
  email: string;
  role: "read" | "edit" | "admin";
}

export interface UpdateOrganizationRequest {
  name: string;
  slug?: string;
}

export interface UpdateOrganizationMemberRequest {
  role?: "read" | "edit" | "admin" | "owner";
  status?: "invited" | "active" | "removed" | "declined" | "expired";
}

export interface JoinOrganizationRequest {
  token: string;
}

export interface UpdateWorkspacePlanRequest {
  plan: "free" | "pro" | "team" | "enterprise";
}

export interface MarketplaceInstall {
  id: string;
  listingId: string;
  workspaceId: string;
  installedByUserId: string;
  installedDeckId: number;
  installedDeckName?: string;
  sourceVersionNumber: number;
  status: "active" | "superseded" | "removed";
  supersededByInstallId?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MarketplaceCreatorAccount {
  id: string;
  userId: string;
  workspaceId: string;
  provider: string;
  providerAccountId: string;
  onboardingStatus: string;
  detailsSubmitted: boolean;
  chargesEnabled: boolean;
  payoutsEnabled: boolean;
  onboardingUrl?: string;
  dashboardUrl?: string;
  onboardingCompletedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MarketplaceOrder {
  id: string;
  listingId: string;
  listingVersionNumber: number;
  buyerUserId: string;
  buyerWorkspaceId: string;
  creatorUserId: string;
  creatorAccountId?: string;
  provider: string;
  providerCheckoutSessionId: string;
  providerPaymentIntentId?: string;
  status: string;
  amountCents: number;
  currency: string;
  platformFeeCents: number;
  creatorAmountCents: number;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface MarketplaceLicense {
  id: string;
  listingId: string;
  buyerUserId: string;
  orderId: string;
  status: string;
  grantedVersionNumber: number;
  createdAt: string;
  updatedAt: string;
}

export interface MarketplaceListingVersion {
  id: string;
  listingId: string;
  versionNumber: number;
  sourceDeckId: number;
  publishedByUserId: string;
  changeSummary: string;
  noteCount: number;
  cardCount: number;
  createdAt: string;
}

export interface MarketplaceListingSummary {
  id: string;
  slug: string;
  title: string;
  summary: string;
  description: string;
  category: string;
  tags: string[];
  coverImageUrl: string;
  creatorUserId: string;
  creatorDisplayName?: string;
  creatorEmail?: string;
  workspaceId: string;
  sourceDeckId: number;
  sourceDeckName: string;
  priceMode: "free" | "premium";
  priceCents: number;
  currency: string;
  status: "draft" | "published" | "archived";
  installCount: number;
  latestVersionNumber: number;
  canEdit: boolean;
  updateAvailable: boolean;
  currentUserLicense?: MarketplaceLicense;
  currentUserInstall?: MarketplaceInstall;
  createdAt: string;
  updatedAt: string;
}

export interface MarketplaceListingDetail {
  listing: MarketplaceListingSummary;
  latestVersion?: MarketplaceListingVersion;
  versions: MarketplaceListingVersion[];
  currentUserLicense?: MarketplaceLicense;
  currentUserInstall?: MarketplaceInstall;
  updateAvailable: boolean;
  canEdit: boolean;
  canPublish: boolean;
  availableWorkspaces: WorkspaceSession[];
}

export interface MarketplaceCreatorAccountStatusResponse {
  account?: MarketplaceCreatorAccount;
  provider: string;
  canSellPremium: boolean;
}

export interface MarketplaceCheckoutResponse {
  provider: string;
  checkoutUrl?: string;
  completed: boolean;
  order: MarketplaceOrder;
  license?: MarketplaceLicense;
}

export interface CreateMarketplaceListingRequest {
  deckId: number;
  title: string;
  slug?: string;
  summary: string;
  description: string;
  category: string;
  tags: string[];
  coverImageUrl: string;
  priceMode?: "free" | "premium";
  priceCents?: number;
  currency?: string;
}

export interface UpdateMarketplaceListingRequest {
  deckId: number;
  title: string;
  slug?: string;
  summary: string;
  description: string;
  category: string;
  tags: string[];
  coverImageUrl: string;
  priceMode?: "free" | "premium";
  priceCents?: number;
  currency?: string;
}

export interface PublishMarketplaceListingRequest {
  changeSummary?: string;
}

export interface InstallMarketplaceListingRequest {
  destinationWorkspaceId?: string;
}

export interface UpdateMarketplaceInstallRequest {
  destinationWorkspaceId?: string;
}

// Deck endpoints
export async function fetchDecks(): Promise<Deck[]> {
  return requestJSON(`${API_BASE}/decks`);
}

export async function fetchDashboard(): Promise<DashboardResponse> {
  return requestJSON(`${API_BASE}/dashboard`);
}

export async function fetchStudyAnalyticsOverview(): Promise<StudyAnalyticsOverview> {
  return requestJSON(`${API_BASE}/analytics/overview`);
}

export async function createDeck(req: CreateDeckRequest): Promise<Deck> {
  return requestJSON(`${API_BASE}/decks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchDeck(
  id: number,
): Promise<{ deck: Deck; stats: DeckStats }> {
  return requestJSON(`${API_BASE}/decks/${id}`);
}

export async function fetchDeckStats(id: number): Promise<DeckStats> {
  return requestJSON(`${API_BASE}/decks/${id}/stats`);
}

export async function fetchDeckNotes(
  deckId: number,
  limit: number = 20,
  cursor?: string,
): Promise<DeckNotesResponse> {
  const params = new URLSearchParams({ limit: String(limit) });
  if (cursor) params.set("cursor", cursor);
  return requestJSON(`${API_BASE}/decks/${deckId}/notes?${params.toString()}`);
}

export async function updateDeck(
  id: number,
  req: UpdateDeckRequest,
): Promise<Deck> {
  return requestJSON(`${API_BASE}/decks/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteDeck(id: number): Promise<void> {
  return requestJSON(`${API_BASE}/decks/${id}`, {
    method: "DELETE",
  });
}

export async function importNotesFile(
  req: ImportFileRequest,
): Promise<ImportNotesResponse> {
  const formData = new FormData();
  formData.append("file", req.file);
  if (req.source) formData.append("source", req.source);
  if (req.deckName) formData.append("deckName", req.deckName);
  if (req.noteType) formData.append("noteType", req.noteType);
  if (req.format) formData.append("format", req.format);

  const res = await fetch(`${API_BASE}/import`, {
    method: "POST",
    credentials: "include",
    body: formData,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || "Failed to import file");
  }
  return res.json();
}

// Note Type endpoints
export async function fetchNoteTypes(): Promise<NoteType[]> {
  return requestJSON(`${API_BASE}/note-types`);
}

export async function fetchNoteType(name: string): Promise<NoteType> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(name)}`);
}

// Field management
export interface AddFieldRequest {
  fieldName: string;
  position?: number;
}

export interface RenameFieldRequest {
  oldName: string;
  newName: string;
}

export interface RemoveFieldRequest {
  fieldName: string;
}

export interface ReorderFieldsRequest {
  fields: string[];
}

export interface FieldsResponse {
  message: string;
  fields: string[];
}

export interface UpdateTemplateRequest {
  name?: string;
  qFmt?: string;
  aFmt?: string;
  styling?: string;
  ifFieldNonEmpty?: string;
  deckOverride?: string;
  browserQFmt?: string;
  browserAFmt?: string;
}

export interface TemplatesResponse {
  message: string;
  templates: CardTemplate[];
}

export async function createTemplate(
  noteTypeName: string,
  req: CreateTemplateRequest,
): Promise<TemplatesResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function addField(
  noteTypeName: string,
  req: AddFieldRequest,
): Promise<FieldsResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function renameField(
  noteTypeName: string,
  req: RenameFieldRequest,
): Promise<FieldsResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/rename`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function removeField(
  noteTypeName: string,
  req: RemoveFieldRequest,
): Promise<FieldsResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`,
    {
      method: "DELETE",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function reorderFields(
  noteTypeName: string,
  req: ReorderFieldsRequest,
): Promise<FieldsResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/reorder`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export interface SetSortFieldRequest {
  fieldIndex: number;
}

export async function setSortField(
  noteTypeName: string,
  req: SetSortFieldRequest,
): Promise<{ message: string; sortFieldIndex: number; sortFieldName: string }> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/sort-field`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function setFieldOptions(
  noteTypeName: string,
  fieldName: string,
  options: FieldOptions,
): Promise<{ message: string; fieldOptions: Record<string, FieldOptions> }> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/options`,
    {
      method: "PUT",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ fieldName, options }),
    },
  );
}

export async function updateTemplate(
  noteTypeName: string,
  templateName: string,
  req: UpdateTemplateRequest,
): Promise<TemplatesResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates/${encodeURIComponent(templateName)}`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function deleteTemplate(
  noteTypeName: string,
  templateName: string,
): Promise<TemplatesResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates/${encodeURIComponent(templateName)}`,
    {
      method: "DELETE",
    },
  );
}

// Note endpoints
export async function createNote(
  req: CreateNoteRequest,
): Promise<{ note: Note; cards: Card[] }> {
  return requestJSON(`${API_BASE}/notes`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchNote(id: number): Promise<Note> {
  return requestJSON(`${API_BASE}/notes/${id}`);
}

export async function fetchNotes(
  params: ListNotesParams = {},
): Promise<ListNotesResponse> {
  const query = new URLSearchParams();
  if (params.deckId) query.set("deckId", String(params.deckId));
  if (params.q) query.set("q", params.q);
  if (params.typeId) query.set("typeId", params.typeId);
  if (params.tag) query.set("tag", params.tag);
  if (params.limit) query.set("limit", String(params.limit));
  if (params.cursor) query.set("cursor", params.cursor);
  const suffix = query.toString() ? `?${query.toString()}` : "";
  return requestJSON(`${API_BASE}/notes${suffix}`);
}

export async function updateNote(
  id: number,
  req: UpdateNoteRequest,
): Promise<{ note: Note; cards: Card[] }> {
  return requestJSON(`${API_BASE}/notes/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function generateAICardSuggestions(
  req: GenerateAICardSuggestionsRequest,
): Promise<AICardSuggestionsResponse> {
  return requestJSON(`${API_BASE}/ai/card-suggestions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteNote(id: number): Promise<void> {
  return requestJSON(`${API_BASE}/notes/${id}`, {
    method: "DELETE",
  });
}

export async function checkDuplicate(
  req: CheckDuplicateRequest,
): Promise<DuplicateResult> {
  return requestJSON(`${API_BASE}/notes/check-duplicate`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

// Card endpoints
export async function fetchDueCards(
  deckId: number,
  limit: number = 10,
): Promise<Card[]> {
  return requestJSON(`${API_BASE}/decks/${deckId}/due?limit=${limit}`);
}

export async function fetchCard(id: number): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}`);
}

export async function answerCard(
  id: number,
  req: AnswerCardRequest,
): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}/answer`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export interface UpdateCardRequest {
  flag?: number; // 0-7 color flags
  marked?: boolean; // toggle marked status
  suspended?: boolean; // toggle suspended status
}

export async function updateCard(
  id: number,
  req: UpdateCardRequest,
): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function createStudySession(
  req: CreateStudySessionRequest = {},
): Promise<StudySession> {
  return requestJSON(`${API_BASE}/study-sessions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateStudySession(
  id: string,
  req: UpdateStudySessionRequest,
): Promise<StudySession> {
  return requestJSON(`${API_BASE}/study-sessions/${encodeURIComponent(id)}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

// Empty cards endpoints
export interface EmptyCardInfo {
  cardId: number;
  noteId: number;
  deckId: number;
  templateName: string;
  ordinal: number;
  front: string;
  back: string;
  reason: string;
}

export interface EmptyCardsResponse {
  count: number;
  emptyCards: EmptyCardInfo[];
}

export interface DeleteEmptyCardsRequest {
  cardIds: number[];
}

export interface DeleteEmptyCardsResponse {
  deleted: number;
  failed?: string[];
}

export async function findEmptyCards(): Promise<EmptyCardsResponse> {
  return requestJSON(`${API_BASE}/cards/empty`);
}

export async function deleteEmptyCards(
  req: DeleteEmptyCardsRequest,
): Promise<DeleteEmptyCardsResponse> {
  return requestJSON(`${API_BASE}/cards/empty/delete`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchStudyGroups(): Promise<StudyGroupSummary[]> {
  return requestJSON(`${API_BASE}/study-groups`);
}

export async function createStudyGroup(
  req: CreateStudyGroupRequest,
): Promise<StudyGroupDetail> {
  return requestJSON(`${API_BASE}/study-groups`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchStudyGroup(id: string): Promise<StudyGroupDetail> {
  return requestJSON(`${API_BASE}/study-groups/${id}`);
}

export async function updateStudyGroup(
  id: string,
  req: UpdateStudyGroupRequest,
): Promise<StudyGroupDetail> {
  return requestJSON(`${API_BASE}/study-groups/${id}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteStudyGroup(id: string): Promise<void> {
  return requestJSON(`${API_BASE}/study-groups/${id}`, {
    method: "DELETE",
  });
}

export async function inviteStudyGroupMember(
  id: string,
  req: InviteStudyGroupMemberRequest,
): Promise<StudyGroupMember> {
  return requestJSON(`${API_BASE}/study-groups/${id}/members`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateStudyGroupMember(
  id: string,
  memberId: string,
  req: UpdateStudyGroupMemberRequest,
): Promise<StudyGroupMember> {
  return requestJSON(`${API_BASE}/study-groups/${id}/members/${memberId}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteStudyGroupMember(
  id: string,
  memberId: string,
): Promise<void> {
  return requestJSON(`${API_BASE}/study-groups/${id}/members/${memberId}`, {
    method: "DELETE",
  });
}

export async function joinStudyGroup(
  req: JoinStudyGroupRequest,
): Promise<StudyGroupDetail> {
  return requestJSON(`${API_BASE}/study-groups/join`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchStudyGroupVersions(
  id: string,
): Promise<StudyGroupVersion[]> {
  return requestJSON(`${API_BASE}/study-groups/${id}/versions`);
}

export async function publishStudyGroupVersion(
  id: string,
  req: PublishStudyGroupVersionRequest,
): Promise<StudyGroupVersion> {
  return requestJSON(`${API_BASE}/study-groups/${id}/versions`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function installStudyGroupDeck(
  id: string,
  req: InstallStudyGroupDeckRequest,
): Promise<StudyGroupInstall> {
  return requestJSON(`${API_BASE}/study-groups/${id}/installs`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateStudyGroupInstall(
  id: string,
  installId: string,
  req: UpdateStudyGroupInstallRequest = {},
): Promise<StudyGroupInstall> {
  return requestJSON(
    `${API_BASE}/study-groups/${id}/installs/${installId}/update`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function removeStudyGroupInstall(
  id: string,
  installId: string,
): Promise<void> {
  return requestJSON(`${API_BASE}/study-groups/${id}/installs/${installId}`, {
    method: "DELETE",
  });
}

export async function fetchStudyGroupDashboard(
  id: string,
): Promise<StudyGroupDashboard> {
  return requestJSON(`${API_BASE}/study-groups/${id}/dashboard`);
}

export async function fetchOrganization(
  orgId: string,
): Promise<OrganizationDetail> {
  return requestJSON(`${API_BASE}/orgs/${encodeURIComponent(orgId)}`);
}

export async function createOrganization(
  req: CreateOrganizationRequest,
): Promise<OrganizationDetail> {
  return requestJSON(`${API_BASE}/orgs`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateOrganization(
  orgId: string,
  req: UpdateOrganizationRequest,
): Promise<OrganizationDetail> {
  return requestJSON(`${API_BASE}/orgs/${encodeURIComponent(orgId)}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteOrganization(orgId: string): Promise<void> {
  return requestJSON(`${API_BASE}/orgs/${encodeURIComponent(orgId)}`, {
    method: "DELETE",
  });
}

export async function addOrganizationMember(
  orgId: string,
  req: AddOrganizationMemberRequest,
): Promise<{ member: OrganizationMember }> {
  return requestJSON(`${API_BASE}/orgs/${encodeURIComponent(orgId)}/members`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateOrganizationMember(
  orgId: string,
  memberId: string,
  req: UpdateOrganizationMemberRequest,
): Promise<{ member: OrganizationMember }> {
  return requestJSON(
    `${API_BASE}/orgs/${encodeURIComponent(orgId)}/members/${encodeURIComponent(memberId)}`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function deleteOrganizationMember(
  orgId: string,
  memberId: string,
): Promise<void> {
  return requestJSON(
    `${API_BASE}/orgs/${encodeURIComponent(orgId)}/members/${encodeURIComponent(memberId)}`,
    {
      method: "DELETE",
    },
  );
}

export async function joinOrganization(
  req: JoinOrganizationRequest,
): Promise<OrganizationDetail> {
  return requestJSON(`${API_BASE}/orgs/join`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateWorkspacePlan(
  workspaceId: string,
  req: UpdateWorkspacePlanRequest,
): Promise<AuthSessionResponse> {
  return requestJSON(
    `${API_BASE}/workspaces/${encodeURIComponent(workspaceId)}/plan`,
    {
      method: "PATCH",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function completeOnboardingPlan(
  req: UpdateWorkspacePlanRequest,
): Promise<AuthSessionResponse> {
  return requestJSON(`${API_BASE}/onboarding/plan`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function startBillingCheckout(req: {
  plan: UpdateWorkspacePlanRequest["plan"];
}): Promise<BillingCheckoutResponse> {
  return requestJSON(`${API_BASE}/billing/checkout`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function openBillingPortal(req?: {
  plan?: UpdateWorkspacePlanRequest["plan"];
}): Promise<BillingPortalResponse> {
  return requestJSON(`${API_BASE}/billing/portal`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req ?? {}),
  });
}

export async function syncBillingCheckoutSession(
  sessionId: string,
): Promise<BillingCheckoutSyncResponse> {
  return requestJSON(
    `${API_BASE}/billing/checkout/sessions/${encodeURIComponent(sessionId)}/sync`,
    {
      method: "POST",
    },
  );
}

export async function fetchMarketplaceListings(
  scope?: "mine",
): Promise<MarketplaceListingSummary[]> {
  const params = new URLSearchParams();
  if (scope) params.set("scope", scope);
  const query = params.toString();
  return requestJSON(
    `${API_BASE}/marketplace/listings${query ? `?${query}` : ""}`,
  );
}

export async function fetchMarketplaceCreatorAccountStatus(): Promise<MarketplaceCreatorAccountStatusResponse> {
  return requestJSON(`${API_BASE}/marketplace/creator-account/status`);
}

export async function startMarketplaceCreatorAccount(): Promise<MarketplaceCreatorAccountStatusResponse> {
  return requestJSON(`${API_BASE}/marketplace/creator-account/start`, {
    method: "POST",
  });
}

export async function createMarketplaceListing(
  req: CreateMarketplaceListingRequest,
): Promise<MarketplaceListingDetail> {
  return requestJSON(`${API_BASE}/marketplace/listings`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function fetchMarketplaceListing(
  ref: string,
): Promise<MarketplaceListingDetail> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}`);
}

export async function updateMarketplaceListing(
  ref: string,
  req: UpdateMarketplaceListingRequest,
): Promise<MarketplaceListingDetail> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function deleteMarketplaceListing(ref: string): Promise<void> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}`, {
    method: "DELETE",
  });
}

export async function publishMarketplaceListing(
  ref: string,
  req: PublishMarketplaceListingRequest = {},
): Promise<MarketplaceListingVersion> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}/publish`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function checkoutMarketplaceListing(
  ref: string,
): Promise<MarketplaceCheckoutResponse> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}/checkout`, {
    method: "POST",
  });
}

export async function syncMarketplaceCheckoutSession(
  sessionId: string,
): Promise<MarketplaceCheckoutResponse> {
  return requestJSON(
    `${API_BASE}/marketplace/checkout/sessions/${encodeURIComponent(sessionId)}/sync`,
    {
      method: "POST",
    },
  );
}

export async function installMarketplaceListing(
  ref: string,
  req: InstallMarketplaceListingRequest = {},
): Promise<MarketplaceInstall> {
  return requestJSON(`${API_BASE}/marketplace/listings/${ref}/installs`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(req),
  });
}

export async function updateMarketplaceInstall(
  ref: string,
  installId: string,
  req: UpdateMarketplaceInstallRequest = {},
): Promise<MarketplaceInstall> {
  return requestJSON(
    `${API_BASE}/marketplace/listings/${ref}/installs/${installId}/update`,
    {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(req),
    },
  );
}

export async function removeMarketplaceInstall(
  ref: string,
  installId: string,
): Promise<void> {
  return requestJSON(
    `${API_BASE}/marketplace/listings/${ref}/installs/${installId}`,
    {
      method: "DELETE",
    },
  );
}

// Backup endpoints
export async function createBackup(): Promise<{
  message: string;
  backupPath: string;
}> {
  return requestJSON(`${API_BASE}/backups`, {
    method: "POST",
  });
}

export async function listBackups(): Promise<
  Array<{ path: string; filename: string; size: number; modified: string }>
> {
  return requestJSON(`${API_BASE}/backups`);
}

// Health check
export async function healthCheck(): Promise<{
  status: string;
  service: string;
  version: string;
}> {
  return requestJSON(`${API_BASE}/health`);
}

export async function fetchSession(): Promise<AuthSessionResponse> {
  return requestJSON(`${API_BASE}/auth/session`);
}

export async function requestOTP(email: string): Promise<OTPRequestResponse> {
  return requestJSON(`${API_BASE}/auth/otp/request`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email }),
  });
}

export async function verifyOTP(
  email: string,
  code: string,
): Promise<AuthSessionResponse> {
  return requestJSON(`${API_BASE}/auth/otp/verify`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, code }),
  });
}

export async function fetchEntitlements(): Promise<Entitlements> {
  return requestJSON(`${API_BASE}/entitlements`);
}

export async function logout(): Promise<{ ok: boolean }> {
  return requestJSON(`${API_BASE}/auth/logout`, {
    method: "POST",
  });
}
