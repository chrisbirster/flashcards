import { describe, expect, it, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import type { ReactElement } from 'react'
import { MemoryRouter, Route, Routes } from 'react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import {
  AddNoteFieldEditorRoutePage,
  AddNoteTemplateEditorRoutePage,
} from '#/pages/NoteTypeEditorRoutes'

const fetchNoteTypesMock = vi.fn()

vi.mock('#/lib/api', () => ({
  fetchNoteTypes: () => fetchNoteTypesMock(),
}))

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
        </Routes>
      </MemoryRouter>
    </QueryClientProvider>,
  )
}

describe('NoteType editor routes', () => {
  beforeEach(() => {
    fetchNoteTypesMock.mockReset()
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
})
