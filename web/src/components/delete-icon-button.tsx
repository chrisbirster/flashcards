import { forwardRef, type JSX } from "react"

const DeleteIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
    </svg>
)

type DeleteButtonProps = JSX.IntrinsicElements['button'] & {
    disabled: boolean;
    onDelete: () => void;
}

export const DeleteButton = forwardRef<HTMLButtonElement, DeleteButtonProps>(function DeleteButton({
    disabled,
    onDelete,
    ...props
}, ref) {
    return (
        <button
            onClick={onDelete}
            disabled={disabled}
            className="p-1 text-red-400 hover:text-red-600 disabled:opacity-30 disabled:cursor-not-allowed"
            title="Remove field"
            ref={ref}
            {...props}
        >
            <DeleteIcon />
        </button>
    )
})
