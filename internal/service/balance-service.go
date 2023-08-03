// Package service contains business logic of a project
package service

import (
	"context"
	"fmt"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/google/uuid"
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

// Deposit is a method of BalanceService that calls  method of Repository
func (b *BalanceService) Deposit(ctx context.Context, balance *model.Balance) error {
	if balance.Operation <= 0 {
		return fmt.Errorf("BalanceService-Deposit: the amount to be deposited must be greater than zero")
	}
	err := b.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return fmt.Errorf("BalanceService-Deposit-BalanceOperation: error:%v", err)
	}
	return nil
}

// Withdraw is a method of BalanceService that calls  method of Repository
func (b *BalanceService) Withdraw(ctx context.Context, balance *model.Balance) error {
	if balance.Operation <= 0 {
		return fmt.Errorf("BalanceService-Withdraw: the amount to be withdrawed must be greater than zero")
	}
	balance.Operation = -balance.Operation
	err := b.bRep.BalanceOperation(ctx, balance)
	if err != nil {
		return fmt.Errorf("BalanceService-Withdraw-BalanceOperation: error:%v", err)
	}
	return nil
}

// GetBalance is a method of BalanceService that calls  method of Repository
func (b *BalanceService) GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error) {
	money, err := b.bRep.GetBalance(ctx, profileID)
	if err != nil {
		return 0, fmt.Errorf("BalanceService-GetBalance-GetBalance: error:%v", err)
	}
	return money, nil
}
