// Package config with environment variables
package config

// Variables is a struct with environment variables
type Variables struct {
	PostgresConnBalance string `env:"POSTGRES_CONN_BALANCE"`
}
