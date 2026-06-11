#!/usr/bin/env bash
# Production deploy on the VM: pull, rebuild, start, prune build cache, health check.
#
# Usage (on the VM, as a user that can run docker — usually with sudo):
#   cd /home/mac/tulip
#   sudo -u mac -H git -C /home/mac/tulip pull    # pull code first
#   bash scripts/vm-deploy.sh
#
# Or from mac user:
#   cd /home/mac/tulip && bash scripts/vm-deploy.sh
set -euo pipefail

project_dir="${PROJECT_DIR:-/home/mac/tulip}"
compose_file="docker-compose.prod.yml"

cd "$project_dir"

echo "=== Tulip production deploy ==="
echo "Project: $project_dir"
echo

echo "--- Disk before deploy ---"
df -h / | tail -1
echo

echo "--- Building and starting containers ---"
sudo docker compose -f "$compose_file" up -d --build

echo
echo "--- Waiting for services ---"
sleep 5

echo "--- Container status ---"
sudo docker compose -f "$compose_file" ps

echo
echo "--- App logs (last 5 lines) ---"
sudo docker logs tulip_app --tail 5 2>&1 || true

echo
echo "--- Postgres logs (last 5 lines, excluding collation noise) ---"
sudo docker logs tulip_db --tail 10 2>&1 | grep -v "collation version" | tail -5 || true

echo
echo "--- Pruning build cache (post-deploy) ---"
sudo docker builder prune -f

echo
echo "--- Disk after deploy ---"
df -h / | tail -1

use_pct=$(df / | tail -1 | awk '{print $5}' | tr -d '%')
if [[ "$use_pct" -ge 80 ]]; then
  echo "WARN  Disk is ${use_pct}% full — consider: bash scripts/vm-maintenance.sh"
fi

echo
echo "--- Quick checks ---"
ss -tlnp 2>/dev/null | grep 5432 && echo "OK  Postgres on 127.0.0.1:5432 (TablePlus tunnel)" || echo "WARN  Postgres port not found"
ss -tlnp 2>/dev/null | grep 8080 && echo "OK  App on :8080" || echo "WARN  App port not found"

echo
echo "Deploy complete."
echo "Monthly cleanup: bash scripts/vm-maintenance.sh"
