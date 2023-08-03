// Package repository is a lower level of project
package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("PgRepository-BalanceOperation: error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				log.Fatalf("PgRepository-BalanceOperation: error in rollback: %v", errRollback)
			}
		} else {
			errCommit := tx.Commit(ctx)
			if errCommit != nil {
				log.Fatalf("PgRepository-BalanceOperation: error committing transaction: %v", errCommit)
			}
		}
	}()
	_, err = tx.Exec(ctx, "INSERT INTO balance (balanceid, profileid, operation) VALUES ($1, $2, $3)",
		balance.BalanceID, balance.ProfileID, balance.Operation)
	if err != nil {
		return fmt.Errorf("PgRepository-BalanceOperation: error in method tx.Exec(): %w", err)
	}

	return nil
}

// GetBalance counted sum of operations and returns balance of profile by him id
func (p *PgRepository) GetBalance(ctx context.Context, profileID uuid.UUID) (float64, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("PgRepository-GetBalance: error starting transaction: %w", err)
	}
	defer func() {
		if err != nil {
			errRollback := tx.Rollback(ctx)
			if errRollback != nil {
				log.Fatalf("PgRepository-BalanceOperation: error in rollback: %v", errRollback)
			}
		} else {
			errCommit := tx.Commit(ctx)
			if errCommit != nil {
				log.Fatalf("PgRepository-BalanceOperation: error committing transaction: %v", errCommit)
			}
		}
	}()

	rows, err := tx.Query(ctx, "SELECT operation FROM balance WHERE profileid = $1 FOR UPDATE", profileID)
	if err != nil {
		return 0, fmt.Errorf("PgRepository-GetBalance: error in method tx.QueryRow(): %w", err)
	}
	defer rows.Close()

	var money float64

	for rows.Next() {
		var operation float64
		err := rows.Scan(&operation)
		if err != nil {
			return 0, fmt.Errorf("BalanceService-GetBalance: error in method rows.Scan:%w", err)
		}
		money += operation
	}

	return money, nil
}
