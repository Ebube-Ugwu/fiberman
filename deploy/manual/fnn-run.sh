#!/bin/sh
set -eu

FIBER_HOME_DIR="${FIBER_HOME:-/home/holo-1/fiber-data}"
FNN_CONFIG_DIR="${FNN_CONFIG_DIR:-/home/holo-1/fnn-config}"
NETWORK="${FNN_NETWORK:-testnet}"
CONFIG_PATH="${FIBER_CONFIG:-$FIBER_HOME_DIR/config.yml}"
SOURCE_CONFIG="$FNN_CONFIG_DIR/$NETWORK/config.yml"
KEY_PATH="$FIBER_HOME_DIR/ckb/key"

mkdir -p "$FIBER_HOME_DIR/ckb"

if [ ! -f "$KEY_PATH" ]; then
  head -c 32 /dev/urandom | od -An -tx1 | tr -d ' \n' > "$KEY_PATH"
  printf '\n' >> "$KEY_PATH"
  chmod 600 "$KEY_PATH"
fi

if [ ! -f "$CONFIG_PATH" ]; then
  cp "$SOURCE_CONFIG" "$CONFIG_PATH"
fi

exec /usr/local/bin/fnn -d "$FIBER_HOME_DIR" -c "$CONFIG_PATH"
