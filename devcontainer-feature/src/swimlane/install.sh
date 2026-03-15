#!/usr/bin/env bash
set -euo pipefail

V="${VERSION:-latest}"
if [ "$V" != "latest" ] && [ "$V" != "current" ]; then
  V="${V#v}"
fi

# Install swimlane via upstream installer, honoring VERSION option
curl -fsSL https://raw.githubusercontent.com/notwillk/swimlane/main/scripts/install.sh | sh -s -- "$V"

swimlane --help >/dev/null || true
