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
PLATFORM="${SERVERBEACON_PLATFORM:-auto}"

## Default to building everything
BUILD_API=true
BUILD_CLI=true

function usage() {
    echo ""
    echo "Usage: ${0} [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --all           Build API + CLI images (default)"
    echo "  --api           Build API image only"  
    echo "  --cli           Build CLI image only"
    echo "  -p, --platform  Platform to build for, i.e. linux/amd64"
    echo "  -t, --tag       Docker tag (default: latest)"
    echo ""
}

function detect_platform() {
    local arch

    arch="$(uname -m)"
    case "$arch" in
        x86_64|amd64)
            echo "linux/amd64"
            ;;
        aarch64|arm64)
            echo "linux/arm64"
            ;;
        armv7l|armv7)
            echo "linux/arm/v7"
            ;;
        armv6l|armv6)
            echo "linux/arm/v6"
            ;;
        *)
            echo "[ERROR] Unsupported architecture: ${arch}" >&2
            exit 1
            ;;
    esac
}

## Parse args
while [[ $# -gt 0 ]]; do
    case $1 in
        --all)
            BUILD_API=true
            BUILD_CLI=true
            shift
            ;;
        --api)
            BUILD_API=true
            BUILD_CLI=false
            shift
            ;;
        --cli)
            BUILD_API=false
            BUILD_CLI=true
            shift
            ;;
        -t|--tag)
            TAG="$2"
            shift 2
            ;;
        -p|--platform)
            PLATFORM="$2"
            shift 2
            ;;
        --help|-h)
            usage
            exit 0
            ;;
        -*)
            echo "Unknown option: $1" >&2
            exit 1
            ;;
        *)
            TAG="$1"
            shift
            ;;
    esac
done

if [[ "$PLATFORM" == "auto" ]] || [[ -z "$PLATFORM" ]]; then
    PLATFORM="$(detect_platform)"
fi

function cleanup() {
  cd "$CWD"
}
trap cleanup EXIT

cd "$REPO_ROOT"

echo "Building container for platform: ${PLATFORM}"

## Build API
if [ "$BUILD_API" = true ]; then
    echo "Building API image: ${IMG}-api:${TAG}"

    if ! docker build \
        --platform "${PLATFORM}" \
        --target serverbeacon-api \
        -t "${IMG}-api:${TAG}" .; then
        
        echo "[ERROR] Failed building API image" >&2
        exit 1
    fi

    echo "API image built: ${IMG}-api:${TAG}"
fi

## Build CLI
if [ "$BUILD_CLI" = true ]; then
    echo ""
    echo "Building CLI image: ${IMG}:${TAG}"
    
    if ! docker build \
        --platform "${PLATFORM}" \
        --target serverbeacon \
        -t "${IMG}:${TAG}" .; then
        
        echo "[ERROR] Failed building CLI image" >&2
        exit 1
    fi

    echo "CLI image built: ${IMG}:${TAG}"
fi

echo ""
echo "Finished building images"
echo "API: ${IMG}-api:${TAG}"
[ "$BUILD_CLI" = true ] && echo "CLI: ${IMG}:${TAG}"
