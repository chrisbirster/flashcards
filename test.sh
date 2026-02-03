#!/bin/bash

set -e

echo "=================================="
echo "Running Microdote Test Suite"
echo "=================================="
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Backend tests
echo -e "${BLUE}Running Backend Tests (Go)...${NC}"
go test -v -count=1
BACKEND_EXIT=$?

if [ $BACKEND_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ Backend tests passed${NC}"
else
    echo -e "${RED}✗ Backend tests failed${NC}"
    exit 1
fi

echo ""

# Check if backend server is running
echo -e "${BLUE}Checking backend server...${NC}"
if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend server is running${NC}"
else
    echo -e "${RED}✗ Backend server not running. Please start it with: go run *.go${NC}"
    exit 1
fi

echo ""

# Frontend E2E tests
echo -e "${BLUE}Running Frontend E2E Tests (Playwright)...${NC}"
cd web
npm run test:e2e
FRONTEND_EXIT=$?
cd ..

if [ $FRONTEND_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ Frontend E2E tests passed${NC}"
else
    echo -e "${RED}✗ Frontend E2E tests failed${NC}"
    exit 1
fi

echo ""
echo "=================================="
echo -e "${GREEN}✓ All tests passed!${NC}"
echo "=================================="
echo ""
echo "Test Summary:"
echo "  - Backend: 15 unit tests (8 M0 + 4 M1 + 3 M2)"
echo "  - Frontend: 37 E2E tests (10 deck management + 14 study screen + 13 add note)"
echo "  - Total: 52 tests"
