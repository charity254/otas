# Testing Guide

This guide explains how to manually test the `POST /transactions` endpoint on your local development server before the authentication (`/login`, `/register`) endpoints are fully complete.

## Prerequisite: Database Setup 

The saving engine requires a valid user with defined saving parameters (`saving_type` and `daily_limit`), plus an initial `accounts` record.

If you don't have these, connect to your PostgreSQL database and run the following dummy data insertions:

```sql
-- 1. Create a dummy user
INSERT INTO users (name, email, phone, password, saving_type, daily_limit) 
VALUES ('Test User', 'test@example.com', '+254700000000', 'password123', 'group', 5) 
RETURNING id; 
-- Assume this query returns ID = 1.

-- 2. Create a main account for that user
INSERT INTO accounts (user_id, type, balance) 
VALUES (1, 'main', 1000.00);

-- OPTIONAL: Create a 'group' saving account for that user
INSERT INTO accounts (user_id, type, balance) 
VALUES (1, 'group', 0.00);
```

## How to Test Without JWT (Local Development Only)

If you just want to test the transaction logic without generating JWT Tokens, you can temporarily hardcode the `AuthMiddleware` in `internal/routes/routes.go`.

**1. Temporarily modify `routes.go`:**
Change `AuthMiddleware()` to return a hardcoded user ID:
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Set("userID", 1) // Force User 1 for testing
        c.Next()
    }
}
```

**2. Start the server:**
```bash
go run main.go
```

**3. Send the API request:**
Run this in another terminal, or use Postman:
```bash
curl -X POST http://localhost:8080/transactions \
-H "Content-Type: application/json" \
-d '{"amount": 150}'
```

*(Make sure to revert `routes.go` before committing!)*


## How to Test With JWT (Proper Setup)

If you want to test the full flow:

1. Ensure your `.env` file has a `JWT_SECRET` defined:
```env
JWT_SECRET=super_secret_key_123
```
2. Generate a token programmatically by calling `jwt.Generate(1)` (where 1 is your User ID).
3. Send the HTTP request with the `Authorization` header:

```bash
curl -X POST http://localhost:8080/transactions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer <YOUR_GENERATED_JWT_TOKEN>" \
-d '{"amount": 150}'
```
