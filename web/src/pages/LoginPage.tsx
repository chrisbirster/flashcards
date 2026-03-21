import { useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Navigate, useLocation, useNavigate } from 'react-router'
import { useAppRepository } from '#/lib/app-repository'

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
      navigate(redirectTo, {replace: true})
    },
    onError: (error) => {
      setVerifyError(error instanceof Error ? error.message : 'Failed to verify code')
    },
  })

  if (sessionQuery.data?.authenticated) {
    return <Navigate to="/decks" replace />
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
    <div className="min-h-screen bg-[radial-gradient(circle_at_top,#f7f1e3,transparent_35%),linear-gradient(180deg,#f8fafc_0%,#eef2ff_100%)] px-4 py-12">
      <div className="mx-auto grid max-w-6xl gap-8 lg:grid-cols-[1.15fr_0.85fr]">
        <section className="rounded-[2rem] border border-white/60 bg-white/80 p-8 shadow-[0_40px_120px_rgba(15,23,42,0.08)] backdrop-blur sm:p-10">
          <span className="inline-flex rounded-full border border-slate-200 bg-slate-50 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-slate-500">
            Vutadex
          </span>
          <h1 className="mt-6 max-w-xl text-4xl font-black tracking-tight text-slate-950 sm:text-5xl">
            Study faster with a browser-first flashcard workflow.
          </h1>
          <p className="mt-5 max-w-2xl text-lg leading-8 text-slate-600">
            Sign in with a one-time code to manage decks, create notes, and keep your collection ready for future desktop and mobile sync.
          </p>
          <div className="mt-10 grid gap-4 sm:grid-cols-3">
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-5">
              <p className="text-sm font-semibold text-slate-900">Build intentionally</p>
              <p className="mt-2 text-sm leading-6 text-slate-600">Recent notes, duplicate checks, and deck-aware editing keep card creation fast and coherent.</p>
            </div>
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-5">
              <p className="text-sm font-semibold text-slate-900">Learn anywhere</p>
              <p className="mt-2 text-sm leading-6 text-slate-600">The app runs on the web now, with desktop and mobile sync planned as Pro extensions.</p>
            </div>
            <div className="rounded-3xl border border-slate-200 bg-slate-50/80 p-5">
              <p className="text-sm font-semibold text-slate-900">Scale by plan</p>
              <p className="mt-2 text-sm leading-6 text-slate-600">Free is tight by design. Pro and Team unlock sharing, larger limits, and organization workflows.</p>
            </div>
          </div>
        </section>

        <section className="rounded-[2rem] bg-slate-950 p-8 text-white shadow-[0_40px_120px_rgba(15,23,42,0.22)] sm:p-10">
          <div className="mx-auto max-w-md">
            <p className="text-sm font-semibold uppercase tracking-[0.28em] text-amber-200/80">Sign in</p>
            <h2 className="mt-3 text-3xl font-bold">Use a one-time code</h2>
            <p className="mt-3 text-sm leading-6 text-slate-300">
              Enter your email address and we’ll send a 6-digit code. Sessions are stored in secure HttpOnly cookies and stay active for 7 days while you use the app.
            </p>

            <form onSubmit={submitEmail} className="mt-8 space-y-4">
              <label className="block text-sm font-medium text-slate-200" htmlFor="email">
                Email
              </label>
              <input
                id="email"
                type="email"
                autoComplete="email"
                value={email}
                onChange={(event) => setEmail(event.target.value)}
                placeholder="you@company.com"
                className="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-base text-white outline-none transition focus:border-amber-300"
              />
              <button
                type="submit"
                disabled={requestMutation.isPending}
                className="w-full rounded-2xl bg-amber-300 px-4 py-3 text-sm font-semibold text-slate-950 transition hover:bg-amber-200 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {requestMutation.isPending ? 'Sending code...' : 'Send login code'}
              </button>
            </form>

            {requestError && (
              <p className="mt-4 rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-3 text-sm text-rose-200">{requestError}</p>
            )}
            {requestMessage && (
              <p className="mt-4 rounded-2xl border border-emerald-400/30 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-200">{requestMessage}</p>
            )}

            {requestedEmail && (
              <form onSubmit={submitCode} className="mt-8 space-y-4 border-t border-slate-800 pt-8">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <p className="text-sm font-medium text-slate-200">Verification code</p>
                    <p className="text-xs text-slate-400">{requestedEmail}</p>
                  </div>
                  <button
                    type="button"
                    onClick={() => requestMutation.mutate(requestedEmail)}
                    className="text-xs font-semibold uppercase tracking-[0.24em] text-amber-200 transition hover:text-amber-100"
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
                  className="w-full rounded-2xl border border-slate-700 bg-slate-900 px-4 py-3 text-center text-2xl tracking-[0.45em] text-white outline-none transition focus:border-amber-300"
                />
                {devCode && (
                  <p className="rounded-2xl border border-emerald-400/30 bg-emerald-400/10 px-4 py-3 text-sm text-emerald-200">
                    Development code: <span className="font-mono text-base font-semibold">{devCode}</span>
                  </p>
                )}
                <button
                  type="submit"
                  disabled={verifyMutation.isPending}
                  className="w-full rounded-2xl border border-white/15 bg-white px-4 py-3 text-sm font-semibold text-slate-950 transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {verifyMutation.isPending ? 'Verifying...' : 'Sign in'}
                </button>
                {verifyError && (
                  <p className="rounded-2xl border border-rose-400/30 bg-rose-400/10 px-4 py-3 text-sm text-rose-200">{verifyError}</p>
                )}
                {expiresAt && <p className="text-xs text-slate-400">Code expires at {formatExpiresAt(expiresAt)}.</p>}
              </form>
            )}

            {sessionQuery.isLoading && <p className="mt-6 text-sm text-slate-400">Checking your session...</p>}
          </div>
        </section>
      </div>
    </div>
  )
}
