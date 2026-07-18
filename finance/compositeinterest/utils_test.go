package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/decimal"
)

func TestGetCompoundingFrequencyDaily(t *testing.T) {
	value, err := Daily.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(365.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyMonthly(t *testing.T) {
	value, err := Monthly.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(12.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyBimonthly(t *testing.T) {
	value, err := Bimonthly.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(6.0)

	assert.Equal(t, expectedValue.InexactFloat64(), value.InexactFloat64())
}

func TestGetCompoundingFrequencyQuarterlyOne(t *testing.T) {
	value, err := QuarterlyOne.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(4.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyQuarterlyTwo(t *testing.T) {
	value, err := QuarterlyTwo.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(3.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencySemiAnnually(t *testing.T) {
	value, err := SemiAnnually.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(2.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyAnnually(t *testing.T) {
	value, err := Annually.getCompoundingFrequency()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(1.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := invalidFreq.getCompoundingFrequency()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetCompoundingFrequencytoMonthsInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := invalidFreq.getCompoundingFrequencytoMonths()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetOrderTimeInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := invalidFreq.getOrderTime()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetCompoundingFrequencytoMonthsEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := emptyFreq.getCompoundingFrequencytoMonths()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
}

func TestGetOrderTimeEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := emptyFreq.getOrderTime()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
}

func TestGetCompoundingFrequencyEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := emptyFreq.getCompoundingFrequency()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
}

func TestGetCompoundingFrequencyAllValidValues(t *testing.T) {
	testCases := []struct {
		name          string
		frequency     CompoundingFrequency
		expectedValue float64
	}{
		{"daily", Daily, 365.0},
		{"monthly", Monthly, 12.0},
		{"bimonthly", Bimonthly, 6.0},
		{"quarterly one", QuarterlyOne, 4.0},
		{"quarterly two", QuarterlyTwo, 3.0},
		{"semi annually", SemiAnnually, 2.0},
		{"annually", Annually, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.frequency.getCompoundingFrequency()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue, value.InexactFloat64())
		})
	}
}

func TestGetCompoundingFrequencyConsistency(t *testing.T) {
	value1, err1 := Monthly.getCompoundingFrequency()
	value2, err2 := Monthly.getCompoundingFrequency()

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, value1, value2)
}

func TestGetCompoundingFrequencyMapValues(t *testing.T) {
	frequencies := []CompoundingFrequency{Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo, SemiAnnually, Annually}

	for _, freq := range frequencies {
		value, err := freq.getCompoundingFrequency()
		require.NoError(t, err)
		assert.True(t, value.InexactFloat64() > 0, "all valid frequencies should have positive values")
	}
}

func TestNewPeriodDaily(t *testing.T) {
	period, err := NewPeriod(decimal.MustFromFloat64(365), Daily)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 365.0, value.InexactFloat64())
	assert.Equal(t, Daily, freq)
}

func TestNewPeriodMonthly(t *testing.T) {
	period, err := NewPeriod(decimal.MustFromFloat64(12), Monthly)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 12.0, value.InexactFloat64())
	assert.Equal(t, Monthly, freq)
}

func TestNewPeriodQuarterly(t *testing.T) {
	period, err := NewPeriod(decimal.MustFromFloat64(4), QuarterlyOne)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 4.0, value.InexactFloat64())
	assert.Equal(t, QuarterlyOne, freq)
}

func TestNewPeriodAnnually(t *testing.T) {
	period, err := NewPeriod(decimal.MustFromFloat64(1), Annually)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 1.0, value.InexactFloat64())
	assert.Equal(t, Annually, freq)
}

func TestNewPeriodErrorWithInvalidFrequency(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	_, err := NewPeriod(decimal.MustFromFloat64(12), invalidFreq)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid compounding frequency")
}

func TestNewPeriodErrorWithEmptyPeriod(t *testing.T) {
	period := &Period{}
	_, _, err := period.getPeriod()
	assert.Error(t, err)
}

func TestNewRateInterest(t *testing.T) {
	rateInterest, err := NewRateInterest(decimal.MustFromFloat64(0.05), Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	assert.Equal(t, 0.05, rateInterest.value.InexactFloat64())
	assert.Equal(t, Monthly, rateInterest.compoundingFrequency)
	assert.Equal(t, RateEffectyPeriodic, rateInterest.typeRate)
}

func TestNewRateInterestWithDifferentTypes(t *testing.T) {
	testCases := []struct {
		name     string
		value    float64
		freq     CompoundingFrequency
		typeRate TypeRate
	}{
		{"periodic", 0.01, Monthly, RateEffectyPeriodic},
		{"nominal", 0.12, Monthly, RateEffectyNominal},
		{"annual", 0.1268, Monthly, RateEffectyAnnually},
		{"anticipate periodic", 0.01, Monthly, RateAnticipateEffectyPeriodic},
		{"anticipate nominal", 0.12, Monthly, RateAnticipateEffectyNominal},
		{"anticipate annual", 0.1268, Monthly, RateEffectyAnnually},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rateInterest, err := NewRateInterest(decimal.MustFromFloat64(tc.value), tc.freq, tc.typeRate)
			require.NoError(t, err)
			assert.Equal(t, tc.value, rateInterest.value.InexactFloat64())
		})
	}
}

func TestCompoundingFrequencyMapComplete(t *testing.T) {
	expectedFrequencies := 7

	count := 0
	for _, v := range countCompoundingFrequency {
		assert.True(t, v.InexactFloat64() > 0.0)
		count++
	}

	assert.Equal(t, expectedFrequencies, count)
}

func TestCompoundingFrequencyValues(t *testing.T) {
	testCases := map[CompoundingFrequency]decimal.Decimal{
		Daily:        decimal.MustFromFloat64(365.0),
		Monthly:      decimal.MustFromFloat64(12.0),
		Bimonthly:    decimal.MustFromFloat64(6.0),
		QuarterlyOne: decimal.MustFromFloat64(4.0),
		QuarterlyTwo: decimal.MustFromFloat64(3.0),
		SemiAnnually: decimal.MustFromFloat64(2.0),
		Annually:     decimal.MustFromFloat64(1.0),
	}

	for freq, expectedValue := range testCases {
		value, err := freq.getCompoundingFrequency()
		require.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}
}
