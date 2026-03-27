import { useEffect, useRef, useState, type ReactNode } from "react";
import { Link, NavLink, Outlet, useLocation } from "react-router";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useAppRepository } from "#/lib/app-repository";
import { appNavigation, pageTitleForPath, type AppNavigationItem } from "#/lib/app-navigation";
import { ThemeToggle } from "#/components/theme-toggle";
import { AppTopBar } from "#/components/app-top-bar";
import { MobileBottomNav } from "#/components/mobile-bottom-nav";
import { MoreSheet } from "#/components/more-sheet";

function SidebarIcon({ item }: { item: AppNavigationItem }) {
  const stroke = {
    home: "M3 10.5 12 3l9 7.5V20a1 1 0 0 1-1 1h-5.5v-6h-5V21H4a1 1 0 0 1-1-1z",
    stats: "M5 19.5V11m7 8.5V6m7 13.5v-5",
    focus: "M12 2.75c.7 2.1 2.35 4.56 4.5 5.6a8.5 8.5 0 1 1-9 0c2.15-1.04 3.8-3.5 4.5-5.6Z",
    notes: "M7 4.5h10A1.5 1.5 0 0 1 18.5 6v12A1.5 1.5 0 0 1 17 19.5H7A1.5 1.5 0 0 1 5.5 18V6A1.5 1.5 0 0 1 7 4.5Zm2 3h6m-6 4h6m-6 4h4",
    marketplace: "M5 8.5 7.2 5h9.6L19 8.5V18a1.5 1.5 0 0 1-1.5 1.5h-11A1.5 1.5 0 0 1 5 18zm0 0h14M9 11.5h6",
    templates: "M4.5 5.5h15v5h-15zm0 8h9v5h-9zm12 0h3v5h-3z",
    decks: "M12 4 4.5 7.5 12 11l7.5-3.5L12 4Zm-7.5 8L12 15.5 19.5 12M4.5 16.5 12 20l7.5-3.5",
    "study-groups": "M9 10.5a3 3 0 1 0 0-6 3 3 0 0 0 0 6Zm6 1.5a2.5 2.5 0 1 0 0-5m-10.5 11v-1.2A4.8 4.8 0 0 1 9.3 12h.4a4.8 4.8 0 0 1 4.8 4.8V18m1-3.5a4.1 4.1 0 0 1 4 4",
  }[item.icon];

  return (
    <svg
      className="h-5 w-5 shrink-0"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      aria-hidden="true"
    >
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d={stroke} />
    </svg>
  );
}

function SidebarTooltip({
  item,
  collapsed,
}: {
  item: AppNavigationItem;
  collapsed: boolean;
}) {
  return (
    <div
      className={[
        "pointer-events-none absolute left-full top-1/2 z-40 ml-3 w-64 -translate-y-1/2 rounded-[1.25rem] border border-[var(--app-line-strong)] bg-[color:var(--app-header)] p-3 opacity-0 shadow-[0_20px_50px_rgba(0,0,0,0.22)] backdrop-blur transition duration-150 group-hover:opacity-100 group-focus-within:opacity-100",
        collapsed ? "translate-x-0" : "md:group-hover:translate-x-1 md:group-focus-within:translate-x-1",
      ].join(" ")}
    >
      <p className="text-sm font-semibold text-[var(--app-text)]">{item.label}</p>
      <p className="mt-1 text-sm leading-6 text-[var(--app-text-soft)]">{item.description}</p>
    </div>
  );
}

function AppMark({ collapsed }: { collapsed: boolean }) {
  return (
    <Link to="/" className="flex items-center gap-3">
      <span className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--app-accent)] text-lg font-semibold text-[var(--app-accent-ink)] shadow-[0_12px_30px_rgba(112,214,108,0.18)]">
        V
      </span>
      {!collapsed ? (
        <div>
          <p className="text-lg font-semibold text-[var(--app-text)]">Vutadex</p>
          <p className="text-xs uppercase tracking-[0.26em] text-[var(--app-muted)]">
            Flashcards workspace
          </p>
        </div>
      ) : null}
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
  const [collapsed, setCollapsed] = useState(() => {
    if (typeof window === "undefined") {
      return false;
    }
    return window.localStorage.getItem("vutadex.sidebar.collapsed") === "true";
  });
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

  useEffect(() => {
    if (typeof window === "undefined") {
      return;
    }
    window.localStorage.setItem("vutadex.sidebar.collapsed", String(collapsed));
  }, [collapsed]);

  return (
    <div
      className={`flex h-full flex-col overflow-visible transition-[width] duration-200 ${collapsed ? "w-24" : "w-80"}`}
    >
      <div className={`flex items-center justify-between ${collapsed ? "px-3" : "px-5"} py-5`}>
        <AppMark collapsed={collapsed} />
        <button
          type="button"
          onClick={() => setCollapsed((current) => !current)}
          className="hidden h-11 w-11 items-center justify-center rounded-2xl border border-[var(--app-line)] bg-[var(--app-card)] text-[var(--app-text-soft)] transition hover:border-[var(--app-line-strong)] hover:text-[var(--app-text)] md:inline-flex"
          aria-label={collapsed ? "Expand sidebar" : "Collapse sidebar"}
          title={collapsed ? "Expand sidebar" : "Collapse sidebar"}
        >
          <svg
            className={`h-4 w-4 transition-transform ${collapsed ? "rotate-180" : ""}`}
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            aria-hidden="true"
          >
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d="M15 6 9 12l6 6" />
          </svg>
        </button>
      </div>

      <nav className={`flex-1 ${collapsed ? "px-2" : "px-3"} pb-4`}>
        <ul className="space-y-1">
          {appNavigation.map((item) => (
            <li key={item.to} className="group relative">
              <NavLink
                to={item.to}
                end={item.to === "/"}
                className={({ isActive }) =>
                  [
                    "relative flex min-h-14 items-center rounded-2xl transition-colors",
                    collapsed ? "justify-center px-3 py-3" : "gap-3 px-4 py-3",
                    isActive
                      ? "bg-[var(--app-accent)] text-[var(--app-accent-ink)] shadow-sm"
                      : "text-[var(--app-text-soft)] hover:bg-[var(--app-muted-surface)] hover:text-[var(--app-text)]",
                  ].join(" ")
                }
                aria-label={item.label}
              >
                <SidebarIcon item={item} />
                {!collapsed ? <div className="text-sm font-semibold">{item.label}</div> : null}
              </NavLink>
              <SidebarTooltip item={item} collapsed={collapsed} />
            </li>
          ))}
        </ul>
      </nav>

      <div className={`border-t border-[var(--app-line)] ${collapsed ? "px-3" : "px-5"} py-5`}>
        <div ref={menuRef} className="relative">
          {menuOpen ? (
            <div
              role="menu"
              className={[
                "absolute bottom-full mb-3 rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-2 shadow-[0_20px_60px_rgba(0,0,0,0.28)]",
                collapsed ? "left-full ml-3 w-64" : "inset-x-0",
              ].join(" ")}
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
            className={[
              "w-full rounded-2xl bg-[var(--app-muted-surface)] text-left transition hover:bg-[var(--app-card)]",
              collapsed ? "flex h-16 items-center justify-center p-0" : "p-4",
            ].join(" ")}
            aria-haspopup="menu"
            aria-expanded={menuOpen}
            aria-label="Open account menu"
          >
            <div className={`flex items-center ${collapsed ? "justify-center" : "gap-3"}`}>
              <span className="flex h-11 w-11 items-center justify-center rounded-full bg-[var(--app-accent)] text-sm font-semibold text-[var(--app-accent-ink)]">
                {userInitial}
              </span>
              {!collapsed ? (
                <>
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
                </>
              ) : null}
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
      className="h-[100dvh] overflow-hidden bg-[var(--app-bg)] text-[var(--app-text)]"
    >
      <div className="flex h-full">
        <aside className="hidden h-full shrink-0 border-r border-[var(--app-line)] bg-[var(--app-panel)] md:block">
          <DesktopSidebar onLogout={() => logoutMutation.mutate()} />
        </aside>

        <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
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

          <main className="app-shell-main min-h-0 flex-1 overflow-y-auto px-4 py-5 md:px-8 md:py-8">
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
