import { useTheme } from '#/lib/theme'

function SunIcon() {
  return (
    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true">
      <circle cx="12" cy="12" r="4.5" strokeWidth="1.8" />
      <path strokeLinecap="round" strokeWidth="1.8" d="M12 2.5v2.2M12 19.3v2.2M4.7 4.7l1.6 1.6M17.7 17.7l1.6 1.6M2.5 12h2.2M19.3 12h2.2M4.7 19.3l1.6-1.6M17.7 6.3l1.6-1.6" />
    </svg>
  )
}

function MoonIcon() {
  return (
    <svg className="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.8"
        d="M20 14.2A7.8 7.8 0 0 1 9.8 4a8.5 8.5 0 1 0 10.2 10.2Z"
      />
    </svg>
  )
}

export function ThemeToggle({ compact = false }: { compact?: boolean }) {
  const { theme, toggleTheme } = useTheme()
  const label = theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className={[
        'inline-flex items-center gap-2 rounded-2xl border px-3 py-2 text-sm font-medium transition',
        'border-[var(--app-line-strong)] bg-[var(--app-card)] text-[var(--app-text-soft)] hover:border-[var(--app-accent)] hover:text-[var(--app-text)]',
      ].join(' ')}
      aria-label={label}
      title={label}
    >
      <span className="inline-flex h-8 w-8 items-center justify-center rounded-full bg-[var(--app-muted-surface)] text-[var(--app-accent)]">
        {theme === 'dark' ? <SunIcon /> : <MoonIcon />}
      </span>
      {!compact && <span>{theme === 'dark' ? 'Light mode' : 'Dark mode'}</span>}
    </button>
  )
}
