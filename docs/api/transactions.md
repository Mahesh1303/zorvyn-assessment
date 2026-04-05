# Transactions

Base URL: `http://localhost:8080/api`

All endpoints require:
```
Authorization: Bearer <token>
```

> Read access: all roles
> Write access: admin + analyst only

---

## POST `/transactions` — create transaction

### Request body
```json
{
    "amount": 5000.00,
    "type": "income",
    "category": "salary",
    "description": "Monthly salary",
    "date": "2026-04-01"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `amount` | number | yes | must be greater than 0 |
| `type` | string | yes | `income` or `expense` |
| `category` | string | yes | free text e.g. salary, rent, groceries |
| `description` | string | no | optional notes |
| `date` | string | yes | format `YYYY-MM-DD` |

### Response `201`
```json
{
    "message": "transaction created"
}
```

### Errors

| Status | Reason |
|---|---|
| `400` | missing fields, invalid type, or invalid date format |
| `403` | viewer role cannot create records |

---

## GET `/transactions` — list transactions

Returns all transactions. Supports optional filters via query params.

### Query params

| Param | Example | Description |
|---|---|---|
| `type` | `expense` | filter by `income` or `expense` |
| `category` | `salary` | filter by category name |
| `from` | `2026-01-01` | filter from date (inclusive) |
| `to` | `2026-03-31` | filter to date (inclusive) |

All params are optional and combinable:
```
GET /api/transactions
GET /api/transactions?type=expense
GET /api/transactions?category=salary
GET /api/transactions?from=2026-01-01&to=2026-03-31
GET /api/transactions?type=expense&category=rent&from=2026-01-01&to=2026-03-31
GET /api/transactions?limit=10&offset=0
```

### Response `200`
```json
{
    "data": [
        {
            "id": "uuid",
            "created_by": "uuid",
            "amount": 1200,
            "type": "expense",
            "category": "rent",
            "description": "Apartment rent",
            "date": "2026-03-01T00:00:00Z",
            "created_at": "2026-03-01T10:00:00Z"
        }
    ],
    "meta": {
        "limit": 10,
        "offset": 0
    }
}
```

---

## GET `/transactions/:id` — get transaction

### Path params

| Param | Description |
|---|---|
| `id` | transaction UUID |

### Response `200`
```json
{
    "id": "uuid",
    "created_by": "uuid",
    "amount": 5000,
    "type": "income",
    "category": "salary",
    "description": "Monthly salary",
    "date": "2026-04-01T00:00:00Z",
    "created_at": "2026-04-05T10:00:00Z"
}
```

### Errors

| Status | Reason |
|---|---|
| `404` | transaction not found |

---

## PUT `/transactions/:id` — update transaction

### Request body

Same fields as POST. All fields required.
```json
{
    "amount": 5500.00,
    "type": "income",
    "category": "salary",
    "description": "Monthly salary — updated",
    "date": "2026-04-01"
}
```

### Response `200`
```json
{
    "data": {
        "id": "uuid",
        "created_by": "uuid",
        "updated_by": "uuid",
        "amount": 5500,
        "type": "income",
        "category": "salary",
        "description": "Monthly salary — updated",
        "date": "2026-04-01T00:00:00Z",
        "updated_at": "2026-04-05T11:00:00Z"
    }
}
```

### Errors

| Status | Reason |
|---|---|
| `400` | invalid fields |
| `403` | viewer role cannot update records |
| `404` | transaction not found |

---

## DELETE `/transactions/:id` — soft delete

Marks the transaction as deleted. Data is retained in the database and never permanently removed.

### Response `200`
```json
{
    "message": "deleted"
}
```

### Errors

| Status | Reason |
|---|---|
| `403` | viewer role cannot delete records |
| `404` | transaction not found |