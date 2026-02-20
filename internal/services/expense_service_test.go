package services

import (
	"testing"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/user/debt-optimization-engine/internal/models"
)

func TestValidateSplits(t *testing.T) {
	u1 := uuid.New()
	u2 := uuid.New()

	tests := []struct {
		name    string
		expense *models.Expense
		wantErr bool
		msg     string
	}{
		{
			name: "Valid splits",
			expense: &models.Expense{
				Amount: decimal.NewFromInt(100),
				Splits: []models.ExpenseSplit{
					{UserID: u1, Amount: decimal.NewFromInt(50)},
					{UserID: u2, Amount: decimal.NewFromInt(50)},
				},
			},
			wantErr: false,
		},
		{
			name: "Sum mismatch",
			expense: &models.Expense{
				Amount: decimal.NewFromInt(100),
				Splits: []models.ExpenseSplit{
					{UserID: u1, Amount: decimal.NewFromInt(50)},
					{UserID: u2, Amount: decimal.NewFromInt(40)},
				},
			},
			wantErr: true,
			msg:     "sum of splits does not equal total amount",
		},
		{
			name: "Duplicate member",
			expense: &models.Expense{
				Amount: decimal.NewFromInt(100),
				Splits: []models.ExpenseSplit{
					{UserID: u1, Amount: decimal.NewFromInt(50)},
					{UserID: u1, Amount: decimal.NewFromInt(50)},
				},
			},
			wantErr: true,
			msg:     "duplicate participant found in splits",
		},
		{
			name: "Negative expense",
			expense: &models.Expense{
				Amount: decimal.NewFromInt(-100),
				Splits: []models.ExpenseSplit{
					{UserID: u1, Amount: decimal.NewFromInt(-100)},
				},
			},
			wantErr: true,
			msg:     "expense amount must be positive",
		},
		{
			name: "Negative split",
			expense: &models.Expense{
				Amount: decimal.NewFromInt(100),
				Splits: []models.ExpenseSplit{
					{UserID: u1, Amount: decimal.NewFromInt(150)},
					{UserID: u2, Amount: decimal.NewFromInt(-50)},
				},
			},
			wantErr: true,
			msg:     "split amount cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateSplits(tt.expense)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.msg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCalculateEqualSplits(t *testing.T) {
	amount := decimal.NewFromFloat(100.00)
	users := []string{"u1", "u2", "u3"}
	
	splits := CalculateEqualSplits(amount, users)
	
	// 100 / 3 = 33.33, 33.33, 33.34
	assert.Equal(t, 3, len(splits))
	assert.True(t, splits[0].Equal(decimal.NewFromFloat(33.33)))
	assert.True(t, splits[1].Equal(decimal.NewFromFloat(33.33)))
	assert.True(t, splits[2].Equal(decimal.NewFromFloat(33.34)))
	
	sum := decimal.Zero
	for _, s := range splits {
		sum = sum.Add(s)
	}
	assert.True(t, sum.Equal(amount))
}
