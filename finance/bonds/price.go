package bonds

import "github.com/yeferson59/gofinance/money"

// cashflowSums discounts the bond's coupon and redemption cash flows at the
// per-period yield y and returns, in one pass, the present value (the bond's
// price) together with Σ t·PV(cashflowₜ) and Σ t(t+1)·PV(cashflowₜ). The latter
// two feed the duration and convexity measures. It validates the frequency and
// period count and returns ErrInvalidYield when 1+y is not positive.
func (b Config) cashflowSums(y money.Decimal) (price, sumT, sumTT money.Decimal, err error) {
	if b.freq < 1 {
		return money.Decimal{}, money.Decimal{}, money.Decimal{}, ErrInvalidFrequency
	}

	if b.periods < 1 {
		return money.Decimal{}, money.Decimal{}, money.Decimal{}, ErrInvalidPeriods
	}

	onePlus := money.One.Add(y)
	if !onePlus.IsPos() {
		return money.Decimal{}, money.Decimal{}, money.Decimal{}, ErrInvalidYield
	}

	coupon, err := b.face.ToDecimal().Mul(b.couponRate).Div(money.MustFromInt64(int64(b.freq), 0))
	if err != nil {
		return money.Decimal{}, money.Decimal{}, money.Decimal{}, err
	}

	face := b.face.ToDecimal()

	price = money.Zero
	sumT = money.Zero
	sumTT = money.Zero
	factor := onePlus // (1+y)^t for the current t, starting at t = 1

	for t := 1; t <= b.periods; t++ {
		cashflow := coupon
		if t == b.periods {
			cashflow = cashflow.Add(face)
		}

		pv, err := cashflow.Div(factor)
		if err != nil {
			return money.Decimal{}, money.Decimal{}, money.Decimal{}, err
		}

		tDec := money.MustFromInt64(int64(t), 0)

		price = price.Add(pv)
		sumT = sumT.Add(pv.Mul(tDec))
		sumTT = sumTT.Add(pv.Mul(tDec.Mul(tDec.Add(money.One))))

		factor = factor.Mul(onePlus)
	}

	return price, sumT, sumTT, nil
}

// periodicYield returns the per-period yield, yield / frequency.
func (b Config) periodicYield() (money.Decimal, error) {
	if b.freq < 1 {
		return money.Decimal{}, ErrInvalidFrequency
	}

	return b.yield.Div(money.MustFromInt64(int64(b.freq), 0))
}

// Price returns the bond's clean price: the present value of its coupons and
// redemption discounted at the configured yield.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, or ErrInvalidYield on
// invalid terms.
func (b Config) Price() (money.Money, error) {
	y, err := b.periodicYield()
	if err != nil {
		return money.Money{}, err
	}

	price, _, _, err := b.cashflowSums(y)
	if err != nil {
		return money.Money{}, err
	}

	return price.ToMoney(b.face.Currency()), nil
}

// MustPrice is like Price but panics on error.
func (b Config) MustPrice() money.Money {
	m, err := b.Price()
	if err != nil {
		panic(err)
	}

	return m
}
