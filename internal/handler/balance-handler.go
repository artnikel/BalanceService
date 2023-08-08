// Package handler is the top level of the application and it contains request handlers
package handler

import (
	"context"
	"fmt"
	"strconv"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// BalanceService is an interface that contains methods of service for balance
type BalanceService interface {
	BalanceOperation(ctx context.Context, balance *model.Balance) error
	GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error)
}

// EntityBalance contains Balance Service interface
type EntityBalance struct {
	srvBalance BalanceService
	validate   *validator.Validate
	proto.UnimplementedBalanceServiceServer
}

// NewEntityBalance accepts User Service interface and returns an object of *EntityUser
func NewEntityBalance(srvBalance BalanceService, validate *validator.Validate) *EntityBalance {
	return &EntityBalance{srvBalance: srvBalance, validate: validate}
}

// BalanceOperation calls BalanceOperation method of Service by handler
func (b *EntityBalance) BalanceOperation(ctx context.Context, req *proto.BalanceOperationRequest) (*proto.BalanceOperationResponse, error) {
	profileid := req.Balance.Profileid
	err := b.validate.VarCtx(ctx, profileid, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.BalanceOperationResponse{}, fmt.Errorf("EntityBalance-BalanceOperation: failed to validate profile id")
	}
	profileUUID, err := uuid.Parse(profileid)
	if err != nil {
		return &proto.BalanceOperationResponse{}, fmt.Errorf("EntityBalance-BalanceOperation: failed to parse Profileid")
	}
	createdOperation := &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: profileUUID,
		Operation: req.Balance.Operation,
	}
	err = b.srvBalance.BalanceOperation(ctx, createdOperation)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.BalanceOperationResponse{}, fmt.Errorf("EntityBalance-BalanceOperation: failed to made opeartion")
	}
	strOperation := strconv.FormatFloat(req.Balance.Operation, 'f', -1, 64)
	return &proto.BalanceOperationResponse{
		Operation: strOperation,
	}, nil
}

// GetBalance is mecalls SignUp method of Service by handler
func (b *EntityBalance) GetBalance(ctx context.Context, req *proto.GetBalanceRequest) (*proto.GetBalanceResponse, error) {
	id := req.Profileid
	err := b.validate.VarCtx(ctx, id, "required,uuid")
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{}, fmt.Errorf("EntityBalance-GetBalance: failed to validate id")
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{}, fmt.Errorf("EntityBalance-GetBalance: failed to parse id")
	}
	money, err := b.srvBalance.GetBalance(ctx, idUUID)
	if err != nil {
		logrus.Errorf("error: %v", err)
		return &proto.GetBalanceResponse{}, fmt.Errorf("EntityBalance-GetBalance: failed to get balance")
	}
	return &proto.GetBalanceResponse{
		Money: money,
	}, nil
}
