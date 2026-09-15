package trading

import "errors"

var (
	ErrAccountAlreadyExists = errors.New("account already exists")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrInsufficientShares   = errors.New("insufficient shares")
	ErrPositionNotFound     = errors.New("position not found")
	ErrInvalidOrderSide     = errors.New("invalid order side")
	ErrPriceNotFound        = errors.New("price not found")
)
