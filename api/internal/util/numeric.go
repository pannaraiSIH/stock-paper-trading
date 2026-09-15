package util

import (
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/shopspring/decimal"
)

func ToDecimal(n pgtype.Numeric) decimal.Decimal {
	return decimal.NewFromBigInt(n.Int, n.Exp)
}

func ToNumeric(c decimal.Decimal) pgtype.Numeric {
	return pgtype.Numeric{
		Int:   new(big.Int).Set(c.Coefficient()),
		Exp:   c.Exponent(),
		Valid: true,
	}
}
