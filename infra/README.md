# Vutadex Infra (SST)

This directory provisions two things with [SST](https://sst.dev/docs/):

1. AWS SES-backed email API used by the Go backend.
2. Cloudflare-hosted Vutadex web app (`vutadex.com`) from the Vite build in `web/`.

## What gets created

- `sst.aws.Email` identity (`VutadexEmail`)
- `sst.Secret` auth key (`EmailApiKey`)
- `sst.aws.Function` URL (`EmailApi`) that sends via SES
- `sst.cloudflare.StaticSite` (`WebSite`) for the browser app

## Prerequisites

1. AWS credentials configured locally.
2. SES configured in the same AWS region as the SST stack.
3. Turso CLI installed and authenticated if you want remote non-prod/prod env hydration.
4. Fly CLI installed and authenticated for secret sync + deploy.
5. Stripe CLI installed if you want local webhook forwarding.
3. Cloudflare API token and account ID:
   - `CLOUDFLARE_API_TOKEN`
   - `CLOUDFLARE_DEFAULT_ACCOUNT_ID`
6. `npm` installed (used to build the Vite app during deploy).

Official references:
- Turso CLI: [docs.turso.tech/cli/introduction](https://docs.turso.tech/cli/introduction)
- Turso DB show URL: [docs.turso.tech/cli/db/show](https://docs.turso.tech/cli/db/show)
- Turso auth tokens: [docs.turso.tech/cli/db/tokens/create](https://docs.turso.tech/cli/db/tokens/create)
- Fly secrets import: [fly.io/docs/flyctl/secrets-import](https://fly.io/docs/flyctl/secrets-import/)
- Stripe customer portal: [docs.stripe.com/customer-management](https://docs.stripe.com/customer-management)
- Stripe customer portal configuration: [docs.stripe.com/customer-management/configure-portal](https://docs.stripe.com/customer-management/configure-portal)
- Stripe CLI webhook forwarding: [docs.stripe.com/stripe-cli/use-cli](https://docs.stripe.com/stripe-cli/use-cli)

## SES requirements

SST provisions the email function and SES-backed sending path, but AWS SES still
has to be made production-ready outside SST before `no-reply@vutadex.com` can
send OTPs.

Required AWS + DNS steps:

1. Create or confirm the SES identity for `vutadex.com` in the same AWS region
   used by this stack.
2. Add the SES verification and DKIM records in DNS manually where
   `vutadex.com` is hosted.
   - If DNS is on Cloudflare, add the SES TXT/CNAME records there.
3. Request SES production access so the region is no longer sandboxed.
4. Add deliverability records:
   - SPF
   - DMARC

Recommended sender settings:

- `VUTADEX_EMAIL_SENDER=no-reply@vutadex.com`
- `VUTADEX_EMAIL_FROM=no-reply@vutadex.com`

If SES verification or sandbox removal is incomplete, the SST function may
deploy successfully but production OTP mail will not send reliably.

## Single-source env workflow

The simplest pattern for this repo is:

1. Keep one env file per environment.
2. Hydrate what we can from CLIs into that file.
3. Push that file into Fly with one command.

Recommended env files:

- `.env.local` for SQLite-only local dev
- `.env.staging` for Turso-backed staging / non-prod
- `.env.production` for Turso-backed production

Important reality check:

- Turso can give us the database URL and auth token from the CLI.
- Fly can import secrets from a file, but **cannot reveal existing secret values back**. Treat your env file as the source of truth.
- Stripe can be partially automated:
  - products/prices can be bootstrapped from the repo
  - webhook forwarding can be done locally with the Stripe CLI
  - secret keys and production webhook endpoint secrets still come from Stripe at bootstrap time

## Environment

Required:

```bash
export CLOUDFLARE_API_TOKEN="<cloudflare-api-token>"
export CLOUDFLARE_DEFAULT_ACCOUNT_ID="<cloudflare-account-id>"
```

Optional:

```bash
export AWS_REGION=us-east-1
export VUTADEX_EMAIL_SENDER=no-reply@vutadex.com
export VUTADEX_EMAIL_FROM=no-reply@vutadex.com
export VUTADEX_EMAIL_API_AUTH_HEADER=Authorization

# Defaults to vutadex.com on production, <stage>.vutadex.com otherwise
export VUTADEX_WEB_DOMAIN=vutadex.com
```

## Local SQLite dev

Create a local env file once:

```bash
task env:init:local
```

Run the app locally with SQLite:

```bash
ENV_FILE=.env.local task dev:sqlite:app
```

Run backend + app + marketing locally with SQLite:

```bash
ENV_FILE=.env.local task dev:sqlite
```

## Staging / production env hydration with Turso

Hydrate a remote env file directly from Turso without opening the dashboard:

```bash
ENV_FILE=.env.staging \
TURSO_DB=vutadex-staging \
APP_ORIGIN=https://staging.app.vutadex.com \
MARKETING_ORIGIN=https://staging.vutadex.com \
VUTADEX_ENV=staging \
task env:hydrate:turso
```

Production example:

```bash
ENV_FILE=.env.production \
TURSO_DB=vutadex-prod \
APP_ORIGIN=https://app.vutadex.com \
MARKETING_ORIGIN=https://vutadex.com \
VUTADEX_ENV=production \
task env:hydrate:turso
```

That task writes:

- `VUTADEX_DATABASE_URL`
- `VUTADEX_DATABASE_AUTH_TOKEN`
- `VUTADEX_APP_ORIGIN`
- `VUTADEX_MARKETING_ORIGIN`
- `VUTADEX_COOKIE_SECURE`
- `VITE_APP_ORIGIN`
- `VUTADEX_SESSION_SECRET` (if missing)

Run the app locally against Turso:

```bash
ENV_FILE=.env.staging task dev:turso:app
```

Or with marketing too:

```bash
ENV_FILE=.env.staging task dev:turso
```

## First-time setup

```bash
cd infra
npm install
npx sst install
npx sst secret set EmailApiKey "<strong-random-token>" --stage production
```

Before production deploy, verify the SES identity and publish the SES + DKIM
records in DNS for `vutadex.com`.

## Stripe setup for subscription billing

### What the app expects

The app uses:

- Stripe Checkout for first paid subscription purchase
- Stripe Customer Portal for existing paid customer changes
- `POST /api/billing/webhook` for subscription sync

Launch pricing:

- Free: `$0`
- Pro: `$12/month`
- Team: `$8/user/month`, 3-seat minimum
- Enterprise: manual / sales-led

Billing behavior:

- upgrades happen immediately
- downgrades happen at period end
- paid access remains enabled through `current_period_end`

### Step 1: put the Stripe secret key into your env file

Add this manually once to `.env.staging` or `.env.production`:

```bash
VUTADEX_STRIPE_SECRET_KEY=sk_live_...
```

This is the one part that still starts in Stripe itself.

### Step 2: bootstrap the Stripe billing product and prices

The repo can create or find the billing product and monthly prices for you and
write the resulting price IDs back into the same env file:

```bash
ENV_FILE=.env.staging task stripe:billing:bootstrap
```

That fills:

- `VUTADEX_STRIPE_BILLING_PRICE_PRO_MONTHLY`
- `VUTADEX_STRIPE_BILLING_PRICE_TEAM_MONTHLY`

### Step 3: configure the Customer Portal in Stripe

In Stripe Dashboard:

1. Turn on subscription management in the Customer Portal.
2. Allow plan switching between prices on the same product.
3. Set downgrades to happen **at period end**.
4. Turn on payment method management and invoice history.
5. Leave seat quantity changes out of the portal. Team seat counts are managed by the app.

### Step 4: create the production webhook endpoint in Stripe

Create a webhook endpoint pointing at:

- `https://app.vutadex.com/api/billing/webhook`

Subscribe it to:

- `checkout.session.completed`
- `customer.subscription.created`
- `customer.subscription.updated`
- `customer.subscription.deleted`

Then copy the webhook signing secret into the env file:

```bash
VUTADEX_STRIPE_WEBHOOK_SECRET=whsec_...
```

### Step 5: local Stripe webhook testing

For local backend testing:

```bash
task stripe:billing:listen
```

That forwards billing events to:

- `http://localhost:8000/api/billing/webhook`

Stripe CLI prints a temporary local webhook secret. Put that into your local env file
as `VUTADEX_STRIPE_WEBHOOK_SECRET` while testing.

## Deploy

```bash
cd infra
npx sst deploy --stage production
```

Capture outputs:
- `emailApiBaseUrl`
- `authHeaderName`
- `webDomain`
- `webUrl`

Then set Fly secrets in the API app:

```bash
fly secrets set \
  VUTADEX_OTP_MAIL_PROVIDER=sst \
  VUTADEX_TEAM_INVITE_MAIL_PROVIDER=sst \
  VUTADEX_EMAIL_SEND_URL="<emailApiBaseUrl>send" \
  VUTADEX_EMAIL_SEND_AUTH_HEADER="<authHeaderName>" \
  VUTADEX_EMAIL_SEND_AUTH_VALUE="<same value you set in EmailApiKey>"
```

Deploy/restart Fly:

```bash
fly deploy
```

## Push env into Fly from one file

Instead of copying secrets one by one, use the env file as the source of truth:

```bash
ENV_FILE=.env.staging FLY_APP=vutadex-staging task secrets:fly:import
```

Production example:

```bash
ENV_FILE=.env.production FLY_APP=vutadex-app task secrets:fly:import
```

Then deploy:

```bash
FLY_APP=vutadex-app task deploy:fly
```

Or do both in one step:

```bash
ENV_FILE=.env.production FLY_APP=vutadex-app task deploy:fly:sync
```

## Stripe subscription billing env (app backend)

Set these Fly secrets for workspace subscription billing:

```bash
fly secrets set \
  VUTADEX_STRIPE_SECRET_KEY="sk_live_..." \
  VUTADEX_STRIPE_WEBHOOK_SECRET="whsec_..." \
  VUTADEX_STRIPE_BILLING_PRICE_PRO_MONTHLY="price_..." \
  VUTADEX_STRIPE_BILLING_PRICE_TEAM_MONTHLY="price_..." \
  VUTADEX_STRIPE_BILLING_CHECKOUT_SUCCESS_URL="https://app.vutadex.com/billing/complete?checkout=success" \
  VUTADEX_STRIPE_BILLING_CHECKOUT_CANCEL_URL="https://app.vutadex.com/billing/complete?checkout=cancelled" \
  VUTADEX_STRIPE_BILLING_PORTAL_RETURN_URL="https://app.vutadex.com/settings?billing=returned"
```

Launch pricing assumptions baked into the app:

- Free: `$0`
- Pro: `$12/month`
- Team: `$8/user/month`, 3-seat minimum
- Enterprise: manual / sales-led

Subscription billing rules:

- upgrades happen immediately
- downgrades happen at period end
- paid features remain enabled through `current_period_end`
- Team seat increases prorate upward immediately
- Team seat decreases reconcile at renewal rather than removing access mid-cycle

Then redeploy app:

```bash
fly deploy
```

## Marketplace commerce env (app backend)

Set these Fly secrets to enable live marketplace creator onboarding and premium deck checkout:

```bash
fly secrets set \
  VUTADEX_STRIPE_SECRET_KEY="sk_live_..." \
  VUTADEX_STRIPE_WEBHOOK_SECRET="whsec_..." \
  VUTADEX_STRIPE_CONNECT_COUNTRY="US" \
  VUTADEX_STRIPE_CONNECT_REFRESH_URL="https://app.vutadex.com/marketplace/publish?creator=refresh" \
  VUTADEX_STRIPE_CONNECT_RETURN_URL="https://app.vutadex.com/marketplace/publish?creator=return" \
  VUTADEX_MARKETPLACE_CHECKOUT_SUCCESS_URL="https://app.vutadex.com/marketplace?checkout=success" \
  VUTADEX_MARKETPLACE_CHECKOUT_CANCEL_URL="https://app.vutadex.com/marketplace?checkout=cancelled" \
  VUTADEX_MARKETPLACE_PLATFORM_FEE_BPS="1500"
```

Recommended Stripe webhook targets:
- `POST /api/marketplace/webhook` for marketplace checkout/account events
- `POST /api/billing/webhook` for subscription billing events

Because Fly secrets are write-only, prefer updating `.env.staging` / `.env.production`
and re-importing with `task secrets:fly:import` instead of trying to treat Fly as
the canonical source for configuration.

Then redeploy app:

```bash
fly deploy
```

## Remove

```bash
cd infra
npx sst remove --stage production
```
