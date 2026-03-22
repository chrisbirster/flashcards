import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Navigate, useLocation, useNavigate } from 'react-router'
import { useAppRepository } from '#/lib/app-repository'
import { ThemeToggle } from '#/components/theme-toggle'

function formatExpiresAt(value?: string) {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return date.toLocaleTimeString([], {hour: 'numeric', minute: '2-digit'})
}

export function LoginPage() {
  const repository = useAppRepository()
  const queryClient = useQueryClient()
  const navigate = useNavigate()
  const location = useLocation()
  const [email, setEmail] = useState('')
  const [requestedEmail, setRequestedEmail] = useState('')
  const [code, setCode] = useState('')
  const [expiresAt, setExpiresAt] = useState('')
  const [devCode, setDevCode] = useState('')
  const [requestError, setRequestError] = useState('')
  const [verifyError, setVerifyError] = useState('')
  const [requestMessage, setRequestMessage] = useState('')

  const redirectTo = useMemo(() => {
    const state = location.state as {from?: string} | null
    return state?.from || '/decks'
  }, [location.state])

  const sessionQuery = useQuery({
    queryKey: ['auth-session'],
    queryFn: () => repository.fetchSession(),
  })

  const requestMutation = useMutation({
    mutationFn: (address: string) => repository.requestOTP(address),
    onSuccess: (response, address) => {
      const normalized = address.trim().toLowerCase()
      setRequestedEmail(normalized)
      setCode('')
      setExpiresAt(response.expiresAt)
      setDevCode(response.devCode || '')
      setRequestError('')
      setVerifyError('')
      if (response.devCode) {
        setRequestMessage(`Development mode: use the code below for ${normalized}. It expires at ${formatExpiresAt(response.expiresAt)}.`)
      } else {
        setRequestMessage(`Code sent to ${normalized}. It expires at ${formatExpiresAt(response.expiresAt)}.`)
      }
    },
    onError: (error) => {
      setDevCode('')
      setRequestMessage('')
      setRequestError(error instanceof Error ? error.message : 'Failed to send code')
    },
  })

  const verifyMutation = useMutation({
    mutationFn: ({email, code}: {email: string; code: string}) => repository.verifyOTP(email, code),
    onSuccess: async (session) => {
      queryClient.setQueryData(['auth-session'], session)
      await queryClient.invalidateQueries({queryKey: ['auth-session']})
      navigate(session.user?.onboarding ? '/onboarding/plan' : redirectTo, {replace: true})
    },
    onError: (error) => {
      setVerifyError(error instanceof Error ? error.message : 'Failed to verify code')
    },
  })

  if (sessionQuery.data?.authenticated) {
    return <Navigate to={sessionQuery.data.user?.onboarding ? '/onboarding/plan' : redirectTo} replace />
  }

  const activeEmail = requestedEmail || email

  function submitEmail(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setRequestError('')
    setVerifyError('')
    requestMutation.mutate(email)
  }

  function submitCode(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setVerifyError('')
    verifyMutation.mutate({email: activeEmail, code})
  }

  return (
    <div className="app-auth-shell min-h-screen px-4 py-12">
      <div className="mx-auto mb-6 flex max-w-6xl justify-end">
        <ThemeToggle />
      </div>
      <div className="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[1.15fr_0.85fr] lg:gap-8">
        <section className="order-2 rounded-[2rem] border border-[var(--app-line-strong)] bg-[color:var(--app-panel)]/92 p-6 shadow-[0_40px_120px_rgba(0,0,0,0.16)] backdrop-blur sm:p-8 lg:order-1 lg:p-10">
          <span className="inline-flex rounded-full border border-[var(--app-line-strong)] bg-[var(--app-muted-surface)] px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-[var(--app-muted)]">
            Vutadex
          </span>
          <h1 className="mt-5 max-w-xl text-3xl font-black tracking-tight text-[var(--app-text)] sm:text-4xl lg:text-5xl">
            Study faster with a browser-first flashcard workflow.
          </h1>
          <p className="mt-4 max-w-2xl text-base leading-7 text-[var(--app-text-soft)] sm:text-lg sm:leading-8">
            Sign in with a one-time code to manage decks, create notes, and keep your collection ready for future desktop and mobile sync.
          </p>
          <div className="mt-8 grid gap-4 sm:grid-cols-3">
            <div className="rounded-3xl border border-[var(--app-line)] bg-[var(--app-card)]/90 p-5">
              <p className="text-sm font-semibold text-[var(--app-text)]">Build intentionally</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">Recent notes, duplicate checks, and deck-aware editing keep card creation fast and coherent.</p>
            </div>
            <div className="rounded-3xl border border-[var(--app-line)] bg-[var(--app-card)]/90 p-5">
              <p className="text-sm font-semibold text-[var(--app-text)]">Learn anywhere</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">The app runs on the web now, with desktop and mobile sync planned as Pro extensions.</p>
            </div>
            <div className="rounded-3xl border border-[var(--app-line)] bg-[var(--app-card)]/90 p-5">
              <p className="text-sm font-semibold text-[var(--app-text)]">Scale by plan</p>
              <p className="mt-2 text-sm leading-6 text-[var(--app-text-soft)]">Free is tight by design. Pro and Team unlock sharing, larger limits, and organization workflows.</p>
            </div>
          </div>
        </section>

        <section className="order-1 rounded-[2rem] bg-[var(--app-card-strong)] p-6 text-[var(--app-text)] shadow-[0_40px_120px_rgba(0,0,0,0.28)] sm:p-8 lg:order-2 lg:p-10">
          <div className="mx-auto max-w-md">
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-[var(--app-accent)]">Sign in</p>
            <h2 className="mt-3 text-2xl font-bold sm:text-3xl">Use a one-time code</h2>
            <p className="mt-3 text-sm leading-6 text-[var(--app-text-soft)]">
              Enter your email address and we’ll send a 6-digit code. Sessions are stored in secure HttpOnly cookies and stay active for 7 days while you use the app.
            </p>

            <form onSubmit={submitEmail} className="mt-8 space-y-4">
              <label className="block text-sm font-medium text-[var(--app-text-soft)]" htmlFor="email">
                Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="you@company.com"
                className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-bg)] px-4 py-3 text-base text-[var(--app-text)] outline-none transition focus:border-[var(--app-accent)]"
              />
              <button
                type="submit"
                disabled={requestMutation.isPending}
                className="w-full rounded-2xl bg-[var(--app-accent)] px-4 py-3 text-sm font-semibold text-[var(--app-accent-ink)] transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {requestMutation.isPending ? 'Sending code...' : 'Send login code'}
              </button>
            </form>

            {requestError && (
              <p className="mt-4 rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 py-3 text-sm text-[var(--app-danger-text)]">{requestError}</p>
            )}
            {requestMessage && (
              <p className="mt-4 rounded-2xl border border-[var(--app-success-line)] bg-[var(--app-success-surface)] px-4 py-3 text-sm text-[var(--app-success-text)]">{requestMessage}</p>
            )}

            {requestedEmail && (
              <form onSubmit={submitCode} className="mt-8 space-y-4 border-t border-[var(--app-line)] pt-8">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium text-[var(--app-text-soft)]">Verification code</p>
                    <p className="text-xs text-[var(--app-muted)]">{requestedEmail}</p>
                  </div>
                  <button
                    type="button"
                    onClick={() => requestMutation.mutate(requestedEmail)}
                    className="text-xs font-semibold uppercase tracking-[0.24em] text-[var(--app-accent)] transition hover:brightness-110"
                  >
                    Resend
                  </button>
                </div>

                <input
                  inputMode="numeric"
                  autoComplete="one-time-code"
                  value={code}
                  onChange={(event) => setCode(event.target.value.replace(/\D/g, '').slice(0, 6))}
                  placeholder="123456"
                  className="w-full rounded-2xl border border-[var(--app-line-strong)] bg-[var(--app-bg)] px-4 py-3 text-center text-2xl tracking-[0.45em] text-[var(--app-text)] outline-none transition focus:border-[var(--app-accent)]"
                />
                {devCode && (
                  <p className="rounded-2xl border border-[var(--app-success-line)] bg-[var(--app-success-surface)] px-4 py-3 text-sm text-[var(--app-success-text)]">
                    Development code: <span className="font-mono text-base font-semibold">{devCode}</span>
                  </p>
                )}
                <button
                  type="submit"
                  disabled={verifyMutation.isPending}
                  className="w-full rounded-2xl bg-[var(--app-accent)] px-4 py-3 text-sm font-semibold text-[var(--app-accent-ink)] transition hover:brightness-105 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {verifyMutation.isPending ? 'Verifying...' : 'Sign in'}
                </button>
                {verifyError && (
                  <p className="rounded-2xl border border-[var(--app-danger-line)] bg-[var(--app-danger-surface)] px-4 py-3 text-sm text-[var(--app-danger-text)]">{verifyError}</p>
                )}
                {expiresAt && <p className="text-xs text-[var(--app-muted)]">Code expires at {formatExpiresAt(expiresAt)}.</p>}
              </form>
            )}

            {sessionQuery.isLoading && <p className="mt-6 text-sm text-[var(--app-muted)]">Checking your session...</p>}
          </div>
        </section>
      </div>
    </div>
  )
}
