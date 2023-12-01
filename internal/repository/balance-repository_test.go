package repository

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/artnikel/BalanceService/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

var (
	pg          *PgRepository
	testBalance = &model.Balance{
		BalanceID: uuid.New(),
		ProfileID: uuid.New(),
		Operation: decimal.NewFromFloat(100.9),
	}
)

func SetupTestPostgres() (*pgxpool.Pool, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("postgres", "latest", []string{
		"POSTGRES_USER=balanceuser",
		"POSTGRES_PASSWORD=balancepassword",
		"POSTGRES_DB=balancedb"})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	err = RunMigrations(resource.GetPort("5432/tcp"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}
	dbURL := fmt.Sprintf("postgres://balanceuser:balancepassword@localhost:%s/balancedb", resource.GetPort("5432/tcp"))
	cfg, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse dbURL: %w", err)
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect pgxpool: %w", err)
	}
	cleanup := func() {
		dbpool.Close()
		pool.Purge(resource)
	}
	return dbpool, cleanup, nil
}

func RunMigrations(port string) error {
	cmd := exec.Command("flyway",
		"-url=jdbc:postgresql://localhost:"+port+"/balancedb",
		"-user=balanceuser",
		"-password=balancepassword",
		"-locations=filesystem:../../migrations",
		"-connectRetries=10",
		"migrate")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	return nil
}

func TestMain(m *testing.M) {
	dbpool, cleanupPostgres, err := SetupTestPostgres()
	if err != nil {
		fmt.Println("Could not construct the pool: ", err)
		cleanupPostgres()
		os.Exit(1)
	}
	pg = NewPgRepository(dbpool)
	exitVal := m.Run()
	cleanupPostgres()
	os.Exit(exitVal)
}

func TestOperationWithGetBalance(t *testing.T) {
	err := pg.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	money, err := pg.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, money, testBalance.Operation.InexactFloat64())
}

func TestBalanceOperations(t *testing.T) {
	testBalance.ProfileID = uuid.New()
	testBalance.BalanceID = uuid.New()
	testBalance.Operation = decimal.NewFromFloat(800.5)
	err := pg.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	testBalance.Operation = decimal.NewFromFloat(-700.5)
	testBalance.BalanceID = uuid.New()
	err = pg.BalanceOperation(context.Background(), testBalance)
	require.NoError(t, err)
	money, err := pg.GetBalance(context.Background(), testBalance.ProfileID)
	require.NoError(t, err)
	require.Equal(t, decimal.NewFromFloat(money), decimal.NewFromFloat(100.0))
}

func TestGetBalanceByFakeID(t *testing.T) {
	money, _ := pg.GetBalance(context.Background(), uuid.Nil)
	require.Empty(t, money)
	money, _ = pg.GetBalance(context.Background(), uuid.New())
	require.Empty(t, money)
	fakeUUID, err := uuid.Parse("00000000-0000-0000-0000-41db8a3d9113")
	require.NoError(t, err)
	money, _ = pg.GetBalance(context.Background(), fakeUUID)
	require.Empty(t, money)
}
