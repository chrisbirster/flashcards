#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${1:?usage: hydrate-turso.sh <env_file> <turso_db> [app_origin] [marketing_origin] [environment]}"
TURSO_DB="${2:?turso db name is required}"
APP_ORIGIN="${3:-https://app.vutadex.com}"
MARKETING_ORIGIN="${4:-https://vutadex.com}"
VUTADEX_ENVIRONMENT="${5:-staging}"
UPSERTER="$ROOT_DIR/scripts/env/upsert_env.py"

if ! command -v turso >/dev/null 2>&1; then
  echo "turso CLI is required. Install it from https://docs.turso.tech/cli/introduction" >&2
  exit 1
fi

mkdir -p "$(dirname "$ENV_FILE")"
touch "$ENV_FILE"

session_secret="$(
  python3 - <<'PY'
import secrets
print(secrets.token_urlsafe(32))
PY
)"

if ! grep -q '^VUTADEX_SESSION_SECRET=' "$ENV_FILE" 2>/dev/null; then
  python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_SESSION_SECRET" "$session_secret"
fi

if db_url="$(turso db show --url "$TURSO_DB" 2>/dev/null)"; then
  :
elif db_url="$(turso db show "$TURSO_DB" --url 2>/dev/null)"; then
  :
else
  echo "Failed to resolve Turso database URL for $TURSO_DB" >&2
  exit 1
fi

db_token="$(turso db tokens create "$TURSO_DB" | tail -n 1 | tr -d '\r')"
db_url="$(printf "%s" "$db_url" | tail -n 1 | tr -d '\r')"

python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_ENV" "$VUTADEX_ENVIRONMENT"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_DATABASE_URL" "$db_url"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_DATABASE_AUTH_TOKEN" "$db_token"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_APP_ORIGIN" "$APP_ORIGIN"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_MARKETING_ORIGIN" "$MARKETING_ORIGIN"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_COOKIE_SECURE" "true"
python3 "$UPSERTER" "$ENV_FILE" "VITE_APP_ORIGIN" "$APP_ORIGIN"

if [[ "$APP_ORIGIN" == https://*.vutadex.com* || "$APP_ORIGIN" == "https://app.vutadex.com" ]]; then
  python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_COOKIE_DOMAIN" ".vutadex.com"
fi

echo "Hydrated Turso connection into $ENV_FILE"
echo "Database URL: $db_url"
