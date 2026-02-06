export function TemplateFieldPreview({ previewContent, label }: { previewContent: string, label: string }) {
    return (
        <div>
            <div className="text-xs text-gray-400 mb-1">{label}</div>
            <div className="p-2 bg-gray-50 rounded text-sm whitespace-pre-wrap">
                {previewContent}
            </div>
        </div>
    )
}