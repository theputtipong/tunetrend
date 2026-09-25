#!/usr/bin/env bash
# Start the TuneTrend mobile app on an Android emulator.
#
# Usage:
#   ./scripts/start.sh                     # run against http://10.0.2.2:8080 (local backend)
#   ./scripts/start.sh --with-backend      # also start local backend (docker db + go run)
#   ./scripts/start.sh --api https://...   # run against another backend
#   ./scripts/start.sh -d macos            # run on a specific Flutter device
#
# Extra args after `--` are passed straight to `flutter run`.
set -euo pipefail

MOBILE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="$MOBILE_DIR/../backend"

API_BASE_URL="${API_BASE_URL:-}"
DEVICE=""
WITH_BACKEND=false
EMULATOR_ID="${EMULATOR_ID:-}"
FLUTTER_ARGS=()

while [ $# -gt 0 ]; do
  case "$1" in
    --api) API_BASE_URL="$2"; shift 2 ;;
    -d|--device) DEVICE="$2"; shift 2 ;;
    --with-backend) WITH_BACKEND=true; shift ;;
    --emulator) EMULATOR_ID="$2"; shift 2 ;;
    -h|--help) sed -n '2,10p' "$0" | sed 's/^# \{0,1\}//'; exit 0 ;;
    --) shift; FLUTTER_ARGS=("$@"); break ;;
    *) echo "❌ Unknown option: $1" >&2; exit 1 ;;
  esac
done

BACKEND_PID=""
cleanup() {
  if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    echo "🛑 Stopping backend (pid $BACKEND_PID)..."
    kill "$BACKEND_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT

start_backend() {
  if [ ! -f "$BACKEND_DIR/.env" ]; then
    echo "❌ $BACKEND_DIR/.env not found. Run apps/backend/setup.sh and fill in the keys first." >&2
    exit 1
  fi
  if curl -s -m 2 -o /dev/null "http://localhost:8080/"; then
    echo "✅ Backend already running on :8080"
    return
  fi

  echo "🐳 Starting database..."
  (cd "$BACKEND_DIR" && docker compose up -d)

  echo "🚀 Starting backend (logs: $BACKEND_DIR/backend.log)..."
  (cd "$BACKEND_DIR" && exec go run ./cmd/api) > "$BACKEND_DIR/backend.log" 2>&1 &
  BACKEND_PID=$!

  for _ in $(seq 1 60); do
    if curl -s -m 2 -o /dev/null "http://localhost:8080/"; then
      echo "✅ Backend is up on :8080"
      return
    fi
    if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
      echo "❌ Backend exited. Last log lines:" >&2
      tail -20 "$BACKEND_DIR/backend.log" >&2
      exit 1
    fi
    sleep 1
  done
  echo "❌ Backend did not respond on :8080 within 60s. See backend.log" >&2
  exit 1
}

first_android_device() {
  adb devices 2>/dev/null | awk 'NR>1 && $2=="device" {print $1; exit}'
}

ensure_android_device() {
  local id
  id="$(first_android_device)"
  if [ -n "$id" ]; then
    DEVICE="$id"
    return
  fi

  if [ -z "$EMULATOR_ID" ]; then
    EMULATOR_ID="$(flutter emulators 2>/dev/null | awk -F' • ' '$4 ~ /android/ {gsub(/ +$/, "", $1); print $1; exit}')"
  fi
  if [ -z "$EMULATOR_ID" ]; then
    echo "❌ No Android device connected and no emulator found. Create one in Android Studio." >&2
    exit 1
  fi

  echo "📱 Launching emulator $EMULATOR_ID..."
  flutter emulators --launch "$EMULATOR_ID"

  for _ in $(seq 1 120); do
    id="$(first_android_device)"
    if [ -n "$id" ] && [ "$(adb -s "$id" shell getprop sys.boot_completed 2>/dev/null | tr -d '\r')" = "1" ]; then
      DEVICE="$id"
      echo "✅ Emulator ready ($DEVICE)"
      return
    fi
    sleep 2
  done
  echo "❌ Emulator did not finish booting in time." >&2
  exit 1
}

if $WITH_BACKEND; then
  start_backend
fi

if [ -z "$DEVICE" ]; then
  ensure_android_device
fi

cd "$MOBILE_DIR"
flutter pub get

RUN_ARGS=(-d "$DEVICE")
if [ -n "$API_BASE_URL" ]; then
  RUN_ARGS+=(--dart-define=API_BASE_URL="$API_BASE_URL")
fi

echo "▶️  flutter run ${RUN_ARGS[*]} ${FLUTTER_ARGS[*]:-}"
flutter run "${RUN_ARGS[@]}" ${FLUTTER_ARGS[@]+"${FLUTTER_ARGS[@]}"}
