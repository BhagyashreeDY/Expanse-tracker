package services

import (
	"errors"

	"github.com/shopspring/decimal"
	"github.com/user/debt-optimization-engine/internal/models"
)

// ValidateSplits ensures that the sum of split amounts matches the total expense amount
// and that there are no duplicate participants or negative amounts.
func ValidateSplits(expense *models.Expense) error {
	if expense.Amount.IsNegative() || expense.Amount.IsZero() {
		return errors.New("expense amount must be positive")
	}

	if len(expense.Splits) == 0 {
		return errors.New("expense must have at least one participant")
	}

	sum := decimal.NewFromInt(0)
	members := make(map[string]bool)

	for _, split := range expense.Splits {
		if split.Amount.IsNegative() {
			return errors.New("split amount cannot be negative")
		}

		uid := split.UserID.String()
		if members[uid] {
			return errors.New("duplicate participant found in splits")
		}
		members[uid] = true
		sum = sum.Add(split.Amount)
	}

	if !sum.Equal(expense.Amount) {
		return errors.New("sum of splits does not equal total amount")
	}

	return nil
}

// CalculateEqualSplits distributes the amount equally among participants, handling rounding drift.
func CalculateEqualSplits(amount decimal.Decimal, userIDs []string) []decimal.Decimal {
	if len(userIDs) == 0 {
		return nil
	}
	numParticipants := decimal.NewFromInt(int64(len(userIDs)))
	
	// Use 2 decimal places for financial calculations
	individualAmount := amount.DivRound(numParticipants, 2)
	
	splits := make([]decimal.Decimal, len(userIDs))
	runningSum := decimal.Zero

	for i := 0; i < len(userIDs)-1; i++ {
		splits[i] = individualAmount
		runningSum = runningSum.Add(individualAmount)
	}

	// Last person takes the remaining to avoid rounding drift
	splits[len(userIDs)-1] = amount.Sub(runningSum)

	return splits
}
