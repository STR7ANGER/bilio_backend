# BillStack

> Invoicing & Expense Management for Agencies and Freelancers

BillStack is a complete billing workspace that helps agencies and freelancers create branded invoices, track payments, manage expenses, and understand profitability—all in one place.

## 🎯 Problem We Solve

Agencies and freelancers waste hours every week juggling:
- Spreadsheets for tracking invoices
- Manual payment reminders via email
- Scattered expense receipts
- Multiple payment gateways
- No clear view of per-client profitability

**BillStack consolidates everything into a single, powerful workspace.**

## ✨ Key Features

### Invoicing
- **Branded Invoice Creation** - Professional PDFs with your logo and company details
- **Payment Links** - Accept payments via Stripe & Razorpay (credit card, UPI, wallets)
- **Recurring Invoices** - Automate monthly/quarterly retainers
- **Smart Reminders** - Automatic overdue notifications with escalating templates
- **Multi-Currency** - Support for INR, USD, and more

### Expense Tracking
- **Receipt Capture** - Upload and attach receipts to expenses
- **Client/Project Linking** - Track expenses per client for accurate profitability
- **Categorization** - Organize expenses by category for tax reporting

### Financial Intelligence
- **Client Profitability** - See P&L breakdown per client
- **Tax Summaries** - Export-ready CSV/PDF reports for accountants
- **Revenue Tracking** - Monitor MRR, outstanding invoices, and cash flow
- **Dashboard Analytics** - Visual overview of your business health

## 🏗️ Tech Stack

### Frontend
- **Next.js** (React) - Modern, fast web application
- **Tailwind CSS** - Beautiful, responsive design
- **shadcn/ui** - Polished component library

### Backend
- **Golang** - High-performance REST API
- **PostgreSQL** - Reliable, scalable database
- **Redis** - Job queue and caching
- **AWS S3** - Secure file storage for receipts and PDFs

### Integrations
- **Stripe** - Global payment processing
- **Razorpay** - India-specific payments (UPI, cards, wallets)
- **SendGrid/SES** - Transactional emails
- **Puppeteer** - Professional PDF generation

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker (optional but recommended)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/billstack.git
cd billstack

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Start with Docker Compose (recommended)
docker-compose up -d

# Or run manually:

# Backend
cd backend
go mod download
go run cmd/server/main.go

# Frontend
cd frontend
npm install
npm run dev
```

### Environment Variables

```env
# Database
DATABASE_URL=postgres://user:pass@localhost:5432/billstack

# Redis
REDIS_URL=redis://localhost:6379

# AWS S3
S3_BUCKET=billstack-dev
S3_REGION=ap-south-1
S3_ACCESS_KEY_ID=your_key
S3_SECRET_ACCESS_KEY=your_secret

# Payment Providers
STRIPE_SECRET_KEY=sk_test_...
RAZORPAY_KEY_ID=rzp_test_...
RAZORPAY_KEY_SECRET=...

# Email
SENDGRID_API_KEY=SG....

# Auth
JWT_SECRET=your_secure_secret

# App
APP_URL=http://localhost:3000
```

## 📚 API Documentation

The REST API is available at `/api/v1/`. Key endpoints:

### Authentication
- `POST /api/v1/auth/register` - Create new account
- `POST /api/v1/auth/login` - Login and get JWT token

### Clients
- `GET /api/v1/clients` - List all clients
- `POST /api/v1/clients` - Create new client
- `GET /api/v1/clients/{id}` - Get client details
- `PUT /api/v1/clients/{id}` - Update client
- `DELETE /api/v1/clients/{id}` - Delete client

### Invoices
- `GET /api/v1/invoices` - List invoices (filter by status, client, date)
- `POST /api/v1/invoices` - Create new invoice
- `GET /api/v1/invoices/{id}` - Get invoice details
- `PUT /api/v1/invoices/{id}` - Update invoice (draft/pending only)
- `POST /api/v1/invoices/{id}/send` - Send invoice via email
- `POST /api/v1/invoices/{id}/mark-paid` - Manually mark as paid
- `GET /api/v1/invoices/{id}/pdf` - Get PDF download link

### Expenses
- `GET /api/v1/expenses` - List all expenses
- `POST /api/v1/expenses` - Record new expense
- `GET /api/v1/expenses/{id}` - Get expense details
- `PUT /api/v1/expenses/{id}` - Update expense

### Reports
- `GET /api/v1/reports/summary` - Revenue, expenses, profit overview
- `GET /api/v1/reports/client-profit/{id}` - Per-client profitability
- `GET /api/v1/reports/tax-summary` - Export for tax filing

Full API documentation: [Link to Swagger/OpenAPI spec]

## 🗄️ Database Schema

Core tables:
- `users` - Workspace owners
- `clients` - Customer records
- `invoices` - Invoice headers
- `invoice_items` - Line items
- `payments` - Payment records
- `expenses` - Expense tracking
- `recurring_invoices` - Automated invoice templates
- `audit_log` - Change tracking

See `/backend/migrations/` for complete schema definitions.

## 🔒 Security

- **HTTPS Only** - All traffic encrypted via TLS
- **JWT Authentication** - Secure token-based auth
- **Password Hashing** - Bcrypt for secure password storage
- **Webhook Verification** - Signature validation for payment webhooks
- **Private File Storage** - S3 presigned URLs with expiration
- **PCI Compliance** - Payment processing via Stripe/Razorpay (no card data stored)
- **Audit Logging** - Track all invoice and payment changes

## 💰 Pricing

### Free
- 3 clients
- 5 invoices/month
- Basic export

### Starter - $20/month
- Unlimited invoices
- 25 clients
- Recurring invoices
- Email support

### Pro - $40/month
- Unlimited clients
- Priority support
- Advanced reports
- Tax export tools

### Enterprise - $200+/month
- Multi-user teams
- Custom branding
- Dedicated onboarding
- API access

**14-day free trial on all paid plans**

## 🛠️ Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Database Migrations

```bash
# Run migrations
migrate -path backend/migrations -database $DATABASE_URL up

# Rollback
migrate -path backend/migrations -database $DATABASE_URL down 1
```

### Local Development

```bash
# Watch mode for backend
cd backend
air # or go run with file watcher

# Watch mode for frontend
cd frontend
npm run dev
```

## 📈 Roadmap

### Q1 2025 (MVP)
- ✅ Core invoicing & expenses
- ✅ Payment integration (Stripe & Razorpay)
- ✅ PDF generation
- ✅ Recurring invoices
- ✅ Basic reports

### Q2 2025
- Multi-user workspaces & roles
- Client portal (view/download invoices)
- Mobile app (React Native)
- More payment gateways

### Q3 2025
- Accounting integrations (QuickBooks, Zoho)
- WhatsApp/SMS reminders
- Advanced tax templates (GST, 1099)
- Time tracking integration

### Q4 2025
- API for third-party integrations
- White-label options
- AI-powered expense categorization
- Payroll integration

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## 📞 Support

- **Documentation**: [docs.billstack.com](https://docs.billstack.com)
- **Email**: support@billstack.com
- **Twitter**: [@billstack](https://twitter.com/billstack)
- **Discord**: [Join our community](https://discord.gg/billstack)

## 🙏 Acknowledgments

Built with ❤️ for agencies and freelancers who deserve better tools.

Special thanks to:
- The open-source community
- Our beta testers
- Early adopters who provided invaluable feedback

---

**Ready to simplify your billing?** [Start your free trial →](https://app.billstack.com/signup)