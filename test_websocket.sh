#!/bin/bash

# DevSync WebSocket Test Script
# This script tests the WebSocket chat functionality

echo "🚀 DevSync WebSocket Chat Test"
echo "================================"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if websocat is installed
if ! command -v websocat &> /dev/null; then
    echo -e "${YELLOW}⚠️  websocat is not installed${NC}"
    echo "Installing websocat..."
    echo ""
    echo "On macOS, run:"
    echo "  brew install websocat"
    echo ""
    echo "Or download from: https://github.com/vi/websocat/releases"
    exit 1
fi

# Configuration
WS_URL="ws://localhost:8080/ws"
USER_ID=1
PROJECT_ID=1

echo -e "${BLUE}📡 Connecting to WebSocket...${NC}"
echo "URL: $WS_URL"
echo "User ID: $USER_ID"
echo "Project ID: $PROJECT_ID"
echo ""
echo -e "${GREEN}✅ Connected! Type messages below (Ctrl+C to exit)${NC}"
echo -e "${YELLOW}Format: type any message and press Enter${NC}"
echo ""

# Connect and test
websocat "$WS_URL" | while read -r line; do
    echo -e "${GREEN}📨 Received:${NC} $line"
done &

# Send test messages
echo '{"type":"chat_message","project_id":'$PROJECT_ID',"user_id":'$USER_ID',"data":{"message":"Hello from terminal!","user_id":'$USER_ID',"project_id":'$PROJECT_ID'}}' | websocat "$WS_URL"

wait
