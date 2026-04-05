# Users

Base URL: `http://localhost:8080/api`

All endpoints require:
```
Authorization: Bearer <token>
```

> Admin role required for all user endpoints.

---

## POST `/users` — create user

Creates a new user with the specified role.

### Request body
```json
{
    "name": "John Doe",
    "email": "john@test.com",
    "password": "password123",
    "role": "analyst"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `name` | string | yes | full name |
| `email` | string | yes | unique email address |
| `password` | string | yes | plain text, hashed before storage |
| `role` | string | yes | one of `admin` `analyst` `viewer` |

### Response `201`
```json
{
    "message": "user created successfully"
}
```

### Errors

| Status | Reason |
|---|---|
| `400` | missing required fields or invalid role |
| `403` | caller is not an admin |
| `409` | email already exists |

---

## GET `/users` — list users

Returns all users. Passwords are never included in responses.

### Response `200`
```json
{
    "data": [
        {
            "id": "e17c1df1-3086-4190-85df-d193517c22ab",
            "name": "Admin User",
            "email": "admin@test.com",
            "role": "admin",
            "is_active": true,
            "created_at": "2026-04-01T10:00:00Z"
        },
        {
            "id": "a82f3bc4-1234-5678-abcd-ef0123456789",
            "name": "Analyst User",
            "email": "analyst@test.com",
            "role": "analyst",
            "is_active": true,
            "created_at": "2026-04-01T10:00:00Z"
        }
    ]
}
```

---

## GET `/users/:id` — get user

Returns a single user by ID.

### Path params

| Param | Description |
|---|---|
| `id` | user UUID |

### Response `200`
```json
{
    "id": "e17c1df1-3086-4190-85df-d193517c22ab",
    "name": "Admin User",
    "email": "admin@test.com",
    "role": "admin",
    "is_active": true,
    "created_at": "2026-04-01T10:00:00Z"
}
```

### Errors

| Status | Reason |
|---|---|
| `404` | user not found |

---

## PATCH `/users/:id/role` — update role

Changes a user's role.

### Request body
```json
{
    "role": "viewer"
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `role` | string | yes | one of `admin` `analyst` `viewer` |

### Response `200`
```json
{
    "message": "role updated"
}
```

### Errors

| Status | Reason |
|---|---|
| `400` | invalid role value |
| `404` | user not found |

---

## PATCH `/users/:id/status` — update status

Activates or deactivates a user. Inactive users cannot log in.

### Request body
```json
{
    "active": false
}
```

| Field | Type | Required | Description |
|---|---|---|---|
| `active` | boolean | yes | `true` to activate, `false` to deactivate |

### Response `200`
```json
{
    "message": "user status updated"
}
```

### Errors

| Status | Reason |
|---|---|
| `404` | user not found |