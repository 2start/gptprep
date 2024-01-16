#!/bin/sh

# Detect the OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

# Use GitHub API to get the latest release data
LATEST_RELEASE_DATA=$(curl -s https://api.github.com/repos/2start/gptprep/releases/latest)

# Extract the tag name (version) from the release data
LATEST_VERSION=$(echo "$LATEST_RELEASE_DATA" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Check if we got a version number
if [ -z "$LATEST_VERSION" ]; then
  echo "Could not find latest version. Exiting."
  exit 1
fi

# Define the base URL for downloading the binaries with the latest version
BASE_URL="https://github.com/2start/gptprep/releases/download/${LATEST_VERSION}"

# Choose the binary based on OS and architecture
case "${OS}_${ARCH}" in
  Linux_x86_64)   BINARY="gptprep-linux-amd64" ;;
  Linux_aarch64)  BINARY="gptprep-linux-arm64" ;;
  Darwin_x86_64)  BINARY="gptprep-darwin-amd64" ;;
  Darwin_arm64)   BINARY="gptprep-darwin-arm64" ;;
  # Add more cases for each supported OS and architecture
  *)              echo "Unsupported OS or architecture"; exit 1 ;;
esac

# Download the binary
curl -L -o gptprep "${BASE_URL}/${BINARY}"

# Make the binary executable
chmod +x gptprep

# Move the binary to a directory in PATH (e.g., /usr/local/bin)
# Ensure that /usr/local/bin exists and is writable, or choose another suitable directory
mv gptprep /usr/local/bin/gptprep

echo "gptprep installed successfully"
