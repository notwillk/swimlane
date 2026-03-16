#!/usr/bin/env sh
# Exit 0 if all .go files are formatted; exit 1 and list files otherwise.
set -e
UNFMT=$(find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -l 2>/dev/null || true)
if [ -n "$UNFMT" ]; then
  echo "$UNFMT"
  exit 1
fi
