import { Outlet } from "react-router";

export function ContentSection() {
    return (
        <main className="py-8 px-4">
            <Outlet />
        </main>
    )
}