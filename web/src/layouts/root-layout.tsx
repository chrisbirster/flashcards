import { useState, type ReactNode } from 'react'
import { Link, NavLink, Outlet, useLocation } from 'react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'
import { appNavigation, pageTitleForPath } from '#/lib/app-navigation'

function AppMark() {
  return (
    <Link to="/" className="flex items-center gap-3">
      <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-950 text-lg font-semibold text-white">
        V
      </span>
      <div>
        <p className="text-lg font-semibold text-slate-950">Vutadex</p>
        <p className="text-xs uppercase tracking-[0.26em] text-slate-400">Flashcards workspace</p>
      </div>
    </Link>
  )
}

function SidebarContent({ onNavigate }: {onNavigate?: () => void}) {
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
  const plan = session?.entitlements?.plan?.toUpperCase() || 'FREE'

  return (
    <div className="flex h-full flex-col">
      <div className="px-5 py-5">
        <AppMark />
      </div>

      <nav className="flex-1 px-3 pb-4">
        <ul className="space-y-1">
          {appNavigation.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                end={item.to === '/'}
                onClick={onNavigate}
                className={({ isActive }) =>
                  [
                    'block rounded-2xl px-4 py-3 transition-colors',
                    isActive
                      ? 'bg-slate-950 text-white shadow-sm'
                      : 'text-slate-600 hover:bg-slate-100 hover:text-slate-950',
                  ].join(' ')
                }
              >
                <div className="text-sm font-semibold">{item.label}</div>
                <div className="mt-1 text-xs opacity-75">{item.description}</div>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>

      <div className="border-t border-slate-200 px-5 py-5">
        <div className="rounded-2xl bg-slate-100 p-4">
          <div className="flex items-center gap-3">
            <span className="flex h-11 w-11 items-center justify-center rounded-full bg-slate-950 text-sm font-semibold text-white">
              {userInitial}
            </span>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-semibold text-slate-950">{userLabel}</p>
              <p className="text-xs uppercase tracking-[0.18em] text-slate-500">{plan} plan</p>
            </div>
          </div>
          <button
            type="button"
            onClick={() => logoutMutation.mutate()}
            className="mt-4 w-full rounded-xl border border-slate-300 px-3 py-2 text-sm font-medium text-slate-700 hover:border-slate-400 hover:bg-white"
          >
            Sign out
          </button>
        </div>
      </div>
    </div>
  )
}

export function Layout({children}: {children?: ReactNode}) {
  const location = useLocation()
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <div className="min-h-screen bg-stone-100 text-slate-950">
      <div className="flex min-h-screen">
        <aside className="hidden w-80 shrink-0 border-r border-slate-200 bg-white md:block">
          <SidebarContent />
        </aside>

        <div className="flex min-h-screen min-w-0 flex-1 flex-col">
          <header className="sticky top-0 z-30 border-b border-slate-200 bg-stone-100/90 backdrop-blur">
            <div className="flex items-center justify-between gap-4 px-4 py-4 md:px-8">
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={() => setMobileOpen(true)}
                  className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-slate-300 bg-white text-slate-700 md:hidden"
                  aria-label="Open navigation"
                >
                  <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M4 7h16M4 12h16M4 17h16" />
                  </svg>
                </button>
                <div>
                  <p className="text-xs uppercase tracking-[0.24em] text-slate-400">Workspace</p>
                  <h1 className="text-2xl font-semibold tracking-tight">{pageTitleForPath(location.pathname)}</h1>
                </div>
              </div>

              <Link
                to="/notes/add"
                className="inline-flex items-center rounded-2xl bg-slate-950 px-4 py-2.5 text-sm font-medium text-white shadow-sm hover:bg-slate-800"
              >
                Add note
              </Link>
            </div>
          </header>

          <main className="flex-1 px-4 py-6 md:px-8 md:py-8">{children ?? <Outlet />}</main>
        </div>
      </div>

      {mobileOpen && (
        <div className="fixed inset-0 z-50 md:hidden">
          <button
            type="button"
            onClick={() => setMobileOpen(false)}
            className="absolute inset-0 bg-slate-950/45"
            aria-label="Close navigation"
          />
          <div className="absolute inset-y-0 left-0 w-[88vw] max-w-sm bg-white shadow-2xl">
            <div className="flex items-center justify-between border-b border-slate-200 px-5 py-4">
              <AppMark />
              <button
                type="button"
                onClick={() => setMobileOpen(false)}
                className="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-slate-300 text-slate-700"
                aria-label="Close navigation"
              >
                <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M6 6l12 12M18 6L6 18" />
                </svg>
              </button>
            </div>
            <SidebarContent onNavigate={() => setMobileOpen(false)} />
          </div>
        </div>
      )}
    </div>
  )
}
