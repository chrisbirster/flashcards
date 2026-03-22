import { type ReactNode, useEffect } from 'react'

function useEscapeToClose(open: boolean, onClose: () => void) {
  useEffect(() => {
    if (!open) return
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [open, onClose])
}

function Overlay({ onClose }: { onClose: () => void }) {
  return (
    <button
      type="button"
      onClick={onClose}
      className="absolute inset-0 bg-black/60"
      aria-label="Close"
    />
  )
}

export function Sheet({
  open,
  onClose,
  title,
  children,
}: {
  open: boolean
  onClose: () => void
  title?: string
  children: ReactNode
}) {
  useEscapeToClose(open, onClose)
  if (!open) return null

  return (
    <div className="fixed inset-0 z-50">
      <Overlay onClose={onClose} />
      <div
        className="absolute inset-x-0 bottom-0 rounded-t-[1.75rem] border border-[var(--app-line)] bg-[var(--app-panel)] shadow-2xl md:hidden"
        style={{ paddingBottom: 'env(safe-area-inset-bottom)' }}
      >
        <div className="mx-auto mt-3 h-1.5 w-14 rounded-full bg-[var(--app-line-strong)]" />
        {title ? <div className="px-5 py-4 text-sm font-semibold text-[var(--app-text)]">{title}</div> : null}
        <div className="max-h-[78dvh] overflow-y-auto px-5 pb-5">{children}</div>
      </div>
      <div className="absolute inset-0 hidden items-center justify-center p-6 md:flex">
        <div className="w-full max-w-lg rounded-[1.75rem] border border-[var(--app-line)] bg-[var(--app-panel)] shadow-2xl">
          {title ? (
            <div className="border-b border-[var(--app-line)] px-6 py-5 text-base font-semibold text-[var(--app-text)]">
              {title}
            </div>
          ) : null}
          <div className="max-h-[min(80vh,48rem)] overflow-y-auto px-6 py-6">{children}</div>
        </div>
      </div>
    </div>
  )
}

export function FullscreenSheet({
  open,
  onClose,
  title,
  children,
}: {
  open: boolean
  onClose: () => void
  title: string
  children: ReactNode
}) {
  useEscapeToClose(open, onClose)
  if (!open) return null

  return (
    <div className="fixed inset-0 z-50 xl:hidden">
      <div className="absolute inset-0 bg-[var(--app-bg)]" />
      <div className="relative flex min-h-[100dvh] flex-col">
        <div
          className="sticky top-0 z-10 flex items-center justify-between border-b border-[var(--app-line)] bg-[color:var(--app-header)] px-4 py-4 backdrop-blur"
          style={{ paddingTop: 'calc(1rem + env(safe-area-inset-top))' }}
        >
          <div>
            <p className="text-xs uppercase tracking-[0.24em] text-[var(--app-muted)]">Editor</p>
            <h2 className="text-lg font-semibold text-[var(--app-text)]">{title}</h2>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="inline-flex min-h-11 min-w-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-3 text-sm font-medium text-[var(--app-text-soft)]"
          >
            Close
          </button>
        </div>
        <div className="flex-1 overflow-y-auto px-4 py-4">{children}</div>
      </div>
    </div>
  )
}

export function ConfirmSheet({
  open,
  onClose,
  title,
  description,
  confirmLabel,
  cancelLabel = 'Cancel',
  onConfirm,
  destructive = false,
}: {
  open: boolean
  onClose: () => void
  title: string
  description: string
  confirmLabel: string
  cancelLabel?: string
  onConfirm: () => void
  destructive?: boolean
}) {
  return (
    <Sheet open={open} onClose={onClose} title={title}>
      <p className="text-sm leading-6 text-[var(--app-text-soft)]">{description}</p>
      <div className="mt-5 flex flex-col gap-3">
        <button
          type="button"
          onClick={() => {
            onConfirm()
            onClose()
          }}
          className={[
            'inline-flex min-h-11 items-center justify-center rounded-2xl px-4 text-sm font-semibold',
            destructive
              ? 'bg-[color:var(--app-danger-text)] text-[var(--app-accent-ink)]'
              : 'bg-[var(--app-accent)] text-[var(--app-accent-ink)]',
          ].join(' ')}
        >
          {confirmLabel}
        </button>
        <button
          type="button"
          onClick={onClose}
          className="inline-flex min-h-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
        >
          {cancelLabel}
        </button>
      </div>
    </Sheet>
  )
}
