import { useMemo, useRef, useState, type ReactNode } from 'react'
import { AddNoteFormContext, type AddNoteFormContextValue } from './add-note-form-context'

export function AddNoteFormProvider({
  deckId,
  children,
}: {
  deckId?: number
  children: ReactNode
}) {
  const [selectedNoteType, setSelectedNoteType] = useState('')
  const [selectedDeckId, setSelectedDeckId] = useState<number>(deckId || 0)
  const [fieldValues, setFieldValues] = useState<Record<string, string>>({})
  const [tags, setTags] = useState('')
  const [activeField, setActiveField] = useState<string | null>(null)
  const textareaRefs = useRef<Record<string, HTMLTextAreaElement | null>>({})

  const value = useMemo<AddNoteFormContextValue>(
    () => ({
      selectedNoteType,
      setSelectedNoteType,
      selectedDeckId,
      setSelectedDeckId,
      fieldValues,
      setFieldValues,
      tags,
      setTags,
      activeField,
      setActiveField,
      textareaRefs,
    }),
    [selectedNoteType, selectedDeckId, fieldValues, tags, activeField],
  )

  return <AddNoteFormContext.Provider value={value}>{children}</AddNoteFormContext.Provider>
}
