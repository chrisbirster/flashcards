import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchDecks, fetchDeckStats, createDeck, type Deck } from '#/lib/api'
import { StudyScreen } from '#/components/StudyScreen'
import { AddNoteScreen } from '#/components/AddNoteScreen'

function DeckItem({ deck, onStudy, onAddCards }: { deck: Deck; onStudy: (deck: Deck) => void; onAddCards: (deck: Deck) => void }) {
  const { data: stats } = useQuery({
    queryKey: ['deck-stats', deck.id],
    queryFn: () => fetchDeckStats(deck.id),
  })

  return (
    <li className="p-6 hover:bg-gray-50 transition-colors">
      <div className="flex items-center justify-between">
        <div className="flex-1">
          <h3 className="text-lg font-medium text-gray-900">{deck.name}</h3>
          <div className="flex gap-4 mt-2 text-sm">
            {stats ? (
              <>
                <span className="text-blue-600 font-medium">
                  {stats.newCards} new
                </span>
                <span className="text-orange-600 font-medium">
                  {stats.learning} learning
                </span>
                <span className="text-green-600 font-medium">
                  {stats.review} review
                </span>
                {stats.suspended > 0 && (
                  <span className="text-gray-500">
                    {stats.suspended} suspended
                  </span>
                )}
                <span className="text-gray-400">
                  ({stats.totalCards} total)
                </span>
              </>
            ) : (
              <span className="text-gray-500">Loading stats...</span>
            )}
          </div>
          {stats && stats.dueToday > 0 && (
            <div className="mt-1 text-sm font-semibold text-indigo-600">
              {stats.dueToday} due today
            </div>
          )}
        </div>
        <div className="flex gap-2">
          <button
            className="px-4 py-2 text-sm text-white bg-blue-600 hover:bg-blue-700 rounded-md disabled:bg-gray-300"
            disabled={!stats || stats.dueToday === 0}
            onClick={() => onStudy(deck)}
          >
            Study Now
          </button>
          <button
            className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-md"
            onClick={() => onAddCards(deck)}
          >
            Add Cards
          </button>
        </div>
      </div>
    </li>
  )
}

export default function App() {
  const [newDeckName, setNewDeckName] = useState('')
  const [studyingDeck, setStudyingDeck] = useState<Deck | null>(null)
  const [addingToDeck, setAddingToDeck] = useState<Deck | null>(null)
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

  const handleCreateDeck = (e: React.FormEvent) => {
    e.preventDefault()
    if (newDeckName.trim()) {
      createDeckMutation.mutate({ name: newDeckName })
    }
  }

  // If studying a deck, show study screen
  if (studyingDeck) {
    return (
      <StudyScreen
        deckId={studyingDeck.id}
        deckName={studyingDeck.name}
        onExit={() => setStudyingDeck(null)}
      />
    )
  }

  // If adding cards to a deck, show add note screen
  if (addingToDeck) {
    return (
      <AddNoteScreen
        deckId={addingToDeck.id}
        onClose={() => setAddingToDeck(null)}
        onSuccess={() => {
          queryClient.invalidateQueries({ queryKey: ['deck-stats', addingToDeck.id] })
        }}
      />
    )
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
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-4xl mx-auto">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Microdote</h1>
          <p className="text-gray-600">Spaced repetition flashcards</p>
        </header>

        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4">Create New Deck</h2>
          <form onSubmit={handleCreateDeck} className="flex gap-2">
            <input
              type="text"
              value={newDeckName}
              onChange={(e) => setNewDeckName(e.target.value)}
              placeholder="Enter deck name..."
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={createDeckMutation.isPending}
            />
            <button
              type="submit"
              disabled={createDeckMutation.isPending || !newDeckName.trim()}
              className="px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              {createDeckMutation.isPending ? 'Creating...' : 'Create'}
            </button>
          </form>
          {createDeckMutation.isError && (
            <p className="mt-2 text-sm text-red-600">
              Failed to create deck. Please try again.
            </p>
          )}
        </div>

        <div className="bg-white rounded-lg shadow">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-xl font-semibold">Your Decks</h2>
          </div>

          {!decks || decks.length === 0 ? (
            <div className="p-8 text-center text-gray-500">
              No decks yet. Create one to get started!
            </div>
          ) : (
            <ul className="divide-y divide-gray-200">
              {decks.map((deck: Deck) => (
                <DeckItem key={deck.id} deck={deck} onStudy={setStudyingDeck} onAddCards={setAddingToDeck} />
              ))}
            </ul>
          )}
        </div>

        <footer className="mt-8 text-center text-sm text-gray-500">
          <p>Milestone M1 — Studying MVP</p>
          <p className="mt-1">Backend: Go + SQLite | Frontend: React + TanStack Query</p>
        </footer>
      </div>
    </div>
  )
}
