import { Link } from "react-router"
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { HeaderMenuItem } from "./header-menu-item"
import { useAppRepository } from '#/lib/app-repository'

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
    const repository = useAppRepository()
    const queryClient = useQueryClient()
    const { data: session } = useQuery({
        queryKey: ['auth-session'],
        queryFn: () => repository.fetchSession(),
    })

    const logoutMutation = useMutation({
        mutationFn: () => repository.logout(),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['auth-session'] })
            queryClient.invalidateQueries({ queryKey: ['entitlements'] })
        },
    })

    const userLabel = session?.user?.displayName || session?.user?.email || 'User'
    const userInitial = userLabel.trim().charAt(0).toUpperCase() || 'U'

    return (
        <div className="flex w-full flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
            <Link to="/decks" className="flex items-center shrink-0">
                <div>
                    <h1 className="text-2xl font-bold text-gray-900">Vutadex</h1>
                    <p className="text-xs uppercase tracking-[0.2em] text-gray-400">Browser-first flashcards</p>
                </div>
            </Link>
            <div className="flex flex-col gap-2 sm:flex-row sm:items-center">
                <div className="flex gap-1 overflow-x-auto pb-1 sm:overflow-visible sm:pb-0">
                    {navItems.map((item) => (
                        <HeaderMenuItem key={item.path} item={item} />
                    ))}
                </div>
                <div className="flex items-center gap-2">
                    {session?.entitlements?.plan && (
                        <span className="inline-flex items-center rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-600">
                            {session.entitlements.plan.toUpperCase()}
                        </span>
                    )}
                    {session?.authenticated ? (
                        <>
                            <span className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-slate-900 text-sm font-semibold text-white">
                                {userInitial}
                            </span>
                            <div className="hidden min-w-0 sm:block">
                                <p className="truncate text-sm font-medium text-gray-700">{userLabel}</p>
                                <p className="text-xs uppercase tracking-[0.24em] text-gray-400">Signed in</p>
                            </div>
                            <button
                                type="button"
                                onClick={() => logoutMutation.mutate()}
                                className="px-3 py-2 text-sm text-gray-600 hover:bg-gray-100 rounded-md"
                            >
                                Sign Out
                            </button>
                        </>
                    ) : (
                        <Link to="/login" className="px-3 py-2 text-sm text-white bg-slate-900 hover:bg-slate-700 rounded-md">
                            Log In
                        </Link>
                    )}
                </div>
            </div>
        </div>
    )
}
