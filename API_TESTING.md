# API Testing Checklist

Use these commands to verify the Expense Tracker API functionality.

## 1. Health Check
```bash
curl -X GET http://localhost:8080/health
```
**Expected Response:**
```json
{
  "database": "connected",
  "status": "ok",
  "time": "2026-02-21T00:05:00Z"
}
```

## 2. Create Users
```bash
# Create Alice
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"username": "Alice", "email": "alice@example.com"}'

# Create Bob
curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{"username": "Bob", "email": "bob@example.com"}'
```

## 3. Create Group
```bash
curl -X POST http://localhost:8080/groups \
-H "Content-Type: application/json" \
-d '{"name": "Trip to Paris"}'
```
*Note: Capture the `id` from the response.*

## 4. Add Members to Group
```bash
# Add Alice to Group
curl -X POST http://localhost:8080/groups/<GROUP_ID>/members \
-H "Content-Type: application/json" \
-d '{"user_id": "<ALICE_ID>"}'

# Add Bob to Group
curl -X POST http://localhost:8080/groups/<GROUP_ID>/members \
-H "Content-Type: application/json" \
-d '{"user_id": "<BOB_ID>"}'
```

## 5. Record Expense
```bash
# Alice pays $100, split equally with Bob
curl -X POST http://localhost:8080/groups/<GROUP_ID>/expenses \
-H "Content-Type: application/json" \
-d '{
    "payer_id": "<ALICE_ID>",
    "amount": "100.00",
    "description": "Dinner",
    "split_type": "EQUAL",
    "splits": [
        {"user_id": "<ALICE_ID>"},
        {"user_id": "<BOB_ID>"}
    ]
}'
```

## 6. Get Balances
```bash
curl -X GET http://localhost:8080/groups/<GROUP_ID>/balances
```
**Expected:** Alice: +50, Bob: -50

## 7. Get Optimized Settlement
```bash
curl -X GET http://localhost:8080/groups/<GROUP_ID>/settlement
```

## 8. Compare Strategies
```bash
curl -X GET http://localhost:8080/groups/<GROUP_ID>/settlement/compare
```

## 9. Negative/Invalid Tests (Expect 400 Errors)
```bash
# Duplicate split member
curl -X POST http://localhost:8080/groups/<GROUP_ID>/expenses ... -d '{"splits": [{"user_id": "U1"}, {"user_id": "U1"}]}'

# Mismatched sum
curl -X POST http://localhost:8080/groups/<GROUP_ID>/expenses ... -d '{"amount": "100", "splits": [{"user_id": "U1", "amount": "10"}]}'
```
