package money

func (m *Money) StringMoney() (string, error) {
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
