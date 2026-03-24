#!/usr/bin/env bash
set -euo pipefail

echo "Building serverbeacon with GoReleaser"
echo ""

# Check if goreleaser is installed
if ! command -v goreleaser &> /dev/null; then
    echo "[ERROR] goreleaser is not installed" >&2
    exit 1
fi

## Clean previous builds and build with snapshot mode (no git tag required)
goreleaser build --snapshot --clean

echo ""
echo "[SUCCESS] Build complete! Binaries are in the dist/ directory:"
echo ""

## Debug API build
echo "API binaries:"
find dist -type f -name "serverbeacon-api*" -print

## Debug CLI build
echo ""
echo "CLI binaries:"
find dist -type f -name "serverbeacon" -o -name "serverbeacon_*" -print
