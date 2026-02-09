import { useRef, useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchDecks, createDeck, importNotesFile, type ImportNotesResponse, type ImportSource } from '#/lib/api'
import { DeckItem } from '#/components/deck-item'

export function DecksPage() {
  const [newDeckName, setNewDeckName] = useState('')
  const [importFile, setImportFile] = useState<File | null>(null)
  const [importSource, setImportSource] = useState<ImportSource>('auto')
  const [importDeckName, setImportDeckName] = useState('')
  const [importNoteType, setImportNoteType] = useState('Basic')
  const [importResult, setImportResult] = useState<ImportNotesResponse | null>(null)
  const importInputRef = useRef<HTMLInputElement>(null)
  const queryClient = useQueryClient()

  const { data: decks, isLoading, error } = useQuery({
    queryKey: ['decks'],
    queryFn: fetchDecks,
  })

  const createDeckMutation = useMutation({
    mutationFn: createDeck,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      setNewDeckName('')
    },
  })

  const importMutation = useMutation({
    mutationFn: importNotesFile,
    onSuccess: (result) => {
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      setImportResult(result)
      setImportFile(null)
      if (importInputRef.current) {
        importInputRef.current.value = ''
      }
    },
  })

  const handleCreateDeck = (e: React.FormEvent) => {
    e.preventDefault()
    if (newDeckName.trim()) {
      createDeckMutation.mutate({ name: newDeckName })
    }
  }

  const handleImportDeck = (e: React.FormEvent) => {
    e.preventDefault()
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
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-red-600">
          Error: {error instanceof Error ? error.message : 'Failed to load decks'}
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">Create New Deck</h2>
        <form onSubmit={handleCreateDeck} className="flex gap-2">
          <input
            type="text"
            value={newDeckName}
            onChange={(e) => setNewDeckName(e.target.value)}
            placeholder="Deck name"
            className="flex-1 px-3 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <button
            type="submit"
            disabled={createDeckMutation.isPending || !newDeckName.trim()}
            className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300"
          >
            {createDeckMutation.isPending ? 'Creating...' : 'Create'}
          </button>
        </form>
        {createDeckMutation.isError && (
          <p className="mt-2 text-red-600 text-sm">
            Error: {createDeckMutation.error instanceof Error ? createDeckMutation.error.message : 'Failed to create deck'}
          </p>
        )}
      </div>

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-xl font-semibold mb-2">Import Notes</h2>
        <p className="text-sm text-gray-600 mb-4">
          Recommended format is JSON/YAML (safe for tabs/newlines). Also supports Anki text/APKG and Quizlet exports.
        </p>
        <form onSubmit={handleImportDeck} className="grid gap-3 md:grid-cols-2">
          <input
            ref={importInputRef}
            type="file"
            accept=".json,.yaml,.yml,.csv,.tsv,.txt,.apkg,.colpkg"
            onChange={(e) => setImportFile(e.target.files?.[0] ?? null)}
            className="w-full px-3 py-2 border rounded-md"
          />
          <select
            value={importSource}
            onChange={(e) => setImportSource(e.target.value as ImportSource)}
            className="w-full px-3 py-2 border rounded-md bg-white"
          >
            <option value="auto">Auto Detect</option>
            <option value="native">Native JSON/YAML</option>
            <option value="anki">Anki</option>
            <option value="quizlet">Quizlet</option>
          </select>
          <input
            type="text"
            value={importDeckName}
            onChange={(e) => setImportDeckName(e.target.value)}
            placeholder="Optional deck override"
            className="w-full px-3 py-2 border rounded-md"
          />
          <select
            value={importNoteType}
            onChange={(e) => setImportNoteType(e.target.value)}
            className="w-full px-3 py-2 border rounded-md bg-white"
          >
            <option value="Basic">Basic</option>
            <option value="Cloze">Cloze</option>
          </select>
          <button
            type="submit"
            disabled={importMutation.isPending || !importFile}
            className="md:col-span-2 px-6 py-2 bg-emerald-600 text-white rounded-md hover:bg-emerald-700 disabled:bg-gray-300"
          >
            {importMutation.isPending ? 'Importing...' : 'Import File'}
          </button>
        </form>

        {importMutation.isError && (
          <p className="mt-3 text-red-600 text-sm">
            Error: {importMutation.error instanceof Error ? importMutation.error.message : 'Failed to import file'}
          </p>
        )}

        {importResult && (
          <div className="mt-3 p-3 rounded-md bg-emerald-50 text-sm text-emerald-900">
            <p>
              Imported {importResult.imported} note(s), skipped {importResult.skipped}. Detected source: {importResult.source}, format: {importResult.format}.
            </p>
            {importResult.decksCreated && importResult.decksCreated.length > 0 && (
              <p className="mt-1">Created decks: {importResult.decksCreated.join(', ')}</p>
            )}
            {importResult.errors && importResult.errors.length > 0 && (
              <p className="mt-1 text-amber-700">Warnings: {importResult.errors.slice(0, 5).join(' | ')}</p>
            )}
          </div>
        )}
      </div>

      <div className="bg-white rounded-lg shadow">
        <h2 className="text-xl font-semibold p-6 pb-4">Your Decks</h2>
        {decks && decks.length > 0 ? (
          <ul className="divide-y divide-gray-200">
            {decks.map((deck) => (
              <DeckItem key={deck.id} deck={deck} />
            ))}
          </ul>
        ) : (
          <p className="p-6 text-gray-500 text-center">
            No decks yet. Create your first deck above!
          </p>
        )}
      </div>
    </div>
  )
}
