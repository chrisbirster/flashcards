import { createContext, useContext, type ReactNode } from 'react'
import {
  addField,
  answerCard,
  checkDuplicate,
  createBackup,
  createDeck,
  createNote,
  deleteEmptyCards,
  fetchCard,
  fetchDeck,
  fetchDeckNotes,
  fetchDecks,
  fetchDeckStats,
  fetchDueCards,
  fetchEntitlements,
  fetchNote,
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
  updateCard,
  updateTemplate,
  type AddFieldRequest,
  type AnswerCardRequest,
  type Card,
  type CheckDuplicateRequest,
  type CreateDeckRequest,
  type CreateNoteRequest,
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
  type Note,
  type NoteType,
  type RenameFieldRequest,
  type ReorderFieldsRequest,
  type TemplatesResponse,
  type UpdateCardRequest,
  type UpdateTemplateRequest,
  type AuthSessionResponse,
  type OTPRequestResponse,
} from '#/lib/api'

export interface AppRepository {
  fetchDecks(): Promise<Deck[]>
  createDeck(req: CreateDeckRequest): Promise<Deck>
  fetchDeck(id: number): Promise<{deck: Deck; stats: DeckStats}>
  fetchDeckStats(id: number): Promise<DeckStats>
  fetchDeckNotes(deckId: number, limit?: number, cursor?: string): Promise<DeckNotesResponse>
  importNotesFile(req: ImportFileRequest): Promise<ImportNotesResponse>
  fetchNoteTypes(): Promise<NoteType[]>
  fetchNoteType(name: string): Promise<NoteType>
  addField(noteTypeName: string, req: AddFieldRequest): Promise<FieldsResponse>
  renameField(noteTypeName: string, req: RenameFieldRequest): Promise<FieldsResponse>
  removeField(noteTypeName: string, req: {fieldName: string}): Promise<FieldsResponse>
  reorderFields(noteTypeName: string, req: ReorderFieldsRequest): Promise<FieldsResponse>
  setSortField(noteTypeName: string, req: {fieldIndex: number}): Promise<{message: string; sortFieldIndex: number; sortFieldName: string}>
  setFieldOptions(noteTypeName: string, fieldName: string, options: FieldOptions): Promise<{message: string; fieldOptions: Record<string, FieldOptions>}>
  updateTemplate(noteTypeName: string, templateName: string, req: UpdateTemplateRequest): Promise<TemplatesResponse>
  createNote(req: CreateNoteRequest): Promise<{note: Note; cards: Card[]}>
  fetchNote(id: number): Promise<Note>
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
  fetchDecks,
  createDeck,
  fetchDeck,
  fetchDeckStats,
  fetchDeckNotes,
  importNotesFile,
  fetchNoteTypes,
  fetchNoteType,
  addField,
  renameField,
  removeField,
  reorderFields,
  setSortField,
  setFieldOptions,
  updateTemplate,
  createNote,
  fetchNote,
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
    fetchDecks: () => notImplemented('fetchDecks'),
    createDeck: () => notImplemented('createDeck'),
    fetchDeck: () => notImplemented('fetchDeck'),
    fetchDeckStats: () => notImplemented('fetchDeckStats'),
    fetchDeckNotes: () => notImplemented('fetchDeckNotes'),
    importNotesFile: () => notImplemented('importNotesFile'),
    fetchNoteTypes: () => notImplemented('fetchNoteTypes'),
    fetchNoteType: () => notImplemented('fetchNoteType'),
    addField: () => notImplemented('addField'),
    renameField: () => notImplemented('renameField'),
    removeField: () => notImplemented('removeField'),
    reorderFields: () => notImplemented('reorderFields'),
    setSortField: () => notImplemented('setSortField'),
    setFieldOptions: () => notImplemented('setFieldOptions'),
    updateTemplate: () => notImplemented('updateTemplate'),
    createNote: () => notImplemented('createNote'),
    fetchNote: () => notImplemented('fetchNote'),
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
            maxSharedDecks: 0,
            maxSyncDevices: 0,
            maxWorkspaces: 1,
          },
          usage: {
            decks: 0,
            notes: 0,
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
          maxSharedDecks: 0,
          maxSyncDevices: 0,
          maxWorkspaces: 1,
        },
        usage: {
          decks: 0,
          notes: 0,
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
