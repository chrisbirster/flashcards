# Microdote — Anki Parity Plan (Tasks + Acceptance Criteria)
_Based on the Anki Manual (print edition)_

---

## 📊 Implementation Progress Summary

| Milestone | Status | Tasks Complete | Notes |
|-----------|--------|----------------|-------|
| **M0** Product Skeleton | ✅ Complete | 4/4 | SQLite storage, profiles, backups |
| **M1** Studying MVP | ✅ Complete | 11/11 | FSRS scheduling, study UI, keyboard shortcuts, daily limits + backlog ordering |
| **M2** Adding/Editing | 🚧 In Progress | 11/17 | Add note + duplicate + cloze flows; field/template editors; conditional/deck override generation |
| **M3** Browser | ⬜ Not Started | 0/11 | — |
| **M4** Deck Options | ⬜ Not Started | 0/8 | — |
| **M5** Statistics | ⬜ Not Started | 0/4 | — |
| **M6** Sync + Safety | ⬜ Not Started | 0/9 | — |
| **M7** Import/Export | ⬜ Not Started | 0/8 | — |
| **M8** Preferences | ⬜ Not Started | 0/19 | — |

**Test Coverage (repo scan):** 125 tests defined (33 backend + 92 E2E)

## 🔍 Repo-Validated Snapshot (2026-02-07)

This plan was cross-checked against current code and tests:

- `Task 0113` is implemented end-to-end (UI timer -> API -> `revlog.time_taken_ms`) and is marked complete.
- `Task 0125` now enforces daily new/review caps with review-backlog priority in due-card selection.
- M2 status markers are synchronized between milestone and flattened list.
- `Task 0203` is partial: tags/flags/marked persist, but browser-queryability depends on M3 search/browser work.
- `Task 0213` and `Task 0224` are partial because browser-mode sorting/rendering is blocked on M3.
- `Task 0214` now includes per-field font/size/RTL plus HTML-editor-default option wiring in add-note editing.
- `Task 0233` now includes cloze-template warning plus auto-fix path.

## 🎯 Next Sprint Checklist (Recommended Order)

### Sprint A — Start M3 Foundation (Unblocks Remaining M2 Criteria)
- [ ] **Task 0301:** Browser shell (sidebar/table/editor skeleton).
- [ ] **Task 0302:** Cards vs Notes mode.
- [ ] **Task 0311:** Search parser + execution (minimum viable operators: text/tag/deck/state).
- [ ] **Task 0321:** Column configuration; include sort-field column path for `Task 0213`.

### Sprint B — M2 Dependency Closures
- [ ] Complete browser-queryability aspect of `Task 0203` (tags/flags/marked filters).
- [ ] Complete browser sorting/display acceptance for `Task 0213`.
- [ ] Complete browser appearance usage path for `Task 0224`.

---

This document translates the **Anki Manual** into an implementation plan for **Microdote** (a web-based reimplementation of Anki functionality).
It is organized as:

- **Milestones** → sets of shippable increments
- Within each milestone: **Features** (human-readable “what”)
- Under each feature: **Tasks** (implementable “how”) with **Acceptance Criteria** and a clear **Definition of Done**
- An **Appendix: Page Coverage Index** that maps the manual’s page ranges to the features below, so *all pages are accounted for*.

> **Notes**
> - “Page” numbers below refer to the **print.html PDF page numbers**.
> - Some Anki desktop features rely on OS-level integrations; Microdote equivalents may be web-friendly variants, but must preserve *user-visible behavior*.

---

## Global Definitions

### Acceptance Criteria
Concrete, testable conditions that must be true for the task/feature to be considered complete.

### Definition of Done (DoD)
A task is “Done” when ALL of these are true:

1. **User flow works end-to-end** in the UI (or API) for the described scenario.
2. **Automated tests** exist for the critical rules/edge cases (unit + integration where relevant).
3. **Data persisted** correctly (reload-safe, session-safe).
4. **Error states** are handled (validation + clear messages).
5. **Performance baseline**: common operations are responsive (no obvious UI jank on typical decks).
6. **Telemetry hooks** (optional) are present for key actions (study answers, sync, import/export) to support later analytics.
7. **Docs updated** (short user-facing help + developer notes for tricky logic).

---

# Milestones

## Milestone M0 — Product Skeleton + Core Data Model ✅ COMPLETE
Goal: Microdote can store decks/notes/cards and render a minimal UI shell.

### Feature 001: Accounts, profiles, and collections
**Summary:** A “collection” is the user’s universe of decks/notes/cards/options. Support multiple profiles (optional in v1).
**Manual pages:** 136–137 (Profiles), plus overall collection concept around 42–43.

**Tasks**
- **Task 0001: Define canonical data model** ✅
  - **Description:** Define DB schema / types for: Collection, Deck, NoteType, Field, Template/CardType, Note, Card, Revlog, Tags, Flags, MediaRefs, Config/Preferences, Sync metadata (USN-like versioning).
  - **Acceptance Criteria:**
    - Schema supports: notes generating multiple cards, deck hierarchy, per-deck options presets, revlog entries with button pressed + time spent.
    - A note can change note type (migration path defined).
    - Cards have stable IDs; revlog references card IDs.
- **Task 0002: Implement collection storage layer** ✅
  - **Description:** Create persistence abstraction (IndexedDB/SQLite-in-browser/remote DB) with transactions and migration system.
  - **Acceptance Criteria:**
    - CRUD for core entities works; migrations can upgrade schema without data loss.
    - Can load a "default collection" on first run.
- **Task 0003: Profile support (phaseable)** ✅
  - **Description:** Allow separate "profiles" with separate collections + preferences; in v1 can be "single profile" but structure must not block multi-profile later.
  - **Acceptance Criteria:**
    - If enabled: switching profiles shows different decks/notes.
    - Only one profile may be connected to one sync account (if implementing Anki-like rule).
- **Task 0004: Backups & restore hooks (foundation)** ✅
  - **Description:** Add local export snapshot capability; wire up UI stubs for "Create Backup" and "Restore Backup".
  - **Acceptance Criteria:**
    - User can create a backup artifact and later restore it.
    - Restore disables auto-sync until user confirms (mirrors safety behavior).

---

## Milestone M1 — Studying MVP (Review Loop) ✅ COMPLETE
Goal: A user can study a deck: show question → reveal answer → grade (Again/Hard/Good/Easy) → scheduling updates.

### Feature 010: Deck list + deck overview
**Summary:** Users choose a deck; overview shows due counts (New/Learning/Review) and a “Study Now” entry point.
**Manual pages:** 48–49.

**Tasks**
- **Task 0101: Deck list UI** ✅
  - **Acceptance Criteria:**
    - Shows all decks (including nested decks).
    - Selecting a deck opens deck overview for that deck.
- **Task 0102: Deck overview due counts** ✅
  - **Acceptance Criteria:**
    - Correctly displays counts split into New/Learning/To Review.
    - If bury siblings is on, show buried count indicator (or equivalent).
- **Task 0103: Study session start/stop** ✅
  - **Acceptance Criteria:**
    - "Study Now" starts a session and continues until daily queue empty.
    - User can return to overview at any time.

### Feature 011: Review screen + answer buttons
**Summary:** Show question first; reveal answer; let user choose Again/Hard/Good/Easy; keyboard shortcuts.
**Manual pages:** 49–50, 53–54.

**Tasks**
- **Task 0111: Question/answer rendering** ✅
  - **Acceptance Criteria:**
    - Question shows first; "Show Answer" reveals answer.
    - Works for HTML templates with CSS styling.
- **Task 0112: Answer buttons + shortcuts** ✅
  - **Acceptance Criteria:**
    - Buttons: Again/Hard/Good/Easy displayed after reveal.
    - Keyboard: Space/Enter = show answer, then "Good"; 1–4 select buttons.
- **Task 0113: Time spent tracking** ✅
  - **Acceptance Criteria:**
    - Record time from question shown to answer selection (ms resolution).
    - Persist in revlog for stats.

### Feature 012: Scheduler v1 (FSRS style with learning steps)
**Summary:** Implement learning/relearning/review states, learning steps, ease factor, lapses.
**Manual pages:** Review behavior described around 49–55 and deck options/scheduler details around 120–130.

**Tasks**
- **Task 0121: Scheduling state machine** ✅ (via FSRS library)
  - **Acceptance Criteria:**
    - Cards exist in: New, Learning, Review, Relearning, Suspended, Buried, Filtered.
    - Transitions correctly on answer buttons.
- **Task 0122: Learning steps + graduating interval** ✅ (via FSRS library)
  - **Acceptance Criteria:**
    - Learning steps schedule in minutes/hours; graduation interval in days.
    - "Again" on learning repeats step; "Good/Easy" advance appropriately.
- **Task 0123: Ease factor + interval updates** ✅ (via FSRS library)
  - **Acceptance Criteria:**
    - Review cards have ease factor; interval changes follow configured multipliers.
    - Enforces minimum interval growth rule ("new interval at least 1 day longer than previous" behavior note).
- **Task 0124: Lapses + leeches** ✅ (via FSRS library)
  - **Acceptance Criteria:**
    - Track lapses; tag as leech at threshold; optional auto-suspend.
    - Leech warnings at half-threshold increments.
- **Task 0125: Daily limits + ordering** ✅
  - **Acceptance Criteria:**
    - New/review daily limits enforced per deck (or preset).
    - If user falls behind, older waiting cards prioritized (backlog behavior).

---

## Milestone M2 — Adding/Editing Content 🚧 IN PROGRESS
Goal: Users can add notes; note types define fields and templates; cards generated automatically.

### Feature 020: Add notes flow + duplicate check
**Summary:** Add Notes window, choose note type + deck, enter fields, duplicate detection.
**Manual pages:** 55–57.

**Tasks**
- **Task 0201: Add note UI** ✅
  - **Acceptance Criteria:**
    - User selects note type and target deck independently.
    - Fields render in configured order; required fields validated.
- **Task 0202: Duplicate check** ✅
  - **Acceptance Criteria:**
    - Configurable duplicate scope (collection/deck) and which field(s) participate.
    - User sees clear warning and can override if allowed.
- **Task 0203: Organizing content with tags/flags/marked** ⏳
  - **Acceptance Criteria:**
    - Add tags on create; toggle "Marked" tag; assign flags.
    - Tags/flags persist and are queryable in browser search.

### Feature 021: Note types, fields, and sort field
**Summary:** Manage note types; add/clone; manage fields; select sort field; RTL editing options.
**Manual pages:** 42–43, 58–60.

**Tasks**
- **Task 0211: Note type manager** ✅ (API + built-in types)
  - **Acceptance Criteria:**
    - Create new note type from built-in base ("Add") or clone existing.
    - Rename note type; delete note type only if safe (or block).
- **Task 0212: Field editor** ✅
  - **Acceptance Criteria:**
    - Add/rename/remove/reorder fields (drag-drop or reposition).
    - Prevent reserved field names (Tags, Type, Deck, Card, FrontSide).
- **Task 0213: Sort field** ⏳
  - **Acceptance Criteria:**
    - Exactly one field can be designated sort field.
    - Browser can display and sort by sort field.
- **Task 0214: Editing options** ✅
  - **Acceptance Criteria:**
    - Per-field font/size, HTML editor default, RTL editing option.

### Feature 022: Card templates (HTML/CSS) + card generation rules
**Summary:** Card templates define front/back and which cards are generated based on field presence; template options like deck override and browser appearance.
**Manual pages:** 72–104 (Card Templates + related troubleshooting), plus 72–73 excerpt.

**Tasks**
- **Task 0221: Template editor UI** ✅
  - **Acceptance Criteria:**
    - Edit Front/Back/Styling with live preview for current note sample.
    - Support multiple card templates per note type.
- **Task 0222: Conditional generation logic** ✅
  - **Acceptance Criteria:**
    - Support “only generate card if field has text” behaviors.
    - Changing templates triggers card regeneration logic safely.
- **Task 0223: Deck override per template** ✅
  - **Acceptance Criteria:**
    - If override set, cards of that template go to specified deck regardless of add-window deck.
- **Task 0224: Browser appearance templates** ⏳
  - **Acceptance Criteria:**
    - Optional simplified Q/A rendering for browser table columns.

### Feature 023: Cloze deletion
**Summary:** Cloze note type with {{cloze:Field}} template and c1/c2… deletions; handle empty cloze cards and cleanup.
**Manual pages:** ~95–105 and cloze troubleshooting around 102–104.

**Tasks**
- **Task 0231: Cloze editor support** ✅
  - **Acceptance Criteria:**
    - User can create cloze deletions with c1/c2 numbering.
    - Each cloze number creates a separate card.
- **Task 0232: Empty cloze card detection + cleanup** ✅
  - **Acceptance Criteria:**
    - If a cloze deletion number removed, resulting blank card is flagged.
    - “Empty Cards” tool deletes blank cards with user confirmation.
- **Task 0233: Cloze template validation** ✅
  - **Acceptance Criteria:**
    - Warn if cloze filter missing; provide auto-fix path.

### Feature 024: Image occlusion (IO) note type
**Summary:** Create cloze-like cards from images (masking areas), plus edit IO notes.
**Manual pages:** referenced in adding/editing section list (55) and in note types list (42).

**Tasks**
- **Task 0241: IO creation flow**
  - **Acceptance Criteria:**
    - User can upload image, draw masks, generate multiple cards.
- **Task 0242: IO editing flow**
  - **Acceptance Criteria:**
    - User can re-open an IO note and adjust masks; cards update accordingly.
- **Task 0243: IO rendering in review**
  - **Acceptance Criteria:**
    - Mask display matches created occlusions; reveals answer correctly.

---

## Milestone M3 — Browser (Search + Bulk Ops)
Goal: Users can browse/search/edit notes/cards, similar to Anki’s Browse window.

### Feature 030: Browse layout (sidebar + table + editor)
**Summary:** Browser has sidebar, card/note table, editing area; resizable panes; cards/notes mode toggle.
**Manual pages:** 138–143.

**Tasks**
- **Task 0301: Browser shell**
  - **Acceptance Criteria:**
    - 3-pane layout: sidebar, table, editor.
    - Resizable panes persist sizes.
- **Task 0302: Cards vs Notes table mode**
  - **Acceptance Criteria:**
    - Toggle between Cards and Notes modes.
    - Mode impacts row identity and available columns.

### Feature 031: Search system (query language + saved searches)
**Summary:** Search box supports Anki-like query grammar; saved searches; sidebar filters.
**Manual pages:** 138–150 (search-related subsections), plus unicode normalization note around 70–71.

**Tasks**
- **Task 0311: Search parser + execution**
  - **Acceptance Criteria:**
    - Support basic operators: text search, tag:, deck:, note:, card state filters, is:suspended/buried, etc. (implement progressively).
    - Correct tokenization + escaping for special characters.
- **Task 0312: Saved searches**
  - **Acceptance Criteria:**
    - Users can save queries; appear in sidebar; click to apply.
- **Task 0313: Unicode normalization toggle (advanced)**
  - **Acceptance Criteria:**
    - Default behavior normalizes text for search consistency; optional advanced toggle disables it.

### Feature 032: Columns, sorting, and browser actions
**Summary:** Configurable columns; sorting; edit actions; find & replace; duplicate finder; suspend/bury; set due date.
**Manual pages:** 141–149.

**Tasks**
- **Task 0321: Column configuration**
  - **Acceptance Criteria:**
    - Right-click config (web equivalent) lets user show/hide/reorder columns.
    - Sorting works for supported columns (not for raw rendered Q/A unless custom format provided).
- **Task 0322: Inline editing + note editor**
  - **Acceptance Criteria:**
    - Selecting row populates editor; edits persist.
    - Supports tag editing and field editing.
- **Task 0323: Find & Replace**
  - **Acceptance Criteria:**
    - Replace across selected notes/cards with preview and undo safety.
- **Task 0324: Find duplicates**
  - **Acceptance Criteria:**
    - User picks field to compare; returns duplicates grouped.
- **Task 0325: Bulk actions**
  - **Acceptance Criteria:**
    - Suspend/unsuspend, toggle marked, set flag, change deck, change note type.
    - “Set Due Date” and “Reposition” supported (if implementing scheduling tools).

---

## Milestone M4 — Deck Options + Advanced Studying
Goal: Provide deck options, filtered decks, custom study, and advanced scheduling configuration.

### Feature 040: Deck options presets
**Summary:** Deck options are presets applied to decks; include new/review limits, steps, lapses, bury, leech, etc.
**Manual pages:** ~110–131 (scheduler/options), plus bury siblings mention in overview (49).

**Tasks**
- **Task 0401: Preset system**
  - **Acceptance Criteria:**
    - Create/edit presets; assign to decks; subdecks inherit unless overridden.
- **Task 0402: New card options**
  - **Acceptance Criteria:**
    - Daily new limit, insertion order, bury related new cards, etc.
- **Task 0403: Review options**
  - **Acceptance Criteria:**
    - Daily review limit, interval modifiers, hard interval, easy bonus, max interval.
- **Task 0404: Lapses options**
  - **Acceptance Criteria:**
    - Relearning steps, leech threshold, leech action (tag/suspend).
- **Task 0405: Custom scheduling hook (advanced)**
  - **Acceptance Criteria:**
    - Optional “custom scheduling” JS hook can override computed states; gated behind “advanced” toggle.

### Feature 041: Filtered decks + custom study sessions
**Summary:** Create filtered decks from search strings; custom study presets; home deck behavior; rebuild/empty.
**Manual pages:** 149–151.

**Tasks**
- **Task 0411: Filtered deck creation**
  - **Acceptance Criteria:**
    - User defines search filter + limits + ordering; “Build” populates deck.
- **Task 0412: Home deck linkage**
  - **Acceptance Criteria:**
    - Cards moved into filtered deck remember home deck; return when emptied/deleted or after study per option.
- **Task 0413: Custom Study presets**
  - **Acceptance Criteria:**
    - “Custom Study” offers preset actions (review forgotten, review ahead, preview new, etc.).
    - Existing “Custom Study Session” replaced unless renamed.

### Feature 042: Reviewing ahead + rescheduling rules
**Summary:** Handling study of cards before due; whether it affects scheduling; cram-like behavior.
**Manual pages:** 149–151, plus revlog type notes (type=3 for early cram in filtered decks).

**Tasks**
- **Task 0421: Early review behavior**
  - **Acceptance Criteria:**
    - Early reviews are tracked distinctly (revlog type “filtered/early”).
    - Configurable whether early reviews reschedule or not.

---

## Milestone M5 — Statistics + Insights
Goal: Provide stats similar to Anki: today summary + graphs + export.

### Feature 050: Stats screen (today + graphs)
**Summary:** Deck/collection selection, history scopes, key metrics, graphs like future due, reviews, calendar.
**Manual pages:** 191–194.

**Tasks**
- **Task 0501: Scope selector (deck vs collection)**
  - **Acceptance Criteria:**
    - User can view stats for current deck, other deck, or whole collection.
- **Task 0502: Today summary**
  - **Acceptance Criteria:**
    - Shows Again count, correct %, counts by state (Learn/Review/Relearn/Filtered).
- **Task 0503: Graphs**
  - **Acceptance Criteria:**
    - At least: Reviews over time, Calendar heatmap, Future Due estimate.
    - Numbers match underlying revlog/card schedule.
- **Task 0504: Export stats**
  - **Acceptance Criteria:**
    - “Save PDF” (web equivalent: export to PDF/print view) generates shareable output.

---

## Milestone M6 — Sync + Safety
Goal: Multi-device sync (Microdote account), conflict handling, media sync semantics, deletion safety.

### Feature 060: Account setup + sync primitives
**Summary:** Sign-in, initial sync, automatic syncing, button state, and network settings.
**Manual pages:** 130–134.

**Tasks**
- **Task 0601: Sync authentication**
  - **Acceptance Criteria:**
    - User can sign in/out; sync is disabled when signed out.
- **Task 0602: Two-way merge sync**
  - **Acceptance Criteria:**
    - Reviews and note edits merge when made on different devices.
    - Uses per-object change tracking (USN-like) to avoid full resyncs.
- **Task 0603: One-way sync enforcement**
  - **Acceptance Criteria:**
    - Non-mergeable changes (note type/template changes) prompt user to choose local vs remote on next sync.
    - User can manually force one-way sync via setting.

### Feature 061: Media sync semantics + safety
**Summary:** Media merges independently; deletions only propagate when fully in sync; restore deleted media by logout trick.
**Manual pages:** 132–133, 198–200.

**Tasks**
- **Task 0611: Media storage + referencing**
  - **Acceptance Criteria:**
    - Media is stored separately; notes reference by filename/ID.
- **Task 0612: Media sync merge**
  - **Acceptance Criteria:**
    - Add/replace media syncs; one-way card sync does not override media merge.
- **Task 0613: Safe deletions**
  - **Acceptance Criteria:**
    - Deleted media does not propagate unless client is fully in sync; otherwise re-download.
    - Logging out and re-sync restores media if still present remotely.

### Feature 062: Firewalls/proxies + self-hosted sync (optional)
**Summary:** Network config for proxies/firewalls; optional self-hosted sync server.
**Manual pages:** 131–137.

**Tasks**
- **Task 0621: Network settings UI**
  - **Acceptance Criteria:**
    - Configure proxy URL and basic network troubleshooting hints.
- **Task 0622: Self-hosted server mode (optional)**
  - **Acceptance Criteria:**
    - Point client at custom sync endpoint; behaves like standard sync.

---

## Milestone M7 — Import/Export + Sharing
Goal: Import/export decks, text files, and full collection packages; share decks; handle updating and merging.

### Feature 070: Import formats (text + packaged decks)
**Summary:** Import TSV/text; import .apkg deck packages and .colpkg collection packages; updating behavior.
**Manual pages:** 174–177.

**Tasks**
- **Task 0701: Text import**
  - **Acceptance Criteria:**
    - Import TSV with field mapping; optionally update existing notes when first field stable.
    - Supports preserving HTML formatting.
- **Task 0702: Deck package import (.apkg)**
  - **Acceptance Criteria:**
    - Imports notes/cards/note types/media into existing collection.
    - Handles “update existing notes if newer” semantics.
- **Task 0703: Collection import (.colpkg)**
  - **Acceptance Criteria:**
    - Replaces current collection with imported one (with explicit warning + backup).
    - Does not delete existing media automatically; provide “check media” tool.

### Feature 071: Export formats (text + deck + collection)
**Summary:** Export notes to text; export deck package; export full collection; include media option.
**Manual pages:** 176–180.

**Tasks**
- **Task 0711: Export notes to text**
  - **Acceptance Criteria:**
    - Exports with tab-separated fields; preserves HTML formatting.
- **Task 0712: Export deck package**
  - **Acceptance Criteria:**
    - Exports selected deck + children + note types + referenced media.
- **Task 0713: Export collection package**
  - **Acceptance Criteria:**
    - Exports entire collection; optional compression; warns about version compatibility.

### Feature 072: Shared decks discovery (optional)
**Summary:** Browse/download shared decks, then import.
**Manual pages:** 43.

**Tasks**
- **Task 0721: Shared deck catalog integration**
  - **Acceptance Criteria:**
    - User can browse a catalog and import deck package into Microdote.

---

## Milestone M8 — Preferences, Media Tools, Add-ons, Troubleshooting
Goal: Provide settings, media management, plugin ecosystem, and operational tools.

### Feature 080: Preferences
**Summary:** Appearance, review, scheduler settings, syncing settings, backups configuration.
**Manual pages:** 104–106 (Preferences headings), 105–106.

**Tasks**
- **Task 0801: Preferences UI + persistence**
  - **Acceptance Criteria:**
    - Settings persist and apply immediately where safe.
- **Task 0802: Review/scheduler preferences**
  - **Acceptance Criteria:**
    - “Next day starts at” affects what counts as “today”.
    - Motion/compact mode toggles affect UI.
- **Task 0803: Backup retention configuration**
  - **Acceptance Criteria:**
    - User can configure backup intervals and retention counts.

### Feature 081: Backups + restore
**Summary:** Automatic backups, manual backups, restoring, deletion log.
**Manual pages:** 178–180.

**Tasks**
- **Task 0811: Automatic backups**
  - **Acceptance Criteria:**
    - Periodic snapshots created; retention policy enforced.
- **Task 0812: Restore workflow**
  - **Acceptance Criteria:**
    - Restoring clearly warns about losing newer changes; disables auto sync until restart/confirm.
- **Task 0813: Deletion log**
  - **Acceptance Criteria:**
    - Deleted items logged (for recovery/audit).

### Feature 082: Media tools
**Summary:** Media folder semantics, check media tool for unused/missing, filename encoding, supported formats.
**Manual pages:** 198–200.

**Tasks**
- **Task 0821: Check media report**
  - **Acceptance Criteria:**
    - Lists unused files and missing references.
    - Allows: delete unused, tag notes missing media, restore from trash if applicable.
- **Task 0822: Filename compatibility**
  - **Acceptance Criteria:**
    - Filenames sanitized/encoded for cross-platform; incompatibles flagged before sync.
- **Task 0823: Static media on every card**
  - **Acceptance Criteria:**
    - Support “_filename” convention (or Microdote equivalent) so checks don’t delete static resources.

### Feature 083: Math and symbols rendering
**Summary:** MathJax/LaTeX support, security warnings, web/mobile considerations.
**Manual pages:** 200–208.

**Tasks**
- **Task 0831: MathJax support**
  - **Acceptance Criteria:**
    - MathJax renders in review and preview; configurable via template.
- **Task 0832: LaTeX pipeline (optional)**
  - **Acceptance Criteria:**
    - If implemented: render LaTeX server-side or via WASM; safe command allowlist.
- **Task 0833: Security posture**
  - **Acceptance Criteria:**
    - User-generated HTML is sandboxed to prevent XSS; math rendering does not allow arbitrary JS execution.

### Feature 084: Add-ons / plugin system
**Summary:** Add-ons extend app; browse/install/remove; safe mode; version compatibility.
**Manual pages:** 210–212 and safe mode/startup options 182–184.

**Tasks**
- **Task 0841: Plugin architecture**
  - **Acceptance Criteria:**
    - Plugins can register: menu actions, UI panels, hooks (review scheduling, editor transforms).
    - Clear permission boundaries; plugin can be disabled.
- **Task 0842: Plugin marketplace + install by code**
  - **Acceptance Criteria:**
    - Install by ID/code; show enabled/disabled list; remove plugin.
- **Task 0843: Safe mode**
  - **Acceptance Criteria:**
    - Start Microdote with plugins disabled; user can re-enable selectively.

### Feature 085: Maintenance & troubleshooting tools
**Summary:** Database check/repair, startup options, portable installs/paths, corrupt collection recovery guidance.
**Manual pages:** 182–187.

**Tasks**
- **Task 0851: Database integrity check**
  - **Acceptance Criteria:**
    - Detect and report corruption; offer restore-from-backup first.
- **Task 0852: Portable / custom data folder (web equivalent)**
  - **Acceptance Criteria:**
    - Support exporting/importing a full “portable” bundle; document limitations for browsers.
- **Task 0853: Diagnostics**
  - **Acceptance Criteria:**
    - Provide debug console/log download; include sync logs and storage stats.

---

# Anki Feature List (Canonical Index)

> This section is intended to match the user’s requested structure.
> Each feature points back to milestones above for detailed tasks.

- **Feature 001:** Accounts, profiles, and collections — _M0_ — pages 42–43, 136–137
- **Feature 010:** Deck list + deck overview — _M1_ — pages 48–49
- **Feature 011:** Review screen + answer buttons — _M1_ — pages 49–50, 53–54
- **Feature 012:** Scheduler v1 (learning/review/relearn) — _M1_ — pages 49–55, 120–130
- **Feature 020:** Add notes + duplicate check — _M2_ — pages 55–57
- **Feature 021:** Note types + fields + sort field — _M2_ — pages 42–43, 58–60
- **Feature 022:** Card templates (HTML/CSS) — _M2_ — pages 72–104
- **Feature 023:** Cloze deletion — _M2_ — pages 95–105
- **Feature 024:** Image occlusion — _M2_ — pages 42, 55+
- **Feature 030:** Browse layout — _M3_ — pages 138–143
- **Feature 031:** Search system — _M3_ — pages 138–150, 70–71
- **Feature 032:** Columns + bulk ops — _M3_ — pages 141–149
- **Feature 040:** Deck options presets — _M4_ — pages 110–131
- **Feature 041:** Filtered decks + custom study — _M4_ — pages 149–151
- **Feature 050:** Statistics — _M5_ — pages 191–194
- **Feature 060:** Sync primitives — _M6_ — pages 130–134
- **Feature 061:** Media sync semantics — _M6_ — pages 132–133, 198–200
- **Feature 070:** Import — _M7_ — pages 174–177
- **Feature 071:** Export — _M7_ — pages 176–180
- **Feature 080:** Preferences — _M8_ — pages 104–106
- **Feature 081:** Backups — _M8_ — pages 178–180
- **Feature 082:** Check media — _M8_ — pages 198–200
- **Feature 083:** Math/symbols — _M8_ — pages 200–208
- **Feature 084:** Add-ons / plugins — _M8_ — pages 182–184, 210–212
- **Feature 085:** Troubleshooting / maintenance — _M8_ — pages 182–187

---

# Anki Task List (Dependency-Oriented)

> Tasks below are a flattened index you can paste into a tracker.
> Dependencies reference Task IDs.
> Status: ✅ = Complete, ⏳ = Partial, (blank) = Not started

- **Task 0001:** Define canonical data model ✅
  - **Dependencies:** —
- **Task 0002:** Implement collection storage layer ✅
  - **Dependencies:** Task 0001
- **Task 0003:** Profile support (phaseable) ✅
  - **Dependencies:** Task 0001, Task 0002
- **Task 0004:** Backups & restore hooks (foundation) ✅
  - **Dependencies:** Task 0002

- **Task 0101:** Deck list UI ✅
  - **Dependencies:** Task 0002
- **Task 0102:** Deck overview due counts ✅
  - **Dependencies:** Task 0101, Task 0121
- **Task 0103:** Study session start/stop ✅
  - **Dependencies:** Task 0102, Task 0111

- **Task 0111:** Question/answer rendering ✅
  - **Dependencies:** Task 0002, Task 0221
- **Task 0112:** Answer buttons + shortcuts ✅
  - **Dependencies:** Task 0111
- **Task 0113:** Time spent tracking ✅
  - **Dependencies:** Task 0111

- **Task 0121:** Scheduling state machine ✅
  - **Dependencies:** Task 0001
- **Task 0122:** Learning steps + graduating interval ✅
  - **Dependencies:** Task 0121
- **Task 0123:** Ease factor + interval updates ✅
  - **Dependencies:** Task 0121
- **Task 0124:** Lapses + leeches ✅
  - **Dependencies:** Task 0121
- **Task 0125:** Daily limits + ordering ✅
  - **Dependencies:** Task 0121, Task 0401

- **Task 0201:** Add note UI ✅
  - **Dependencies:** Task 0211, Task 0212
- **Task 0202:** Duplicate check ✅
  - **Dependencies:** Task 0201, Task 0311
- **Task 0203:** Tags/flags/marked ⏳
  - **Dependencies:** Task 0201, Task 0325

- **Task 0211:** Note type manager ✅
  - **Dependencies:** Task 0001
- **Task 0212:** Field editor ✅
  - **Dependencies:** Task 0211
- **Task 0213:** Sort field ⏳
  - **Dependencies:** Task 0212, Task 0302
- **Task 0214:** Editing options ✅
  - **Dependencies:** Task 0212

- **Task 0221:** Template editor UI ✅
  - **Dependencies:** Task 0211, Task 0212
- **Task 0222:** Conditional generation logic ✅
  - **Dependencies:** Task 0221, Task 0001
- **Task 0223:** Deck override per template ✅
  - **Dependencies:** Task 0221, Task 0101
- **Task 0224:** Browser appearance templates ⏳
  - **Dependencies:** Task 0221, Task 0301

- **Task 0231:** Cloze editor support ✅
  - **Dependencies:** Task 0221, Task 0201
- **Task 0232:** Empty cloze detection + cleanup ✅
  - **Dependencies:** Task 0231
- **Task 0233:** Cloze template validation ✅
  - **Dependencies:** Task 0221

- **Task 0241:** IO creation flow  
  - **Dependencies:** Task 0201, Task 0821
- **Task 0242:** IO editing flow  
  - **Dependencies:** Task 0241
- **Task 0243:** IO rendering in review  
  - **Dependencies:** Task 0111, Task 0241

- **Task 0301:** Browser shell  
  - **Dependencies:** Task 0002
- **Task 0302:** Cards vs Notes mode  
  - **Dependencies:** Task 0301, Task 0001

- **Task 0311:** Search parser + execution  
  - **Dependencies:** Task 0002
- **Task 0312:** Saved searches  
  - **Dependencies:** Task 0311
- **Task 0313:** Unicode normalization toggle  
  - **Dependencies:** Task 0311

- **Task 0321:** Column configuration  
  - **Dependencies:** Task 0302
- **Task 0322:** Inline editing + note editor  
  - **Dependencies:** Task 0301, Task 0201
- **Task 0323:** Find & Replace  
  - **Dependencies:** Task 0322
- **Task 0324:** Find duplicates  
  - **Dependencies:** Task 0311, Task 0213
- **Task 0325:** Bulk actions  
  - **Dependencies:** Task 0301, Task 0121

- **Task 0401:** Preset system  
  - **Dependencies:** Task 0001, Task 0101
- **Task 0402:** New card options  
  - **Dependencies:** Task 0401
- **Task 0403:** Review options  
  - **Dependencies:** Task 0401, Task 0123
- **Task 0404:** Lapses options  
  - **Dependencies:** Task 0401, Task 0124
- **Task 0405:** Custom scheduling hook  
  - **Dependencies:** Task 0401, Task 0121

- **Task 0411:** Filtered deck creation  
  - **Dependencies:** Task 0311, Task 0401
- **Task 0412:** Home deck linkage  
  - **Dependencies:** Task 0411
- **Task 0413:** Custom study presets  
  - **Dependencies:** Task 0411

- **Task 0421:** Early review behavior  
  - **Dependencies:** Task 0411, Task 0121

- **Task 0501:** Scope selector  
  - **Dependencies:** Task 0101
- **Task 0502:** Today summary  
  - **Dependencies:** Task 0113, Task 0121
- **Task 0503:** Graphs  
  - **Dependencies:** Task 0502
- **Task 0504:** Export stats  
  - **Dependencies:** Task 0503

- **Task 0601:** Sync authentication  
  - **Dependencies:** Task 0002
- **Task 0602:** Two-way merge sync  
  - **Dependencies:** Task 0601, Task 0001
- **Task 0603:** One-way sync enforcement  
  - **Dependencies:** Task 0602

- **Task 0611:** Media storage + referencing  
  - **Dependencies:** Task 0001
- **Task 0612:** Media sync merge  
  - **Dependencies:** Task 0611, Task 0602
- **Task 0613:** Safe deletions  
  - **Dependencies:** Task 0612, Task 0821

- **Task 0621:** Network settings UI  
  - **Dependencies:** Task 0601
- **Task 0622:** Self-hosted server mode  
  - **Dependencies:** Task 0602

- **Task 0701:** Text import  
  - **Dependencies:** Task 0201, Task 0211
- **Task 0702:** Deck package import  
  - **Dependencies:** Task 0701, Task 0611
- **Task 0703:** Collection import  
  - **Dependencies:** Task 0702, Task 0004

- **Task 0711:** Export notes to text  
  - **Dependencies:** Task 0701
- **Task 0712:** Export deck package  
  - **Dependencies:** Task 0702
- **Task 0713:** Export collection package  
  - **Dependencies:** Task 0703

- **Task 0721:** Shared deck catalog integration  
  - **Dependencies:** Task 0702

- **Task 0801:** Preferences UI + persistence  
  - **Dependencies:** Task 0002
- **Task 0802:** Review/scheduler preferences  
  - **Dependencies:** Task 0121
- **Task 0803:** Backup retention configuration  
  - **Dependencies:** Task 0811

- **Task 0811:** Automatic backups  
  - **Dependencies:** Task 0002
- **Task 0812:** Restore workflow  
  - **Dependencies:** Task 0811
- **Task 0813:** Deletion log  
  - **Dependencies:** Task 0002

- **Task 0821:** Check media report  
  - **Dependencies:** Task 0611
- **Task 0822:** Filename compatibility  
  - **Dependencies:** Task 0821
- **Task 0823:** Static media convention  
  - **Dependencies:** Task 0821

- **Task 0831:** MathJax support  
  - **Dependencies:** Task 0221
- **Task 0832:** LaTeX pipeline  
  - **Dependencies:** Task 0831
- **Task 0833:** Security posture  
  - **Dependencies:** Task 0221, Task 0601

- **Task 0841:** Plugin architecture  
  - **Dependencies:** Task 0002
- **Task 0842:** Plugin marketplace + install by code  
  - **Dependencies:** Task 0841
- **Task 0843:** Safe mode  
  - **Dependencies:** Task 0842

- **Task 0851:** Database integrity check  
  - **Dependencies:** Task 0002
- **Task 0852:** Portable bundle export/import  
  - **Dependencies:** Task 0713
- **Task 0853:** Diagnostics  
  - **Dependencies:** Task 0851, Task 0602

---

# Appendix A — Page Coverage Index (All Pages)

This index maps **print.html pages 1–223** into coverage buckets. The goal is **no gaps**.

> **Key:** each range maps to one or more Features above.

- **1–41:** Intro, fundamentals, main window/decks basics, core concepts (cards/notes/decks/collection) → Features 001, 010, 020, 021
- **42–46:** Note types, collection, shared decks → Features 021, 070, 072
- **47–54:** Studying basics, keyboard shortcuts, falling behind → Features 010, 011, 012
- **55–71:** Adding/editing, fields, unicode/input/normalization → Features 020, 021, 031
- **72–104:** Card templates, styling, typing answers, cloze + template troubleshooting → Features 022, 023, 011
- **105–110:** Preferences overview & early deck options → Features 080, 040
- **111–130:** Deck options, scheduler parameters, custom scheduling → Features 040, 012, 042
- **131–137:** Syncing with AnkiWeb, conflicts, proxies, profiles → Features 060, 061, 062, 001
- **138–151:** Browsing + search + filtered decks/custom study → Features 030, 031, 032, 041, 042
- **152–173:** Additional management tools & workflows (editing/browsing extensions) → Features 032, 085
- **174–180:** Import/export + backups → Features 070, 071, 081
- **181–190:** Install/launcher/program files/startup options/portable troubleshooting → Features 085, 084
- **191–197:** Statistics + revlog/data interpretation → Features 050
- **198–200:** Media & check media → Features 082, 061
- **200–209:** Math/symbols + leeches guidance → Features 083, 012
- **210–223:** Add-ons, troubleshooting, maintenance → Features 084, 085

If you want, I can generate a **second appendix** that breaks down the “152–173” and “210–223” blocks into more granular sub-features (it’s largely operational/troubleshooting text in the manual, and web implementations often consolidate it into fewer UI surfaces).
