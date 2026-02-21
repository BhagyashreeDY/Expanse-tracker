# Design & Tech Choices

This document explains why we built the system this way and the trade-offs we made.

## 1. Clean Architecture
We followed a layered approach to keep the code easy to maintain and test:
- **Handlers**: The "front door" (Gin routes). They just parse JSON and send it down.
- **Services**: The "brain". This is where rounding drift is handled and group logic lives.
- **Algorithms**: Pure math. These don't know about databases; they just take numbers and return transaction lists.
- **Repositories**: The "memory". All the SQL queries live here.

## 2. Why we avoid "Float" for money
In programming, basic decimals (floats) are stored as binary fractions. This causes tiny errors (like `0.1 + 0.2` resulting in `0.30000000000000004`). In a debt app, these errors accumulate and eventually, the group doesn't balance to zero.

**Our Choice**: We use `shopspring/decimal`. It treats numbers as integers with a decimal point, meaning `0.1 + 0.2` is exactly `0.3`. No rounding errors, ever.

## 3. Handling the "Penny Problem"
When you split a $10 bill 3 ways, it becomes $3.3333...
If you give everyone $3.33, you've only accounted for $9.99. Someone "lost" a penny.

**Our Solution**: Our `CalculateEqualSplits` function gives the first (N-1) people the rounded amount and gives the last person the remainder.
- Person 1: $3.33
- Person 2: $3.33
- Person 3: $3.34
Total: $10.00. This ensures the database always stays in balance.

## 4. Database Integrity
We use PostgreSQL because of its strong consistency and we use `pgxpool` for connection management.
- Every expense insertion is wrapped in a **SQL Transaction**. If the expense split data fails to save, the main expense isn't saved either. This avoids "zombie" expenses with no splits.
- We use `DECIMAL(18,2)` in the DB to match our Go-side decimal logic.

## 5. Filtering logic
The balance calculation is dynamic. Instead of storing a "running total" for each user (which can get out of sync), we calculate the net balance on the fly from the raw expense records. This allows us to easily add **Date Filtering**—you can ask "what do I owe for only the trip in June?" and the engine will calculate it perfectly.
