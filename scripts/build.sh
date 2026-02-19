#!/bin/bash

# build.sh - Build ImpactBench for multiple platforms

APP_NAME="impactbench"
SRC_PATH="./cmd/impactbench"
DIST_DIR="./dist"

# Create dist directory
mkdir -p $DIST_DIR

echo "Building $APP_NAME..."

# Linux
echo "  - Linux amd64"
GOOS=linux GOARCH=amd64 /usr/local/go/bin/go build -o $DIST_DIR/$APP_NAME-linux-amd64 $SRC_PATH

# macOS
echo "  - macOS amd64"
GOOS=darwin GOARCH=amd64 /usr/local/go/bin/go build -o $DIST_DIR/$APP_NAME-darwin-amd64 $SRC_PATH
echo "  - macOS arm64"
GOOS=darwin GOARCH=arm64 /usr/local/go/bin/go build -o $DIST_DIR/$APP_NAME-darwin-arm64 $SRC_PATH

# Windows
echo "  - Windows amd64"
GOOS=windows GOARCH=amd64 /usr/local/go/bin/go build -o $DIST_DIR/$APP_NAME-windows-amd64.exe $SRC_PATH

echo "Builds complete! Files located in $DIST_DIR"
