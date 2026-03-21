import { useParams, useNavigate } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import { StudyScreen } from '#/components/StudyScreen'
import { useAppRepository } from '#/lib/app-repository'
import { EmptyState, PageContainer, PageSection } from '#/components/page-layout'

export function StudyPage() {
  const { deckId } = useParams<{ deckId: string }>()
  const navigate = useNavigate()
  const repository = useAppRepository()

  const { data: decks, isLoading } = useQuery({
    queryKey: ['decks'],
    queryFn: () => repository.fetchDecks(),
  })

  if (isLoading) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-text-soft)]">
          Loading study deck...
        </PageSection>
      </PageContainer>
    )
  }

  const deck = decks?.find(d => d.id === Number(deckId))

  if (!deck) {
    return (
      <PageContainer>
        <EmptyState
          title="Deck not found"
          description="The deck you tried to study is missing or no longer accessible."
          action={
            <button
              type="button"
              onClick={() => navigate('/decks')}
              className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
            >
              Back to Decks
            </button>
          }
        />
      </PageContainer>
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
