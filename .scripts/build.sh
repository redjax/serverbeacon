#!/usr/bin/env bash
set -eo pipefail

## Default to building everything
BUILD_API=true
BUILD_CLI=true

THIS_DIR="$(dirname "$0")"
REPO_ROOT=$(realpath -m "${THIS_DIR}/..")
CWD=$(pwd)

function usage() {
    echo ""
    echo "Usage: ${0} [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Print help menu"
    echo "  --all           Build all apps (default)"
    echo "  --api           Build API app"
    echo "  --cli           Build CLI app"
    echo ""
    echo "Note: You can pass multiple apps to build, like --api --cli"
    echo ""
}

function cleanup() {
    cd "$CWD"
}
trap cleanup EXIT

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
        --help|-h)
            usage
            exit 0
            ;;
        *)
            echo "Unknown option: $1" &>2
            exit 1
            ;;
    esac
done

## Get OS and ARCH
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
echo "Target: ${OS}/${ARCH}"

## Create bin directory
mkdir -p bin

## Helper function to build a binary
build_binary() {
    local entrypoint=$1
    local binary_name=$2
    local output="bin/${binary_name}"
    
    if [ "$OS" = "windows" ]; then
        output="bin/${binary_name}.exe"
    fi
    
    echo "Building ${binary_name} → ${output}"
    go build -ldflags="-s -w" -o "${output}" "./cmd/${entrypoint}"
    
    if [ -f "${output}" ]; then
        echo "✓ ${binary_name} built successfully"
        ls -lh "${output}"
    else
        echo "✗ ${binary_name} build failed"
        exit 1
    fi
    echo ""
}

## Build selected binaries
if [ "$BUILD_API" = true ]; then
    build_binary "api" "serverbeacon-api"
fi

if [ "$BUILD_CLI" = true ]; then
    build_binary "cli" "serverbeacon"
fi

echo "All builds complete"
echo "Binaries in bin/:"
ls -lh bin/serverbeacon* 2>/dev/null || echo "No binaries found"
