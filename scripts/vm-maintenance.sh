#!/usr/bin/env bash
# Safe VM disk + Docker cleanup. Keeps running containers and database volumes.
#
# Usage (on the VM):
#   cd /home/mac/tulip
#   bash scripts/vm-maintenance.sh          # monthly / when disk is getting full
#   bash scripts/vm-maintenance.sh --report # disk/memory report only, no changes
#
# Reclaims: build cache, dangling images/containers, old journal logs, apt cache.
# Never deletes: pgdata volume, server.crt/key, authorized_keys, project files.
set -euo pipefail

report_only=false
for arg in "$@"; do
  if [[ "$arg" == "--report" ]]; then
    report_only=true
  fi
done

project_dir="${PROJECT_DIR:-/home/mac/tulip}"
compose_file="docker-compose.prod.yml"

section() {
  echo
  echo "=== $1 ==="
}

warn_disk() {
  local use_pct="${1:-0}"
  if [[ "$use_pct" -ge 90 ]]; then
    echo "WARN  Disk is ${use_pct}% full — run cleanup soon or resize the boot disk."
  elif [[ "$use_pct" -ge 80 ]]; then
    echo "NOTE  Disk is ${use_pct}% full — consider monthly cleanup or a larger disk."
  else
    echo "OK    Disk usage is healthy (${use_pct}%)."
  fi
}

section "Tulip VM maintenance"
echo "Project: $project_dir"
echo "Mode:    $(if $report_only; then echo 'report only'; else echo 'cleanup'; fi)"

section "Disk (before)"
df -h / | tail -1
use_pct=$(df / | tail -1 | awk '{print $5}' | tr -d '%')
warn_disk "$use_pct"

section "Docker usage (before)"
sudo docker system df || true

section "Memory"
free -h | head -2
sudo docker stats --no-stream 2>/dev/null || echo "SKIP  containers not running"

if [[ -d "$project_dir" ]]; then
  section "Container status"
  sudo docker compose -f "$project_dir/$compose_file" ps || true
  ss -tlnp 2>/dev/null | grep 5432 || echo "NOTE  Postgres not listening on 127.0.0.1:5432"
else
  echo "WARN  Project dir not found: $project_dir"
fi

if $report_only; then
  echo
  echo "Report only — no changes made."
  exit 0
fi

section "Cleanup"

echo "→ Docker: remove stopped containers, unused networks, dangling images..."
sudo docker system prune -f

echo "→ Docker: remove unused images (keeps images used by running containers)..."
sudo docker image prune -a -f

echo "→ Docker: remove build cache (safe after deploys; ~500MB typical)..."
sudo docker builder prune -f

echo "→ Journal: keep last 50MB of system logs..."
sudo journalctl --vacuum-size=50M 2>/dev/null || true

echo "→ Apt cache..."
sudo apt-get clean -y 2>/dev/null || true
sudo apt-get autoremove -y 2>/dev/null || true

# Optional old cert backups from generate-db-certs.sh --force
if [[ -d "$project_dir" ]]; then
  shopt -s nullglob
  backups=("$project_dir"/certs-backup-*)
  if [[ ${#backups[@]} -gt 0 ]]; then
    echo "→ Removing ${#backups[@]} old cert backup folder(s)..."
    sudo rm -rf "${backups[@]}"
  fi
  shopt -u nullglob
fi

section "Disk (after)"
df -h / | tail -1
use_pct=$(df / | tail -1 | awk '{print $5}' | tr -d '%')
warn_disk "$use_pct"

section "Docker usage (after)"
sudo docker system df || true

section "Health check"
if [[ -d "$project_dir" ]]; then
  sudo docker compose -f "$project_dir/$compose_file" ps
  if sudo docker ps --format '{{.Names}}' | grep -q '^tulip_db$'; then
    sudo docker logs tulip_db --tail 3 2>&1 | grep -v "collation version" | tail -3 || true
  fi
  if sudo docker ps --format '{{.Names}}' | grep -q '^tulip_app$'; then
    sudo docker logs tulip_app --tail 3 2>&1 || true
  fi
fi

echo
echo "Done. Database volumes were NOT pruned."
echo "Tip: run 'bash scripts/vm-maintenance.sh --report' anytime to check disk/memory."
