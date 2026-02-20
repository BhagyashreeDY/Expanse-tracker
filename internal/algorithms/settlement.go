package algorithms

import (
	"sort"

	"github.com/shopspring/decimal"
)

// Settlement represents a transaction to settle a debt.
type Settlement struct {
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

// Balance represents the net balance of a user.
type Balance struct {
	User    string
	Amount  decimal.Decimal
}

// SettleOptimized implements the greedy debt minimization algorithm.
func SettleOptimized(balances map[string]decimal.Decimal) []Settlement {
	var debtors []Balance
	var creditors []Balance

	for user, amount := range balances {
		if amount.IsNegative() {
			debtors = append(debtors, Balance{User: user, Amount: amount.Abs()})
		} else if amount.IsPositive() {
			creditors = append(creditors, Balance{User: user, Amount: amount})
		}
	}

	// Sort to match largest debtor with largest creditor (Greedy approach)
	sortBalances := func(b []Balance) {
		sort.Slice(b, func(i, j int) bool {
			return b[i].Amount.GreaterThan(b[j].Amount)
		})
	}

	sortBalances(debtors)
	sortBalances(creditors)

	return match(debtors, creditors)
}

// SettleNaive implements a simple pairwise settlement without optimization (FIFO matching).
func SettleNaive(balances map[string]decimal.Decimal) []Settlement {
	var debtors []Balance
	var creditors []Balance

	for user, amount := range balances {
		if amount.IsNegative() {
			debtors = append(debtors, Balance{User: user, Amount: amount.Abs()})
		} else if amount.IsPositive() {
			creditors = append(creditors, Balance{User: user, Amount: amount})
		}
	}

	// No sorting - just pure naive sequential matching
	return match(debtors, creditors)
}

func match(debtors []Balance, creditors []Balance) []Settlement {
	var settlements []Settlement
	i, j := 0, 0

	// Use pointers to allow modification during iteration
	for i < len(debtors) && j < len(creditors) {
		debtor := &debtors[i]
		creditor := &creditors[j]

		settleAmount := decimal.Min(debtor.Amount, creditor.Amount)

		if settleAmount.IsPositive() {
			settlements = append(settlements, Settlement{
				From:   debtor.User,
				To:     creditor.User,
				Amount: settleAmount,
			})
		}

		debtor.Amount = debtor.Amount.Sub(settleAmount)
		creditor.Amount = creditor.Amount.Sub(settleAmount)

		if debtor.Amount.IsZero() {
			i++
		}
		if creditor.Amount.IsZero() {
			j++
		}
	}

	return settlements
}
