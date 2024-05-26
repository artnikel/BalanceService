package service

import (
	"context"
	"testing"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/artnikel/BalanceService/internal/service/mocks"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(200.5),
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


