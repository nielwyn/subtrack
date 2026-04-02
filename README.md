# subtrack

A Go backend for tracking SaaS subscriptions — manage vendors, plans, seats, costs, and renewal dates through a REST API.

[![CI/CD Pipeline](https://github.com/nielwyn/subtrack/actions/workflows/ci.yml/badge.svg)](https://github.com/nielwyn/subtrack/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nielwyn/subtrack)](https://goreportcard.com/report/github.com/nielwyn/subtrack)

## Features

- JWT-based authentication with bcrypt password hashing
- Full CRUD for SaaS subscriptions
- Filter by vendor, status, billing cycle, seat capacity, and upcoming renewals
- Pagination on all list endpoints
- Per-IP rate limiting
- Request ID tracing on every request
- Structured logging via Zap
- Prometheus metrics endpoint
- Health and readiness endpoints
- Graceful shutdown
- Docker and Docker Compose support
- GitHub Actions CI/CD

## Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin) v1.9.1
- **ORM**: [GORM](https://gorm.io/) v1.25.5
- **Database**: PostgreSQL
- **Auth**: [JWT](https://github.com/golang-jwt/jwt) v5.2.0
- **Logging**: [Zap](https://github.com/uber-go/zap) v1.26.0
- **Metrics**: [Prometheus](https://github.com/prometheus/client_golang) v1.18.0
- **Rate Limiting**: [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time/rate)

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker and Docker Compose (optional)

### Local Setup

1. Clone the repo
   ```bash
   git clone https://github.com/nielwyn/subtrack.git
   cd subtrack
   ```

2. Set up environment variables
   ```bash
   cp .env.example .env
   # Edit .env — set DB_PASSWORD and JWT_SECRET
   ```

3. Install dependencies
   ```bash
   make deps
   ```

4. Start PostgreSQL and run
   ```bash
   make run
   ```

API available at `http://localhost:8080`

### Docker

```bash
cp .env.example .env
make docker-run     # Start PostgreSQL + API
make docker-logs    # View logs
make docker-down    # Stop services
```

## API

### Base URL
```
http://localhost:8080
```

### Response Format

```json
{ "success": true, "message": "...", "data": { ... } }
{ "success": false, "code": "NOT_FOUND", "message": "subscription not found" }
```

### Error Codes

| Code | Meaning |
|------|---------|
| `INVALID_INPUT` | Validation failed |
| `UNAUTHORIZED` | Missing or invalid token |
| `NOT_FOUND` | Resource does not exist |
| `RATE_LIMITED` | Too many requests |
| `USERNAME_EXISTS` | Username already taken |
| `EMAIL_EXISTS` | Email already registered |
| `INTERNAL_ERROR` | Server error |

### Health

| Method | Endpoint | Auth |
|--------|----------|------|
| GET | /health | No |
| GET | /ready | No |
| GET | /metrics | No |

### Auth

| Method | Endpoint | Auth |
|--------|----------|------|
| POST | /api/v1/auth/register | No |
| POST | /api/v1/auth/login | No |

### Subscriptions (JWT required)

Include `Authorization: Bearer <token>` on all requests.

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/subscriptions | Create subscription |
| GET | /api/v1/subscriptions | List subscriptions |
| GET | /api/v1/subscriptions/:id | Get by ID |
| PUT | /api/v1/subscriptions/:id | Update subscription |
| DELETE | /api/v1/subscriptions/:id | Delete subscription |

#### Subscription fields

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | e.g. "GitHub Teams" |
| `vendor` | string | e.g. "GitHub" |
| `plan` | string | e.g. "Teams", "Enterprise" |
| `billing_cycle` | string | `monthly` or `annual` |
| `cost_per_seat` | float | Cost per user per billing cycle |
| `seats` | int | Total seats purchased |
| `current_users` | int | Active users assigned |
| `auto_renews` | bool | Whether it auto-renews |
| `renewal_date` | datetime | Next renewal date |
| `status` | string | `active`, `trial`, or `cancelled` |
| `notes` | string | Free text |

#### Filters

```
GET /api/v1/subscriptions?vendor=GitHub&status=active&billing_cycle=annual&low_capacity=5&renewing_soon=30
```

| Param | Description |
|-------|-------------|
| `page` | Page number (default: 1) |
| `limit` | Page size (default: 20, max: 100) |
| `vendor` | Filter by vendor name |
| `status` | Filter by status |
| `billing_cycle` | `monthly` or `annual` |
| `min_cost` | Minimum cost per seat |
| `max_cost` | Maximum cost per seat |
| `low_capacity` | Subscriptions where available seats ≤ this value |
| `renewing_soon` | Subscriptions renewing within N days |

#### Example

```bash
# Create a subscription
curl -X POST http://localhost:8080/api/v1/subscriptions \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "GitHub Teams",
    "vendor": "GitHub",
    "plan": "Teams",
    "billing_cycle": "annual",
    "cost_per_seat": 4.00,
    "seats": 50,
    "current_users": 38,
    "auto_renews": true,
    "renewal_date": "2026-12-01T00:00:00Z",
    "status": "active"
  }'
```

## Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| SERVER_HOST | Server host | 0.0.0.0 | No |
| SERVER_PORT | Server port | 8080 | No |
| GIN_MODE | `debug` or `release` | debug | No |
| DB_HOST | PostgreSQL host | localhost | Yes |
| DB_PORT | PostgreSQL port | 5432 | No |
| DB_USER | Database user | postgres | Yes |
| DB_PASSWORD | Database password | — | Yes |
| DB_NAME | Database name | inventory_db | Yes |
| DB_SSLMODE | SSL mode | disable | No |
| JWT_SECRET | JWT signing secret | — | Yes |
| JWT_EXPIRY_HOURS | Token expiry in hours | 24 | No |
| LOG_LEVEL | `debug` / `info` / `error` | debug | No |
| LOG_ENCODING | `json` or `console` | json | No |

## Make Commands

```bash
make run          # Run locally
make build        # Build binary
make test         # Run tests
make deps         # Download dependencies
make docker-run   # Start with Docker Compose
make docker-down  # Stop Docker Compose
make docker-logs  # View logs
make lint         # Run linter
make fmt          # Format code
make vet          # Run go vet
```

