#!/usr/bin/env bash
set -euo pipefail

if ! command -v docker >&/dev/null; then
  echo "[ERROR] Docker is not installed." 2>&1
  exit 1
fi

declare -a CONTAINER_ARGS=()
PORT=18080
NAME="serverbeacon"
IMG="serverbeacon"
TAG="latest"
## Interactive disabled, Dockerfile is distroless
#  and does not have an interactive shell.
# INTERACTIVE=false
DRY_RUN=false

cmd=(docker run)

function usage() {
  echo ""
  echo "Usage: ${0} [OPTIONS] [-- [CONTAINER ARGS...]]"
  echo ""
  echo "Options:"
  echo "  -p, --port         <int>     Set host port"
  echo "  -n, --name         <string>  Set container's name when running"
  echo "  -h, --help                   Print this help menu"
  # echo "  -i, --interactive            Run container in interactive mode"
  echo "  -i, --img          <string>  Set the container image name to run"
  echo "  -t, --tag          <string>  Set the container image's tag to target"
  echo "  --dry-run                    Enable dry-run mode, state actions without taking them"
  echo ""
  echo "Examples:"
  echo "  ${0} -i serverbeacon -t latest -- start rest-api"
  echo "  ${0} -i serverbeacon-api -t latest"
  echo ""
}

while [[ $# -gt 0 ]]; do
  case $1 in
    -p|--port)
      PORT="${2}"
      shift 2
      ;;
    -n|--name)
      NAME="${2}"
      shift 2
      ;;
    # -i|--interactive)
    #   INTERACTIVE=true
    #   shift
    #   ;;
    -i|--img)
      IMG="${2}"
      shift 2
      ;;
    -t|--tag)
      TAG="${2}"
      shift 2
      ;;
    --dry-run)
      DRY_RUN=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      CONTAINER_ARGS=("$@")
      break
      ;;
    *)
      echo "[ERROR] Invalid arg: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ -z "${PORT}" ]]; then
  echo "[ERROR] Missing port value" >&2
  exit 1
fi

if [[ -z "${NAME}" ]]; then
  echo "[ERROR] Missing container name" >&2
  exit 1
fi

if [[ -z "${IMG}" ]]; then
  echo "[ERROR] Missing container image" >&2
  exit 1
fi

if [[ -z "${TAG}" ]]; then
  echo "[ERROR] Missing container image tag" >&2
  exit 1
fi

cmd+=(-d -p "${PORT}:18080" --name "${NAME}")

if [[ "${#CONTAINER_ARGS[@]}" -gt 0 ]]; then
  cmd+=("${IMG}:${TAG}")
  cmd+=("${CONTAINER_ARGS[@]}")
else
  cmd+=("${IMG}:${TAG}")
fi

echo "Running command:"
echo "  ${cmd[*]}"

if ! "${cmd[@]}"; then
  echo "[ERROR] Failed to run Docker container '${IMG}:${TAG}'" >&2
  exit 1
fi
