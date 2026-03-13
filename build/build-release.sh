#!/bin/bash

# Exit immediately if a command exits with a non-zero status.
set -e

VERSION=$(git describe --tags --always --dirty --long)
MAIN_PACKAGE="../cmd/main.go"
OUTPUT_NAME="go-mirror-zig"

mkdir -p release

# Target platforms.
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "windows/amd64"
    "darwin/amd64"
    "darwin/arm64"
)

echo "Building release artifacts for version ${VERSION}..."

for platform in "${PLATFORMS[@]}"; do
    GOOS=$(dirname "${platform}")
    GOARCH=$(basename "${platform}")

    OUTPUT_FILENAME="${OUTPUT_NAME}-${GOOS}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        OUTPUT_FILENAME="${OUTPUT_FILENAME}.exe"
    fi

    echo "Building for ${GOOS}/${GOARCH}..."

    # Build the binary with version information.
    # The -s and -w flags strip debugging information, reducing the binary size.
    env GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-s -w -X main.version=${VERSION}" -o "release/${OUTPUT_FILENAME}" ${MAIN_PACKAGE}

    cd release
    tar -czvf "${OUTPUT_NAME}-${VERSION}-${GOOS}-${GOARCH}.tar.gz" "${OUTPUT_FILENAME}"
    rm "${OUTPUT_FILENAME}"
    cd ..
done

echo "Build complete. Artifacts are in the 'release' directory."