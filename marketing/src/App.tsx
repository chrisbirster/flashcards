import { useEffect, useMemo, useState } from 'react'

type SessionResponse = {
  authenticated: boolean
  user?: {
    displayName?: string
    email: string
  }
}

const appOrigin = (
  import.meta.env.VITE_APP_ORIGIN ??
  (import.meta.env.DEV ? 'http://localhost:8000' : 'https://app.vutadex.com')
).replace(/\/$/, '')

function userInitial(user?: SessionResponse['user']) {
  const label = user?.displayName || user?.email || ''
  return label.trim().charAt(0).toUpperCase() || 'V'
}

export function App() {
  const [session, setSession] = useState<SessionResponse | null>(null)

  useEffect(() => {
    let cancelled = false

    fetch(`${appOrigin}/api/auth/session`, {
      credentials: 'include',
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error('Session request failed')
        }
        return response.json() as Promise<SessionResponse>
      })
      .then((data) => {
        if (!cancelled) setSession(data)
      })
      .catch(() => {
        if (!cancelled) setSession({authenticated: false})
      })

    return () => {
      cancelled = true
    }
  }, [])

  const topRight = useMemo(() => {
    if (session?.authenticated) {
      return (
        <a
          href={`${appOrigin}/decks`}
          className="inline-flex items-center gap-3 rounded-full border px-3 py-2 text-sm font-semibold shadow-sm backdrop-blur transition hover:-translate-y-0.5"
          style={{borderColor: 'var(--line)', backgroundColor: 'rgba(255,252,247,0.86)', color: 'var(--ink)'}}
        >
          <span className="inline-flex h-9 w-9 items-center justify-center rounded-full text-white" style={{backgroundColor: 'var(--forest)'}}>
            {userInitial(session.user)}
          </span>
          <span className="hidden sm:inline">Open app</span>
        </a>
      )
    }

    return (
      <a
        href={`${appOrigin}/login`}
        className="inline-flex items-center rounded-full px-5 py-3 text-sm font-semibold text-white transition"
        style={{backgroundColor: 'var(--forest)', boxShadow: '0 10px 24px rgba(35,67,56,0.16)'}}
      >
        Log in
      </a>
    )
  }, [session])

  return (
    <div
      className="min-h-screen"
      style={{
        background:
          'radial-gradient(circle at top, var(--mist) 0%, transparent 32%), radial-gradient(circle at 86% 16%, var(--clay-wash) 0%, transparent 24%), linear-gradient(180deg, var(--canvas) 0%, var(--canvas-soft) 100%)',
      }}
    >
      <div className="mx-auto max-w-7xl px-5 py-6 sm:px-8">
        <header className="flex items-center justify-between gap-4">
          <a href="/" className="flex items-center gap-3">
            <span
              className="inline-flex h-11 w-11 items-center justify-center rounded-2xl text-base font-black text-white"
              style={{backgroundColor: 'var(--forest)'}}
            >
              V
            </span>
            <div>
              <p className="text-lg font-black tracking-tight" style={{color: 'var(--ink)'}}>Vutadex</p>
              <p className="text-xs uppercase tracking-[0.26em]" style={{color: 'var(--muted)'}}>Flashcards, refined</p>
            </div>
          </a>
          {topRight}
        </header>

        <main className="mt-10 grid gap-8 lg:grid-cols-[1.1fr_0.9fr]">
          <section
            className="rounded-[2.5rem] border p-8 shadow-[0_28px_80px_rgba(24,35,29,0.08)] backdrop-blur sm:p-12"
            style={{borderColor: 'var(--line)', backgroundColor: 'var(--surface)'}}
          >
            <p className="text-sm font-semibold uppercase tracking-[0.34em]" style={{color: 'var(--accent)'}}>The serious flashcard workspace</p>
            <h1 className="mt-5 max-w-3xl text-5xl font-black tracking-[-0.05em] sm:text-7xl" style={{color: 'var(--ink)'}}>
              Build decks with momentum. Study with less friction.
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8" style={{color: 'var(--muted)'}}>
              Vutadex gives you a clean note-building workflow, deck-level structure, and a browser-first foundation that is ready for desktop and mobile sync next.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <a
                href={`${appOrigin}/login`}
                className="inline-flex items-center rounded-full px-6 py-3 text-sm font-semibold text-white transition"
                style={{backgroundColor: 'var(--forest)', boxShadow: '0 14px 28px rgba(35,67,56,0.16)'}}
              >
                Start in the app
              </a>
              <a
                href="#pricing"
                className="inline-flex items-center rounded-full border px-6 py-3 text-sm font-semibold transition"
                style={{borderColor: 'var(--line-strong)', backgroundColor: 'rgba(255,252,247,0.75)', color: 'var(--ink)'}}
              >
                View plans
              </a>
            </div>

            <div className="mt-12 grid gap-4 sm:grid-cols-3">
              <article className="rounded-3xl p-5 text-white" style={{backgroundColor: 'var(--surface-deep)'}}>
                <p className="text-sm font-semibold" style={{color: '#f2d6c7'}}>Recent-note context</p>
                <p className="mt-2 text-sm leading-6 text-white/75">Keep note creation grounded by seeing the latest items in the deck while you build.</p>
              </article>
              <article className="rounded-3xl border p-5" style={{borderColor: 'var(--line)', backgroundColor: 'var(--surface-muted)'}}>
                <p className="text-sm font-semibold" style={{color: 'var(--ink)'}}>OTP sign-in</p>
                <p className="mt-2 text-sm leading-6" style={{color: 'var(--muted)'}}>Fast account access without a password reset maze or brittle local-only sessions.</p>
              </article>
              <article className="rounded-3xl border p-5" style={{borderColor: 'var(--line)', backgroundColor: 'var(--surface-strong)'}}>
                <p className="text-sm font-semibold" style={{color: 'var(--ink)'}}>Sync-ready foundation</p>
                <p className="mt-2 text-sm leading-6" style={{color: 'var(--muted)'}}>Web today, with desktop and mobile sync designed as a Pro extension rather than an afterthought.</p>
              </article>
            </div>
          </section>

          <section className="grid gap-5">
            <article
              className="rounded-[2rem] p-8 text-white shadow-[0_24px_70px_rgba(27,52,43,0.18)]"
              style={{background: 'linear-gradient(180deg, var(--surface-deep) 0%, var(--surface-deeper) 100%)'}}
            >
              <p className="text-xs font-semibold uppercase tracking-[0.28em]" style={{color: '#f2d6c7'}}>Why Vutadex</p>
              <ul className="mt-6 space-y-5 text-sm leading-7 text-white/78">
                <li>Create notes with structure, not a cluttered spreadsheet interface.</li>
                <li>Study from a focused browser app instead of juggling exports and local-only data silos.</li>
                <li>Adopt sharing, larger limits, and teams when the collection actually needs them.</li>
              </ul>
            </article>

            <article
              id="pricing"
              className="rounded-[2rem] border p-8 shadow-[0_24px_70px_rgba(24,35,29,0.08)]"
              style={{borderColor: 'var(--line)', backgroundColor: 'rgba(255,252,247,0.92)'}}
            >
              <p className="text-xs font-semibold uppercase tracking-[0.28em]" style={{color: 'var(--muted)'}}>Pricing direction</p>
              <div className="mt-6 grid gap-4">
                <div className="rounded-3xl border p-5" style={{borderColor: 'var(--line)', backgroundColor: 'rgba(255,252,247,0.72)'}}>
                  <p className="text-sm font-semibold" style={{color: 'var(--ink)'}}>Free</p>
                  <p className="mt-2 text-sm leading-6" style={{color: 'var(--muted)'}}>2 decks, 10 notes, browser study, imports, and OTP account access.</p>
                </div>
                <div className="rounded-3xl border p-5" style={{borderColor: 'rgba(182,95,63,0.35)', backgroundColor: 'var(--accent-soft)'}}>
                  <p className="text-sm font-semibold" style={{color: 'var(--ink)'}}>Pro</p>
                  <p className="mt-2 text-sm leading-6" style={{color: 'var(--muted)'}}>Larger limits, deck sharing, backups, and future cross-device sync entitlements.</p>
                </div>
                <div className="rounded-3xl border p-5" style={{borderColor: 'var(--line)', backgroundColor: 'var(--surface-muted)'}}>
                  <p className="text-sm font-semibold" style={{color: 'var(--ink)'}}>Team</p>
                  <p className="mt-2 text-sm leading-6" style={{color: 'var(--muted)'}}>Organization workspaces, seats, centralized billing, and shared libraries.</p>
                </div>
              </div>
            </article>
          </section>
        </main>
      </div>
    </div>
  )
}
