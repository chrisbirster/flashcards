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
            className="w-4 h-4 text-blue-600 rounded border-gray-300 focus:ring-blue-500"
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
            <label className="block text-xs text-gray-600 mb-1">Direction</label>
            <label className="flex items-center gap-2 cursor-pointer">
                <RtlOptionInput
                    isChecked={isChecked}
                    handleChange={handleChange}
                    datatestid={datatestid}
                    isPending={isPending}
                />
                <span className="text-xs text-gray-700">RTL</span>
            </label>
        </div>
    )
}