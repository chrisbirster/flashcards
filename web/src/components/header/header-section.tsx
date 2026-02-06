import { HeaderMenuLinks } from "./header-menu-links";

export function Headersection() {
    return (
        <nav className="bg-white shadow-sm border-b border-gray-200">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between items-center h-16">
                    <HeaderMenuLinks />
                </div>
            </div>
        </nav>
    )
}