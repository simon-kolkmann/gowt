#!/bin/bash

echo 'Cleaning build directory'
rm -r build

echo '[1/2] Building linux/amd64...'
GOARCH=amd64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/gowt-linux-amd64 .

echo '[2/2] Building windows/amd64...'
GOARCH=amd64 GOOS=windows go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o ./build/gowt-win-amd64 .

echo 'Done!'
