import { describe, expect, it } from 'vitest'
import { fireEvent, render, screen } from '@testing-library/react'
import { AddNoteFormProvider } from '#/components/add-note-form-provider'
import { useAddNoteFormContext } from '#/components/add-note-form-context'

function AddNoteContextProbe() {
  const {
    selectedDeckId,
    selectedNoteType,
    fieldValues,
    tags,
    setSelectedDeckId,
    setSelectedNoteType,
    setFieldValues,
    setTags,
  } = useAddNoteFormContext()

  return (
    <div>
      <div data-testid="selected-deck-id">{selectedDeckId}</div>
      <div data-testid="selected-note-type">{selectedNoteType}</div>
      <div data-testid="field-front">{fieldValues.Front || ''}</div>
      <div data-testid="tags">{tags}</div>
      <button type="button" onClick={() => setSelectedDeckId(99)}>
        set deck
      </button>
      <button type="button" onClick={() => setSelectedNoteType('Basic')}>
        set note type
      </button>
      <button type="button" onClick={() => setFieldValues((prev) => ({ ...prev, Front: 'Question' }))}>
        set field
      </button>
      <button type="button" onClick={() => setTags('tag-a tag-b')}>
        set tags
      </button>
    </div>
  )
}

describe('AddNoteFormProvider', () => {
  it('initializes selected deck from prop and updates context state', () => {
    render(
      <AddNoteFormProvider deckId={42}>
        <AddNoteContextProbe />
      </AddNoteFormProvider>,
    )

    expect(screen.getByTestId('selected-deck-id')).toHaveTextContent('42')
    expect(screen.getByTestId('selected-note-type')).toHaveTextContent('')
    expect(screen.getByTestId('field-front')).toHaveTextContent('')
    expect(screen.getByTestId('tags')).toHaveTextContent('')

    fireEvent.click(screen.getByRole('button', { name: 'set deck' }))
    fireEvent.click(screen.getByRole('button', { name: 'set note type' }))
    fireEvent.click(screen.getByRole('button', { name: 'set field' }))
    fireEvent.click(screen.getByRole('button', { name: 'set tags' }))

    expect(screen.getByTestId('selected-deck-id')).toHaveTextContent('99')
    expect(screen.getByTestId('selected-note-type')).toHaveTextContent('Basic')
    expect(screen.getByTestId('field-front')).toHaveTextContent('Question')
    expect(screen.getByTestId('tags')).toHaveTextContent('tag-a tag-b')
  })

  it('throws when hook is used outside provider', () => {
    function OutsideProviderProbe() {
      useAddNoteFormContext()
      return null
    }

    expect(() => render(<OutsideProviderProbe />)).toThrow(
      'useAddNoteFormContext must be used within AddNoteFormProvider',
    )
  })
})
