package money

import "strings"

func (m *Money) String() (string, error) {
	var builder strings.Builder

	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	builder.Grow(20)
	builder.WriteString(m.value.String())
	builder.WriteString(" ")
	builder.WriteString(isoCode)
	return builder.String(), nil
}
