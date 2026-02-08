import { createContext, useContext } from 'react'
import type { Dispatch, MutableRefObject, SetStateAction } from 'react'

export interface AddNoteFormContextValue {
  selectedNoteType: string
  setSelectedNoteType: Dispatch<SetStateAction<string>>
  selectedDeckId: number
  setSelectedDeckId: Dispatch<SetStateAction<number>>
  fieldValues: Record<string, string>
  setFieldValues: Dispatch<SetStateAction<Record<string, string>>>
  tags: string
  setTags: Dispatch<SetStateAction<string>>
  activeField: string | null
  setActiveField: Dispatch<SetStateAction<string | null>>
  textareaRefs: MutableRefObject<Record<string, HTMLTextAreaElement | null>>
}

export const AddNoteFormContext = createContext<AddNoteFormContextValue | null>(null)

export function useAddNoteFormContext() {
  const context = useContext(AddNoteFormContext)

  if (!context) {
    throw new Error('useAddNoteFormContext must be used within AddNoteFormProvider')
  }

  return context
}
