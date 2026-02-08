#!/usr/bin/env bash

set -euo pipefail

STAGE="${1:-baseline}"

case "$STAGE" in
  baseline)
    GO_THRESHOLD="40"
    WEB_THRESHOLD="0"
    ;;
  m3)
    GO_THRESHOLD="55"
    WEB_THRESHOLD="60"
    ;;
  m4)
    GO_THRESHOLD="75"
    WEB_THRESHOLD="80"
    ;;
  release)
    GO_THRESHOLD="95"
    WEB_THRESHOLD="95"
    ;;
  *)
    echo "Unknown coverage stage: $STAGE"
    echo "Valid stages: baseline | m3 | m4 | release"
    exit 1
    ;;
esac

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_COVER_PROFILE="/tmp/flashcards-go.cover"
export GOCACHE="/tmp/go-build"
mkdir -p "$GOCACHE"

echo "== Coverage gate ($STAGE) =="
echo "Backend threshold: ${GO_THRESHOLD}%"
echo "Frontend threshold: ${WEB_THRESHOLD}%"

echo ""
echo "Running backend coverage..."
go test ./... -coverprofile="$GO_COVER_PROFILE" >/tmp/flashcards-go-test.log
GO_COVERAGE="$(go tool cover -func="$GO_COVER_PROFILE" | awk '/^total:/ {gsub("%","",$3); print $3}')"
echo "Backend coverage: ${GO_COVERAGE}%"

if ! awk "BEGIN { exit !(${GO_COVERAGE} >= ${GO_THRESHOLD}) }"; then
  echo "Backend coverage gate failed: ${GO_COVERAGE}% < ${GO_THRESHOLD}%"
  exit 1
fi

if awk "BEGIN { exit !(${WEB_THRESHOLD} <= 0) }"; then
  echo "Frontend coverage gate skipped at this stage."
  exit 0
fi

echo ""
echo "Running frontend coverage..."
if ! npm --prefix "$ROOT_DIR/web" run test:coverage >/tmp/flashcards-web-coverage.log 2>&1; then
  echo "Frontend coverage run failed."
  cat /tmp/flashcards-web-coverage.log
  exit 1
fi

WEB_SUMMARY="$ROOT_DIR/web/coverage/coverage-summary.json"
if [[ ! -f "$WEB_SUMMARY" ]]; then
  echo "Frontend coverage summary not found: $WEB_SUMMARY"
  exit 1
fi

WEB_COVERAGE="$(
  node -e "const s=require(process.argv[1]); console.log(s.total.lines.pct)" "$WEB_SUMMARY"
)"
echo "Frontend coverage: ${WEB_COVERAGE}%"

if ! awk "BEGIN { exit !(${WEB_COVERAGE} >= ${WEB_THRESHOLD}) }"; then
  echo "Frontend coverage gate failed: ${WEB_COVERAGE}% < ${WEB_THRESHOLD}%"
  exit 1
fi

echo ""
echo "Coverage gate passed."
