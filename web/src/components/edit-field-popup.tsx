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
        <div className="ml-4 p-3 bg-blue-50 border border-blue-200 rounded-md" data-testid={`field-options-panel-${field}`}
            {...props}
            ref={ref}
        >
            <div className="text-xs font-medium text-gray-700 mb-2">Field Options: {field}</div>
            <div className="grid grid-cols-3 gap-3">
                {children}
            </div>
        </div>
    )
})
