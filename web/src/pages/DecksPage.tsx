import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchDecks, fetchDeckStats, createDeck, type Deck } from '#/lib/api'

function DeckItem({ deck }: { deck: Deck }) {
  const navigate = useNavigate()
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
            onClick={() => navigate(`/study/${deck.id}`)}
          >
            Study Now
          </button>
          <button
            className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-md"
            onClick={() => navigate(`/notes/add?deckId=${deck.id}`)}
          >
            Add Cards
          </button>
        </div>
      </div>
    </li>
  )
}

export function DecksPage() {
  const [newDeckName, setNewDeckName] = useState('')
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
            disabled={createDeckMutation.isPending}
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
