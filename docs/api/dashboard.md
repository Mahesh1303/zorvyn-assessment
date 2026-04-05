# Dashboard

Base URL: `http://localhost:8080/api`

All endpoints require:
```
Authorization: Bearer <token>
```

> All authenticated roles can access dashboard endpoints except `/dashboard/analytics`, which is restricted to admin and analyst.

---

## Shared Query Params

Every dashboard endpoint accepts these optional filters. All params are combinable.

| Param | Example | Description |
|---|---|---|
| `from` | `2026-01-01` | filter from date (inclusive) |
| `to` | `2026-03-31` | filter to date (inclusive) |
| `type` | `expense` | `income` or `expense` only |
| `category` | `rent` | single category |
| `category` | `rent&category=salary` | multiple categories (repeat param) |
| `limit` | `10` | number of recent transactions (default 10, max 100) |
| `offset` | `0` | offset for recent transactions pagination |

---

## GET `/dashboard` — full dashboard

Returns all dashboard data in a single call — summary, category totals, monthly trends, and recent transactions. Use this to load a full dashboard page in one request.

### Examples
```
GET /api/dashboard
GET /api/dashboard?type=expense
GET /api/dashboard?from=2026-01-01&to=2026-03-31
GET /api/dashboard/categories?type=expense&type=rent
GET /api/dashboard/analytics?from=2026-01-01&to=2026-03-31
GET /api/dashboard?category=rent&category=salary&category=groceries
GET /api/dashboard/summary?from=2026-01-01&to=2026-03-31&type=expense
GET /api/dashboard?from=2026-01-01&to=2026-03-31&type=expense&category=rent&limit=5&offset=0


```

### Response `200`
```json
{
    "data": {
        "summary": {
            "total_income": 17150.00,
            "total_expenses": 6080.00,
            "net_balance": 11070.00
        },
        "categories": [
            { "category": "salary",        "type": "income",  "total": 12850.00 },
            { "category": "rent",          "type": "expense", "total": 3350.00  },
            { "category": "freelance",     "type": "income",  "total": 3200.00  },
            { "category": "groceries",     "type": "expense", "total": 1050.00  },
            { "category": "software",      "type": "expense", "total": 650.00   }
        ],
        "trends": [
            { "month": "2026-04", "income": 4900.00,  "expense": 1520.00 },
            { "month": "2026-03", "income": 8250.00,  "expense": 2850.00 },
            { "month": "2026-02", "income": 5000.00,  "expense": 1750.00 }
        ],
        "recent": [
            {
                "id": "uuid",
                "created_by": "uuid",
                "amount": 220.00,
                "type": "expense",
                "category": "entertainment",
                "description": "Weekend outing",
                "date": "2026-04-02T00:00:00Z",
                "created_at": "2026-04-02T10:00:00Z"
            }
        ]
    }
}
```

---

## GET `/dashboard/summary` — totals only

Returns income, expense, and net balance totals. Accepts all shared query params.

### Examples
```
GET /api/dashboard/summary
GET /api/dashboard/summary?from=2026-01-01&to=2026-03-31
GET /api/dashboard/summary?type=expense
```

### Response `200`
```json
{
    "data": {
        "total_income": 17150.00,
        "total_expenses": 6080.00,
        "net_balance": 11070.00
    }
}
```

---

## GET `/dashboard/categories` — category breakdown

Returns totals grouped by category and type, ordered by total descending. Accepts all shared query params.

### Examples
```
GET /api/dashboard/categories
GET /api/dashboard/categories?type=expense
GET /api/dashboard/categories?from=2026-01-01&to=2026-03-31
GET /api/dashboard/categories?type=expense&from=2026-01-01
```

### Response `200`
```json
{
    "data": [
        { "category": "salary",    "type": "income",  "total": 12850.00 },
        { "category": "rent",      "type": "expense", "total": 3350.00  },
        { "category": "freelance", "type": "income",  "total": 3200.00  },
        { "category": "groceries", "type": "expense", "total": 1050.00  }
    ]
}
```

---

## GET `/dashboard/trends` — monthly trends

Returns income vs expense totals per month for the last 12 months, ordered most recent first. Accepts all shared query params.

### Examples
```
GET /api/dashboard/trends
GET /api/dashboard/trends?from=2026-01-01
GET /api/dashboard/trends?category=salary&category=freelance
```

### Response `200`
```json
{
    "data": [
        { "month": "2026-04", "income": 4900.00,  "expense": 1520.00 },
        { "month": "2026-03", "income": 8250.00,  "expense": 2850.00 },
        { "month": "2026-02", "income": 5000.00,  "expense": 1750.00 }
    ]
}
```