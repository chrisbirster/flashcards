import { useEffect, useMemo, useState } from 'react'

type SessionResponse = {
  authenticated: boolean
  user?: {
    displayName?: string
    email: string
  }
}

const appOrigin = (import.meta.env.VITE_APP_ORIGIN ?? 'https://app.vutadex.com').replace(/\/$/, '')

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
          className="inline-flex items-center gap-3 rounded-full border border-slate-200 bg-white/85 px-3 py-2 text-sm font-semibold text-slate-900 shadow-sm backdrop-blur transition hover:-translate-y-0.5"
        >
          <span className="inline-flex h-9 w-9 items-center justify-center rounded-full bg-slate-950 text-white">
            {userInitial(session.user)}
          </span>
          <span className="hidden sm:inline">Open app</span>
        </a>
      )
    }

    return (
      <a
        href={`${appOrigin}/login`}
        className="inline-flex items-center rounded-full bg-slate-950 px-5 py-3 text-sm font-semibold text-white transition hover:bg-slate-800"
      >
        Log in
      </a>
    )
  }, [session])

  return (
    <div className="min-h-screen bg-[radial-gradient(circle_at_top,#fff4d6,transparent_28%),linear-gradient(180deg,#f6f0e8_0%,#f8fafc_100%)]">
      <div className="mx-auto max-w-7xl px-5 py-6 sm:px-8">
        <header className="flex items-center justify-between gap-4">
          <a href="/" className="flex items-center gap-3">
            <span className="inline-flex h-11 w-11 items-center justify-center rounded-2xl bg-slate-950 text-base font-black text-white">
              V
            </span>
            <div>
              <p className="text-lg font-black tracking-tight text-slate-950">Vutadex</p>
              <p className="text-xs uppercase tracking-[0.26em] text-slate-500">Flashcards, refined</p>
            </div>
          </a>
          {topRight}
        </header>

        <main className="mt-10 grid gap-8 lg:grid-cols-[1.1fr_0.9fr]">
          <section className="rounded-[2.5rem] border border-white/60 bg-white/75 p-8 shadow-[0_40px_120px_rgba(15,23,42,0.08)] backdrop-blur sm:p-12">
            <p className="text-sm font-semibold uppercase tracking-[0.34em] text-amber-700">The serious flashcard workspace</p>
            <h1 className="mt-5 max-w-3xl text-5xl font-black tracking-tight text-slate-950 sm:text-7xl">
              Build decks with momentum. Study with less friction.
            </h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-600">
              Vutadex gives you a clean note-building workflow, deck-level structure, and a browser-first foundation that is ready for desktop and mobile sync next.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <a
                href={`${appOrigin}/login`}
                className="inline-flex items-center rounded-full bg-slate-950 px-6 py-3 text-sm font-semibold text-white transition hover:bg-slate-800"
              >
                Start in the app
              </a>
              <a
                href="#pricing"
                className="inline-flex items-center rounded-full border border-slate-300 bg-white px-6 py-3 text-sm font-semibold text-slate-900 transition hover:border-slate-400"
              >
                View plans
              </a>
            </div>

            <div className="mt-12 grid gap-4 sm:grid-cols-3">
              <article className="rounded-3xl bg-slate-950 p-5 text-white">
                <p className="text-sm font-semibold text-amber-200">Recent-note context</p>
                <p className="mt-2 text-sm leading-6 text-slate-300">Keep note creation grounded by seeing the latest items in the deck while you build.</p>
              </article>
              <article className="rounded-3xl border border-slate-200 bg-slate-50/90 p-5">
                <p className="text-sm font-semibold text-slate-900">OTP sign-in</p>
                <p className="mt-2 text-sm leading-6 text-slate-600">Fast account access without a password reset maze or brittle local-only sessions.</p>
              </article>
              <article className="rounded-3xl border border-slate-200 bg-slate-50/90 p-5">
                <p className="text-sm font-semibold text-slate-900">Sync-ready foundation</p>
                <p className="mt-2 text-sm leading-6 text-slate-600">Web today, with desktop and mobile sync designed as a Pro extension rather than an afterthought.</p>
              </article>
            </div>
          </section>

          <section className="grid gap-5">
            <article className="rounded-[2rem] bg-[#0f172a] p-8 text-white shadow-[0_35px_100px_rgba(15,23,42,0.22)]">
              <p className="text-xs font-semibold uppercase tracking-[0.28em] text-amber-200/80">Why Vutadex</p>
              <ul className="mt-6 space-y-5 text-sm leading-7 text-slate-300">
                <li>Create notes with structure, not a cluttered spreadsheet interface.</li>
                <li>Study from a focused browser app instead of juggling exports and local-only data silos.</li>
                <li>Adopt sharing, larger limits, and teams when the collection actually needs them.</li>
              </ul>
            </article>

            <article id="pricing" className="rounded-[2rem] border border-slate-200 bg-white/90 p-8 shadow-[0_24px_80px_rgba(15,23,42,0.08)]">
              <p className="text-xs font-semibold uppercase tracking-[0.28em] text-slate-500">Pricing direction</p>
              <div className="mt-6 grid gap-4">
                <div className="rounded-3xl border border-slate-200 p-5">
                  <p className="text-sm font-semibold text-slate-900">Free</p>
                  <p className="mt-2 text-sm leading-6 text-slate-600">2 decks, 10 notes, browser study, imports, and OTP account access.</p>
                </div>
                <div className="rounded-3xl border border-amber-300 bg-amber-50/80 p-5">
                  <p className="text-sm font-semibold text-slate-900">Pro</p>
                  <p className="mt-2 text-sm leading-6 text-slate-600">Larger limits, deck sharing, backups, and future cross-device sync entitlements.</p>
                </div>
                <div className="rounded-3xl border border-slate-200 p-5">
                  <p className="text-sm font-semibold text-slate-900">Team</p>
                  <p className="mt-2 text-sm leading-6 text-slate-600">Organization workspaces, seats, centralized billing, and shared libraries.</p>
                </div>
              </div>
            </article>
          </section>
        </main>
      </div>
    </div>
  )
}
