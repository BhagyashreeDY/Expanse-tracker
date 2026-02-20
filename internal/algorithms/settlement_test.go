package algorithms

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestSettlementAlgorithms(t *testing.T) {
	balances := map[string]decimal.Decimal{
		"Alice":   decimal.NewFromInt(100),
		"Bob":     decimal.NewFromInt(50),
		"Charlie": decimal.NewFromInt(-80),
		"David":   decimal.NewFromInt(-70),
	}

	t.Run("Optimized (Greedy)", func(t *testing.T) {
		result := SettleOptimized(balances)
		// With greedy on these balances:
		// Debtors: Charlie(80), David(70)
		// Creditors: Alice(100), Bob(50)
		// 1. Charlie pays Alice 80 -> Alice needs 20, Charlie settled
		// 2. David pays Alice 20 -> Alice settled, David needs 50
		// 3. David pays Bob 50 -> David settled, Bob settled
		// Total 3 transactions.
		assert.LessOrEqual(t, len(result), 3)
		
		vol := decimal.Zero
		for _, tx := range result {
			vol = vol.Add(tx.Amount)
		}
		assert.True(t, vol.Equal(decimal.NewFromInt(150)))
	})

	t.Run("Naive", func(t *testing.T) {
		result := SettleNaive(balances)
		// Naive matching should also result in correct volume settlement
		vol := decimal.Zero
		for _, tx := range result {
			vol = vol.Add(tx.Amount)
		}
		assert.True(t, vol.Equal(decimal.NewFromInt(150)))
	})
}

func TestSettleNaiveConsistency(t *testing.T) {
	balances := map[string]decimal.Decimal{
		"A": decimal.NewFromInt(10),
		"B": decimal.NewFromInt(-10),
	}
	result := SettleNaive(balances)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "B", result[0].From)
	assert.Equal(t, "A", result[0].To)
	assert.True(t, result[0].Amount.Equal(decimal.NewFromInt(10)))
}
