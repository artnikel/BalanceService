// Package errors contains business errors
package errors

const (
	// NotEnoughMoney is error code if user don`t have enough money
	NotEnoughMoney = "NOT_ENOUGH_MONEY"
)

// BusinessError is struct for business errors
type BusinessError struct {
	Code string
}

// New is constructor for manage business errors
func New(code string) *BusinessError {
	return &BusinessError{Code: code}
}

// Error is method for creating business errors
func (bs *BusinessError) Error() string {
	return bs.Code
}
