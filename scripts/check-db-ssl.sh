#!/usr/bin/env bash
# Run this ON THE VM from your tulip project directory.
# Example: cd ~/tulip && bash scripts/check-db-ssl.sh
set -euo pipefail

echo "=== Tulip DB / SSL diagnostics ==="
echo

PROJECT_DIR="${1:-.}"
cd "$PROJECT_DIR"

echo "--- 1. Certificate files on disk ---"
for f in server.crt server.key; do
  if [[ -f "$f" ]]; then
    size=$(stat -c%s "$f" 2>/dev/null || stat -f%z "$f" 2>/dev/null || echo "?")
    echo "OK  $f exists (${size} bytes)"
  else
    echo "MISSING  $f  (generate with: bash scripts/generate-db-certs.sh --force)"
  fi
done

if [[ -f server.crt ]]; then
  openssl x509 -in server.crt -noout -subject -dates 2>/dev/null || echo "WARN  server.crt is not a valid PEM certificate"
  if openssl x509 -in server.crt -noout -checkend 0 2>/dev/null; then
    echo "OK  server.crt is currently valid"
  else
    echo "FAIL  server.crt is EXPIRED — run: bash scripts/generate-db-certs.sh --force"
  fi
fi

if [[ -f server.key && ! -r server.key ]]; then
  echo "NOTE  server.key is not readable by $(whoami) (owned by mac) — that is OK if Postgres container can read it"
fi
echo

echo "--- 2. App container environment ---"
if docker ps --format '{{.Names}}' | grep -q '^tulip_app$'; then
  docker exec tulip_app env | grep -E '^(ENVIRONMENT|POSTGRES_HOST|POSTGRES_USER|POSTGRES_DB|AUTH_ORIGIN)=' | sort
else
  echo "WARN  tulip_app container is not running"
fi
echo

echo "--- 3. Can the app reach Postgres? ---"
if docker ps --format '{{.Names}}' | grep -q '^tulip_app$'; then
  if docker exec tulip_app ping -c 1 tulip_db >/dev/null 2>&1; then
    echo "OK  tulip_app can ping tulip_db"
  else
    echo "FAIL  tulip_app cannot resolve/ping tulip_db — check POSTGRES_HOST=tulip_db"
  fi
else
  echo "SKIP  tulip_app not running"
fi
echo

echo "--- 4. Postgres SSL settings (inside container) ---"
if docker ps --format '{{.Names}}' | grep -q '^tulip_db$'; then
  docker exec tulip_db psql -U "${POSTGRES_USER:-publicacionestulip}" -d "${POSTGRES_DB:-tulipdb}" -c "SHOW ssl;" 2>/dev/null || \
    echo "WARN  could not query SHOW ssl (check POSTGRES_USER / POSTGRES_DB)"
  echo
  echo "pg_hba.conf (last 10 lines):"
  docker exec tulip_db tail -10 /var/lib/postgresql/data/pg_hba.conf 2>/dev/null || true
else
  echo "WARN  tulip_db container is not running"
fi
echo

echo "--- 5. Recent Postgres SSL / auth errors ---"
if docker ps --format '{{.Names}}' | grep -q '^tulip_db$'; then
  docker logs tulip_db --tail 30 2>&1 | grep -E 'SSL|ssl|FATAL|bad certificate|no encryption' || echo "OK  no recent SSL errors in last 30 log lines"
else
  echo "SKIP  tulip_db not running"
fi
echo

echo "--- 6. API smoke test (from inside app container) ---"
if docker ps --format '{{.Names}}' | grep -q '^tulip_app$'; then
  docker exec tulip_app sh -c '
    ORIGIN="${AUTH_ORIGIN:-http://localhost:8080}"
    TOKEN="${AUTH_TOKEN}"
    if [ -z "$TOKEN" ]; then echo "SKIP  AUTH_TOKEN not set"; exit 0; fi
    wget -qO- "$ORIGIN/api/book/letter" \
      --header="Authorization: Bearer $TOKEN" \
      --header="Origin: $ORIGIN" 2>&1 | head -c 200
    echo
  ' || echo "FAIL  API request failed — DB or SSL issue likely"
else
  echo "SKIP  tulip_app not running"
fi
echo

echo "=== Done ==="
echo "Expected prod values:"
echo "  ENVIRONMENT=docker"
echo "  POSTGRES_HOST=tulip_db"
echo "After updating .env.docker: docker compose -f docker-compose.prod.yml up -d --build"
