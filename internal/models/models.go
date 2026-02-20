package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type GroupMember struct {
	GroupID  uuid.UUID `json:"group_id"`
	UserID   uuid.UUID `json:"user_id"`
	JoinedAt time.Time `json:"joined_at"`
}

type SplitType string

const (
	SplitEqual      SplitType = "EQUAL"
	SplitPercentage SplitType = "PERCENTAGE"
	SplitExact      SplitType = "EXACT"
)

type Expense struct {
	ID          uuid.UUID       `json:"id"`
	GroupID     uuid.UUID       `json:"group_id"`
	PayerID     uuid.UUID       `json:"payer_id"`
	Amount      decimal.Decimal `json:"amount"`
	Description string          `json:"description"`
	SplitType   SplitType       `json:"split_type"`
	CreatedAt   time.Time       `json:"created_at"`
	Splits      []ExpenseSplit  `json:"splits"`
}

type ExpenseSplit struct {
	ExpenseID uuid.UUID       `json:"expense_id"`
	UserID    uuid.UUID       `json:"user_id"`
	Amount    decimal.Decimal `json:"amount"`
}

type GroupBalances struct {
	GroupID  uuid.UUID                  `json:"group_id"`
	Balances map[string]decimal.Decimal `json:"balances"` // Username -> Balance
}

type SettlementResponse struct {
	Transactions         interface{} `json:"transactions"`
	TotalTransactions    int         `json:"total_transactions"`
	OptimizationGain     string      `json:"optimization_gain"`
	RawBalances          interface{} `json:"raw_balances,omitempty"`
}

type SettlementComparison struct {
	Greedy   SettlementStats `json:"greedy"`
	Baseline SettlementStats `json:"baseline"`
}

type SettlementStats struct {
	TransactionCount int             `json:"transaction_count"`
	TotalVolume      decimal.Decimal `json:"total_volume"`
	OptimizationGain string          `json:"optimization_gain"`
}

type SettlementPayment struct {
	ID         uuid.UUID       `json:"id"`
	GroupID    uuid.UUID       `json:"group_id"`
	FromUserID uuid.UUID       `json:"from_user_id"`
	ToUserID   uuid.UUID       `json:"to_user_id"`
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  time.Time       `json:"created_at"`
}

func ParseUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
