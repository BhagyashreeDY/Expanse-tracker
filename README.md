# Debt Optimization & Intelligent Settlement Engine

## Project Overview
A production-grade REST API built in Go for tracking shared expenses and optimizing debt settlements. The engine ensures clinical financial accuracy using fixed-point arithmetic and minimizes settlement overhead using an optimized greedy algorithm.

## Core Features
- **Deterministic Settlement**: Minimizes total transaction count using an O(n log n) greedy algorithm.
- **Strategy Comparison**: Built-in benchmarking to compare Optimized (Greedy) vs Naive (Pairwise) settlement strategies.
- **Financial Precision**: Uses `shopspring/decimal` for all monetary logic, strictly avoiding `float64` to prevent precision loss.
- **Time-Filtered Analytics**: Support for `from` and `to` date filters on all balance and settlement queries.
- **Integrity Validation**: Strict rules for split consistency, non-negative amounts, and rounding drift protection.

## Why `float64` is Unsafe for Financial Systems
Standard floating-point numbers (`float64`) use binary representation (IEEE 754), which cannot accurately represent many decimal values (like $0.1$ or $0.01$). Over thousands of transactions, these tiny inaccuracies accumulate into significant "rounding errors," leading to data corruption and financial discrepancies.

**Solution**: This project uses `DECIMAL(18,2)` in PostgreSQL and `shopspring/decimal` in Go to handle money as fixed-precision integers internally.

## Algorithmic Strength
### Greedy Settlement Logic
The core engine solves the "Minimum Cash Flow" problem by:
1. **Aggregating Net Balances**: Computing exactly how much each user is "in the red" or "in the black" across all group expenses.
2. **Sorting by Magnitude**: Sorting debtors and creditors by their absolute values.
3. **Optimized Matching**: Greedily matching the largest debtor with the largest creditor iteratively.

### Why Greedy Reduces Transactions?
In a naive system, every expense might result in a separate transaction. In a group of $N$ people, if everyone pays for everyone else once, you could have $N(N-1)$ transactions. 

The Greedy algorithm proves that any group of $N$ people can be settled in at most $N-1$ transactions, drastically reducing the mental and financial overhead of bank transfers.

### Complexity Analysis
- **Time Complexity**: $O(N \log N)$ where $N$ is the number of users in the group. The dominant factor is the sorting of the balances.
- **Space Complexity**: $O(N)$ to store the balance mapping and the final transaction list.

## Before vs After Optimization
**Scenario**: Alice, Bob, Charlie, and David go on a trip.

**Raw Expenses (Pre-Optimization):**
- Alice pays $100 (others owe her $25 each)
- Bob pays $50 (others owe him $12.5 each)
- *Total raw splits: 6 transactions if everyone pays everyone back individually.*

**Optimized Result (Greedy):**
- Charlie pays Alice $37.50
- David pays Alice $37.50
- Alice pays Bob $0 (Net calculation cancels this out as Alice is still a net creditor)
- *Total optimized transactions: 2.*

**Efficiency Gain**: 66.6% reduction in manual payments.

## Design Trade-offs
- **Greedy vs Optimal**: While the greedy approach is $O(N \log N)$ and highly efficient, it doesn't always find the *theoretical* minimum if sub-groups within the group could balance themselves independently (a variation of the subset sum problem, which is NP-Hard). However, for group expenses, the difference is negligible, and the greedy approach is the industry standard for its performance and reliability.
- **Real-time Recalculation**: Every query for settlements recalculates current balances from raw expenses. This ensures 100% data integrity but might require caching for extremely large groups.

## API Endpoints
- `POST /users`: Register a new user.
- `POST /groups`: Create a new expense group.
- `POST /groups/:id/expenses`: Add an expense (EQUAL, PERCENTAGE, or EXACT split).
- `GET /groups/:id/balances`: View current net balances.
- `GET /groups/:id/settlement`: Get the optimized settlement plan.
- `GET /groups/:id/settlement/compare`: Compare Greedy vs Naive strategies.
- `POST /groups/:id/settlement/record`: Mark a payment as completed.

## Architecture
- `/cmd`: Main entry point.
- `/internal/handlers`: Controller layer (Gin).
- `/internal/services`: Domain services (Splitting logic, validation).
- `/internal/algorithms`: Pure algorithmic engine.
- `/internal/repositories`: Database persistence (Postgres/pgx).
