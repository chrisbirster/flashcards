import { useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate, useParams } from 'react-router'
import { FieldEditor } from '#/components/FieldEditor'
import { TemplateEditor } from '#/components/TemplateEditor'
import { useAppRepository } from '#/lib/app-repository'

interface EditorRouteStateProps {
  message: string
  onClose: () => void
}

function EditorRouteState({ message, onClose }: EditorRouteStateProps) {
  return (
    <div className="fixed inset-0 z-50 flex items-end justify-center bg-black/60 p-2 sm:items-center sm:p-0">
      <div className="max-h-[95dvh] w-full max-w-md rounded-t-[1.75rem] border border-[var(--app-line)] bg-[var(--app-panel)] p-4 shadow-xl sm:mx-4 sm:rounded-[1.75rem] sm:p-6">
        <p className="text-sm leading-6 text-[var(--app-text-soft)]">{message}</p>
        <div className="mt-4 flex justify-end">
          <button
            type="button"
            onClick={onClose}
            className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
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
  const repository = useAppRepository()

  const { data: noteTypes, isLoading, error } = useQuery({
    queryKey: ['note-types'],
    queryFn: () => repository.fetchNoteTypes(),
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
