package compositeinterest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yeferson59/gofinance/money"
)

func TestNewPeriodSuccess(t *testing.T) {
	period, err := NewPeriod(money.MustFromFloat64(2), QuarterlyTwo)

	tx := assert.New(t)

	tx.NotNil(period, "period mustn't be nil")
	tx.Nil(err)
}

func TestGetValuePeriod(t *testing.T) {
	numberPeriod := 2.0
	period, err := NewPeriod(money.MustFromFloat64(numberPeriod), QuarterlyTwo)

	tx := assert.New(t)

	tx.NotNil(period, "period mustn't be nil")
	tx.Nil(err)

	value, _, err := period.getPeriod()

	tx.Equal(numberPeriod, value.InexactFloat64())
	tx.Nil(err)
	tx.NoError(err)
}

func TestNewPeriodZeroValueIsValid(t *testing.T) {
	// A zero period is allowed: it signals that the period should be derived.
	period, err := NewPeriod(money.MustFromFloat64(0), QuarterlyTwo)

	tx := assert.New(t)

	tx.NoError(err)

	value, _, err := period.getPeriod()
	tx.NoError(err)
	tx.Equal(0.0, value.InexactFloat64())
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
		period, err := NewPeriod(money.MustFromFloat64(tt.value), tt.compoundingFrequency)

		tx := assert.New(t)

		tx.NotNil(period, "period mustn't be nil")
		tx.Nil(err)
	}
}

func TestNewCompositeInterestZeroValue(t *testing.T) {
	presentMoney, err := money.New(0, 2, money.USD)
	tx := assert.New(t)
	tx.NoError(err)
	futureMoney, err := money.New(0, 2, money.USD)
	tx.NoError(err)
	cr, err := New(presentMoney, futureMoney, RateInterest{}, Period{})

	tx.NotNil(cr)
	tx.NoError(err)
}

func TestGetEqualsRateInterestPeriods_differentTime(t *testing.T) {
	period, err := NewPeriod(money.MustFromFloat64(2), Bimonthly)

	tx := assert.New(t)

	tx.NotNil(period)
	tx.NoError(err)

	rate, err := NewRateInterest(money.MustFromFloat64(0.05), Monthly, RateEffectyNominal)

	tx.NotNil(rate)
	tx.NoError(err)

	presentMoney, err := money.New(500000, 2, money.USD)
	tx.NoError(err)
	futureMoney, err := money.New(0, 2, money.USD)
	tx.NoError(err)
	compositeInterest, err := New(presentMoney, futureMoney, rate, period)

	tx.NotNil(compositeInterest)
	tx.NoError(err)

	future, err := compositeInterest.Future()

	tx.NoError(err)
	// 2 bimonthly periods = 4 months; 5% nominal monthly => i = 0.05/12
	// FV = 5000 × (1 + 0.05/12)^4 = 5083.8556
	tx.InDelta(5083.8556, future.ToDecimal().InexactFloat64(), 0.01)
}
