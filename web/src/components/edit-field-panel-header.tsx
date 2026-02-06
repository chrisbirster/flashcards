import { EditFieldPanelCloseButton } from "./edit-field-panel-close-button";

type EditFieldPanelHeaderProps = {
    noteTypeName: string;
    onClose: () => void;
}

export function EditFieldPanelHeader({ noteTypeName, onClose }: EditFieldPanelHeaderProps) {
    return (
        <div className="flex items-center justify-between p-4 border-b">
            <h2 className="text-lg font-semibold text-gray-900">
                Edit Fields: {noteTypeName}
            </h2>
            <EditFieldPanelCloseButton onClick={onClose} />
        </div>
    )

}
