package money

import "github.com/quagmt/udecimal"

type Decimal struct {
	udecimal.Decimal
}

func NewFromFloat64(f float64) (Decimal, error) {
	decimal, err := udecimal.NewFromFloat64(f)
	return Decimal{
		decimal,
	}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	decimal, err := udecimal.NewFromInt64(coef, prec)
	return Decimal{
		decimal,
	}, err
}

func MustFromFloat64(f float64) Decimal {
	decimal := udecimal.MustFromFloat64(f)
	return Decimal{
		decimal,
	}
}

func MustFromInt64(coef int64, prec uint8) Decimal {
	decimal := udecimal.MustFromInt64(coef, prec)
	return Decimal{
		decimal,
	}
}
