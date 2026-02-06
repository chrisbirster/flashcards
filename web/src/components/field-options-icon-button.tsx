import { forwardRef, type JSX } from "react"

const FieldOptionsIcon = () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4" />
    </svg>
)

type FieldOptionsIconButtonProps = JSX.IntrinsicElements['button'] & {
    handleClick: () => void;
    isEditing: boolean;
    isPending: boolean;
    datatestid?: string;
}

export const FieldOptionsIconButton = forwardRef<HTMLButtonElement, FieldOptionsIconButtonProps>(function FieldOptionsIconButton({
    handleClick,
    isEditing,
    isPending,
    datatestid,
    ...props
}, ref) {
    return (
        <button
            onClick={handleClick}
            className={`p-1 ${isEditing? 'text-blue-600' : 'text-gray-400 hover:text-gray-600'}`}
            title="Field options (font, size, RTL)"
            disabled={isPending}
            data-testid={datatestid}
            {...props}
            ref={ref}
        >
            <FieldOptionsIcon />
        </button>
    )
})