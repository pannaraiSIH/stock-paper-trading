package market

import "errors"

var (
	ErrInvalidQueryParameters    = errors.New("invalid query parameters")
	ErrStockNotFound             = errors.New("stock not found")
	ErrMarketProviderUnavailable = errors.New("market provider unavailable")
)
