package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCompoundingFrequencyDaily(t *testing.T) {
	value, err := getCompoundingFrequency(Daily)
	require.NoError(t, err)

	assert.Equal(t, 365.0, value)
}

func TestGetCompoundingFrequencyMonthly(t *testing.T) {
	value, err := getCompoundingFrequency(Monthly)
	require.NoError(t, err)

	assert.Equal(t, 12.0, value)
}

func TestGetCompoundingFrequencyBimonthly(t *testing.T) {
	value, err := getCompoundingFrequency(Bimonthly)
	require.NoError(t, err)

	assert.Equal(t, 6.0, value)
}

func TestGetCompoundingFrequencyQuarterlyOne(t *testing.T) {
	value, err := getCompoundingFrequency(QuarterlyOne)
	require.NoError(t, err)

	assert.Equal(t, 4.0, value)
}

func TestGetCompoundingFrequencyQuarterlyTwo(t *testing.T) {
	value, err := getCompoundingFrequency(QuarterlyTwo)
	require.NoError(t, err)

	assert.Equal(t, 3.0, value)
}

func TestGetCompoundingFrequencySemiAnnually(t *testing.T) {
	value, err := getCompoundingFrequency(SemiAnnually)
	require.NoError(t, err)

	assert.Equal(t, 2.0, value)
}

func TestGetCompoundingFrequencyAnnually(t *testing.T) {
	value, err := getCompoundingFrequency(Annually)
	require.NoError(t, err)

	assert.Equal(t, 1.0, value)
}

func TestGetCompoundingFrequencyInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := getCompoundingFrequency(invalidFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetCompoundingFrequencytoMonthsInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := getCompoundingFrequencytoMonths(invalidFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetOrderTimeInvalid(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	value, err := getOrderTime(invalidFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
	assert.Equal(t, "invalid value compounding frequency", err.Error())
}

func TestGetCompoundingFrequencytoMonthsEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := getCompoundingFrequencytoMonths(emptyFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
}

func TestGetOrderTimeEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := getOrderTime(emptyFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
}

func TestGetCompoundingFrequencyEmpty(t *testing.T) {
	emptyFreq := CompoundingFrequency("")
	value, err := getCompoundingFrequency(emptyFreq)

	assert.Error(t, err)
	assert.Equal(t, 0.0, value)
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
			value, err := getCompoundingFrequency(tc.frequency)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedValue, value)
		})
	}
}

func TestGetCompoundingFrequencyConsistency(t *testing.T) {
	value1, err1 := getCompoundingFrequency(Monthly)
	value2, err2 := getCompoundingFrequency(Monthly)

	require.NoError(t, err1)
	require.NoError(t, err2)

	assert.Equal(t, value1, value2)
}

func TestGetCompoundingFrequencyMapValues(t *testing.T) {
	frequencies := []CompoundingFrequency{Daily, Monthly, Bimonthly, QuarterlyOne, QuarterlyTwo, SemiAnnually, Annually}

	for _, freq := range frequencies {
		value, err := getCompoundingFrequency(freq)
		require.NoError(t, err)
		assert.True(t, value > 0, "all valid frequencies should have positive values")
	}
}

func TestNewPeriodDaily(t *testing.T) {
	period, err := NewPeriod(365, Daily)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 365.0, value)
	assert.Equal(t, Daily, freq)
}

func TestNewPeriodMonthly(t *testing.T) {
	period, err := NewPeriod(12, Monthly)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 12.0, value)
	assert.Equal(t, Monthly, freq)
}

func TestNewPeriodQuarterly(t *testing.T) {
	period, err := NewPeriod(4, QuarterlyOne)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 4.0, value)
	assert.Equal(t, QuarterlyOne, freq)
}

func TestNewPeriodAnnually(t *testing.T) {
	period, err := NewPeriod(1, Annually)
	require.NoError(t, err)

	value, freq, err := period.getPeriod()
	require.NoError(t, err)

	assert.Equal(t, 1.0, value)
	assert.Equal(t, Annually, freq)
}

func TestNewPeriodErrorWithInvalidFrequency(t *testing.T) {
	invalidFreq := CompoundingFrequency("invalid")
	_, err := NewPeriod(12, invalidFreq)
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid compounding frequency")
}

func TestNewPeriodErrorWithEmptyPeriod(t *testing.T) {
	period := &Period{}
	_, _, err := period.getPeriod()
	assert.Error(t, err)
}

func TestNewRateInterest(t *testing.T) {
	rateInterest, err := NewRateInterest(0.05, Monthly, RateEffectyPeriodic)
	require.NoError(t, err)

	assert.Equal(t, 0.05, rateInterest.value)
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
			rateInterest, err := NewRateInterest(tc.value, tc.freq, tc.typeRate)
			require.NoError(t, err)
			assert.Equal(t, tc.value, rateInterest.value)
		})
	}
}

func TestCompoundingFrequencyMapComplete(t *testing.T) {
	expectedFrequencies := 7

	count := 0
	for _, v := range countCompoundingFrequency {
		assert.True(t, v > 0)
		count++
	}

	assert.Equal(t, expectedFrequencies, count)
}

func TestCompoundingFrequencyValues(t *testing.T) {
	testCases := map[CompoundingFrequency]float64{
		Daily:        365.0,
		Monthly:      12.0,
		Bimonthly:    6.0,
		QuarterlyOne: 4.0,
		QuarterlyTwo: 3.0,
		SemiAnnually: 2.0,
		Annually:     1.0,
	}

	for freq, expectedValue := range testCases {
		value, err := getCompoundingFrequency(freq)
		require.NoError(t, err)
		assert.Equal(t, expectedValue, value)
	}
}
