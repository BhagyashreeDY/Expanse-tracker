# System Design: Debt Optimization Engine

## 1. Architectural Layers
The project follows a **Clean Layered Architecture** to ensure modularity, testability, and clear separation of concerns.

- **`cmd/` (Entry Point)**: Handles dependency injection, configuration loading, and starts the Gin server.
- **`internal/handlers/` (Interface Layer)**: Manages HTTP request parsing, input validation (Gin bindings), and response formatting.
- **`internal/services/` (Business Logic)**: Orchestrates complex operations like calculating balances across groups, handling expense split logic, and enforcing financial integrity rules.
- **`internal/algorithms/` (Logic Layer)**: Pure, stateless implementation of settlement algorithms (Greedy and Naive).
- **`internal/repositories/` (Data Access)**: Handles all PostgreSQL interactions using the high-performance `pgxpool` driver.

## 2. Database Schema
The schema is optimized for **atomic financial transactions** and data integrity.

- **`users` & `groups`**: Core domain entities.
- **`expenses`**: Records the summary of a spending event (total amount, payer, split type).
- **`expense_splits`**: Normalizes the participant data. Each row represents exactly how much a specific user owes for a specific expense.
- **`DECIMAL(18,2)`**: All monetary columns use fixed-point decimals. This prevents precision loss and ensures that the sum of splits always matches the database state of the parent expense.

## 3. Financial Precision: The Decimal Strategy
Standard floating-point numbers (`float64`) are avoided entirely in the financial logic.

- **The Problem**: Float64 uses binary fractions (base-2), which cannot represent common decimal values (base-10) like 0.1 exactly.
- **The Solution**: We use the `shopspring/decimal` library in Go. It stores numbers as arbitrary-precision integers with an exponent, allowing for exact addition, subtraction, and comparison.

## 4. Validation & Integrity Strategy
The system implements a multi-stage validation pipeline:
1. **Request Validation**: Gin ensures that UUIDs and positive amounts are sent in the JSON body.
2. **Business Rules (`ValidateSplits`)**:
   - Rejects negative split amounts.
   - Rejects duplicate users in a single expense.
   - **Rounding Drift Guard**: In `CalculateEqualSplits`, the last participant is assigned `Total - (N-1)*Split` to ensure the sum is exactly equal to the total, even when the division is not clean (e.g., $10 / 3$).
3. **Database Constraints**: SQL-level `CHECK` constraints prevent negative amounts from being stored.

## 5. Time-Filtered Settlement Logic
The repository layer supports optional `from` and `to` timestamps for all financial queries.
- This allows users to view their debt status for a specific month, trip, or custom range.
- The service layer recalculates balances dynamically based on the filtered set of expenses, providing a point-in-time financial snapshot without needing separate ledger entries.
