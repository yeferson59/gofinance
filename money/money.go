package money

import "strings"

func (m *Money) String() (string, error) {
	var b strings.Builder

	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	b.Grow(20)
	b.WriteString(m.value.String())
	b.WriteString(" ")
	b.WriteString(isoCode)
	return b.String(), nil
}
