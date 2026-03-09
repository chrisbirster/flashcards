import { Link } from "react-router"
import { HeaderMenuItem } from "./header-menu-item"

const navItems = [
    { path: '/decks', label: 'Decks' },
    { path: '/notes/add', label: 'Add Note' },
    { path: '/templates', label: 'Templates' },
    { path: '/tools/empty-cards', label: 'Empty Cards' },
    // Future routes
    // { path: '/browse', label: 'Browse' },
    // { path: '/stats', label: 'Statistics' },
    // { path: '/settings', label: 'Settings' },
]

export function HeaderMenuLinks() {
    return (
        <div className="flex w-full flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <Link to="/decks" className="flex items-center shrink-0">
                <h1 className="text-2xl font-bold text-gray-900">Microdote</h1>
            </Link>
            <div className="flex gap-1 overflow-x-auto pb-1 sm:overflow-visible sm:pb-0">
                {navItems.map((item) => (
                    <HeaderMenuItem key={item.path} item={item} />
                ))}
            </div>
        </div>
    )
}
