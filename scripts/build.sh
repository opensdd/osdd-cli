#!/bin/bash
set -euo pipefail

# Configuration
APP_NAME="osdd"
VERSION="${VERSION:-}"
OUTPUT_DIR="${OUTPUT_DIR:-out}"
PACKAGE="github.com/opensdd/osdd-cli"
MAIN_PKG="./internal/cmd/osdd"

# Validate VERSION is provided
if [ -z "$VERSION" ]; then
  echo "Error: VERSION environment variable is not set" >&2
  echo "Usage: VERSION=vX.Y.Z ./scripts/build.sh" >&2
  exit 1
fi

# Validate VERSION format (vX.Y.Z or "dev")
if [ "$VERSION" != "dev" ] && ! echo "$VERSION" | grep -qE '^v[0-9]+\.[0-9]+\.[0-9]+$'; then
  echo "Error: VERSION must be 'dev' or in format vX.Y.Z (e.g., v1.0.0)" >&2
  exit 1
fi

# ldflags for version injection (preserve user-provided LDFLAGS)
USER_LDFLAGS="${LDFLAGS:-}"
LDFLAGS="${USER_LDFLAGS} -X ${PACKAGE}/internal/version.Version=${VERSION}"

# Target platforms
TARGETS=(
  "linux:amd64"
  "linux:arm64"
  "darwin:amd64"
  "darwin:arm64"
  "windows:amd64"
  "windows:arm64"
)

mkdir -p "$OUTPUT_DIR"

echo "Building ${APP_NAME} ${VERSION} for: ${TARGETS[*]}"

build_one() {
  local os="$1"
  local arch="$2"
  local ext=""
  if [ "$os" = "windows" ]; then
    ext=".exe"
  fi
  local out="${OUTPUT_DIR}/${APP_NAME}-${os}-${arch}${ext}"
  echo "→ $os/$arch -> $out"
  GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 \
    go build -ldflags "$LDFLAGS" -o "$out" "$MAIN_PKG"
  if [ "$os" != "windows" ]; then
    chmod +x "$out"
  fi
}

# Build in parallel
pids=()
for t in "${TARGETS[@]}"; do
  IFS=: read -r os arch <<<"$t"
  build_one "$os" "$arch" &
  pids+=($!)
done

# Wait for all builds to complete
fail=0
for pid in "${pids[@]}"; do
  if ! wait "$pid"; then
    fail=1
  fi
done

if [ "$fail" -ne 0 ]; then
  echo "One or more builds failed." >&2
  exit 1
fi

echo "All builds completed. Artifacts are in $OUTPUT_DIR"
echo "Examples:"
echo "  $OUTPUT_DIR/${APP_NAME}-darwin-arm64"
echo "  $OUTPUT_DIR/${APP_NAME}-windows-amd64.exe"
