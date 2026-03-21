import { createContext, useContext, type ReactNode } from 'react'
import {
  addField,
  answerCard,
  checkDuplicate,
  createBackup,
  fetchDashboard,
  createDeck,
  createTemplate,
  createNote,
  deleteDeck,
  deleteNote,
  deleteTemplate,
  deleteEmptyCards,
  fetchCard,
  fetchDeck,
  fetchDeckNotes,
  fetchDecks,
  fetchDeckStats,
  fetchDueCards,
  fetchEntitlements,
  fetchNote,
  fetchNotes,
  fetchNoteType,
  fetchNoteTypes,
  fetchSession,
  requestOTP,
  verifyOTP,
  findEmptyCards,
  healthCheck,
  importNotesFile,
  listBackups,
  logout,
  removeField,
  renameField,
  reorderFields,
  setFieldOptions,
  setSortField,
  updateDeck,
  updateNote,
  updateCard,
  updateTemplate,
  type AddFieldRequest,
  type AnswerCardRequest,
  type Card,
  type CheckDuplicateRequest,
  type CreateTemplateRequest,
  type CreateDeckRequest,
  type CreateNoteRequest,
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
  type ImportFileRequest,
  type ImportNotesResponse,
  type ListNotesParams,
  type ListNotesResponse,
  type Note,
  type NoteType,
  type RenameFieldRequest,
  type ReorderFieldsRequest,
  type TemplatesResponse,
  type UpdateCardRequest,
  type UpdateDeckRequest,
  type UpdateNoteRequest,
  type UpdateTemplateRequest,
  type AuthSessionResponse,
  type OTPRequestResponse,
} from '#/lib/api'

export interface AppRepository {
  fetchDashboard(): Promise<DashboardResponse>
  fetchDecks(): Promise<Deck[]>
  createDeck(req: CreateDeckRequest): Promise<Deck>
  fetchDeck(id: number): Promise<{deck: Deck; stats: DeckStats}>
  fetchDeckStats(id: number): Promise<DeckStats>
  fetchDeckNotes(deckId: number, limit?: number, cursor?: string): Promise<DeckNotesResponse>
  updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck>
  deleteDeck(id: number): Promise<void>
  importNotesFile(req: ImportFileRequest): Promise<ImportNotesResponse>
  fetchNoteTypes(): Promise<NoteType[]>
  fetchNoteType(name: string): Promise<NoteType>
  addField(noteTypeName: string, req: AddFieldRequest): Promise<FieldsResponse>
  renameField(noteTypeName: string, req: RenameFieldRequest): Promise<FieldsResponse>
  removeField(noteTypeName: string, req: {fieldName: string}): Promise<FieldsResponse>
  reorderFields(noteTypeName: string, req: ReorderFieldsRequest): Promise<FieldsResponse>
  setSortField(noteTypeName: string, req: {fieldIndex: number}): Promise<{message: string; sortFieldIndex: number; sortFieldName: string}>
  setFieldOptions(noteTypeName: string, fieldName: string, options: FieldOptions): Promise<{message: string; fieldOptions: Record<string, FieldOptions>}>
  createTemplate(noteTypeName: string, req: CreateTemplateRequest): Promise<TemplatesResponse>
  updateTemplate(noteTypeName: string, templateName: string, req: UpdateTemplateRequest): Promise<TemplatesResponse>
  deleteTemplate(noteTypeName: string, templateName: string): Promise<TemplatesResponse>
  createNote(req: CreateNoteRequest): Promise<{note: Note; cards: Card[]}>
  fetchNote(id: number): Promise<Note>
  fetchNotes(params?: ListNotesParams): Promise<ListNotesResponse>
  updateNote(id: number, req: UpdateNoteRequest): Promise<{note: Note; cards: Card[]}>
  deleteNote(id: number): Promise<void>
  checkDuplicate(req: CheckDuplicateRequest): Promise<DuplicateResult>
  fetchDueCards(deckId: number, limit?: number): Promise<Card[]>
  fetchCard(id: number): Promise<Card>
  answerCard(id: number, req: AnswerCardRequest): Promise<Card>
  updateCard(id: number, req: UpdateCardRequest): Promise<Card>
  findEmptyCards(): Promise<EmptyCardsResponse>
  deleteEmptyCards(req: DeleteEmptyCardsRequest): Promise<DeleteEmptyCardsResponse>
  createBackup(): Promise<{message: string; backupPath: string}>
  listBackups(): Promise<Array<{path: string; filename: string; size: number; modified: string}>>
  healthCheck(): Promise<{status: string; service: string; version: string}>
  fetchSession(): Promise<AuthSessionResponse>
  requestOTP(email: string): Promise<OTPRequestResponse>
  verifyOTP(email: string, code: string): Promise<AuthSessionResponse>
  fetchEntitlements(): Promise<Entitlements>
  logout(): Promise<{ok: boolean}>
}

export const remoteRepository: AppRepository = {
  fetchDashboard,
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
  fetchNote,
  fetchNotes,
  updateNote,
  deleteNote,
  checkDuplicate,
  fetchDueCards,
  fetchCard,
  answerCard,
  updateCard,
  findEmptyCards,
  deleteEmptyCards,
  createBackup,
  listBackups,
  healthCheck,
  fetchSession,
  requestOTP,
  verifyOTP,
  fetchEntitlements,
  logout,
}

function notImplemented<T>(method: string): Promise<T> {
  return Promise.reject(new Error(`${method} is not implemented for the local repository yet.`))
}

export function createLocalRepository(): AppRepository {
  return {
    fetchDashboard: () => notImplemented('fetchDashboard'),
    fetchDecks: () => notImplemented('fetchDecks'),
    createDeck: () => notImplemented('createDeck'),
    fetchDeck: () => notImplemented('fetchDeck'),
    fetchDeckStats: () => notImplemented('fetchDeckStats'),
    fetchDeckNotes: () => notImplemented('fetchDeckNotes'),
    updateDeck: () => notImplemented('updateDeck'),
    deleteDeck: () => notImplemented('deleteDeck'),
    importNotesFile: () => notImplemented('importNotesFile'),
    fetchNoteTypes: () => notImplemented('fetchNoteTypes'),
    fetchNoteType: () => notImplemented('fetchNoteType'),
    addField: () => notImplemented('addField'),
    renameField: () => notImplemented('renameField'),
    removeField: () => notImplemented('removeField'),
    reorderFields: () => notImplemented('reorderFields'),
    setSortField: () => notImplemented('setSortField'),
    setFieldOptions: () => notImplemented('setFieldOptions'),
    createTemplate: () => notImplemented('createTemplate'),
    updateTemplate: () => notImplemented('updateTemplate'),
    deleteTemplate: () => notImplemented('deleteTemplate'),
    createNote: () => notImplemented('createNote'),
    fetchNote: () => notImplemented('fetchNote'),
    fetchNotes: () => notImplemented('fetchNotes'),
    updateNote: () => notImplemented('updateNote'),
    deleteNote: () => notImplemented('deleteNote'),
    checkDuplicate: () => notImplemented('checkDuplicate'),
    fetchDueCards: () => notImplemented('fetchDueCards'),
    fetchCard: () => notImplemented('fetchCard'),
    answerCard: () => notImplemented('answerCard'),
    updateCard: () => notImplemented('updateCard'),
    findEmptyCards: () => notImplemented('findEmptyCards'),
    deleteEmptyCards: () => notImplemented('deleteEmptyCards'),
    createBackup: () => notImplemented('createBackup'),
    listBackups: () => notImplemented('listBackups'),
    healthCheck: () => Promise.resolve({status: 'ok', service: 'vutadex-local', version: 'local'}),
    fetchSession: () =>
      Promise.resolve({
        authenticated: false,
        googleAuthConfigured: false,
        otpAuthEnabled: true,
        entitlements: {
          plan: 'guest',
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
    requestOTP: () => Promise.resolve({ok: true, expiresAt: new Date().toISOString(), retryAfterSeconds: 60}),
    verifyOTP: () => notImplemented('verifyOTP'),
    fetchEntitlements: () =>
      Promise.resolve({
        plan: 'guest',
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
    logout: () => Promise.resolve({ok: true}),
  }
}

const AppRepositoryContext = createContext<AppRepository>(remoteRepository)

export function AppRepositoryProvider({
  children,
  repository = remoteRepository,
}: {
  children: ReactNode
  repository?: AppRepository
}) {
  return <AppRepositoryContext.Provider value={repository}>{children}</AppRepositoryContext.Provider>
}

export function useAppRepository() {
  return useContext(AppRepositoryContext)
}
