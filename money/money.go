package money

import "strings"

// StringMoney renders m as its ISO 4217 code and the amount at the currency's
// precision, e.g. "USD 10.00" or "JPY 1501". An amount carrying more digits is
// rounded half away from zero. Scan reads this form back.
func (m Money) StringMoney() (string, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return "", err
	}

	return string(isoCode[:]) + " " + m.StringFixed(prec), nil
}

// Format renders m using its currency's display symbol instead of its ISO
// code, e.g. "$10.00" rather than "USD 10.00", with any minus sign before the
// symbol ("-$10.00"). Like StringMoney, it shows the currency's precision and
// rounds half away from zero.
func (m Money) Format() (string, error) {
	symbol, err := m.currency.Symbol()
	if err != nil {
		return "", err
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return "", err
	}

	amount := m.StringFixed(prec)
	if digits, negative := strings.CutPrefix(amount, "-"); negative {
		return "-" + symbol + digits, nil
	}

	return symbol + amount, nil
}
