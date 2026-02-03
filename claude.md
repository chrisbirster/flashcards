# Microdote — Project Context

## Project Overview

**Microdote** is a web-based reimplementation of Anki's spaced repetition flashcard system. The goal is to achieve feature parity with Anki while providing a modern, cross-platform experience through a React Native mobile app with a Go backend.

The complete feature roadmap and task breakdown is maintained in `microdote_anki_tasks.md`, which maps all Anki Manual features into 8 major milestones (M0-M8) with detailed acceptance criteria and dependencies.

## Tech Stack

### Frontend (Mobile App)
- **React Native 0.81** with **Expo 54**
- **TypeScript** for type safety
- **React Navigation** for routing and navigation
- **TanStack Query (React Query)** for data fetching and cache management

### Backend
- **Go** (see `go.mod` for version and dependencies)
- RESTful API serving the React Native frontend

### Storage
- Planning to use IndexedDB/SQLite-in-browser or remote database (see Task 0002 in task plan)

## Project Structure

```
/flashcards
├── main.go                      # Go backend server
├── go.mod, go.sum              # Go dependencies
├── web/                        # React Native frontend
├── microdote_anki_tasks.md     # Complete feature & task roadmap
└── claude.md                   # This file
```

## Architecture Guidelines

### Code Conventions
- Use TypeScript strict mode; avoid `any` types
- Follow React hooks best practices (useEffect dependencies, etc.)
- Use TanStack Query for all server state; local state for UI-only concerns
- Keep components small and focused; extract custom hooks for complex logic
- Name files with PascalCase for components, camelCase for utilities

### Test Coverage Requirements
**CRITICAL:** Maintain **≥90% test coverage** for all features going forward.

- **Backend (Go):** Unit tests for all CRUD operations, business logic, and edge cases
- **Frontend (React):** E2E tests with Playwright for all user workflows
- **Before merging any feature:** Run `./test.sh` and ensure all tests pass
- **Test-Driven Development:** Write tests alongside implementation, not as an afterthought
- **Coverage areas:**
  - Happy path (normal user flow)
  - Edge cases (empty states, validation, errors)
  - Data persistence (survives reloads)
  - Concurrent operations (rapid clicks, simultaneous edits)

**Test Types:**
- Backend: Unit tests (`*_test.go`)
- Frontend: E2E tests (`web/e2e/*.spec.ts`)
- Integration: API endpoint tests (via E2E tests calling real backend)

### Data Flow
- Backend (Go) serves RESTful API endpoints
- Frontend uses TanStack Query to fetch/mutate data
- Follow the data model defined in Task 0001 (Collection, Deck, NoteType, Card, etc.)

### Important Constraints
- Must support offline-first functionality (local storage + sync)
- Review scheduling must match Anki's behavior (FSRS-style algorithm)
- Import/export must be compatible with Anki's .apkg format
- Security: Sanitize user-generated HTML to prevent XSS (see Task 0833)

## Frontend Patterns & Conventions

### Preferences (React Context)
Use custom hooks for reading and updating preferences:

```typescript
// Simple boolean preference pattern
import {useAutoplayDisabled, useSetAutoplayDisabled} from '#/state/preferences'

function SettingsScreen() {
  const autoplayDisabled = useAutoplayDisabled()
  const setAutoplayDisabled = useSetAutoplayDisabled()

  return (
    <Toggle
      value={autoplayDisabled}
      onValueChange={setAutoplayDisabled}
    />
  )
}
```

### Navigation
Navigation uses **React Navigation** with type-safe route parameters.

#### Screen Components
```typescript
// Screen component
import {type NativeStackScreenProps} from '@react-navigation/native-stack'
import {type CommonNavigatorParams} from '#/lib/routes/types'

type Props = NativeStackScreenProps<CommonNavigatorParams, 'Profile'>

export function ProfileScreen({route, navigation}: Props) {
  const {name} = route.params  // Type-safe params

  return (
    <Layout.Screen>
      {/* Screen content */}
    </Layout.Screen>
  )
}
```

#### Programmatic Navigation
```typescript
// Using useNavigation hook
import {useNavigation} from '@react-navigation/native'

const navigation = useNavigation()
navigation.navigate('Profile', {name: 'alice.bsky.social'})

// Or use the navigate helper
import {navigate} from '#/Navigation'
navigate('Profile', {name: 'alice.bsky.social'})
```

### Platform-Specific Code
Use file extensions for platform-specific implementations:

```
Component.tsx          # Shared/default
Component.web.tsx      # Web-only
Component.native.tsx   # iOS + Android
Component.ios.tsx      # iOS-only
Component.android.tsx  # Android-only
```

### Import Aliases
**Always use the `#/` alias for absolute imports:**

```typescript
// Good ✅
import {useSession} from '#/state/session'
import {Button} from '#/components/Button'

// Avoid ❌
import {useSession} from '../../../state/session'
```

### React Compiler is Enabled
This codebase uses **React Compiler**, so don't proactively add `useMemo` or `useCallback`. The compiler handles memoization automatically.

```typescript
// UNNECESSARY ❌ - React Compiler handles this
const handlePress = useCallback(() => {
  doSomething()
}, [doSomething])

// JUST WRITE THIS ✅
const handlePress = () => {
  doSomething()
}
```

**Only use `useMemo`/`useCallback` when you have a specific reason:**
- The value is immediately used in an effect's dependency array
- You're passing a callback to a non-React library that needs referential stability

### Best Practices

#### Accessibility
- Always provide `label` prop for interactive elements
- Use `accessibilityHint` where helpful

#### Translations
- Wrap ALL user-facing strings with `msg()` or `<Trans>`

#### Styling
- Combine static atoms with theme atoms
- Use platform utilities for platform-specific styles

#### State Management
- Use **TanStack Query** for server state
- Use **React Context** for UI preferences

#### Components
- Check if a component exists in `#/components/` before creating new ones

#### Types
- Define explicit types for props
- Use `NativeStackScreenProps` for screens

#### Testing
- Components should have `testID` props for E2E testing

## Development Workflow

### Common Commands
```bash
# Backend
go run main.go              # Start Go server
go test ./...              # Run backend tests

# Frontend (from web/ directory)
npm start                  # Start Expo development server
npm run ios                # Run on iOS simulator
npm run android            # Run on Android emulator
npm test                   # Run frontend tests
```

### Testing Strategy
- Unit tests for scheduling algorithm (critical path)
- Integration tests for import/export
- E2E tests for core study flow (add card → study → review)

## Milestone Progress Tracking

This section tracks completed tasks and implementation plans. Each time a significant plan is created or milestone work is completed, document it here.

### Current Status: M2 — Adding/Editing Content (In Progress)

**Previous Milestones:**
- M0 — Product Skeleton + Core Data Model ✅ COMPLETE
- M1 — Studying MVP ✅ COMPLETE (36 tests passing, 90%+ coverage)

**Current Milestone:** M2 — Adding/Editing Content
**M2 Goal:** Users can create, edit, and delete notes/cards within the app.

**M2 Tasks Completed:**
- ✅ Task 0201: Add note UI (AddNoteScreen component with field inputs, preview)
- ✅ Task 0211: Note type manager + API (GET /api/note-types endpoint)

**M2 Tasks Remaining:**
- Task 0202: Duplicate check
- Task 0203: Tags/flags/marked support (tags working, flags/marked pending)
- Task 0212: Field editor (edit note type fields)
- Task 0221: Template editor UI
- Task 0231: Cloze editor support

---

### Completed Plans & Tasks

#### 2026-02-01 - Plan: Milestone M0 Implementation
**Related Tasks:** 0001, 0002, 0003, 0004
**Accomplishments:**
- ✅ **Task 0001:** Enhanced data model with tags, flags, media refs, sync metadata (USN), deck options, deck hierarchy
- ✅ **Task 0002:** Implemented SQLite storage layer with Store interface, CRUD operations, migrations, and transactions
- ✅ **Task 0003:** Added profile support foundation (Profile struct, CRUD, active profile tracking)
- ✅ **Task 0004:** Implemented backup/restore functionality (ZIP-based backups, cleanup, retention policy)
- ✅ **REST API:** Built Go HTTP server with chi router serving JSON endpoints for decks, notes, cards, backups
- ✅ **Frontend:** Created React web app with TanStack Query, deck list view, create deck functionality

**Key Decisions:**
- Used SQLite for storage (file-based, Anki-compatible, easy backup/restore)
- All business logic in Go backend, React frontend is view-only calling REST API
- REST API with JSON over HTTP (simple, well-supported)
- chi router for clean, type-safe routing
- TanStack Query for server state management on frontend

**Database Schema:**
- `collections`, `profiles`, `decks`, `deck_options`, `note_types`, `notes`, `cards`, `revlog`, `media`, `metadata`
- Proper foreign keys, indexes for performance
- Migration system with versioning

**API Endpoints:**
```
GET    /api/health
GET    /api/collection
GET    /api/decks
POST   /api/decks
GET    /api/decks/{id}
GET    /api/decks/{deckId}/due
POST   /api/notes
GET    /api/notes/{id}
GET    /api/cards/{id}
POST   /api/cards/{id}/answer
POST   /api/backups
GET    /api/backups
```

**Status:** ✅ Complete

---

#### 2026-02-02 - Plan: Milestone M1 Implementation (Studying MVP)
**Related Tasks:** 0101, 0102, 0111, 0112
**Accomplishments:**
- ✅ **Task 0102:** Deck stats endpoint - shows New/Learning/Review/Suspended card counts
- ✅ **Task 0111:** Question/answer rendering - StudyScreen component with card display
- ✅ **Task 0112:** Answer buttons + keyboard shortcuts (Space/Enter to show answer, 1-4 for ratings)
- ✅ **Study Flow:** Complete study workflow from fetching due cards → answering → FSRS scheduling

**Key Features Implemented:**
- `GET /api/decks/{id}/stats` endpoint returns card counts by state
- `DeckStats` struct with NewCards, Learning, Review, Relearning, Suspended, Buried, TotalCards, DueToday
- `StudyScreen.tsx` component with full study UI
- Keyboard shortcuts: Space/Enter for show answer or "Good", 1-4 for Again/Hard/Good/Easy
- Time tracking for each card review
- Progress indicator showing current card position
- Exit button to return to deck list
- "All done!" completion state when no cards are due
- Stats refresh after answering cards

**Test Coverage:**
- **Backend:** 12 unit tests (+4 for M1)
  - TestGetDeckStats: Card counting by state
  - TestStudyFlow: Complete study workflow
  - TestAnswerCardMultipleRatings: All 4 FSRS ratings
  - TestEmptyDeckStudy: Empty deck edge case
- **Frontend:** 24 E2E tests (+14 for M1)
  - Study screen rendering and navigation
  - Keyboard shortcut handling
  - Progress indicators and stats display
  - Exit functionality
  - Completion state
- **Total:** 36 passing tests

**Status:** ✅ Complete

---

**Verification (M0):**
- Backend server runs on :8080
- Frontend connects via proxy (Vite dev server)
- Can create decks via API and UI
- Data persists to SQLite database
- Backup creation works
- Profile system initialized

**Bug Fixes:**
- Fixed Collection ID counter initialization (nextDeckID, nextNoteID, nextCardID now load from DB max values)
- Deck creation now works correctly with persisted data

**Test Coverage (M0):**
- **Backend:** 8 unit tests (Go testing framework)
  - Deck CRUD operations
  - Note CRUD operations
  - Card CRUD with FSRS state
  - Due card filtering
  - Profile management
  - Collection ID counter initialization
  - Backup creation
- **Frontend:** 10 E2E tests (Playwright)
  - Application rendering
  - Deck creation and listing
  - UI state management
  - Form validation
  - Data persistence across reloads
  - Rapid operations handling

**Test Commands:**
```bash
# Run all tests
./test.sh

# Backend only
go test -v

# Frontend only (requires backend running)
cd web && npm run test:e2e

# Frontend with UI
cd web && npm run test:e2e:ui
```

---

## Task Plan Reference

The complete task breakdown is in `microdote_anki_tasks.md`. Key milestones:

- **M0:** Product Skeleton + Core Data Model (Tasks 0001-0004)
- **M1:** Studying MVP - Review Loop (Tasks 0101-0125)
- **M2:** Adding/Editing Content (Tasks 0201-0243)
- **M3:** Browser - Search + Bulk Ops (Tasks 0301-0325)
- **M4:** Deck Options + Advanced Studying (Tasks 0401-0421)
- **M5:** Statistics + Insights (Tasks 0501-0504)
- **M6:** Sync + Safety (Tasks 0601-0622)
- **M7:** Import/Export + Sharing (Tasks 0701-0721)
- **M8:** Preferences, Media Tools, Add-ons, Troubleshooting (Tasks 0801-0853)

## Things to Avoid

- Don't break Anki import/export compatibility
- Don't skip validation on user input (especially in search/template editing)
- Don't implement custom scheduling without maintaining FSRS compatibility
- Don't add UI features that aren't in the task plan without explicit approval
- Don't commit sensitive data (API keys, user data) to the repository

## Open Questions / Decisions Needed

_(Track architectural decisions that need user input here)_

- [ ] Which database to use for production? (IndexedDB, SQLite WASM, remote PostgreSQL)
- [ ] Self-hosted sync server vs cloud service?
- [ ] Native mobile features to prioritize (offline, notifications, widgets)?

## Resources

- **Anki Manual:** The canonical reference for expected behavior
- **Task Plan:** `microdote_anki_tasks.md` - detailed implementation roadmap
- **FSRS Algorithm:** https://github.com/open-spaced-repetition/fsrs4anki (for scheduler implementation)

---

_This file should be updated as the project evolves. When completing tasks or creating implementation plans, document them in the "Completed Plans & Tasks" section above._
