# Vutadex Infra Cheatsheet

This is the short version of the commands we've been using most often.

Assumptions:
- run commands from the repo root unless noted otherwise
- `.env.local` is for local SQLite
- `.env.production` is the production source of truth
- Fly app name is `vutadex-app`

## Local dev

Create the local env file:

```bash
task env:init:local
```

Run app + backend locally with SQLite:

```bash
ENV_FILE=.env.local task dev:sqlite:app
```

Run app + backend + marketing locally with SQLite:

```bash
ENV_FILE=.env.local task dev:sqlite
```

## Turso

Create a Turso database:

```bash
turso db create vutadex-prod
```

Hydrate `.env.production` from Turso:

```bash
TURSO_DB=vutadex-prod task env:hydrate:production
```

Run locally against the production-style Turso config:

```bash
ENV_FILE=.env.production task dev:turso:app
```

## SST / marketing / email

Set the SST email API secret once:

```bash
cd infra
bunx sst secret set EmailApiKey "<strong-random-token>" --stage production
```

Deploy the SST stack using `.env.production`:

```bash
task deploy:infra
```

Alias if you only want the email/SST path:

```bash
task deploy:email
```

Alias if you mean the marketing site deploy:

```bash
task deploy:marketing
```

Notes:
- all three tasks above currently run the same SST deploy
- `task deploy:marketing` and `task deploy:email` both source `.env.production`

## Fly

Import `VUTADEX_*` secrets from `.env.production` into Fly:

```bash
ENV_FILE=.env.production FLY_APP=vutadex-app task secrets:fly:import
```

Deploy the app server to Fly:

```bash
FLY_APP=vutadex-app task deploy:fly
```

Import secrets and deploy in one step:

```bash
ENV_FILE=.env.production FLY_APP=vutadex-app task deploy:fly:sync
```

Create the Fly app if it doesn't exist yet:

```bash
fly apps create vutadex-app
```

Check app status:

```bash
fly status -a vutadex-app
```

List machines:

```bash
fly machines list -a vutadex-app
```

Read recent logs:

```bash
fly logs -a vutadex-app
```

Quick public health check:

```bash
curl -I https://app.vutadex.com
```

## Stripe

Bootstrap billing products and prices into `.env.production`:

```bash
ENV_FILE=.env.production task stripe:billing:bootstrap
```

Forward Stripe billing webhooks to local backend:

```bash
task stripe:billing:listen
```

Billing values expected in `.env.production`:

```env
VUTADEX_STRIPE_SECRET_KEY=sk_live_...
VUTADEX_STRIPE_WEBHOOK_SECRET=whsec_...
VUTADEX_STRIPE_BILLING_PRICE_PRO_MONTHLY=price_...
VUTADEX_STRIPE_BILLING_PRICE_TEAM_MONTHLY=price_...
VUTADEX_STRIPE_BILLING_CHECKOUT_SUCCESS_URL=https://app.vutadex.com/billing/complete?checkout=success
VUTADEX_STRIPE_BILLING_CHECKOUT_CANCEL_URL=https://app.vutadex.com/billing/complete?checkout=cancelled
VUTADEX_STRIPE_BILLING_PORTAL_RETURN_URL=https://app.vutadex.com/settings?billing=returned
```

## Useful production checks

Confirm the marketing site:

```bash
curl -I https://vutadex.com
```

Confirm the app domain:

```bash
curl -I https://app.vutadex.com
```

Confirm Fly machines wake on traffic:

```bash
curl -I https://app.vutadex.com
fly machines list -a vutadex-app
```

## Current production URLs

- Marketing: `https://vutadex.com`
- App: `https://app.vutadex.com`
- Fly fallback: `https://vutadex-app.fly.dev`

