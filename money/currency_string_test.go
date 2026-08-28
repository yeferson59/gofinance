package money

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCurrencyString covers the String method added to close the gap recorded
// in TESTING_PLAN.md §2.9: Currency is an integer type, so printing one
// produced its number rather than its ISO code.
func TestCurrencyString(t *testing.T) {
	cases := map[Currency]string{
		USD: "USD",
		EUR: "EUR",
		JPY: "JPY",
		BHD: "BHD",
		XXX: "XXX",
	}

	for currency, expected := range cases {
		assert.Equal(t, expected, currency.String())
		assert.Equal(t, expected, fmt.Sprintf("%v", currency))
		assert.Equal(t, expected, fmt.Sprint(currency))
	}
}

// TestCurrencyStringUnknown checks an unrecognised currency names the problem
// instead of printing a bare integer that means nothing to a reader.
func TestCurrencyStringUnknown(t *testing.T) {
	unknown := Currency(255)

	assert.Equal(t, "Currency(255)", unknown.String())

	// The accessor still reports the error; String is for humans reading
	// output, not a replacement for validation.
	_, err := unknown.GetCurrencyISOCode()
	assert.Error(t, err)
}

// TestCurrencyStringAgreesWithISOCode checks the two never disagree across
// every declared currency.
func TestCurrencyStringAgreesWithISOCode(t *testing.T) {
	for currency := range currencyCode {
		isoCode, err := currency.GetCurrencyISOCode()
		assert.NoError(t, err)
		assert.Equal(t, string(isoCode[:]), currency.String())
	}
}
