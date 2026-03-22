import { useEffect, useRef, useState, type ReactNode } from "react";
import { Link, NavLink, Outlet, useLocation } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAppRepository } from "#/lib/app-repository";
import { appNavigation, pageTitleForPath } from "#/lib/app-navigation";
import { ThemeToggle } from "#/components/theme-toggle";
import { AppTopBar } from "#/components/app-top-bar";
import { MobileBottomNav } from "#/components/mobile-bottom-nav";
import { MoreSheet } from "#/components/more-sheet";

function AppMark() {
  return (
    <Link to="/" className="flex items-center gap-3">
      <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] text-lg font-semibold text-[var(--app-accent-ink)] shadow-[0_12px_30px_rgba(112,214,108,0.18)]">
        V
      </span>
      <div>
        <p className="text-lg font-semibold text-[var(--app-text)]">Vutadex</p>
        <p className="text-xs uppercase tracking-[0.26em] text-[var(--app-muted)]">
          Flashcards workspace
        </p>
      </div>
    </Link>
  );
}

function DesktopSidebar({ onLogout }: { onLogout: () => void }) {
  const repository = useAppRepository();
  const { data: session } = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const userLabel =
    session?.user?.displayName || session?.user?.email || "User";
  const userInitial = userLabel.trim().charAt(0).toUpperCase() || "U";
  const plan = session?.entitlements?.plan?.toUpperCase() || "FREE";
  const [menuOpen, setMenuOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    if (!menuOpen) {
      return;
    }

    function handlePointerDown(event: PointerEvent) {
      if (!menuRef.current?.contains(event.target as Node)) {
        setMenuOpen(false);
      }
    }

    function handleEscape(event: KeyboardEvent) {
      if (event.key === "Escape") {
        setMenuOpen(false);
      }
    }

    window.addEventListener("pointerdown", handlePointerDown);
    window.addEventListener("keydown", handleEscape);
    return () => {
      window.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("keydown", handleEscape);
    };
  }, [menuOpen]);

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
                end={item.to === "/"}
                className={({ isActive }) =>
                  [
                    "block rounded-2xl px-4 py-3 transition-colors",
                    isActive
                      ? "bg-[var(--app-accent)] text-[var(--app-accent-ink)] shadow-sm"
                      : "text-[var(--app-text-soft)] hover:bg-[var(--app-muted-surface)] hover:text-[var(--app-text)]",
                  ].join(" ")
                }
              >
                <div className="text-sm font-semibold">{item.label}</div>
                <div className="mt-1 text-xs opacity-75">
                  {item.description}
                </div>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>

      <div className="border-t border-[var(--app-line)] px-5 py-5">
        <div ref={menuRef} className="relative">
          {menuOpen ? (
            <div
              role="menu"
              className="absolute inset-x-0 bottom-full mb-3 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-2 shadow-[0_20px_60px_rgba(0,0,0,0.28)]"
            >
              <Link
                to="/settings"
                onClick={() => setMenuOpen(false)}
                role="menuitem"
                className="flex min-h-11 items-center rounded-[1rem] px-3 py-3 text-sm font-medium text-[var(--app-text)] transition hover:bg-[var(--app-muted-surface)]"
              >
                User settings
              </Link>
              {session?.workspace?.organizationId ? (
                <Link
                  to="/team"
                  onClick={() => setMenuOpen(false)}
                  role="menuitem"
                  className="mt-1 flex min-h-11 items-center rounded-[1rem] px-3 py-3 text-sm font-medium text-[var(--app-text)] transition hover:bg-[var(--app-muted-surface)]"
                >
                  Team
                </Link>
              ) : null}
              <button
                type="button"
                onClick={() => {
                  setMenuOpen(false);
                  onLogout();
                }}
                role="menuitem"
                className="mt-1 flex min-h-11 w-full items-center rounded-[1rem] px-3 py-3 text-left text-sm font-medium text-[var(--app-text-soft)] transition hover:bg-[var(--app-muted-surface)] hover:text-[var(--app-text)]"
              >
                Sign out
              </button>
            </div>
          ) : null}

          <button
            type="button"
            onClick={() => setMenuOpen((current) => !current)}
            className="w-full rounded-2xl bg-[var(--app-muted-surface)] p-4 text-left transition hover:bg-[var(--app-card)]"
            aria-haspopup="menu"
            aria-expanded={menuOpen}
          >
            <div className="flex items-center gap-3">
              <span className="flex h-11 w-11 items-center justify-center rounded-full bg-[var(--app-accent)] text-sm font-semibold text-[var(--app-accent-ink)]">
                {userInitial}
              </span>
              <div className="min-w-0 flex-1">
                <p className="truncate text-sm font-semibold text-[var(--app-text)]">
                  {userLabel}
                </p>
                <p className="text-xs uppercase tracking-[0.18em] text-[var(--app-muted)]">
                  {plan} plan
                </p>
              </div>
              <svg
                className={`h-4 w-4 text-[var(--app-muted)] transition-transform ${menuOpen ? "rotate-180" : ""}`}
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.8"
                  d="m6 9 6 6 6-6"
                />
              </svg>
            </div>
          </button>
        </div>
      </div>
    </div>
  );
}

export function Layout({ children }: { children?: ReactNode }) {
  const location = useLocation();
  const repository = useAppRepository();
  const queryClient = useQueryClient();
  const [moreOpen, setMoreOpen] = useState(false);

  const { data: session } = useQuery({
    queryKey: ["auth-session"],
    queryFn: () => repository.fetchSession(),
  });

  const logoutMutation = useMutation({
    mutationFn: () => repository.logout(),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["auth-session"] });
      queryClient.invalidateQueries({ queryKey: ["entitlements"] });
      queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    },
  });

  const title = pageTitleForPath(location.pathname);

  return (
    <div
      className="bg-[var(--app-bg)] text-[var(--app-text)]"
      style={{ minHeight: "100dvh" }}
    >
      <div className="flex min-h-[100dvh]">
        <aside className="hidden w-80 shrink-0 border-r border-[var(--app-line)] bg-[var(--app-panel)] md:block">
          <DesktopSidebar onLogout={() => logoutMutation.mutate()} />
        </aside>

        <div className="flex min-h-[100dvh] min-w-0 flex-1 flex-col">
          <AppTopBar
            title={title}
            trailing={
              <div className="flex items-center gap-3">
                <ThemeToggle compact />
                <Link
                  to="/notes/add"
                  className="hidden items-center rounded-2xl bg-[var(--app-accent)] px-4 py-2.5 text-sm font-medium text-[var(--app-accent-ink)] shadow-sm transition hover:brightness-105 md:inline-flex"
                >
                  Add note
                </Link>
              </div>
            }
          />

          <main className="app-shell-main flex-1 px-4 py-5 md:px-8 md:py-8">
            {children ?? <Outlet />}
          </main>
        </div>
      </div>

      <MobileBottomNav onOpenMore={() => setMoreOpen(true)} />
      <MoreSheet
        open={moreOpen}
        onClose={() => setMoreOpen(false)}
        session={session}
        onLogout={() => logoutMutation.mutate()}
      />
    </div>
  );
}
