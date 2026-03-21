import { forwardRef, type JSX } from "react";

type EditFieldPopupProps = JSX.IntrinsicElements['div'] & {
    field: string;
}

export const EditFieldPopup = forwardRef<HTMLDivElement, EditFieldPopupProps>(function EditFieldPopup({
    field,
    children,
    ...props
}, ref) {
    return (
        <div className="ml-0 rounded-[1.25rem] border border-[var(--app-line)] bg-[var(--app-card)] p-3 sm:ml-4" data-testid={`field-options-panel-${field}`}
            {...props}
            ref={ref}
        >
            <div className="mb-2 text-xs font-medium uppercase tracking-[0.16em] text-[var(--app-muted)]">Field Options: {field}</div>
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                {children}
            </div>
        </div>
    )
})
