package compoundinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
)

// allTypeRates is every declared TypeRate. Conversions are expected to accept
// all of them, so the tests below sweep the whole set rather than sampling it:
// the silent-zero bug this file guards against lived precisely in the
// combinations no single-case test visited.
var allTypeRates = []TypeRate{
	RateEffectyPeriodic,
	RateEffectyAnnually,
	RateEffectyNominal,
	RateAnticipateEffectyPeriodic,
	RateAnticipateEffectyNominal,
}

var allFrequencies = []CompoundingFrequency{
	Daily, Monthly, Bimonthly, Quarterly, FourMonthly, SemiAnnually, Annually,
}

// conversions is every rate conversion, by name, so a test can call them all.
func conversions(rt RateInterest) map[string]func() (decimal.Decimal, error) {
	return map[string]func() (decimal.Decimal, error){
		"RatePeriodic":                  rt.RatePeriodic,
		"RateNominal":                   rt.RateNominal,
		"RateEffectyAnnually":           rt.RateEffectyAnnually,
		"RateAnticipatePeriodic":        rt.RateAnticipatePeriodic,
		"RateAnticipateNominal":         rt.RateAnticipateNominal,
		"RateAnticipateEffectyAnnually": rt.RateAnticipateEffectyAnnually,
		"ToPeriodic":                    rt.ToPeriodic,
		"ToNominal":                     rt.ToNominal,
		"ToAnticipatePeriodic":          rt.ToAnticipatePeriodic,
		"ToAnticipateNominal":           rt.ToAnticipateNominal,
	}
}

// TestRateConversionNeverSilentlyZero sweeps the full matrix of rate types,
// frequencies and conversions. A positive rate has a positive equivalent in
// every form, so returning zero with a nil error means the function fell
// through its branches without computing anything.
//
// Regression test for the defect recorded in TESTING_PLAN.md §2.2: 12 of the
// 25 type/conversion combinations returned 0 with err == nil, because the
// conversions were chains of ifs with no else or default branch.
func TestRateConversionNeverSilentlyZero(t *testing.T) {
	for _, typeRate := range allTypeRates {
		for _, frequency := range allFrequencies {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), frequency, typeRate)
			require.NoError(t, err)

			for name, convert := range conversions(rateInterest) {
				t.Run(string(typeRate)+"/"+string(frequency)+"/"+name, func(t *testing.T) {
					value, err := convert()
					require.NoError(t, err)
					assert.True(t, value.IsPos(),
						"a positive rate must convert to a positive rate, got %v", value)
				})
			}
		}
	}
}

// TestEffectiveAnnualIsIndependentOfQuotedForm takes one economic rate — 1%
// effective per month — expresses it in all five TypeRate forms, and checks
// that every form converts back to the same effective annual rate. The
// effective annual rate is a property of the money, not of how the rate was
// quoted.
func TestEffectiveAnnualIsIndependentOfQuotedForm(t *testing.T) {
	periodic := decimal.MustFromFloat64(0.01)

	base, err := NewRateInterest(periodic, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	expected, err := base.RateEffectyAnnually()
	require.NoError(t, err)

	nominal, err := base.RateNominal()
	require.NoError(t, err)

	anticipatePeriodic, err := base.RateAnticipatePeriodic()
	require.NoError(t, err)

	anticipateNominal, err := base.RateAnticipateNominal()
	require.NoError(t, err)

	sameRate := []struct {
		typeRate TypeRate
		value    decimal.Decimal
	}{
		{RateEffectyPeriodic, periodic},
		{RateEffectyAnnually, expected},
		{RateEffectyNominal, nominal},
		{RateAnticipateEffectyPeriodic, anticipatePeriodic},
		{RateAnticipateEffectyNominal, anticipateNominal},
	}

	for _, form := range sameRate {
		t.Run(string(form.typeRate), func(t *testing.T) {
			rateInterest, err := NewRateInterest(form.value, Monthly, form.typeRate)
			require.NoError(t, err)

			effectiveAnnual, err := rateInterest.RateEffectyAnnually()
			require.NoError(t, err)

			assert.InDelta(t, expected.InexactFloat64(), effectiveAnnual.InexactFloat64(), 1e-12)
		})
	}
}

// TestAnticipatedOrdinaryReciprocity checks the identity connecting the two
// rate families, d = i/(1+i), in both directions and across frequencies.
func TestAnticipatedOrdinaryReciprocity(t *testing.T) {
	for _, frequency := range allFrequencies {
		t.Run(string(frequency), func(t *testing.T) {
			ordinary, err := NewRateInterest(decimal.MustFromFloat64(0.02), frequency, RateEffectyPeriodic)
			require.NoError(t, err)

			discount, err := ordinary.RateAnticipatePeriodic()
			require.NoError(t, err)

			// d = i/(1+i) = 0.02/1.02
			assert.InDelta(t, 0.02/1.02, discount.InexactFloat64(), 1e-12)

			// Going back the other way must recover the ordinary rate.
			anticipated, err := NewRateInterest(discount, frequency, RateAnticipateEffectyPeriodic)
			require.NoError(t, err)

			recovered, err := anticipated.RatePeriodic()
			require.NoError(t, err)

			assert.InDelta(t, 0.02, recovered.InexactFloat64(), 1e-12)
		})
	}
}

// TestRateConversionRoundTrip converts a rate from each form into every other
// form and back, checking the original value is recovered.
func TestRateConversionRoundTrip(t *testing.T) {
	roundTrips := []struct {
		name     string
		from     TypeRate
		to       TypeRate
		value    float64
		forward  func(RateInterest) (decimal.Decimal, error)
		backward func(RateInterest) (decimal.Decimal, error)
	}{
		{
			"periodic <-> nominal", RateEffectyPeriodic, RateEffectyNominal, 0.01,
			RateInterest.RateNominal, RateInterest.RatePeriodic,
		},
		{
			"periodic <-> effective annual", RateEffectyPeriodic, RateEffectyAnnually, 0.01,
			RateInterest.RateEffectyAnnually, RateInterest.RatePeriodic,
		},
		{
			"periodic <-> anticipated periodic", RateEffectyPeriodic, RateAnticipateEffectyPeriodic, 0.01,
			RateInterest.RateAnticipatePeriodic, RateInterest.RatePeriodic,
		},
		{
			"periodic <-> anticipated nominal", RateEffectyPeriodic, RateAnticipateEffectyNominal, 0.01,
			RateInterest.RateAnticipateNominal, RateInterest.RatePeriodic,
		},
		{
			"effective annual <-> anticipated nominal", RateEffectyAnnually, RateAnticipateEffectyNominal, 0.1268,
			RateInterest.RateAnticipateNominal, RateInterest.RateEffectyAnnually,
		},
		{
			"nominal <-> anticipated periodic", RateEffectyNominal, RateAnticipateEffectyPeriodic, 0.12,
			RateInterest.RateAnticipatePeriodic, RateInterest.RateNominal,
		},
	}

	for _, roundTrip := range roundTrips {
		t.Run(roundTrip.name, func(t *testing.T) {
			original, err := NewRateInterest(decimal.MustFromFloat64(roundTrip.value), Monthly, roundTrip.from)
			require.NoError(t, err)

			converted, err := roundTrip.forward(original)
			require.NoError(t, err)

			target, err := NewRateInterest(converted, Monthly, roundTrip.to)
			require.NoError(t, err)

			recovered, err := roundTrip.backward(target)
			require.NoError(t, err)

			assert.InDelta(t, roundTrip.value, recovered.InexactFloat64(), 1e-12)
		})
	}
}

// TestRateConversionCrossFamilyValues pins the arithmetic of the conversions
// that used to return zero, against values computed by hand.
func TestRateConversionCrossFamilyValues(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		typeRate TypeRate
		convert  func(RateInterest) (decimal.Decimal, error)
		expected float64
	}{
		// i = 12% per month -> d = 0.12/1.12
		{"periodic to anticipated periodic", 0.12, RateEffectyPeriodic, RateInterest.RateAnticipatePeriodic, 0.12 / 1.12},
		// d × 12
		{"periodic to anticipated nominal", 0.12, RateEffectyPeriodic, RateInterest.RateAnticipateNominal, 12 * 0.12 / 1.12},
		// j = 12% nominal -> i = 1% -> d = 0.01/1.01
		{"nominal to anticipated periodic", 0.12, RateEffectyNominal, RateInterest.RateAnticipatePeriodic, 0.01 / 1.01},
		{"nominal to anticipated nominal", 0.12, RateEffectyNominal, RateInterest.RateAnticipateNominal, 12 * 0.01 / 1.01},
		// d = 12% per month -> i = 0.12/0.88
		{"anticipated periodic to periodic", 0.12, RateAnticipateEffectyPeriodic, RateInterest.RatePeriodic, 0.12 / 0.88},
		{"anticipated periodic to nominal", 0.12, RateAnticipateEffectyPeriodic, RateInterest.RateNominal, 12 * 0.12 / 0.88},
		// f = 12% anticipated nominal -> d = 1% -> i = 0.01/0.99
		{"anticipated nominal to periodic", 0.12, RateAnticipateEffectyNominal, RateInterest.RatePeriodic, 0.01 / 0.99},
		{"anticipated nominal to nominal", 0.12, RateAnticipateEffectyNominal, RateInterest.RateNominal, 12 * 0.01 / 0.99},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(test.value), Monthly, test.typeRate)
			require.NoError(t, err)

			got, err := test.convert(rateInterest)
			require.NoError(t, err)

			assert.InDelta(t, test.expected, got.InexactFloat64(), 1e-12)
		})
	}
}

// TestRateConversionInvalidTypeRate checks that an unrecognised rate type is
// reported instead of silently producing a zero.
func TestRateConversionInvalidTypeRate(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, TypeRate("not-a-rate-type"))
	require.NoError(t, err)

	for name, convert := range conversions(rateInterest) {
		t.Run(name, func(t *testing.T) {
			_, err := convert()
			assert.ErrorIs(t, err, ErrInvalidTypeRate)
		})
	}

	_, err = rateInterest.RatePeriodicToPeriodic(Quarterly)
	assert.ErrorIs(t, err, ErrInvalidTypeRate)

	_, err = rateInterest.RateNominalToNominal(Quarterly)
	assert.ErrorIs(t, err, ErrInvalidTypeRate)
}

// TestAnticipatedRateAtOrAboveOne checks the degenerate case: an anticipated
// rate that discounts a period's whole value has no finite ordinary
// equivalent, and must be reported rather than divided by zero.
func TestAnticipatedRateAtOrAboveOne(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		typeRate TypeRate
	}{
		{"100% anticipated periodic", 1, RateAnticipateEffectyPeriodic},
		{"150% anticipated periodic", 1.5, RateAnticipateEffectyPeriodic},
		// f = 12 with monthly compounding means d = 1 per period.
		{"anticipated nominal implying d = 1", 12, RateAnticipateEffectyNominal},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(test.value), Monthly, test.typeRate)
			require.NoError(t, err)

			_, err = rateInterest.RatePeriodic()
			assert.ErrorIs(t, err, ErrInvalidAnticipatedRate)

			_, err = rateInterest.RateEffectyAnnually()
			assert.ErrorIs(t, err, ErrInvalidAnticipatedRate)
		})
	}
}

// TestRateConversionZeroRate checks that a 0% rate converts to 0% in every
// form instead of erroring: a rate of zero is a legitimate input.
func TestRateConversionZeroRate(t *testing.T) {
	for _, typeRate := range allTypeRates {
		rateInterest, err := NewRateInterest(decimal.Zero, Monthly, typeRate)
		require.NoError(t, err)

		for name, convert := range conversions(rateInterest) {
			t.Run(string(typeRate)+"/"+name, func(t *testing.T) {
				value, err := convert()
				require.NoError(t, err)
				assert.True(t, value.IsZero(), "0%% must convert to 0%%, got %v", value)
			})
		}
	}
}
