#!/bin/bash
# PicoClaw Build and Run Script
# Usage: ./build-and-run.sh [command] [args...]

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Project directory
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_DIR"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🦞 PicoClaw Build and Run${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Step 1: Copy workspace directory for go:embed
echo -e "${YELLOW}📦 Preparing workspace...${NC}"
if [ ! -d "cmd/picoclaw/workspace" ]; then
    cp -r workspace cmd/picoclaw/workspace
    echo -e "${GREEN}✓ Workspace copied${NC}"
else
    echo -e "${GREEN}✓ Workspace already exists${NC}"
fi

# Step 2: Build
echo -e "${YELLOW}🔨 Building PicoClaw...${NC}"
go build -o picoclaw cmd/picoclaw/main.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful${NC}"

    # Show binary info
    SIZE=$(ls -lh picoclaw | awk '{print $5}')
    echo -e "${BLUE}   Binary size: ${SIZE}${NC}"
else
    echo -e "${RED}✗ Build failed${NC}"
    exit 1
fi

# Step 3: Run
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

# Default command is "agent" with a test message
COMMAND="${1:-agent}"
shift || true  # Remove first argument if exists

if [ "$COMMAND" = "agent" ] && [ $# -eq 0 ]; then
    # Default test
    echo -e "${YELLOW}🚀 Running test: What is 2+2?${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    ./picoclaw agent -m "What is 2+2? Reply with just the number."
else
    # Run with user-provided command and args
    echo -e "${YELLOW}🚀 Running: picoclaw $COMMAND $@${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    ./picoclaw "$COMMAND" "$@"
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✓ Done!${NC}"
