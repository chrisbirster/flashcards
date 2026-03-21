import { EditFieldPanelCloseButton } from "./edit-field-panel-close-button";

type EditFieldPanelHeaderProps = {
    noteTypeName: string;
    onClose: () => void;
}

export function EditFieldPanelHeader({ noteTypeName, onClose }: EditFieldPanelHeaderProps) {
  return (
    <div className="flex items-start justify-between gap-3 border-b border-[var(--app-line)] bg-[color:var(--app-header)]/95 p-3 backdrop-blur sm:items-center sm:p-4">
      <h2 className="text-base font-semibold text-[var(--app-text)] sm:text-lg">
        Edit Fields: {noteTypeName}
      </h2>
      <EditFieldPanelCloseButton onClick={onClose} data-testid="close-field-editor" />
    </div>
  )
}
