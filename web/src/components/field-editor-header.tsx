import { forwardRef, type JSX } from "react";

type FieldEditorHeaderProps = JSX.IntrinsicElements['div'] & {
    sortField: string;
}

export const FieldEditorHeader = forwardRef<HTMLDivElement, FieldEditorHeaderProps>(function FieldEditorHeader(
    { sortField, ...props }, ref) {
    return (
        <div className="flex items-center justify-between mb-2"
            ref={ref}
            {...props}
        >
            <h3 className="text-sm font-medium text-gray-700">Fields</h3>
            <div className="text-xs text-gray-500">
                Sort field: <span className="font-medium text-blue-600">{sortField}</span>
            </div>
        </div>
    )
})