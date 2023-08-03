// Package model contains models of using entities
package model

import (
	"github.com/google/uuid"
)

// Balance contains an info about the balance and will be written in a balance table
type Balance struct {
	BalanceID uuid.UUID `json:"balanceid" validate:"required,uuid"`
	ProfileID uuid.UUID `json:"profileid" validate:"required,uuid"`
	Operation float64   `json:"operation" validate:"required,gt=0"`
}
