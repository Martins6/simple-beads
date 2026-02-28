#!/bin/bash
# Example: Web Development Project
# This example shows a typical web development workflow with sbeads

set -e

echo "=== Web Development Project Example ==="
echo ""

# Initialize project
sb init

# Create main epic
echo "Creating project epic..."
sb create "Build E-commerce Website" -d "Full-stack e-commerce platform with React and Node.js" -p 0
# Assume we get sb-epic1

echo ""
echo "=== Phase 1: Design & Planning ==="
sb create "Create wireframes in Figma" -d "Homepage, product page, cart, checkout" -p 0 --parent sb-epic1
# sb-design1

sb create "Design database schema" -d "Users, products, orders, cart tables" -p 0 --parent sb-epic1
# sb-db1

sb create "Setup project repository" -d "Initialize Git, setup folder structure" -p 1 --parent sb-epic1
# sb-repo

echo ""
echo "=== Phase 2: Backend Development ==="
sb create "Setup PostgreSQL database" -d "Create tables and initial schema" -p 1 --parent sb-epic1 --deps sb-db1
# sb-setup-db

sb create "Implement user authentication" -d "JWT-based auth, login/signup endpoints" -p 1 --parent sb-epic1 --deps sb-setup-db
# sb-auth

sb create "Build products API" -d "CRUD endpoints for products" -p 1 --parent sb-epic1 --deps sb-setup-db
# sb-products

sb create "Build cart API" -d "Add/remove items, calculate totals" -p 1 --parent sb-epic1 --deps sb-setup-db,sb-auth
# sb-cart

sb create "Build checkout API" -d "Process orders, integrate payment" -p 1 --parent sb-epic1 --deps sb-cart
# sb-checkout

echo ""
echo "=== Phase 3: Frontend Development ==="
sb create "Setup React application" -d "Initialize with Vite, configure routing" -p 2 --parent sb-epic1 --deps sb-repo
# sb-react

sb create "Build authentication UI" -d "Login, signup, password reset forms" -p 2 --parent sb-epic1 --deps sb-react,sb-auth
# sb-auth-ui

sb create "Build product catalog" -d "Product grid, filters, search" -p 2 --parent sb-epic1 --deps sb-react,sb-products
# sb-catalog

sb create "Build shopping cart" -d "Cart page, add/remove items" -p 2 --parent sb-epic1 --deps sb-react,sb-cart
# sb-cart-ui

sb create "Build checkout flow" -d "Checkout form, payment integration" -p 2 --parent sb-epic1 --deps sb-react,sb-checkout,sb-cart-ui
# sb-checkout-ui

echo ""
echo "=== Phase 4: Testing & Deployment ==="
sb create "Write API tests" -d "Unit and integration tests for backend" -p 2 --parent sb-epic1 --deps sb-checkout
# sb-api-tests

sb create "Write E2E tests" -d "Cypress tests for critical paths" -p 2 --parent sb-epic1 --deps sb-checkout-ui
# sb-e2e-tests

sb create "Performance optimization" -d "Optimize queries, add caching" -p 3 --parent sb-epic1 --deps sb-api-tests
# sb-perf

sb create "Deploy to production" -d "Setup CI/CD, deploy to AWS" -p 0 --parent sb-epic1 --deps sb-e2e-tests,sb-checkout-ui
# sb-deploy

echo ""
echo "=== Project Status ==="
echo ""
echo "All tasks created! Here's your workflow:"
echo ""
echo "Ready to start (no blockers):"
sb ready | head -10

echo ""
echo "Blocked (waiting on dependencies):"
sb blocked

echo ""
echo "Start with: sb ready"
echo "Then: sb close <task-id> as you complete work"
