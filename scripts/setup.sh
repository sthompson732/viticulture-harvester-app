/*
 * setup.sh: Shell script for setting up the application environment.
 * Runs commands to configure the database and other services.
 * Usage: Run this script to prepare the application for first-time use.
 * Author(s): Shannon Thompson
 * Created on: 04/10/2024
 */
#!/bin/bash

# Ensuring that the Go is installed
if ! command -v go &> /dev/null
then
    echo "Go could not be found, please install it."
    exit
fi

echo "Go is installed, version: $(go version)"

# Navigate to the project directory; adjust this as necessary
cd "$(dirname "$0")"

# Setup GOPATH and GO111MODULE
export GOPATH=$(go env GOPATH)
export GO111MODULE=on
echo "Set GOPATH=$(go env GOPATH)"
echo "Enabled Go Modules"

# Ensure all dependencies are installed
echo "Installing dependencies..."
go mod tidy
if [ $? -ne 0 ]; then
    echo "Failed to install dependencies"
    exit 1
fi

echo "Dependencies installed successfully."

# Build the project
echo "Building the project..."
go build -o harvester ./cmd/harvester
if [ $? -ne 0 ]; then
    echo "Build failed"
    exit 1
fi

echo "Build successful. Binary created: harvester"

# Setting up runtime environment variables
echo "Setting up environment variables..."
export CONFIG_PATH="./config.yaml" # Adjust this path according to your config file placement
echo "Set CONFIG_PATH to ${CONFIG_PATH}"

# Instructions to run the application
cat << EOF

Setup completed successfully.

You can now run the application using:
./harvester

EOF
