package money

import "errors"

var (
	// ErrNoAllocationRatios is returned by Allocate when called without any
	// ratios to split by.
	ErrNoAllocationRatios = errors.New("money: no allocation ratios given")

	// ErrZeroAllocationRatios is returned by Allocate when every ratio is zero,
	// making the split undefined.
	ErrZeroAllocationRatios = errors.New("money: allocation ratios sum to zero")

	// ErrInvalidAllocationCount is returned by AllocateEvenly when asked to
	// split into zero or fewer parts.
	ErrInvalidAllocationCount = errors.New("money: allocation count must be positive")

	// ErrInvalidExchangeRate is returned by Convert when given a rate that is
	// zero or negative.
	ErrInvalidExchangeRate = errors.New("money: exchange rate must be positive")

	// ErrCurrencyMismatch is returned by operations that require both operands
	// to share the same currency, such as SafeAdd and SafeSub.
	ErrCurrencyMismatch = errors.New("money: currency mismatch")
)
