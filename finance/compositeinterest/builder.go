package compositeinterest

import (
	"github.com/yeferson59/gofinance/money"
)

type CompositeConfig struct {
	present   money.Money
	future    money.Money
	rate      money.Decimal
	periods   int
	frequency CompoundingFrequency
	rateType  TypeRate
}

func NewComposite() *CompositeConfig {
	return &CompositeConfig{}
}

func (c *CompositeConfig) Present(amount float64, currency money.Currency) *CompositeConfig {
	c.present = money.MustMoneyFromFloat64(amount, currency)
	return c
}

func (c *CompositeConfig) PresentMoney(m money.Money) *CompositeConfig {
	c.present = m
	return c
}

func (c *CompositeConfig) Future(amount float64, currency money.Currency) *CompositeConfig {
	c.future = money.MustMoneyFromFloat64(amount, currency)
	return c
}

func (c *CompositeConfig) FutureMoney(m money.Money) *CompositeConfig {
	c.future = m
	return c
}

func (c *CompositeConfig) Rate(rate float64) *CompositeConfig {
	c.rate = money.MustFromFloat64(rate)
	return c
}

func (c *CompositeConfig) RateMoney(r money.Decimal) *CompositeConfig {
	c.rate = r
	return c
}

func (c *CompositeConfig) Periods(n int) *CompositeConfig {
	c.periods = n
	return c
}

func (c *CompositeConfig) Frequency(f CompoundingFrequency) *CompositeConfig {
	c.frequency = f
	return c
}

func (c *CompositeConfig) RateType(t TypeRate) *CompositeConfig {
	c.rateType = t
	return c
}

func (c *CompositeConfig) Monthly() *CompositeConfig {
	c.frequency = Monthly
	return c
}

func (c *CompositeConfig) Annually() *CompositeConfig {
	c.frequency = Annually
	return c
}

func (c *CompositeConfig) Quarterly() *CompositeConfig {
	c.frequency = QuarterlyOne
	return c
}

func (c *CompositeConfig) Daily() *CompositeConfig {
	c.frequency = Daily
	return c
}

func (c *CompositeConfig) Build() (CompositeInterest, error) {
	period, err := NewPeriod(money.MustFromFloat64(float64(c.periods)), c.frequency)
	if err != nil {
		return CompositeInterest{}, err
	}

	rate, err := NewRateInterest(c.rate, c.frequency, c.rateType)
	if err != nil {
		return CompositeInterest{}, err
	}

	return New(c.present, c.future, rate, period)
}

func (c *CompositeConfig) MustBuild() CompositeInterest {
	ci, err := c.Build()
	if err != nil {
		panic(err)
	}
	return ci
}
