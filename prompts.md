# AI Prompts Used (Technical Details & Core Logic)

This document provides complete transparency into the AI-assisted development process. The following sequence of prompts reflects the evolution of the **Debt Optimization & Intelligent Settlement Engine**.

## 1. Initial Architecture & Foundation
"Design a production-grade Go REST API for a debt-settlement engine using Clean Architecture principles. The project must exclude float64 for financial data and use 'shopspring/decimal' instead. Organize the codebase into correctly separated layers (Handlers, Services, Algorithms, Repositories). Provide a skeleton that ensures complete separation of concerns and dependency injection via pgxpool for PostgreSQL."

## 2. Algorithmic Implementation (Greedy Logic)
"Implement a Greedy Debt-Minimization algorithm in Go. 
**The Logic:**
1. Input: A map of `username` to `net_balance` (Decimal).
2. Separate into `debtors` (negative) and `creditors` (positive).
3. Sort both by magnitude (Greedy).
4. Match largest debtor with largest creditor iteratively until balances are zero.
5. Return the list of transactions.
Ensure all calculations use `shopspring/decimal` for financial accuracy."

## 3. Financial Integrity & Rounding Drift
"Create a service layer function `CalculateEqualSplits` that handles rounding drift (e.g., $100 split 3 ways). The last person must take the remainder (`Total - Sum`) to ensure the total is exactly correct. Implement `ValidateSplits` to reject expenses where the sum doesn't match the total, or where there are duplicate members or negative amounts."

## 4. Feature Enhancement (Analytics & Modeling)
"Extend the engine with advanced modeling:
1. **Strategy Comparison**: Implement an 'Optimized' vs 'Naive' comparison endpoint showing transaction count, volume, and percentage gain.
2. **Time-Series Support**: Add `from` and `to` filtering to all balance and settlement queries to allow historical snapshots."

## 5. Production Audit & Stabilization
"Perform a full backend audit:
1. **Sanitize Routes**: Check all paths for 404 risks and ensure method consistency.
2. **Middleware**: Add structured panic recovery and structured logging.
3. **Health Check**: Add a `/health` endpoint that pings the DB and reports status.
4. **Documentation**: Generate a README explaining the 'Why float is unsafe' rationale and the 'Greedy Algorithm' complexity (O(N log N))."

## 6. Documentation Overhaul
"Refactor all Markdown documentation (README, DESIGN, ALGORITHM, EXAMPLES) for academic and industry-level evaluation. Ensure clear structural links, pseudocode for algorithms, realistic usage scenarios, and a clean professional layout. Focus on explaining the design trade-offs and complexity analysis."
