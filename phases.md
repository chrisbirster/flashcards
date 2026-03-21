# Vutadex Phases

Last updated: 2026-03-21

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
| Phase 2 | `planned` | Real Study Groups CRUD, membership, invites, dashboards, and Team/Enterprise gating. |
| Phase 3 | `planned` | Marketplace foundation: listings, publish flow, install flow, and listing detail pages. |
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

## Next Up

### Phase 2
- Replace the Study Groups placeholder with real routes and APIs.
- Add:
  - group list/detail/create/edit
  - invites and join flow
  - member management
  - roles: `owner`, `admin`, `member`
  - group dashboard with shared decks, activity, and due-card totals
- Gate creation and management to Team and Enterprise.

## Notes

- The marketing site can evolve independently, but the app roadmap above is the current implementation priority.
- Every new feature should be specified mobile-first before desktop expansion.
- When a phase is started or completed, update both the table and the detailed section in this file.
