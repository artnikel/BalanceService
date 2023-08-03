// Package main of a project
package main

import (
	"context"
	"log"
	"net"

	"github.com/artnikel/BalanceService/internal/config"
	"github.com/artnikel/BalanceService/internal/handler"
	"github.com/artnikel/BalanceService/internal/repository"
	"github.com/artnikel/BalanceService/internal/service"
	"github.com/artnikel/BalanceService/proto"
	"github.com/caarlos0/env"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func connectPostgres(cfg *config.Variables) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(cfg.PostgresConnBalance)
	if err != nil {
		return nil, err
	}
	dbpool, err := pgxpool.NewWithConfig(context.Background(), cfgPostgres)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

// nolint gocritic
func main() {
	var (
		cfg config.Variables
		v   = validator.New()
	)
	if err := env.Parse(&cfg); err != nil {
		log.Fatal("could not parse config: ", err)
	}
	dbpool, errPool := connectPostgres(&cfg)
	if errPool != nil {
		log.Fatal("could not construct the pool: ", errPool)
	}
	defer dbpool.Close()
	pgRep := repository.NewPgRepository(dbpool)
	pgServ := service.NewBalanceService(pgRep)
	pgHandl := handler.NewEntityBalance(pgServ, v)
	lis, err := net.Listen("tcp", "localhost:8095")
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterUserServiceServer(grpcServer, pgHandl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve listener: %s", err)
	}
}
