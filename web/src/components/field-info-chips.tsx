import { forwardRef, type JSX } from "react"

type FieldInfoChipProps = JSX.IntrinsicElements['div'] & {
    booleanIndicator: boolean;
    datatestid?: string;
}

const FieldInfoChip = forwardRef<HTMLDivElement, FieldInfoChipProps>(function FieldInfoChips({
    booleanIndicator,
    datatestid,
    children,
    ...props
}, ref) {
    return (
        <div className="flex items-center gap-1" data-testid={datatestid} {...props} ref={ref}>
            {booleanIndicator && children}
        </div>
    )
})

const RtlChip = () => (<span className="text-xs text-purple-600 bg-purple-100 px-2 py-0.5 rounded">RTL</span>)

const RTLFieldInfoChip = forwardRef<HTMLDivElement, FieldInfoChipProps>(function FieldInfoChips({
    booleanIndicator,
    datatestid,
    ...props
}, ref) {
    return (
        <FieldInfoChip
            booleanIndicator={booleanIndicator}
            datatestid={datatestid}
            {...props}
            ref={ref}
        >
            <RtlChip />
        </FieldInfoChip>
    )
})

const SortChip = () => (<span className="text-xs text-blue-600 bg-blue-100 px-2 py-0.5 rounded">Sort</span>)

const SortFieldInfoChip = forwardRef<HTMLDivElement, FieldInfoChipProps>(function FieldInfoChips({
    booleanIndicator,
    datatestid,
    ...props
}, ref) {
    return (
        <FieldInfoChip
            booleanIndicator={booleanIndicator}
            datatestid={datatestid}
            {...props}
            ref={ref}
        >   
            <SortChip />
        </FieldInfoChip>
    )
})

export { RTLFieldInfoChip, SortFieldInfoChip }