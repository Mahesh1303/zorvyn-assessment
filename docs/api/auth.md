# Auth

Base URL: `http://localhost:8080`

No authentication required for these endpoints.

---

## POST `/auth/register-admin`

One-time endpoint to create the first admin account. Once an admin exists this endpoint permanently returns `403`. All subsequent users are created by the admin via `POST /api/users`.

### Request body
```json
{
    "name": "Admin User",
    "email": "admin@finance.com",
    "password": "password123"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | full name |
| `email` | string | yes | admin email address |
| `password` | string | yes | plain text, hashed before storage |

> Role is always set to `admin` — it cannot be changed via this endpoint.

### Response `201`
```json
{
    "message": "admin registered successfully",
    "data": {
        "id": "uuid",
        "name": "Admin User",
        "email": "admin@finance.com",
        "role": "admin",
        "is_active": true,
        "created_at": "2026-04-05T10:00:00Z"
    }
}
```

### Errors

| Status | Reason |
|---|---|
| `400` | missing name, email or password |
| `403` | an admin already exists — system already initialized |
| `409` | email already taken |

---

## POST `/auth/login`

Authenticates a user and returns a JWT token.

### Request body
```json
{
    "email": "admin@finance.com",
    "password": "password123"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `email` | string | yes | user email address |
| `password` | string | yes | user password |

### Response `200`
```json
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

Use this token in the `Authorization` header for all protected endpoints:
```
Authorization: Bearer <token>
```

Tokens expire after **24 hours**.

### Errors

| Status | Reason |
|---|---|
| `400` | missing email or password |
| `401` | invalid credentials |
| `403` | account is inactive |

---

## Typical first-time flow
```
1. POST /auth/register-admin   → create the first admin (one time only)
2. POST /auth/login            → get JWT token
3. POST /api/users             → admin creates analyst and viewer accounts
4. POST /auth/login            → analysts and viewers log in with their credentials
```