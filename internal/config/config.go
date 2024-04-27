// Package config with environment variables
package config

import "github.com/caarlos0/env"

// Variables is a struct with environment variables
type Variables struct {
	PostgresConnBalance string `env:"POSTGRES_CONN_BALANCE"`
	BalanceAddress      string `env:"BALANCE_ADDRESS"`
}

// New returns parsed object of config
func New() (*Variables, error) {
	cfg := &Variables{}
	err := env.Parse(cfg)
	return cfg, err
}
