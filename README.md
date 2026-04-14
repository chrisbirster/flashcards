# Vutadex

Vutadex is a browser-first flashcard and study workspace with OTP sign-in, FSRS-based spaced repetition, per-deck workload controls, analytics, focus sessions, and paid Pro/Team plans.

- Marketing: [https://vutadex.com](https://vutadex.com)
- App: [https://app.vutadex.com](https://app.vutadex.com)

## What’s in this repo

- Go backend and API
- React app in [web](./web)
- Marketing site in [marketing](./marketing)
- SST infrastructure in [infra](./infra)

## Core product areas

- OTP-based sign-in
- Notes, cards, templates, and decks
- FSRS scheduling
- Per-deck new/review caps and priority controls
- Study analytics and focus sessions
- Team administration and billing

## Quick start

Create a local SQLite env file:

```bash
task env:init:local
```

Run the app locally:

```bash
ENV_FILE=.env.local task dev:sqlite:app
```

Run the app plus marketing site locally:

```bash
ENV_FILE=.env.local task dev:sqlite
```

## Production workflow

Hydrate `.env.production` from Turso:

```bash
TURSO_DB=vutadex-prod task env:hydrate:production
```

Deploy the Fly app:

```bash
ENV_FILE=.env.production FLY_APP=vutadex-app task deploy:fly:sync
```

Deploy the SST infra stack:

```bash
task deploy:infra
```

## Billing

The app supports:

- Free
- Pro: $12/month
- Team: $8/user/month, 3-seat minimum
- Enterprise: contact sales

Bootstrap Stripe billing prices into `.env.production`:

```bash
ENV_FILE=.env.production task stripe:billing:bootstrap
```

## Documentation

- Infra setup: [infra/README.md](./infra/README.md)
- Infra command cheatsheet: [infra/CHEATSHEET.md](./infra/CHEATSHEET.md)
- Product roadmap: [phases.md](./phases.md)

## License

This repository is publicly visible for reference and transparency, but it is
not released under an open source license.

Copyright © 2026 Chris Birster. All rights reserved.

You may not copy, modify, redistribute, or use this code to operate a
competing service without prior written permission.

## Notes

- Study Groups is not part of the current public release and is currently marked as coming soon.
- The app uses `.env.local` for local SQLite development and `.env.production` as the production source of truth.
