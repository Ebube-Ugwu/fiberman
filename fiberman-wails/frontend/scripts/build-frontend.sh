#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
WRAPPER_FRONTEND_DIR="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
WAILS_DIR="$(cd -- "${WRAPPER_FRONTEND_DIR}/.." && pwd)"
REPO_DIR="$(cd -- "${WAILS_DIR}/.." && pwd)"
ANGULAR_DIR="${REPO_DIR}/fiberman-frontend"
DIST_DIR="${WRAPPER_FRONTEND_DIR}/dist"

cd "${ANGULAR_DIR}"
npm run build

rm -rf "${DIST_DIR}"
mkdir -p "${DIST_DIR}"
cp -r "${ANGULAR_DIR}/dist/fiberman-frontend/." "${DIST_DIR}/"
