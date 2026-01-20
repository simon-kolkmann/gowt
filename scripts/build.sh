#!/bin/bash

echo 'Cleaning build directory'
rm -rf build

mkdir build
cd build

VERSION=$(git describe --always --tags --dirty)
BINARY_BASE_NAME=gowt-$VERSION

# build targets
declare -a TARGETS=(
	"linux/amd64"
	"windows/amd64"
)

for TARGET in "${TARGETS[@]}"; do
	# split the target into GOOS and GOARCH
	IFS='/' read -r GOOS GOARCH <<< "$TARGET"

	BINARY_NAME=$BINARY_BASE_NAME-$GOOS-$GOARCH

	echo "Building $TARGET..."

	# compile
	GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-X main.version=$VERSION" -o "$BINARY_NAME" ..

	# make executable
	chmod +x "$BINARY_NAME"

	# create archive
	tar -czvf "$BINARY_NAME.tar.gz" "$BINARY_NAME"
done

echo 'Done!'
