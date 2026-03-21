package annuities

import (
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

type AnnuityConfig struct {
	value     money.Money
	present   money.Money
	future    money.Money
	periods   int
	rate      float64
	frequency compositeinterest.CompoundingFrequency
	rateType  compositeinterest.TypeRate
}

func NewAnnuity() *AnnuityConfig {
	return &AnnuityConfig{
		frequency: compositeinterest.Monthly,
		rateType:  compositeinterest.RateEffectyPeriodic,
	}
}

func (a *AnnuityConfig) Value(amount float64, currency money.Currency) *AnnuityConfig {
	a.value = money.MustMoneyFromFloat64(amount, currency)
	return a
}

func (a *AnnuityConfig) Present(amount float64, currency money.Currency) *AnnuityConfig {
	a.present = money.MustMoneyFromFloat64(amount, currency)
	return a
}

func (a *AnnuityConfig) Future(amount float64, currency money.Currency) *AnnuityConfig {
	a.future = money.MustMoneyFromFloat64(amount, currency)
	return a
}

func (a *AnnuityConfig) Periods(n int) *AnnuityConfig {
	a.periods = n
	return a
}

func (a *AnnuityConfig) Years(n int) *AnnuityConfig {
	a.periods = n * 12
	return a
}

func (a *AnnuityConfig) Rate(r float64) *AnnuityConfig {
	a.rate = r
	return a
}

func (a *AnnuityConfig) AnnualRate(r float64) *AnnuityConfig {
	divisor := 12.0
	switch a.frequency {
	case compositeinterest.Daily:
		divisor = 365.0
	case compositeinterest.Bimonthly:
		divisor = 6.0
	case compositeinterest.QuarterlyOne, compositeinterest.QuarterlyTwo:
		divisor = 4.0
	case compositeinterest.SemiAnnually:
		divisor = 2.0
	case compositeinterest.Annually:
		divisor = 1.0
	}
	a.rate = r / divisor
	return a
}

func (a *AnnuityConfig) Monthly() *AnnuityConfig {
	a.frequency = compositeinterest.Monthly
	return a
}

func (a *AnnuityConfig) Annually() *AnnuityConfig {
	a.frequency = compositeinterest.Annually
	return a
}

func (a *AnnuityConfig) Quarterly() *AnnuityConfig {
	a.frequency = compositeinterest.QuarterlyOne
	return a
}

func (a *AnnuityConfig) Build() (Annuity, error) {
	rate, err := compositeinterest.NewRateInterest(
		money.MustFromFloat64(a.rate),
		a.frequency,
		a.rateType,
	)
	if err != nil {
		return Annuity{}, err
	}

	period, err := compositeinterest.NewPeriod(
		money.MustFromFloat64(float64(a.periods)),
		a.frequency,
	)
	if err != nil {
		return Annuity{}, err
	}

	return New(a.value, a.present, a.future, period, rate)
}

func (a *AnnuityConfig) MustBuild() Annuity {
	annuity, err := a.Build()
	if err != nil {
		panic(err)
	}
	return annuity
}

func (a *AnnuityConfig) Payment() (money.Money, error) {
	annuity, err := a.Build()
	if err != nil {
		return money.Money{}, err
	}
	return annuity.PaymentFromPresentValue()
}

func (a *AnnuityConfig) MustPayment() money.Money {
	m, err := a.Payment()
	if err != nil {
		panic(err)
	}
	return m
}
