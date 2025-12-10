#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔═══════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║         CleanGo Installation Script v1.0              ║${NC}"
echo -e "${BLUE}╚═══════════════════════════════════════════════════════╝${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.20 or higher.${NC}"
    echo -e "${YELLOW}Visit: https://golang.org/dl/${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo -e "${GREEN}✓${NC} Go version: ${GO_VERSION}"

# Get GOPATH
GOPATH=$(go env GOPATH)
GOBIN="${GOPATH}/bin"

echo -e "${BLUE}📦 Installing cleango...${NC}"
echo ""

# Install using go install
go install github.com/YeridStick/cleango/cmd/cleango@latest

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✅ CleanGo installed successfully!${NC}"
    echo ""

    # Check if GOBIN is in PATH
    if [[ ":$PATH:" != *":${GOBIN}:"* ]]; then
        echo -e "${YELLOW}⚠️  Warning: ${GOBIN} is not in your PATH${NC}"
        echo ""
        echo -e "Add this line to your ${BLUE}~/.bashrc${NC} or ${BLUE}~/.zshrc${NC}:"
        echo -e "${GREEN}export PATH=\$PATH:${GOBIN}${NC}"
        echo ""
        echo "Then run: source ~/.bashrc  (or source ~/.zshrc)"
        echo ""
    fi

    # Verify installation
    if command -v cleango &> /dev/null; then
        CLEANGO_VERSION=$(cleango --version 2>&1 || echo "unknown")
        echo -e "${GREEN}✓${NC} cleango is now available"
        echo -e "${GREEN}✓${NC} Version: ${CLEANGO_VERSION}"
    else
        echo -e "${YELLOW}⚠️  cleango command not found in PATH. Please add ${GOBIN} to PATH.${NC}"
    fi

    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}🚀 Quick Start:${NC}"
    echo ""
    echo -e "  ${YELLOW}1.${NC} Create a new project:"
    echo -e "     ${GREEN}cleango new my-api${NC}"
    echo ""
    echo -e "  ${YELLOW}2.${NC} Navigate to your project:"
    echo -e "     ${GREEN}cd my-api${NC}"
    echo ""
    echo -e "  ${YELLOW}3.${NC} Run the application:"
    echo -e "     ${GREEN}go run ./cmd/api${NC}"
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "For more information, visit: ${BLUE}https://github.com/YeridStick/cleango${NC}"
    echo ""
else
    echo ""
    echo -e "${RED}❌ Installation failed. Please check the error messages above.${NC}"
    exit 1
fi
