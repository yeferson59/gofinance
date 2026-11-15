package annuities

import (
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

// The Money/Decimal setters configure exactly what the float64 ones do, so the
// same terms give the same result through either.
func TestDecimalSettersMatchFloatSetters(t *testing.T) {
	viaFloat, err := NewAnnuity().Present(10000, money.USD).Quarterly().AnnualRate(0.08).Periods(8).Payment()
	if err != nil {
		t.Fatal(err)
	}

	viaDecimal, err := NewAnnuity().
		PresentMoney(money.MustMoneyFromString("10000", money.USD)).
		Quarterly().
		AnnualRateDecimal(decimal.MustFromString("0.08")).
		Periods(8).
		Payment()
	if err != nil {
		t.Fatal(err)
	}

	if !viaFloat.Equal(viaDecimal) {
		t.Errorf("float setters gave %s, decimal setters gave %s", viaFloat, viaDecimal)
	}

	periodic, err := NewAnnuity().ValueMoney(money.MustMoneyFromString("500", money.USD)).RateDecimal(decimal.MustFromString("0.01")).Periods(12).FutureValue()
	if err != nil {
		t.Fatal(err)
	}

	periodicFloat, err := NewAnnuity().Value(500, money.USD).Rate(0.01).Periods(12).FutureValue()
	if err != nil {
		t.Fatal(err)
	}

	if !periodic.Equal(periodicFloat) {
		t.Errorf("RateDecimal gave %s, Rate gave %s", periodic, periodicFloat)
	}

	effective, err := NewAnnuity().FutureMoney(money.MustMoneyFromString("100000", money.COP)).EffectiveAnnualRateDecimal(decimal.MustFromString("0.2668")).Periods(12).Payment()
	if err != nil {
		t.Fatal(err)
	}

	effectiveFloat, err := NewAnnuity().Future(100000, money.COP).EffectiveAnnualRate(0.2668).Periods(12).Payment()
	if err != nil {
		t.Fatal(err)
	}

	if !effective.Equal(effectiveFloat) {
		t.Errorf("EffectiveAnnualRateDecimal gave %s, EffectiveAnnualRate gave %s", effective, effectiveFloat)
	}
}

// An amount no float64 can hold exactly keeps every digit through the Money
// setter.
func TestMoneySetterKeepsExactAmount(t *testing.T) {
	exact := money.MustMoneyFromString("12345678901234567.89", money.COP)

	config := NewAnnuity().PresentMoney(exact)
	if !config.present.Equal(exact) {
		t.Errorf("PresentMoney stored %s, want %s", config.present, exact)
	}
}
