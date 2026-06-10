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
  for f in server.crt server.key; do
    if [[ ! -f "$f" ]]; then
      continue
    fi
    if cp "$f" "$backup/" 2>/dev/null; then
      continue
    fi
    if command -v sudo >/dev/null && sudo cp "$f" "$backup/" 2>/dev/null; then
      echo "NOTE  backed up $f using sudo (file was not readable as $(whoami))"
      continue
    fi
    echo "WARN  could not backup $f — it will be overwritten"
  done
  echo "Backed up old certs to $backup/ (where permissions allowed)"
fi

echo "Generating self-signed certificate (valid 10 years)..."
openssl req -x509 -newkey rsa:2048 -nodes \
  -keyout server.key \
  -out server.crt \
  -days 3650 \
  -subj "/CN=tulip_db/O=PublicacionesTulip"

chmod 600 server.key
chmod 644 server.crt
# Postgres Docker (uid 999) must own bind-mounted TLS files — not mac, not root-only 600.
postgres_uid=999
if command -v sudo >/dev/null; then
  sudo chown "${postgres_uid}:${postgres_uid}" server.key server.crt
else
  chown "${postgres_uid}:${postgres_uid}" server.key server.crt 2>/dev/null || \
    echo "WARN  run: sudo chown 999:999 server.key server.crt"
fi

echo "Created server.crt and server.key"
openssl x509 -in server.crt -noout -subject -dates
echo
echo "Next:"
echo "  docker compose -f docker-compose.prod.yml restart tulip_db"
echo "  docker compose -f docker-compose.prod.yml up -d --build"
