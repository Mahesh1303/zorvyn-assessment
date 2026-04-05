# Finance Processing API

A backend REST API for a finance dashboard system built with Go. Supports role-based access control, financial record management, and rich dashboard analytics with flexible filtering.

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.26.1 |
| Framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Auth | JWT (golang-jwt) |
| Migrations | golang-migrate |
| Config | koanf |
| Logging | zerolog |
| Containers | Docker + Docker Compose |

---

## Project Structure
```
.
├── cmd/
│   ├── server/main.go        # server entry point
│   └── seed/main.go          # seed entry point
├── docker/
│   ├── Dockerfile.server
│   ├── Dockerfile.seed
├── internal/
│   ├── config/               # env config loading
│   ├── database/             # db connection + migrations
│   ├── handlers/             # HTTP request handlers
│   ├── middleware/           # JWT auth + logging
│   ├── models/               # GORM models + response DTOs
│   ├── policy/               # role-based access control
│   ├── repository/           # database queries
│   ├── routes/               # route registration
│   ├── seed/                 # seed logic
│   ├── server/               # server wiring
│   └── services/             # business logic
├── docs/                     # documentation
│   └── api/                  # API reference docs
├── migrations/               # SQL migration files (000001_init.up.sql etc)
├── docker-compose.yml
├── Makefile
└── .env.example
```

---

## Request Flow Architecture

Each API request follows a layered flow:

`route -> handler -> service -> repository`

- `routes` map HTTP endpoints to handler functions
- `handlers` parse/validate request input and build HTTP responses
- `services` enforce business rules and role/policy checks
- `repositories` run GORM queries against PostgreSQL

This separation keeps transport logic, business logic, and data access cleanly isolated.

---

## Middleware And Security

The app applies middleware in this order:

1. `Rate Limiter` (global): sliding window, max `20` requests per `30s`
2. `Logging` (global): logs request `method`, `path`, `status`, and `latency`
3. `Error Handler` (global): converts unhandled errors into JSON responses
4. `Auth` (for `/api/*`): validates Bearer JWT and injects user context

Route protection:

- Public routes: `/health`, `/auth/login`, `/auth/register-admin`
- Protected routes: all `/api/*` endpoints

---

## Operational Defaults

- Pagination defaults: `limit=10`, `offset=0`
- Pagination bounds: `limit` max `100`, `offset` cannot be negative
- JWT token expiry: `24 hours`
- Database migrations run automatically when the server starts

---

## Container Architecture
```
┌─────────────────────────────────────────────────────┐
│                  docker network                      │
│                                                      │
│  ┌──────────────────┐       ┌───────────────────┐   │
│  │  finance_server  │──────▶│  finance_postgres  │   │
│  │  ./server :8080  │       │   postgres:16      │   │
│  └──────────────────┘       │   :5432            │   │
│                             │                    │   │
│  ┌──────────────────┐       │  [postgres_data]   │   │
│  │  finance_seed    │──────▶│   volume persists  │   │
│  │  exits after run │       └───────────────────┘   │
│  └──────────────────┘                               │
│   * profile: seed only                               │
└─────────────────────────────────────────────────────┘
          ▲
          │ HTTP :8080
          │
       client
```

- `finance_server` runs migrations on startup then serves the API
- `finance_seed` only runs when explicitly called — never auto-starts with the server
- Data persists in `postgres_data` volume across restarts
- `make down` keeps data, `make down-v` wipes everything

---

## Roles

| Role | Permissions |
|---|---|
| `admin` | full access — manage users, records, and dashboard |
| `analyst` | create and manage records, view dashboard |
| `viewer` | read-only access to records and dashboard |

---

## Getting Started

### Docker based

- Docker and Docker Compose installed

### 1. Clone the repo
```bash
git clone https://github.com/Mahesh1303/zorvyn-assessment.git
cd zorvyn-assessment
```

### 2. Set up environment
```bash
cp .env.example .env
```

Open `.env` and fill in your values. At minimum set `AUTH_JWT_SECRET` to a long random string.

### 3. Build and start
```bash
make build   # build docker images
make run     # start postgres + server
```

Migrations run automatically on first start.

### 4. Verify
```bash
curl http://localhost:8080/health
# {"status":"Hello from server"}
```

### 5. Seed test data
```bash
make seed
```

Creates 3 test users and 20 financial records across different categories and date ranges.

---



## Without Docker (local dev)
```bash
# make sure postgres is running locally and .env has correct DB_URL
make dev        # runs server locally It runs seed first and then server
```

---

## Test Credentials

After running `make seed`:

| Role | Email | Password |
|---|---|---|
| Admin | admin@test.com | password123 |
| Analyst | analyst@test.com | password123 |
| Viewer | viewer@test.com | password123 |

---

## Commands
```bash
make build      # build docker images
make up         # start postgres + server in background
make seed       # run seed container (test data)
make down       # stop containers, keep database volume
make down-v     # stop containers, wipe database volume
make logs       # tail server logs
make dev        # run server locally without docker (seed + server)
make help       # help
```
---

## API Docs

---

## API Docs

Full API documentation is in [`docs/api/`](docs/api/):

- [Auth](docs/api/auth.md) — login
- [Users](docs/api/users.md) — user management (admin only)
- [Transactions](docs/api/transactions.md) — financial records CRUD + filters
- [Dashboard](docs/api/dashboard.md) — analytics and summaries
- [Errors](docs/api/errors.md) — error response reference

---

## Access Control Summary

| Endpoint | viewer | analyst | admin |
|---|---|---|---|
| `POST /auth/login` | yes | yes | yes |
| `POST /auth/register-admin` | yes | yes | yes |
| `POST /api/users` | no | no | yes |
| `GET /api/users` | no | no | yes |
| `PATCH /api/users/:id/role` | no | no | yes |
| `PATCH /api/users/:id/status` | no | no | yes |
| `GET /api/transactions` | yes | yes | yes |
| `GET /api/transactions/:id` | yes | yes | yes |
| `POST /api/transactions` | no | yes | yes |
| `PUT /api/transactions/:id` | no | yes | yes |
| `DELETE /api/transactions/:id` | no | yes | yes |
| `GET /api/dashboard/*` | yes | yes | yes |

---

## Assumptions

- Financial records are **shared across all users** — this is a company finance tool, not a personal finance app
- `created_by` tracks who entered the record, not ownership — all analysts and admins can view and edit all records
- Soft delete is used everywhere — data is never permanently removed
- JWT tokens expire after 24 hours
- The first admin is created via `make seed` — after that, admins create users via `POST /api/users`
- Rate limiting is applied globally at 20 requests per 30 second sliding window

---

## Tradeoffs

**GORM over raw SQL** — chosen for development speed and readability. For high-traffic production, `sqlc` with raw SQL would give better performance and compile-time query safety.

**Single `/api/dashboard` endpoint** — returns summary, categories, trends, and recent transactions in one response to minimize frontend network calls. Individual sub-endpoints are also available for targeted queries.

**In-process migrations** — `golang-migrate` runs on server startup automatically. In a production CI/CD pipeline these would run as a separate step before deployment.

**Seed via separate binary** — seed runs as its own Docker container with its own entrypoint. This keeps the server binary clean and makes seeding explicit rather than automatic.

---

## What Could Be Added

- **Natural language queries** — `POST /api/dashboard/ask` accepting `{"query": "show me rent expenses from last quarter"}` — the Claude API would parse the intent and map it to the existing dashboard filter system automatically
- CSV export for records
- Unit and integration tests
- Email alerts on large transactions
- Webhook support for real-time updates