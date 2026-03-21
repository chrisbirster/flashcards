export function TemplateFieldPreview({ previewContent, label }: { previewContent: string, label: string }) {
  return (
    <div className="space-y-2">
      <div className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">{label}</div>
      <div className="rounded-2xl border border-[var(--app-line)] bg-[var(--app-card-strong)] p-3 text-sm whitespace-pre-wrap text-[var(--app-text)]">
        {previewContent}
      </div>
    </div>
  )
}
