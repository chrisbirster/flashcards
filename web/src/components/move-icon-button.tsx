import { forwardRef, type JSX } from "react";

const UpArrowIcon = forwardRef<SVGSVGElement, JSX.IntrinsicElements['svg']>(function UpArrowIcon(props, ref) {
    return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ref={ref}
            {...props}
        >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
        </svg>
    )
})

const DownArrowIcon = forwardRef<SVGSVGElement, JSX.IntrinsicElements['svg']>(function UpArrowIcon(props, ref) {
    return (
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"
            ref={ref}
            {...props}
        >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
    )
})

type IconButtonProps = JSX.IntrinsicElements['button'] & {
    disabled: boolean;
    handleClick: () => void;
    title?: string;
    svg: typeof UpArrowIcon | typeof DownArrowIcon;
}

const IconButton = forwardRef<HTMLButtonElement, IconButtonProps>(function IconButton({
    disabled,
    handleClick,
    title,
    children,
    ...props
}, ref
) {
    return (
        <button
            onClick={handleClick}
            disabled={disabled}
            className="rounded-xl p-1 text-[var(--app-text-soft)] hover:text-[var(--app-text)] disabled:cursor-not-allowed disabled:opacity-30"
            title={title}
            ref={ref}
            {...props}
        >
            {children}
        </button>
    )
})

function MoveUpIconButton({
    disabled,
    handleClick,
}: {
    disabled: boolean;
    handleClick: () => void;
}) {
    return (
        <IconButton disabled={disabled} handleClick={handleClick} title={"Move up"} svg={UpArrowIcon}>
            <UpArrowIcon />
        </IconButton>
    )
}

function MoveDownIconButton({
    disabled,
    handleClick,
}: {
    disabled: boolean;
    handleClick: () => void;
}) {
    return (
        <IconButton disabled={disabled} handleClick={handleClick} title={"Move down"} svg={UpArrowIcon}>
            <DownArrowIcon />
        </IconButton>
    )
}

export {
    MoveDownIconButton,
    MoveUpIconButton,
}
