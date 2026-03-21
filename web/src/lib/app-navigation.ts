export interface AppNavigationItem {
  to: string
  label: string
  description: string
}

export const appNavigation: AppNavigationItem[] = [
  { to: '/', label: 'Home', description: 'Overview and usage' },
  { to: '/notes/view', label: 'Notes', description: 'Browse and edit notes' },
  { to: '/templates', label: 'Templates', description: 'Manage card templates' },
  { to: '/decks', label: 'Decks', description: 'Organize study decks' },
  { to: '/study-groups', label: 'Study Groups', description: 'Shared learning spaces' },
]

export function pageTitleForPath(pathname: string): string {
  if (pathname === '/') return 'Home'
  if (pathname.startsWith('/notes/view')) return 'Notes'
  if (pathname.startsWith('/notes/add')) return 'Add Note'
  if (pathname.startsWith('/templates')) return 'Templates'
  if (pathname.startsWith('/decks')) return 'Decks'
  if (pathname.startsWith('/study-groups')) return 'Study Groups'
  if (pathname.startsWith('/study/')) return 'Study'
  if (pathname.startsWith('/tools/empty-cards')) return 'Empty Cards'
  return 'Vutadex'
}
