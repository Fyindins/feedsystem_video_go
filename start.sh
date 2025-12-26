#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)"
cd "$ROOT_DIR"

BACKEND_DIR="${BACKEND_DIR:-$ROOT_DIR/backend}"
FRONTEND_DIR="${FRONTEND_DIR:-$ROOT_DIR/frontend}"
RUN_DIR="${RUN_DIR:-$BACKEND_DIR/.run}"

START_REDIS="${START_REDIS:-1}"
START_BACKEND="${START_BACKEND:-1}"
START_FRONTEND="${START_FRONTEND:-1}"

# frontend:
# - FRONTEND_INSTALL=auto|1|0 (default auto, only when node_modules missing)
# - FRONTEND_SCRIPT=dev|preview|build (default dev)
FRONTEND_INSTALL="${FRONTEND_INSTALL:-auto}"
FRONTEND_SCRIPT="${FRONTEND_SCRIPT:-dev}"

# redis:
# - If you already run Redis elsewhere, this will detect it (when redis-cli exists) and skip.
REDIS_HOST="${REDIS_HOST:-127.0.0.1}"
REDIS_PORT="${REDIS_PORT:-6379}"
REDIS_CONF="${REDIS_CONF:-}"

require_dir() {
  if [ ! -d "$2" ]; then
    echo "[start.sh] $1 dir not found: $2"
    exit 1
  fi
}

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "[start.sh] command not found: $1"
    exit 1
  fi
}

BACKEND_PID=""
FRONTEND_PID=""
cleanup() {
  set +e
  if [ -n "${FRONTEND_PID:-}" ]; then
    echo "[start.sh] Stopping frontend (pid=$FRONTEND_PID)"
    kill "$FRONTEND_PID" >/dev/null 2>&1 || true
  fi
  if [ -n "${BACKEND_PID:-}" ]; then
    echo "[start.sh] Stopping backend (pid=$BACKEND_PID)"
    kill "$BACKEND_PID" >/dev/null 2>&1 || true
  fi
  # On Windows Git Bash, child processes sometimes don't receive signals reliably.
  # Fall back to taskkill when available.
  if command -v taskkill >/dev/null 2>&1; then
    if [ -n "${FRONTEND_PID:-}" ]; then
      taskkill //PID "$FRONTEND_PID" //T //F >/dev/null 2>&1 || true
    fi
    if [ -n "${BACKEND_PID:-}" ]; then
      taskkill //PID "$BACKEND_PID" //T //F >/dev/null 2>&1 || true
    fi
  fi
}
trap cleanup INT TERM EXIT

if [ "$START_BACKEND" = "1" ]; then
  require_dir backend "$BACKEND_DIR"
  mkdir -p "$RUN_DIR"
  require_cmd go
fi

if [ "$START_FRONTEND" = "1" ]; then
  require_dir frontend "$FRONTEND_DIR"
  mkdir -p "$RUN_DIR"
  require_cmd npm
fi

start_redis() {
  if ! command -v redis-server >/dev/null 2>&1; then
    echo "[start.sh] redis-server not found; skip starting Redis"
    return 0
  fi

  if command -v redis-cli >/dev/null 2>&1; then
    if redis-cli -h "$REDIS_HOST" -p "$REDIS_PORT" ping >/dev/null 2>&1; then
      echo "[start.sh] Redis already running at $REDIS_HOST:$REDIS_PORT"
      return 0
    fi
  fi

  echo "[start.sh] Starting Redis at $REDIS_HOST:$REDIS_PORT"
  if [ -n "$REDIS_CONF" ]; then
    nohup redis-server "$REDIS_CONF" >"$RUN_DIR/redis.log" 2>&1 &
  else
    nohup redis-server --bind "$REDIS_HOST" --port "$REDIS_PORT" >"$RUN_DIR/redis.log" 2>&1 &
  fi
  echo $! >"$RUN_DIR/redis.pid"
}

start_backend_bg() {
  echo "[start.sh] Starting backend (background)"
  (cd "$BACKEND_DIR" && go run ./cmd) &
  BACKEND_PID=$!
  echo "$BACKEND_PID" >"$RUN_DIR/backend.pid"
  echo "[start.sh] Backend PID: $BACKEND_PID"
}

start_frontend_bg() {
  if [ "$FRONTEND_INSTALL" = "1" ] || { [ "$FRONTEND_INSTALL" = "auto" ] && [ ! -d "$FRONTEND_DIR/node_modules" ]; }; then
    echo "[start.sh] Installing frontend deps"
    (cd "$FRONTEND_DIR" && npm install)
  fi

  echo "[start.sh] Starting frontend (background, npm run $FRONTEND_SCRIPT)"
  (cd "$FRONTEND_DIR" && npm run "$FRONTEND_SCRIPT") &
  FRONTEND_PID=$!
  echo "$FRONTEND_PID" >"$RUN_DIR/frontend.pid"
  echo "[start.sh] Frontend PID: $FRONTEND_PID"
}

if [ "$START_REDIS" = "1" ] && [ "$START_BACKEND" = "1" ]; then
  start_redis
fi

if [ "$START_BACKEND" = "1" ]; then
  start_backend_bg
fi

if [ "$START_FRONTEND" = "1" ]; then
  start_frontend_bg
fi

if [ "$START_BACKEND" = "1" ] || [ "$START_FRONTEND" = "1" ]; then
  echo "[start.sh] Press Ctrl+C to stop."
  wait
else
  echo "[start.sh] Nothing to start. Set START_BACKEND=1 and/or START_FRONTEND=1."
fi
