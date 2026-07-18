package compoundinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/term"
)

func TestGetCompoundingFrequencyDaily(t *testing.T) {
	value, err := Daily.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(365.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyMonthly(t *testing.T) {
	value, err := Monthly.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(12.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyBimonthly(t *testing.T) {
	value, err := Bimonthly.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(6.0)

	assert.Equal(t, expectedValue.InexactFloat64(), value.InexactFloat64())
}

func TestGetCompoundingFrequencyQuarterly(t *testing.T) {
	value, err := Quarterly.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(4.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyFourMonthly(t *testing.T) {
	value, err := FourMonthly.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(3.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencySemiAnnually(t *testing.T) {
	value, err := SemiAnnually.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(2.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyAnnually(t *testing.T) {
	value, err := Annually.PeriodsPerYear()
	require.NoError(t, err)

	expectedValue := decimal.MustFromFloat64(1.0)

	assert.Equal(t, expectedValue, value)
}

func TestGetCompoundingFrequencyInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := invalidFreq.PeriodsPerYear()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.ErrorIs(t, err, term.ErrInvalidFrequency)
}

func TestGetCompoundingFrequencytoMonthsInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := invalidFreq.MonthsPerPeriod()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.ErrorIs(t, err, term.ErrInvalidFrequency)
}

func TestGetOrderTimeInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := getOrderTime(invalidFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetCompoundingFrequencytoMonthsEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := emptyFreq.MonthsPerPeriod()

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
}

func TestGetOrderTimeEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := getOrderTime(emptyFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value.InexactFloat64())
}

func TestGetCompoundingFrequencyEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := emptyFreq.PeriodsPerYear()

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
		{"quarterly one", Quarterly, 4.0},
		{"quarterly two", FourMonthly, 3.0},
		{"semi annually", SemiAnnually, 2.0},
		{"annually", Annually, 1.0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			value, err := tc.frequency.PeriodsPerYear()
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue, value.InexactFloat64())
		})
	}
}

func TestGetCompoundingFrequencyConsistency(t *testing.T) {
	value1, err1 := Monthly.PeriodsPerYear()
	value2, err2 := Monthly.PeriodsPerYear()

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, value1, value2)
}

func TestGetCompoundingFrequencyMapValues(t *testing.T) {
	frequencies := []CompoundingFrequency{Daily, Monthly, Bimonthly, Quarterly, FourMonthly, SemiAnnually, Annually}

	for _, freq := range frequencies {
		value, err := freq.PeriodsPerYear()
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
	period, err := NewPeriod(decimal.MustFromFloat64(4), Quarterly)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 4.0, value.InexactFloat64())
	assert.Equal(t, Quarterly, freq)
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
	for _, v := range orderTime {
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
		Quarterly:    decimal.MustFromFloat64(4.0),
		FourMonthly:  decimal.MustFromFloat64(3.0),
		SemiAnnually: decimal.MustFromFloat64(2.0),
		Annually:     decimal.MustFromFloat64(1.0),
	}

	for freq, expectedValue := range testCases {
		value, err := freq.PeriodsPerYear()
		require.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}
}
