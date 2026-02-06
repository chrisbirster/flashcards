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
        <div className="flex items-center gap-8">
            <Link to="/decks" className="flex items-center">
                <h1 className="text-2xl font-bold text-gray-900">Microdote</h1>
            </Link>
            <div className="flex gap-1">
                {navItems.map((item) => (
                    <HeaderMenuItem key={item.path} item={item} />
                ))}
            </div>
        </div>
    )
}
