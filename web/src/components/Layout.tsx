import { Outlet, Link, useLocation } from 'react-router-dom'

export function Layout() {
  const location = useLocation()

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

  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-8">
              <Link to="/decks" className="flex items-center">
                <h1 className="text-2xl font-bold text-gray-900">Microdote</h1>
              </Link>
              <div className="flex gap-1">
                {navItems.map((item) => (
                  <Link
                    key={item.path}
                    to={item.path}
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      location.pathname === item.path || location.pathname.startsWith(item.path)
                        ? 'bg-blue-100 text-blue-700'
                        : 'text-gray-600 hover:bg-gray-100 hover:text-gray-900'
                    }`}
                  >
                    {item.label}
                  </Link>
                ))}
              </div>
            </div>
          </div>
        </div>
      </nav>
      <main className="py-8 px-4">
        <Outlet />
      </main>
    </div>
  )
}
