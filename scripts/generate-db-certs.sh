#!/usr/bin/env bash
# Generate a self-signed TLS certificate for Postgres.
# Run ON THE VM in the tulip project directory, then restart compose.
set -euo pipefail

cd "${1:-.}"

if [[ -f server.crt && -f server.key ]]; then
  echo "server.crt and server.key already exist."
  echo "To regenerate, delete them first: rm server.crt server.key"
  openssl x509 -in server.crt -noout -subject -dates
  exit 0
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
echo "Next: docker compose -f docker-compose.prod.yml up -d"
