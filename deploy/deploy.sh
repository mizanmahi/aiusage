#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 ]]; then
  echo "usage: $0 <version>" >&2
  exit 1
fi

version="${1#v}"
if [[ ! "$version" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "version must be semver, for example 0.4.0" >&2
  exit 1
fi

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
env_file="$root_dir/.env"
image="mizanph/aiusage:$version"

if [[ ! -f "$env_file" ]]; then
  echo "missing $env_file" >&2
  exit 1
fi

cd "$root_dir"
export IMAGE_TAG="$version"

docker pull "$image"
docker run --rm --env-file "$env_file" --entrypoint sh "$image" \
  -c 'goose -dir /app/migrations postgres "$DATABASE_URL" up'
docker compose -f compose.yml up -d --no-deps --force-recreate app

for _ in {1..20}; do
  if curl --fail --silent --show-error http://127.0.0.1:8080/health >/dev/null; then
    echo "deployed $image"
    exit 0
  fi
  sleep 1
done

docker compose -f compose.yml logs --tail=100 app >&2
echo "health check failed for $image" >&2
exit 1
