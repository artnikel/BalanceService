// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/artnikel/BalanceService/bproto"
	"github.com/artnikel/BalanceService/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// BalanceService is an interface that contains methods of service for balance
type BalanceService interface {
	Deposit(ctx context.Context, balance *model.Balance) error
	// Withdraw(ctx context.Context, balance *model.Balance) error
	GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error)
}

// EntityBalance contains Balance Service interface
type EntityBalance struct {
	srvBalance BalanceService
	validate   *validator.Validate
	bproto.UnimplementedBalanceServiceServer
}

// NewEntityBalance accepts User Service interface and returns an object of *EntityUser
func NewEntityBalance(srvBalance BalanceService, validate *validator.Validate) *EntityBalance {
	return &EntityBalance{srvBalance: srvBalance, validate: validate}
}

// nolint dupl
// Deposit calls SignUp method of Service by handler
func (b *EntityBalance) Deposit(ctx context.Context, req *bproto.DepositRequest) (*bproto.DepositResponse, error) {
	profileUUID, err := uuid.Parse(req.Balance.Profileid)
	if err != nil {
		return &bproto.DepositResponse{}, fmt.Errorf("failed to parse Profileid")
	}
	createdOperation := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: profileUUID,
		Operation: req.Balance.Operation,
	}

	// err = b.validate.VarCtx(ctx, createdOperation.Operation, "required,gt=0")
	// if err != nil {
	// 	logrus.Errorf("error: %v", err)
	// 	return &bproto.DepositResponse{}, fmt.Errorf("failed to validate")
	// }

	err = b.srvBalance.Deposit(ctx, createdOperation)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &bproto.DepositResponse{}, fmt.Errorf("failed to Deposit")
	}
	strOperation := strconv.FormatFloat(req.Balance.Operation, 'f', -1, 64)
	return &bproto.DepositResponse{
		Operation: strOperation,
	}, nil
}

// // nolint dupl
// // Withdraw calls SignUp method of Service by handler
// func (b *EntityBalance) Withdraw(ctx context.Context, req *bproto.WithdrawRequest) (*bproto.WithdrawResponse, error) {
// 	profileUUID, err := uuid.Parse(req.Balance.Profileid)
// 	if err != nil {
// 		return &bproto.WithdrawResponse{}, fmt.Errorf("failed to parse profileid")
// 	}
// 	createdOperation := &model.Balance{
// 		BalanceID: uuid.New(),
// 		ProfileID: profileUUID,
// 		Operation: req.Balance.Operation,
// 	}

// 	err = b.validate.VarCtx(ctx, createdOperation.Operation, "required,gt=0")
// 	if err != nil {
// 		logrus.Errorf("error: %v", err)
// 		return &bproto.WithdrawResponse{}, fmt.Errorf("failed to validate")
// 	}

// 	err = b.srvBalance.Withdraw(ctx, createdOperation)
// 	if err != nil {
// 		logrus.Errorf("error: %v", err)
// 		return &bproto.WithdrawResponse{}, fmt.Errorf("failed to withdraw")
// 	}
// 	strOperation := strconv.FormatFloat(req.Balance.Operation, 'f', -1, 64)
// 	return &bproto.WithdrawResponse{
// 		Operation: strOperation,
// 	}, nil
// }

// GetBalance calls SignUp method of Service by handler
func (b *EntityBalance) GetBalance(ctx context.Context, req *bproto.GetBalanceRequest) (*bproto.GetBalanceResponse, error) {
	id := req.Profileid
	// err := b.validate.VarCtx(ctx, id, "required,uuid")
	// if err != nil {
	// 	logrus.Errorf("error: %v", err)
	// 	return &bproto.GetBalanceResponse{}, fmt.Errorf("failed to validate")
	// }
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &bproto.GetBalanceResponse{}, fmt.Errorf("failed to parse id")
	}
	money, err := b.srvBalance.GetBalance(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &bproto.GetBalanceResponse{}, fmt.Errorf("failed to Get Balance")
	}
	return &bproto.GetBalanceResponse{
		Money: money,
	}, nil
}
