import { HeaderMenuLinks } from "./header-menu-links";

export function Headersection() {
    return (
        <nav className="bg-white shadow-sm border-b border-gray-200">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="py-3 sm:py-0 sm:h-16 flex items-center">
                    <HeaderMenuLinks />
                </div>
            </div>
        </nav>
    )
}
