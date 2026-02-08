# Coverage Ramp Plan (Toward 95%)

This repo currently has:
- Backend (Go) statement coverage: **40.2%** (`go test ./... -coverprofile`)
- Frontend unit coverage: **not yet measurable in this environment** (npm registry blocked for Vitest/RTL install)

## Stage Gates

Use `./scripts/check_coverage.sh <stage>`.

| Stage | Backend | Frontend | Purpose |
|---|---:|---:|---|
| `baseline` | 40% | 0% | Lock current minimums and prevent regressions |
| `m3` | 55% | 60% | Browser/search foundation with test-backed API/UI changes |
| `m4` | 75% | 80% | Deck options + advanced studying coverage |
| `release` | 95% | 95% | Final parity-quality gate |

## Immediate Next Tasks

1. Unblock frontend unit tests:
   - Run in `web/`:  
     `npm install -D vitest @vitest/coverage-v8 @testing-library/react @testing-library/jest-dom @testing-library/user-event jsdom`
2. Start frontend unit coverage:
   - `npm --prefix web run test:unit`
   - `npm --prefix web run test:coverage`
3. Expand backend API handler tests with `httptest`:
   - Deck CRUD error paths
   - Note type field/template update validation
   - Card answer/update endpoints
4. Raise stage from `baseline` to `m3` once both backend + frontend thresholds pass consistently.

## Notes

- Frontend unit-test scaffolding is checked in:
  - `web/vitest.config.ts`
  - `web/test/setup.ts`
  - `web/test/add-note-form-provider.test.tsx`
  - `web/test/note-type-editor-routes.test.tsx`
- Until npm registry access is available, frontend coverage remains blocked in this environment.
