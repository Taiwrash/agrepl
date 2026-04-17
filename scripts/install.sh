#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing agrepl...${NC}"

# Check for Go
if ! command -v go &> /dev/null
then
    echo -e "${RED}Error: Go is not installed. Please install Go to build agrepl.${NC}"
    exit 1
fi

# Determine OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

echo -e "System: ${OS}/${ARCH}"

# Build
echo -e "Building agrepl..."
go build -o agrepl main.go

# Install
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    echo -e "${BLUE}Requesting sudo to install to $INSTALL_DIR${NC}"
    sudo mv agrepl "$INSTALL_DIR/agrepl"
else
    mv agrepl "$INSTALL_DIR/agrepl"
fi

echo -e "${GREEN}✓ agrepl installed successfully to $INSTALL_DIR/agrepl${NC}"
echo -e ""
echo -e "Try it out:"
echo -e "  agrepl --help"
echo -e ""
echo -e "${BLUE}Note: To intercept HTTPS, agrepl generates a local Root CA.${NC}"
echo -e "The first time you run 'record', it will be created in .agent-replay/ca/"
