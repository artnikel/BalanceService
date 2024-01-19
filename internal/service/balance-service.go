// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	berrors "github.com/artnikel/BalanceService/internal/errors"
)

// BalanceRepository is interface with methods for balance operations
type BalanceRepository interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) error
	GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error)
}

// BalanceService contains BalanceRepository interface
type BalanceService struct {
	bRep BalanceRepository
}

// NewBalanceService accepts BalanceRepository object and returnes an object of type *BalanceService
func NewBalanceService(bRep BalanceRepository) *BalanceService {
	return &BalanceService{bRep: bRep}
}

// BalanceOperation is a method of BalanceService that calls  method of Repository
func (b *BalanceService) BalanceOperation(ctx context.Context, balance *model.Balance) error {
	if balance.Operation.IsNegative() {
		money, err := b.GetBalance(ctx, balance.ProfileID)
		if err != nil {
			return fmt.Errorf("getBalance %w", err)
		}
		if decimal.NewFromFloat(money).Cmp(balance.Operation.Abs()) == 1 {
			err = b.bRep.BalanceOperation(ctx, balance)
			if err != nil {
				return fmt.Errorf("balanceOperation %w", err)
			}
			return nil
		}
		return berrors.New(berrors.NotEnoughMoney)
	}
	err := b.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return fmt.Errorf("balanceOperation %w", err)
	}
	return nil
}

// GetBalance is a method of BalanceService that calls  method of Repository
func (b *BalanceService) GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error) {
	money, err := b.bRep.GetBalance(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("getBalance %w", err)
	}
	return money, nil
}
