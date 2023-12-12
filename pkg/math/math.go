package math

import (
	"math"

	"github.com/shopspring/decimal"
)

func AbsInt(x int) int {
	y := x >> 31
	return (x ^ y) - y
}

func AbsInt32(x int32) int32 {
	y := x >> 31
	return (x ^ y) - y
}

func AbsInt64(x int64) int64 {
	y := x >> 63
	return (x ^ y) - y
}

// DecimalPlacesWithDigitsNumber Round to decimal places by number of digits
func DecimalPlacesWithDigitsNumber(value float64, digitsNumber int) float64 {
	//if Inf
	if math.IsInf(value, 0) {
		return 0
	}
	f, _ := decimal.NewFromFloat(value).Round(int32(digitsNumber)).Float64()
	return f
}
