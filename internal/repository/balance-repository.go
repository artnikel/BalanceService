// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

// PgRepository represents the PostgreSQL repository implementation.
type PgRepository struct {
	pool *pgxpool.Pool
}

// NewPgRepository creates and returns a new instance of PgRepository, using the provided pgxpool.Pool.
func NewPgRepository(pool *pgxpool.Pool) *PgRepository {
	return &PgRepository{
		pool: pool,
	}
}

// BalanceOperation allows to record a deposit or withdrawal transaction in the database
func (p *PgRepository) BalanceOperation(ctx context.Context, balance *model.Balance) error {
	_, err := p.pool.Exec(ctx, "INSERT INTO balance (balanceid, profileid, operation) VALUES ($1, $2, $3)",
		balance.BalanceID, balance.ProfileID, balance.Operation)
	if err != nil {
		return fmt.Errorf("exec %w", err)
	}

	return nil
}

// GetBalance counted sum of operations and returns balance of profile by him id
func (p *PgRepository) GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error) {
	rows, err := p.pool.Query(ctx, "SELECT operation FROM balance WHERE profileid = $1 FOR UPDATE", profileID)
	if err != nil {
		return 0, fmt.Errorf("query %w", err)
	}
	defer rows.Close()

	var money decimal.Decimal

	for rows.Next() {
		var operation decimal.Decimal
		err := rows.Scan(&operation)
		if err != nil {
			return 0, fmt.Errorf("scan %w", err)
		}
		money = money.Add(operation)
	}

	return money.InexactFloat64(), nil
}
