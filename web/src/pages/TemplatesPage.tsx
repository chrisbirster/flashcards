import { Outlet, useNavigate } from 'react-router'
import { useQuery } from '@tanstack/react-query'
import type { NoteType } from '#/lib/api'
import { useAppRepository } from '#/lib/app-repository'
import { EmptyState, PageContainer, PageSection, SurfaceCard } from '#/components/page-layout'

export function TemplatesPage() {
  const navigate = useNavigate()
  const repository = useAppRepository()

  const { data: noteTypes, isLoading, error } = useQuery({
    queryKey: ['note-types'],
    queryFn: () => repository.fetchNoteTypes(),
  })

  const handleEditTemplate = (noteType: NoteType) => {
    navigate(encodeURIComponent(noteType.name))
  }

  if (isLoading) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-text-soft)]">
          Loading template workspace...
        </PageSection>
      </PageContainer>
    )
  }

  if (error) {
    return (
      <PageContainer>
        <PageSection className="px-5 py-16 text-center text-sm text-[var(--app-danger-text)]">
          Error: {error instanceof Error ? error.message : 'Failed to load note types'}
        </PageSection>
      </PageContainer>
    )
  }

  return (
    <PageContainer className="space-y-4">
      <PageSection className="p-4 sm:p-5">
        <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">Templates</p>
        <h1 className="mt-2 text-2xl font-semibold tracking-tight text-[var(--app-text)]">Card Templates</h1>
        <p className="mt-2 max-w-3xl text-sm leading-6 text-[var(--app-text-soft)]">
          Manage note-type templates, tune conditional generation, and preview how cards render on both mobile and desktop.
        </p>
      </PageSection>

      {!noteTypes || noteTypes.length === 0 ? (
        <EmptyState
          title="No note types found"
          description="Create a note type first, then return here to edit fields and card templates."
        />
      ) : (
        <div className="grid gap-4">
          {noteTypes.map((noteType) => (
            <PageSection key={noteType.name} className="p-4 sm:p-5">
              <div className="flex flex-col gap-4 md:flex-row md:items-start md:justify-between">
                <div className="min-w-0">
                  <h2 className="text-lg font-semibold text-[var(--app-text)]">{noteType.name}</h2>
                  <p className="mt-2 text-sm text-[var(--app-text-soft)]">
                    {noteType.fields.length} field{noteType.fields.length !== 1 ? 's' : ''},{' '}
                    {noteType.templates.length} template{noteType.templates.length !== 1 ? 's' : ''}
                  </p>
                </div>
                <button
                  type="button"
                  onClick={() => handleEditTemplate(noteType)}
                  className="inline-flex min-h-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] px-4 text-sm font-semibold text-[var(--app-accent-ink)]"
                >
                  Edit templates
                </button>
              </div>

              <div className="mt-5 grid gap-4 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)]">
                <SurfaceCard className="space-y-3 border-none bg-[var(--app-card-strong)] p-4">
                  <p className="text-[11px] uppercase tracking-[0.18em] text-[var(--app-muted)]">Fields</p>
                  <div className="flex flex-wrap gap-2">
                    {noteType.fields.map((field) => (
                      <span
                        key={field}
                        className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-3 py-1 text-xs font-medium text-[var(--app-text)]"
                      >
                        {field}
                      </span>
                    ))}
                  </div>
                </SurfaceCard>

                <div className="grid gap-3">
                  {noteType.templates.map((template) => (
                    <SurfaceCard key={template.name} className="space-y-2 border-none bg-[var(--app-card-strong)] p-4">
                      <div className="flex flex-wrap items-center justify-between gap-2">
                        <span className="text-sm font-semibold text-[var(--app-text)]">{template.name}</span>
                        <div className="flex flex-wrap items-center gap-2">
                          {template.isCloze ? (
                            <span className="rounded-full bg-[var(--app-accent)] px-2.5 py-1 text-[11px] font-semibold uppercase tracking-[0.14em] text-[var(--app-accent-ink)]">
                              Cloze
                            </span>
                          ) : null}
                          {template.deckOverride ? (
                            <span className="rounded-full border border-[var(--app-line)] bg-[var(--app-card)] px-2.5 py-1 text-[11px] text-[var(--app-text-soft)]">
                              {template.deckOverride}
                            </span>
                          ) : null}
                        </div>
                      </div>
                      {template.ifFieldNonEmpty ? (
                        <p className="text-xs leading-5 text-[var(--app-text-soft)]">
                          Conditional on field <span className="font-medium text-[var(--app-text)]">{template.ifFieldNonEmpty}</span>
                        </p>
                      ) : (
                        <p className="text-xs leading-5 text-[var(--app-text-soft)]">
                          Generated whenever the note has enough content for this template.
                        </p>
                      )}
                    </SurfaceCard>
                  ))}
                </div>
              </div>
            </PageSection>
          ))}
        </div>
      )}

      <Outlet />
    </PageContainer>
  )
}
