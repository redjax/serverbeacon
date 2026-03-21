#!/usr/bin/env bash
set -euo pipefail

if ! command -v docker >&/dev/null; then
  echo "[ERROR] Docker is not installed." >&2
  exit 1
fi

THIS_DIR="$(dirname "${0}")"
REPO_ROOT=$(realpath -m "${THIS_DIR}/../..")
CWD=$(pwd)

IMG="${SERVERBEACON_IMG_NAME:-serverbeacon}"
TAG="${SERVERBEACON_IMG_TAG:-latest}"

function cleanup() {
  cd "$CWD"
}
trap cleanup EXIT

cmd=(docker build -t "${IMG}:${TAG}" .)

cd "$REPO_ROOT"

echo "Running command:"
echo "  ${cmd[*]}"

if ! "${cmd[@]}"; then
  LAST_EXIT=$?
  echo "[ERROR] Failed building Docker container" >&2
  exit $LAST_EXIT
fi
