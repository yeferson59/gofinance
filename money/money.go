package money

func (m *Money) String() (string, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	prec, err := m.currency.GetCurrencyPrecisionCode()
	if err != nil {
		return "", err
	}

	return isoCode + " " + m.value.StringFixed(prec), nil
}
