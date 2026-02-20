# System Design & Implementation Details

## Database Schema
The schema is designed for high data integrity:
- **`users` & `groups`**: Core entities.
- **`expenses`**: Stores the event, payer, and split type.
- **`expense_splits`**: Normalizes the relationship between an expense and its participants.
- **`DECIMAL(18,2)`**: Chosen for all currency columns to provide 18 digits of precision with 2 decimal places.

## Money Handling
### Why not float64?
Floating-point numbers (IEEE 754) represent decimals as binary fractions, which leads to precision loss (e.g., `0.1 + 0.2 != 0.3`). In financial systems, this results in "rounding drift" where a few cents disappear or appear unexpectedly.
**Solution**: We use `shopspring/decimal`, which performs exact decimal arithmetic.

### Rounding Strategy
- **Division**: When splitting $10 among 3 people, we get $3.33, $3.33, and $3.34. 
- **Drift Handling**: The service layer calculates the `total - (n-1)*split` for the final participant to ensure the sum always equals the original expense.

## Advanced Analytical Features
### Graph Modeling
The system maps debts into a **Weighted Directed Graph**. Nodes represent users, and directed edges represent the amount owed between individuals. This allows for pathfinding-based optimizations and better visualization of group financial health.

### Strategy Benchmarking
We compare our **Greedy Optimization** against **Random** and **FIFO** baselines.
- **Greedy**: Always reduces the transaction count to at most `N-1`.
- **Random/FIFO**: Baseline matching that often results in redundant transactions.
The API provides an "Optimization Gain" metric to prove the algorithm's efficiency in a production environment.

### Anomaly Detection Rule Engine
1. **Outlier Detection**: Flags any expense > 3x the group's average.
2. **Duplication Guard**: Identifies identical amounts and descriptions recorded in proximity.
3. **Behavioral Imbalance**: Detects "Perpetual Payers" where one user funds >80% of expenses, indicating potential tax or reimbursement risks.

### Transaction Cost Optimization
When a `transaction_fee` is specified, the engine calculates whether the debt value outweighs the transfer cost. It automatically prunes micro-transactions that would cost more in fees than the value they transfer (a constrained optimization problem).
