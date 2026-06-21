#!/bin/bash

# NAS Dashboard Start Script

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Starting NAS Dashboard System...${NC}"

# Start Backend
echo -e "${GREEN}Starting Backend...${NC}"
cd backend
if [ -f "main" ]; then
    ./main &
    BACKEND_PID=$!
else
    # Try go run if binary not found
    go run cmd/server/main.go &
    BACKEND_PID=$!
fi

# Start Frontend
echo -e "${GREEN}Starting Frontend...${NC}"
cd ../frontend
npm run dev -- --host &
FRONTEND_PID=$!

echo -e "${BLUE}NAS Dashboard is starting!${NC}"
echo -e "Backend PID: $BACKEND_PID"
echo -e "Frontend PID: $FRONTEND_PID"
echo -e "Press Ctrl+C to stop all."

# Wait for both processes
trap "kill $BACKEND_PID $FRONTEND_PID; exit" INT
wait
