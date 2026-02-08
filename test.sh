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
if curl -s http://localhost:8000/api/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓ Backend server is running${NC}"
else
    echo -e "${RED}✗ Backend server not running. Please start it with: go run .${NC}"
    exit 1
fi

echo ""

# Coverage gate (baseline)
echo -e "${BLUE}Running Coverage Gate (baseline)...${NC}"
./scripts/check_coverage.sh baseline
COVERAGE_EXIT=$?

if [ $COVERAGE_EXIT -eq 0 ]; then
    echo -e "${GREEN}✓ Coverage baseline passed${NC}"
else
    echo -e "${RED}✗ Coverage baseline failed${NC}"
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
echo "  - Backend: go test ./... (see output for current test count)"
echo "  - Coverage: baseline gate via scripts/check_coverage.sh"
echo "  - Frontend: Playwright E2E suite"
