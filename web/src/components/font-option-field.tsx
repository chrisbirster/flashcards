import { forwardRef, type JSX } from "react"

const FONT_OPTIONS = [
    { value: '', label: 'Default' },
    { value: 'Arial', label: 'Arial' },
    { value: 'Times New Roman', label: 'Times New Roman' },
    { value: 'Georgia', label: 'Georgia' },
    { value: 'Verdana', label: 'Verdana' },
    { value: 'Courier New', label: 'Courier New' },
    { value: 'Comic Sans MS', label: 'Comic Sans MS' },
]

const FONT_SIZE_OPTIONS = [
    { value: 0, label: 'Default' },
    { value: 12, label: '12px' },
    { value: 14, label: '14px' },
    { value: 16, label: '16px' },
    { value: 18, label: '18px' },
    { value: 20, label: '20px' },
    { value: 24, label: '24px' },
    { value: 28, label: '28px' },
    { value: 32, label: '32px' },
]

type FontOptionFieldProps = JSX.IntrinsicElements['select'] & {
    fieldValue: number | undefined;
    handleChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
    isPending: boolean;
    datatestid: string;
    options: typeof FONT_SIZE_OPTIONS | typeof FONT_OPTIONS;
}

const FontOptionFieldSelect = forwardRef<HTMLSelectElement, FontOptionFieldProps>(function FontOptionFieldSelect({
    fieldValue,
    handleChange,
    isPending,
    datatestid,
    options,
    ...props
}, ref) {

    return (
        <select
            value={fieldValue ?? 0}
            onChange={handleChange}
            className="w-full px-2 py-1 text-xs border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-blue-500"
            disabled={isPending}
            data-testid={datatestid}
            {...props}
            ref={ref}
        >
            {options.map(opt => (
                <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
        </select>
    )
})

function FontTypeOptionField({
    fieldValue,
    handleChange,
    isPending,
    datatestid,
}: Omit<FontOptionFieldProps, 'options'>
) {
    return (
        <div>
            <label className="block text-xs text-gray-600 mb-1">Size</label>
            <FontOptionFieldSelect
                fieldValue={fieldValue}
                handleChange={handleChange}
                isPending={isPending}
                datatestid={datatestid}
                options={FONT_OPTIONS}
            />
        </div>
    )
}

function FontSizeOptionField({
    fieldValue,
    handleChange,
    isPending,
    datatestid,
}: Omit<FontOptionFieldProps, 'options'>
) {
    return (
        <div>
            <label className="block text-xs text-gray-600 mb-1">Size</label>
            <FontOptionFieldSelect
                fieldValue={fieldValue}
                handleChange={handleChange}
                isPending={isPending}
                datatestid={datatestid}
                options={FONT_SIZE_OPTIONS}
            />
        </div>
    )
}

export { FontTypeOptionField, FontSizeOptionField }