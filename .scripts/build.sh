#!/bin/bash
set -eo pipefail

echo "Building serverbeacon"

## Get OS and ARCH
OS=$(go env GOOS)
ARCH=$(go env GOARCH)

echo "Target: ${OS}/${ARCH}"

## Create dist directory
mkdir -p dist

## Set output name based on OS
OUTPUT="dist/serverbeacon"
if [ "$OS" = "windows" ]; then
    OUTPUT="dist/serverbeacon.exe"
fi

## Build
echo "Building to ${OUTPUT}"
go build -ldflags="-s -w" -o "${OUTPUT}" ./cmd/serverbeacon

## Check if build was successful
if [ -f "${OUTPUT}" ]; then
    echo "Build successful!"
    echo "Binary: ${OUTPUT}"
    ls -lh "${OUTPUT}"
else
    echo "Build failed!"
    exit 1
fi
