# Debt Optimization & Intelligent Settlement Engine

## Project Overview
A production-grade REST API built in Go for tracking shared expenses and optimizing debt settlements. The engine ensures high financial accuracy using fixed-point arithmetic and minimizes settlement overhead using an optimized greedy algorithm.

### Key References
- **[System Design](DESIGN.md)**: Deep dive into architecture and precision handling.
- **[Algorithm Logic](ALGORITHM.md)**: Step-by-step breakdown of the Greedy settlement engine.
- **[Real-World Examples](EXAMPLES.md)**: Scenarios showing the algorithm in action.
- **[Development Prompts](prompts.md)**: Transparent list of all AI interactions.
- **[Testing Guide](API_TESTING.md)**: Manual verification commands and payloads.

---

## 🚀 Core Features
- **Deterministic Settlement**: Minimizes transaction count using an $O(N \log N)$ greedy algorithm.
- **Strategy Benchmarking**: Real-time comparison between Optimized and Naive settlement strategies.
- **Fixed-Point Precision**: Strictly uses the `decimal` type to prevent binary rounding errors found in `float64`.
- **Time-Series Filtering**: Snapshot your balances and debts across any specific date range.
- **Rounding Drift Protection**: Guaranteed sum consistency for equal and percentage splits.

---

## 📡 API Endpoints Summary

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/users` | Register a new user with email and username. |
| `POST` | `/groups` | Create a new expense group. |
| `POST` | `/groups/:id/members` | Add a user to an existing group. |
| `POST` | `/groups/:id/expenses` | Record a new expense (supports EQUAL splits). |
| `GET` | `/groups/:id/balances` | View net balances for all group members. |
| `GET` | `/groups/:id/settlement` | Get the optimized settlement transaction list. |
| `GET` | `/groups/:id/settlement/compare`| Compare Optimized vs. Naive metrics. |
| `GET` | `/health` | Check API and Database connectivity status. |

---

## 🛠️ Quick Start & Usage Snippet

### 1. Requirements
- Go 1.22+
- PostgreSQL
- `.env` file with `DB_USER`, `DB_PASSWORD`, `DB_NAME`

### 2. Installation
```bash
# Clone the repository
git clone https://github.com/BhagyashreeDY/Expanse-tracker.git
cd Expanse-tracker

# Run migrations (PostgreSQL)
psql -d expense_tracker -f migrations/001_init_schema.sql

# Start the server
go run cmd/main.go
```

### 3. Example Request (Record Expense)
```bash
curl -X POST http://localhost:8080/groups/<GROUP_ID>/expenses \
-H "Content-Type: application/json" \
-d '{
    "payer_id": "<USER_UUID>",
    "amount": "120.00",
    "description": "Shared Dinner",
    "split_type": "EQUAL",
    "splits": [
        {"user_id": "<USER_1>"},
        {"user_id": "<USER_2>"},
        {"user_id": "<USER_3>"}
    ]
}'
```

---

## 📊 Algorithmic Complexity
- **Time**: $O(N \log N)$ where $N$ is the number of participants (due to sorting creditors and debtors).
- **Space**: $O(N)$ for balance aggregation and matching pools.

---

## 🏗 Architecture
The project follows **Clean Architecture** patterns, ensuring that business logic in `internal/services` is decoupled from HTTP handlers in `internal/handlers` and database persistence in `internal/repositories`.
