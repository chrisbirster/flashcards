import type { ReactNode } from 'react'

export function ActionBar({
  children,
  className = '',
}: {
  children: ReactNode
  className?: string
}) {
  return (
    <div
      className={[
        'sticky bottom-0 z-10 border-t border-[var(--app-line)] bg-[color:var(--app-header)]/95 px-4 py-3 backdrop-blur',
        className,
      ].join(' ')}
      style={{ paddingBottom: 'calc(0.75rem + env(safe-area-inset-bottom))' }}
    >
      {children}
    </div>
  )
}

export function FormActions({
  children,
  className = '',
}: {
  children: ReactNode
  className?: string
}) {
  return (
    <div className={['flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-end', className].join(' ')}>
      {children}
    </div>
  )
}
