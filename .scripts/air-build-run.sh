#!/usr/bin/env bash
set -euo pipefail

if ! command -v air >&/dev/null; then
  echo "[ERROR] air is not installed" >&2
  exit 1
fi

THIS_DIR="$(dirname "$0")"
REPO_ROOT=$(realpath -m "${THIS_DIR}/..")
CWD=$(pwd)

BUILD_TARGET="cmd/serverbeacon/main.go"
OUTPUT_DIR="bin"
BIN_NAME="serverbeacon"
OUTPUT_TARGET="${OUTPUT_DIR}/${BIN_NAME}"

function cleanup() {
  cd "${CWD}"
}
trap cleanup EXIT

cmd=(air --build.cmd "go build -o ${OUTPUT_TARGET} ${BUILD_TARGET}" --build.entrypoint "./${OUTPUT_TARGET}")

echo "Running command:"
echo "  ${cmd[*]}"

if ! "${cmd[@]}" 2>&1; then
  LAST_EXIT=$?
  echo "[ERROR] Failure running air hot reload server." 2>&1
  exit $LAST_EXIT
fi
