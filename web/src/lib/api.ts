const API_BASE = (import.meta.env.VITE_API_BASE ?? '/api').replace(/\/$/, '')

export class APIError extends Error {
  status: number
  code?: string

  constructor(message: string, status: number, code?: string) {
    super(message)
    this.name = 'APIError'
    this.status = status
    this.code = code
  }
}

async function requestJSON<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    credentials: 'include',
    ...init,
  })

  if (!res.ok) {
    if (res.status === 401 && typeof window !== 'undefined' && !window.location.pathname.startsWith('/login')) {
      window.location.assign('/login')
    }
    const contentType = res.headers.get('content-type') || ''
    if (contentType.includes('application/json')) {
      const payload = await res.json().catch(() => null) as {message?: string; code?: string} | null
      throw new APIError(payload?.message || 'Request failed', res.status, payload?.code)
    }
    const text = await res.text()
    throw new APIError(text || 'Request failed', res.status)
  }

  if (res.status === 204) {
    return undefined as T
  }

  return res.json()
}

export interface Deck {
  id: number
  name: string
  parentId?: number
  cardIds: number[]
}

export interface DeckStats {
  deckId: number
  newCards: number
  learning: number
  review: number
  relearning: number
  suspended: number
  buried: number
  totalCards: number
  dueToday: number
}

export interface Note {
  id: number
  typeId: string
  fieldVals: Record<string, string>
  tags: string[]
  createdAt: string
  modifiedAt: string
  deckId?: number
  cardCount?: number
}

export interface Card {
  id: number
  noteId: number
  deckId: number
  templateName: string
  ordinal: number
  front: string
  back: string
  flag: number
  marked: boolean
  suspended: boolean
}

export interface CardTemplate {
  name: string
  qFmt: string
  aFmt: string
  styling: string
  ifFieldNonEmpty?: string
  isCloze: boolean
  deckOverride?: string
  browserQFmt?: string
  browserAFmt?: string
}

export interface FieldOptions {
  font?: string
  fontSize?: number
  rtl?: boolean
  htmlEditor?: boolean
}

export interface NoteType {
  name: string
  fields: string[]
  templates: CardTemplate[]
  sortFieldIndex: number
  fieldOptions?: Record<string, FieldOptions>
}

export interface CreateDeckRequest {
  name: string
}

export type ImportSource = 'auto' | 'native' | 'anki' | 'quizlet'

export interface ImportFileRequest {
  file: File
  source?: ImportSource
  deckName?: string
  noteType?: string
  format?: string
}

export interface ImportNotesResponse {
  imported: number
  skipped: number
  source: string
  format: string
  decksCreated?: string[]
  errors?: string[]
}

export interface CreateNoteRequest {
  typeId: string
  deckId: number
  fieldVals: Record<string, string>
  tags?: string[]
}

export interface UpdateNoteRequest {
  typeId: string
  deckId: number
  fieldVals: Record<string, string>
  tags?: string[]
}

export interface UpdateDeckRequest {
  name: string
}

export interface CreateTemplateRequest {
  name: string
  sourceTemplateName?: string
}

export interface AnswerCardRequest {
  rating: number // 1=Again, 2=Hard, 3=Good, 4=Easy
  timeTakenMs?: number // Time spent on the card in milliseconds
}

export interface CheckDuplicateRequest {
  typeId: string
  fieldName: string
  value: string
  deckId?: number
}

export interface NoteBrief {
  id: number
  typeId: string
  fieldVals: Record<string, string>
  deckId?: number
}

export interface DuplicateResult {
  isDuplicate: boolean
  duplicates?: NoteBrief[]
}

export interface RecentDeckNoteSummary {
  noteId: number
  noteType: string
  createdAt: string
  modifiedAt: string
  tags: string[]
  fieldPreview: string
  cardCountInDeck: number
}

export interface DeckNotesResponse {
  notes: RecentDeckNoteSummary[]
}

export interface NoteListItem {
  id: number
  typeId: string
  fieldVals: Record<string, string>
  fieldPreview: string
  tags: string[]
  createdAt: string
  modifiedAt: string
  deckId?: number
  deckName?: string
  cardCount: number
}

export interface ListNotesResponse {
  notes: NoteListItem[]
  total: number
  nextCursor?: string
  prevCursor?: string
}

export interface ListNotesParams {
  deckId?: number
  q?: string
  typeId?: string
  tag?: string
  limit?: number
  cursor?: string
}

export interface PlanLimits {
  maxDecks: number
  maxNotes: number
  maxCardsTotal: number
  maxSharedDecks: number
  maxSyncDevices: number
  maxWorkspaces: number
}

export interface EntitlementUsage {
  decks: number
  notes: number
  cardsTotal: number
  sharedDecks: number
  syncDevices: number
  workspaces: number
}

export interface EntitlementFeatures {
  googleLogin: boolean
  accountBacked: boolean
  sync: boolean
  shareDecks: boolean
  organizations: boolean
  studyGroups: boolean
  marketplacePublish: boolean
  enterprise: boolean
}

export interface Entitlements {
  plan: 'guest' | 'free' | 'pro' | 'team' | 'enterprise'
  limits: PlanLimits
  usage: EntitlementUsage
  features: EntitlementFeatures
}

export interface AccountUser {
  id: string
  email: string
  displayName: string
  avatarUrl?: string
}

export interface WorkspaceSession {
  id: string
  name: string
  slug: string
  collectionId: string
}

export interface AuthSessionResponse {
  authenticated: boolean
  googleAuthConfigured: boolean
  otpAuthEnabled: boolean
  user?: AccountUser
  workspace?: WorkspaceSession
  entitlements: Entitlements
}

export interface OTPRequestResponse {
  ok: boolean
  expiresAt: string
  retryAfterSeconds: number
  delivery?: 'email' | 'dev-inline'
  devCode?: string
}

// Deck endpoints
export async function fetchDecks(): Promise<Deck[]> {
  return requestJSON(`${API_BASE}/decks`)
}

export async function createDeck(req: CreateDeckRequest): Promise<Deck> {
  return requestJSON(`${API_BASE}/decks`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function fetchDeck(id: number): Promise<{deck: Deck; stats: DeckStats}> {
  return requestJSON(`${API_BASE}/decks/${id}`)
}

export async function fetchDeckStats(id: number): Promise<DeckStats> {
  return requestJSON(`${API_BASE}/decks/${id}/stats`)
}

export async function fetchDeckNotes(deckId: number, limit: number = 20, cursor?: string): Promise<DeckNotesResponse> {
  const params = new URLSearchParams({limit: String(limit)})
  if (cursor) params.set('cursor', cursor)
  return requestJSON(`${API_BASE}/decks/${deckId}/notes?${params.toString()}`)
}

export async function updateDeck(id: number, req: UpdateDeckRequest): Promise<Deck> {
  return requestJSON(`${API_BASE}/decks/${id}`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function deleteDeck(id: number): Promise<void> {
  return requestJSON(`${API_BASE}/decks/${id}`, {
    method: 'DELETE',
  })
}

export async function importNotesFile(req: ImportFileRequest): Promise<ImportNotesResponse> {
  const formData = new FormData()
  formData.append('file', req.file)
  if (req.source) formData.append('source', req.source)
  if (req.deckName) formData.append('deckName', req.deckName)
  if (req.noteType) formData.append('noteType', req.noteType)
  if (req.format) formData.append('format', req.format)

  const res = await fetch(`${API_BASE}/import`, {
    method: 'POST',
    credentials: 'include',
    body: formData,
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to import file')
  }
  return res.json()
}

// Note Type endpoints
export async function fetchNoteTypes(): Promise<NoteType[]> {
  return requestJSON(`${API_BASE}/note-types`)
}

export async function fetchNoteType(name: string): Promise<NoteType> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(name)}`)
}

// Field management
export interface AddFieldRequest {
  fieldName: string
  position?: number
}

export interface RenameFieldRequest {
  oldName: string
  newName: string
}

export interface RemoveFieldRequest {
  fieldName: string
}

export interface ReorderFieldsRequest {
  fields: string[]
}

export interface FieldsResponse {
  message: string
  fields: string[]
}

export interface UpdateTemplateRequest {
  name?: string
  qFmt?: string
  aFmt?: string
  styling?: string
  ifFieldNonEmpty?: string
  deckOverride?: string
  browserQFmt?: string
  browserAFmt?: string
}

export interface TemplatesResponse {
  message: string
  templates: CardTemplate[]
}

export async function createTemplate(noteTypeName: string, req: CreateTemplateRequest): Promise<TemplatesResponse> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function addField(noteTypeName: string, req: AddFieldRequest): Promise<FieldsResponse> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function renameField(noteTypeName: string, req: RenameFieldRequest): Promise<FieldsResponse> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/rename`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function removeField(noteTypeName: string, req: RemoveFieldRequest): Promise<FieldsResponse> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`, {
    method: 'DELETE',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function reorderFields(noteTypeName: string, req: ReorderFieldsRequest): Promise<FieldsResponse> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/reorder`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export interface SetSortFieldRequest {
  fieldIndex: number
}

export async function setSortField(noteTypeName: string, req: SetSortFieldRequest): Promise<{message: string; sortFieldIndex: number; sortFieldName: string}> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/sort-field`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function setFieldOptions(
  noteTypeName: string,
  fieldName: string,
  options: FieldOptions
): Promise<{message: string; fieldOptions: Record<string, FieldOptions>}> {
  return requestJSON(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/options`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({fieldName, options}),
  })
}

export async function updateTemplate(
  noteTypeName: string,
  templateName: string,
  req: UpdateTemplateRequest
): Promise<TemplatesResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates/${encodeURIComponent(templateName)}`,
    {
      method: 'PATCH',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(req),
    }
  )
}

export async function deleteTemplate(noteTypeName: string, templateName: string): Promise<TemplatesResponse> {
  return requestJSON(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates/${encodeURIComponent(templateName)}`,
    {
      method: 'DELETE',
    }
  )
}

// Note endpoints
export async function createNote(req: CreateNoteRequest): Promise<{note: Note; cards: Card[]}> {
  return requestJSON(`${API_BASE}/notes`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function fetchNote(id: number): Promise<Note> {
  return requestJSON(`${API_BASE}/notes/${id}`)
}

export async function fetchNotes(params: ListNotesParams = {}): Promise<ListNotesResponse> {
  const query = new URLSearchParams()
  if (params.deckId) query.set('deckId', String(params.deckId))
  if (params.q) query.set('q', params.q)
  if (params.typeId) query.set('typeId', params.typeId)
  if (params.tag) query.set('tag', params.tag)
  if (params.limit) query.set('limit', String(params.limit))
  if (params.cursor) query.set('cursor', params.cursor)
  const suffix = query.toString() ? `?${query.toString()}` : ''
  return requestJSON(`${API_BASE}/notes${suffix}`)
}

export async function updateNote(id: number, req: UpdateNoteRequest): Promise<{note: Note; cards: Card[]}> {
  return requestJSON(`${API_BASE}/notes/${id}`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export async function deleteNote(id: number): Promise<void> {
  return requestJSON(`${API_BASE}/notes/${id}`, {
    method: 'DELETE',
  })
}

export async function checkDuplicate(req: CheckDuplicateRequest): Promise<DuplicateResult> {
  return requestJSON(`${API_BASE}/notes/check-duplicate`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

// Card endpoints
export async function fetchDueCards(deckId: number, limit: number = 10): Promise<Card[]> {
  return requestJSON(`${API_BASE}/decks/${deckId}/due?limit=${limit}`)
}

export async function fetchCard(id: number): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}`)
}

export async function answerCard(id: number, req: AnswerCardRequest): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}/answer`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

export interface UpdateCardRequest {
  flag?: number    // 0-7 color flags
  marked?: boolean // toggle marked status
  suspended?: boolean // toggle suspended status
}

export async function updateCard(id: number, req: UpdateCardRequest): Promise<Card> {
  return requestJSON(`${API_BASE}/cards/${id}`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

// Empty cards endpoints
export interface EmptyCardInfo {
  cardId: number
  noteId: number
  deckId: number
  templateName: string
  ordinal: number
  front: string
  back: string
  reason: string
}

export interface EmptyCardsResponse {
  count: number
  emptyCards: EmptyCardInfo[]
}

export interface DeleteEmptyCardsRequest {
  cardIds: number[]
}

export interface DeleteEmptyCardsResponse {
  deleted: number
  failed?: string[]
}

export async function findEmptyCards(): Promise<EmptyCardsResponse> {
  return requestJSON(`${API_BASE}/cards/empty`)
}

export async function deleteEmptyCards(req: DeleteEmptyCardsRequest): Promise<DeleteEmptyCardsResponse> {
  return requestJSON(`${API_BASE}/cards/empty/delete`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
}

// Backup endpoints
export async function createBackup(): Promise<{message: string; backupPath: string}> {
  return requestJSON(`${API_BASE}/backups`, {
    method: 'POST',
  })
}

export async function listBackups(): Promise<Array<{path: string; filename: string; size: number; modified: string}>> {
  return requestJSON(`${API_BASE}/backups`)
}

// Health check
export async function healthCheck(): Promise<{status: string; service: string; version: string}> {
  return requestJSON(`${API_BASE}/health`)
}

export async function fetchSession(): Promise<AuthSessionResponse> {
  return requestJSON(`${API_BASE}/auth/session`)
}

export async function requestOTP(email: string): Promise<OTPRequestResponse> {
  return requestJSON(`${API_BASE}/auth/otp/request`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email}),
  })
}

export async function verifyOTP(email: string, code: string): Promise<AuthSessionResponse> {
  return requestJSON(`${API_BASE}/auth/otp/verify`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({email, code}),
  })
}

export async function fetchEntitlements(): Promise<Entitlements> {
  return requestJSON(`${API_BASE}/entitlements`)
}

export async function logout(): Promise<{ok: boolean}> {
  return requestJSON(`${API_BASE}/auth/logout`, {
    method: 'POST',
  })
}
