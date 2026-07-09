package money

import "errors"

// ErrNoAllocationRatios is returned by Allocate when called without any
// ratios to split by.
var ErrNoAllocationRatios = errors.New("money: no allocation ratios given")

// ErrZeroAllocationRatios is returned by Allocate when every ratio is zero,
// making the split undefined.
var ErrZeroAllocationRatios = errors.New("money: allocation ratios sum to zero")

// ErrInvalidAllocationCount is returned by AllocateEvenly when asked to
// split into zero or fewer parts.
var ErrInvalidAllocationCount = errors.New("money: allocation count must be positive")

// Allocate splits m into len(ratios) parts proportional to ratios. The
// parts always sum back to exactly m: any smallest-currency-unit remainder
// left over by rounding is distributed one unit at a time to the earliest
// ratios, following the classic "Fowler" money allocation algorithm. This
// is the safe way to divide an amount (e.g. a shared bill or a profit
// split) among several parties without losing or inventing money through
// rounding.
func (m Money) Allocate(ratios ...uint32) ([]Money, error) {
	if len(ratios) == 0 {
		return nil, ErrNoAllocationRatios
	}

	var total uint64
	for _, r := range ratios {
		total += uint64(r)
	}
	if total == 0 {
		return nil, ErrZeroAllocationRatios
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return nil, err
	}

	totalDec, err := decFromUint64(total, 0)
	if err != nil {
		return nil, err
	}

	results := make([]Money, len(ratios))
	allocated := decZero

	for i, r := range ratios {
		rDec, err := decFromUint64(uint64(r), 0)
		if err != nil {
			return nil, err
		}

		share, err := m.value.Mul(rDec)
		if err != nil {
			return nil, err
		}

		share, err = share.Div(totalDec)
		if err != nil {
			return nil, err
		}

		share = share.Trunc(prec)

		results[i] = Money{value: share, currency: m.currency}

		allocated, err = allocated.Add(share)
		if err != nil {
			return nil, err
		}
	}

	remainder, err := m.value.Sub(allocated)
	if err != nil {
		return nil, err
	}

	unit, err := decFromInt64(1, prec)
	if err != nil {
		return nil, err
	}
	if remainder.IsNeg() {
		unit = unit.Neg()
	}

	for i := 0; !remainder.IsZero() && i < len(results); i++ {
		results[i].value, err = results[i].value.Add(unit)
		if err != nil {
			return nil, err
		}

		remainder, err = remainder.Sub(unit)
		if err != nil {
			return nil, err
		}
	}

	return results, nil
}

// AllocateEvenly splits m into n roughly-equal parts whose sum is exactly
// m, distributing any rounding remainder across the first parts.
func (m Money) AllocateEvenly(n int) ([]Money, error) {
	if n <= 0 {
		return nil, ErrInvalidAllocationCount
	}

	ratios := make([]uint32, n)
	for i := range ratios {
		ratios[i] = 1
	}

	return m.Allocate(ratios...)
}
