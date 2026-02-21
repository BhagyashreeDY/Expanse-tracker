# Common Scenarios & Use Cases

Here are a few ways the optimization engine handles real-world debt.

## Scenario 1: The Group Dinner (3 People)
**Alice, Bob, and Charlie** go out.
1. Alice pays ₹120 for the meal (₹40 each split).
2. Bob pays ₹90 for wine (₹30 each split).
3. Charlie pays ₹30 for desserts (₹10 each split).

### The Math:
- Alice spent ₹120, owes ₹80 (her share of the wine & dessert). Net: **+₹40**
- Bob spent ₹90, owes ₹80 (his share of dinner & dessert). Net: **+₹10**
- Charlie spent ₹30, owes ₹80 (his share of dinner & wine). Net: **-₹50**

### The Solution:
Instead of 6 different payments between everyone, the engine gives you just 2:
- **Charlie pays Alice ₹40**
- **Charlie pays Bob ₹10**

Total transactions reduced by **66%**.

---

## Scenario 2: One High Spender (5 People)
Alice pays for everything on a weekend trip (₹500 total).
David also pays ₹50 just for Alice and Bob (₹25 each).

### The Optimization:
Instead of Everyone paying Alice back separately, and Alice paying David back:
1. **Bob** pays **Alice** ₹125.00
2. **Charlie** pays **Alice** ₹100.00
3. **Eve** pays **Alice** ₹100.00
4. **David** pays **Alice** ₹50.00

Wait, why did Bob pay ₹125? Because the engine looks at the *net* debt, it cancels out the money Alice owed David and simplifies it into fewer, larger payments.

---

## Comparison Output Example
When you call our comparison endpoint, you get a report like this:

```json
{
  "greedy": {
    "transaction_count": 2,
    "total_volume": "50.00",
    "optimization_gain": "66.7%"
  },
  "baseline": {
    "transaction_count": 6,
    "total_volume": "50.00"
  }
}
```
This tells you exactly how many manual transfers you saved by using the optimization engine.
