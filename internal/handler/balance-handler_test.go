package handler

import (
	"context"
	"testing"

	"github.com/artnikel/BalanceService/internal/handler/mocks"
	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/proto"
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
	protoBalance := &proto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation,
	}
	srv.On("Deposit", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	resp, err := hndl.Deposit(context.Background(), &proto.DepositRequest{
		Balance: protoBalance,
	})
	if resp.Error != "" {
		t.Errorf("error %v:", resp.Error)
	}
	require.NoError(t, err)
	srv.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	srv := new(mocks.BalanceService)
	hndl := NewEntityBalance(srv, v)
	protoBalance := &proto.Balance{
		Balanceid: testBalance.BalanceID.String(),
		Profileid: testBalance.ProfileID.String(),
		Operation: testBalance.Operation,
	}
	srv.On("Deposit", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	resp, err := hndl.Deposit(context.Background(), &proto.DepositRequest{
		Balance: protoBalance,
	})
	require.NoError(t, err)
	if resp.Error != "" {
		t.Errorf("error %v:", resp.Error)
	}
	srv.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()
	resp2, err := hndl.GetBalance(context.Background(), &proto.GetBalanceRequest{
		Profileid: protoBalance.Profileid,
	})
	if resp2.Error != "" {
		t.Errorf("error %v:", resp2.Error)
	}
	require.Equal(t, resp2.Money, testBalance.Operation)
	require.NoError(t, err)
	srv.AssertExpectations(t)
}
