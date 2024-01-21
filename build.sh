#!/bin/sh

set -o pipefail
set -u

APP_NAME="drive-health"
DIST_DIR="${DIST_DIR:-dist}"

# Create the dist directory if it doesn't exist
mkdir -p $DIST_DIR

# Build the application
echo "Building the application..."
GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o $DIST_DIR/$APP_NAME

# echo "Copying additional resources..."
cp -r static templates $DIST_DIR/

echo "Compilation and packaging completed."
