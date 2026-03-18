import { describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { RecentDeckNotesPanel } from '#/components/recent-deck-notes-panel'
import { AppRepositoryProvider, type AppRepository } from '#/lib/app-repository'

function renderPanel(repository: AppRepository) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <AppRepositoryProvider repository={repository}>
        <RecentDeckNotesPanel deckId={1} />
      </AppRepositoryProvider>
    </QueryClientProvider>,
  )
}

describe('RecentDeckNotesPanel', () => {
  it('shows a loading state while recent notes are fetching', () => {
    const repository = {
      fetchDeckNotes: vi.fn(() => new Promise(() => {})),
    } as unknown as AppRepository

    renderPanel(repository)

    expect(screen.getByText('Loading recent notes...')).toBeInTheDocument()
  })

  it('shows an empty state when the selected deck has no notes', async () => {
    const repository = {
      fetchDeckNotes: vi.fn().mockResolvedValue({ notes: [] }),
    } as unknown as AppRepository

    renderPanel(repository)

    expect(await screen.findByText('No notes in this deck yet.')).toBeInTheDocument()
  })

  it('renders the latest notes when the deck has existing notes', async () => {
    const repository = {
      fetchDeckNotes: vi.fn().mockResolvedValue({
        notes: [
          {
            noteId: 9,
            noteType: 'Basic',
            createdAt: '2026-03-15T12:30:00Z',
            modifiedAt: '2026-03-15T12:30:00Z',
            tags: ['aws'],
            fieldPreview: 'Newest recent note',
            cardCountInDeck: 1,
          },
        ],
      }),
    } as unknown as AppRepository

    renderPanel(repository)

    expect(await screen.findByText('Newest recent note')).toBeInTheDocument()
    expect(screen.getByText('#aws')).toBeInTheDocument()
  })
})
