package simpleinterest

import (
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestDecimalSettersMatchFloatSetters(t *testing.T) {
	want, err := NewSimple().Present(5000, money.USD).AnnualRate(0.12).Periods(18).Months().FutureValue()
	if err != nil {
		t.Fatal(err)
	}

	got, err := NewSimple().
		PresentMoney(money.MustMoneyFromString("5000", money.USD)).
		AnnualRateDecimal(decimal.MustFromString("0.12")).
		Periods(18).
		Months().
		FutureValue()
	if err != nil {
		t.Fatal(err)
	}

	if !want.Equal(got) {
		t.Errorf("float setters gave %s, decimal setters gave %s", want, got)
	}

	present, err := NewSimple().FutureMoney(money.MustMoneyFromString("5900", money.USD)).RateDecimal(decimal.MustFromString("0.01")).Periods(18).PresentValue()
	if err != nil {
		t.Fatal(err)
	}

	presentFloat, err := NewSimple().Future(5900, money.USD).Rate(0.01).Periods(18).PresentValue()
	if err != nil {
		t.Fatal(err)
	}

	if !present.Equal(presentFloat) {
		t.Errorf("RateDecimal gave %s, Rate gave %s", present, presentFloat)
	}

	interest := money.MustMoneyFromString("900", money.USD)
	if config := NewSimple().InterestMoney(interest); !config.interest.Equal(interest) {
		t.Errorf("InterestMoney stored %s, want %s", config.interest, interest)
	}
}
