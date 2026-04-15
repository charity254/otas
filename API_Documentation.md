# API Documentation

## Transactions

### Process a Transaction
Records a new transaction for the authenticated user and triggers the OTAS micro-saving engine. The engine will calculate the appropriate round-up/deduction and route the funds to either the `group`, `locked`, `flexible`, or `main` account based on the user's configured Saving Type.

**Endpoint:**
`POST /transactions`

**Headers:**
- `Content-Type: application/json`
- `Authorization: Bearer <JWT_TOKEN>` (Required)

**Request Body:**
```json
{
  "amount": 150.50
}
```

**Parameters:**
- `amount` (float, required): The total amount of the transaction. Must be strictly greater than 0.

**Success Response:**
`HTTP 201 Created`
```json
{
  "message": "transaction processed successfully",
  "transaction": {
    "id": 1,
    "user_id": 1,
    "amount": 150.5,
    "deduction": 15.05,
    "allocated_to": "group",
    "created_at": "2026-04-15T19:45:00Z"
  }
}
```

**Error Responses:**
- `HTTP 400 Bad Request` - If the request body is missing or `amount` is invalid.
- `HTTP 401 Unauthorized` - If the JWT token is missing, expired, or invalid.
- `HTTP 500 Internal Server Error` - If the database fails to process the transaction.
