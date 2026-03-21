import { useState, type ReactNode } from 'react'
import { Link, NavLink, Outlet, useLocation } from 'react-router'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useAppRepository } from '#/lib/app-repository'
import { appNavigation, pageTitleForPath } from '#/lib/app-navigation'
import { ThemeToggle } from '#/components/theme-toggle'

function AppMark() {
  return (
    <Link to="/" className="flex items-center gap-3">
      <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] text-lg font-semibold text-[var(--app-accent-ink)] shadow-[0_12px_30px_rgba(112,214,108,0.18)]">
        V
      </span>
      <div>
        <p className="text-lg font-semibold text-[var(--app-text)]">Vutadex</p>
        <p className="text-xs uppercase tracking-[0.26em] text-[var(--app-muted)]">Flashcards workspace</p>
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
                      ? 'bg-[var(--app-accent)] text-[var(--app-accent-ink)] shadow-sm'
                      : 'text-[var(--app-text-soft)] hover:bg-[var(--app-muted-surface)] hover:text-[var(--app-text)]',
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

      <div className="border-t border-[var(--app-line)] px-5 py-5">
        <div className="rounded-2xl bg-[var(--app-muted-surface)] p-4">
          <div className="flex items-center gap-3">
            <span className="flex h-11 w-11 items-center justify-center rounded-full bg-[var(--app-accent)] text-sm font-semibold text-[var(--app-accent-ink)]">
              {userInitial}
            </span>
            <div className="min-w-0 flex-1">
              <p className="truncate text-sm font-semibold text-[var(--app-text)]">{userLabel}</p>
              <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">{plan} plan</p>
            </div>
          </div>
          <button
            type="button"
            onClick={() => logoutMutation.mutate()}
            className="mt-4 w-full rounded-xl border border-[var(--app-line-strong)] px-3 py-2 text-sm font-medium text-[var(--app-text-soft)] hover:border-[var(--app-accent)] hover:bg-[var(--app-card)] hover:text-[var(--app-text)]"
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
    <div className="min-h-screen bg-[var(--app-bg)] text-[var(--app-text)]">
      <div className="flex min-h-screen">
        <aside className="hidden w-80 shrink-0 border-r border-[var(--app-line)] bg-[var(--app-panel)] md:block">
          <SidebarContent />
        </aside>

        <div className="flex min-h-screen min-w-0 flex-1 flex-col">
          <header className="sticky top-0 z-30 border-b border-[var(--app-line)] bg-[color:var(--app-header)] backdrop-blur">
            <div className="flex items-center justify-between gap-4 px-4 py-4 md:px-8">
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  onClick={() => setMobileOpen(true)}
                  className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text-soft)] md:hidden"
                  aria-label="Open navigation"
                >
                  <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.8} d="M4 7h16M4 12h16M4 17h16" />
                  </svg>
                </button>
                <div>
                  <p className="text-xs uppercase tracking-[0.24em] text-[var(--app-muted)]">Workspace</p>
                  <h1 className="text-2xl font-semibold tracking-tight">{pageTitleForPath(location.pathname)}</h1>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <ThemeToggle compact />
                <Link
                  to="/notes/add"
                  className="inline-flex items-center rounded-2xl bg-[var(--app-accent)] px-4 py-2.5 text-sm font-medium text-[var(--app-accent-ink)] shadow-sm transition hover:brightness-105"
                >
                  Add note
                </Link>
              </div>
            </div>
          </header>

          <main className="flex-1 px-4 py-6 md:px-8 md:py-8">{children ?? <Outlet />}</main>
        </div>
      </div>

      {mobileOpen && (
        <div className="fixed inset-0 z-50 md:hidden">
          <button type="button" onClick={() => setMobileOpen(false)} className="absolute inset-0 bg-black/55" aria-label="Close navigation" />
          <div className="absolute inset-y-0 left-0 w-[88vw] max-w-sm bg-[var(--app-panel)] shadow-2xl">
            <div className="flex items-center justify-between border-b border-[var(--app-line)] px-5 py-4">
              <AppMark />
              <button
                type="button"
                onClick={() => setMobileOpen(false)}
                className="inline-flex h-10 w-10 items-center justify-center rounded-2xl border border-[var(--app-line-strong)] text-[var(--app-text-soft)]"
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
