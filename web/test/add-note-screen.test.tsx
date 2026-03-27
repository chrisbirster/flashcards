import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { MemoryRouter } from 'react-router'
import { AddNoteScreen } from '#/components/AddNoteScreen'
import { AppRepositoryProvider, type AppRepository } from '#/lib/app-repository'

vi.mock('#/components/ai-suggestion-panel', () => ({
  AICardSuggestionPanel: ({
    open,
    onApplySuggestion,
  }: {
    open: boolean
    onApplySuggestion: (suggestion: {title: string; rationale: string; fieldVals: Record<string, string>}) => void
  }) =>
    open ? (
      <div data-testid="ai-suggestion-panel">
        <button
          type="button"
          onClick={() =>
            onApplySuggestion({
              title: 'Mitochondria',
              rationale: 'Mocked AI suggestion',
              fieldVals: {
                Front: 'Mitochondria',
                Back: 'The powerhouse of the cell',
              },
            })
          }
        >
          Apply mocked AI suggestion
        </button>
      </div>
    ) : null,
}))

afterEach(() => {
  cleanup()
})

beforeEach(() => {
  vi.useRealTimers()
})

function renderAddNoteScreen(repository: AppRepository) {
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
        <MemoryRouter>
          <AddNoteScreen deckId={1} onClose={() => {}} />
        </MemoryRouter>
      </AppRepositoryProvider>
    </QueryClientProvider>,
  )
}

describe('AddNoteScreen recent notes', () => {
  it('prepends the newly created note to the recent deck notes panel', async () => {
    const repository = {
      fetchNoteTypes: vi.fn().mockResolvedValue([
        {
          name: 'Basic',
          fields: ['Front', 'Back'],
          templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
          sortFieldIndex: 0,
          fieldOptions: {},
        },
      ]),
      fetchDecks: vi.fn().mockResolvedValue([{ id: 1, name: 'Deck 1', parentId: null, cardIds: [] }]),
      fetchDeckNotes: vi.fn().mockResolvedValue({
        notes: [
          {
            noteId: 1,
            noteType: 'Basic',
            createdAt: '2026-03-15T11:00:00Z',
            modifiedAt: '2026-03-15T11:00:00Z',
            tags: ['existing'],
            fieldPreview: 'Existing recent note',
            cardCountInDeck: 1,
          },
        ],
      }),
      checkDuplicate: vi.fn().mockResolvedValue({ isDuplicate: false, duplicates: [] }),
      createNote: vi.fn().mockResolvedValue({
        note: {
          id: 2,
          typeId: 'Basic',
          fieldVals: {
            Front: 'Newest front',
            Back: 'Newest back',
          },
          tags: ['new'],
          createdAt: '2026-03-15T12:00:00Z',
          modifiedAt: '2026-03-15T12:00:00Z',
        },
        cards: [
          {
            id: 20,
            deckId: 1,
          },
        ],
      }),
    } as unknown as AppRepository

    renderAddNoteScreen(repository)

    expect(await screen.findByText('Existing recent note')).toBeInTheDocument()

    fireEvent.change(await screen.findByPlaceholderText('Enter front...'), { target: { value: 'Newest front' } })
    fireEvent.change(screen.getByPlaceholderText('Enter back...'), { target: { value: 'Newest back' } })
    fireEvent.click(screen.getByRole('button', { name: 'Add Note' }))

    await waitFor(() => expect(repository.createNote).toHaveBeenCalled())

    const items = within(screen.getByTestId('recent-deck-notes')).getAllByRole('listitem')
    expect(items[0]).toHaveTextContent('Newest front')
    expect(items[1]).toHaveTextContent('Existing recent note')
  })

  it('applies an AI suggestion into the current note form before save', async () => {
    const originalMatchMedia = window.matchMedia
    window.matchMedia = vi.fn().mockImplementation((query: string) => ({
      matches: query === '(min-width: 768px)',
      media: query,
      onchange: null,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
      addListener: vi.fn(),
      removeListener: vi.fn(),
      dispatchEvent: vi.fn(),
    }))

    const repository = {
      fetchNoteTypes: vi.fn().mockResolvedValue([
        {
          name: 'Basic',
          fields: ['Front', 'Back'],
          templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
          sortFieldIndex: 0,
          fieldOptions: {},
        },
      ]),
      fetchDecks: vi.fn().mockResolvedValue([{ id: 1, name: 'Deck 1', parentId: null, cardIds: [] }]),
      fetchDeckNotes: vi.fn().mockResolvedValue({ notes: [] }),
      checkDuplicate: vi.fn().mockResolvedValue({ isDuplicate: false, duplicates: [] }),
      createNote: vi.fn(),
    } as unknown as AppRepository

    try {
      renderAddNoteScreen(repository)

      fireEvent.click(await screen.findByRole('button', { name: 'AI suggestions' }))
      fireEvent.click(await screen.findByRole('button', { name: 'Apply mocked AI suggestion' }))

      await waitFor(() => {
        expect(screen.getAllByPlaceholderText('Enter front...').some((field) => (field as HTMLTextAreaElement).value === 'Mitochondria')).toBe(true)
        expect(
          screen
            .getAllByPlaceholderText('Enter back...')
            .some((field) => (field as HTMLTextAreaElement).value === 'The powerhouse of the cell'),
        ).toBe(true)
      })
    } finally {
      window.matchMedia = originalMatchMedia
    }
  })

  it('only checks duplicates when the primary field changes', async () => {
    const checkDuplicate = vi.fn().mockResolvedValue({ isDuplicate: false, duplicates: [] })
    const repository = {
      fetchNoteTypes: vi.fn().mockResolvedValue([
        {
          name: 'Basic',
          fields: ['Front', 'Back'],
          templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
          sortFieldIndex: 0,
          fieldOptions: {},
        },
      ]),
      fetchDecks: vi.fn().mockResolvedValue([{ id: 1, name: 'Deck 1', parentId: null, cardIds: [] }]),
      fetchDeckNotes: vi.fn().mockResolvedValue({ notes: [] }),
      checkDuplicate,
      createNote: vi.fn(),
    } as unknown as AppRepository

    renderAddNoteScreen(repository)

    fireEvent.change(await screen.findByPlaceholderText('Enter front...'), { target: { value: 'AWS region' } })

    await waitFor(() => expect(checkDuplicate).toHaveBeenCalledTimes(1), { timeout: 1200 })

    fireEvent.change(screen.getByPlaceholderText('Enter back...'), { target: { value: 'A group of data centers' } })

    await new Promise((resolve) => setTimeout(resolve, 700))

    expect(checkDuplicate).toHaveBeenCalledTimes(1)
  })
})
