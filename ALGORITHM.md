# How we settle debts: The Greedy Matching Logic

We use a "Greedy" algorithm to figure out the best way to move money around. The goal is simple: everyone gets their money back, but with the fewest number of Venmo or bank transfers possible.

## The Problem
If 5 people go on a trip and everyone pays for something, you could end up with dozens of tiny transfers (Alice pays Bob for coffee, Bob pays Charlie for gas, etc.). That's annoying.

**Our Goal**: Clear all debts in **N-1 transactions or less** (where N is the number of people).

## How it works (Step-by-Step)

### 1. Find the "Bottom Line"
First, we calculate a single net balance for every person:
`Net Balance = (Total you paid) - (Total you owe)`
- If your balance is **positive**, you are owed money.
- If it's **negative**, you owe money.

### 2. Sort the Pools
We split the group into two lists:
1. **Debtors**: People who need to pay.
2. **Creditors**: People who need to get paid.

We sort both lists from largest to smallest. This is the "greedy" part—we always try to settle the biggest debts first to knock them out of the equation quickly.

### 3. Match and Settle
While there are people left in both lists:
1. Take the person who owes the most (Debtor A) and the person who is owed the most (Creditor B).
2. The transfer amount is the smaller of the two values.
3. Debtor A pays Creditor B.
4. Update their balances. If someone hits zero, they are done and removed from the list.
5. Repeat until everyone is at zero.

## Example Flow
Imagine Alice is owed $40, and Bob owes $30, while Charlie owes $10.
1. The engine matches Bob with Alice ($30). Bob is done. Alice still needs $10.
2. The engine matches Charlie with Alice ($10). Both are done.
3. Total transfers: 2.

## Why this works
By matching the biggest debts first, we avoid "fragmenting" the money. Any group of N people can always be settled in at most N-1 steps.

## Complexity
The algorithm is very fast—**O(n log n)**. The only "slow" part is the sorting, which is negligible for any realistic group size (even hundreds of people). We use fixed-precision math (no floats!) to make sure not a single cent is lost in the process.
