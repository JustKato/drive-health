#!/bin/bash

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

# Run tests before proceeding
echo_color yellow "Running tests..."
if ! go test; then
    echo_color red "Tests failed. Cancelling build process."
    exit 1
fi
echo_color green "All tests passed successfully."

echo_color green "Starting the Docker build process with version $GIT_VERSION..."

LATEST_IMAGE_NAME="ghcr.io/justkato/drive-health:latest"
IMAGE_NAME="ghcr.io/justkato/drive-health:$GIT_VERSION"
echo_color yellow "Image to be built: $IMAGE_NAME"

# Confirmation to build
read -p "Are you sure you want to build an image? (y/N) " response
if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    # Building the Docker image
    echo "Building Docker image: $IMAGE_NAME"
    docker build --no-cache -t $IMAGE_NAME .

    # Also tag this build as 'latest'
    echo "Tagging image as latest: $LATEST_IMAGE_NAME"
    docker tag $IMAGE_NAME $LATEST_IMAGE_NAME
else
    echo_color red "Build cancelled."
    exit 1
fi

# Prompt to push the image
read -p "Push image to repository? (y/N) " push_response
if [[ "$push_response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    # Pushing the image
    echo "Pushing image: $IMAGE_NAME"
    docker push $IMAGE_NAME

    # Pushing the 'latest' image
    echo "Pushing latest image: $LATEST_IMAGE_NAME"
    docker push $LATEST_IMAGE_NAME
else
    echo_color red "Push cancelled."
fi

echo_color green "Ending the Docker build process..."