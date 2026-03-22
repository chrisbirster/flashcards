import { useEffect, useMemo, useRef, useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import type { AICardSuggestion, NoteType } from '#/lib/api'
import { useAppRepository } from '#/lib/app-repository'
import { EmptyState, SurfaceCard } from '#/components/page-layout'

interface AICardSuggestionPanelProps {
  open: boolean
  noteType?: NoteType
  initialSourceText?: string
  existingFieldVals?: Record<string, string>
  onApplySuggestion: (suggestion: AICardSuggestion) => void
}

export function AICardSuggestionPanel({
  open,
  noteType,
  initialSourceText = '',
  existingFieldVals = {},
  onApplySuggestion,
}: AICardSuggestionPanelProps) {
  const repository = useAppRepository()
  const [sourceText, setSourceText] = useState(initialSourceText)
  const [hasSubmitted, setHasSubmitted] = useState(false)
  const wasOpenRef = useRef(open)

  useEffect(() => {
    const wasOpen = wasOpenRef.current
    if (open && !wasOpen) {
      setSourceText(initialSourceText)
      setHasSubmitted(false)
    }
    wasOpenRef.current = open
  }, [initialSourceText, open])

  const fieldList = useMemo(() => noteType?.fields ?? [], [noteType])

  const generateMutation = useMutation({
    mutationFn: () =>
      repository.generateAICardSuggestions({
        sourceText,
        noteType: noteType?.name ?? '',
        existingFieldVals,
        maxSuggestions: 3,
      }),
  })

  const suggestions = generateMutation.data?.suggestions ?? []

  return (
    <div className="space-y-4" data-testid="ai-suggestion-panel">
      <div className="space-y-2">
        <p className="text-sm font-semibold text-[var(--app-text)]">AI note suggestions</p>
        <p className="text-sm leading-6 text-[var(--app-text-soft)]">
          Paste rough notes, lecture bullets, or source material. We&apos;ll suggest {noteType?.name ?? 'the selected'} note
          field values, and nothing is saved until you apply a suggestion and review it.
        </p>
      </div>

      <label className="block space-y-2">
        <span className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">Source material</span>
        <textarea
          value={sourceText}
          onChange={(event) => setSourceText(event.target.value)}
          rows={7}
          placeholder="Paste notes, rough study bullets, or source material..."
          className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 py-3 text-sm text-[var(--app-text)] outline-none focus:border-[var(--app-accent)]"
        />
      </label>

      {noteType ? (
        <div className="flex flex-wrap gap-2">
          {fieldList.map((field) => (
            <span
              key={field}
              className="rounded-full border border-[var(--app-line)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-medium text-[var(--app-text-soft)]"
            >
              {field}
            </span>
          ))}
        </div>
      ) : (
        <div className="rounded-2xl border border-[var(--app-warning-line)] bg-[var(--app-warning-surface)] px-4 py-3 text-sm text-[var(--app-warning-text)]">
          Choose a note type first so the suggestions match the required fields.
        </div>
      )}

      <div className="flex flex-wrap items-center gap-3">
        <button
          type="button"
          onClick={() => {
            setHasSubmitted(true)
            generateMutation.mutate()
          }}
          disabled={generateMutation.isPending || !noteType || sourceText.trim().length === 0}
          className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)] disabled:cursor-not-allowed disabled:opacity-60"
        >
          {generateMutation.isPending ? 'Generating...' : 'Generate suggestions'}
        </button>
        {generateMutation.data?.provider ? (
          <span className="text-xs text-[var(--app-muted)]">
            Provider: {generateMutation.data.provider === 'dev' ? 'local preview mode' : `${generateMutation.data.provider}${generateMutation.data.model ? ` • ${generateMutation.data.model}` : ''}`}
          </span>
        ) : null}
      </div>

      {generateMutation.isError ? (
        <div className="rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 py-3 text-sm text-[var(--app-danger-text)]">
          {generateMutation.error instanceof Error ? generateMutation.error.message : 'AI suggestions failed.'}
        </div>
      ) : null}

      {hasSubmitted && !generateMutation.isPending && suggestions.length === 0 ? (
        <EmptyState
          title="No suggestions yet"
          description="Try a clearer prompt, paste a more specific excerpt, or switch to a note type that fits the source material."
        />
      ) : null}

      {suggestions.length > 0 ? (
        <div className="space-y-3">
          {suggestions.map((suggestion, index) => (
            <SurfaceCard key={`${suggestion.title}-${index}`} className="space-y-4">
              <div className="space-y-2">
                <div className="flex flex-wrap items-center justify-between gap-3">
                  <p className="text-base font-semibold text-[var(--app-text)]">{suggestion.title}</p>
                  <button
                    type="button"
                    onClick={() => onApplySuggestion(suggestion)}
                    className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card-strong)] px-4 text-sm font-medium text-[var(--app-text)]"
                  >
                    Use suggestion
                  </button>
                </div>
                {suggestion.rationale ? (
                  <p className="text-sm leading-6 text-[var(--app-text-soft)]">{suggestion.rationale}</p>
                ) : null}
              </div>
              <div className="space-y-3">
                {fieldList.map((field) => (
                  <div key={field} className="space-y-1">
                    <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">{field}</p>
                    <div className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card-strong)] px-4 py-3 text-sm leading-6 text-[var(--app-text)]">
                      {suggestion.fieldVals[field] || <span className="text-[var(--app-muted)]">Empty</span>}
                    </div>
                  </div>
                ))}
              </div>
            </SurfaceCard>
          ))}
        </div>
      ) : null}
    </div>
  )
}
