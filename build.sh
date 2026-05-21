#!/bin/bash
set -e

ROOT="$(cd "$(dirname "$(readlink -f "$0")")" && pwd)"
TMP="$ROOT/.tmp"
PKG_NAME="simple_nas"
VERSION="${VERSION:-$(date +%Y%m%d)}"
OUTPUT="${ROOT}/${PKG_NAME}-${VERSION}.tgz"

echo "=== Building NAS ${VERSION} ==="

# Clean
rm -rf "$TMP"
mkdir -p "$TMP/pkg"

# 1. Build backend
echo "--- Building backend ---"
cd "$ROOT/backend"
go build -ldflags="-s -w" -o "$TMP/nas" .

# 2. Build frontend
echo "--- Building frontend ---"
cd "$ROOT/frontend"
npm run build

# 3. Assemble package
echo "--- Assembling package ---"

# Copy backend binary
mkdir -p "$TMP/pkg/usr/local/nas/bin"
cp "$TMP/nas" "$TMP/pkg/usr/local/nas/bin/"

# Copy frontend pages
mkdir -p "$TMP/pkg/usr/local/nas/pages"
cp -r "$ROOT/frontend/dist/"* "$TMP/pkg/usr/local/nas/pages/"

# Create data directories
mkdir -p "$TMP/pkg/data/data"

# Copy rootfs files (excluding test/dev configs)
rsync -a --exclude 'config.json' --exclude 'config_test.yaml' \
    "$ROOT/backend/rootfs/" "$TMP/pkg/"

# 4. Create archive
echo "--- Creating archive ---"
cd "$TMP/pkg"
tar czf "$OUTPUT" etc usr data
cd "$ROOT"

# 5. Cleanup
rm -rf "$TMP"

echo "=== Done: $OUTPUT ==="
echo ""
echo "Deploy:"
echo "  tar xzf ${PKG_NAME}-${VERSION}.tgz -C /"
echo "  systemctl daemon-reload && systemctl restart simple_nas"