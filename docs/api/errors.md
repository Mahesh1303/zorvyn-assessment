# Errors

All error responses follow this consistent shape:
```json
{
    "error": "description of what went wrong"
}
```

---

## Status Codes

| Code | Name | When it happens |
|---|---|---|
| `400` | Bad Request | missing required fields, invalid type value, invalid date format, invalid UUID |
| `401` | Unauthorized | missing `Authorization` header or expired/invalid JWT token |
| `403` | Forbidden | authenticated but role does not have permission for this action |
| `404` | Not Found | resource does not exist or has been soft deleted |
| `409` | Conflict | resource already exists e.g. duplicate email |
| `500` | Internal Server Error | unexpected server-side failure |

---

## Examples

### 400 — missing field
```json
{ "error": "name, email and password are required" }
```

### 400 — invalid type
```json
{ "error": "type must be income or expense" }
```

### 400 — invalid date
```json
{ "error": "date must be YYYY-MM-DD" }
```

### 401 — no token
```json
{ "error": "missing or invalid token" }
```

### 403 — wrong role
```json
{ "error": "forbidden: insufficient permissions" }
```

### 404 — not found
```json
{ "error": "user not found" }
```

### 409 — already exists
```json
{ "error": "user already exists" }
```

### 500 — server error
```json
{ "error": "internal server error" }
```

---

## Rate Limiting

All endpoints are rate limited at **20 requests per 30 seconds** using a sliding window algorithm. When the limit is exceeded the server returns:
```
HTTP 429 Too Many Requests
```