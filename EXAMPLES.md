# Settlement Examples: Before vs. After Optimization

This document provides concrete scenarios to demonstrate the efficiency gains of the Greedy Settlement Engine compared to naive pairwise matching.

## Scenario 1: The Weekend Trip (Circular Debt)
**Participants**: Alice, Bob, Charlie

### Raw Activity:
1. **Dinner**: Alice pays $120 for everyone ($40 each).
2. **Fuel**: Bob pays $90 for everyone ($30 each).
3. **Parking**: Charlie pays $30 for everyone ($10 each).

### Net Balances:
- **Alice**: Pays $120, owes $40 (Own) + $30 (Fuel) + $10 (Parking). Net: **+$40**
- **Bob**: Pays $90, owes $40 (Dinner) + $30 (Own) + $10 (Parking). Net: **+$10**
- **Charlie**: Pays $30, owes $40 (Dinner) + $30 (Fuel) + $10 (Own). Net: **-$50**

### Comparison:
| Strategy | Transactions | Total Volume | Efficiency |
| :--- | :--- | :--- | :--- |
| **Naive** | 6 individual payments | $160 | - |
| **Optimized** | 2 payments | $50 | **66% Fewer Tx** |

**Optimized Result**: 
- Charlie pays Alice $40
- Charlie pays Bob $10

---

## Scenario 2: Large Group Outing
**Participants**: Alice, Bob, Charlie, David, Eve

### Raw Activity:
- Alice pays $500 for the group ($100 each).
- David pays $50 for Alice and Bob ($25 each).

### Net Balances:
- **Alice**: Pays $500, owes $100 (Own) + $25 (David). Net: **+$375**
- **Bob**: Pays $0, owes $100 (Alice) + $25 (David). Net: **-$125**
- **Charlie**: Pays $0, owes $100 (Alice). Net: **-$100**
- **David**: Pays $50, owes $100 (Alice). Net: **-$50**
- **Eve**: Pays $0, owes $100 (Alice). Net: **-$100**

### Optimized Settlement result:
1. **Bob** pays **Alice** $125.00
2. **Charlie** pays **Alice** $100.00
3. **Eve** pays **Alice** $100.00
4. **David** pays **Alice** $50.00

---

## Scenario 3: Real-World Comparison Output
When you query the `/settlement/compare` endpoint, the system returns a payload similar to this:

```json
{
  "greedy": {
    "transaction_count": 4,
    "total_volume": "375.00",
    "optimization_gain": "50.0%"
  },
  "baseline": {
    "transaction_count": 8,
    "total_volume": "375.00"
  }
}
```

### Key Observations:
1. **Transaction Count**: The Greedy algorithm consistently produces fewer transactions than the Baseline.
2. **Volume Consistency**: The total money moved stays identical between both strategies, ensuring no funds are lost or created during optimization.
3. **Gain Percentage**: Calculated as `(BaselineCount - GreedyCount) / BaselineCount`.
