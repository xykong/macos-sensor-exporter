#!/bin/bash
# 
# macOS Sensor Exporter Installation Script
# This script downloads and installs the precompiled binary without requiring build tools
#

set -e

VERSION="1.1.0"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="macos-sensor-exporter"

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "arm64" ]; then
    DOWNLOAD_URL="https://github.com/xykong/macos-sensor-exporter/releases/download/v${VERSION}/macos-sensor-exporter_${VERSION}_Darwin_arm64.tar.gz"
    SHA256="1160c195442373c60891ea342b02b3d363a33ed90016a151e8638942edf88c12"
elif [ "$ARCH" = "x86_64" ]; then
    DOWNLOAD_URL="https://github.com/xykong/macos-sensor-exporter/releases/download/v${VERSION}/macos-sensor-exporter_${VERSION}_Darwin_x86_64.tar.gz"
    SHA256="8f4a4d09fb84c2c11c8868eb3911982e3d6642d16897fdfc3803ff38038b7103"
else
    echo "Error: Unsupported architecture: $ARCH"
    exit 1
fi

echo "macOS Sensor Exporter v${VERSION} Installation"
echo "=============================================="
echo "Architecture: $ARCH"
echo "Download URL: $DOWNLOAD_URL"
echo ""

# Create temp directory
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Download
echo "Downloading..."
cd "$TEMP_DIR"
if command -v curl &> /dev/null; then
    curl -sL "$DOWNLOAD_URL" -o archive.tar.gz
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O archive.tar.gz
else
    echo "Error: Neither curl nor wget is available"
    exit 1
fi

# Verify checksum
echo "Verifying checksum..."
if command -v shasum &> /dev/null; then
    echo "$SHA256  archive.tar.gz" | shasum -a 256 -c - || {
        echo "Error: Checksum verification failed"
        exit 1
    }
else
    echo "Warning: shasum not available, skipping checksum verification"
fi

# Extract
echo "Extracting..."
tar -xzf archive.tar.gz

# Install
echo "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Need sudo permission to install to $INSTALL_DIR"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Verify installation
if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
    echo ""
    echo "✓ Installation successful!"
    echo ""
    echo "Usage:"
    echo "  $BINARY_NAME start          # Start the exporter (default port: 9101)"
    echo "  $BINARY_NAME show           # Show sensor information"
    echo "  $BINARY_NAME show -o json   # Show sensor information in JSON format"
    echo "  $BINARY_NAME --help         # Show help"
    echo ""
    echo "To uninstall:"
    echo "  sudo rm $INSTALL_DIR/$BINARY_NAME"
else
    echo "Error: Installation failed"
    exit 1
fi
