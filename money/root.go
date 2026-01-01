package money

import (
	"sync"

	"github.com/quagmt/udecimal"
)

type Money struct {
	value    udecimal.Decimal
	currency Currency
	mutex    sync.RWMutex
}

func New(value int64, precision uint8, currency Currency) (*Money, error) {
	parsedValue, err := udecimal.NewFromInt64(value, precision)
	if err != nil {
		return nil, err
	}

	return &Money{
		value:    parsedValue,
		currency: currency,
	}, nil
}
