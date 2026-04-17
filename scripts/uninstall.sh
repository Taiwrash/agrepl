#!/bin/bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Uninstalling agrepl...${NC}"

BINARY_PATH="/usr/local/bin/agrepl"

# 1. Remove binary
if [ -f "$BINARY_PATH" ]; then
    echo -e "Removing binary from $BINARY_PATH..."
    if [ ! -w "$BINARY_PATH" ]; then
        echo -e "${BLUE}Requesting sudo to remove $BINARY_PATH${NC}"
        sudo rm "$BINARY_PATH"
    else
        rm "$BINARY_PATH"
    fi
else
    echo -e "${RED}agrepl binary not found at $BINARY_PATH${NC}"
fi

# 2. Optional: Clear global config
read -p "Do you also want to remove all global configuration and auth data (~/.agrepl)? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "Removing ~/.agrepl..."
    rm -rf "$HOME/.agrepl"
fi

echo -e "${GREEN}✓ agrepl has been uninstalled.${NC}"
