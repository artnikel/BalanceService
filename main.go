// Package main of a project
package main

import (
	"context"
	"fmt"
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
		log.Fatalf("could not parse config: %v", err)
	}
	dbpool, errPool := connectPostgres(cfg.PostgresConnBalance)
	if errPool != nil {
		log.Fatalf("could not construct the pool: %v", errPool)
	}
	defer dbpool.Close()
	pgRep := repository.NewPgRepository(dbpool)
	pgServ := service.NewBalanceService(pgRep)
	pgHandl := handler.NewEntityBalance(pgServ, v)
	lis, err := net.Listen("tcp", cfg.BalanceAddress)
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	fmt.Println("Balance Service started")
	grpcServer := grpc.NewServer()
	proto.RegisterBalanceServiceServer(grpcServer, pgHandl)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve listener: %s", err)
	}
}
