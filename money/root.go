package money

import (
	"github.com/quagmt/udecimal"
)

type Money struct {
	udecimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := udecimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{
		Decimal:  parsedValue,
		currency: currency,
	}, nil
}
