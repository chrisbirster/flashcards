import { createContext, useContext, type ReactNode } from "react";
import {
  addField,
  addOrganizationMember,
  answerCard,
  openBillingPortal,
  checkDuplicate,
  checkoutMarketplaceListing,
  completeOnboardingPlan,
  createBackup,
  createDeck,
  createMarketplaceListing,
  createNote,
  createOrganization,
  createStudyGroup,
  createStudySession,
  createTemplate,
  deleteStudyGroup,
  deleteDeck,
  deleteEmptyCards,
  deleteMarketplaceListing,
  deleteNote,
  deleteOrganization,
  deleteOrganizationMember,
  deleteStudyGroupMember,
  deleteTemplate,
  fetchCard,
  fetchDashboard,
  fetchDeck,
  fetchDeckNotes,
  fetchDecks,
  fetchDeckStats,
  fetchDueCards,
  fetchEntitlements,
  fetchMarketplaceCreatorAccountStatus,
  fetchMarketplaceListing,
  fetchMarketplaceListings,
  fetchNote,
  fetchNotes,
  fetchNoteType,
  fetchNoteTypes,
  fetchOrganization,
  fetchSession,
  fetchStudyAnalyticsOverview,
  fetchStudyGroup,
  fetchStudyGroupDashboard,
  fetchStudyGroups,
  fetchStudyGroupVersions,
  findEmptyCards,
  generateAICardSuggestions,
  healthCheck,
  importNotesFile,
  installMarketplaceListing,
  installStudyGroupDeck,
  inviteStudyGroupMember,
  joinOrganization,
  joinStudyGroup,
  listBackups,
  logout,
  publishMarketplaceListing,
  publishStudyGroupVersion,
  removeField,
  removeMarketplaceInstall,
  removeStudyGroupInstall,
  renameField,
  reorderFields,
  requestOTP,
  setFieldOptions,
  setSortField,
  startBillingCheckout,
  startMarketplaceCreatorAccount,
  syncBillingCheckoutSession,
  syncMarketplaceCheckoutSession,
  updateCard,
  updateDeck,
  updateMarketplaceInstall,
  updateMarketplaceListing,
  updateNote,
  updateOrganization,
  updateOrganizationMember,
  updateStudyGroup,
  updateStudyGroupInstall,
  updateStudyGroupMember,
  updateStudySession,
  updateTemplate,
  updateWorkspacePlan,
  verifyOTP,
  type AddFieldRequest,
  type AddOrganizationMemberRequest,
  type AICardSuggestionsResponse,
  type AnswerCardRequest,
  type AuthSessionResponse,
  type BillingCheckoutResponse,
  type BillingCheckoutSyncResponse,
  type BillingPortalResponse,
  type Card,
  type CheckDuplicateRequest,
  type CreateDeckRequest,
  type CreateMarketplaceListingRequest,
  type CreateNoteRequest,
  type CreateOrganizationRequest,
  type CreateStudyGroupRequest,
  type CreateStudySessionRequest,
  type CreateTemplateRequest,
  type DashboardResponse,
  type Deck,
  type DeckNotesResponse,
  type DeckStats,
  type DeleteEmptyCardsRequest,
  type DeleteEmptyCardsResponse,
  type DuplicateResult,
  type EmptyCardsResponse,
  type Entitlements,
  type FieldOptions,
  type FieldsResponse,
  type GenerateAICardSuggestionsRequest,
  type ImportFileRequest,
  type ImportNotesResponse,
  type InstallMarketplaceListingRequest,
  type InstallStudyGroupDeckRequest,
  type InviteStudyGroupMemberRequest,
  type JoinOrganizationRequest,
  type JoinStudyGroupRequest,
  type ListNotesParams,
  type ListNotesResponse,
  type MarketplaceCheckoutResponse,
  type MarketplaceCreatorAccountStatusResponse,
  type MarketplaceInstall,
  type MarketplaceListingDetail,
  type MarketplaceListingSummary,
  type MarketplaceListingVersion,
  type Note,
  type NoteType,
  type OrganizationDetail,
  type OrganizationMember,
  type OTPRequestResponse,
  type PublishMarketplaceListingRequest,
  type PublishStudyGroupVersionRequest,
  type RenameFieldRequest,
  type ReorderFieldsRequest,
  type StudyAnalyticsOverview,
  type StudyGroupDashboard,
  type StudyGroupDetail,
  type StudyGroupInstall,
  type StudyGroupMember,
  type StudyGroupSummary,
  type StudyGroupVersion,
  type StudySession,
  type TemplatesResponse,
  type UpdateCardRequest,
  type UpdateDeckRequest,
  type UpdateMarketplaceInstallRequest,
  type UpdateMarketplaceListingRequest,
  type UpdateNoteRequest,
  type UpdateOrganizationMemberRequest,
  type UpdateOrganizationRequest,
  type UpdateStudyGroupInstallRequest,
  type UpdateStudyGroupMemberRequest,
  type UpdateStudyGroupRequest,
  type UpdateStudySessionRequest,
  type UpdateTemplateRequest,
  type UpdateWorkspacePlanRequest,
} from "#/lib/api";

export interface AppRepository {
  fetchDashboard(): Promise<DashboardResponse>;
  fetchStudyAnalyticsOverview(): Promise<StudyAnalyticsOverview>;
  fetchDecks(): Promise<Deck[]>;
  createDeck(req: CreateDeckRequest): Promise<Deck>;
  fetchDeck(id: number): Promise<{ deck: Deck; stats: DeckStats }>;
  fetchDeckStats(id: number): Promise<DeckStats>;
  fetchDeckNotes(
    deckId: number,
    limit?: number,
    cursor?: string,
  ): Promise<DeckNotesResponse>;
  updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck>;
  deleteDeck(id: number): Promise<void>;
  importNotesFile(req: ImportFileRequest): Promise<ImportNotesResponse>;
  fetchNoteTypes(): Promise<NoteType[]>;
  fetchNoteType(name: string): Promise<NoteType>;
  addField(noteTypeName: string, req: AddFieldRequest): Promise<FieldsResponse>;
  renameField(
    noteTypeName: string,
    req: RenameFieldRequest,
  ): Promise<FieldsResponse>;
  removeField(
    noteTypeName: string,
    req: { fieldName: string },
  ): Promise<FieldsResponse>;
  reorderFields(
    noteTypeName: string,
    req: ReorderFieldsRequest,
  ): Promise<FieldsResponse>;
  setSortField(
    noteTypeName: string,
    req: { fieldIndex: number },
  ): Promise<{
    message: string;
    sortFieldIndex: number;
    sortFieldName: string;
  }>;
  setFieldOptions(
    noteTypeName: string,
    fieldName: string,
    options: FieldOptions,
  ): Promise<{ message: string; fieldOptions: Record<string, FieldOptions> }>;
  createTemplate(
    noteTypeName: string,
    req: CreateTemplateRequest,
  ): Promise<TemplatesResponse>;
  updateTemplate(
    noteTypeName: string,
    templateName: string,
    req: UpdateTemplateRequest,
  ): Promise<TemplatesResponse>;
  deleteTemplate(
    noteTypeName: string,
    templateName: string,
  ): Promise<TemplatesResponse>;
  createNote(req: CreateNoteRequest): Promise<{ note: Note; cards: Card[] }>;
  generateAICardSuggestions(
    req: GenerateAICardSuggestionsRequest,
  ): Promise<AICardSuggestionsResponse>;
  fetchNote(id: number): Promise<Note>;
  fetchNotes(params?: ListNotesParams): Promise<ListNotesResponse>;
  updateNote(
    id: number,
    req: UpdateNoteRequest,
  ): Promise<{ note: Note; cards: Card[] }>;
  deleteNote(id: number): Promise<void>;
  checkDuplicate(req: CheckDuplicateRequest): Promise<DuplicateResult>;
  fetchDueCards(deckId: number, limit?: number): Promise<Card[]>;
  fetchCard(id: number): Promise<Card>;
  answerCard(id: number, req: AnswerCardRequest): Promise<Card>;
  updateCard(id: number, req: UpdateCardRequest): Promise<Card>;
  createStudySession(req?: CreateStudySessionRequest): Promise<StudySession>;
  updateStudySession(
    id: string,
    req: UpdateStudySessionRequest,
  ): Promise<StudySession>;
  findEmptyCards(): Promise<EmptyCardsResponse>;
  deleteEmptyCards(
    req: DeleteEmptyCardsRequest,
  ): Promise<DeleteEmptyCardsResponse>;
  createBackup(): Promise<{ message: string; backupPath: string }>;
  listBackups(): Promise<
    Array<{ path: string; filename: string; size: number; modified: string }>
  >;
  healthCheck(): Promise<{ status: string; service: string; version: string }>;
  fetchSession(): Promise<AuthSessionResponse>;
  requestOTP(email: string): Promise<OTPRequestResponse>;
  verifyOTP(email: string, code: string): Promise<AuthSessionResponse>;
  completeOnboardingPlan(
    req: UpdateWorkspacePlanRequest,
  ): Promise<AuthSessionResponse>;
  startBillingCheckout(req: {
    plan: UpdateWorkspacePlanRequest["plan"];
  }): Promise<BillingCheckoutResponse>;
  openBillingPortal(req?: {
    plan?: UpdateWorkspacePlanRequest["plan"];
  }): Promise<BillingPortalResponse>;
  syncBillingCheckoutSession(
    sessionId: string,
  ): Promise<BillingCheckoutSyncResponse>;
  fetchEntitlements(): Promise<Entitlements>;
  createOrganization(req: CreateOrganizationRequest): Promise<OrganizationDetail>;
  fetchOrganization(orgId: string): Promise<OrganizationDetail>;
  updateOrganization(
    orgId: string,
    req: UpdateOrganizationRequest,
  ): Promise<OrganizationDetail>;
  deleteOrganization(orgId: string): Promise<void>;
  addOrganizationMember(
    orgId: string,
    req: AddOrganizationMemberRequest,
  ): Promise<{ member: OrganizationMember }>;
  updateOrganizationMember(
    orgId: string,
    memberId: string,
    req: UpdateOrganizationMemberRequest,
  ): Promise<{ member: OrganizationMember }>;
  deleteOrganizationMember(orgId: string, memberId: string): Promise<void>;
  joinOrganization(req: JoinOrganizationRequest): Promise<OrganizationDetail>;
  updateWorkspacePlan(
    workspaceId: string,
    req: UpdateWorkspacePlanRequest,
  ): Promise<AuthSessionResponse>;
  fetchStudyGroups(): Promise<StudyGroupSummary[]>;
  createStudyGroup(req: CreateStudyGroupRequest): Promise<StudyGroupDetail>;
  fetchStudyGroup(id: string): Promise<StudyGroupDetail>;
  updateStudyGroup(
    id: string,
    req: UpdateStudyGroupRequest,
  ): Promise<StudyGroupDetail>;
  deleteStudyGroup(id: string): Promise<void>;
  inviteStudyGroupMember(
    id: string,
    req: InviteStudyGroupMemberRequest,
  ): Promise<StudyGroupMember>;
  updateStudyGroupMember(
    id: string,
    memberId: string,
    req: UpdateStudyGroupMemberRequest,
  ): Promise<StudyGroupMember>;
  deleteStudyGroupMember(id: string, memberId: string): Promise<void>;
  joinStudyGroup(req: JoinStudyGroupRequest): Promise<StudyGroupDetail>;
  fetchStudyGroupVersions(id: string): Promise<StudyGroupVersion[]>;
  publishStudyGroupVersion(
    id: string,
    req: PublishStudyGroupVersionRequest,
  ): Promise<StudyGroupVersion>;
  installStudyGroupDeck(
    id: string,
    req: InstallStudyGroupDeckRequest,
  ): Promise<StudyGroupInstall>;
  updateStudyGroupInstall(
    id: string,
    installId: string,
    req?: UpdateStudyGroupInstallRequest,
  ): Promise<StudyGroupInstall>;
  removeStudyGroupInstall(id: string, installId: string): Promise<void>;
  fetchStudyGroupDashboard(id: string): Promise<StudyGroupDashboard>;
  fetchMarketplaceListings(
    scope?: "mine",
  ): Promise<MarketplaceListingSummary[]>;
  fetchMarketplaceCreatorAccountStatus(): Promise<MarketplaceCreatorAccountStatusResponse>;
  startMarketplaceCreatorAccount(): Promise<MarketplaceCreatorAccountStatusResponse>;
  createMarketplaceListing(
    req: CreateMarketplaceListingRequest,
  ): Promise<MarketplaceListingDetail>;
  fetchMarketplaceListing(ref: string): Promise<MarketplaceListingDetail>;
  updateMarketplaceListing(
    ref: string,
    req: UpdateMarketplaceListingRequest,
  ): Promise<MarketplaceListingDetail>;
  deleteMarketplaceListing(ref: string): Promise<void>;
  publishMarketplaceListing(
    ref: string,
    req?: PublishMarketplaceListingRequest,
  ): Promise<MarketplaceListingVersion>;
  checkoutMarketplaceListing(ref: string): Promise<MarketplaceCheckoutResponse>;
  syncMarketplaceCheckoutSession(
    sessionId: string,
  ): Promise<MarketplaceCheckoutResponse>;
  installMarketplaceListing(
    ref: string,
    req?: InstallMarketplaceListingRequest,
  ): Promise<MarketplaceInstall>;
  updateMarketplaceInstall(
    ref: string,
    installId: string,
    req?: UpdateMarketplaceInstallRequest,
  ): Promise<MarketplaceInstall>;
  removeMarketplaceInstall(ref: string, installId: string): Promise<void>;
  logout(): Promise<{ ok: boolean }>;
}

export const remoteRepository: AppRepository = {
  fetchDashboard,
  fetchStudyAnalyticsOverview,
  fetchDecks,
  createDeck,
  fetchDeck,
  fetchDeckStats,
  fetchDeckNotes,
  updateDeck,
  deleteDeck,
  importNotesFile,
  fetchNoteTypes,
  fetchNoteType,
  addField,
  renameField,
  removeField,
  reorderFields,
  setSortField,
  setFieldOptions,
  createTemplate,
  updateTemplate,
  deleteTemplate,
  createNote,
  generateAICardSuggestions,
  fetchNote,
  fetchNotes,
  updateNote,
  deleteNote,
  checkDuplicate,
  fetchDueCards,
  fetchCard,
  answerCard,
  updateCard,
  createStudySession,
  updateStudySession,
  findEmptyCards,
  deleteEmptyCards,
  createBackup,
  listBackups,
  healthCheck,
  fetchSession,
  requestOTP,
  verifyOTP,
  completeOnboardingPlan,
  startBillingCheckout,
  openBillingPortal,
  syncBillingCheckoutSession,
  fetchEntitlements,
  createOrganization,
  fetchOrganization,
  updateOrganization,
  deleteOrganization,
  addOrganizationMember,
  updateOrganizationMember,
  deleteOrganizationMember,
  joinOrganization,
  updateWorkspacePlan,
  fetchStudyGroups,
  createStudyGroup,
  fetchStudyGroup,
  updateStudyGroup,
  deleteStudyGroup,
  inviteStudyGroupMember,
  updateStudyGroupMember,
  deleteStudyGroupMember,
  joinStudyGroup,
  fetchStudyGroupVersions,
  publishStudyGroupVersion,
  installStudyGroupDeck,
  updateStudyGroupInstall,
  removeStudyGroupInstall,
  fetchStudyGroupDashboard,
  fetchMarketplaceListings,
  fetchMarketplaceCreatorAccountStatus,
  startMarketplaceCreatorAccount,
  createMarketplaceListing,
  fetchMarketplaceListing,
  updateMarketplaceListing,
  deleteMarketplaceListing,
  publishMarketplaceListing,
  checkoutMarketplaceListing,
  syncMarketplaceCheckoutSession,
  installMarketplaceListing,
  updateMarketplaceInstall,
  removeMarketplaceInstall,
  logout,
};

function notImplemented<T>(method: string): Promise<T> {
  return Promise.reject(
    new Error(`${method} is not implemented for the local repository yet.`),
  );
}

export function createLocalRepository(): AppRepository {
  return {
    fetchDashboard: () => notImplemented("fetchDashboard"),
    fetchStudyAnalyticsOverview: () =>
      notImplemented("fetchStudyAnalyticsOverview"),
    fetchDecks: () => notImplemented("fetchDecks"),
    createDeck: () => notImplemented("createDeck"),
    fetchDeck: () => notImplemented("fetchDeck"),
    fetchDeckStats: () => notImplemented("fetchDeckStats"),
    fetchDeckNotes: () => notImplemented("fetchDeckNotes"),
    updateDeck: () => notImplemented("updateDeck"),
    deleteDeck: () => notImplemented("deleteDeck"),
    importNotesFile: () => notImplemented("importNotesFile"),
    fetchNoteTypes: () => notImplemented("fetchNoteTypes"),
    fetchNoteType: () => notImplemented("fetchNoteType"),
    addField: () => notImplemented("addField"),
    renameField: () => notImplemented("renameField"),
    removeField: () => notImplemented("removeField"),
    reorderFields: () => notImplemented("reorderFields"),
    setSortField: () => notImplemented("setSortField"),
    setFieldOptions: () => notImplemented("setFieldOptions"),
    createTemplate: () => notImplemented("createTemplate"),
    updateTemplate: () => notImplemented("updateTemplate"),
    deleteTemplate: () => notImplemented("deleteTemplate"),
    createNote: () => notImplemented("createNote"),
    generateAICardSuggestions: () =>
      notImplemented("generateAICardSuggestions"),
    fetchNote: () => notImplemented("fetchNote"),
    fetchNotes: () => notImplemented("fetchNotes"),
    updateNote: () => notImplemented("updateNote"),
    deleteNote: () => notImplemented("deleteNote"),
    checkDuplicate: () => notImplemented("checkDuplicate"),
    fetchDueCards: () => notImplemented("fetchDueCards"),
    fetchCard: () => notImplemented("fetchCard"),
    answerCard: () => notImplemented("answerCard"),
    updateCard: () => notImplemented("updateCard"),
    createStudySession: () => notImplemented("createStudySession"),
    updateStudySession: () => notImplemented("updateStudySession"),
    findEmptyCards: () => notImplemented("findEmptyCards"),
    deleteEmptyCards: () => notImplemented("deleteEmptyCards"),
    createBackup: () => notImplemented("createBackup"),
    listBackups: () => notImplemented("listBackups"),
    healthCheck: () =>
      Promise.resolve({
        status: "ok",
        service: "vutadex-local",
        version: "local",
      }),
    fetchSession: () =>
      Promise.resolve({
        authenticated: false,
        googleAuthConfigured: false,
        otpAuthEnabled: true,
        entitlements: {
          plan: "guest",
          limits: {
            maxDecks: 2,
            maxNotes: 10,
            maxCardsTotal: 100,
            maxSharedDecks: 0,
            maxSyncDevices: 0,
            maxWorkspaces: 1,
          },
          usage: {
            decks: 0,
            notes: 0,
            cardsTotal: 0,
            sharedDecks: 0,
            syncDevices: 0,
            workspaces: 1,
          },
          features: {
            googleLogin: false,
            accountBacked: false,
            sync: false,
            shareDecks: false,
            organizations: false,
            studyGroups: false,
            marketplacePublish: false,
            enterprise: false,
          },
        },
      }),
    requestOTP: () =>
      Promise.resolve({
        ok: true,
        expiresAt: new Date().toISOString(),
        retryAfterSeconds: 60,
      }),
    verifyOTP: () => notImplemented("verifyOTP"),
    completeOnboardingPlan: () => notImplemented("completeOnboardingPlan"),
    startBillingCheckout: () => notImplemented("startBillingCheckout"),
    openBillingPortal: () => notImplemented("openBillingPortal"),
    syncBillingCheckoutSession: () =>
      notImplemented("syncBillingCheckoutSession"),
    fetchEntitlements: () =>
      Promise.resolve({
        plan: "guest",
        limits: {
          maxDecks: 2,
          maxNotes: 10,
          maxCardsTotal: 100,
          maxSharedDecks: 0,
          maxSyncDevices: 0,
          maxWorkspaces: 1,
        },
        usage: {
          decks: 0,
          notes: 0,
          cardsTotal: 0,
          sharedDecks: 0,
          syncDevices: 0,
          workspaces: 1,
        },
        features: {
          googleLogin: false,
          accountBacked: false,
          sync: false,
          shareDecks: false,
          organizations: false,
          studyGroups: false,
          marketplacePublish: false,
          enterprise: false,
        },
      }),
    createOrganization: () => notImplemented("createOrganization"),
    fetchOrganization: () => notImplemented("fetchOrganization"),
    updateOrganization: () => notImplemented("updateOrganization"),
    deleteOrganization: () => notImplemented("deleteOrganization"),
    addOrganizationMember: () => notImplemented("addOrganizationMember"),
    updateOrganizationMember: () =>
      notImplemented("updateOrganizationMember"),
    deleteOrganizationMember: () =>
      notImplemented("deleteOrganizationMember"),
    joinOrganization: () => notImplemented("joinOrganization"),
    updateWorkspacePlan: () => notImplemented("updateWorkspacePlan"),
    fetchStudyGroups: () => notImplemented("fetchStudyGroups"),
    createStudyGroup: () => notImplemented("createStudyGroup"),
    fetchStudyGroup: () => notImplemented("fetchStudyGroup"),
    updateStudyGroup: () => notImplemented("updateStudyGroup"),
    deleteStudyGroup: () => notImplemented("deleteStudyGroup"),
    inviteStudyGroupMember: () => notImplemented("inviteStudyGroupMember"),
    updateStudyGroupMember: () => notImplemented("updateStudyGroupMember"),
    deleteStudyGroupMember: () => notImplemented("deleteStudyGroupMember"),
    joinStudyGroup: () => notImplemented("joinStudyGroup"),
    fetchStudyGroupVersions: () => notImplemented("fetchStudyGroupVersions"),
    publishStudyGroupVersion: () => notImplemented("publishStudyGroupVersion"),
    installStudyGroupDeck: () => notImplemented("installStudyGroupDeck"),
    updateStudyGroupInstall: () => notImplemented("updateStudyGroupInstall"),
    removeStudyGroupInstall: () => notImplemented("removeStudyGroupInstall"),
    fetchStudyGroupDashboard: () => notImplemented("fetchStudyGroupDashboard"),
    fetchMarketplaceListings: () => notImplemented("fetchMarketplaceListings"),
    fetchMarketplaceCreatorAccountStatus: () =>
      notImplemented("fetchMarketplaceCreatorAccountStatus"),
    startMarketplaceCreatorAccount: () =>
      notImplemented("startMarketplaceCreatorAccount"),
    createMarketplaceListing: () => notImplemented("createMarketplaceListing"),
    fetchMarketplaceListing: () => notImplemented("fetchMarketplaceListing"),
    updateMarketplaceListing: () => notImplemented("updateMarketplaceListing"),
    deleteMarketplaceListing: () => notImplemented("deleteMarketplaceListing"),
    publishMarketplaceListing: () =>
      notImplemented("publishMarketplaceListing"),
    checkoutMarketplaceListing: () =>
      notImplemented("checkoutMarketplaceListing"),
    syncMarketplaceCheckoutSession: () =>
      notImplemented("syncMarketplaceCheckoutSession"),
    installMarketplaceListing: () =>
      notImplemented("installMarketplaceListing"),
    updateMarketplaceInstall: () => notImplemented("updateMarketplaceInstall"),
    removeMarketplaceInstall: () => notImplemented("removeMarketplaceInstall"),
    logout: () => Promise.resolve({ ok: true }),
  };
}

const AppRepositoryContext = createContext<AppRepository>(remoteRepository);

export function AppRepositoryProvider({
  children,
  repository = remoteRepository,
}: {
  children: ReactNode;
  repository?: AppRepository;
}) {
  return (
    <AppRepositoryContext.Provider value={repository}>
      {children}
    </AppRepositoryContext.Provider>
  );
}

export function useAppRepository() {
  return useContext(AppRepositoryContext);
}
