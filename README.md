# Plandalf Web

This repository contains the hosted web client for Plandalf.

The first goal is simple: make the same Plandalf study cycle available from a phone or browser without requiring the computer that owns a local SQLite file to stay online.

The hosted app uses:

- Go for the HTTP API
- Turso/libSQL for hosted SQLite
- regular SQLite for local development
- SolidJS 2 for the web client
- StyleX for component styling
- Plandalf's FSRS-7 scheduling rules and immutable review history

The current web app is intentionally focused on studying. Deck and note authoring, local-to-hosted sync, and the later StingJS native client will build on this foundation.

## Architecture

```text
SolidJS 2 + StyleX web app
            |
            v
       Plandalf API
            |
            v
      Turso / libSQL
```

The hosted scheduler is a direct port of Plandalf's current Zig FSRS-7 implementation, including its 35 default weights, fractional-day calculations, parameter-set identity, review audit fields, and rebuildable scheduler state.

## Current API

```text
GET  /api/v1/health
GET  /api/v1/decks
POST /api/v1/decks
POST /api/v1/decks/:deckID/cards
GET  /api/v1/decks/:deckID/study/next
GET  /api/v1/cards/:cardID
GET  /api/v1/cards/:cardID/study/preview
POST /api/v1/cards/:cardID/reviews
```

The deck/card POST routes are useful for bootstrapping the first hosted database. The long-term authoring path will work with Plandalf notes and synchronization rather than treating raw card creation as the primary workflow.

## Run locally

Build the Solid app first because the Go binary embeds `web/dist`:

```bash
cd web
bun install
bun run build
cd ..

go run .
```

The default local database is:

```text
./data/plandalf.db
```

For frontend development, run the Go API on port 8000 and Vite separately:

```bash
# terminal 1
mkdir -p data
go run .

# terminal 2
cd web
bun run dev
```

Vite proxies `/api` to `http://localhost:8000`.

## Turso

Production uses the same connection pattern the previous Vutadex application used: the Go `database/sql` layer switches from the local SQLite driver to the libSQL driver when a database URL is configured.

Set:

```bash
export PLANDALF_DATABASE_URL='libsql://your-database.turso.io'
export PLANDALF_DATABASE_AUTH_TOKEN='your-turso-token'
export PLANDALF_API_TOKEN='your-private-api-token'
export PLANDALF_APP_ORIGIN='https://your-plandalf-app.example'
```

Optional additional browser origins can be provided as a comma-separated list:

```bash
export PLANDALF_ALLOWED_ORIGINS='https://example.com,https://preview.example.com'
```

If `PLANDALF_API_TOKEN` is configured, clients send it as a bearer token. The browser stores the token locally on that device after you connect.

## Bootstrap a deck

For local development with no API token:

```bash
curl -X POST http://localhost:8000/api/v1/decks \
  -H 'Content-Type: application/json' \
  -d '{"name":"MongoDB"}'
```

Then add a raw test card using the returned deck ID:

```bash
curl -X POST http://localhost:8000/api/v1/decks/1/cards \
  -H 'Content-Type: application/json' \
  -d '{"question":"What does $match do?","answer":"Filters documents in an aggregation pipeline."}'
```

Open the app and the deck will be available to study.

## Study behavior

The current hosted study queue keeps the Plandalf behavior we already established:

- due reviews are shown before unseen cards
- unseen cards are currently unlimited unless a later study-policy layer limits their introduction
- revealing a card does not change scheduler state
- selecting Again, Hard, Good, or Easy creates an immutable review event
- FSRS-7 calculates and stores the next due time
- the next eligible card is fetched immediately

The API includes an expected review count with each submission so stale or duplicate browser reviews can be rejected safely.

## SolidJS 2 and StingJS

The web client uses SolidJS 2 directly rather than a React compatibility layer. Component styling is done with StyleX.

The UI starts with a small set of primitives such as `Stack`, `Text`, `Button`, and `Surface`. Keeping application screens expressed through a small primitive layer gives us a clean seam for the planned StingJS native client later, while the browser implementation remains ordinary SolidJS 2 today.

## Next work

The most important next step is synchronization between the normal Plandalf database and the hosted Turso database. That is what will let decks and immutable review history move safely between the CLI/desktop workflow and phone study without maintaining two independent scheduling histories.

## License

See [LICENSE.md](./LICENSE.md) for the repository's current license terms.
