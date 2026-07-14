#!/bin/sh
set -eu

FIBER_HOME_DIR="${FIBER_HOME:-/fiber}"
CONFIG_TEMPLATE="${FIBER_CONFIG_TEMPLATE:-/usr/local/share/fiber/config/testnet/config.yml}"
CONFIG_PATH="${FIBER_CONFIG:-$FIBER_HOME_DIR/config.yml}"
KEY_PATH="$FIBER_HOME_DIR/ckb/key"

mkdir -p "$FIBER_HOME_DIR/ckb"

if [ ! -f "$KEY_PATH" ]; then
  head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' > "$KEY_PATH"
  printf '\n' >> "$KEY_PATH"
  chmod 600 "$KEY_PATH"
  echo "Generated disposable testnet key at $KEY_PATH"
fi

if [ ! -f "$CONFIG_PATH" ]; then
  cp "$CONFIG_TEMPLATE" "$CONFIG_PATH"
  echo "Copied bundled Fiber config to $CONFIG_PATH"
fi

# Keep FNN RPC private to the container so the Java backend can use loopback
# without enabling public biscuit-authenticated RPC.
if grep -q 'listening_addr: "0.0.0.0:8227"' "$CONFIG_PATH"; then
  sed -i 's/listening_addr: "0.0.0.0:8227"/listening_addr: "127.0.0.1:8227"/' "$CONFIG_PATH"
fi

fnn -d "$FIBER_HOME_DIR" -c "$CONFIG_PATH" &
FNN_PID=$!

java -jar /app/fiberman.jar &
APP_PID=$!

cleanup() {
  kill "$APP_PID" "$FNN_PID" 2>/dev/null || true
}

trap cleanup INT TERM

while true; do
  if ! kill -0 "$FNN_PID" 2>/dev/null; then
    wait "$FNN_PID" || true
    kill "$APP_PID" 2>/dev/null || true
    wait "$APP_PID" 2>/dev/null || true
    exit 1
  fi

  if ! kill -0 "$APP_PID" 2>/dev/null; then
    wait "$APP_PID"
    APP_STATUS=$?
    kill "$FNN_PID" 2>/dev/null || true
    wait "$FNN_PID" 2>/dev/null || true
    exit "$APP_STATUS"
  fi

  sleep 2
done
