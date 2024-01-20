#!/bin/bash

APP_NAME="drive-health"
DIST_DIR="dist"

# Create the dist directory if it doesn't exist
mkdir -p $DIST_DIR

# Build the application
echo "Building the application..."
GOOS=linux GOARCH=amd64 go build -o $DIST_DIR/$APP_NAME

# Copy additional resources (like .env, static files, templates) to the dist directory
echo "Copying additional resources..."
cp -r .env static templates $DIST_DIR/

echo "Compilation and packaging completed."
