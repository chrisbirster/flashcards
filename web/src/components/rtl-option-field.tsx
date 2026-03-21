import { forwardRef, type JSX } from "react"

type RtlOptionFieldProps = JSX.IntrinsicElements['input'] & {
    isChecked: boolean;
    handleChange: () => void;
    datatestid: string;
    isPending: boolean;
}

const RtlOptionInput = forwardRef<HTMLInputElement, RtlOptionFieldProps>(function RtlOptionInput({
    isChecked,
    handleChange,
    datatestid,
    isPending,
    ...props
}, ref) {
    return (
        <input
            type="checkbox"
            checked={isChecked}
            onChange={handleChange}
            className="h-4 w-4 rounded border-[var(--app-line-strong)] text-[var(--app-accent)]"
            disabled={isPending}
            data-testid={datatestid}
            {...props}
            ref={ref}
        />
    )
})


export const RtlOptionField = ({
    isChecked,
    handleChange,
    datatestid,
    isPending,
}: RtlOptionFieldProps 
) => {
    return (
        <div>
            <label className="mb-1 block text-xs text-[var(--app-text-soft)]">Direction</label>
            <label className="flex items-center gap-2 cursor-pointer">
                <RtlOptionInput
                    isChecked={isChecked}
                    handleChange={handleChange}
                    datatestid={datatestid}
                    isPending={isPending}
                />
                <span className="text-xs text-[var(--app-text)]">RTL</span>
            </label>
        </div>
    )
}
