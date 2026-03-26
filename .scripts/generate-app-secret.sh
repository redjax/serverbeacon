#!/usr/bin/env bash
set -euo pipefail

if ! command -v openssl &>/dev/null; then
  echo "[ERROR] openssl is not installed" 2>&1
  exit 1
fi

SECRET="$(openssl rand -base64 32)"

echo ""
echo "Serverbeacon app secret:"
echo "$SECRET"
echo ""
echo "!! BACK THIS KEY UP SOMEWHERE !!"
echo "If you lose it, you will not be able to retrieve data from the database."
echo ""
echo "Set an env var APP_SECRET=<your-secret> to provide the secret to the app."
echo "Do not store this secret in plaintext."
echo ""
