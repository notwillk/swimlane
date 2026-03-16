#!/usr/bin/env sh
# Install swimlane from GitHub Releases.
# Usage: curl -sSL https://raw.githubusercontent.com/notwillk/swimlane/main/scripts/install.sh | sh
# Or: curl -sSL ... | sh -s -- v1.0.0

set -e

REPO="notwillk/swimlane"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY="swimlane"

detect_os_arch() {
  OS=$(uname -s)
  ARCH=$(uname -m)
  case "$OS" in
    Linux)  OS=linux ;;
    Darwin) OS=darwin ;;
    *)
      echo "Unsupported OS: $OS. Supported: linux, darwin." >&2
      exit 1
      ;;
  esac
  case "$ARCH" in
    x86_64|amd64) ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *)
      echo "Unsupported architecture: $ARCH. Supported: amd64, arm64." >&2
      exit 1
      ;;
  esac
}

main() {
  detect_os_arch

  VERSION="${1:-latest}"
  if [ "$VERSION" = "latest" ]; then
    URL="https://github.com/${REPO}/releases/latest/download/swimlane_${OS}_${ARCH}.tar.gz"
  else
    URL="https://github.com/${REPO}/releases/download/${VERSION}/swimlane_${OS}_${ARCH}.tar.gz"
  fi

  echo "Installing swimlane to $INSTALL_DIR"
  tmpdir=$(mktemp -d)
  trap 'rm -rf "$tmpdir"' EXIT
  if command -v curl >/dev/null 2>&1; then
    curl -sSLf "$URL" -o "$tmpdir/archive.tar.gz"
  elif command -v wget >/dev/null 2>&1; then
    wget -q "$URL" -O "$tmpdir/archive.tar.gz"
  else
    echo "Need curl or wget to download." >&2
    exit 1
  fi
  tar -xzf "$tmpdir/archive.tar.gz" -C "$tmpdir"
  mkdir -p "$INSTALL_DIR"
  mv "$tmpdir/swimlane" "$INSTALL_DIR/$BINARY"
  chmod +x "$INSTALL_DIR/$BINARY"
  echo "Installed $INSTALL_DIR/$BINARY"
}

main "$@"
