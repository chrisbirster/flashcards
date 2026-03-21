import { useState, useEffect, useRef, useCallback } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import DOMPurify from 'dompurify'
import { useAppRepository } from '#/lib/app-repository'
import { ActionBar } from '#/components/action-bar'
import { EmptyState, PageContainer, PageSection } from '#/components/page-layout'

const FLAG_COLORS = [
  { id: 0, name: 'None', color: 'bg-[var(--app-card-strong)] text-[var(--app-text-soft)] border-[var(--app-line-strong)]' },
  { id: 1, name: 'Red', color: 'bg-rose-500 text-white border-rose-500' },
  { id: 2, name: 'Orange', color: 'bg-orange-500 text-white border-orange-500' },
  { id: 3, name: 'Green', color: 'bg-emerald-500 text-white border-emerald-500' },
  { id: 4, name: 'Blue', color: 'bg-sky-500 text-white border-sky-500' },
  { id: 5, name: 'Pink', color: 'bg-pink-500 text-white border-pink-500' },
  { id: 6, name: 'Turquoise', color: 'bg-teal-500 text-white border-teal-500' },
  { id: 7, name: 'Purple', color: 'bg-violet-500 text-white border-violet-500' },
]

const ANSWER_BUTTONS = [
  {
    rating: 1,
    label: 'Again',
    shortcut: '1',
    className: 'border-rose-500/40 bg-rose-500/12 text-rose-200',
  },
  {
    rating: 2,
    label: 'Hard',
    shortcut: '2',
    className: 'border-orange-500/40 bg-orange-500/12 text-orange-200',
  },
  {
    rating: 3,
    label: 'Good',
    shortcut: '3',
    className: 'border-emerald-500/40 bg-emerald-500/12 text-emerald-200',
  },
  {
    rating: 4,
    label: 'Easy',
    shortcut: '4',
    className: 'border-sky-500/40 bg-sky-500/12 text-sky-200',
  },
]

interface StudyScreenProps {
  deckId: number
  deckName: string
  onExit: () => void
}

export function StudyScreen({ deckId, deckName, onExit }: StudyScreenProps) {
  const [currentCardIndex, setCurrentCardIndex] = useState(0)
  const [showAnswer, setShowAnswer] = useState(false)
  const [showFlagMenu, setShowFlagMenu] = useState(false)
  const questionStartTimeRef = useRef(0)
  const queryClient = useQueryClient()
  const repository = useAppRepository()

  const { data: cards, isLoading } = useQuery({
    queryKey: ['due-cards', deckId],
    queryFn: () => repository.fetchDueCards(deckId, 50),
  })

  const answerMutation = useMutation({
    mutationFn: ({ cardId, rating, timeTakenMs }: { cardId: number; rating: number; timeTakenMs: number }) =>
      repository.answerCard(cardId, { rating, timeTakenMs }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['due-cards', deckId] })
      queryClient.invalidateQueries({ queryKey: ['deck-stats', deckId] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })

      if (cards && currentCardIndex < cards.length - 1) {
        setCurrentCardIndex(currentCardIndex + 1)
        setShowAnswer(false)
        questionStartTimeRef.current = Date.now()
      } else {
        onExit()
      }
    },
  })

  const updateCardMutation = useMutation({
    mutationFn: ({ cardId, flag, marked, suspended }: { cardId: number; flag?: number; marked?: boolean; suspended?: boolean }) =>
      repository.updateCard(cardId, { flag, marked, suspended }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['due-cards', deckId] })
      queryClient.invalidateQueries({ queryKey: ['dashboard'] })
      setShowFlagMenu(false)
    },
  })

  const currentCard = cards?.[currentCardIndex]
  const progress = `${currentCardIndex + 1} / ${cards?.length ?? 0}`

  const handleAnswer = useCallback((rating: number) => {
    if (!currentCard) return

    const startedAt = questionStartTimeRef.current || Date.now()
    const timeTakenMs = Math.max(0, Date.now() - startedAt)
    answerMutation.mutate({
      cardId: currentCard.id,
      rating,
      timeTakenMs,
    })
  }, [answerMutation, currentCard])

  useEffect(() => {
    const handleKeyPress = (event: KeyboardEvent) => {
      if (!currentCard) return

      if (event.key === 'm' || event.key === 'M') {
        event.preventDefault()
        updateCardMutation.mutate({ cardId: currentCard.id, marked: !currentCard.marked })
        return
      }

      if (event.key === 'Escape' && showFlagMenu) {
        setShowFlagMenu(false)
        return
      }

      if (event.key === '@') {
        event.preventDefault()
        updateCardMutation.mutate({ cardId: currentCard.id, suspended: true })
        return
      }

      if (event.key === ' ' || event.key === 'Enter') {
        event.preventDefault()
        if (!showAnswer) {
          setShowAnswer(true)
        } else {
          handleAnswer(3)
        }
        return
      }

      if (showAnswer && ['1', '2', '3', '4'].includes(event.key)) {
        handleAnswer(parseInt(event.key))
      }
    }

    window.addEventListener('keydown', handleKeyPress)
    return () => window.removeEventListener('keydown', handleKeyPress)
  }, [currentCard, handleAnswer, showAnswer, showFlagMenu, updateCardMutation])

  useEffect(() => {
    if (currentCard) {
      questionStartTimeRef.current = Date.now()
    }
  }, [currentCard])

  if (isLoading) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-text-soft)]">
          Loading study queue...
        </PageSection>
      </PageContainer>
    )
  }

  if (!cards || cards.length === 0) {
    return (
      <PageContainer>
        <EmptyState
          title="All caught up"
          description="There are no due cards in this deck right now."
          action={
            <button
              type="button"
              onClick={onExit}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Back to Decks
            </button>
          }
        />
      </PageContainer>
    )
  }

  if (!currentCard) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-danger-text)]">
          Error loading card.
        </PageSection>
      </PageContainer>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-4 sm:p-5">
        <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
          <div>
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Study session</p>
            <h1 className="mt-2 text-2xl font-semibold tracking-tight text-[var(--app-text)]">{deckName}</h1>
            <p className="mt-2 text-sm text-[var(--app-text-soft)]">Card {progress}</p>
          </div>
          <button
            type="button"
            onClick={onExit}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Exit
          </button>
        </div>
      </PageSection>

      <PageSection className="p-4 sm:p-5">
        <div className="flex flex-wrap items-center gap-2">
          <div className="relative">
            <button
              type="button"
              onClick={() => setShowFlagMenu((current) => !current)}
              className={`inline-flex min-h-11 items-center gap-2 rounded-2xl border px-4 text-sm font-medium ${
                currentCard.flag > 0
                  ? FLAG_COLORS[currentCard.flag].color
                  : 'border-[var(--app-line-strong)] bg-[var(--app-card-strong)] text-[var(--app-text-soft)]'
              }`}
              data-testid="flag-button"
            >
              <svg className="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M3 6a3 3 0 013-3h10a1 1 0 01.8 1.6L14.25 8l2.55 3.4A1 1 0 0116 13H6a1 1 0 00-1 1v3a1 1 0 11-2 0V6z" clipRule="evenodd" />
              </svg>
              Flag
            </button>
            {showFlagMenu ? (
              <div className="absolute left-0 top-full z-10 mt-2 min-w-[12rem] overflow-hidden rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-panel)] shadow-2xl">
                {FLAG_COLORS.map((flag) => (
                  <button
                    key={flag.id}
                    type="button"
                    onClick={() => updateCardMutation.mutate({ cardId: currentCard.id, flag: flag.id })}
                    className="flex w-full items-center gap-2 border-b border-[var(--app-line)] px-4 py-3 text-left text-sm text-[var(--app-text)] last:border-b-0 hover:bg-[var(--app-card-strong)]"
                  >
                    <span className={`h-3 w-3 rounded-full ${flag.id === 0 ? 'border border-[var(--app-line-strong)] bg-transparent' : flag.color.split(' ')[0]}`} />
                    {flag.name}
                  </button>
                ))}
              </div>
            ) : null}
          </div>

          <button
            type="button"
            onClick={() => updateCardMutation.mutate({ cardId: currentCard.id, marked: !currentCard.marked })}
            className={`inline-flex min-h-11 items-center gap-2 rounded-2xl border px-4 text-sm font-medium ${
              currentCard.marked
                ? 'border-amber-500/40 bg-amber-500/12 text-amber-200'
                : 'border-[var(--app-line-strong)] bg-[var(--app-card-strong)] text-[var(--app-text-soft)]'
            }`}
            data-testid="mark-button"
          >
            <svg className="h-4 w-4" fill={currentCard.marked ? 'currentColor' : 'none'} stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
            Mark
          </button>

          <button
            type="button"
            onClick={() => updateCardMutation.mutate({ cardId: currentCard.id, suspended: true })}
            className="inline-flex min-h-11 items-center gap-2 rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
            data-testid="suspend-button"
          >
            <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Suspend
          </button>
        </div>
      </PageSection>

      <PageSection className="overflow-hidden">
        <div className="border-b border-[var(--app-line)] px-4 py-3 text-xs uppercase tracking-[0.18em] text-[var(--app-muted)] sm:px-5">
          Question
        </div>
        <div className="px-4 py-6 sm:px-5 sm:py-8">
          <div
            className="prose prose-invert max-w-none text-base leading-8 text-[var(--app-text)] sm:text-xl"
            dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(currentCard.front) }}
          />
        </div>

        {showAnswer ? (
          <>
            <div className="border-b border-t border-[var(--app-line)] px-4 py-3 text-xs uppercase tracking-[0.18em] text-[var(--app-muted)] sm:px-5">
              Answer
            </div>
            <div className="px-4 py-6 sm:px-5 sm:py-8">
              <div
                className="prose prose-invert max-w-none text-base leading-8 text-[var(--app-text)] sm:text-xl"
                dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(currentCard.back) }}
              />
            </div>
          </>
        ) : null}
      </PageSection>

      {!showAnswer ? (
        <ActionBar>
          <div className="space-y-3">
            <button
              type="button"
              onClick={() => setShowAnswer(true)}
              className="inline-flex min-h-12 w-full items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-base font-semibold text-[var(--app-accent-ink)]"
            >
              Show Answer
            </button>
            <p className="text-center text-xs text-[var(--app-muted)]">
              Press <kbd className="rounded-lg border border-[var(--app-line-strong)] bg-[var(--app-card)] px-2 py-1">Space</kbd> or{' '}
              <kbd className="rounded-lg border border-[var(--app-line-strong)] bg-[var(--app-card)] px-2 py-1">Enter</kbd>
            </p>
          </div>
        </ActionBar>
      ) : (
        <ActionBar>
          <div className="space-y-3">
            <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
              {ANSWER_BUTTONS.map((button) => (
                <button
                  key={button.rating}
                  type="button"
                  onClick={() => handleAnswer(button.rating)}
                  disabled={answerMutation.isPending}
                  className={`inline-flex min-h-14 flex-col items-center justify-center rounded-2xl border px-4 py-3 text-sm font-semibold ${button.className} disabled:cursor-not-allowed disabled:opacity-50`}
                >
                  <span className="text-base">{button.label}</span>
                  <span className="mt-1 text-xs opacity-80">{button.shortcut}</span>
                </button>
              ))}
            </div>
            <p className="text-center text-xs text-[var(--app-muted)]">
              Press <kbd className="rounded-lg border border-[var(--app-line-strong)] bg-[var(--app-card)] px-2 py-1">1</kbd>-
              <kbd className="rounded-lg border border-[var(--app-line-strong)] bg-[var(--app-card)] px-2 py-1">4</kbd> or{' '}
              <kbd className="rounded-lg border border-[var(--app-line-strong)] bg-[var(--app-card)] px-2 py-1">Space</kbd> for Good
            </p>
          </div>
        </ActionBar>
      )}
    </PageContainer>
  )
}
