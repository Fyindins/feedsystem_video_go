#!/usr/bin/env sh
set -eu

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
cleanup() {
  if [ -n "${BACKEND_PID:-}" ]; then
    echo "[start.sh] Stopping backend (pid=$BACKEND_PID)"
    kill "$BACKEND_PID" >/dev/null 2>&1 || true
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

start_backend_fg() {
  echo "[start.sh] Starting backend"
  cd "$BACKEND_DIR"
  go run ./cmd
}

start_frontend_fg() {
  if [ "$FRONTEND_INSTALL" = "1" ] || { [ "$FRONTEND_INSTALL" = "auto" ] && [ ! -d "$FRONTEND_DIR/node_modules" ]; }; then
    echo "[start.sh] Installing frontend deps"
    (cd "$FRONTEND_DIR" && npm install)
  fi

  echo "[start.sh] Starting frontend (npm run $FRONTEND_SCRIPT)"
  cd "$FRONTEND_DIR"
  npm run "$FRONTEND_SCRIPT"
}

if [ "$START_REDIS" = "1" ] && [ "$START_BACKEND" = "1" ]; then
  start_redis
fi

if [ "$START_BACKEND" = "1" ] && [ "$START_FRONTEND" = "1" ]; then
  start_backend_bg
  start_frontend_fg
elif [ "$START_BACKEND" = "1" ]; then
  start_backend_fg
elif [ "$START_FRONTEND" = "1" ]; then
  start_frontend_fg
else
  echo "[start.sh] Nothing to start. Set START_BACKEND=1 and/or START_FRONTEND=1."
fi
