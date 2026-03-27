#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
ENV_FILE="${1:-$ROOT_DIR/.env.local}"
UPSERTER="$ROOT_DIR/scripts/env/upsert_env.py"

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

python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_ENV" "development"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_DATABASE_PATH" "./data/microdote.db"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_APP_ORIGIN" "http://localhost:3000"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_MARKETING_ORIGIN" "http://localhost:4173"
python3 "$UPSERTER" "$ENV_FILE" "VUTADEX_COOKIE_SECURE" "false"
python3 "$UPSERTER" "$ENV_FILE" "VITE_APP_ORIGIN" "http://localhost:3000"

echo "Initialized local SQLite env at $ENV_FILE"
