# Development History: AI Prompt Engineering

This document provides a transparent log of the logical prompts used to build the **Debt Optimization & Intelligent Settlement Engine**. The development followed a rigorous, step-by-step engineering approach, focusing on architectural integrity, financial precision, and algorithmic optimization.

## Step 1: Core Architecture & Foundation
> "Initialize a Go project for a production-grade Debt Optimization Engine using Clean Architecture. The project should be structured with `/cmd` for the entry point and `/internal` for handlers, services, algorithms, and repositories. Use the Gin framework for routing and `pgxpool` for PostgreSQL connectivity. 
> 
> **Critical Requirement**: All monetary logic must use the `shopspring/decimal` library to ensure fixed-point financial precision; the use of `float64` is strictly prohibited to prevent data corruption. Implement Dependency Injection to wire the layers together during startup."

## Step 2: Core Domain Logic (The Greedy Engine)
> "Implement a deterministic **Greedy Debt Minimization** algorithm within the `/internal/algorithms` package.
> 
> **Logic Requirements**: 
> 1. Input: A mapping of net balances across a group.
> 2. Segregation: Separate users into 'debtors' (negative) and 'creditors' (positive) pools.
> 3. Optimization: Perform $O(N \log N)$ matching by sorting both pools by magnitude and iteratively settling the largest debtor against the largest creditor.
> 4. Goal: Solve the 'Minimum Cash Flow' problem by reducing the total transaction count to at most $N-1$. 
> Include unit tests for circular debt scenarios and zero-balance edge cases."

## Step 3: Service Layer & Financial Integrity
> "Develop the service layer to handle internal expense logic and validation. 
> 
> **Feature Set**:
> 1. **Rounding Drift Protection**: In the `CalculateEqualSplits` function, ensure that the last participant absorbs the remainder (`Total - Sum`) so the final group total is 100% accurate down to the ₹0.01.
> 2. **Validation Pipeline**: Build a `ValidateSplits` function to reject expenses with negative amounts, duplicate participant IDs, or mismatched sum totals.
> 3. **Transactional Safety**: Ensure all repository operations for creating expenses and splits are wrapped in a single PostgreSQL transaction to guarantee atomicity."

## Step 4: Performance Analytics & Benchmarking
> "Extend the engine to support strategy comparison and historical queries:
> 1. **Baseline Strategy**: Implement a 'Naive' settlement strategy (pairwise matching) to serve as a baseline for performance metrics.
> 2. **Comparison Endpoint**: Build a `/settlement/compare` endpoint that returns transaction count, total volume, and an 'Optimization Gain' percentage comparing Naive vs. Greedy strategies.
> 3. **Time-Series Filtering**: Update all repository methods to support optional `from` and `to` timestamps, allowing users to generate settlements for specific date ranges."

## Step 5: Production Hardening & Audit
> "Conduct a final production audit and hardening phase:
> 1. **Middleware**: Implement a professional panic recovery middleware in Gin that returns a structured JSON 500 error instead of a server stack trace.
> 2. **Health Monitoring**: Add a `/health` endpoint that performs a real-time database ping to report system status.
> 3. **Documentation**: Refactor the project documentation to explain the $O(N \log N)$ complexity of the greedy engine, the 'Penny Problem' in splitting logic, and the engineering rationale for avoiding floating-point math in financial systems."
