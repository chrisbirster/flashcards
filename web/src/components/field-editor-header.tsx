import { forwardRef, type JSX } from "react";

type FieldEditorHeaderProps = JSX.IntrinsicElements['div'] & {
    sortField: string;
}

export const FieldEditorHeader = forwardRef<HTMLDivElement, FieldEditorHeaderProps>(function FieldEditorHeader(
    { sortField, ...props }, ref) {
    return (
        <div className="mb-3 flex items-center justify-between"
            ref={ref}
            {...props}
        >
            <h3 className="text-sm font-medium text-[var(--app-text)]">Fields</h3>
            <div className="text-xs text-[var(--app-text-soft)]">
                Sort field: <span className="font-medium text-[var(--app-accent)]">{sortField}</span>
            </div>
        </div>
    )
})
