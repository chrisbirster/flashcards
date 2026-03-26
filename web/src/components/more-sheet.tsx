import { Link } from "react-router";
import type { AuthSessionResponse } from "#/lib/api";
import { ThemeToggle } from "#/components/theme-toggle";
import { Sheet } from "#/components/sheet";

export function MoreSheet({
  open,
  onClose,
  session,
  onLogout,
}: {
  open: boolean;
  onClose: () => void;
  session?: AuthSessionResponse;
  onLogout: () => void;
}) {
  const userLabel =
    session?.user?.displayName || session?.user?.email || "User";
  const userInitial = userLabel.trim().charAt(0).toUpperCase() || "U";
  const plan = session?.entitlements.plan?.toUpperCase() || "FREE";

  return (
    <Sheet open={open} onClose={onClose} title="More">
      <div className="space-y-4">
        <div className="rounded-[1.5rem] border border-[var(--app-line)] bg-[var(--app-card)] p-4">
          <div className="flex items-center gap-3">
            <span className="flex h-12 w-12 items-center justify-center rounded-full bg-[var(--app-accent)] text-base font-semibold text-[var(--app-accent-ink)]">
              {userInitial}
            </span>
            <div className="min-w-0">
              <p className="truncate text-sm font-semibold text-[var(--app-text)]">
                {userLabel}
              </p>
              <p className="text-xs uppercase tracking-[0.2em] text-[var(--app-muted)]">
                {plan} plan
              </p>
            </div>
          </div>
        </div>

        <div className="grid gap-3">
          <Link
            to="/settings"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            User settings
          </Link>
          {session?.workspace?.organizationId ? (
            <Link
              to="/team"
              onClick={onClose}
              className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
            >
              Team
            </Link>
          ) : null}
          <Link
            to="/stats"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Stats
          </Link>
          <Link
            to="/focus"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Focus
          </Link>
          <Link
            to="/marketplace"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Marketplace
          </Link>
          <Link
            to="/templates"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Templates
          </Link>
          <Link
            to="/study-groups"
            onClick={onClose}
            className="inline-flex min-h-11 items-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text)]"
          >
            Study Groups
          </Link>
        </div>

        <ThemeToggle />

        <button
          type="button"
          onClick={() => {
            onLogout();
            onClose();
          }}
          className="inline-flex min-h-11 w-full items-center justify-center rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-card)] px-4 text-sm font-medium text-[var(--app-text-soft)]"
        >
          Sign out
        </button>
      </div>
    </Sheet>
  );
}
