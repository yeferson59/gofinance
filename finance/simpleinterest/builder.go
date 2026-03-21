package simpleinterest

import (
	"github.com/yeferson59/gofinance/money"
)

type SimpleConfig struct {
	present    money.Money
	future     money.Money
	interest   money.Money
	rate       money.Decimal
	periods    int
	periodType Periods
}

func NewSimple() *SimpleConfig {
	return &SimpleConfig{
		periodType: Months,
	}
}

func (s *SimpleConfig) Present(amount float64, currency money.Currency) *SimpleConfig {
	s.present = money.MustMoneyFromFloat64(amount, currency)
	return s
}

func (s *SimpleConfig) Future(amount float64, currency money.Currency) *SimpleConfig {
	s.future = money.MustMoneyFromFloat64(amount, currency)
	return s
}

func (s *SimpleConfig) Interest(amount float64, currency money.Currency) *SimpleConfig {
	s.interest = money.MustMoneyFromFloat64(amount, currency)
	return s
}

func (s *SimpleConfig) Rate(r float64) *SimpleConfig {
	s.rate = money.MustFromFloat64(r)
	return s
}

func (s *SimpleConfig) AnnualRate(r float64) *SimpleConfig {
	divisor := 12.0
	switch s.periodType {
	case Days:
		divisor = 365.0
	case Weeks:
		divisor = 52.0
	case Years:
		divisor = 1.0
	}
	s.rate = money.MustFromFloat64(r / divisor)
	return s
}

func (s *SimpleConfig) Periods(n int) *SimpleConfig {
	s.periods = n
	return s
}

func (s *SimpleConfig) PeriodType(p Periods) *SimpleConfig {
	s.periodType = p
	return s
}

func (s *SimpleConfig) Months() *SimpleConfig {
	s.periodType = Months
	return s
}

func (s *SimpleConfig) Years() *SimpleConfig {
	s.periodType = Years
	return s
}

func (s *SimpleConfig) Days() *SimpleConfig {
	s.periodType = Days
	return s
}

func (s *SimpleConfig) Weeks() *SimpleConfig {
	s.periodType = Weeks
	return s
}

func (s *SimpleConfig) Build() SimpleInterest {
	period := NewPeriod(money.MustFromFloat64(float64(s.periods)), s.periodType)
	return New(s.future, s.present, s.interest, s.rate, period)
}

func (s *SimpleConfig) FutureValue() (money.Money, error) {
	si := s.Build()
	return si.FutureWithRateInterest()
}
