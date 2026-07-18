package simpleinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func NewPeriodTest(t *testing.T) Period {
	numberPeriod := decimal.MustFromFloat64(2.0)
	period := NewPeriod(numberPeriod, Days)

	tx := assert.New(t)

	require.NotNil(t, period, "period should not be nil")
	periodValue, _ := period.getPeriod()
	expectedValue := decimal.MustFromFloat64(2.0)
	tx.Equal(expectedValue.String(), periodValue.String(), "days value should match input")

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
			numberDecimal, _ := decimal.NewFromFloat64(tt.number)
			period := NewPeriod(numberDecimal, tt.Periods)

			tx := assert.New(t)

			switch tt.Periods {
			case Days:
				tx.NotEqual(decimal.Decimal{}, period.days, "days should not be zero for Days period")
				periodValue, _ := period.getPeriod()
				expectedValue, _ := decimal.NewFromFloat64(tt.number)
				tx.Equal(expectedValue.String(), periodValue.String())
			case Weeks:
				tx.NotEqual(decimal.Decimal{}, period.weeks, "weeks should not be zero for Weeks period")
				periodValue, _ := period.getPeriod()
				expectedValue, _ := decimal.NewFromFloat64(tt.number)
				tx.Equal(expectedValue.String(), periodValue.String())
			case Months:
				tx.NotEqual(decimal.Decimal{}, period.months, "months should not be zero for Months period")
				periodValue, _ := period.getPeriod()
				expectedValue, _ := decimal.NewFromFloat64(tt.number)
				tx.Equal(expectedValue.String(), periodValue.String())
			case Years:
				tx.NotEqual(decimal.Decimal{}, period.years, "years should not be zero for Years period")
				periodValue, _ := period.getPeriod()
				expectedValue, _ := decimal.NewFromFloat64(tt.number)
				tx.Equal(expectedValue.String(), periodValue.String())
			case "":
				zero, _ := decimal.NewFromFloat64(0.0)
				tx.Equal(zero.String(), period.days.String(), "all fields should be zero for invalid period")
				tx.Equal(zero.String(), period.weeks.String())
				tx.Equal(zero.String(), period.months.String())
				tx.Equal(zero.String(), period.years.String())
			}
		}
	})

	t.Run("get valid value for period", func(t *testing.T) {
		for _, tt := range tests {
			numberDecimal, _ := decimal.NewFromFloat64(tt.number)
			period := NewPeriod(numberDecimal, tt.Periods)
			valuePeriod, err := period.getPeriod()

			if tt.Periods == "" {
				zero, _ := decimal.NewFromFloat64(0.0)
				assert.Equal(t, zero.String(), valuePeriod.String(), "value should be zero for invalid period")
				assert.Error(t, err, "should return error for invalid period")
			} else {
				require.NotNil(t, valuePeriod, "value should not be nil for valid period")
				expectedValue, _ := decimal.NewFromFloat64(tt.number)
				assert.Equal(t, expectedValue.String(), valuePeriod.String(), "period value should match input")
				assert.NoError(t, err, "should not return error for valid period")
			}
		}
	})
}

func TestGetPeriod(t *testing.T) {
	period := NewPeriodTest(t)
	expectedPeriod, _ := decimal.NewFromFloat64(2.0)

	numberPeriod, err := period.getPeriod()
	require.NoError(t, err, "should not return error")
	require.NotNil(t, numberPeriod, "should return a value")
	assert.Equal(t, expectedPeriod.String(), numberPeriod.String(), "period value should match expected")
}

func TestNewSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)

	zero, _ := money.New(0, 0, money.COP)
	zeroDecimal, _ := decimal.NewFromFloat64(0)
	simpleInterest := New(zero, zero, zero, zeroDecimal, period)

	tx := assert.New(t)

	tx.NotNil(simpleInterest, "simple interest should not be nil")
	tx.Equal(zero.String(), simpleInterest.future.String(), "future should be 0")
	tx.Equal(zero.String(), simpleInterest.present.String(), "present should be 0")
	tx.Equal(zero.String(), simpleInterest.interest.String(), "interest should be 0")
	tx.Equal(zeroDecimal.String(), simpleInterest.rateInterest.String(), "rateInterest should be 0")
	tx.NotNil(simpleInterest.periods, "periods should not be nil")
}

func TestGetPeriodWithSimpleInterest(t *testing.T) {
	period := NewPeriodTest(t)
	zero, _ := money.New(0, 0, money.COP)
	zeroDecimal, _ := decimal.NewFromFloat64(0)
	simpleInterest := New(zero, zero, zero, zeroDecimal, period)

	valuePeriod, err := simpleInterest.GetPeriods()

	require.NoError(t, err, "should not return error")
	expectedPeriod, _ := decimal.NewFromFloat64(2.0)
	assert.Equal(t, expectedPeriod.String(), valuePeriod.String(), "period value should equal 2.0")
}
