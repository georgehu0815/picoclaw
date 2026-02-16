#!/bin/bash
# Run PicoClaw in gateway mode (for Telegram, Discord, etc.)

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Check if binary exists
if [ ! -f "./picoclaw" ]; then
    echo -e "${YELLOW}Binary not found. Building...${NC}"
    ./build-and-run.sh gateway "$@"
    exit $?
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🦞 PicoClaw Gateway Mode${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Starting gateway for chat channels...${NC}"
echo -e "${YELLOW}Enabled channels will be listed below:${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"

./picoclaw gateway "$@"
