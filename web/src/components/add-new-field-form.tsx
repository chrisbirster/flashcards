import { forwardRef, type JSX } from "react";

type AddNewFieldFormProps = JSX.IntrinsicElements['form'] & {
    newFieldName: string;
    setNewFieldName: (name: string) => void;
    handleAddField: (e: React.FormEvent) => void;
    isPending: boolean;
}
          <form onSubmit={handleAddField} className="flex gap-2">
            <input
              id="new-field-name"
              type="text"
              value={newFieldName}
              onChange={(e) => setNewFieldName(e.target.value)}
              placeholder="New field name..."
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              disabled={isPending}
            />
            <button
              type="submit"
              disabled={!newFieldName.trim() || isPending}
              className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed"
            >
              Add
            </button>
          </form>