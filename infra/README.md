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
2. SES sender identity verified in your AWS region.
3. Cloudflare API token and account ID:
   - `CLOUDFLARE_API_TOKEN`
   - `CLOUDFLARE_DEFAULT_ACCOUNT_ID`
4. `npm` installed (used to build the Vite app during deploy).

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

## First-time setup

```bash
cd infra
npm install
npx sst install
npx sst secret set EmailApiKey "<strong-random-token>" --stage production
```

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

## Stripe billing env (app backend)

Set these Fly secrets for Stripe checkout + webhook handling:

```bash
fly secrets set \
  VUTADEX_STRIPE_SECRET_KEY="sk_live_..." \
  VUTADEX_STRIPE_WEBHOOK_SECRET="whsec_..." \
  VUTADEX_STRIPE_PRICE_PRO="price_..." \
  VUTADEX_STRIPE_CHECKOUT_SUCCESS_URL="https://app.vutadex.com/team/settings?billing=success" \
  VUTADEX_STRIPE_CHECKOUT_CANCEL_URL="https://app.vutadex.com/team/settings?billing=canceled"
```

Then redeploy app:

```bash
fly deploy
```

## Remove

```bash
cd infra
npx sst remove --stage production
```
