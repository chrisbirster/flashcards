# Vutadex Phases

Last updated: 2026-03-22

This file tracks delivery status for the current app roadmap. Update it whenever a phase meaningfully changes state.

Status legend:
- `done`: shipped in the repo
- `in_progress`: active implementation
- `planned`: agreed next step, not started
- `later`: intentionally deferred

## Current Status

| Phase | Status | Summary |
| --- | --- | --- |
| Phase 0 | `done` | Mobile-first app foundation, Vutadex dark/light theme system, mobile shell, dashboard aggregation, and responsive rewrites for core app routes. |
| Phase 1 | `done` | Shared-content / per-user-review-state split for authenticated study, due queues, deck stats, and revlog ownership. |
| Phase 2 | `done` | Study Groups foundation: canonical source deck, published versions, invites, personal installs, fork/update state, and cross-collection installs. |
| Phase 3 | `done` | Marketplace foundation: listing CRUD, publish flow, listing detail pages, and free versioned installs with source attribution. |
| Phase 4 | `later` | Paid marketplace and creator payouts via Stripe Connect Express. |
| Phase 5 | `later` | AI note-to-card generation, analytics, and study protocols. |
| Phase 6 | `later` | Real-time collaborative editing using Hocuspocus/Yjs. |

## Completed

### Phase 0
- Added the mobile-first app shell:
  - top app bar
  - mobile bottom nav
  - More sheet
  - sticky action bars
  - shared page/sheet/surface primitives
- Rebuilt these screens for handheld layouts:
  - `Login`
  - `Home`
  - `Add Note`
  - `Notes`
  - `Study`
  - `Decks`
  - `Templates`
  - `Study Groups` placeholder
- Added `GET /api/dashboard` for a single mobile-friendly dashboard payload.
- Extended deck responses with summary fields used by mobile deck cards.
- Moved the app onto the Vutadex dark/light theme tokens instead of the old generic Tailwind palette.

### Phase 1
- Added per-user review state storage with `card_review_states`.
- Added `user_id` support on `revlog`.
- Kept notes, cards, decks, and templates as shared content.
- Moved authenticated review behavior to user-scoped state for:
  - due cards
  - answering cards
  - flags
  - marked
  - suspended
  - deck stats
  - dashboard due counts
- Added integration coverage proving one user answering a shared card does not change another user’s due queue.

### Phase 2
- Added Study Groups foundation flows:
  - group CRUD
  - invite and join flow
  - role management
  - explicit source version publishing
  - personal installs
  - fresh-copy install updates
  - install removal
  - fork detection
  - lightweight dashboard
- Implemented true cross-collection installs for personal Study Group copies.
- Scoped note type identities by collection and rewired request handlers so deck, note, template, and study flows resolve the active workspace collection correctly.
- Added regression coverage for cross-collection Study Group installs and preserved per-user review isolation after install.

### Phase 3
- Added Marketplace foundation flows:
  - listing CRUD
  - explicit publish flow
  - published listing detail pages
  - creator management surface
  - free install, update, and removal flows
- Added marketplace listing versions so installs point at published source versions instead of mutable raw deck content.
- Reused the Phase 2 copy/install model so marketplace installs create workspace-local copies with private review history.
- Added integration coverage for:
  - Pro-plus publishing entitlement gates
  - draft vs published visibility
  - premium install blocking until Phase 4
  - free cross-collection installs
  - fresh-copy version updates
  - source deck review isolation after install study
- Tightened request-scoped collection loading so cross-collection installs cannot leave stale note/card ID counters in memory.

## Next Up

### Phase 4
- Paid marketplace and creator economy is the next implementation target.
- Current implementation target:
  - Stripe Connect Express onboarding
  - premium listing checkout
  - order and license records
  - platform fee tracking
  - payout bookkeeping
- Explicitly out of scope before Phase 4 starts:
  - real payment capture for premium marketplace listings
  - creator payouts
  - premium install entitlement after checkout

## Notes

- The marketing site can evolve independently, but the app roadmap above is the current implementation priority.
- Every new feature should be specified mobile-first before desktop expansion.
- When a phase is started or completed, update both the table and the detailed section in this file.

## Detailed Phase Definitions

### Phase 0: Mobile-First App Foundation

Status: `done`

Objective:
- Make the main app usable and consistent on phones before expanding product scope.
- Establish the Vutadex dark/light theme system as the base UI language for all future work.

Primary scope:
- Rebuild the app shell mobile-first.
- Add mobile navigation patterns and sheet primitives.
- Remove the old desktop-first layout assumptions from the core routes.
- Add a single dashboard payload for the home view.

Screens in scope:
- `/`
- `/login`
- `/notes/add`
- `/notes/view`
- `/decks`
- `/templates`
- `/study-groups`
- `/study/:deckId`

Core deliverables:
- Mobile top bar and bottom navigation.
- `More` sheet for secondary routes and account actions.
- Sticky action bars for dense forms and study actions.
- Fullscreen mobile sheets for editing-heavy views.
- Vutadex dark/light semantic theme tokens.
- Mobile-safe layout primitives shared across pages.
- `GET /api/dashboard`.
- Deck list payload enriched for mobile summaries.

Technical outcomes:
- No horizontal scrolling on supported app routes.
- Primary actions remain visible and thumb-reachable on phone widths.
- Mobile layout becomes the default implementation target for future work.

Dependencies:
- None. This is the foundational UX tranche.

Exit criteria:
- Core app flows are clean on a narrow viewport.
- Theme system is unified.
- Dashboard no longer requires multi-request fanout for basic home data.

### Phase 1: Shared Content / Per-User Review State

Status: `done`

Objective:
- Separate shared study content from personal scheduling state so collaboration and marketplace installs are possible later.

Problem it solves:
- Before this phase, answering a card mutated shared `cards` review state directly.
- That model makes shared decks impossible because all users would share the same due queue.

Primary scope:
- Keep notes, cards, decks, and templates as shared content.
- Move due/scheduling/review metadata to per-user storage.
- Keep current app behavior intact for single-user flows.

Data model changes:
- Add `card_review_states`.
- Add `user_id` to `revlog`.
- Backfill initial per-user review states from current card state for existing users/workspace owners.

Behavioral changes:
- Authenticated due queues become user-specific.
- Answering a card updates only the active user’s review state.
- Flags, marked, and suspended become user-specific.
- Deck stats and dashboard due counts become user-specific.
- Shared card content remains shared.

Core deliverables:
- Migration for per-user review state.
- Store methods for user-scoped card loading, due selection, stats, and revlog writes.
- API handlers switched to user-scoped study behavior.
- Integration coverage proving one user’s review does not affect another user’s due queue.

Dependencies:
- Phase 0 shell and app routes already in place.

Exit criteria:
- Two authenticated users can study the same shared card content and get separate due behavior.
- Current mobile-first app routes continue working without contract regressions.

### Phase 2: Study Groups

Status: `done`

Objective:
- Deliver safe group-based distribution and accountability without shared review-state side effects.

User-facing goal:
- Let teams publish source-deck updates and invite members to install personal study copies without affecting each other’s review cycles.

Primary scope:
- Replace the placeholder `/study-groups` page with:
  - list view
  - detail view
  - create/edit flows
  - invite/join flows
  - member management
  - install/update flows
  - group dashboard

Core product rules:
- Study Groups are owned by Team or Enterprise workspaces.
- Study Groups are invite-only in this phase.
- Groups are source-deck centric.
- Owners and admins maintain a canonical source deck in the owning workspace.
- Members study their own installed copies, not the source deck directly.
- Updates are opt-in.
- Updating installs a fresh new local copy and keeps the old copy intact.
- Member edits to an installed copy mark that install as `forked`.
- No review-history preservation across source-version updates in this phase.
- Invited users can join regardless of their own plan.
- Roles:
  - `owner`
  - `admin`
  - `member`

Backend scope:
- Add Study Group APIs for CRUD, invites, join, membership changes, version publishing, installs, and install updates.
- Add Study Group data model support for:
  - published source versions
  - member installs
  - fork state
  - audit events
- Add group dashboard payloads focused on lightweight activity and version adoption.
- Enforce plan gating and role permissions.

Frontend scope:
- Mobile-first Study Groups routes and sheets.
- Member list and invite flows.
- Source version and install status cards.
- Join flow with workspace selection.
- Dashboard cards for activity, latest version, and adoption.

Current implementation snapshot:
- Implemented:
  - group CRUD
  - invite and join flow
  - role management
  - explicit source version publishing
  - personal installs
  - fresh-copy install updates
  - install removal
  - fork detection on installed copies
  - lightweight dashboard
  - mobile-first `/study-groups`, `/study-groups/:groupId`, and join flow UI
  - true cross-collection installs between distinct workspace collections

Dependencies:
- Phase 1 is retained and required so group members can install shared content without sharing due queues.

Exit criteria:
- Team/Enterprise workspace can create a group.
- Invited user can join.
- Group membership and role management work.
- Owner/admin can publish a source version.
- Member can install the latest version into a workspace.
- Member can update to a newer version as a fresh copy.
- One member studying an installed group deck does not affect another member’s due queue.
- Editing an installed member copy marks that install as `forked`.

### Phase 3: Marketplace Foundation

Status: `done`

Objective:
- Launch the non-payment marketplace surface so decks can be listed, browsed, published, and installed.

User-facing goal:
- Make expert-made or community decks discoverable and installable, even before paid checkout exists.

Primary scope:
- Add marketplace routes:
  - `/marketplace`
  - `/marketplace/:slug`
  - `/marketplace/publish`
- Add listing metadata:
  - title
  - slug
  - summary
  - long description
  - author / creator identity
  - category
  - tags
  - cover image
  - price mode
  - install count
  - status
- Add explicit marketplace listing versions so installs point at a published source version instead of raw mutable deck content.

Publishing rules:
- Free users can browse and install free listings.
- Pro, Team, and Enterprise users can publish.
- Team and Enterprise can publish under workspace identity for now.
- Premium purchase flows remain deferred to Phase 4, so this tranche only guarantees free listing installs.

Install model:
- Installing a listing creates a workspace-local deck install linked back to the source listing.
- Review history does not copy.
- Installed content keeps source attribution and version metadata.
- Listing updates are versioned.
- If a newer free version exists, users can install the newer version as a fresh copy.
- Marketplace should reuse the source-version and personal-install model introduced in Phase 2 where possible.

Dependencies:
- Phase 1 per-user review state.
- Phase 2 copy/install pipeline and cross-collection support.

Exit criteria:
- Users can browse listings.
- Users can open listing detail pages.
- Eligible users can create and publish listings.
- Free users can install free listings into a workspace.
- Installed marketplace copies keep source attribution and published version metadata.

### Phase 4: Paid Marketplace and Creator Economy

Status: `later`

Objective:
- Add real paid transactions, creator onboarding, payouts, and licensing to the marketplace.

Primary scope:
- Stripe Connect Express onboarding for creators.
- One-time paid deck purchases.
- Platform fee collection.
- License grant and payout bookkeeping.

Planned billing rules:
- Initial rollout:
  - USD only
  - one-time purchases only
  - default platform fee of `15%`

Data model scope:
- Creator accounts
- Orders
- Licenses
- Payout records
- Listing versioning

Dependencies:
- Phase 3 marketplace foundation must exist first.

Exit criteria:
- Creator can onboard.
- Buyer can purchase a listing.
- License is granted exactly once.
- Platform fee and payout records are correct and webhook-safe.

### Phase 5: AI Generation, Analytics, and Study Protocols

Status: `later`

Objective:
- Add AI-assisted note-to-card generation and richer learning/retention analytics without removing user control.

Primary scope:
- AI-assisted generation from notes using structured model outputs.
- Review/accept flow before generated cards are saved.
- Analytics for:
  - streaks
  - due trends
  - weak-topic detection
  - retention behavior
  - study-session metrics
- Study protocol support:
  - Pomodoro
  - focus sessions
  - gap-based session structure

Planned product rules:
- AI generation is assistive, not auto-publishing.
- User approval is required before cards are persisted.
- Likely gated to Pro and above, with Enterprise overrides later.

Dependencies:
- Marketplace and Study Groups do not strictly block this phase.
- It benefits from Phase 1 because analytics should be tied to per-user study state.

Exit criteria:
- Users can generate structured card suggestions from notes.
- Study dashboards exist at user/deck/group levels.
- Focus sessions and analytics events are persisted and queryable.

### Phase 6: Real-Time Collaborative Editing

Status: `later`

Objective:
- Add live multi-user editing for shared/team-owned learning content.

Primary scope:
- Real-time editing for:
  - note fields
  - deck descriptions
  - template text
- Presence and live cursor/member visibility.
- Conflict-safe syncing using CRDT-based collaboration.
- Version history and rollback support.

Planned architecture:
- Separate collaboration service.
- Hocuspocus/Yjs for real-time document syncing.
- Team and Enterprise shared content first.

Important constraint:
- This phase is intentionally late because real-time collaboration raises the bar on permissions, rollback, and content ownership. The shared-content/per-user-review split from Phase 1 is the prerequisite that makes this tractable.

Dependencies:
- Phase 1 is mandatory.
- Phase 2 likely provides the first meaningful collaborative surface.

Exit criteria:
- Two authorized users can edit the same supported document type live.
- Changes converge correctly.
- Unauthorized users are blocked.
- Version history and rollback exist for collaborative surfaces.

## Sequence Rationale

Why this order:
- Phase 0 fixes the app foundation first.
- Phase 1 fixes the data model needed for any multi-user product.
- Phase 2 and Phase 3 expose the first collaborative and distribution surfaces.
- Phase 4 monetizes Marketplace once the catalog/publish flow is stable.
- Phase 5 adds intelligence and deeper analytics after the core product model is solid.
- Phase 6 adds real-time collaboration last because it has the highest implementation and correctness complexity.

## Update Rules

When this document should be updated:
- when a phase changes status
- when a phase gains or loses scope
- when dependencies change
- when a new delivery rule, plan gate, or sequencing decision is made

Preferred update style:
- keep the phase table short
- keep the detailed phase sections as the source of truth
- record shipped work under the matching phase rather than creating ad hoc notes elsewhere
