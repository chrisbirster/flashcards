import { useParams, useNavigate } from 'react-router-dom'
import { useQuery } from '@tanstack/react-query'
import { fetchDecks } from '#/lib/api'
import { StudyScreen } from '#/components/StudyScreen'

export function StudyPage() {
  const { deckId } = useParams<{ deckId: string }>()
  const navigate = useNavigate()

  const { data: decks, isLoading } = useQuery({
    queryKey: ['decks'],
    queryFn: fetchDecks,
  })

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading...</div>
      </div>
    )
  }

  const deck = decks?.find(d => d.id === Number(deckId))

  if (!deck) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="text-red-600 mb-4">Deck not found</div>
          <button
            onClick={() => navigate('/decks')}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Back to Decks
          </button>
        </div>
      </div>
    )
  }

  return (
    <StudyScreen
      deckId={deck.id}
      deckName={deck.name}
      onExit={() => navigate('/decks')}
    />
  )
}
