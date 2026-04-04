# Finance Processing API

A backend REST API for a finance dashboard system built with Go, Fiber, GORM, and PostgreSQL. Supports role-based access control, financial record management, and rich dashboard analytics with flexible filtering.


## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| Framework | Fiber v2 |
| ORM | GORM |
| Database | PostgreSQL 16 |
| Auth | JWT (golang-jwt) |
| Migrations | golang-migrate |
| Config | koanf |
| Logging | zerolog |
| Containerization | Docker + Docker Compose |
---



---

## Project Structure
```
.
├── cmd/
│   ├── server/main.go            # server entry point
│   └── seed/main.go              # seed entry point
├── docker/
│   ├── Dockerfile.server
│   ├── Dockerfile.seed
│   └── .env.docker               # env for docker
├── internal/
│   ├── config/                   # env config loading  
│   ├── database/
│   │   └── db.go                 # db connection + migrations
│   ├── handlers/                 # HTTP request handlers
│   ├── middleware/               # JWT auth + logging
│   ├── models/                   # GORM models + response DTOs
│   ├── policy/                   # role-based access control
│   ├── repository/               # database queries
│   ├── routes/                   # route registration
│   ├── seed/                     # seed logic
│   ├── server/                   # server wiring
│   └── services/                 # business logic
├── migrations/
├── docker-compose.yml
├── Makefile
└── .env.example
```




## containers
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
- `finance_seed` only runs when explicitly called — never auto-runs with the server
- Data persists in `postgres_data` volume across restarts
- `make down` keeps data intact, `make down-v` wipes it

---



---

## Roles

| Role | Permissions |
|---|---|
| `admin` | full access — manage users, create/edit/delete records, view dashboard |
| `analyst` | create/edit/delete records, view dashboard |
| `viewer` | read-only access to records and dashboard |

---


## Getting Started

### Prerequisites

- Docker and Docker Compose

### 1. Clone the repo
```bash
git clone https://github.com/yourusername/finance-processing.git
cd finance-processing
```

### 2. Setup environment
```bash
cp .env.example .env
```

Edit `.env` and set env variables to your credentials eg`AUTH_JWT_SECRET` to a long random string.

### 3. Build images
```bash
make build
```

### 4. Start the server
```bashmake dev        # run seeding and server  (no docker)
make run
```

This starts PostgreSQL and the API server. Migrations run automatically.

### 5. Verify it's running
```bash
curl http://localhost:8080/health
# {"status":"ok"}
```

### 6. Seed test data (optional)
```bash
make seed
```

Creates 3 test users and 20 financial records.

---

## Test Credentials (after seeding)

| Role | Email | Password |
|---|---|---|
| Admin | admin@test.com | password123 |
| Analyst | analyst@test.com | password123 |
| Viewer | viewer@test.com | password123 |

---

## Common Commands
```bash
make up         # start postgres + server
make seed       # run seed job (docker)
make build      # build docker images
make down       # stop containers (keep data)
make down-v     # stop containers and wipe database
make logs       # view server logs

make dev        # run seeding and server  (no docker)
```

---

``` bash
make dev        # run this for complete setup in one command after entering the env
```
## API Reference

All protected endpoints require:
```
Authorization: Bearer <token>
```

---

### Auth

#### POST `/auth/login`
```json
{
    "email": "admin@test.com",
    "password": "password123"
}
```
Response `200`:
```json
{
    "token": "eyJhbGci..."
}
```

---

### Users
> Admin only

#### POST `/api/users` — create user
```json
{
    "name": "John Doe",
    "email": "john@test.com",
    "password": "password123",
    "role": "analyst"
}
```
Valid roles: `admin` `analyst` `viewer`

Response `201`:
```json
{
    "message": "user created successfully"
}
```

#### GET `/api/users` — list all users

No body. Returns array of users without passwords.

#### GET `/api/users/:id` — get single user

No body.

#### PATCH `/api/users/:id/role` — update role
```json
{ "role": "viewer" }
```

#### PATCH `/api/users/:id/status` — activate or deactivate
```json
{ "active": false }
```

---

### Transactions
> Read: all roles — Write: admin + analyst

#### POST `/api/transactions` — create transaction
```json
{
    "amount": 5000.00,
    "type": "income",
    "category": "salary",
    "description": "Monthly salary",
    "date": "2026-04-01"
}
```
Valid types: `income` `expense`

Response `201`:
```json
{
    "data": {
        "id": "uuid",
        "created_by": "uuid",
        "amount": 5000,
        "type": "income",
        "category": "salary",
        "description": "Monthly salary",
        "date": "2026-04-01T00:00:00Z",
        "created_at": "2026-04-05T10:00:00Z"
    }
}
```

#### GET `/api/transactions` — list with filters

All params optional and combinable:
```
?type=expense
?type=income
?category=salary
?from=2026-01-01
?to=2026-03-31
?type=expense&category=rent&from=2026-01-01&to=2026-03-31
```

#### GET `/api/transactions/:id` — get single transaction

#### PUT `/api/transactions/:id` — update transaction

Same body as POST.

#### DELETE `/api/transactions/:id` — soft delete

No body. Record is never permanently removed.

---

### Dashboard
> All authenticated roles

All dashboard endpoints accept these optional query params:

| Param | Example | Description |
|---|---|---|
| `from` | `2026-01-01` | filter from date |
| `to` | `2026-03-31` | filter to date |
| `type` | `expense` | income or expense only |
| `category` | `rent` | single category |
| `category` | `rent&category=salary` | multiple categories |
| `limit` | `10` | recent transactions limit |
| `offset` | `0` | recent transactions offset |

#### GET `/api/dashboard` — full dashboard in one call
```
GET /api/dashboard
GET /api/dashboard?from=2026-01-01&to=2026-03-31
GET /api/dashboard?type=expense&category=rent&category=groceries
GET /api/dashboard?from=2026-01-01&to=2026-03-31&type=expense&limit=5
```

Response `200`:
```json
{
    "data": {
        "summary": {
            "total_income": 17150.00,
            "total_expenses": 6080.00,
            "net_balance": 11070.00
        },
        "categories": [
            { "category": "salary",    "type": "income",  "total": 12850.00 },
            { "category": "rent",      "type": "expense", "total": 3350.00  },
            { "category": "freelance", "type": "income",  "total": 3200.00  }
        ],
        "trends": [
            { "month": "2026-04", "income": 4900.00, "expense": 1520.00 },
            { "month": "2026-03", "income": 8250.00, "expense": 2850.00 }
        ],
        "recent": [
            {
                "id": "uuid",
                "amount": 220.00,
                "type": "expense",
                "category": "entertainment",
                "date": "2026-04-02T00:00:00Z"
            }
        ]
    }
}
```

#### GET `/api/dashboard/summary` — totals only
#### GET `/api/dashboard/categories` — category breakdown
#### GET `/api/dashboard/trends` — monthly income vs expense

---

## Access Control

| Endpoint | viewer | analyst | admin |
|---|---|---|---|
| `POST /auth/login` | yes | yes | yes |
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

## Error Responses

All errors follow this shape:
```json
{ "error": "description of what went wrong" }
```

| Status | When |
|---|---|
| `400` | invalid input or missing required fields |
| `401` | missing or invalid JWT token |
| `403` | role does not have permission |
| `404` | resource not found |
| `409` | resource already exists |
| `500` | internal server error |

---

## Assumptions

- Financial records are **shared across all users** — this is a company finance tool, not a personal finance app
- `created_by` on a transaction tracks who entered the record, not ownership — any authenticated analyst or admin can view and edit all records
- Soft delete is used for both users and transactions — data is never permanently removed from the database
- JWT tokens expire after 24 hours
- The first admin is created via `make seed` — after that, admins create additional users via `POST /api/users`
- Rate limiting is applied globally at 20 requests per 30 second sliding window

---

## Tradeoffs

**GORM over raw SQL** — chosen for development speed and readability. For high-traffic production, `sqlc` with raw SQL would give better performance and compile-time query safety.

**Single `/api/dashboard` endpoint** — returns summary, categories, trends, and recent in one response to minimize frontend network calls. Individual sub-endpoints (`/summary`, `/categories`, `/trends`) are also available for targeted queries.

**In-process migrations** — `golang-migrate` runs automatically on server startup. In a production CI/CD pipeline these would run as a separate step before deployment to avoid race conditions on multi-instance deploys.

**Profiles for seed** — the seed container uses Docker Compose profiles so it never runs automatically with the server. This prevents accidental re-seeding on every `make run`.

---

## What Could Be Added

- **Natural language queries** — `POST /api/dashboard/ask` with `{"query": "show me rent expenses from last quarter"}` — Claude API parses intent and maps to existing dashboard filters
- Pagination on transaction listing
- Export to CSV
- Email notifications on large transactions
- Webhook support for real-time updates
- Unit and integration tests

