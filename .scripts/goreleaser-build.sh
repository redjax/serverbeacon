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
ls -lh dist/*/serverbeacon* 2>/dev/null || echo "No binaries found"
echo ""
echo "Build artifacts:"
find dist -name "serverbeacon*" -type f -exec echo "  {}" \;