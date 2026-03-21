import { fireEvent, render, screen, within } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { describe, expect, it, vi } from 'vitest'
import { TemplateEditor } from '#/components/TemplateEditor'
import { AppRepositoryProvider, type AppRepository } from '#/lib/app-repository'

function renderTemplateEditor(repository: AppRepository) {
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
        <TemplateEditor
          noteType={{
            name: 'Basic',
            fields: ['Front', 'Back'],
            templates: [
              {
                name: 'Card 1',
                qFmt: '<div class="front">{{Front}}</div>',
                aFmt: '<div class="back">{{FrontSide}}<hr id="answer">{{Back}}</div>',
                styling: '.card { background: black; color: white; padding: 16px; }',
                isCloze: false,
              },
            ],
            sortFieldIndex: 0,
            fieldOptions: {},
          }}
          onClose={() => {}}
        />
      </AppRepositoryProvider>
    </QueryClientProvider>,
  )
}

describe('TemplateEditor styling preview', () => {
  it('renders styled front and back previews instead of echoing raw CSS', async () => {
    const repository = {
      fetchDecks: vi.fn().mockResolvedValue([]),
    } as unknown as AppRepository

    renderTemplateEditor(repository)

    fireEvent.click(screen.getByTestId('tab-styling'))

    const preview = await screen.findByTestId('styling-preview')
    expect(within(preview).getByText('Front card')).toBeInTheDocument()
    expect(within(preview).getByText('Back card')).toBeInTheDocument()
    expect(screen.getByTestId('preview-content-front')).toHaveTextContent('[Front sample 1]')
    expect(screen.getByTestId('preview-content-back')).toHaveTextContent('[Front sample 1]')
    expect(screen.getByTestId('preview-content-back')).toHaveTextContent('[Back sample 2]')
  })
})
