#!/bin/bash
# Quick test without rebuild (if binary already exists)

if [ ! -f "./picoclaw" ]; then
    echo "❌ Binary not found. Run ./build-and-run.sh first"
    exit 1
fi

MESSAGE="${1:-What is 2+2? Reply with just the number.}"

echo "🦞 Testing PicoClaw with Claude..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
./picoclaw agent -m "$MESSAGE"
