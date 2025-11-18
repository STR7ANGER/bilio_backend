# API Testing Guide

This document contains all curl commands to test the BillStack API endpoints.

## Prerequisites

1. Start the server: `go run cmd/server/main.go`
2. Run migrations: `migrate -path migrations -database $DATABASE_URL up`
3. Set environment variable: `export JWT_SECRET=your-secret-key`

## Base URL

```
http://localhost:8080/api/v1
```

## 1. Authentication

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User",
    "workspace_name": "Test Workspace"
  }'
```

**Response (201 Created):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com",
    "name": "Test User",
    "workspace_name": "Test Workspace",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzA1MzI1ODAwLCJpYXQiOjE3MDQ3MjEwMDB9.signature"
}
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Response (200 OK):**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com",
    "name": "Test User",
    "workspace_name": "Test Workspace",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNTUwZTg0MDAtZTI5Yi00MWQ0LWE3MTYtNDQ2NjU1NDQwMDAwIiwiZW1haWwiOiJ0ZXN0QGV4YW1wbGUuY29tIiwiZXhwIjoxNzA1MzI1ODAwLCJpYXQiOjE3MDQ3MjEwMDB9.signature"
}
```

**Error Response (400/401):**
```json
{
  "error": "invalid credentials"
}
```

**Save the token from the response for authenticated requests.**

---

## 2. Clients

Replace `YOUR_TOKEN` with the JWT token from login/register.

### List All Clients
```bash
curl -X GET http://localhost:8080/api/v1/clients \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
[
  {
    "id": "660e8400-e29b-41d4-a716-446655440001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Acme Corporation",
    "email": "contact@acme.com",
    "company": "Acme Corp",
    "phone": "+1-555-0123",
    "address": "123 Business St, City, State 12345",
    "tax_id": "TAX-123456",
    "currency": "USD",
    "created_at": "2024-01-15T10:35:00Z",
    "updated_at": "2024-01-15T10:35:00Z"
  }
]
```

### Create Client
```bash
curl -X POST http://localhost:8080/api/v1/clients \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Acme Corporation",
    "email": "contact@acme.com",
    "company": "Acme Corp",
    "phone": "+1-555-0123",
    "address": "123 Business St, City, State 12345",
    "tax_id": "TAX-123456",
    "currency": "USD"
  }'
```

**Response (201 Created):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "email": "contact@acme.com",
  "company": "Acme Corp",
  "phone": "+1-555-0123",
  "address": "123 Business St, City, State 12345",
  "tax_id": "TAX-123456",
  "currency": "USD",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

### Get Client by ID
```bash
curl -X GET http://localhost:8080/api/v1/clients/CLIENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation",
  "email": "contact@acme.com",
  "company": "Acme Corp",
  "phone": "+1-555-0123",
  "address": "123 Business St, City, State 12345",
  "tax_id": "TAX-123456",
  "currency": "USD",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "client not found"
}
```

### Update Client
```bash
curl -X PUT http://localhost:8080/api/v1/clients/CLIENT_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Acme Corporation Updated",
    "email": "newcontact@acme.com",
    "company": "Acme Corp",
    "phone": "+1-555-0123",
    "address": "456 New St, City, State 12345",
    "tax_id": "TAX-123456",
    "currency": "USD"
  }'
```

**Response (200 OK):**
```json
{
  "id": "660e8400-e29b-41d4-a716-446655440001",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Acme Corporation Updated",
  "email": "newcontact@acme.com",
  "company": "Acme Corp",
  "phone": "+1-555-0123",
  "address": "456 New St, City, State 12345",
  "tax_id": "TAX-123456",
  "currency": "USD",
  "created_at": "2024-01-15T10:35:00Z",
  "updated_at": "2024-01-15T10:40:00Z"
}
```

### Delete Client
```bash
curl -X DELETE http://localhost:8080/api/v1/clients/CLIENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (204 No Content):**
```
(No response body)
```

---

## 3. Invoices

### List All Invoices
```bash
curl -X GET http://localhost:8080/api/v1/invoices \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
[
  {
    "id": "770e8400-e29b-41d4-a716-446655440002",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "client_id": "660e8400-e29b-41d4-a716-446655440001",
    "invoice_number": "INV-001",
    "status": "draft",
    "issue_date": "2024-01-15T00:00:00Z",
    "due_date": "2024-02-15T00:00:00Z",
    "currency": "USD",
    "subtotal": 5500.00,
    "tax_rate": 10.0,
    "tax_amount": 550.00,
    "total": 6050.00,
    "notes": "Payment terms: Net 30",
    "payment_link": null,
    "created_at": "2024-01-15T10:45:00Z",
    "updated_at": "2024-01-15T10:45:00Z"
  }
]
```

### List Invoices with Filters
```bash
# Filter by status
curl -X GET "http://localhost:8080/api/v1/invoices?status=pending" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by client
curl -X GET "http://localhost:8080/api/v1/invoices?client_id=CLIENT_ID" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/invoices?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Combine filters
curl -X GET "http://localhost:8080/api/v1/invoices?status=pending&client_id=CLIENT_ID&from_date=2024-01-01" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):** Same format as List All Invoices, but filtered.

### Create Invoice
```bash
curl -X POST http://localhost:8080/api/v1/invoices \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "client_id": "CLIENT_ID",
    "invoice_number": "INV-001",
    "status": "draft",
    "issue_date": "2024-01-15T00:00:00Z",
    "due_date": "2024-02-15T00:00:00Z",
    "currency": "USD",
    "tax_rate": 10.0,
    "notes": "Payment terms: Net 30",
    "items": [
      {
        "description": "Web Development Services",
        "quantity": 40,
        "unit_price": 100.00
      },
      {
        "description": "Design Services",
        "quantity": 20,
        "unit_price": 75.00
      }
    ]
  }'
```

**Response (201 Created):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "invoice_number": "INV-001",
  "status": "draft",
  "issue_date": "2024-01-15T00:00:00Z",
  "due_date": "2024-02-15T00:00:00Z",
  "currency": "USD",
  "subtotal": 5500.00,
  "tax_rate": 10.0,
  "tax_amount": 550.00,
  "total": 6050.00,
  "notes": "Payment terms: Net 30",
  "payment_link": null,
  "created_at": "2024-01-15T10:45:00Z",
  "updated_at": "2024-01-15T10:45:00Z",
  "items": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "invoice_id": "770e8400-e29b-41d4-a716-446655440002",
      "description": "Web Development Services",
      "quantity": 40,
      "unit_price": 100.00,
      "amount": 4000.00,
      "created_at": "2024-01-15T10:45:00Z",
      "updated_at": "2024-01-15T10:45:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440004",
      "invoice_id": "770e8400-e29b-41d4-a716-446655440002",
      "description": "Design Services",
      "quantity": 20,
      "unit_price": 75.00,
      "amount": 1500.00,
      "created_at": "2024-01-15T10:45:00Z",
      "updated_at": "2024-01-15T10:45:00Z"
    }
  ]
}
```

### Get Invoice by ID
```bash
curl -X GET http://localhost:8080/api/v1/invoices/INVOICE_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "invoice_number": "INV-001",
  "status": "pending",
  "issue_date": "2024-01-15T00:00:00Z",
  "due_date": "2024-02-15T00:00:00Z",
  "currency": "USD",
  "subtotal": 6875.00,
  "tax_rate": 10.0,
  "tax_amount": 687.50,
  "total": 7562.50,
  "notes": "Updated payment terms",
  "payment_link": null,
  "created_at": "2024-01-15T10:45:00Z",
  "updated_at": "2024-01-15T10:50:00Z",
  "items": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440003",
      "invoice_id": "770e8400-e29b-41d4-a716-446655440002",
      "description": "Web Development Services",
      "quantity": 50,
      "unit_price": 100.00,
      "amount": 5000.00,
      "created_at": "2024-01-15T10:45:00Z",
      "updated_at": "2024-01-15T10:50:00Z"
    },
    {
      "id": "880e8400-e29b-41d4-a716-446655440004",
      "invoice_id": "770e8400-e29b-41d4-a716-446655440002",
      "description": "Design Services",
      "quantity": 25,
      "unit_price": 75.00,
      "amount": 1875.00,
      "created_at": "2024-01-15T10:45:00Z",
      "updated_at": "2024-01-15T10:50:00Z"
    }
  ],
  "payments": []
}
```

### Update Invoice
```bash
curl -X PUT http://localhost:8080/api/v1/invoices/INVOICE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "status": "pending",
    "issue_date": "2024-01-15T00:00:00Z",
    "due_date": "2024-02-15T00:00:00Z",
    "currency": "USD",
    "tax_rate": 10.0,
    "notes": "Updated payment terms",
    "items": [
      {
        "description": "Web Development Services",
        "quantity": 50,
        "unit_price": 100.00
      },
      {
        "description": "Design Services",
        "quantity": 25,
        "unit_price": 75.00
      }
    ]
  }'
```

**Response (200 OK):** Same format as Get Invoice by ID, with updated values.

**Error Response (400):**
```json
{
  "error": "can only update draft or pending invoices"
}
```

### Mark Invoice as Paid
```bash
curl -X POST http://localhost:8080/api/v1/invoices/INVOICE_ID/mark-paid \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "amount": 6875.00,
    "currency": "USD",
    "payment_method": "stripe",
    "payment_date": "2024-01-20T00:00:00Z",
    "transaction_id": "txn_123456789",
    "notes": "Payment received via Stripe"
  }'
```

**Response (200 OK):**
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440002",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "invoice_number": "INV-001",
  "status": "paid",
  "issue_date": "2024-01-15T00:00:00Z",
  "due_date": "2024-02-15T00:00:00Z",
  "currency": "USD",
  "subtotal": 6875.00,
  "tax_rate": 10.0,
  "tax_amount": 687.50,
  "total": 7562.50,
  "notes": "Updated payment terms",
  "payment_link": null,
  "created_at": "2024-01-15T10:45:00Z",
  "updated_at": "2024-01-15T10:55:00Z",
  "items": [...],
  "payments": [
    {
      "id": "990e8400-e29b-41d4-a716-446655440005",
      "invoice_id": "770e8400-e29b-41d4-a716-446655440002",
      "amount": 6875.00,
      "currency": "USD",
      "payment_method": "stripe",
      "payment_date": "2024-01-20T00:00:00Z",
      "transaction_id": "txn_123456789",
      "notes": "Payment received via Stripe",
      "created_at": "2024-01-15T10:55:00Z",
      "updated_at": "2024-01-15T10:55:00Z"
    }
  ]
}
```

### Send Invoice (Placeholder)
```bash
curl -X POST http://localhost:8080/api/v1/invoices/INVOICE_ID/send \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "Invoice send functionality will be implemented"
}
```

### Get Invoice PDF (Placeholder)
```bash
curl -X GET http://localhost:8080/api/v1/invoices/INVOICE_ID/pdf \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "message": "PDF generation will be implemented"
}
```

---

## 4. Expenses

### List All Expenses
```bash
curl -X GET http://localhost:8080/api/v1/expenses \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
[
  {
    "id": "aa0e8400-e29b-41d4-a716-446655440006",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "client_id": "660e8400-e29b-41d4-a716-446655440001",
    "description": "Office Supplies",
    "amount": 150.00,
    "currency": "USD",
    "category": "office",
    "expense_date": "2024-01-10T00:00:00Z",
    "receipt_url": "https://example.com/receipts/receipt-001.pdf",
    "notes": "Purchased office supplies for project",
    "created_at": "2024-01-15T11:00:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
]
```

### List Expenses with Filters
```bash
# Filter by client
curl -X GET "http://localhost:8080/api/v1/expenses?client_id=CLIENT_ID" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by category
curl -X GET "http://localhost:8080/api/v1/expenses?category=travel" \
  -H "Authorization: Bearer YOUR_TOKEN"

# Filter by date range
curl -X GET "http://localhost:8080/api/v1/expenses?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):** Same format as List All Expenses, but filtered.

### Create Expense
```bash
curl -X POST http://localhost:8080/api/v1/expenses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "client_id": "CLIENT_ID",
    "description": "Office Supplies",
    "amount": 150.00,
    "currency": "USD",
    "category": "office",
    "expense_date": "2024-01-10T00:00:00Z",
    "receipt_url": "https://example.com/receipts/receipt-001.pdf",
    "notes": "Purchased office supplies for project"
  }'
```

**Response (201 Created):**
```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440006",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "description": "Office Supplies",
  "amount": 150.00,
  "currency": "USD",
  "category": "office",
  "expense_date": "2024-01-10T00:00:00Z",
  "receipt_url": "https://example.com/receipts/receipt-001.pdf",
  "notes": "Purchased office supplies for project",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

### Get Expense by ID
```bash
curl -X GET http://localhost:8080/api/v1/expenses/EXPENSE_ID \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440006",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "description": "Office Supplies",
  "amount": 150.00,
  "currency": "USD",
  "category": "office",
  "expense_date": "2024-01-10T00:00:00Z",
  "receipt_url": "https://example.com/receipts/receipt-001.pdf",
  "notes": "Purchased office supplies for project",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:00:00Z"
}
```

**Error Response (404):**
```json
{
  "error": "expense not found"
}
```

### Update Expense
```bash
curl -X PUT http://localhost:8080/api/v1/expenses/EXPENSE_ID \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "client_id": "CLIENT_ID",
    "description": "Office Supplies (Updated)",
    "amount": 175.00,
    "currency": "USD",
    "category": "office",
    "expense_date": "2024-01-10T00:00:00Z",
    "receipt_url": "https://example.com/receipts/receipt-001-updated.pdf",
    "notes": "Updated expense with additional items"
  }'
```

**Response (200 OK):**
```json
{
  "id": "aa0e8400-e29b-41d4-a716-446655440006",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "description": "Office Supplies (Updated)",
  "amount": 175.00,
  "currency": "USD",
  "category": "office",
  "expense_date": "2024-01-10T00:00:00Z",
  "receipt_url": "https://example.com/receipts/receipt-001-updated.pdf",
  "notes": "Updated expense with additional items",
  "created_at": "2024-01-15T11:00:00Z",
  "updated_at": "2024-01-15T11:05:00Z"
}
```

---

## 5. Reports

### Get Summary Report
```bash
# All time
curl -X GET http://localhost:8080/api/v1/reports/summary \
  -H "Authorization: Bearer YOUR_TOKEN"

# With date range
curl -X GET "http://localhost:8080/api/v1/reports/summary?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "total_revenue": 7562.50,
  "total_expenses": 325.00,
  "net_profit": 7237.50,
  "outstanding_invoices": 0,
  "paid_invoices": 1,
  "total_invoices": 1
}
```

### Get Client Profitability Report
```bash
# All time
curl -X GET http://localhost:8080/api/v1/reports/client-profit/CLIENT_ID \
  -H "Authorization: Bearer YOUR_TOKEN"

# With date range
curl -X GET "http://localhost:8080/api/v1/reports/client-profit/CLIENT_ID?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "client_id": "660e8400-e29b-41d4-a716-446655440001",
  "client_name": "Acme Corporation",
  "total_revenue": 7562.50,
  "total_expenses": 150.00,
  "net_profit": 7412.50,
  "profit_margin": 97.98
}
```

**Error Response (404):**
```json
{
  "error": "client not found"
}
```

### Get Tax Summary
```bash
curl -X GET "http://localhost:8080/api/v1/reports/tax-summary?from_date=2024-01-01&to_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Response (200 OK):**
```json
{
  "period": "2024-01 to 2024-01",
  "total_revenue": 7562.50,
  "total_expenses": 325.00,
  "net_income": 7237.50,
  "invoices": [
    {
      "invoice_number": "INV-001",
      "date": "2024-01-15T00:00:00Z",
      "client_name": "",
      "amount": 7562.50,
      "tax_amount": 687.50
    }
  ],
  "expenses": [
    {
      "description": "Office Supplies",
      "date": "2024-01-10T00:00:00Z",
      "category": "office",
      "amount": 150.00
    },
    {
      "description": "Travel Expenses",
      "date": "2024-01-12T00:00:00Z",
      "category": "travel",
      "amount": 500.00
    }
  ]
}
```

**Error Response (400):**
```json
{
  "error": "from_date and to_date are required"
}
```

---

## Quick Test Script

You can also use the automated test script:

```bash
chmod +x test_api.sh
./test_api.sh
```

This script will:
1. Register/login and get a token
2. Test all endpoints in sequence
3. Use the IDs from previous responses for subsequent requests

---

## Common Error Responses

### Authentication Error (401)
```json
{
  "error": "unauthorized"
}
```

### Not Found Error (404)
```json
{
  "error": "client not found"
}
```

### Bad Request Error (400)
```json
{
  "error": "invalid payload"
}
```

or

```json
{
  "error": "email is required"
}
```

### Internal Server Error (500)
```json
{
  "error": "internal server error"
}
```

---

## Notes

- All protected endpoints require the `Authorization: Bearer YOUR_TOKEN` header
- Date formats: Use ISO 8601 format (e.g., `2024-01-15T00:00:00Z`) for datetime fields
- Date filters: Use `YYYY-MM-DD` format (e.g., `2024-01-15`)
- Invoice statuses: `draft`, `pending`, `paid`, `overdue`, `cancelled`
- Currency defaults to `USD` if not specified
- Invoice items are required when creating/updating invoices

