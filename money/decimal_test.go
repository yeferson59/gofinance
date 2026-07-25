package money

import (
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
)

func TestToDecimal(t *testing.T) {
	n := "200000"

	d := decimal.MustFromString(n)

	if n != d.String() {
		t.Errorf("%s must be equal %s", n, d.String())
	}

	m := FromDecimal(d, USD)

	if m.String() != d.String() {
		t.Errorf("%s must be equal to %s", n, d.String())
	}
}
