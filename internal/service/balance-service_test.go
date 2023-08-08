package service

import (
	"context"
	"testing"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: 200.5,
	}
)

func TestBalanceOperation(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	rep.AssertExpectations(t)
}

func TestGetBalance(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	testBalance.BalanceID = uuid.New()
	testBalance.ProfileID = uuid.New()
	testBalance.Operation = 254.7
	rep.On("BalanceOperation", mock.Anything, mock.AnythingOfType("*model.Balance")).Return(nil).Once()
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(testBalance.Operation, nil).Once()

	err := srv.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)

	money, err := srv.GetBalance(context.Background(), testBalance.BalanceID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation)
	rep.AssertExpectations(t)
}

func TestGetBalanceByWrongID(t *testing.T) {
	rep := new(mocks.BalanceRepository)
	srv := NewBalanceService(rep)
	test := &model.Balance{}
	rep.On("GetBalance", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(test.Operation, nil).Once()
	money, _ := srv.GetBalance(context.Background(), uuid.Nil)
	require.Empty(t, money)
	rep.AssertExpectations(t)
}
