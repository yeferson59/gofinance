package money

import (
	"github.com/quagmt/udecimal"
)

type Money struct {
	value    udecimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := udecimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}
