import { useEffect, useMemo, useState, type ReactNode } from 'react'

type SessionResponse = {
  authenticated: boolean
  user?: {
    displayName?: string
    email: string
  }
}

type IconName =
  | 'sparkles'
  | 'menu'
  | 'arrow-right'
  | 'play'
  | 'note'
  | 'template'
  | 'deck'
  | 'market'
  | 'repeat'
  | 'users'
  | 'chart'
  | 'mail'
  | 'shield'
  | 'upload'
  | 'clock'
  | 'check'
  | 'star'
  | 'crown'
  | 'trophy'
  | 'bell'
  | 'message'
  | 'target'
  | 'brain'
  | 'chevron-down'

type FeatureCard = {
  title: string
  description: string
  icon: IconName
  status?: string
}

type MarketplaceCard = {
  category: string
  title: string
  author: string
  price: string
  installs: string
  rating: string
  premium?: boolean
}

type ScienceCard = {
  stat: string
  label: string
  title: string
  description: string
  icon: IconName
}

type PricingTier = {
  name: string
  price: string
  cadence: string
  description: string
  features: string[]
  cta: string
  featured?: boolean
}

type Testimonial = {
  quote: string
  name: string
  role: string
}

const appOrigin = (
  import.meta.env.VITE_APP_ORIGIN ??
  (import.meta.env.DEV ? 'http://localhost:8000' : 'https://app.vutadex.com')
).replace(/\/$/, '')

const navLinks = [
  {label: 'Features', href: '#features'},
  {label: 'Marketplace', href: '#marketplace'},
  {label: 'Pricing', href: '#pricing'},
  {label: 'FAQ', href: '#faq'},
] as const

const builtFor = [
  'Medical prep',
  'Certification tracks',
  'Law review',
  'Language learning',
  'Engineering teams',
  'Exam cohorts',
]

const featureCards: FeatureCard[] = [
  {
    title: 'Edit Notes',
    description: 'Capture knowledge with deck-aware note creation, tags, template-backed fields, and recent note context while you build.',
    icon: 'note',
  },
  {
    title: 'Card Templates',
    description: 'Move beyond a single flashcard shape with reusable templates, cloze support, and live styling previews.',
    icon: 'template',
  },
  {
    title: 'Deck Management',
    description: 'Organize decks, rename and clean them up safely, and keep the structure of your study system coherent.',
    icon: 'deck',
  },
  {
    title: 'Marketplace',
    description: 'Publish decks with descriptions, creator identity, and pricing so high-quality collections can travel further.',
    icon: 'market',
    status: 'Roadmap',
  },
  {
    title: 'Spaced Repetition',
    description: 'Stay focused on what is actually due instead of manually deciding what to review next.',
    icon: 'repeat',
  },
  {
    title: 'OTP Sign-In',
    description: 'Use one-time codes and secure server-side sessions instead of password resets and local-only auth flows.',
    icon: 'mail',
  },
  {
    title: 'Study Groups',
    description: 'Bring deck sharing, member management, and group workflows into the product as Vutadex expands to teams.',
    icon: 'users',
    status: 'Coming soon',
  },
  {
    title: 'Analytics',
    description: 'Track decks, note volume, due work, and plan usage from a browser-first workspace overview.',
    icon: 'chart',
  },
]

const steps = [
  {
    title: 'Capture notes',
    description: 'Write or import source material and keep it organized by deck, tags, and note type.',
    icon: 'note',
  },
  {
    title: 'Turn it into cards',
    description: 'Generate cards from templates so the note remains the source of truth and cards stay consistent.',
    icon: 'deck',
  },
  {
    title: 'Review with spacing',
    description: 'Study what is due instead of what happens to be in front of you.',
    icon: 'repeat',
  },
  {
    title: 'Grow into groups',
    description: 'Share decks, invite collaborators, and give serious learners a team workflow when they need it.',
    icon: 'users',
  },
  {
    title: 'Track progress',
    description: 'See deck totals, due work, usage against plan limits, and how the collection evolves over time.',
    icon: 'chart',
  },
] as const

const marketplaceCards: MarketplaceCard[] = [
  {
    category: 'Cloud',
    title: 'AWS Solutions Architect Associate',
    author: 'Vutadex Creator',
    price: '$29',
    installs: '8.9k installs',
    rating: '4.8',
    premium: true,
  },
  {
    category: 'Medicine',
    title: 'Upper Limb Anatomy',
    author: 'Shared deck collection',
    price: 'Free',
    installs: '5.4k installs',
    rating: '4.9',
  },
  {
    category: 'Language',
    title: 'Japanese N3 Vocabulary',
    author: 'Community author',
    price: '$19',
    installs: '12.1k installs',
    rating: '4.9',
    premium: true,
  },
]

const scienceCards: ScienceCard[] = [
  {
    stat: 'Less guesswork',
    label: 'Deck-driven review',
    title: 'Spaced repetition',
    description: 'Review timing stays attached to cards so the next session is obvious instead of manual.',
    icon: 'repeat',
  },
  {
    stat: 'High signal',
    label: 'Prompted retrieval',
    title: 'Active recall',
    description: 'Templates force a question-and-answer discipline that turns notes into actual recall practice.',
    icon: 'brain',
  },
  {
    stat: 'Fast loop',
    label: 'Create, edit, study',
    title: 'Tight feedback cycle',
    description: 'Recent notes, note editing, and deck context shorten the gap between making cards and refining them.',
    icon: 'clock',
  },
  {
    stat: 'Built to scale',
    label: 'Free to enterprise',
    title: 'System progression',
    description: 'Start simple, then add larger limits, sharing, and team structures as the collection grows.',
    icon: 'shield',
  },
]

const studyGroupBenefits = [
  {
    title: 'Shared decks',
    description: 'Pool contributions into a common collection instead of copying files around.',
    icon: 'deck',
  },
  {
    title: 'Leaderboards',
    description: 'Turn accountability into a visible habit, not a vague intention.',
    icon: 'trophy',
  },
  {
    title: 'Reminders',
    description: 'Keep the group moving with nudges, streaks, and lightweight coordination.',
    icon: 'bell',
  },
  {
    title: 'Discussion',
    description: 'Tie difficult concepts back to the deck and clarify them where the study work lives.',
    icon: 'message',
  },
] as const

const pricingTiers: PricingTier[] = [
  {
    name: 'Free',
    price: '$0',
    cadence: 'forever',
    description: 'For solo learners getting started inside the browser app.',
    features: ['2 decks', '10 notes', '100 total cards', 'OTP account access', 'Import/export basics'],
    cta: 'Start Free',
  },
  {
    name: 'Pro',
    price: '$12',
    cadence: '/month',
    description: 'For serious learners who need larger limits and publishing paths.',
    features: ['100 decks', '50,000 notes', 'Deck sharing', 'Backups', 'Marketplace publishing'],
    cta: 'Start Pro',
    featured: true,
  },
  {
    name: 'Team',
    price: '$8',
    cadence: '/user/month',
    description: 'For study groups, shared libraries, and organization workflows with a 3-seat minimum.',
    features: ['Org workspaces', 'Shared libraries', 'Study group creation', 'Member management', '3-seat minimum'],
    cta: 'Start Team',
  },
  {
    name: 'Enterprise',
    price: 'Custom',
    cadence: '',
    description: 'For high-scale rollouts, governance controls, and contract-backed limits.',
    features: ['Custom entitlements', 'SSO-ready path', 'Admin controls', 'Volume planning', 'Priority support'],
    cta: 'Talk to Sales',
  },
]

const testimonials: Testimonial[] = [
  {
    quote: 'Vutadex is the first flashcard workflow I have used that keeps note creation and study structure aligned.',
    name: 'Chris B.',
    role: 'Cloud certification prep',
  },
  {
    quote: 'The note editor, template model, and deck context make it much easier to refine cards after the first pass.',
    name: 'Emma R.',
    role: 'Law exam review',
  },
  {
    quote: 'I like that the product is browser-first now but clearly designed to grow into sync and team use later.',
    name: 'Jordan P.',
    role: 'Technical training lead',
  },
] as const

const faqItems = [
  {
    question: 'What makes Vutadex different from a basic flashcard app?',
    answer:
      'Vutadex is built around note-centric card creation, deck-aware workflows, and a browser-first product path that can grow into marketplace, teams, and sync instead of starting as a one-off local tool.',
  },
  {
    question: 'Can I start for free?',
    answer:
      'Yes. The free tier is intentionally small so you can validate the workflow before committing to larger limits or collaboration features.',
  },
  {
    question: 'Will deck publishing and marketplace support be part of the product?',
    answer:
      'Yes. The marketing direction already accounts for marketplace publishing, creator identity, and free or paid deck distribution as the next product layers roll out.',
  },
  {
    question: 'Are study groups available yet?',
    answer:
      'Study Groups are part of the near-term direction. The marketing site positions them clearly, while the app currently exposes the placeholder route and the team-oriented product path.',
  },
  {
    question: 'How does login work?',
    answer:
      'The app uses one-time passcodes sent by email with secure HttpOnly cookies managed by the server. In local development, the code is exposed directly in the UI.',
  },
] as const

function userInitial(user?: SessionResponse['user']) {
  const label = user?.displayName || user?.email || ''
  return label.trim().charAt(0).toUpperCase() || 'V'
}

function Icon({name, className = 'h-5 w-5'}: {name: IconName; className?: string}) {
  switch (name) {
    case 'sparkles':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M12 3l1.9 5.1L19 10l-5.1 1.9L12 17l-1.9-5.1L5 10l5.1-1.9L12 3Z" />
          <path d="M19 3v4" />
          <path d="M21 5h-4" />
        </svg>
      )
    case 'menu':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M4 6h16" />
          <path d="M4 12h16" />
          <path d="M4 18h16" />
        </svg>
      )
    case 'arrow-right':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M5 12h14" />
          <path d="m12 5 7 7-7 7" />
        </svg>
      )
    case 'play':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M8 6.5 18 12 8 17.5V6.5Z" />
        </svg>
      )
    case 'note':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M7 3h7l5 5v13H7a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2Z" />
          <path d="M14 3v5h5" />
          <path d="M9 13h6" />
          <path d="M9 17h6" />
        </svg>
      )
    case 'template':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <rect x="3" y="4" width="18" height="6" rx="1.5" />
          <rect x="3" y="14" width="10" height="6" rx="1.5" />
          <rect x="16" y="14" width="5" height="6" rx="1.5" />
        </svg>
      )
    case 'deck':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="m12 4 8 3.5-8 3.5L4 7.5 12 4Z" />
          <path d="m4 12 8 3.5 8-3.5" />
          <path d="m4 16.5 8 3.5 8-3.5" />
        </svg>
      )
    case 'market':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M4 8h16l-1.2 11.2a2 2 0 0 1-2 1.8H7.2a2 2 0 0 1-2-1.8L4 8Z" />
          <path d="M9 8V6a3 3 0 0 1 6 0v2" />
        </svg>
      )
    case 'repeat':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="m17 2 4 4-4 4" />
          <path d="M3 11v-1a4 4 0 0 1 4-4h14" />
          <path d="m7 22-4-4 4-4" />
          <path d="M21 13v1a4 4 0 0 1-4 4H3" />
        </svg>
      )
    case 'users':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
          <circle cx="9" cy="7" r="4" />
          <path d="M22 21v-2a4 4 0 0 0-3-3.9" />
          <path d="M16 3.1a4 4 0 0 1 0 7.8" />
        </svg>
      )
    case 'chart':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M4 4v15a1 1 0 0 0 1 1h15" />
          <path d="M8 16V9" />
          <path d="M12 16V6" />
          <path d="M16 16v-4" />
        </svg>
      )
    case 'mail':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <rect x="3" y="5" width="18" height="14" rx="2" />
          <path d="m4 7 8 6 8-6" />
        </svg>
      )
    case 'shield':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="m12 3 7 3v6c0 4.5-2.9 7.7-7 9-4.1-1.3-7-4.5-7-9V6l7-3Z" />
          <path d="m9.5 12 1.7 1.7 3.3-3.7" />
        </svg>
      )
    case 'upload':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M12 16V4" />
          <path d="m7 9 5-5 5 5" />
          <path d="M5 20h14" />
        </svg>
      )
    case 'clock':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <circle cx="12" cy="13" r="8" />
          <path d="M12 9v4l3 2" />
          <path d="M9 2h6" />
        </svg>
      )
    case 'check':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.2" className={className} aria-hidden="true">
          <path d="m5 12 4 4 10-10" />
        </svg>
      )
    case 'star':
      return (
        <svg viewBox="0 0 24 24" fill="currentColor" className={className} aria-hidden="true">
          <path d="m12 2.8 2.9 5.9 6.5.9-4.7 4.6 1.1 6.5-5.8-3.1-5.8 3.1 1.1-6.5L2.6 9.6l6.5-.9L12 2.8Z" />
        </svg>
      )
    case 'crown':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="m4 18 2-10 6 4 6-4 2 10H4Z" />
          <path d="M7 18h10" />
        </svg>
      )
    case 'trophy':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M8 3h8v4a4 4 0 0 1-8 0V3Z" />
          <path d="M6 5H4a2 2 0 0 0 2 4" />
          <path d="M18 5h2a2 2 0 0 1-2 4" />
          <path d="M12 11v4" />
          <path d="M8 21h8" />
          <path d="M10 15h4a2 2 0 0 1 2 2v2H8v-2a2 2 0 0 1 2-2Z" />
        </svg>
      )
    case 'bell':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M6 16h12l-1.5-2.2A6 6 0 0 1 15.5 10V9a3.5 3.5 0 0 0-7 0v1a6 6 0 0 1-1 3.8L6 16Z" />
          <path d="M10 19a2 2 0 0 0 4 0" />
        </svg>
      )
    case 'message':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M5 6h14a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H9l-4 3v-3H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2Z" />
        </svg>
      )
    case 'target':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <circle cx="12" cy="12" r="9" />
          <circle cx="12" cy="12" r="5" />
          <circle cx="12" cy="12" r="1.8" fill="currentColor" stroke="none" />
        </svg>
      )
    case 'brain':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="M12 4v16" />
          <path d="M9 8a3 3 0 1 0 0 8" />
          <path d="M15 8a3 3 0 1 1 0 8" />
          <path d="M9 8a3 3 0 0 1 6 0" />
          <path d="M9 16a3 3 0 0 0 6 0" />
        </svg>
      )
    case 'chevron-down':
      return (
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8" className={className} aria-hidden="true">
          <path d="m6 9 6 6 6-6" />
        </svg>
      )
  }
}

function Logo() {
  return (
    <a href="/" className="flex items-center gap-3">
      <span className="inline-flex h-11 w-11 items-center justify-center rounded-2xl bg-[var(--accent)] text-base font-black text-black shadow-[0_12px_30px_rgba(112,214,108,0.16)]">
        V
      </span>
      <span className="text-2xl font-black tracking-tight text-white">Vutadex</span>
    </a>
  )
}

function ButtonLink({
  href,
  children,
  variant = 'primary',
  className = '',
}: {
  href: string
  children: ReactNode
  variant?: 'primary' | 'secondary' | 'ghost'
  className?: string
}) {
  const variantClass =
    variant === 'primary'
      ? 'hover:-translate-y-0.5 hover:shadow-[0_18px_40px_rgba(255,255,255,0.08)]'
      : variant === 'secondary'
        ? 'border border-[var(--line-strong)] bg-white/5 text-white hover:border-[var(--accent)] hover:bg-white/8'
        : 'text-[var(--muted)] hover:text-white'

  const variantStyle =
    variant === 'primary'
      ? {
          backgroundColor: '#f3f5f2',
          color: '#060806',
        }
      : undefined

  return (
    <a
      href={href}
      className={`inline-flex items-center justify-center gap-2 rounded-xl px-5 py-3 text-sm font-semibold transition ${variantClass} ${className}`}
      style={variantStyle}
    >
      {children}
    </a>
  )
}

function SectionEyebrow({children}: {children: ReactNode}) {
  return <p className="text-sm font-semibold uppercase tracking-[0.32em] text-[var(--accent)]">{children}</p>
}

function SectionHeading({
  eyebrow,
  title,
  description,
  center = false,
}: {
  eyebrow: string
  title: string
  description: string
  center?: boolean
}) {
  return (
    <div className={center ? 'mx-auto max-w-3xl text-center' : 'max-w-3xl'}>
      <SectionEyebrow>{eyebrow}</SectionEyebrow>
      <h2 className="mt-4 text-balance text-4xl font-black tracking-[-0.04em] text-white sm:text-5xl">{title}</h2>
      <p className="mt-4 text-lg leading-8 text-[var(--muted)]">{description}</p>
    </div>
  )
}

export function App() {
  const [session, setSession] = useState<SessionResponse | null>(null)
  const [menuOpen, setMenuOpen] = useState(false)
  const [openFaq, setOpenFaq] = useState<number | null>(0)

  useEffect(() => {
    let cancelled = false

    fetch(`${appOrigin}/api/auth/session`, {
      credentials: 'include',
    })
      .then(async (response) => {
        if (!response.ok) throw new Error('Session request failed')
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

  const primaryHref = session?.authenticated ? `${appOrigin}/decks` : `${appOrigin}/login`
  const authLabel = session?.authenticated ? 'Open app' : 'Start Studying Free'

  const desktopActions = useMemo(() => {
    if (session?.authenticated) {
      return (
        <div className="hidden items-center gap-3 md:flex">
          <a
            href={`${appOrigin}/decks`}
            className="inline-flex h-11 w-11 items-center justify-center rounded-full border border-[var(--line-strong)] bg-white/5 text-sm font-bold text-[var(--accent)] transition hover:border-[var(--accent)] hover:bg-white/8"
            aria-label="Open app"
          >
            {userInitial(session.user)}
          </a>
          <ButtonLink href={`${appOrigin}/decks`} variant="primary">
            Open app
          </ButtonLink>
        </div>
      )
    }

    return (
      <div className="hidden items-center gap-3 md:flex">
        <ButtonLink href={`${appOrigin}/login`} variant="ghost">
          Log in
        </ButtonLink>
        <ButtonLink href={`${appOrigin}/login`} variant="primary">
          Start Studying Free
        </ButtonLink>
      </div>
    )
  }, [session])

  return (
    <div className="min-h-screen bg-[var(--bg)] text-[var(--text)]">
      <header className="fixed inset-x-0 top-0 z-50 border-b border-[var(--line)] bg-[color:rgba(6,8,6,0.84)] backdrop-blur-xl">
        <nav className="mx-auto flex h-20 max-w-7xl items-center justify-between px-6">
          <Logo />

          <div className="hidden items-center gap-10 md:flex">
            {navLinks.map((link) => (
              <a key={link.href} href={link.href} className="text-sm text-[var(--muted)] transition hover:text-white">
                {link.label}
              </a>
            ))}
          </div>

          {desktopActions}

          <button
            type="button"
            onClick={() => setMenuOpen((open) => !open)}
            className="inline-flex h-11 w-11 items-center justify-center rounded-xl border border-[var(--line-strong)] bg-white/5 text-[var(--muted)] transition hover:border-[var(--accent)] hover:text-white md:hidden"
            aria-label="Toggle menu"
            aria-expanded={menuOpen}
          >
            <Icon name="menu" />
          </button>
        </nav>

        {menuOpen ? (
          <div className="border-t border-[var(--line)] px-6 py-5 md:hidden">
            <div className="flex flex-col gap-4">
              {navLinks.map((link) => (
                <a key={link.href} href={link.href} className="text-sm text-[var(--muted)] transition hover:text-white" onClick={() => setMenuOpen(false)}>
                  {link.label}
                </a>
              ))}
            </div>
            <div className="mt-5 flex flex-col gap-3">
              {session?.authenticated ? (
                <>
                  <a
                    href={`${appOrigin}/decks`}
                    className="inline-flex items-center justify-center rounded-xl border border-[var(--line-strong)] bg-white/5 px-4 py-3 text-sm font-semibold text-white"
                  >
                    Open app as {userInitial(session.user)}
                  </a>
                </>
              ) : (
                <>
                  <ButtonLink href={`${appOrigin}/login`} variant="ghost" className="justify-center">
                    Log in
                  </ButtonLink>
                  <ButtonLink href={`${appOrigin}/login`} variant="primary" className="justify-center">
                    Start Studying Free
                  </ButtonLink>
                </>
              )}
            </div>
          </div>
        ) : null}
      </header>

      <main className="pt-20">
        <section className="relative overflow-hidden border-b border-[var(--line)]">
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(112,214,108,0.18),transparent_30%),radial-gradient(circle_at_85%_12%,rgba(112,214,108,0.14),transparent_20%),linear-gradient(180deg,rgba(8,11,8,0.96),rgba(6,8,6,1))]" />
          <div className="relative mx-auto max-w-7xl px-6 py-20 sm:py-24 lg:py-28">
            <div className="grid items-center gap-16 lg:grid-cols-[1.02fr_0.98fr]">
              <div>
                <div className="inline-flex items-center gap-2 rounded-full border border-[var(--line-strong)] bg-white/5 px-4 py-2 text-sm text-[var(--muted)]">
                  <Icon name="sparkles" className="h-4 w-4 text-[var(--accent)]" />
                  Browser-first flashcard workspace
                </div>

                <h1 className="mt-8 text-balance text-5xl font-black tracking-[-0.05em] text-white sm:text-6xl lg:text-7xl">
                  Learn faster.
                  <br />
                  <span className="text-[var(--text-soft)]">Remember longer.</span>
                </h1>

                <p className="mt-8 max-w-2xl text-xl leading-9 text-[var(--muted)]">
                  Turn notes into structured decks, keep recent context visible while you create, and grow into marketplace and team workflows without leaving the browser.
                </p>

                <div className="mt-10 flex flex-col gap-4 sm:flex-row">
                  <ButtonLink href={primaryHref} variant="primary" className="sm:px-6 sm:py-3.5">
                    {authLabel}
                    <Icon name="arrow-right" className="h-4 w-4" />
                  </ButtonLink>
                  <ButtonLink href="#marketplace" variant="secondary" className="sm:px-6 sm:py-3.5">
                    <Icon name="play" className="h-4 w-4" />
                    Explore Marketplace
                  </ButtonLink>
                </div>

                <div className="mt-12 grid gap-6 sm:grid-cols-3">
                  <div>
                    <p className="text-lg font-semibold text-white">Recent-note context</p>
                    <p className="mt-2 text-sm leading-6 text-[var(--muted)]">Keep deck-building grounded by seeing the latest notes while you add more.</p>
                  </div>
                  <div>
                    <p className="text-lg font-semibold text-white">OTP access</p>
                    <p className="mt-2 text-sm leading-6 text-[var(--muted)]">Sign in with secure one-time codes and server-managed sessions instead of passwords.</p>
                  </div>
                  <div>
                    <p className="text-lg font-semibold text-white">Sync-ready path</p>
                    <p className="mt-2 text-sm leading-6 text-[var(--muted)]">Web first now, with Pro sync, desktop, and mobile growth designed into the platform.</p>
                  </div>
                </div>
              </div>

              <div className="relative">
                <div className="rounded-[2rem] border border-[var(--line-strong)] bg-[var(--card)] p-4 shadow-[0_32px_90px_rgba(0,0,0,0.45)]">
                  <div className="rounded-[1.7rem] border border-[var(--line)] bg-[var(--card-strong)] p-6">
                    <div className="mb-6 flex items-center justify-between gap-4">
                      <div className="flex items-center gap-3">
                        <div className="flex h-12 w-12 items-center justify-center rounded-2xl bg-[var(--accent)] text-black">
                          <Icon name="brain" className="h-6 w-6" />
                        </div>
                        <div>
                          <p className="text-lg font-semibold text-white">Study Session</p>
                          <p className="text-sm text-[var(--muted)]">AWS Solutions Architect</p>
                        </div>
                      </div>
                      <div className="inline-flex items-center gap-2 rounded-full bg-[rgba(112,214,108,0.12)] px-3 py-1 text-sm text-[var(--accent)]">
                        <Icon name="users" className="h-4 w-4" />
                        Study Group
                      </div>
                    </div>

                    <div className="overflow-hidden rounded-[1.5rem] border border-[var(--line)] bg-black">
                      <div className="border-b border-[var(--line)] px-6 py-8">
                        <p className="text-center text-2xl font-medium text-white">What is a region in AWS?</p>
                      </div>
                      <div className="bg-[rgba(112,214,108,0.1)] px-6 py-8">
                        <p className="text-center text-3xl font-black text-[var(--accent)]">A geographic area with multiple AZs</p>
                      </div>
                    </div>

                    <div className="mt-6 flex items-center justify-between gap-4 rounded-2xl border border-[var(--line)] bg-black/40 p-5">
                      <div className="flex items-center gap-5">
                        <div>
                          <p className="text-sm text-[var(--muted)]">Cards reviewed</p>
                          <p className="mt-1 text-2xl font-black text-white">24/50</p>
                        </div>
                        <div className="h-12 w-px bg-[var(--line)]" />
                        <div>
                          <p className="text-sm text-[var(--muted)]">Accuracy</p>
                          <p className="mt-1 text-2xl font-black text-[var(--accent)]">92%</p>
                        </div>
                      </div>
                      <div className="flex h-[4.5rem] w-[4.5rem] items-center justify-center rounded-full border-[5px] border-[var(--accent)] bg-[rgba(112,214,108,0.12)] text-lg font-black text-white">
                        48%
                      </div>
                    </div>
                  </div>
                </div>

                <div className="absolute -right-4 top-12 hidden rounded-2xl border border-[var(--line-strong)] bg-[var(--card)] px-5 py-4 shadow-[0_20px_55px_rgba(0,0,0,0.35)] lg:block">
                  <div className="flex items-center gap-3">
                    <div className="flex h-11 w-11 items-center justify-center rounded-2xl bg-[rgba(112,214,108,0.12)] text-[var(--accent)]">
                      <Icon name="sparkles" className="h-5 w-5" />
                    </div>
                    <div>
                      <p className="text-base font-semibold text-white">Template-powered</p>
                      <p className="text-sm text-[var(--muted)]">Cards stay tied to note structure</p>
                    </div>
                  </div>
                </div>

                <div className="absolute -bottom-5 left-8 hidden rounded-2xl border border-[var(--line-strong)] bg-[var(--card)] px-5 py-4 shadow-[0_20px_55px_rgba(0,0,0,0.35)] lg:block">
                  <div className="flex items-center gap-3">
                    <div className="flex -space-x-2">
                      <span className="h-9 w-9 rounded-full bg-[rgba(112,214,108,0.55)] ring-2 ring-[var(--card)]" />
                      <span className="h-9 w-9 rounded-full bg-[rgba(112,214,108,0.38)] ring-2 ring-[var(--card)]" />
                      <span className="h-9 w-9 rounded-full bg-[rgba(112,214,108,0.22)] ring-2 ring-[var(--card)]" />
                    </div>
                    <div>
                      <p className="text-base font-semibold text-white">Team-ready path</p>
                      <p className="text-sm text-[var(--muted)]">Study groups and shared decks</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className="border-b border-[var(--line)] bg-[var(--panel)] py-12">
          <div className="mx-auto max-w-7xl px-6">
            <p className="mb-6 text-center text-sm text-[var(--muted)]">Built for learners working across serious subject domains</p>
            <div className="flex flex-wrap items-center justify-center gap-x-10 gap-y-4">
              {builtFor.map((item) => (
                <span key={item} className="text-lg font-semibold tracking-tight text-[var(--text-soft)] transition hover:text-white">
                  {item}
                </span>
              ))}
            </div>
          </div>
        </section>

        <section id="features" className="border-b border-[var(--line)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <SectionHeading
              eyebrow="Features"
              title="Everything you need to turn notes into durable recall"
              description="The current product already covers structured note creation, template-driven cards, OTP accounts, and deck management, with marketplace and study groups designed into the next layers."
              center
            />

            <div className="mt-16 grid gap-6 sm:grid-cols-2 xl:grid-cols-4">
              {featureCards.map((card) => (
                <article
                  key={card.title}
                  className="group rounded-[1.6rem] border border-[var(--line)] bg-[var(--card)] p-6 transition hover:-translate-y-1 hover:border-[var(--accent)] hover:bg-[var(--card-strong)]"
                >
                  <div className="mb-5 flex items-start justify-between gap-4">
                    <span className="inline-flex h-14 w-14 items-center justify-center rounded-2xl bg-[rgba(112,214,108,0.1)] text-[var(--accent)] transition group-hover:bg-[var(--accent)] group-hover:text-black">
                      <Icon name={card.icon} className="h-7 w-7" />
                    </span>
                    {card.status ? (
                      <span className="rounded-full border border-[var(--line-strong)] bg-white/5 px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-[var(--muted)]">
                        {card.status}
                      </span>
                    ) : null}
                  </div>
                  <h3 className="text-2xl font-bold tracking-tight text-white">{card.title}</h3>
                  <p className="mt-3 text-base leading-8 text-[var(--muted)]">{card.description}</p>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="border-b border-[var(--line)] bg-[var(--panel)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <SectionHeading
              eyebrow="How It Works"
              title="From notes to mastery in five clean steps"
              description="The flow stays simple: capture, template, review, collaborate, and keep refining the collection instead of starting over."
              center
            />

            <div className="mt-16 grid gap-8 md:grid-cols-5">
              {steps.map((step, index) => (
                <div key={step.title} className="relative text-center">
                  {index < steps.length - 1 ? (
                    <div className="absolute left-1/2 top-10 hidden h-px w-full bg-[var(--line)] md:block" />
                  ) : null}
                  <div className="relative mx-auto flex h-20 w-20 items-center justify-center rounded-full border border-[var(--line-strong)] bg-[var(--card)] text-[var(--accent)] shadow-[0_16px_36px_rgba(0,0,0,0.2)]">
                    <Icon name={step.icon} className="h-9 w-9" />
                    <span className="absolute -right-1 -top-1 inline-flex h-9 w-9 items-center justify-center rounded-full bg-[var(--accent)] text-sm font-black text-black">
                      {String(index + 1).padStart(2, '0')}
                    </span>
                  </div>
                  <h3 className="mt-8 text-2xl font-bold tracking-tight text-white">{step.title}</h3>
                  <p className="mt-3 text-base leading-8 text-[var(--muted)]">{step.description}</p>
                </div>
              ))}
            </div>
          </div>
        </section>

        <section id="marketplace" className="border-b border-[var(--line)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <div className="flex flex-col gap-8 md:flex-row md:items-end md:justify-between">
              <SectionHeading
                eyebrow="Marketplace"
                title="A deck marketplace that feels native to the workflow"
                description="Publish free or premium decks with creator identity, descriptions, and install metadata so high-value study collections can be discovered instead of buried."
              />
              <ButtonLink href={primaryHref} variant="secondary">
                Browse the app
                <Icon name="arrow-right" className="h-4 w-4" />
              </ButtonLink>
            </div>

            <div className="mt-12 grid gap-6 lg:grid-cols-3">
              {marketplaceCards.map((deck) => (
                <article key={deck.title} className="rounded-[1.8rem] border border-[var(--line)] bg-[var(--card)] p-6 transition hover:border-[var(--accent)] hover:bg-[var(--card-strong)]">
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <p className="text-sm font-semibold uppercase tracking-[0.18em] text-[var(--accent)]">{deck.category}</p>
                      <h3 className="mt-3 text-3xl font-bold tracking-tight text-white">{deck.title}</h3>
                      <p className="mt-2 text-base text-[var(--muted)]">by {deck.author}</p>
                    </div>
                    {deck.premium ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-[var(--accent)] px-3 py-1 text-sm font-semibold text-black">
                        <Icon name="crown" className="h-4 w-4" />
                        Premium
                      </span>
                    ) : (
                      <span className="rounded-full border border-[var(--line-strong)] bg-white/5 px-3 py-1 text-sm font-semibold text-[var(--text-soft)]">
                        Free
                      </span>
                    )}
                  </div>

                  <div className="mt-6 flex items-center gap-5 text-sm text-[var(--muted)]">
                    <div className="inline-flex items-center gap-2">
                      <Icon name="star" className="h-4 w-4 text-[var(--accent)]" />
                      <span>{deck.rating}</span>
                    </div>
                    <div className="inline-flex items-center gap-2">
                      <Icon name="users" className="h-4 w-4" />
                      <span>{deck.installs}</span>
                    </div>
                  </div>

                  <div className="mt-8 flex items-center justify-between border-t border-[var(--line)] pt-5">
                    <span className="text-4xl font-black tracking-tight text-white">{deck.price}</span>
                    <ButtonLink href={primaryHref} variant="primary" className="px-5 py-2.5">
                      {deck.price === 'Free' ? 'Get deck' : 'Buy Now'}
                    </ButtonLink>
                  </div>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="border-b border-[var(--line)] bg-[var(--panel)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <SectionHeading
              eyebrow="Learning Science"
              title="Built around recall, repetition, and system growth"
              description="The product copy, UI shape, and roadmap all tie back to a simple idea: make memory work easier to start, easier to refine, and easier to sustain."
              center
            />

            <div className="mt-16 grid gap-6 md:grid-cols-2">
              {scienceCards.map((card) => (
                <article key={card.title} className="rounded-[2rem] border border-[var(--line)] bg-[var(--card)] p-8 transition hover:border-[var(--accent)]">
                  <div className="flex items-start justify-between gap-6">
                    <div className="inline-flex h-16 w-16 items-center justify-center rounded-2xl bg-[rgba(112,214,108,0.1)] text-[var(--accent)]">
                      <Icon name={card.icon} className="h-8 w-8" />
                    </div>
                    <div className="text-right">
                      <p className="text-4xl font-black tracking-tight text-[var(--accent)]">{card.stat}</p>
                      <p className="mt-1 text-sm text-[var(--muted)]">{card.label}</p>
                    </div>
                  </div>
                  <h3 className="mt-8 text-3xl font-bold tracking-tight text-white">{card.title}</h3>
                  <p className="mt-4 text-lg leading-8 text-[var(--muted)]">{card.description}</p>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="border-b border-[var(--line)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <div className="grid gap-12 lg:grid-cols-[0.95fr_1.05fr] lg:items-start">
              <div>
                <SectionHeading
                  eyebrow="Study Groups"
                  title="A collaboration layer that fits the deck model"
                  description="Study Groups are the natural extension of Vutadex: shared decks, member roles, accountability, and collaborative study surfaces built around the same collection model."
                />

                <div className="mt-10 grid gap-6 sm:grid-cols-2">
                  {studyGroupBenefits.map((benefit) => (
                    <div key={benefit.title} className="flex gap-4">
                      <span className="mt-1 inline-flex h-12 w-12 shrink-0 items-center justify-center rounded-2xl bg-[rgba(112,214,108,0.1)] text-[var(--accent)]">
                        <Icon name={benefit.icon} className="h-6 w-6" />
                      </span>
                      <div>
                        <h3 className="text-xl font-bold text-white">{benefit.title}</h3>
                        <p className="mt-2 text-base leading-7 text-[var(--muted)]">{benefit.description}</p>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="mt-10">
                  <ButtonLink href={primaryHref} variant="primary">
                    Create a Study Group
                  </ButtonLink>
                </div>
              </div>

              <div className="rounded-[2rem] border border-[var(--line)] bg-[var(--card)] p-6 shadow-[0_24px_60px_rgba(0,0,0,0.32)]">
                <div className="flex items-center justify-between gap-4">
                  <div className="flex items-center gap-4">
                    <div className="flex h-14 w-14 items-center justify-center rounded-2xl bg-[var(--accent)] text-xl font-black text-black">MS</div>
                    <div>
                      <h3 className="text-2xl font-bold text-white">Med School Study Group</h3>
                      <p className="text-base text-[var(--muted)]">12 members • 8 online</p>
                    </div>
                  </div>
                  <ButtonLink href={primaryHref} variant="secondary" className="px-4 py-2.5">
                    Invite
                  </ButtonLink>
                </div>

                <div className="mt-6 rounded-[1.5rem] border border-[var(--line)] bg-black/35 p-5">
                  <h4 className="text-xl font-bold text-white">Weekly leaderboard</h4>
                  <div className="mt-4 space-y-4">
                    {[
                      {rank: '1', name: 'Alex K.', cards: '847 cards', streak: '14 day streak'},
                      {rank: '2', name: 'Sarah M.', cards: '723 cards', streak: '9 day streak'},
                      {rank: '3', name: 'Jordan P.', cards: '651 cards', streak: '7 day streak'},
                    ].map((person) => (
                      <div key={person.name} className="flex items-center justify-between gap-4">
                        <div className="flex items-center gap-3">
                          <span className={`inline-flex h-9 w-9 items-center justify-center rounded-full text-sm font-black ${person.rank === '1' ? 'bg-[var(--accent)] text-black' : 'bg-white/6 text-[var(--text-soft)]'}`}>
                            {person.rank}
                          </span>
                          <span className="text-xl font-semibold text-white">{person.name}</span>
                        </div>
                        <div className="flex items-center gap-5 text-sm text-[var(--muted)]">
                          <span>{person.cards}</span>
                          <span className="text-[var(--accent)]">{person.streak}</span>
                        </div>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="mt-6 rounded-[1.5rem] border border-[var(--line)] bg-black/35 p-5">
                  <h4 className="text-xl font-bold text-white">Shared decks</h4>
                  <div className="mt-4 space-y-4">
                    {[
                      {name: 'Anatomy - Upper Limb', meta: '234 cards • updated 2h ago'},
                      {name: 'Pharmacology Basics', meta: '189 cards • updated 1d ago'},
                    ].map((deck) => (
                      <div key={deck.name} className="flex items-center justify-between rounded-2xl border border-[var(--line)] bg-white/4 px-4 py-4">
                        <div>
                          <p className="text-lg font-semibold text-white">{deck.name}</p>
                          <p className="mt-1 text-sm text-[var(--muted)]">{deck.meta}</p>
                        </div>
                        <ButtonLink href={primaryHref} variant="ghost" className="px-4 py-2">
                          Study
                        </ButtonLink>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section id="pricing" className="border-b border-[var(--line)] bg-[var(--panel)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <SectionHeading
              eyebrow="Pricing"
              title="Simple pricing that maps to how the product grows"
              description="Start small, unlock more when the collection deserves it, and move into organization workflows without switching tools."
              center
            />

            <div className="mt-16 grid gap-8 lg:grid-cols-4">
              {pricingTiers.map((tier) => (
                <article
                  key={tier.name}
                  className={`relative rounded-[2rem] border p-8 ${tier.featured ? 'border-[var(--accent)] bg-[var(--card-strong)] shadow-[0_20px_60px_rgba(112,214,108,0.1)]' : 'border-[var(--line)] bg-[var(--card)]'}`}
                >
                  {tier.featured ? (
                    <div className="absolute -top-4 left-1/2 -translate-x-1/2 rounded-full bg-[var(--accent)] px-4 py-1 text-sm font-bold text-black">
                      Most popular
                    </div>
                  ) : null}

                  <h3 className="text-2xl font-bold text-white">{tier.name}</h3>
                  <div className="mt-5 flex items-baseline gap-2">
                    <span className="text-5xl font-black tracking-tight text-white">{tier.price}</span>
                    {tier.cadence ? <span className="text-xl text-[var(--muted)]">{tier.cadence}</span> : null}
                  </div>
                  <p className="mt-4 text-base leading-7 text-[var(--muted)]">{tier.description}</p>

                  <ul className="mt-8 space-y-4">
                    {tier.features.map((feature) => (
                      <li key={feature} className="flex items-start gap-3 text-base text-[var(--text-soft)]">
                        <span className="mt-0.5 inline-flex h-6 w-6 items-center justify-center rounded-full bg-[rgba(112,214,108,0.1)] text-[var(--accent)]">
                          <Icon name="check" className="h-3.5 w-3.5" />
                        </span>
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>

                  <div className="mt-8">
                    <ButtonLink href={primaryHref} variant={tier.featured ? 'primary' : 'secondary'} className="w-full justify-center">
                      {tier.cta}
                    </ButtonLink>
                  </div>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section className="border-b border-[var(--line)] py-24">
          <div className="mx-auto max-w-7xl px-6">
            <SectionHeading
              eyebrow="Testimonials"
              title="A sharper pitch for people who outgrow generic flashcard tools"
              description="The strongest reactions to Vutadex come from people who want structure, not just another empty card grid."
              center
            />

            <div className="mt-16 grid gap-6 lg:grid-cols-3">
              {testimonials.map((item) => (
                <article key={item.name} className="rounded-[1.8rem] border border-[var(--line)] bg-[var(--card)] p-6 transition hover:border-[var(--accent)]">
                  <div className="mb-5 flex gap-1 text-[var(--accent)]">
                    {Array.from({length: 5}).map((_, index) => (
                      <Icon key={index} name="star" className="h-4 w-4" />
                    ))}
                  </div>
                  <blockquote className="text-lg leading-8 text-[var(--text-soft)]">“{item.quote}”</blockquote>
                  <div className="mt-8 flex items-center gap-4">
                    <div className="inline-flex h-12 w-12 items-center justify-center rounded-full bg-[rgba(112,214,108,0.12)] text-lg font-bold text-[var(--accent)]">
                      {item.name
                        .split(' ')
                        .map((part) => part[0])
                        .join('')
                        .slice(0, 2)}
                    </div>
                    <div>
                      <p className="text-xl font-semibold text-white">{item.name}</p>
                      <p className="text-sm text-[var(--muted)]">{item.role}</p>
                    </div>
                  </div>
                </article>
              ))}
            </div>
          </div>
        </section>

        <section id="faq" className="border-b border-[var(--line)] bg-[var(--panel)] py-24">
          <div className="mx-auto max-w-4xl px-6">
            <SectionHeading
              eyebrow="FAQ"
              title="The questions serious users ask before they commit"
              description="The landing page should make it clear where Vutadex already helps and where the broader platform is heading next."
              center
            />

            <div className="mt-16 overflow-hidden rounded-[2rem] border border-[var(--line)] bg-[var(--card)]">
              {faqItems.map((item, index) => {
                const open = openFaq === index
                return (
                  <div key={item.question} className={index === faqItems.length - 1 ? '' : 'border-b border-[var(--line)]'}>
                    <button
                      type="button"
                      onClick={() => setOpenFaq((current) => (current === index ? null : index))}
                      className="flex w-full items-center justify-between gap-4 px-6 py-5 text-left"
                      aria-expanded={open}
                    >
                      <span className="text-xl font-semibold text-white">{item.question}</span>
                      <Icon
                        name="chevron-down"
                        className={`h-5 w-5 text-[var(--muted)] transition ${open ? 'rotate-180' : ''}`}
                      />
                    </button>
                    {open ? <div className="px-6 pb-6 text-base leading-8 text-[var(--muted)]">{item.answer}</div> : null}
                  </div>
                )
              })}
            </div>
          </div>
        </section>

        <section className="py-24">
          <div className="mx-auto max-w-5xl px-6">
            <div className="relative overflow-hidden rounded-[2.5rem] border border-[var(--line)] bg-[var(--card)] px-8 py-14 text-center shadow-[0_28px_80px_rgba(0,0,0,0.35)] sm:px-14 sm:py-18">
              <div className="absolute inset-0 bg-[radial-gradient(circle_at_center,rgba(112,214,108,0.12),transparent_58%)]" />
              <div className="relative">
                <div className="mx-auto inline-flex h-16 w-16 items-center justify-center rounded-3xl bg-[var(--accent)] text-black">
                  <Icon name="sparkles" className="h-8 w-8" />
                </div>
                <h2 className="mt-8 text-balance text-5xl font-black tracking-[-0.05em] text-white sm:text-6xl">
                  Start building with Vutadex
                </h2>
                <p className="mx-auto mt-6 max-w-3xl text-xl leading-9 text-[var(--muted)]">
                  Launch from the browser today, use the free tier to test the workflow, and grow into publishing, teams, and sync as the product expands.
                </p>
                <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
                  <ButtonLink href={primaryHref} variant="primary" className="sm:px-6 sm:py-3.5">
                    Start Studying Free
                    <Icon name="arrow-right" className="h-4 w-4" />
                  </ButtonLink>
                  <ButtonLink href="#marketplace" variant="secondary" className="sm:px-6 sm:py-3.5">
                    Explore Marketplace
                  </ButtonLink>
                </div>
                <p className="mt-8 text-sm text-[var(--muted)]">Free tier available • OTP login • Browser-first • Team path ready</p>
              </div>
            </div>
          </div>
        </section>
      </main>

      <footer className="border-t border-[var(--line)] bg-[var(--panel)]">
        <div className="mx-auto max-w-7xl px-6 py-16">
          <div className="grid gap-12 md:grid-cols-2 lg:grid-cols-5">
            <div className="lg:col-span-2">
              <Logo />
              <p className="mt-5 max-w-sm text-base leading-8 text-[var(--muted)]">
                Vutadex is the browser-first flashcard workspace for people who want note structure, deck clarity, and a clear path into publishing and team study.
              </p>
            </div>

            <div>
              <h3 className="text-lg font-semibold text-white">Product</h3>
              <ul className="mt-4 space-y-3 text-[var(--muted)]">
                <li><a href="#features" className="transition hover:text-white">Features</a></li>
                <li><a href="#marketplace" className="transition hover:text-white">Marketplace</a></li>
                <li><a href="#pricing" className="transition hover:text-white">Pricing</a></li>
                <li><a href={`${appOrigin}/decks`} className="transition hover:text-white">Open app</a></li>
              </ul>
            </div>

            <div>
              <h3 className="text-lg font-semibold text-white">Platform</h3>
              <ul className="mt-4 space-y-3 text-[var(--muted)]">
                <li><a href={`${appOrigin}/notes/view`} className="transition hover:text-white">Notes</a></li>
                <li><a href={`${appOrigin}/templates`} className="transition hover:text-white">Templates</a></li>
                <li><a href={`${appOrigin}/decks`} className="transition hover:text-white">Decks</a></li>
                <li><a href={`${appOrigin}/study-groups`} className="transition hover:text-white">Study Groups</a></li>
              </ul>
            </div>

            <div>
              <h3 className="text-lg font-semibold text-white">Legal</h3>
              <ul className="mt-4 space-y-3 text-[var(--muted)]">
                <li><a href="#faq" className="transition hover:text-white">FAQ</a></li>
                <li><a href={primaryHref} className="transition hover:text-white">Log in</a></li>
                <li><a href={primaryHref} className="transition hover:text-white">Start free</a></li>
              </ul>
            </div>
          </div>

          <div className="mt-14 flex flex-col items-center justify-between gap-4 border-t border-[var(--line)] pt-8 text-sm text-[var(--muted)] md:flex-row">
            <p>© 2026 Vutadex. All rights reserved.</p>
            <div className="flex items-center gap-6">
              <a href="#pricing" className="transition hover:text-white">Pricing</a>
              <a href="#faq" className="transition hover:text-white">FAQ</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
