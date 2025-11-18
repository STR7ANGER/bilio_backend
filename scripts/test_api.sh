#!/bin/bash

# BillStack API Test Script
# Make sure your server is running on http://localhost:8080
# Or update BASE_URL below

BASE_URL="http://localhost:8080/api/v1"

echo "=========================================="
echo "BillStack API Test Script"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ============================================
# 1. AUTHENTICATION
# ============================================
echo -e "${BLUE}=== 1. AUTHENTICATION ===${NC}"
echo ""

echo -e "${YELLOW}1.1 Register a new user${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User",
    "workspace_name": "Test Workspace"
  }')

echo "$REGISTER_RESPONSE" | jq '.'
echo ""

# Extract token from register response (if successful)
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.token // empty')

# If registration didn't return token, try login
if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo -e "${YELLOW}1.2 Login${NC}"
  LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d '{
      "email": "test@example.com",
      "password": "password123"
    }')
  
  echo "$LOGIN_RESPONSE" | jq '.'
  TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token // empty')
fi

if [ -z "$TOKEN" ] || [ "$TOKEN" == "null" ]; then
  echo -e "${RED}ERROR: Could not get authentication token. Please check your credentials.${NC}"
  exit 1
fi

echo ""
echo -e "${GREEN}✓ Authentication successful. Token: ${TOKEN:0:20}...${NC}"
echo ""

# ============================================
# 2. CLIENTS
# ============================================
echo -e "${BLUE}=== 2. CLIENTS ===${NC}"
echo ""

echo -e "${YELLOW}2.1 Create a client${NC}"
CLIENT_RESPONSE=$(curl -s -X POST "$BASE_URL/clients" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "Acme Corporation",
    "email": "contact@acme.com",
    "company": "Acme Corp",
    "phone": "+1-555-0123",
    "address": "123 Business St, City, State 12345",
    "tax_id": "TAX-123456",
    "currency": "USD"
  }')

echo "$CLIENT_RESPONSE" | jq '.'
CLIENT_ID=$(echo "$CLIENT_RESPONSE" | jq -r '.id // empty')
echo ""

echo -e "${YELLOW}2.2 List all clients${NC}"
curl -s -X GET "$BASE_URL/clients" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

if [ ! -z "$CLIENT_ID" ] && [ "$CLIENT_ID" != "null" ]; then
  echo -e "${YELLOW}2.3 Get client by ID${NC}"
  curl -s -X GET "$BASE_URL/clients/$CLIENT_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  echo -e "${YELLOW}2.4 Update client${NC}"
  curl -s -X PUT "$BASE_URL/clients/$CLIENT_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d '{
      "name": "Acme Corporation Updated",
      "email": "newcontact@acme.com",
      "company": "Acme Corp",
      "phone": "+1-555-0123",
      "address": "456 New St, City, State 12345",
      "tax_id": "TAX-123456",
      "currency": "USD"
    }' | jq '.'
  echo ""
fi

# ============================================
# 3. INVOICES
# ============================================
echo -e "${BLUE}=== 3. INVOICES ===${NC}"
echo ""

if [ ! -z "$CLIENT_ID" ] && [ "$CLIENT_ID" != "null" ]; then
  echo -e "${YELLOW}3.1 Create an invoice${NC}"
  INVOICE_RESPONSE=$(curl -s -X POST "$BASE_URL/invoices" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"client_id\": \"$CLIENT_ID\",
      \"invoice_number\": \"INV-001\",
      \"status\": \"draft\",
      \"issue_date\": \"2024-01-15T00:00:00Z\",
      \"due_date\": \"2024-02-15T00:00:00Z\",
      \"currency\": \"USD\",
      \"tax_rate\": 10.0,
      \"notes\": \"Payment terms: Net 30\",
      \"items\": [
        {
          \"description\": \"Web Development Services\",
          \"quantity\": 40,
          \"unit_price\": 100.00
        },
        {
          \"description\": \"Design Services\",
          \"quantity\": 20,
          \"unit_price\": 75.00
        }
      ]
    }")

  echo "$INVOICE_RESPONSE" | jq '.'
  INVOICE_ID=$(echo "$INVOICE_RESPONSE" | jq -r '.id // empty')
  echo ""

  echo -e "${YELLOW}3.2 List all invoices${NC}"
  curl -s -X GET "$BASE_URL/invoices" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  echo -e "${YELLOW}3.3 List invoices with filters (status=pending)${NC}"
  curl -s -X GET "$BASE_URL/invoices?status=pending" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  echo -e "${YELLOW}3.4 List invoices filtered by client${NC}"
  curl -s -X GET "$BASE_URL/invoices?client_id=$CLIENT_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  if [ ! -z "$INVOICE_ID" ] && [ "$INVOICE_ID" != "null" ]; then
    echo -e "${YELLOW}3.5 Get invoice by ID${NC}"
    curl -s -X GET "$BASE_URL/invoices/$INVOICE_ID" \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo ""

    echo -e "${YELLOW}3.6 Update invoice${NC}"
    curl -s -X PUT "$BASE_URL/invoices/$INVOICE_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"status\": \"pending\",
        \"issue_date\": \"2024-01-15T00:00:00Z\",
        \"due_date\": \"2024-02-15T00:00:00Z\",
        \"currency\": \"USD\",
        \"tax_rate\": 10.0,
        \"notes\": \"Updated payment terms\",
        \"items\": [
          {
            \"description\": \"Web Development Services\",
            \"quantity\": 50,
            \"unit_price\": 100.00
          },
          {
            \"description\": \"Design Services\",
            \"quantity\": 25,
            \"unit_price\": 75.00
          }
        ]
      }" | jq '.'
    echo ""

    echo -e "${YELLOW}3.7 Mark invoice as paid${NC}"
    curl -s -X POST "$BASE_URL/invoices/$INVOICE_ID/mark-paid" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"amount\": 6875.00,
        \"currency\": \"USD\",
        \"payment_method\": \"stripe\",
        \"payment_date\": \"2024-01-20T00:00:00Z\",
        \"transaction_id\": \"txn_123456789\",
        \"notes\": \"Payment received via Stripe\"
      }" | jq '.'
    echo ""

    echo -e "${YELLOW}3.8 Send invoice (placeholder)${NC}"
    curl -s -X POST "$BASE_URL/invoices/$INVOICE_ID/send" \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo ""

    echo -e "${YELLOW}3.9 Get invoice PDF (placeholder)${NC}"
    curl -s -X GET "$BASE_URL/invoices/$INVOICE_ID/pdf" \
      -H "Authorization: Bearer $TOKEN" | jq '.'
    echo ""
  fi
fi

# ============================================
# 4. EXPENSES
# ============================================
echo -e "${BLUE}=== 4. EXPENSES ===${NC}"
echo ""

echo -e "${YELLOW}4.1 Create an expense${NC}"
EXPENSE_RESPONSE=$(curl -s -X POST "$BASE_URL/expenses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"client_id\": \"$CLIENT_ID\",
    \"description\": \"Office Supplies\",
    \"amount\": 150.00,
    \"currency\": \"USD\",
    \"category\": \"office\",
    \"expense_date\": \"2024-01-10T00:00:00Z\",
    \"receipt_url\": \"https://example.com/receipts/receipt-001.pdf\",
    \"notes\": \"Purchased office supplies for project\"
  }")

echo "$EXPENSE_RESPONSE" | jq '.'
EXPENSE_ID=$(echo "$EXPENSE_RESPONSE" | jq -r '.id // empty')
echo ""

echo -e "${YELLOW}4.2 Create another expense${NC}"
curl -s -X POST "$BASE_URL/expenses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "description": "Travel Expenses",
    "amount": 500.00,
    "currency": "USD",
    "category": "travel",
    "expense_date": "2024-01-12T00:00:00Z",
    "notes": "Flight tickets for client meeting"
  }' | jq '.'
echo ""

echo -e "${YELLOW}4.3 List all expenses${NC}"
curl -s -X GET "$BASE_URL/expenses" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo -e "${YELLOW}4.4 List expenses filtered by client${NC}"
curl -s -X GET "$BASE_URL/expenses?client_id=$CLIENT_ID" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo -e "${YELLOW}4.5 List expenses filtered by category${NC}"
curl -s -X GET "$BASE_URL/expenses?category=travel" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

if [ ! -z "$EXPENSE_ID" ] && [ "$EXPENSE_ID" != "null" ]; then
  echo -e "${YELLOW}4.6 Get expense by ID${NC}"
  curl -s -X GET "$BASE_URL/expenses/$EXPENSE_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  echo -e "${YELLOW}4.7 Update expense${NC}"
  curl -s -X PUT "$BASE_URL/expenses/$EXPENSE_ID" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"client_id\": \"$CLIENT_ID\",
      \"description\": \"Office Supplies (Updated)\",
      \"amount\": 175.00,
      \"currency\": \"USD\",
      \"category\": \"office\",
      \"expense_date\": \"2024-01-10T00:00:00Z\",
      \"receipt_url\": \"https://example.com/receipts/receipt-001-updated.pdf\",
      \"notes\": \"Updated expense with additional items\"
    }" | jq '.'
  echo ""
fi

# ============================================
# 5. REPORTS
# ============================================
echo -e "${BLUE}=== 5. REPORTS ===${NC}"
echo ""

echo -e "${YELLOW}5.1 Get summary report${NC}"
curl -s -X GET "$BASE_URL/reports/summary" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo -e "${YELLOW}5.2 Get summary report with date range${NC}"
curl -s -X GET "$BASE_URL/reports/summary?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

if [ ! -z "$CLIENT_ID" ] && [ "$CLIENT_ID" != "null" ]; then
  echo -e "${YELLOW}5.3 Get client profitability report${NC}"
  curl -s -X GET "$BASE_URL/reports/client-profit/$CLIENT_ID" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""

  echo -e "${YELLOW}5.4 Get client profitability with date range${NC}"
  curl -s -X GET "$BASE_URL/reports/client-profit/$CLIENT_ID?from_date=2024-01-01&to_date=2024-01-31" \
    -H "Authorization: Bearer $TOKEN" | jq '.'
  echo ""
fi

echo -e "${YELLOW}5.5 Get tax summary${NC}"
curl -s -X GET "$BASE_URL/reports/tax-summary?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

echo -e "${GREEN}=========================================="
echo "All API tests completed!"
echo "==========================================${NC}"

