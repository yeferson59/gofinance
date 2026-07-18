package compoundinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
)

func TestRatePeriodicFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	assert.InDelta(t, 0.01, ratePeriodic.InexactFloat64(), 0.0001)
}

func TestRatePeriodicFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	// r = 1.1268^(1/12) - 1 = 0.0099981
	assert.InDelta(t, 0.0099981, ratePeriodic.InexactFloat64(), 0.0000001)
}

func TestRatePeriodicAlreadyPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	assert.Equal(t, 0.01, ratePeriodic.InexactFloat64())
}

func TestRateNominalFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	// nominal = 12 × (1.1268^(1/12) - 1) = 0.1199776
	assert.InDelta(t, 0.1199776, rateNominal.InexactFloat64(), 0.0000001)
}

func TestRateNominalFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.0001)
}

func TestRateNominalAlreadyNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.Equal(t, 0.12, rateNominal.InexactFloat64())
}

func TestRateEffectyAnnuallyFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	// EA = 1.01^12 - 1 = 0.1268250
	assert.InDelta(t, 0.1268250, rateAnnually.InexactFloat64(), 0.0000001)
}

func TestRateEffectyAnnuallyFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	// EA = (1 + 0.12/12)^12 - 1 = 0.1268250
	assert.InDelta(t, 0.1268250, rateAnnually.InexactFloat64(), 0.0000001)
}

func TestRateEffectyAnnuallyAlreadyAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	assert.Equal(t, 0.1268, rateAnnually.InexactFloat64())
}

func TestRatePeriodicToPeriodicMonthlyToQuarterly(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(QuarterlyOne)
	require.NoError(t, err)

	// r = 1.01^3 - 1 = 0.030301
	assert.InDelta(t, 0.030301, newRatePeriodic.InexactFloat64(), 0.0000001)
}

func TestRatePeriodicToPeriodicMonthlyToDaily(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(Daily)
	require.NoError(t, err)

	// r = 1.01^(12/365) - 1 = 0.0003272
	assert.InDelta(t, 0.0003272, newRatePeriodic.InexactFloat64(), 0.0000001)
}

func TestRatePeriodicToPeriodicSameFrequency(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(Monthly)
	require.NoError(t, err)

	assert.InDelta(t, 0.01, newRatePeriodic.InexactFloat64(), 0.0001)
}

func TestRateNominalToNominalMonthlyToQuarterly(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	newRateNominal, err := rateInterest.RateNominalToNominal(QuarterlyOne)
	require.NoError(t, err)

	// quarterly periodic = 1.01^3 - 1 = 0.030301, nominal = 0.030301 × 4 = 0.121204
	assert.InDelta(t, 0.121204, newRateNominal.InexactFloat64(), 0.0000001)
}

func TestRateNominalToNominalSameFrequency(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	newRateNominal, err := rateInterest.RateNominalToNominal(Monthly)
	require.NoError(t, err)

	assert.InDelta(t, 0.12, newRateNominal.InexactFloat64(), 0.0001)
}

func TestRateAnticipateEffectyAnnuallyFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	// anticipated d = 0.12/12 = 0.01, EA = (1 - 0.01)^-12 - 1 = 0.1281781
	assert.InDelta(t, 0.1281781, rateAnnually.InexactFloat64(), 0.0000001)
}

func TestRateAnticipateEffectyAnnuallyFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	// EA = (1 - 0.01)^-12 - 1 = 0.1281781
	assert.InDelta(t, 0.1281781, rateAnnually.InexactFloat64(), 0.0000001)
}

func TestRateAnticipateEffectyAnnuallyAlreadyAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	assert.Equal(t, 0.1268, rateAnnually.InexactFloat64())
}

func TestRateAnticipateNominalFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	// d = 12 × (1 - 1.1268^(-1/12)) = 0.1187899
	assert.InDelta(t, 0.1187899, rateNominal.InexactFloat64(), 0.0000001)
}

func TestRateAnticipateNominalFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	// d = 0.01 × 12 = 0.12
	assert.Equal(t, 0.12, rateNominal.InexactFloat64())
}

func TestRateAnticipateNominalAlreadyNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	assert.Equal(t, 0.12, rateNominal.InexactFloat64())
}

func TestRateAnticipatePeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RateAnticipatePeriodic()
	require.NoError(t, err)

	// d = 0.12 / 12 = 0.01
	assert.Equal(t, 0.01, ratePeriodic.InexactFloat64())
}

func TestRateAnticipatePeriodicAlreadyPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RateAnticipatePeriodic()
	require.NoError(t, err)

	assert.Equal(t, 0.01, ratePeriodic.InexactFloat64())
}

func TestConversionChainPeriodicToNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.0001)
}

func TestConversionChainNominalToAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	// EA = (1 + 0.12/12)^12 - 1 = 0.1268250
	assert.InDelta(t, 0.1268250, rateAnnually.InexactFloat64(), 0.0000001)
}

func TestToAnticipateNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.ToAnticipateNominal()
	require.NoError(t, err)

	// d = 12 × (1 - 1.1268^(-1/12)) = 0.1187899
	assert.InDelta(t, 0.1187899, rateNominal.InexactFloat64(), 0.0000001)
}

func TestToAnticipatePeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.ToAnticipatePeriodic()
	require.NoError(t, err)

	// d = 1 - 1.1268^(-1/12) = 0.0098992
	assert.InDelta(t, 0.0098992, ratePeriodic.InexactFloat64(), 0.0000001)
}

func TestToNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.ToNominal()
	require.NoError(t, err)

	// nominal = 12 × (1.1268^(1/12) - 1) = 0.1199776
	assert.InDelta(t, 0.1199776, rateNominal.InexactFloat64(), 0.0000001)
}

func TestToPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.ToPeriodic()
	require.NoError(t, err)

	// r = 1.1268^(1/12) - 1 = 0.0099981
	assert.InDelta(t, 0.0099981, ratePeriodic.InexactFloat64(), 0.0000001)
}

func TestRateConversionWithDifferentFrequencies(t *testing.T) {
	testCases := []struct {
		name       string
		rate       float64
		freq       CompoundingFrequency
		typeRate   TypeRate
		targetFreq CompoundingFrequency
		expected   float64
	}{
		// r = 1.01^(12/365) - 1
		{"monthly to daily", 0.01, Monthly, RateEffectyPeriodic, Daily, 0.0003272},
		// r = 1.03^(4/12) - 1
		{"quarterly to monthly", 0.03, QuarterlyOne, RateEffectyPeriodic, Monthly, 0.0099016},
		// r = 1.10^(1/2) - 1
		{"annual to semiannual", 0.10, Annually, RateEffectyPeriodic, SemiAnnually, 0.0488088},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(tc.rate), tc.freq, tc.typeRate)
			require.NoError(t, err)

			newRate, err := rateInterest.RatePeriodicToPeriodic(tc.targetFreq)
			require.NoError(t, err)

			assert.InDelta(t, tc.expected, newRate.InexactFloat64(), 0.0000001)
		})
	}
}

func TestRateConversionPreservesValue(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	rateInterest2, err := NewRateInterest(rateAnnually, Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest2.RatePeriodic()
	require.NoError(t, err)

	assert.InDelta(t, 0.01, ratePeriodic.InexactFloat64(), 0.0001)
}

func TestRateConversionConsistency(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateQuarterly, err := rateInterest.RateNominalToNominal(QuarterlyOne)
	require.NoError(t, err)

	rateInterest2, err := NewRateInterest(rateQuarterly, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	rateMonthly, err := rateInterest2.RateNominalToNominal(Monthly)
	require.NoError(t, err)

	// Converting monthly -> quarterly -> monthly must recover the original rate
	assert.InDelta(t, 0.12, rateMonthly.InexactFloat64(), 0.0000001)
}
