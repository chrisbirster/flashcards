import type { CSSProperties, ReactNode } from 'react'

export function PageContainer({
  children,
  className = '',
}: {
  children: ReactNode
  className?: string
}) {
  return <div className={`mx-auto w-full max-w-7xl ${className}`.trim()}>{children}</div>
}

export function PageSection({
  children,
  className = '',
}: {
  children: ReactNode
  className?: string
}) {
  return (
    <section
      className={[
        'rounded-[1.75rem] border border-[var(--app-line)] bg-[var(--app-card)] shadow-sm',
        className,
      ].join(' ')}
    >
      {children}
    </section>
  )
}

export function SurfaceCard({
  children,
  className = '',
}: {
  children: ReactNode
  className?: string
}) {
  return (
    <div
      className={[
        'rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4 shadow-sm sm:p-5',
        className,
      ].join(' ')}
    >
      {children}
    </div>
  )
}

export function StatCard({
  label,
  value,
  detail,
  accent,
}: {
  label: string
  value: string | number
  detail: string
  accent?: CSSProperties
}) {
  return (
    <SurfaceCard className="min-h-[10rem]">
      <p className="text-[11px] uppercase tracking-[0.24em] text-[var(--app-muted)]">{label}</p>
      <p className="mt-4 text-3xl font-semibold tracking-tight text-[var(--app-text)]" style={accent}>
        {value}
      </p>
      <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">{detail}</p>
    </SurfaceCard>
  )
}

export function EmptyState({
  title,
  description,
  action,
}: {
  title: string
  description: string
  action?: ReactNode
}) {
  return (
    <SurfaceCard className="px-5 py-10 text-center">
      <p className="text-lg font-semibold text-[var(--app-text)]">{title}</p>
      <p className="mx-auto mt-3 max-w-xl text-sm leading-6 text-[var(--app-text-soft)]">{description}</p>
      {action ? <div className="mt-6 flex justify-center">{action}</div> : null}
    </SurfaceCard>
  )
}
