# subtrack

REST API for tracking personal subscriptions — vendors, costs, billing platforms, renewal dates. Go + Gin + GORM + PostgreSQL, JWT auth.

## Run

```bash
cp .env.example .env   # set DB_PASSWORD and JWT_SECRET
make run               # needs Go 1.21+ and a running PostgreSQL
# or
make docker-run
```

Server listens on `http://localhost:8080`.

## API

Register at `POST /api/v1/auth/register`, login at `POST /api/v1/auth/login`, then use the token as `Authorization: Bearer <token>` for CRUD on `/api/v1/subscriptions`. Listing supports pagination and filters like `vendor`, `status`, `billing_cycle`, `min_amount`, and `renewing_soon=30` (days).

Health checks at `/health` and `/ready`, Prometheus metrics at `/metrics`.
