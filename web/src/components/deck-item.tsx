import { useState } from 'react'
import { useNavigate } from 'react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import type { Deck } from '#/lib/api'
import { DeckStatItem } from './deck-stat-item'
import { useAppRepository } from '#/lib/app-repository'

export function DeckItem({ deck }: { deck: Deck }) {
  const navigate = useNavigate()
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const [isEditing, setIsEditing] = useState(false)
  const [draftName, setDraftName] = useState(deck.name)
  const [actionError, setActionError] = useState<string | null>(null)

  const { data: stats } = useQuery({
    queryKey: ['deck-stats', deck.id],
    queryFn: () => repository.fetchDeckStats(deck.id),
  })

  const renameMutation = useMutation({
    mutationFn: (name: string) => repository.updateDeck(deck.id, { name }),
    onSuccess: () => {
      setIsEditing(false)
      setActionError(null)
      queryClient.invalidateQueries({ queryKey: ['decks'] })
    },
    onError: (error: Error) => setActionError(error.message),
  })

  const deleteMutation = useMutation({
    mutationFn: () => repository.deleteDeck(deck.id),
    onSuccess: () => {
      setActionError(null)
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['entitlements'] })
    },
    onError: (error: Error) => setActionError(error.message),
  })

  const isEmptyDeck = (stats?.totalCards ?? 0) === 0

  return (
    <li className="p-4 sm:p-6 hover:bg-gray-50 transition-colors">
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
          <div className="flex-1">
            {isEditing ? (
              <form
                onSubmit={(event) => {
                  event.preventDefault()
                  if (!draftName.trim()) return
                  renameMutation.mutate(draftName.trim())
                }}
                className="flex flex-col gap-2 sm:flex-row"
              >
                <input
                  type="text"
                  value={draftName}
                  onChange={(event) => setDraftName(event.target.value)}
                  className="w-full rounded-xl border border-gray-300 px-3 py-2 text-sm focus:border-gray-500 focus:outline-none sm:max-w-sm"
                />
                <div className="flex gap-2">
                  <button
                    type="submit"
                    disabled={renameMutation.isPending || !draftName.trim()}
                    className="rounded-xl bg-slate-950 px-3 py-2 text-sm font-medium text-white hover:bg-slate-800 disabled:bg-gray-300"
                  >
                    {renameMutation.isPending ? 'Saving...' : 'Save'}
                  </button>
                  <button
                    type="button"
                    onClick={() => {
                      setIsEditing(false)
                      setDraftName(deck.name)
                    }}
                    className="rounded-xl border border-gray-300 px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
                  >
                    Cancel
                  </button>
                </div>
              </form>
            ) : (
              <h3 className="text-lg font-medium text-gray-900">{deck.name}</h3>
            )}

            <div className="mt-2 flex flex-wrap gap-x-4 gap-y-1 text-sm">
              {stats ? (
                <>
                  <DeckStatItem stat={stats.newCards} label={'new'} color={'text-blue-600'} />
                  <DeckStatItem stat={stats.learning} label={'learning'} color={'text-orange-600'} />
                  <DeckStatItem stat={stats.review} label={'review'} color={'text-green-600'} />

                  {stats.suspended > 0 && <span className="text-gray-500">{stats.suspended} suspended</span>}
                  <span className="text-gray-400">({stats.totalCards} total)</span>
                </>
              ) : (
                <span className="text-gray-500">Loading stats...</span>
              )}
            </div>

            {stats && stats.dueToday > 0 && (
              <div className="mt-1 text-sm font-semibold text-indigo-600">{stats.dueToday} due today</div>
            )}
            {!isEmptyDeck && (
              <div className="mt-2 text-xs text-gray-500">
                Delete is disabled until this deck is empty.
              </div>
            )}
            {actionError && (
              <div className="mt-2 text-sm text-red-600">
                {actionError}
              </div>
            )}
          </div>

          <div className="flex w-full flex-col gap-2 sm:w-auto sm:flex-row sm:flex-wrap sm:justify-end">
            <button
              className="w-full rounded-md bg-blue-600 px-4 py-2 text-sm text-white hover:bg-blue-700 disabled:bg-gray-300 sm:w-auto"
              disabled={!stats || stats.dueToday === 0}
              onClick={() => navigate(`/study/${deck.id}`)}
            >
              Study
            </button>
            <button
              className="w-full rounded-md px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 sm:w-auto"
              onClick={() => navigate(`/notes/add?deckId=${deck.id}`)}
            >
              Add Note
            </button>
            <button
              className="w-full rounded-md border border-gray-300 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 sm:w-auto"
              onClick={() => {
                setActionError(null)
                setIsEditing(true)
              }}
            >
              Rename
            </button>
            <button
              className="w-full rounded-md border border-red-200 px-4 py-2 text-sm text-red-700 hover:bg-red-50 disabled:cursor-not-allowed disabled:border-gray-200 disabled:text-gray-400 sm:w-auto"
              disabled={!isEmptyDeck || deleteMutation.isPending}
              onClick={() => {
                if (!window.confirm(`Delete the deck "${deck.name}"?`)) {
                  return
                }
                deleteMutation.mutate()
              }}
            >
              {deleteMutation.isPending ? 'Deleting...' : 'Delete'}
            </button>
          </div>
        </div>
      </div>
    </li>
  )
}
