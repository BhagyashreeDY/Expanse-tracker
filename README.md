# Debt Optimization & Settlement Engine

## What is this?
This is a straightforward, production-ready Go API for tracking group expenses and settling debts. The main goal is to make sure everyone gets paid back with the fewest number of bank transfers possible. We built this to be exact—using proper decimal math instead of risky floats—and efficient.

### System Architecture

![System Architecture](https://mermaid.ink/img/pako:eNptkctuwjAQRX-FclZIVfEByC6K-gE7pC4QE0_iNo4jXyeqivjvNXm0EkSOnXvujK9H7EwNFSYF89mU4M4LUKDqA7fK9i0N94E6r56_P7D66xPrF66v2nNlE0HkE6m3wLALf7zNOnf06F_Vq29VlyvO2XREUe-fUKDToOf9T-iCHu-p66r6E_Xo7F7S9VWVz6p7_m-Vf99mH12zXlFqSDApjF9NCe68AAVqPnCrbN-T4D5S58Xz9wdWf31m_dL1VXutLCOIPJI2Bwy78MfbrHMHj_5dvfoWdbni6E0H7PX-CQVaDXre_4Qu6fGeuq6qP1KPzu4lXV-W-ay6Z_5Xmf-8zZ66Zp2h1JBgcjD8XkzwwAtQoOYDt8L2jYf7SJ3n7h9fWf_1m_mBq6u-m_D0v060lDk?type=png)

*(Note: If the image above is not visible, it is likely due to an internet connection issue. The diagram illustrates the data flow from the Client through the Gin framework, Service Layer, and Greedy Engine, down to the PostgreSQL database.)*

### Quick Links
- [How the matching works](ALGORITHM.md)
- [Design choices & Tech stack](DESIGN.md)
- [Example scenarios](EXAMPLES.md)
- [How to test the API](API_TESTING.md)
- [Project development notes (AI Prompts)](prompts.md)

---

## Key Features
- **Smart Debt Matching**: Uses a "greedy" approach to settle group debts in N-1 transactions or less.
- **Accuracy First**: We use `shopspring/decimal` for every calculation. No rounding errors, no missing cents.
- **Flexible Filters**: You can filter balances and settlements by date (using `from` and `to` query params).
- **Strategy Benchmarking**: Check how much better the optimized math is compared to basic pairwise settlement.
- **Clean Code**: Standard Go project structure with clear separation between routing, logic, and database.

---

## API At a Glance

| Method | Endpoint | What it does |
| :--- | :--- | :--- |
| `POST` | `/users` | Create a new user. |
| `POST` | `/groups` | Create an expense group. |
| `POST` | `/groups/:id/members` | Add a user to a group. |
| `POST` | `/groups/:id/expenses` | Add a bill (auto-split supported). |
| `GET` | `/groups/:id/balances` | See who is in the red or black. |
| `GET` | `/groups/:id/settlement` | Get the payment plan. |
| `GET` | `/groups/:id/settlement/compare` | Compare matching strategies. |
| `GET` | `/health` | Check if the API and DB are alive. |

---

## Getting Started

### 1. Setup
- Make sure you have **Go 1.22+** and **PostgreSQL** installed.
- Set up your `.env` file with your database credentials (check `config/config.go` for the keys).

### 2. Run Migrations
Import the schema into your Postgres database:
```bash
psql -d your_db_name -f migrations/001_init_schema.sql
```

### 3. Start the Engine
```bash
go run cmd/main.go
```

### 4. Try it out
Here is how you'd add a ₹120 dinner split between three people:
```bash
curl -X POST http://localhost:8080/groups/<GROUP_ID>/expenses \
-H "Content-Type: application/json" \
-d '{
    "payer_id": "<PAYER_UUID>",
    "amount": "120.00",
    "description": "Team Dinner",
    "split_type": "EQUAL",
    "splits": [
        {"user_id": "<USER_1>"},
        {"user_id": "<USER_2>"},
        {"user_id": "<USER_3>"}
    ]
}'
```

---

## Performance Notes
- **Time Complexity**: O(n log n). The bottleneck is just sorting the people by how much they owe or are owed.
- **Consistency**: We use database transactions to make sure that if a split fails, the whole expense isn't saved.
