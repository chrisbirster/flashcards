import type { ReactNode } from 'react'

export function FieldRow({
  label,
  hint,
  children,
}: {
  label: string
  hint?: string
  children: ReactNode
}) {
  return (
    <div className="space-y-2">
      <div className="flex flex-col gap-1 sm:flex-row sm:items-end sm:justify-between">
        <label className="text-sm font-medium text-[var(--app-text)]">{label}</label>
        {hint ? <span className="text-xs text-[var(--app-muted)]">{hint}</span> : null}
      </div>
      {children}
    </div>
  )
}
