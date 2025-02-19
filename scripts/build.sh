#!/bin/bash

echo 'Cleaning build directory'
rm -r build

mkdir build
cd build

echo '[1/2] Building linux/amd64...'
GOARCH=amd64 GOOS=linux go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o gowt-linux-amd64 ..
tar -czvf gowt-linux-amd64.tar.gz gowt-linux-amd64

echo '[2/2] Building windows/amd64...'
GOARCH=amd64 GOOS=windows go build -ldflags="-X main.version=$(git describe --always --tags --dirty)" -o gowt-win-amd64 ..
tar -czvf gowt-win-amd64.tar.gz gowt-win-amd64

echo 'Done!'
