package money

func (m Money) StringMoney() (string, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return "", err
	}

	return isoCode + " " + m.StringFixed(prec), nil
}

// Format renders m using its currency's display symbol instead of its ISO
// code, e.g. "$10.00" rather than "USD 10.00".
func (m Money) Format() (string, error) {
	symbol, err := m.currency.Symbol()
	if err != nil {
		return "", err
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return "", err
	}

	return symbol + m.StringFixed(prec), nil
}
