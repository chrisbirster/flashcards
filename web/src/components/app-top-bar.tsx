import type { ReactNode } from 'react'

export function AppTopBar({
  title,
  subtitle = 'Workspace',
  leading,
  trailing,
}: {
  title: string
  subtitle?: string
  leading?: ReactNode
  trailing?: ReactNode
}) {
  return (
    <header className="sticky top-0 z-30 border-b border-[var(--app-line)] bg-[color:var(--app-header)] backdrop-blur">
      <div
        className="flex items-center justify-between gap-3 px-4 py-4 md:px-8"
        style={{ paddingTop: 'calc(1rem + env(safe-area-inset-top))' }}
      >
        <div className="flex min-w-0 items-center gap-3">
          {leading}
          <div className="min-w-0">
            <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">{subtitle}</p>
            <h1 className="truncate text-xl font-semibold tracking-tight text-[var(--app-text)] sm:text-2xl">{title}</h1>
          </div>
        </div>
        {trailing}
      </div>
    </header>
  )
}
