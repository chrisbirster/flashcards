import { NavLink } from 'react-router';

type HeaderMenuItemProps = {
    item: { path: string; label: string };
}

export function HeaderMenuItem({
    item,
}: HeaderMenuItemProps
) {

    return (
        <NavLink
            key={item.path}
            to={item.path}
            className={({ isActive }) => `px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                isActive ? 'bg-blue-100 text-blue-700' : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
            }`}
        >
            {item.label}
        </NavLink>
    )
}