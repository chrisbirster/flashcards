import { describe, expect, it, vi } from 'vitest'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { AICardSuggestionPanel } from '#/components/ai-suggestion-panel'
import { AppRepositoryProvider, type AppRepository } from '#/lib/app-repository'

function renderPanel(repository: AppRepository, onApplySuggestion = vi.fn()) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  render(
    <QueryClientProvider client={queryClient}>
      <AppRepositoryProvider repository={repository}>
        <AICardSuggestionPanel
          open
          noteType={{
            name: 'Basic',
            fields: ['Front', 'Back'],
            templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
            sortFieldIndex: 0,
            fieldOptions: {},
          }}
          initialSourceText=""
          existingFieldVals={{}}
          onApplySuggestion={onApplySuggestion}
        />
      </AppRepositoryProvider>
    </QueryClientProvider>,
  )

  return { onApplySuggestion }
}

describe('AICardSuggestionPanel', () => {
  it('generates suggestions and applies a selected result', async () => {
    const repository = {
      generateAICardSuggestions: vi.fn().mockResolvedValue({
        provider: 'dev',
        model: 'heuristic',
        suggestions: [
          {
            title: 'Mitochondria',
            rationale: 'Generated from the pasted notes.',
            fieldVals: {
              Front: 'Mitochondria',
              Back: 'The powerhouse of the cell',
            },
          },
        ],
      }),
    } as unknown as AppRepository

    const { onApplySuggestion } = renderPanel(repository)

    fireEvent.change(screen.getByPlaceholderText('Paste notes, rough study bullets, or source material...'), {
      target: { value: 'Mitochondria: The powerhouse of the cell' },
    })
    fireEvent.click(screen.getByRole('button', { name: 'Generate suggestions' }))

    await waitFor(() => expect(repository.generateAICardSuggestions).toHaveBeenCalled())
    expect(await screen.findByRole('button', { name: 'Use suggestion' })).toBeInTheDocument()

    fireEvent.click(screen.getByRole('button', { name: 'Use suggestion' }))

    expect(onApplySuggestion).toHaveBeenCalledWith({
      title: 'Mitochondria',
      rationale: 'Generated from the pasted notes.',
      fieldVals: {
        Front: 'Mitochondria',
        Back: 'The powerhouse of the cell',
      },
    })
  })
})
