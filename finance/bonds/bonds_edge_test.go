package bonds

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/daycount"
	"github.com/yeferson59/gofinance/v2/money"
)

// TestLongBondYTMNeverPanics is the regression test for the worst defect this
// package carried: YTM scans candidate yields up to 10000%, and at those
// yields a long bond's discount factor (1 + y/f)^n overflows the decimal
// engine. The overflow was raised through the panicking Add/Mul helpers, and
// the scan aborted on any candidate that failed to evaluate — so YTM panicked
// for every bond longer than about eleven years, including the 30-year
// semiannual bond, the most common instrument there is.
func TestLongBondYTMNeverPanics(t *testing.T) {
	maturities := []struct {
		name    string
		periods int
		freq    int
	}{
		{"12-year semiannual", 24, 2},
		{"30-year semiannual", 60, 2},
		{"30-year monthly", 360, 12},
		{"100-year monthly", 1200, 12},
	}

	for _, maturity := range maturities {
		t.Run(maturity.name, func(t *testing.T) {
			priced := NewBond().Face(1000, money.USD).CouponRate(0.05).
				Frequency(maturity.freq).Periods(maturity.periods).Yield(0.06)

			price, err := priced.Price()
			require.NoError(t, err)

			assert.NotPanics(t, func() {
				ytm, err := NewBond().Face(1000, money.USD).CouponRate(0.05).
					Frequency(maturity.freq).Periods(maturity.periods).
					MarketPrice(price.InexactFloat64()).YTM()
				require.NoError(t, err)
				assert.InDelta(t, 0.06, ytm.InexactFloat64(), 1e-6)
			})
		})
	}
}

// TestPriceAtExtremeYieldReportsError checks the other half of the same fix: a
// yield high enough to overflow the discount factor must come back as an
// error, not a panic or a wrong price.
func TestPriceAtExtremeYieldReportsError(t *testing.T) {
	bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(60).Yield(500)

	assert.NotPanics(t, func() {
		_, err := bond.Price()
		assert.Error(t, err)
	})
}

// TestPriceYTMRoundTripGrid checks the two solvers invert each other across a
// grid of coupons, maturities and frequencies rather than at a single point.
func TestPriceYTMRoundTripGrid(t *testing.T) {
	coupons := []float64{0, 0.03, 0.05, 0.09}
	yields := []float64{0.01, 0.05, 0.08, 0.15}
	terms := []struct {
		periods int
		freq    int
	}{
		{2, 1}, {10, 2}, {24, 2}, {40, 4},
	}

	for _, coupon := range coupons {
		for _, yield := range yields {
			for _, term := range terms {
				bond := NewBond().Face(1000, money.USD).CouponRate(coupon).
					Frequency(term.freq).Periods(term.periods)

				price, err := bond.Yield(yield).Price()
				require.NoError(t, err)
				require.Positive(t, price.InexactFloat64())

				recovered, err := bond.MarketPrice(price.InexactFloat64()).YTM()
				require.NoError(t, err)

				assert.InDelta(t, yield, recovered.InexactFloat64(), 1e-6,
					"coupon %v yield %v periods %d freq %d", coupon, yield, term.periods, term.freq)
			}
		}
	}
}

// TestPriceRelationToPar pins the three standard cases: a bond trades at par
// when its coupon equals its yield, at a premium below it, and at a discount
// above it.
func TestPriceRelationToPar(t *testing.T) {
	bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20)

	atPar, err := bond.Yield(0.05).Price()
	require.NoError(t, err)
	assert.InDelta(t, 1000.0, atPar.InexactFloat64(), 1e-6)

	premium, err := bond.Yield(0.03).Price()
	require.NoError(t, err)
	assert.Greater(t, premium.InexactFloat64(), 1000.0)

	discount, err := bond.Yield(0.07).Price()
	require.NoError(t, err)
	assert.Less(t, discount.InexactFloat64(), 1000.0)
}

// TestPriceFallsAsYieldRises checks the monotonicity the yield search depends
// on: price is strictly decreasing in yield, which is what makes a single
// bracketed root exist.
func TestPriceFallsAsYieldRises(t *testing.T) {
	bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20)

	previous := 0.0

	for i, yield := range []float64{0, 0.02, 0.05, 0.10, 0.25, 0.50} {
		price, err := bond.Yield(yield).Price()
		require.NoError(t, err)

		if i > 0 {
			assert.Less(t, price.InexactFloat64(), previous)
		}

		previous = price.InexactFloat64()
	}
}

// TestNegativeYield covers a bond priced at a negative yield, which real
// sovereign markets have traded at: the price rises above the undiscounted sum
// of the cash flows.
func TestNegativeYield(t *testing.T) {
	bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(10).Yield(-0.01)

	price, err := bond.Price()
	require.NoError(t, err)

	// Coupons total 250 and the face is 1000; discounting at a negative rate
	// must value them above that undiscounted 1250.
	assert.Greater(t, price.InexactFloat64(), 1250.0)
}

// TestYieldAtNegativeFrequencyBound checks the documented lower limit: a yield
// of exactly −frequency makes the per-period discount factor zero.
func TestYieldAtNegativeFrequencyBound(t *testing.T) {
	_, err := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(10).
		Yield(-2).Price()
	assert.ErrorIs(t, err, ErrInvalidYield)

	_, err = NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(10).
		Yield(-3).Price()
	assert.ErrorIs(t, err, ErrInvalidYield)
}

// TestModifiedDurationIdentity checks the definition connecting the two
// duration measures: modified = Macaulay / (1 + y/f).
func TestModifiedDurationIdentity(t *testing.T) {
	for _, yield := range []float64{0.02, 0.06, 0.12} {
		bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20).Yield(yield)

		macaulay, err := bond.MacaulayDuration()
		require.NoError(t, err)

		modified, err := bond.ModifiedDuration()
		require.NoError(t, err)

		assert.InDelta(t, macaulay.InexactFloat64()/(1+yield/2), modified.InexactFloat64(), 1e-9)
		assert.Less(t, modified.InexactFloat64(), macaulay.InexactFloat64())
	}
}

// TestDurationBelowMaturity checks that a coupon-paying bond's duration is
// strictly shorter than its maturity — cash arrives before the redemption —
// while a zero-coupon bond's equals it exactly.
func TestDurationBelowMaturity(t *testing.T) {
	coupon := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20).Yield(0.05)

	macaulay, err := coupon.MacaulayDuration()
	require.NoError(t, err)
	assert.Less(t, macaulay.InexactFloat64(), 10.0)

	zero := NewBond().Face(1000, money.USD).CouponRate(0).Frequency(2).Periods(20).Yield(0.05)

	zeroMacaulay, err := zero.MacaulayDuration()
	require.NoError(t, err)
	assert.InDelta(t, 10.0, zeroMacaulay.InexactFloat64(), 1e-9)
}

// TestConvexityIsPositive checks convexity is positive for a plain bond and
// grows with maturity.
func TestConvexityIsPositive(t *testing.T) {
	short, err := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(4).Yield(0.05).Convexity()
	require.NoError(t, err)
	assert.Positive(t, short.InexactFloat64())

	long, err := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(60).Yield(0.05).Convexity()
	require.NoError(t, err)
	assert.Greater(t, long.InexactFloat64(), short.InexactFloat64())
}

// TestDurationConvexityApproximation checks the risk measures against what they
// are for: the second-order Taylor estimate of a price change,
//
//	ΔP/P ≈ −D_mod·Δy + ½·C·Δy²
//
// must track the actual repricing for a small yield move.
func TestDurationConvexityApproximation(t *testing.T) {
	const (
		baseYield = 0.06
		shift     = 0.0025 // 25 basis points
	)

	bond := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20)

	base, err := bond.Yield(baseYield).Price()
	require.NoError(t, err)

	shifted, err := bond.Yield(baseYield + shift).Price()
	require.NoError(t, err)

	modified, err := bond.Yield(baseYield).ModifiedDuration()
	require.NoError(t, err)

	convexity, err := bond.Yield(baseYield).Convexity()
	require.NoError(t, err)

	actual := (shifted.InexactFloat64() - base.InexactFloat64()) / base.InexactFloat64()
	durationOnly := -modified.InexactFloat64() * shift
	estimated := durationOnly + 0.5*convexity.InexactFloat64()*shift*shift

	// The residual is the omitted third-order term, a couple of parts per
	// million of the price for a 25bp move on a 10-year bond.
	assert.InDelta(t, actual, estimated, 1e-5)

	// The point of carrying convexity: it must be a real improvement over the
	// duration-only estimate, not a rounding-level correction.
	assert.Less(t, abs(actual-estimated), abs(actual-durationOnly)/10)
}

// abs is the float64 absolute value, used to compare approximation errors.
func abs(v float64) float64 {
	if v < 0 {
		return -v
	}

	return v
}

// TestRiskMeasureMustVariants covers the panicking risk helpers, which had no
// coverage at all: they return the value on valid terms and panic on invalid
// ones.
func TestRiskMeasureMustVariants(t *testing.T) {
	valid := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(20).Yield(0.05)

	assert.NotPanics(t, func() {
		assert.Positive(t, valid.MustMacaulayDuration().InexactFloat64())
		assert.Positive(t, valid.MustModifiedDuration().InexactFloat64())
		assert.Positive(t, valid.MustConvexity().InexactFloat64())
	})

	invalid := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(0).Periods(20).Yield(0.05)

	assert.Panics(t, func() { invalid.MustMacaulayDuration() })
	assert.Panics(t, func() { invalid.MustModifiedDuration() })
	assert.Panics(t, func() { invalid.MustConvexity() })
}

// TestRiskMeasuresRejectInvalidTerms checks the error paths of the risk
// measures, which share the cash-flow builder's validation.
func TestRiskMeasuresRejectInvalidTerms(t *testing.T) {
	measures := map[string]func(Config) (decimal.Decimal, error){
		"MacaulayDuration": Config.MacaulayDuration,
		"ModifiedDuration": Config.ModifiedDuration,
		"Convexity":        Config.Convexity,
	}

	invalidTerms := map[string]Config{
		"zero frequency": NewBond().Face(1000, money.USD).Frequency(0).Periods(10).Yield(0.05),
		"zero periods":   NewBond().Face(1000, money.USD).Frequency(2).Periods(0).Yield(0.05),
		"yield at bound": NewBond().Face(1000, money.USD).Frequency(2).Periods(10).Yield(-2),
	}

	for measureName, measure := range measures {
		for termName, config := range invalidTerms {
			t.Run(measureName+"/"+termName, func(t *testing.T) {
				_, err := measure(config)
				assert.Error(t, err)
			})
		}
	}
}

// TestCouponPaymentAcrossFrequencies checks the per-period coupon is the annual
// coupon split by frequency, and that an invalid frequency is reported.
func TestCouponPaymentAcrossFrequencies(t *testing.T) {
	for _, test := range []struct {
		freq     int
		expected float64
	}{{1, 60}, {2, 30}, {4, 15}, {12, 5}} {
		coupon, err := NewBond().Face(1000, money.USD).CouponRate(0.06).Frequency(test.freq).CouponPayment()
		require.NoError(t, err)
		assert.InDelta(t, test.expected, coupon.InexactFloat64(), 1e-9)
	}

	_, err := NewBond().Face(1000, money.USD).CouponRate(0.06).Frequency(0).CouponPayment()
	assert.ErrorIs(t, err, ErrInvalidFrequency)
}

// TestAccruedInterestBoundaries checks the two ends of a coupon period:
// settling on the last coupon date accrues nothing, and settling on the next
// one accrues the whole coupon.
func TestAccruedInterestBoundaries(t *testing.T) {
	last := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	next := time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)
	coupon := usd(25)

	for _, convention := range []daycount.Convention{
		daycount.Thirty360, daycount.Actual360, daycount.Actual365Fixed, daycount.ActualActualISDA,
	} {
		t.Run(convention.String(), func(t *testing.T) {
			atStart, err := AccruedInterest(coupon, last, last, next, convention)
			require.NoError(t, err)
			assert.InDelta(t, 0.0, atStart.InexactFloat64(), 1e-9)

			atEnd, err := AccruedInterest(coupon, last, next, next, convention)
			require.NoError(t, err)
			assert.InDelta(t, 25.0, atEnd.InexactFloat64(), 1e-9)

			// Halfway through, the accrual sits between the two ends.
			middle, err := AccruedInterest(coupon, last, time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC), next, convention)
			require.NoError(t, err)
			assert.Greater(t, middle.InexactFloat64(), 0.0)
			assert.Less(t, middle.InexactFloat64(), 25.0)
		})
	}
}

// TestAccruedInterestErrors covers the failure paths: a settlement before the
// last coupon, and a coupon period of zero length.
func TestAccruedInterestErrors(t *testing.T) {
	last := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	next := time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)

	_, err := AccruedInterest(usd(25), last, last.AddDate(0, 0, -1), next, daycount.Thirty360)
	assert.ErrorIs(t, err, daycount.ErrEndBeforeStart)

	_, err = AccruedInterest(usd(25), last, last, last, daycount.Thirty360)
	assert.ErrorIs(t, err, ErrInvalidPeriods)

	_, err = AccruedInterest(usd(25), last, next, next, daycount.Convention(99))
	assert.ErrorIs(t, err, daycount.ErrInvalidConvention)
}

// TestAccruedInterestPreservesCurrency checks the accrual carries the coupon's
// currency rather than defaulting.
func TestAccruedInterestPreservesCurrency(t *testing.T) {
	last := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	next := time.Date(2024, 7, 15, 0, 0, 0, 0, time.UTC)

	accrued, err := AccruedInterest(
		money.MustMoneyFromFloat64(25, money.EUR),
		last, time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC), next, daycount.Thirty360)
	require.NoError(t, err)
	assert.Equal(t, money.EUR, accrued.GetCurrency())
}

// TestSinglePeriodBond covers the shortest possible bond: one period left, so
// its price is the final coupon plus face discounted once.
func TestSinglePeriodBond(t *testing.T) {
	price, err := NewBond().Face(1000, money.USD).CouponRate(0.06).Frequency(2).Periods(1).Yield(0.06).Price()
	require.NoError(t, err)

	// (1000 + 30) / 1.03 = 1000
	assert.InDelta(t, 1000.0, price.InexactFloat64(), 1e-9)

	duration, err := NewBond().Face(1000, money.USD).CouponRate(0.06).Frequency(2).Periods(1).Yield(0.06).MacaulayDuration()
	require.NoError(t, err)
	assert.InDelta(t, 0.5, duration.InexactFloat64(), 1e-9)
}

// TestYTMOnUnreachablePrice checks the exhaustion path: a price no yield in the
// candidate range can produce must be reported, not approximated.
func TestYTMOnUnreachablePrice(t *testing.T) {
	// Far above the undiscounted sum of the cash flows, which even a −99%
	// yield cannot reach from below on this bond.
	_, err := NewBond().Face(1000, money.USD).CouponRate(0.05).Frequency(2).Periods(10).
		MarketPrice(1e12).YTM()
	assert.ErrorIs(t, err, ErrNoConvergence)
}

// TestPriceCurrencyIsPreserved checks the price carries the face value's
// currency.
func TestPriceCurrencyIsPreserved(t *testing.T) {
	price, err := NewBond().Face(1000, money.JPY).CouponRate(0.05).Frequency(2).Periods(10).Yield(0.06).Price()
	require.NoError(t, err)
	assert.Equal(t, money.JPY, price.GetCurrency())
}
