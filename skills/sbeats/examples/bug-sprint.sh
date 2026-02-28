#!/bin/bash
# Example: Bug Sprint
# This example shows how to manage a bug-fixing sprint

set -e

echo "=== Bug Sprint Example ==="
echo ""

sb init

# Create sprint epic
echo "Creating sprint..."
sb create "Sprint 24: Bug Fixes" -d "Fix critical bugs reported this week" -p 0
# sb-sprint

echo ""
echo "=== Critical Bugs (P0) ==="
sb create "Fix login crash on mobile" -d "App crashes when users login on iOS Safari" -p 0 --parent sb-sprint

sb create "Database connection leak" -d "Connections not being closed properly" -p 0 --parent sb-sprint

echo ""
echo "=== High Priority (P1) ==="
sb create "Fix checkout calculation bug" -d "Tax calculation wrong for international orders" -p 1 --parent sb-sprint

sb create "Images not loading on slow connections" -d "Add proper loading states" -p 1 --parent sb-sprint

echo ""
echo "=== Medium Priority (P2) ==="
sb create "Update error messages" -d "Make error messages more user-friendly" -p 2 --parent sb-sprint

sb create "Fix pagination on search" -d "Page numbers don't update correctly" -p 2 --parent sb-sprint

echo ""
echo "=== Low Priority (P3-P4) ==="
sb create "Clean up console warnings" -d "Remove React warnings in dev tools" -p 3 --parent sb-sprint

sb create "Update README" -d "Add setup instructions for new devs" -p 4 --parent sb-sprint

echo ""
echo "=== Sprint Dashboard ==="
echo ""
echo "Critical bugs (work on these first!):"
sb list -p 0

echo ""
echo "All sprint tasks:"
sb list --parent sb-sprint

echo ""
echo "Workflow:"
echo "1. sb ready - See what's actionable"
echo "2. sb list -p 0 - Focus on critical bugs"
echo "3. sb close <id> - Mark fixed bugs as closed"
