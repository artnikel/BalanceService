package handler

import (
	"context"
	"testing"

	"github.com/artnikel/BalanceService/internal/handler/mocks"
	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/proto"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(111.1),
	}
	v = validator.New()
)

func TestBalanceOperation(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &proto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation.InexactFloat64(),
	}
	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	_, err := hndl.BalanceOperation(context.Background(), &proto.BalanceOperationRequest{
		Balance: protoBalance,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestWrnogBalanceOperation(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &proto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: "",
		Operation: testBalance.Operation.InexactFloat64(),
	}
	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	_, err := hndl.BalanceOperation(context.Background(), &proto.BalanceOperationRequest{
		Balance: protoBalance,
	})
	require.Error(t, err)
}

func TestGetBalance(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &proto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation.InexactFloat64(),
	}
	srv.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	_, err := hndl.BalanceOperation(context.Background(), &proto.BalanceOperationRequest{
		Balance: protoBalance,
	})
	require.NoError(t, err)
	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation.InexactFloat64(), nil).Once()
	resp, err := hndl.GetBalance(context.Background(), &proto.GetBalanceRequest{
		Profileid: protoBalance.Profileid,
	})

	require.Equal(t, resp.Money, testBalance.Operation.InexactFloat64())
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetBalanceByWrongID(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &proto.Balance{
		Profileid: "",
	}
	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()
	resp, err := hndl.GetBalance(context.Background(), &proto.GetBalanceRequest{
		Profileid: protoBalance.Profileid,
	})
	require.Error(t, err)
	require.Equal(t, resp.Money, 0.0)
}
