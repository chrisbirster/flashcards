const API_BASE = 'http://localhost:8080/api'

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

export interface CreateNoteRequest {
  typeId: string
  deckId: number
  fieldVals: Record<string, string>
  tags?: string[]
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

// Deck endpoints
export async function fetchDecks(): Promise<Deck[]> {
  const res = await fetch(`${API_BASE}/decks`)
  if (!res.ok) throw new Error('Failed to fetch decks')
  return res.json()
}

export async function createDeck(req: CreateDeckRequest): Promise<Deck> {
  const res = await fetch(`${API_BASE}/decks`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to create deck')
  return res.json()
}

export async function fetchDeck(id: number): Promise<{deck: Deck; stats: DeckStats}> {
  const res = await fetch(`${API_BASE}/decks/${id}`)
  if (!res.ok) throw new Error('Failed to fetch deck')
  return res.json()
}

export async function fetchDeckStats(id: number): Promise<DeckStats> {
  const res = await fetch(`${API_BASE}/decks/${id}/stats`)
  if (!res.ok) throw new Error('Failed to fetch deck stats')
  return res.json()
}

// Note Type endpoints
export async function fetchNoteTypes(): Promise<NoteType[]> {
  const res = await fetch(`${API_BASE}/note-types`)
  if (!res.ok) throw new Error('Failed to fetch note types')
  return res.json()
}

export async function fetchNoteType(name: string): Promise<NoteType> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(name)}`)
  if (!res.ok) throw new Error('Failed to fetch note type')
  return res.json()
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

export async function addField(noteTypeName: string, req: AddFieldRequest): Promise<FieldsResponse> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to add field')
  }
  return res.json()
}

export async function renameField(noteTypeName: string, req: RenameFieldRequest): Promise<FieldsResponse> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/rename`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to rename field')
  }
  return res.json()
}

export async function removeField(noteTypeName: string, req: RemoveFieldRequest): Promise<FieldsResponse> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields`, {
    method: 'DELETE',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to remove field')
  }
  return res.json()
}

export async function reorderFields(noteTypeName: string, req: ReorderFieldsRequest): Promise<FieldsResponse> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/reorder`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to reorder fields')
  }
  return res.json()
}

export interface SetSortFieldRequest {
  fieldIndex: number
}

export async function setSortField(noteTypeName: string, req: SetSortFieldRequest): Promise<{message: string; sortFieldIndex: number; sortFieldName: string}> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/sort-field`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to set sort field')
  }
  return res.json()
}

export async function setFieldOptions(
  noteTypeName: string,
  fieldName: string,
  options: FieldOptions
): Promise<{message: string; fieldOptions: Record<string, FieldOptions>}> {
  const res = await fetch(`${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/fields/options`, {
    method: 'PUT',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({fieldName, options}),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to set field options')
  }
  return res.json()
}

export async function updateTemplate(
  noteTypeName: string,
  templateName: string,
  req: UpdateTemplateRequest
): Promise<TemplatesResponse> {
  const res = await fetch(
    `${API_BASE}/note-types/${encodeURIComponent(noteTypeName)}/templates/${encodeURIComponent(templateName)}`,
    {
      method: 'PATCH',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(req),
    }
  )
  if (!res.ok) {
    const text = await res.text()
    throw new Error(text || 'Failed to update template')
  }
  return res.json()
}

// Note endpoints
export async function createNote(req: CreateNoteRequest): Promise<{note: Note; cards: Card[]}> {
  const res = await fetch(`${API_BASE}/notes`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to create note')
  return res.json()
}

export async function fetchNote(id: number): Promise<Note> {
  const res = await fetch(`${API_BASE}/notes/${id}`)
  if (!res.ok) throw new Error('Failed to fetch note')
  return res.json()
}

export async function checkDuplicate(req: CheckDuplicateRequest): Promise<DuplicateResult> {
  const res = await fetch(`${API_BASE}/notes/check-duplicate`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to check duplicate')
  return res.json()
}

// Card endpoints
export async function fetchDueCards(deckId: number, limit: number = 10): Promise<Card[]> {
  const res = await fetch(`${API_BASE}/decks/${deckId}/due?limit=${limit}`)
  if (!res.ok) throw new Error('Failed to fetch due cards')
  return res.json()
}

export async function fetchCard(id: number): Promise<Card> {
  const res = await fetch(`${API_BASE}/cards/${id}`)
  if (!res.ok) throw new Error('Failed to fetch card')
  return res.json()
}

export async function answerCard(id: number, req: AnswerCardRequest): Promise<Card> {
  const res = await fetch(`${API_BASE}/cards/${id}/answer`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to answer card')
  return res.json()
}

export interface UpdateCardRequest {
  flag?: number    // 0-7 color flags
  marked?: boolean // toggle marked status
  suspended?: boolean // toggle suspended status
}

export async function updateCard(id: number, req: UpdateCardRequest): Promise<Card> {
  const res = await fetch(`${API_BASE}/cards/${id}`, {
    method: 'PATCH',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to update card')
  return res.json()
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
  const res = await fetch(`${API_BASE}/cards/empty`)
  if (!res.ok) throw new Error('Failed to find empty cards')
  return res.json()
}

export async function deleteEmptyCards(req: DeleteEmptyCardsRequest): Promise<DeleteEmptyCardsResponse> {
  const res = await fetch(`${API_BASE}/cards/empty/delete`, {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify(req),
  })
  if (!res.ok) throw new Error('Failed to delete empty cards')
  return res.json()
}

// Backup endpoints
export async function createBackup(): Promise<{message: string; backupPath: string}> {
  const res = await fetch(`${API_BASE}/backups`, {
    method: 'POST',
  })
  if (!res.ok) throw new Error('Failed to create backup')
  return res.json()
}

export async function listBackups(): Promise<Array<{path: string; filename: string; size: number; modified: string}>> {
  const res = await fetch(`${API_BASE}/backups`)
  if (!res.ok) throw new Error('Failed to list backups')
  return res.json()
}

// Health check
export async function healthCheck(): Promise<{status: string; service: string; version: string}> {
  const res = await fetch(`${API_BASE}/health`)
  if (!res.ok) throw new Error('Health check failed')
  return res.json()
}
