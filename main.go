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
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func connectPostgres(connString string) (*pgxpool.Pool, error) {
	cfgPostgres, err := pgxpool.ParseConfig(connString)
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
	v := validator.New()
	cfg, err := config.New()
	if err != nil {
		log.Fatal("Could not parse config: ", err)
	}
	dbpool, errPool := connectPostgres(cfg.PostgresConnBalance)
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
	proto.RegisterBalanceServiceServer(grpcServer, pgHandl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve listener: %s", err)
	}
}
