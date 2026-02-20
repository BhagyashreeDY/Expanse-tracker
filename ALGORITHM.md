# Settlement Algorithm: Optimized Greedy Approach

The Debt Optimization Engine uses a deterministic **Greedy Matching Algorithm** to solve the "Minimum Cash Flow" problem. This document provides a deep dive into the logic, pseudocode, and mathematical properties of the implementation.

## 1. Problem Definition
In a group of $N$ people with shared expenses, individuals often owe varying amounts to multiple others. A naive approach (paying back every individual split) can lead to $N(N-1)$ transactions.

**Goal**: Settle all outstanding debts using the minimum number of transactions possible, while ensuring 100% financial accuracy (zero net balance for all).

## 2. Step-by-Step Logic

### Step 1: Calculate Net Balances
For every user in the group, we calculate a single `NetBalance`:
$$\text{NetBalance} = \sum \text{PaidAmounts} - \sum \text{OwedAmounts}$$
- **Positive Balance**: The user is a **Creditor** (owed money).
- **Negative Balance**: The user is a **Debtor** (owes money).

### Step 2: Segregate and Sort
Users are separated into two pools:
1. **Debtors**: Those with negative balances.
2. **Creditors**: Those with positive balances.

Both lists are sorted by their **magnitude** (absolute value) in descending order. This "Greedy" choice ensures we attempt to clear the largest debts first.

### Step 3: Iterative Matching
While both pools are non-empty:
1. Select the largest debtor ($D$) and the largest creditor ($C$).
2. Calculate the settlement amount: $S = \min(|D.\text{balance}|, C.\text{balance})$.
3. Create a transaction: **$D$ pays $C$ amount $S$**.
4. Update balances:
   - $D.\text{balance} = D.\text{balance} + S$
   - $C.\text{balance} = C.\text{balance} - S$
5. If a user's balance reaches exactly zero, remove them from their respective pool.

## 3. Pseudocode
```text
FUNCTION SettleOptimized(balances):
    debtors = []
    creditors = []
    
    FOR EACH user, amount IN balances:
        IF amount < 0:
            debtors.push({user, abs(amount)})
        ELSE IF amount > 0:
            creditors.push({user, amount})
            
    SORT debtors DESCENDING BY amount
    SORT creditors DESCENDING BY amount
    
    transactions = []
    
    WHILE debtors IS NOT EMPTY AND creditors IS NOT EMPTY:
        d = debtors[0]
        c = creditors[0]
        
        amount = MIN(d.amount, c.amount)
        transactions.push({From: d.user, To: c.user, Amount: amount})
        
        d.amount -= amount
        c.amount -= amount
        
        IF d.amount == 0: debtors.pop_front()
        IF c.amount == 0: creditors.pop_front()
        
    RETURN transactions
```

## 4. Why Greedy Reduces Transactions?
In a group of $N$ people, the maximum number of transactions required to settle all debts is $N-1$. 

In the real world, a single large creditor might be owed money by five different people. Instead of that creditor receiving five small payments, the greedy algorithm tries to match them with the largest debtors first, often closing out multiple users' balances in single steps.

## 5. Complexity Analysis
- **Time Complexity**: $O(N \log N)$ where $N$ is the number of participants. The dominant cost is the sorting of the balances.
- **Space Complexity**: $O(N)$ to store the segregated debtor/creditor lists and the resulting transaction list.

## 6. Mathematical Correctness
The algorithm is **deterministic** and **correct**. It guarantees:
1. **Convergence**: Every iteration reduces at least one person's balance to zero. Since there are $N$ people, it terminates in at most $N-1$ steps.
2. **Precision**: By using fixed-point decimal arithmetic, the sum of all payments will exactly match the total debt volume, leaving everyone with a perfect $0.00$ balance.
