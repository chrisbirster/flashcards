import type { ReactNode } from 'react'
import { Outlet } from 'react-router'
import { Headersection } from '../components/header/header-section'

export function Layout({children}: {children?: ReactNode}) {

    return (
        <div className="min-h-screen bg-gray-50">
            <Headersection />
            {children ?? <Outlet />}
        </div>
    )
}
