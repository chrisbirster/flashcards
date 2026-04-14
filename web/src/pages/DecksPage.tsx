import { useRef, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useNavigate } from 'react-router'
import { APIError, type ImportNotesResponse, type ImportSource } from '#/lib/api'
import { useAppRepository } from '#/lib/app-repository'
import { DeckItem } from '#/components/deck-item'
import { EmptyState, PageContainer, PageSection } from '#/components/page-layout'

function SectionToggle({
  title,
  description,
  expanded,
  onToggle,
}: {
  title: string
  description: string
  expanded: boolean
  onToggle: () => void
}) {
  return (
    <button
      type="button"
      onClick={onToggle}
      className="flex w-full items-center justify-between gap-4 rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-card)] px-4 py-4 text-left md:hidden"
    >
      <div>
        <p className="text-sm font-semibold text-[var(--app-text)]">{title}</p>
        <p className="mt-1 text-sm text-[var(--app-text-soft)]">{description}</p>
      </div>
      <span className="text-xs font-semibold uppercase tracking-[0.18em] text-[var(--app-accent)]">
        {expanded ? 'Hide' : 'Open'}
      </span>
    </button>
  )
}

export function DecksPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const [newDeckName, setNewDeckName] = useState('')
  const [importFile, setImportFile] = useState<File | null>(null)
  const [importSource, setImportSource] = useState<ImportSource>('auto')
  const [importDeckName, setImportDeckName] = useState('')
  const [importNoteType, setImportNoteType] = useState('Basic')
  const [importResult, setImportResult] = useState<ImportNotesResponse | null>(null)
  const [showCreateSection, setShowCreateSection] = useState(false)
  const [showImportSection, setShowImportSection] = useState(false)
  const importInputRef = useRef<HTMLInputElement>(null)

  const { data: decks, isLoading, error } = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })

  const createDeckMutation = useMutation({
    mutationFn: (req: { name: string }) => repository.createDeck(req),
    onSuccess: (deck) => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
      setNewDeckName('')
      setShowCreateSection(false)
      navigate(`/notes/add?deckId=${deck.id}`)
    },
  })

  const importMutation = useMutation({
    mutationFn: (req: Parameters<typeof repository.importNotesFile>[0]) => repository.importNotesFile(req),
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
      setImportResult(result)
      setImportFile(null)
      setShowImportSection(false)
      if (importInputRef.current) {
        importInputRef.current.value = ''
      }
    },
  })

  const handleCreateDeck = (event: React.FormEvent) => {
    event.preventDefault()
    if (newDeckName.trim()) {
      createDeckMutation.mutate({ name: newDeckName })
    }
  }

  const handleImportDeck = (event: React.FormEvent) => {
    event.preventDefault()
    if (!importFile) return

    setImportResult(null)
    importMutation.mutate({
      file: importFile,
      source: importSource,
      deckName: importDeckName.trim() || undefined,
      noteType: importNoteType.trim() || undefined,
    })
  }

  if (isLoading) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="p-5 text-sm text-[var(--app-text-soft)]">Loading decks...</PageSection>
      </PageContainer>
    )
  }

  if (error) {
    return (
      <PageContainer className="space-y-4">
        <PageSection className="border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] p-5 text-sm text-[var(--app-danger-text)]">
          Error: {error instanceof Error ? error.message : 'Failed to load decks'}
        </PageSection>
      </PageContainer>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <div className="grid gap-4 md:grid-cols-2">
        <div className="rounded-[1.75rem] border border-[var(--app-line)] bg-[var(--app-card-strong)] p-5 shadow-sm sm:p-6">
          <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Deck workspace</p>
          <h2 className="mt-3 text-2xl font-semibold tracking-tight text-[var(--app-text)]">Organize the decks you study from every day.</h2>
          <p className="mt-3 text-sm leading-6 text-[var(--app-text-soft)]">
            Create new decks, import source material, and jump straight into study or note creation from a phone-friendly deck list.
          </p>
        </div>

        <PageSection className="p-5 sm:p-6">
          <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Current inventory</p>
          <div className="mt-4 grid grid-cols-3 gap-3">
            <div className="rounded-2xl bg-[var(--app-muted-surface)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Decks</p>
              <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">{decks?.length ?? 0}</p>
            </div>
            <div className="rounded-2xl bg-[var(--app-muted-surface)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Notes</p>
              <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                {(decks ?? []).reduce((sum, deck) => sum + deck.noteCount, 0)}
              </p>
            </div>
            <div className="rounded-2xl bg-[var(--app-muted-surface)] p-4">
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Due</p>
              <p className="mt-2 text-2xl font-semibold text-[var(--app-text)]">
                {(decks ?? []).reduce((sum, deck) => sum + deck.dueToday, 0)}
              </p>
            </div>
          </div>
        </PageSection>
      </div>

      <SectionToggle
        title="Create a deck"
        description="Open a compact mobile form for new deck creation."
        expanded={showCreateSection}
        onToggle={() => setShowCreateSection((current) => !current)}
      />
      <PageSection className={`${showCreateSection ? 'block' : 'hidden'} p-5 sm:p-6 md:block`}>
        <h3 className="text-lg font-semibold text-[var(--app-text)]">Create new deck</h3>
        <p className="mt-2 text-sm text-[var(--app-text-soft)]">Deck names are lightweight on purpose. You can rename later once the structure settles.</p>
        <form onSubmit={handleCreateDeck} className="mt-4 flex flex-col gap-3 sm:flex-row">
          <input
            type="text"
            value={newDeckName}
            onChange={(event) => setNewDeckName(event.target.value)}
            placeholder="Deck name"
            className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
          />
          <button
            type="submit"
            disabled={createDeckMutation.isPending || !newDeckName.trim()}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60"
          >
            {createDeckMutation.isPending ? 'Creating...' : 'Create'}
          </button>
        </form>
        {createDeckMutation.isError ? (
          <p className="mt-3 text-sm text-[var(--app-danger-text)]">
            Error: {createDeckMutation.error instanceof Error ? createDeckMutation.error.message : 'Failed to create deck'}
          </p>
        ) : null}
      </PageSection>

      {createDeckMutation.isError && createDeckMutation.error instanceof APIError && createDeckMutation.error.code === 'plan_limit_exceeded' ? (
        <PageSection className="border-[var(--app-warning-line)] bg-[var(--app-warning-surface)] p-5 text-sm text-[var(--app-warning-text)]">
          <p className="font-medium">Deck limit reached.</p>
          <p className="mt-1">{createDeckMutation.error.message} Sign in or upgrade when billing is configured to unlock more decks.</p>
        </PageSection>
      ) : null}

      <SectionToggle
        title="Import notes"
        description="Bring in notes from JSON, YAML, text, or Anki exports."
        expanded={showImportSection}
        onToggle={() => setShowImportSection((current) => !current)}
      />
      <PageSection className={`${showImportSection ? 'block' : 'hidden'} p-5 sm:p-6 md:block`}>
        <h3 className="text-lg font-semibold text-[var(--app-text)]">Import notes</h3>
        <p className="mt-2 text-sm text-[var(--app-text-soft)]">
          Recommended format is JSON or YAML. Text, TSV, and supported Anki exports still work when you need a quick migration.
        </p>
        <form onSubmit={handleImportDeck} className="mt-4 grid gap-3 md:grid-cols-2">
          <input
            ref={importInputRef}
            type="file"
            accept=".json,.yaml,.yml,.csv,.tsv,.txt,.apkg,.colpkg"
            onChange={(event) => setImportFile(event.target.files?.[0] ?? null)}
            className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)]"
          />
          <select
            value={importSource}
            onChange={(event) => setImportSource(event.target.value as ImportSource)}
            className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)]"
          >
            <option value="auto">Auto detect</option>
            <option value="native">Native JSON/YAML</option>
            <option value="anki">Anki</option>
            <option value="quizlet">Quizlet</option>
          </select>
          <input
            type="text"
            value={importDeckName}
            onChange={(event) => setImportDeckName(event.target.value)}
            placeholder="Optional deck override"
            className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
          />
          <select
            value={importNoteType}
            onChange={(event) => setImportNoteType(event.target.value)}
            className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)]"
          >
            <option value="Basic">Basic</option>
            <option value="Cloze">Cloze</option>
          </select>
          <button
            type="submit"
            disabled={importMutation.isPending || !importFile}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-5 text-sm font-semibold text-[var(--app-accent-ink)] disabled:opacity-60 md:col-span-2"
          >
            {importMutation.isPending ? 'Importing...' : 'Import file'}
          </button>
        </form>

        {importMutation.isError ? (
          <p className="mt-3 text-sm text-[var(--app-danger-text)]">
            Error: {importMutation.error instanceof Error ? importMutation.error.message : 'Failed to import file'}
          </p>
        ) : null}

        {importResult ? (
          <div className="mt-4 rounded-2xl border border-[var(--app-success-line)] bg-[var(--app-success-surface)] p-4 text-sm text-[var(--app-success-text)]">
            <p>
              Imported {importResult.imported} note(s), skipped {importResult.skipped}. Detected source: {importResult.source}, format: {importResult.format}.
            </p>
            {importResult.decksCreated?.length ? <p className="mt-1">Created decks: {importResult.decksCreated.join(', ')}</p> : null}
            {importResult.errors?.length ? <p className="mt-1">{importResult.errors.slice(0, 5).join(' | ')}</p> : null}
          </div>
        ) : null}
      </PageSection>

      <PageSection className="p-4 sm:p-5">
        <div className="flex items-center justify-between gap-3 border-b border-[var(--app-line)] px-1 pb-4">
          <div>
            <h3 className="text-lg font-semibold text-[var(--app-text)]">Your decks</h3>
            <p className="text-sm text-[var(--app-text-soft)]">Study, rename, or add notes without leaving the list.</p>
          </div>
        </div>

        {decks?.length ? (
          <ul className="mt-4 space-y-4">
            {decks.map((deck) => (
              <DeckItem key={deck.id} deck={deck} />
            ))}
          </ul>
        ) : (
          <div className="mt-4">
            <EmptyState
              title="No decks yet"
              description="Create a deck or import source notes above and your workspace will start filling out here."
            />
          </div>
        )}
      </PageSection>
    </PageContainer>
  )
}
