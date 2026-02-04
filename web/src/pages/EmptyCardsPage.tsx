import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { findEmptyCards, deleteEmptyCards } from '#/lib/api'
import type { EmptyCardInfo } from '#/lib/api'
import DOMPurify from 'dompurify'

export function EmptyCardsPage() {
  const queryClient = useQueryClient()
  const [selectedCards, setSelectedCards] = useState<Set<number>>(new Set())
  const [showConfirmDialog, setShowConfirmDialog] = useState(false)

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ['empty-cards'],
    queryFn: findEmptyCards,
  })

  const deleteCardsMutation = useMutation({
    mutationFn: (cardIds: number[]) => deleteEmptyCards({ cardIds }),
    onSuccess: () => {
      setSelectedCards(new Set())
      setShowConfirmDialog(false)
      // Refetch to update the list
      refetch()
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['decks'] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats'] })
    },
  })

  const handleSelectAll = () => {
    if (data && selectedCards.size === data.emptyCards.length) {
      setSelectedCards(new Set())
    } else if (data) {
      setSelectedCards(new Set(data.emptyCards.map(c => c.cardId)))
    }
  }

  const handleSelectCard = (cardId: number) => {
    const newSelected = new Set(selectedCards)
    if (newSelected.has(cardId)) {
      newSelected.delete(cardId)
    } else {
      newSelected.add(cardId)
    }
    setSelectedCards(newSelected)
  }

  const handleDelete = () => {
    if (selectedCards.size > 0) {
      setShowConfirmDialog(true)
    }
  }

  const confirmDelete = () => {
    deleteCardsMutation.mutate(Array.from(selectedCards))
  }

  if (isLoading) {
    return (
      <div className="max-w-5xl mx-auto">
        <div className="text-gray-600">Scanning for empty cards...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="max-w-5xl mx-auto">
        <div className="text-red-600">
          Error: {error instanceof Error ? error.message : 'Failed to find empty cards'}
        </div>
      </div>
    )
  }

  return (
    <div className="max-w-5xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-2">Empty Cards</h1>
        <p className="text-gray-600">
          Cards with no content or missing cloze deletions. These can be safely deleted.
        </p>
      </div>

      {data && data.count > 0 ? (
        <>
          <div className="bg-white rounded-lg shadow border border-gray-200">
            <div className="p-4 border-b border-gray-200 flex items-center justify-between">
              <div className="flex items-center gap-4">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input
                    type="checkbox"
                    checked={selectedCards.size === data.emptyCards.length && data.emptyCards.length > 0}
                    onChange={handleSelectAll}
                    className="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
                  />
                  <span className="text-sm text-gray-700">
                    Select All ({data.count} card{data.count !== 1 ? 's' : ''})
                  </span>
                </label>
                {selectedCards.size > 0 && (
                  <span className="text-sm text-blue-600 font-medium">
                    {selectedCards.size} selected
                  </span>
                )}
              </div>
              <button
                onClick={handleDelete}
                disabled={selectedCards.size === 0 || deleteCardsMutation.isPending}
                className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-sm font-medium"
              >
                {deleteCardsMutation.isPending ? 'Deleting...' : `Delete Selected (${selectedCards.size})`}
              </button>
            </div>

            <div className="divide-y divide-gray-200">
              {data.emptyCards.map((card) => (
                <EmptyCardItem
                  key={card.cardId}
                  card={card}
                  selected={selectedCards.has(card.cardId)}
                  onSelect={handleSelectCard}
                />
              ))}
            </div>
          </div>

          {deleteCardsMutation.isError && (
            <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm">
              Error deleting cards: {deleteCardsMutation.error instanceof Error ? deleteCardsMutation.error.message : 'Unknown error'}
            </div>
          )}
        </>
      ) : (
        <div className="bg-white rounded-lg shadow p-8 text-center">
          <svg className="mx-auto h-12 w-12 text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <h3 className="text-lg font-medium text-gray-900 mb-2">No Empty Cards Found</h3>
          <p className="text-gray-500">All your cards have content. Great job!</p>
        </div>
      )}

      {/* Confirmation Dialog */}
      {showConfirmDialog && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
            <h3 className="text-lg font-semibold text-gray-900 mb-2">Confirm Deletion</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to delete {selectedCards.size} empty card{selectedCards.size !== 1 ? 's' : ''}? 
              This action cannot be undone.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setShowConfirmDialog(false)}
                className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                disabled={deleteCardsMutation.isPending}
              >
                Cancel
              </button>
              <button
                onClick={confirmDelete}
                className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700"
                disabled={deleteCardsMutation.isPending}
              >
                {deleteCardsMutation.isPending ? 'Deleting...' : 'Delete'}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

interface EmptyCardItemProps {
  card: EmptyCardInfo
  selected: boolean
  onSelect: (cardId: number) => void
}

function EmptyCardItem({ card, selected, onSelect }: EmptyCardItemProps) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="p-4 hover:bg-gray-50">
      <div className="flex items-start gap-3">
        <input
          type="checkbox"
          checked={selected}
          onChange={() => onSelect(card.cardId)}
          className="mt-1 w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
        />
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2 mb-2">
            <div>
              <span className="text-sm font-medium text-gray-900">
                Card #{card.cardId}
              </span>
              {card.ordinal > 0 && (
                <span className="ml-2 text-xs text-purple-600 bg-purple-100 px-2 py-0.5 rounded">
                  Cloze c{card.ordinal}
                </span>
              )}
              <p className="text-sm text-gray-600 mt-1">
                Template: {card.templateName}
              </p>
            </div>
            <button
              onClick={() => setExpanded(!expanded)}
              className="text-sm text-blue-600 hover:text-blue-700"
            >
              {expanded ? 'Hide' : 'Show'} Preview
            </button>
          </div>

          <div className="p-2 bg-amber-50 border border-amber-200 rounded text-sm text-amber-800">
            <span className="font-medium">Reason:</span> {card.reason}
          </div>

          {expanded && (
            <div className="mt-3 space-y-2">
              <div className="border rounded-md overflow-hidden">
                <div className="bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 border-b">
                  Front
                </div>
                <div 
                  className="p-3 text-sm"
                  dangerouslySetInnerHTML={{
                    __html: DOMPurify.sanitize(card.front || '<span class="text-gray-400">(empty)</span>')
                  }}
                />
              </div>
              <div className="border rounded-md overflow-hidden">
                <div className="bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600 border-b">
                  Back
                </div>
                <div 
                  className="p-3 text-sm"
                  dangerouslySetInnerHTML={{
                    __html: DOMPurify.sanitize(card.back || '<span class="text-gray-400">(empty)</span>')
                  }}
                />
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
