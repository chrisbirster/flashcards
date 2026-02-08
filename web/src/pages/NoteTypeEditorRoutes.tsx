import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router'
import { FieldEditor } from '#/components/FieldEditor'
import { TemplateEditor } from '#/components/TemplateEditor'
import { fetchNoteTypes } from '#/lib/api'

interface EditorRouteStateProps {
  message: string
  onClose: () => void
}

function EditorRouteState({ message, onClose }: EditorRouteStateProps) {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
        <p className="text-gray-700">{message}</p>
        <div className="mt-4 flex justify-end">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}

function decodeNoteTypeName(rawName?: string): string {
  if (!rawName) return ''

  try {
    return decodeURIComponent(rawName)
  } catch {
    return rawName
  }
}

function useCloseWithFallback(fallbackPath: string) {
  const navigate = useNavigate()

  return () => {
    if (window.history.length > 1) {
      navigate(-1)
      return
    }

    navigate(fallbackPath, { replace: true })
  }
}

function useNoteTypeFromRoute() {
  const { noteTypeName } = useParams<{ noteTypeName: string }>()
  const decodedNoteTypeName = useMemo(() => decodeNoteTypeName(noteTypeName), [noteTypeName])

  const { data: noteTypes, isLoading, error } = useQuery({
    queryKey: ['note-types'],
    queryFn: fetchNoteTypes,
  })

  const noteType = noteTypes?.find((candidate) => candidate.name === decodedNoteTypeName)

  return {
    decodedNoteTypeName,
    noteType,
    isLoading,
    error,
  }
}

function FieldEditorRoute({ fallbackPath }: { fallbackPath: string }) {
  const closeEditor = useCloseWithFallback(fallbackPath)
  const { decodedNoteTypeName, noteType, isLoading, error } = useNoteTypeFromRoute()

  if (isLoading) {
    return <EditorRouteState message="Loading note type..." onClose={closeEditor} />
  }

  if (error) {
    return <EditorRouteState message="Failed to load note type." onClose={closeEditor} />
  }

  if (!noteType) {
    return (
      <EditorRouteState
        message={`Note type "${decodedNoteTypeName || 'Unknown'}" was not found.`}
        onClose={closeEditor}
      />
    )
  }

  return <FieldEditor noteType={noteType} onClose={closeEditor} />
}

function TemplateEditorRoute({ fallbackPath }: { fallbackPath: string }) {
  const closeEditor = useCloseWithFallback(fallbackPath)
  const { decodedNoteTypeName, noteType, isLoading, error } = useNoteTypeFromRoute()

  if (isLoading) {
    return <EditorRouteState message="Loading note type..." onClose={closeEditor} />
  }

  if (error) {
    return <EditorRouteState message="Failed to load note type." onClose={closeEditor} />
  }

  if (!noteType) {
    return (
      <EditorRouteState
        message={`Note type "${decodedNoteTypeName || 'Unknown'}" was not found.`}
        onClose={closeEditor}
      />
    )
  }

  return <TemplateEditor noteType={noteType} onClose={closeEditor} />
}

export function AddNoteFieldEditorRoutePage() {
  return <FieldEditorRoute fallbackPath="/notes/add" />
}

export function AddNoteTemplateEditorRoutePage() {
  return <TemplateEditorRoute fallbackPath="/notes/add" />
}

export function TemplatesTemplateEditorRoutePage() {
  return <TemplateEditorRoute fallbackPath="/templates" />
}
