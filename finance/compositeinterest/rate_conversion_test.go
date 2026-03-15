package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/money"
)

func TestRatePeriodicFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	assert.InDelta(t, 0.01, ratePeriodic.InexactFloat64(), 0.0001)
}

func TestRatePeriodicFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	assert.True(t, ratePeriodic.InexactFloat64() > 0)
	assert.InDelta(t, 0.01, ratePeriodic.InexactFloat64(), 0.0001)
	assert.True(t, ratePeriodic.InexactFloat64() < 0.02)
}

func TestRatePeriodicAlreadyPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RatePeriodic()
	require.NoError(t, err)

	assert.Equal(t, 0.01, ratePeriodic.InexactFloat64())
}

func TestRateNominalFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.001)
}

func TestRateNominalFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.0001)
}

func TestRateNominalAlreadyNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.Equal(t, 0.12, rateNominal.InexactFloat64())
}

func TestRateEffectyAnnuallyFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	assert.InDelta(t, 0.1268, rateAnnually.InexactFloat64(), 0.001)
}

func TestRateEffectyAnnuallyFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	assert.InDelta(t, 0.1268, rateAnnually.InexactFloat64(), 0.001)
}

func TestRateEffectyAnnuallyAlreadyAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	assert.Equal(t, 0.1268, rateAnnually.InexactFloat64())
}

func TestRatePeriodicToPeriodicMonthlyToQuarterly(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(QuarterlyOne)
	require.NoError(t, err)

	assert.True(t, newRatePeriodic.InexactFloat64() > 0.01)
	assert.True(t, newRatePeriodic.InexactFloat64() < 0.035)
}

func TestRatePeriodicToPeriodicMonthlyToDaily(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(Daily)
	require.NoError(t, err)

	assert.True(t, newRatePeriodic.InexactFloat64() < 0.01)
	assert.True(t, newRatePeriodic.InexactFloat64() > 0)
}

func TestRatePeriodicToPeriodicSameFrequency(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	newRatePeriodic, err := rateInterest.RatePeriodicToPeriodic(Monthly)
	require.NoError(t, err)

	assert.InDelta(t, 0.01, newRatePeriodic.InexactFloat64(), 0.0001)
}

func TestRateNominalToNominalMonthlyToQuarterly(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	newRateNominal, err := rateInterest.RateNominalToNominal(QuarterlyOne)
	require.NoError(t, err)

	assert.True(t, newRateNominal.InexactFloat64() > 0)
}

func TestRateNominalToNominalSameFrequency(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	newRateNominal, err := rateInterest.RateNominalToNominal(Monthly)
	require.NoError(t, err)

	assert.InDelta(t, 0.12, newRateNominal.InexactFloat64(), 0.0001)
}

func TestRateAnticipateEffectyAnnuallyFromNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	assert.True(t, rateAnnually.InexactFloat64() > 0)
}

func TestRateAnticipateEffectyAnnuallyFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	assert.True(t, rateAnnually.InexactFloat64() > 0)
}

func TestRateAnticipateEffectyAnnuallyAlreadyAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateAnticipateEffectyAnnually()
	require.NoError(t, err)

	assert.Equal(t, 0.1268, rateAnnually.InexactFloat64())
}

func TestRateAnticipateNominalFromAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	assert.True(t, rateNominal.InexactFloat64() > 0)
}

func TestRateAnticipateNominalFromPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.001)
}

func TestRateAnticipateNominalAlreadyNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateAnticipateNominal()
	require.NoError(t, err)

	assert.Equal(t, 0.12, rateNominal.InexactFloat64())
}

func TestRateAnticipatePeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateAnticipateEffectyNominal)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RateAnticipatePeriodic()
	require.NoError(t, err)

	assert.True(t, ratePeriodic.InexactFloat64() > 0)
	assert.True(t, ratePeriodic.InexactFloat64() < 0.02)
}

func TestRateAnticipatePeriodicAlreadyPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateAnticipateEffectyPeriodic)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.RateAnticipatePeriodic()
	require.NoError(t, err)

	assert.Equal(t, 0.01, ratePeriodic.InexactFloat64())
}

func TestConversionChainPeriodicToNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	rateNominal, err := rateInterest.RateNominal()
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateNominal.InexactFloat64(), 0.0001)
}

func TestConversionChainNominalToAnnually(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateAnnually, err := rateInterest.RateEffectyAnnually()
	require.NoError(t, err)

	assert.InDelta(t, 0.1268, rateAnnually.InexactFloat64(), 0.001)
}

func TestToAnticipateNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.ToAnticipateNominal()
	require.NoError(t, err)

	assert.True(t, rateNominal.InexactFloat64() > 0)
}

func TestToAnticipatePeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.ToAnticipatePeriodic()
	require.NoError(t, err)

	assert.True(t, ratePeriodic.InexactFloat64() > 0)
	assert.True(t, ratePeriodic.InexactFloat64() < 0.02)
}

func TestToNominal(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	rateNominal, err := rateInterest.ToNominal()
	require.NoError(t, err)

	assert.True(t, rateNominal.InexactFloat64() > 0)
}

func TestToPeriodic(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.1268), Monthly, RateEffectyAnnually)
	require.NoError(t, err)

	ratePeriodic, err := rateInterest.ToPeriodic()
	require.NoError(t, err)

	assert.True(t, ratePeriodic.InexactFloat64() > 0)
}

func TestRateConversionWithDifferentFrequencies(t *testing.T) {
	testCases := []struct {
		name       string
		rate       float64
		freq       CompoundingFrequency
		typeRate   TypeRate
		targetFreq CompoundingFrequency
	}{
		{"monthly to daily", 0.01, Monthly, RateEffectyPeriodic, Daily},
		{"quarterly to monthly", 0.03, QuarterlyOne, RateEffectyPeriodic, Monthly},
		{"annual to semiannual", 0.10, Annually, RateEffectyPeriodic, SemiAnnually},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(money.MustFromFloat64(tc.rate), tc.freq, tc.typeRate)
			require.NoError(t, err)

			newRate, err := rateInterest.RatePeriodicToPeriodic(tc.targetFreq)
			require.NoError(t, err)

			assert.True(t, newRate.InexactFloat64() > 0)
		})
	}
}

func TestRateConversionPreservesValue(t *testing.T) {
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.01), Monthly, RateEffectyPeriodic)
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
	rateInterest, err := NewRateInterest(money.MustFromFloat64(0.12), Monthly, RateEffectyNominal)
	require.NoError(t, err)

	rateQuarterly, err := rateInterest.RateNominalToNominal(QuarterlyOne)
	require.NoError(t, err)

	rateInterest2, err := NewRateInterest(rateQuarterly, QuarterlyOne, RateEffectyNominal)
	require.NoError(t, err)

	rateMonthly, err := rateInterest2.RateNominalToNominal(Monthly)
	require.NoError(t, err)

	assert.InDelta(t, 0.12, rateMonthly.InexactFloat64(), 0.001)
}
