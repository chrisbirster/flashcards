import { useNavigate } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import type { Deck } from '#/lib/api'
import { DeckStatItem } from './deck-stat-item'
import { useAppRepository } from '#/lib/app-repository'

export function DeckItem({ deck }: { deck: Deck }) {
  const navigate = useNavigate()
  const repository = useAppRepository()
  const { data: stats } = useQuery({
    queryKey: ['deck-stats', deck.id],
    queryFn: () => repository.fetchDeckStats(deck.id),
  })

  return (
    <li className="p-4 sm:p-6 hover:bg-gray-50 transition-colors">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex-1">
          <h3 className="text-lg font-medium text-gray-900">{deck.name}</h3>
          <div className="flex flex-wrap gap-x-4 gap-y-1 mt-2 text-sm">
            {stats ? (
              <>
                <DeckStatItem stat={stats.newCards} label={"new"} color={"text-blue-600"} />
                <DeckStatItem stat={stats.learning} label={"review"} color={"text-orange-600"} />
                <DeckStatItem stat={stats.review} label={"review"} color={"text-green-600"} />

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
        <div className="flex w-full flex-col gap-2 sm:w-auto sm:flex-row">
          <button
            className="px-4 py-2 text-sm text-white bg-blue-600 hover:bg-blue-700 rounded-md disabled:bg-gray-300 w-full sm:w-auto"
            disabled={!stats || stats.dueToday === 0}
            onClick={() => navigate(`/study/${deck.id}`)}
          >
            Study Now
          </button>
          <button
            className="px-4 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-md w-full sm:w-auto"
            onClick={() => navigate(`/notes/add?deckId=${deck.id}`)}
          >
            Add Cards
          </button>
        </div>
      </div>
    </li>
  )
}
