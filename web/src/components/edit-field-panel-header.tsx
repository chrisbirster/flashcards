import { EditFieldPanelCloseButton } from "./edit-field-panel-close-button";

type EditFieldPanelHeaderProps = {
    noteTypeName: string;
    onClose: () => void;
}

export function EditFieldPanelHeader({ noteTypeName, onClose }: EditFieldPanelHeaderProps) {
    return (
        <div className="flex items-start sm:items-center justify-between gap-3 p-3 sm:p-4 border-b">
            <h2 className="text-base sm:text-lg font-semibold text-gray-900">
                Edit Fields: {noteTypeName}
            </h2>
            <EditFieldPanelCloseButton onClick={onClose} data-testid="close-field-editor" />
        </div>
    )

}
