package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPeriodSuccess(t *testing.T) {
	period, err := NewPeriod(2, QuarterlyTwo)

	tx := assert.New(t)

	tx.NotNil(period, "period mustn't be nil")
	tx.Nil(err)
}

func TestGetValuePeriod(t *testing.T) {
	numberPeriod := 2.0
	period, err := NewPeriod(numberPeriod, QuarterlyTwo)

	tx := assert.New(t)

	tx.NotNil(period, "period mustn't be nil")
	tx.Nil(err)

	value, _, err := period.getPeriod()

	tx.Equal(numberPeriod, value)
	tx.Nil(err)
	tx.NoError(err)
}

func TestNewPeriodFailed(t *testing.T) {
	period, err := NewPeriod(0, QuarterlyTwo)

	tx := assert.New(t)

	tx.NotNil(period, "period must be not nil")
	tx.Nil(err)
}

func TestNewPeriodWithDifferentsValues(t *testing.T) {
	dataTest := []struct {
		value                float64
		compoundingFrequency CompoundingFrequency
	}{
		{
			value:                2,
			compoundingFrequency: Daily,
		},
		{
			value:                5,
			compoundingFrequency: Monthly,
		},
		{
			value:                9,
			compoundingFrequency: Bimonthly,
		},
		{
			value:                1,
			compoundingFrequency: QuarterlyOne,
		},
		{
			value:                3,
			compoundingFrequency: QuarterlyTwo,
		},
		{
			value:                11,
			compoundingFrequency: SemiAnnually,
		},
		{
			value:                0.5,
			compoundingFrequency: Annually,
		},
	}

	for _, tt := range dataTest {
		period, err := NewPeriod(tt.value, tt.compoundingFrequency)

		tx := assert.New(t)

		tx.NotNil(period, "period mustn't be nil")
		tx.Nil(err)
	}
}

func TestNewCompositeInterestZeroValue(t *testing.T) {
	cr, err := New(0, 0, RateInterest{}, Period{})

	tx := assert.New(t)

	tx.NotNil(cr)
	tx.NoError(err)
}

func TestGetEqualsRateInterestPeriods_differentTime(t *testing.T) {
	period, err := NewPeriod(2, Bimonthly)

	tx := assert.New(t)

	tx.NotNil(period)
	tx.NoError(err)

	rate, err := NewRateInterest(0.05, Monthly, RateEffectyNominal)

	tx.NotNil(rate)
	tx.NoError(err)

	compositeInterest, err := New(5_000, 0, rate, period)

	tx.NotNil(compositeInterest)
	tx.NoError(err)

	future, err := compositeInterest.Future()

	tx.NotNil(future)
	tx.NoError(err)
}
