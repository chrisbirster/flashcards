import type { ReactNode } from 'react'

export function IconButton({
    title,
    testId,
    icon,
    handleClick,
}: {
    title: string,
    testId?: string
   icon: ReactNode
   handleClick: () => void
}) {
    return (
        <button
            type="button"
            onClick={handleClick}
            className="px-3 py-2 text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200"
            title={title}
            data-testid={testId}
        >
            {icon}
        </button>
    )
}