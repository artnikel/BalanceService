package handler

import (
	"context"
	"testing"

	"github.com/artnikel/BalanceService/bproto"
	"github.com/artnikel/BalanceService/internal/handler/mocks"
	"github.com/artnikel/BalanceService/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: 111.1,
	}
	v = validator.New()
)

func TestDepositAndWithDraw(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &bproto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation,
	}
	srv.On("Deposit", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	_, err := hndl.Deposit(context.Background(), &bproto.DepositRequest{
		Balance: protoBalance,
	})
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &bproto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation,
	}
	srv.On("Deposit", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	_, err := hndl.Deposit(context.Background(), &bproto.DepositRequest{
		Balance: protoBalance,
	})
	require.NoError(t, err)
	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()
	resp, err := hndl.GetBalance(context.Background(), &bproto.GetBalanceRequest{
		Profileid: protoBalance.Profileid,
	})

	require.Equal(t, resp.Money, testBalance.Operation)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}
