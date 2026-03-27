#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${1:?usage: fly-secrets-import.sh <env_file> <fly_app>}"
FLY_APP="${2:?fly app name is required}"

if ! command -v fly >/dev/null 2>&1; then
  echo "flyctl is required. Install it from https://fly.io/docs/flyctl/" >&2
  exit 1
fi

if [ ! -f "$ENV_FILE" ]; then
  echo "Env file not found: $ENV_FILE" >&2
  exit 1
fi

python3 - "$ENV_FILE" <<'PY' | fly secrets import -a "$FLY_APP"
import sys
from pathlib import Path

env_file = Path(sys.argv[1])
skip = {
    "VUTADEX_DATABASE_PATH",
    "VITE_APP_ORIGIN",
}

for raw in env_file.read_text().splitlines():
    line = raw.strip()
    if not line or line.startswith("#") or "=" not in line:
        continue
    key = line.split("=", 1)[0].strip()
    if key in skip:
        continue
    if key.startswith("VUTADEX_") or key == "PORT":
        print(raw)
PY
