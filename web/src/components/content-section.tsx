import { Outlet } from "react-router";

export function ContentSection() {
    return (
        <main className="py-4 sm:py-8 px-3 sm:px-4">
            <Outlet />
        </main>
    )
}
