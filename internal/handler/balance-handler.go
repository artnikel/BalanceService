// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"strconv"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// BalanceService is an interface that contains methods of service for balance
type BalanceService interface {
	Deposit(ctx context.Context, balance *model.Balance) error
	Withdraw(ctx context.Context, balance *model.Balance) error
	GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error)
}

// EntityBalance contains Balance Service interface
type EntityBalance struct {
	srvBalance BalanceService
	validate   *validator.Validate
	proto.UnimplementedUserServiceServer
}

// NewEntityBalance accepts User Service interface and returns an object of *EntityUser
func NewEntityBalance(srvBalance BalanceService, validate *validator.Validate) *EntityBalance {
	return &EntityBalance{srvBalance: srvBalance, validate: validate}
}

// nolint dupl
// Deposit calls SignUp method of Service by handler
func (b *EntityBalance) Deposit(ctx context.Context, req *proto.DepositRequest) (*proto.DepositResponse, error) {
	createdOperation := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: req.Balance.Operation,
	}

	err := b.validate.VarCtx(ctx, createdOperation.Operation, "required,gt=0")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DepositResponse{
			Error: "failed to validate",
		}, nil
	}

	err = b.srvBalance.Deposit(ctx, createdOperation)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.DepositResponse{
			Error: "failed to Deposit",
		}, nil
	}
	strOperation := strconv.FormatFloat(req.Balance.Operation, 'f', -1, 64)
	return &proto.DepositResponse{
		Operation: "Deposit of " + strOperation + " successfully made",
	}, nil
}

// nolint dupl
// Withdraw calls SignUp method of Service by handler
func (b *EntityBalance) Withdraw(ctx context.Context, req *proto.WithdrawRequest) (*proto.WithdrawResponse, error) {
	createdOperation := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: req.Balance.Operation,
	}

	err := b.validate.VarCtx(ctx, createdOperation.Operation, "required,gt=0")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.WithdrawResponse{
			Error: "failed to validate",
		}, nil
	}

	err = b.srvBalance.Withdraw(ctx, createdOperation)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.WithdrawResponse{
			Error: "failed to Withdraw",
		}, nil
	}
	strOperation := strconv.FormatFloat(req.Balance.Operation, 'f', -1, 64)
	return &proto.WithdrawResponse{
		Operation: "Withdraw of " + strOperation + " successfully made",
	}, nil
}

// GetBalance calls SignUp method of Service by handler
func (b *EntityBalance) GetBalance(ctx context.Context, req *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	id := req.Profileid
	err := b.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{
			Error: "failed to validate",
		}, nil
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{
			Error: "failed to parse id",
		}, nil
	}
	money, err := b.srvBalance.GetBalance(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{
			Error: "failed to Get Balance",
		}, nil
	}
	return &proto.GetBalanceResponse{
		Money: money,
	}, nil
}
