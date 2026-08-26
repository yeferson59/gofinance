package money

import (
	"github.com/yeferson59/gofinance/v2/decimal"
)

// Allocate splits m into len(ratios) parts proportional to ratios. The
// parts always sum back to exactly m: any smallest-currency-unit remainder
// left over by rounding is distributed one unit at a time to the earliest
// ratios, following the classic "Fowler" money allocation algorithm. This
// is the safe way to divide an amount (e.g. a shared bill or a profit
// split) among several parties without losing or inventing money through
// rounding.
//
// Money is not forced to its currency's precision, so m may carry a finer
// amount than the currency can express — fractional yen, or an amount computed
// from a rate and not yet rounded to cents. The parts still sum back to
// exactly m in that case; the sub-unit residue goes to the first part, which
// is therefore the only one that can end up finer than the currency's unit.
// Round m first if every part must be a whole unit.
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

	totalDec, err := decimal.NewFromUint64(total, 0)
	if err != nil {
		return nil, err
	}

	results := make([]Money, len(ratios))
	allocated := decimal.Zero

	for i, r := range ratios {
		rDec, err := decimal.NewFromUint64(uint64(r), 0)
		if err != nil {
			return nil, err
		}

		share, err := m.value.TryMul(rDec)
		if err != nil {
			return nil, err
		}

		share, err = share.Div(totalDec)
		if err != nil {
			return nil, err
		}

		share = share.Trunc(prec)

		results[i] = Money{value: share, currency: m.currency}

		allocated, err = allocated.TryAdd(share)
		if err != nil {
			return nil, err
		}
	}

	remainder, err := m.value.TrySub(allocated)
	if err != nil {
		return nil, err
	}

	unit, err := decimal.NewFromInt64(1, prec)
	if err != nil {
		return nil, err
	}
	if remainder.IsNeg() {
		unit = unit.Neg()
	}

	// Hand out one smallest currency unit at a time to the earliest ratios,
	// for as long as a whole unit is left to give.
	for i := 0; i < len(results) && remainder.Abs().GreaterThanOrEqual(unit.Abs()); i++ {
		results[i].value, err = results[i].value.TryAdd(unit)
		if err != nil {
			return nil, err
		}

		remainder, err = remainder.TrySub(unit)
		if err != nil {
			return nil, err
		}
	}

	// An amount carrying finer precision than its own currency — fractional
	// yen, or a computed amount not yet rounded to cents — leaves a residue
	// smaller than one unit, which no whole unit can cover. Give it to the
	// first part so the split still sums back to exactly m. For an amount
	// already at its currency's precision this residue is zero and nothing
	// happens here.
	if !remainder.IsZero() {
		results[0].value, err = results[0].value.TryAdd(remainder)
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
