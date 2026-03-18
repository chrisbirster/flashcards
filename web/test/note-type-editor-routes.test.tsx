import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest'
import { cleanup, fireEvent, render, screen } from '@testing-library/react'
import type { ReactElement } from 'react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import {
  AddNoteFieldEditorRoutePage,
  AddNoteTemplateEditorRoutePage,
  TemplatesTemplateEditorRoutePage,
} from '#/pages/NoteTypeEditorRoutes'

const fetchNoteTypesMock = vi.fn()

vi.mock('#/lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('#/lib/api')>()
  return {
    ...actual,
    fetchNoteTypes: () => fetchNoteTypesMock(),
  }
})

vi.mock('#/components/FieldEditor', () => ({
  FieldEditor: ({ noteType }: { noteType: { name: string } }) => (
    <div data-testid="field-editor">FieldEditor:{noteType.name}</div>
  ),
}))

vi.mock('#/components/TemplateEditor', () => ({
  TemplateEditor: ({ noteType }: { noteType: { name: string } }) => (
    <div data-testid="template-editor">TemplateEditor:{noteType.name}</div>
  ),
}))

function renderRoute(path: string, element: ReactElement) {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        retry: false,
      },
    },
  })

  return render(
    <QueryClientProvider client={queryClient}>
      <MemoryRouter initialEntries={[path]}>
        <Routes>
          <Route path="/notes/add/note-types/:noteTypeName/fields" element={element} />
          <Route path="/notes/add/note-types/:noteTypeName/templates" element={element} />
          <Route path="/templates/note-types/:noteTypeName/templates" element={element} />
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

describe('NoteType editor routes', () => {
  beforeEach(() => {
    fetchNoteTypesMock.mockReset()
  })

  afterEach(() => {
    cleanup()
  })

  it('renders field editor for note type route', async () => {
    fetchNoteTypesMock.mockResolvedValue([
      {
        name: 'Basic',
        fields: ['Front', 'Back'],
        templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
        sortFieldIndex: 0,
      },
    ])

    renderRoute('/notes/add/note-types/Basic/fields', <AddNoteFieldEditorRoutePage />)

    expect(await screen.findByTestId('field-editor')).toHaveTextContent('FieldEditor:Basic')
  })

  it('renders template editor for template route', async () => {
    fetchNoteTypesMock.mockResolvedValue([
      {
        name: 'Basic',
        fields: ['Front', 'Back'],
        templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
        sortFieldIndex: 0,
      },
    ])

    renderRoute('/notes/add/note-types/Basic/templates', <AddNoteTemplateEditorRoutePage />)

    expect(await screen.findByTestId('template-editor')).toHaveTextContent('TemplateEditor:Basic')
  })

  it('shows not-found state when note type does not exist', async () => {
    fetchNoteTypesMock.mockResolvedValue([
      {
        name: 'Basic',
        fields: ['Front', 'Back'],
        templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
        sortFieldIndex: 0,
      },
    ])

    renderRoute('/notes/add/note-types/UnknownType/fields', <AddNoteFieldEditorRoutePage />)

    expect(await screen.findByText(/was not found/i)).toHaveTextContent('UnknownType')
  })

  it('shows loading state while note types are fetching', async () => {
    fetchNoteTypesMock.mockReturnValue(new Promise(() => {}))

    renderRoute('/notes/add/note-types/Basic/fields', <AddNoteFieldEditorRoutePage />)

    expect(screen.getByText('Loading note type...')).toBeInTheDocument()
  })

  it('shows error state for field route when query fails', async () => {
    fetchNoteTypesMock.mockRejectedValue(new Error('boom'))

    renderRoute('/notes/add/note-types/Basic/fields', <AddNoteFieldEditorRoutePage />)

    expect(await screen.findByText('Failed to load note type.')).toBeInTheDocument()
  })

  it('shows not-found state for template route when note type is missing', async () => {
    fetchNoteTypesMock.mockResolvedValue([
      {
        name: 'Basic',
        fields: ['Front', 'Back'],
        templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
        sortFieldIndex: 0,
      },
    ])

    renderRoute('/notes/add/note-types/UnknownTemplateType/templates', <AddNoteTemplateEditorRoutePage />)

    expect(await screen.findByText(/was not found/i)).toHaveTextContent('UnknownTemplateType')
  })

  it('shows error state for template route when query fails', async () => {
    fetchNoteTypesMock.mockRejectedValue(new Error('boom-template'))

    renderRoute('/notes/add/note-types/Basic/templates', <AddNoteTemplateEditorRoutePage />)

    expect(await screen.findByText('Failed to load note type.')).toBeInTheDocument()
  })

  it('falls back to raw route param when note type decode fails', async () => {
    fetchNoteTypesMock.mockResolvedValue([])

    renderRoute('/notes/add/note-types/%E0%A4%A/fields', <AddNoteFieldEditorRoutePage />)

    expect(await screen.findByText(/was not found/i)).toHaveTextContent('%E0%A4%A')
  })

  it('renders templates-page route variant', async () => {
    fetchNoteTypesMock.mockResolvedValue([
      {
        name: 'Basic',
        fields: ['Front', 'Back'],
        templates: [{ name: 'Card 1', qFmt: '{{Front}}', aFmt: '{{Back}}', styling: '', isCloze: false }],
        sortFieldIndex: 0,
      },
    ])

    renderRoute('/templates/note-types/Basic/templates', <TemplatesTemplateEditorRoutePage />)

    expect(await screen.findByTestId('template-editor')).toHaveTextContent('TemplateEditor:Basic')
  })

  it('handles close action with fallback navigation branch', async () => {
    const historyLengthSpy = vi.spyOn(window.history, 'length', 'get').mockReturnValue(1)
    fetchNoteTypesMock.mockResolvedValue([])

    renderRoute('/notes/add/note-types/Nope/fields', <AddNoteFieldEditorRoutePage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Close' }))

    historyLengthSpy.mockRestore()
  })

  it('handles close action with back-navigation branch', async () => {
    const historyLengthSpy = vi.spyOn(window.history, 'length', 'get').mockReturnValue(2)
    fetchNoteTypesMock.mockResolvedValue([])

    renderRoute('/notes/add/note-types/Nope/templates', <AddNoteTemplateEditorRoutePage />)
    fireEvent.click(await screen.findByRole('button', { name: 'Close' }))

    historyLengthSpy.mockRestore()
  })
})
