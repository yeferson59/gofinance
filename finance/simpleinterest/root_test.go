package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewPeriodTest(t *testing.T) *Period {
	numberPeriod := 2.0
	period := NewPeriod(numberPeriod, Days)

	require.NotNil(t, period, "period should not be nil")
	assert.Nil(t, period.weeks, "weeks should be nil for Days period")
	assert.Nil(t, period.months, "months should be nil for Days period")
	assert.Nil(t, period.years, "years should be nil for Days period")
	assert.NotNil(t, period.days, "days should not be nil for Days period")
	assert.Equal(t, numberPeriod, *period.days, "days value should match input")

	return period
}

func TestNewPeriod(t *testing.T) {
	_ = NewPeriodTest(t)
}

func TestNewPeriodWithDifferentsTimes(t *testing.T) {
	type DataTest struct {
		number  float64
		Periods Periods
	}

	tests := []DataTest{
		{number: 2, Periods: Days},
		{number: 3, Periods: Weeks},
		{number: 3, Periods: Months},
		{number: 4, Periods: Years},
		{number: 4, Periods: ""},
	}

	t.Run("valid new correct period", func(t *testing.T) {
		for _, tt := range tests {
			period := NewPeriod(tt.number, tt.Periods)

			switch tt.Periods {
			case Days:
				assert.NotNil(t, period.days, "days should not be nil for Days period")
				assert.Equal(t, tt.number, *period.days)
			case Weeks:
				assert.NotNil(t, period.weeks, "weeks should not be nil for Weeks period")
				assert.Equal(t, tt.number, *period.weeks)
			case Months:
				assert.NotNil(t, period.months, "months should not be nil for Months period")
				assert.Equal(t, tt.number, *period.months)
			case Years:
				assert.NotNil(t, period.years, "years should not be nil for Years period")
				assert.Equal(t, tt.number, *period.years)
			case "":
				assert.Nil(t, period.days, "all fields should be nil for invalid period")
				assert.Nil(t, period.weeks)
				assert.Nil(t, period.months)
				assert.Nil(t, period.years)
			}
		}
	})

	t.Run("get valid value for period", func(t *testing.T) {
		for _, tt := range tests {
			period := NewPeriod(tt.number, tt.Periods)
			valuePeriod, err := period.getPeriod()

			if tt.Periods == "" {
				assert.Nil(t, valuePeriod, "value should be nil for invalid period")
				assert.Error(t, err, "should return error for invalid period")
			} else {
				require.NotNil(t, valuePeriod, "value should not be nil for valid period")
				assert.Equal(t, tt.number, *valuePeriod, "period value should match input")
				assert.NoError(t, err, "should not return error for valid period")
			}
		}
	})
}

func TestGetPeriod(t *testing.T) {
	period := NewPeriodTest(t)
	expectedPeriod := 2.0

	numberPeriod, err := period.getPeriod()
	require.NoError(t, err, "should not return error")
	require.NotNil(t, numberPeriod, "should return a value")
	assert.Equal(t, expectedPeriod, *numberPeriod, "period value should match expected")
}

func TestNewSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)

	simpleInterest := New(0, 0, 0, 0, period)

	require.NotNil(t, simpleInterest, "simple interest should not be nil")
	assert.Equal(t, 0.0, simpleInterest.future, "future should be 0")
	assert.Equal(t, 0.0, simpleInterest.present, "present should be 0")
	assert.Equal(t, 0.0, simpleInterest.interest, "interest should be 0")
	assert.Equal(t, 0.0, simpleInterest.rateInterest, "rateInterest should be 0")
	assert.NotNil(t, simpleInterest.periods, "periods should not be nil")
}

func TestGetPeriodWithSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)
	simpleInterest := New(0, 0, 0, 0, period)

	valuePeriod, err := simpleInterest.GetPeriods()

	require.NoError(t, err, "should not return error")
	assert.Equal(t, 2.0, valuePeriod, "period value should equal 2.0")
}

