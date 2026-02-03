import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { fetchDueCards, answerCard, updateCard } from '#/lib/api'

// Flag colors matching Anki
const FLAG_COLORS = [
  { id: 0, name: 'None', color: 'bg-gray-200', textColor: 'text-gray-600' },
  { id: 1, name: 'Red', color: 'bg-red-500', textColor: 'text-white' },
  { id: 2, name: 'Orange', color: 'bg-orange-500', textColor: 'text-white' },
  { id: 3, name: 'Green', color: 'bg-green-500', textColor: 'text-white' },
  { id: 4, name: 'Blue', color: 'bg-blue-500', textColor: 'text-white' },
  { id: 5, name: 'Pink', color: 'bg-pink-500', textColor: 'text-white' },
  { id: 6, name: 'Turquoise', color: 'bg-teal-500', textColor: 'text-white' },
  { id: 7, name: 'Purple', color: 'bg-purple-500', textColor: 'text-white' },
]

interface StudyScreenProps {
  deckId: number
  deckName: string
  onExit: () => void
}

export function StudyScreen({ deckId, deckName, onExit }: StudyScreenProps) {
  const [currentCardIndex, setCurrentCardIndex] = useState(0)
  const [showAnswer, setShowAnswer] = useState(false)
  const [questionStartTime, setQuestionStartTime] = useState<number>(Date.now())
  const [showFlagMenu, setShowFlagMenu] = useState(false)
  const queryClient = useQueryClient()

  const { data: cards, isLoading } = useQuery({
    queryKey: ['due-cards', deckId],
    queryFn: () => fetchDueCards(deckId, 50),
  })

  const answerMutation = useMutation({
    mutationFn: ({ cardId, rating, timeTakenMs }: { cardId: number; rating: number; timeTakenMs: number }) =>
      answerCard(cardId, { rating, timeTakenMs }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['due-cards', deckId] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats', deckId] })

      // Move to next card
      if (cards && currentCardIndex < cards.length - 1) {
        setCurrentCardIndex(currentCardIndex + 1)
        setShowAnswer(false)
        setQuestionStartTime(Date.now()) // Reset timer for next card
      } else {
        // No more cards, exit study
        onExit()
      }
    },
  })

  const updateCardMutation = useMutation({
    mutationFn: ({ cardId, flag, marked, suspended }: { cardId: number; flag?: number; marked?: boolean; suspended?: boolean }) =>
      updateCard(cardId, { flag, marked, suspended }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['due-cards', deckId] })
      setShowFlagMenu(false)
    },
  })

  const currentCard = cards?.[currentCardIndex]

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyPress = (e: KeyboardEvent) => {
      if (!currentCard) return

      // M: toggle mark
      if (e.key === 'm' || e.key === 'M') {
        e.preventDefault()
        updateCardMutation.mutate({ cardId: currentCard.id, marked: !currentCard.marked })
        return
      }

      // Escape: close flag menu
      if (e.key === 'Escape' && showFlagMenu) {
        setShowFlagMenu(false)
        return
      }

      // @ or Shift+2: suspend card
      if (e.key === '@') {
        e.preventDefault()
        updateCardMutation.mutate({ cardId: currentCard.id, suspended: true })
        return
      }

      // Space or Enter: show answer or select "Good"
      if (e.key === ' ' || e.key === 'Enter') {
        e.preventDefault()
        if (!showAnswer) {
          setShowAnswer(true)
        } else {
          handleAnswer(3) // Good
        }
        return
      }

      // Number keys: select answer
      if (showAnswer && ['1', '2', '3', '4'].includes(e.key)) {
        handleAnswer(parseInt(e.key))
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [currentCard, showAnswer, showFlagMenu])

  const handleAnswer = (rating: number) => {
    if (!currentCard) return

    const timeTakenMs = Date.now() - questionStartTime
    answerMutation.mutate({
      cardId: currentCard.id,
      rating,
      timeTakenMs,
    })
  }

  // Reset timer when current card changes
  useEffect(() => {
    if (currentCard) {
      setQuestionStartTime(Date.now())
    }
  }, [currentCard?.id])

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-gray-600">Loading cards...</div>
      </div>
    )
  }

  if (!cards || cards.length === 0) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center bg-gray-50">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-gray-900 mb-4">All done!</h2>
          <p className="text-gray-600 mb-8">No more cards due for this deck.</p>
          <button
            onClick={onExit}
            className="px-6 py-3 bg-blue-600 text-white rounded-md hover:bg-blue-700"
          >
            Back to Decks
          </button>
        </div>
      </div>
    )
  }

  if (!currentCard) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-red-600">Error loading card</div>
      </div>
    )
  }

  const progress = `${currentCardIndex + 1} / ${cards.length}`

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-3xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-xl font-semibold text-gray-900">{deckName}</h1>
            <p className="text-sm text-gray-600">Card {progress}</p>
          </div>
          <button
            onClick={onExit}
            className="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded-md"
          >
            Exit
          </button>
        </div>

        {/* Card Tools */}
        <div className="flex items-center gap-2 mb-4">
          {/* Flag button with dropdown */}
          <div className="relative">
            <button
              onClick={() => setShowFlagMenu(!showFlagMenu)}
              className={`px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1 ${
                currentCard.flag > 0
                  ? `${FLAG_COLORS[currentCard.flag].color} ${FLAG_COLORS[currentCard.flag].textColor}`
                  : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
              }`}
              data-testid="flag-button"
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M3 6a3 3 0 013-3h10a1 1 0 01.8 1.6L14.25 8l2.55 3.4A1 1 0 0116 13H6a1 1 0 00-1 1v3a1 1 0 11-2 0V6z" clipRule="evenodd" />
              </svg>
              Flag
            </button>
            {showFlagMenu && (
              <div className="absolute left-0 top-full mt-1 bg-white rounded-lg shadow-lg border z-10 py-1 min-w-[120px]">
                {FLAG_COLORS.map((flag) => (
                  <button
                    key={flag.id}
                    onClick={() => {
                      updateCardMutation.mutate({ cardId: currentCard.id, flag: flag.id })
                    }}
                    className="w-full px-3 py-1.5 text-left text-sm flex items-center gap-2 hover:bg-gray-100"
                  >
                    <span className={`w-3 h-3 rounded-full ${flag.color}`}></span>
                    {flag.name}
                  </button>
                ))}
              </div>
            )}
          </div>

          {/* Mark button */}
          <button
            onClick={() => updateCardMutation.mutate({ cardId: currentCard.id, marked: !currentCard.marked })}
            className={`px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1 ${
              currentCard.marked
                ? 'bg-yellow-500 text-white'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
            data-testid="mark-button"
          >
            <svg className="w-4 h-4" fill={currentCard.marked ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
            Mark
          </button>

          {/* Suspend button */}
          <button
            onClick={() => {
              updateCardMutation.mutate({ cardId: currentCard.id, suspended: true })
            }}
            className="px-3 py-1.5 rounded-md text-sm font-medium bg-gray-100 text-gray-600 hover:bg-gray-200 flex items-center gap-1"
            data-testid="suspend-button"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Suspend
          </button>
        </div>

        {/* Card */}
        <div className="bg-white rounded-lg shadow-lg p-8 mb-6">
          {/* Question */}
          <div className="mb-8">
            <div className="text-sm text-gray-500 mb-2">Question</div>
            <div
              className="text-xl prose max-w-none"
              dangerouslySetInnerHTML={{ __html: currentCard.front }}
            />
          </div>

          {/* Answer (revealed) */}
          {showAnswer && (
            <div className="border-t pt-8">
              <div className="text-sm text-gray-500 mb-2">Answer</div>
              <div
                className="text-xl prose max-w-none"
                dangerouslySetInnerHTML={{ __html: currentCard.back }}
              />
            </div>
          )}
        </div>

        {/* Actions */}
        {!showAnswer ? (
          <div className="text-center">
            <button
              onClick={() => setShowAnswer(true)}
              className="px-8 py-4 bg-blue-600 text-white text-lg rounded-lg hover:bg-blue-700 font-medium"
            >
              Show Answer
            </button>
            <p className="mt-4 text-sm text-gray-500">
              Press <kbd className="px-2 py-1 bg-gray-200 rounded">Space</kbd> or{' '}
              <kbd className="px-2 py-1 bg-gray-200 rounded">Enter</kbd>
            </p>
          </div>
        ) : (
          <div>
            <div className="flex gap-4 justify-center">
              <button
                onClick={() => handleAnswer(1)}
                disabled={answerMutation.isPending}
                className="flex-1 max-w-xs px-6 py-4 bg-red-100 text-red-800 rounded-lg hover:bg-red-200 disabled:opacity-50 font-medium"
              >
                <div className="text-lg">Again</div>
                <div className="text-xs mt-1">1</div>
              </button>
              <button
                onClick={() => handleAnswer(2)}
                disabled={answerMutation.isPending}
                className="flex-1 max-w-xs px-6 py-4 bg-orange-100 text-orange-800 rounded-lg hover:bg-orange-200 disabled:opacity-50 font-medium"
              >
                <div className="text-lg">Hard</div>
                <div className="text-xs mt-1">2</div>
              </button>
              <button
                onClick={() => handleAnswer(3)}
                disabled={answerMutation.isPending}
                className="flex-1 max-w-xs px-6 py-4 bg-green-100 text-green-800 rounded-lg hover:bg-green-200 disabled:opacity-50 font-medium"
              >
                <div className="text-lg">Good</div>
                <div className="text-xs mt-1">3</div>
              </button>
              <button
                onClick={() => handleAnswer(4)}
                disabled={answerMutation.isPending}
                className="flex-1 max-w-xs px-6 py-4 bg-blue-100 text-blue-800 rounded-lg hover:bg-blue-200 disabled:opacity-50 font-medium"
              >
                <div className="text-lg">Easy</div>
                <div className="text-xs mt-1">4</div>
              </button>
            </div>
            <p className="mt-4 text-center text-sm text-gray-500">
              Press <kbd className="px-2 py-1 bg-gray-200 rounded">1</kbd>-
              <kbd className="px-2 py-1 bg-gray-200 rounded">4</kbd> or{' '}
              <kbd className="px-2 py-1 bg-gray-200 rounded">Space</kbd> for Good
            </p>
          </div>
        )}
      </div>
    </div>
  )
}
