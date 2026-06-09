#!/usr/bin/env bash
# Generate a self-signed TLS certificate for Postgres.
# Run ON THE VM in the tulip project directory, then restart compose.
#
# Usage:
#   bash scripts/generate-db-certs.sh          # create only if missing
#   bash scripts/generate-db-certs.sh --force  # backup old certs and regenerate
set -euo pipefail

force=false
project_dir="."
for arg in "$@"; do
  if [[ "$arg" == "--force" ]]; then
    force=true
  elif [[ "$arg" != --* ]]; then
    project_dir="$arg"
  fi
done
cd "$project_dir"

if [[ -f server.crt && -f server.key && "$force" == false ]]; then
  echo "server.crt and server.key already exist."
  openssl x509 -in server.crt -noout -subject -dates
  if ! openssl x509 -in server.crt -noout -checkend 0 2>/dev/null; then
    echo "WARN  Certificate is EXPIRED. Regenerate with: bash scripts/generate-db-certs.sh --force"
  fi
  exit 0
fi

if [[ -f server.crt || -f server.key ]]; then
  backup="certs-backup-$(date +%Y%m%d-%H%M%S)"
  mkdir -p "$backup"
  [[ -f server.crt ]] && cp server.crt "$backup/"
  [[ -f server.key ]] && cp server.key "$backup/"
  echo "Backed up old certs to $backup/"
fi

echo "Generating self-signed certificate (valid 10 years)..."
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout server.key \
  -out server.crt \
  -days 3650 \
  -subj "/CN=tulip_db/O=PublicacionesTulip"

chmod 600 server.key
chmod 644 server.crt

echo "Created server.crt and server.key"
openssl x509 -in server.crt -noout -subject -dates
echo
echo "Next:"
echo "  docker compose -f docker-compose.prod.yml restart tulip_db"
echo "  docker compose -f docker-compose.prod.yml up -d --build"
