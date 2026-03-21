import { forwardRef, type JSX } from 'react'

const CloseSVG = () => {
  return (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
    </svg>
  )
}

export const EditFieldPanelCloseButton = forwardRef<HTMLButtonElement, JSX.IntrinsicElements['button']>(function EditFieldPanelCloseButton({
  onClick,
  ...props
}, ref) {
  return (
    <button
      onClick={onClick}
      className="rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] p-2 text-[var(--app-text-soft)]"
      {...props}
      ref={ref}
    >
      <CloseSVG />
    </button>
  )
})
