#!/usr/bin/env bash

set -o pipefail
set -u

# Function to display messages in color
echo_color() {
    color=$1
    text=$2
    case $color in
        "green") echo -e "\033[0;32m$text\033[0m" ;;
        "yellow") echo -e "\033[0;33m$text\033[0m" ;;
        "red") echo -e "\033[0;31m$text\033[0m" ;;
        *) echo "$text" ;;
    esac
}

# Getting GIT_VERSION from the most recent tag or commit hash
GIT_VERSION=$(git describe --tags --always)
if [ -z "$GIT_VERSION" ]; then
    echo_color red "Error: Unable to determine GIT_VERSION."
    exit 1
fi

APP_NAME="drive-health"
DIST_DIR="${DIST_DIR:-dist}"
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# make sure we are in the source dir
cd $SCRIPT_DIR;

# Create the dist directory if it doesn't exist
mkdir -p $DIST_DIR

# Build the application
echo_color yellow "[🦝] Building the application..."
GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o $DIST_DIR/$APP_NAME

# Copying additional resources...
cp -r static templates $DIST_DIR/

echo_color yellow "[🦝] Compilation and packaging completed, archiving..."

cd $DIST_DIR/

zip "drive-health_$GIT_VERSION.zip" -r .

# TODO: Add reliable method of cleaning up the compiled files optionally

cd $SCRIPT_DIR;