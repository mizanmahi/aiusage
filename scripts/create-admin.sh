#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

if [[ -f "$ROOT_DIR/server/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT_DIR/server/.env"
  set +a
fi

if [[ -z "${DATABASE_URL:-}" ]]; then
  echo "DATABASE_URL is required. Export it or add it to server/.env." >&2
  exit 1
fi

ADMIN_EMAIL="${ADMIN_EMAIL:-mizan@example.com}"
ADMIN_NAME="${ADMIN_NAME:-Mizan}"
API_KEY="${AIUSAGE_ADMIN_API_KEY:-ak_$(openssl rand -hex 24)}"

psql "$DATABASE_URL" \
  -v email="$ADMIN_EMAIL" \
  -v name="$ADMIN_NAME" \
  -v api_key="$API_KEY" <<'SQL'
INSERT INTO users (email, name, api_key_hash, is_admin)
VALUES (
  :'email',
  :'name',
  encode(digest(:'api_key', 'sha256'), 'hex'),
  TRUE
);
SQL

cat <<EOF
Admin user created.

Email: $ADMIN_EMAIL
Name:  $ADMIN_NAME
API key:
$API_KEY

Save this key now. It is not stored in plain text.
EOF
