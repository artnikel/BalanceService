// Package model contains models of using entities
package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID       `json:"balanceid" validate:"required,uuid"`
	ProfileID uuid.UUID       `json:"profileid" validate:"required,uuid"`
	Operation decimal.Decimal `json:"operation" validate:"required"`
}
