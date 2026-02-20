package services

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/user/debt-optimization-engine/internal/algorithms"
	"github.com/user/debt-optimization-engine/internal/models"
	"github.com/user/debt-optimization-engine/internal/repositories"
)

type SettlementService struct {
	repo repositories.Repository
}

func NewSettlementService(repo repositories.Repository) *SettlementService {
	return &SettlementService{repo: repo}
}

func (s *SettlementService) CalculateBalances(ctx context.Context, groupID string, from, to *time.Time) (map[string]decimal.Decimal, error) {
	expenses, err := s.repo.GetExpensesByGroup(ctx, groupID, from, to)
	if err != nil { return nil, err }

	members, err := s.repo.GetGroupMembers(ctx, groupID)
	if err != nil { return nil, err }

	userMap := make(map[string]string)
	balances := make(map[string]decimal.Decimal)
	for _, m := range members {
		userMap[m.ID.String()] = m.Username
		balances[m.Username] = decimal.Zero
	}

	for _, exp := range expenses {
		payerName := userMap[exp.PayerID.String()]
		balances[payerName] = balances[payerName].Add(exp.Amount)
		for _, split := range exp.Splits {
			userName := userMap[split.UserID.String()]
			balances[userName] = balances[userName].Sub(split.Amount)
		}
	}

	return balances, nil
}

func (s *SettlementService) GetSettlement(ctx context.Context, groupID string, from, to *time.Time) (*models.SettlementResponse, error) {
	balances, err := s.CalculateBalances(ctx, groupID, from, to)
	if err != nil { return nil, err }

	optimized := algorithms.SettleOptimized(balances)

	rawCount := 0
	expenses, _ := s.repo.GetExpensesByGroup(ctx, groupID, from, to)
	for _, e := range expenses { rawCount += len(e.Splits) }

	gain := "0%"
	if rawCount > 0 {
		reduction := float64(rawCount-len(optimized)) / float64(rawCount) * 100
		if reduction < 0 { reduction = 0 }
		gain = fmt.Sprintf("%.1f%%", reduction)
	}

	return &models.SettlementResponse{
		Transactions:      optimized,
		TotalTransactions: len(optimized),
		OptimizationGain:  gain,
		RawBalances:       balances,
	}, nil
}

func (s *SettlementService) CompareStrategies(ctx context.Context, groupID string) (*models.SettlementComparison, error) {
	balances, err := s.CalculateBalances(ctx, groupID, nil, nil)
	if err != nil { return nil, err }

	optimized := algorithms.SettleOptimized(balances)
	naive := algorithms.SettleNaive(balances)

	optStats := s.calculateStats(optimized)
	naiveStats := s.calculateStats(naive)

	if naiveStats.TransactionCount > 0 {
		gain := float64(naiveStats.TransactionCount-optStats.TransactionCount) / float64(naiveStats.TransactionCount) * 100
		if gain < 0 { gain = 0 }
		optStats.OptimizationGain = fmt.Sprintf("%.1f%%", gain)
	}

	return &models.SettlementComparison{
		Greedy:   optStats,
		Baseline: naiveStats,
	}, nil
}

func (s *SettlementService) calculateStats(txs []algorithms.Settlement) models.SettlementStats {
	vol := decimal.Zero
	for _, t := range txs {
		vol = vol.Add(t.Amount)
	}
	return models.SettlementStats{
		TransactionCount: len(txs),
		TotalVolume:      vol,
	}
}
