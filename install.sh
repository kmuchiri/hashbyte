#!/bin/bash
set -e

# Detect OS and architecture
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

if [ "$ARCH" = "x86_64" ]; then
    ARCH="x86_64"
elif [ "$ARCH" = "aarch64" ] || [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
elif [ "$ARCH" = "i386" ]; then
    ARCH="i386"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

REPO="kmuchiri/hashbyte"
echo "Fetching latest release for $REPO..."

# Get latest release tag from GitHub API
LATEST_TAG=$(curl -sL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_TAG" ]; then
    echo "Error: Could not determine latest release tag. Have you published a release yet?"
    exit 1
fi

# Map to goreleaser default naming format (e.g., hashbyte_Linux_x86_64.tar.gz)
if [ "$OS" = "linux" ]; then
    OS_TITLE="Linux"
elif [ "$OS" = "darwin" ]; then
    OS_TITLE="Darwin"
else
    echo "Unsupported OS: $OS"
    exit 1
fi

FILENAME="hashbyte_${OS_TITLE}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$REPO/releases/download/${LATEST_TAG}/${FILENAME}"

echo "Downloading $DOWNLOAD_URL..."
curl -sL "$DOWNLOAD_URL" -o "/tmp/$FILENAME"

echo "Extracting..."
tar -xzf "/tmp/$FILENAME" -C /tmp hashbyte

echo "Installing to /usr/local/bin/hashbyte..."
sudo mv /tmp/hashbyte /usr/local/bin/hashbyte
sudo chmod +x /usr/local/bin/hashbyte

rm "/tmp/$FILENAME"

echo "hashbyte successfully installed!"
echo "Run 'hashbyte' to get started."
