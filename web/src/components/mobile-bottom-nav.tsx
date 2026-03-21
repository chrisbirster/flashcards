import { NavLink } from 'react-router'

function navItemBase(isActive: boolean) {
  return [
    'flex min-h-11 flex-1 flex-col items-center justify-center rounded-2xl px-2 py-2 text-[11px] font-medium transition',
    isActive
      ? 'bg-[var(--app-accent)] text-[var(--app-accent-ink)]'
      : 'text-[var(--app-text-soft)] hover:bg-[var(--app-muted-surface)] hover:text-[var(--app-text)]',
  ].join(' ')
}

function NavGlyph({ path }: { path: string }) {
  return (
    <svg className="mb-1 h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" aria-hidden="true">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="1.8" d={path} />
    </svg>
  )
}

export function MobileBottomNav({ onOpenMore }: { onOpenMore: () => void }) {
  return (
    <nav
      className="fixed inset-x-0 bottom-0 z-40 border-t border-[var(--app-line)] bg-[color:var(--app-header)] px-3 py-3 backdrop-blur md:hidden"
      style={{ paddingBottom: 'calc(0.75rem + env(safe-area-inset-bottom))' }}
    >
      <div className="flex items-center gap-2">
        <NavLink to="/" end className={({ isActive }) => navItemBase(isActive)}>
          <NavGlyph path="M3 10.5 12 3l9 7.5V20a1 1 0 0 1-1 1h-5.5v-6h-5V21H4a1 1 0 0 1-1-1z" />
          Home
        </NavLink>
        <NavLink to="/notes/view" className={({ isActive }) => navItemBase(isActive)}>
          <NavGlyph path="M7 4.5h10A1.5 1.5 0 0 1 18.5 6v12A1.5 1.5 0 0 1 17 19.5H7A1.5 1.5 0 0 1 5.5 18V6A1.5 1.5 0 0 1 7 4.5Zm2 3h6m-6 4h6m-6 4h4" />
          Notes
        </NavLink>
        <NavLink to="/notes/add" className={({ isActive }) => navItemBase(isActive)}>
          <NavGlyph path="M12 5v14M5 12h14" />
          Add
        </NavLink>
        <NavLink to="/decks" className={({ isActive }) => navItemBase(isActive)}>
          <NavGlyph path="M12 4 4.5 7.5 12 11l7.5-3.5L12 4Zm-7.5 8L12 15.5 19.5 12M4.5 16.5 12 20l7.5-3.5" />
          Decks
        </NavLink>
        <button
          type="button"
          onClick={onOpenMore}
          className={navItemBase(false)}
        >
          <NavGlyph path="M5 12h.01M12 12h.01M19 12h.01" />
          More
        </button>
      </div>
    </nav>
  )
}
