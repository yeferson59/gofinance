package money

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.TextMarshaler   = Currency(0)
	_ encoding.TextAppender    = Currency(0)
	_ encoding.TextUnmarshaler = (*Currency)(nil)
)

// TestCurrencyMarshalText checks the text form is the bare ISO code: no
// quotes, no padding, and identical to what String reports.
func TestCurrencyMarshalText(t *testing.T) {
	for _, currency := range []Currency{USD, EUR, JPY, BHD, XXX} {
		text, err := currency.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, currency.String(), string(text))
	}
}

// TestCurrencyMarshalTextAllCurrencies checks every declared currency has a
// text form and that it agrees with the ISO accessor.
func TestCurrencyMarshalTextAllCurrencies(t *testing.T) {
	for currency, isoCode := range currencyCode {
		text, err := currency.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, string(isoCode[:]), string(text))
		assert.Len(t, text, 3)
	}
}

// TestCurrencyMarshalTextUnknown checks an unrecognised currency is an error
// rather than the "Currency(<n>)" form String uses for humans, which no
// decoder could read back.
func TestCurrencyMarshalTextUnknown(t *testing.T) {
	unknown := Currency(255)

	_, err := unknown.MarshalText()
	assert.Error(t, err)
}

// TestCurrencyAppendText checks AppendText writes into a caller's buffer
// instead of allocating, and leaves that buffer untouched when it fails.
func TestCurrencyAppendText(t *testing.T) {
	buf := []byte("rate=")

	buf, err := USD.AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "rate=USD", string(buf))

	unchanged, err := Currency(255).AppendText(buf)
	assert.Error(t, err)
	assert.Equal(t, "rate=USD", string(unchanged))
}

// TestCurrencyUnmarshalTextRoundTrip checks every declared currency survives
// the encode/decode cycle unchanged.
func TestCurrencyUnmarshalTextRoundTrip(t *testing.T) {
	for currency := range currencyCode {
		text, err := currency.MarshalText()
		require.NoError(t, err)

		var decoded Currency
		require.NoError(t, decoded.UnmarshalText(text))
		assert.Equal(t, currency, decoded)
	}
}

// TestCurrencyUnmarshalTextNormalises checks decoding is lenient about case
// and surrounding space the way Scan and GetCurrencyFromISOCode are, while
// MarshalText only ever writes the canonical uppercase form.
func TestCurrencyUnmarshalTextNormalises(t *testing.T) {
	for _, input := range []string{"USD", "usd", "uSd", " USD ", "\tUSD\n"} {
		var c Currency
		require.NoError(t, c.UnmarshalText([]byte(input)), input)
		assert.Equal(t, USD, c, input)
	}
}

// TestCurrencyUnmarshalTextInvalid checks malformed input is rejected rather
// than silently decoding to XXX, which would hide a truncated document.
func TestCurrencyUnmarshalTextInvalid(t *testing.T) {
	cases := map[string][]byte{
		"empty":     []byte(""),
		"nil":       nil,
		"blank":     []byte("   "),
		"too short": []byte("US"),
		"too long":  []byte("USDD"),
		"unknown":   []byte("ZZZ"),
		"numeric":   []byte("840"),
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			c := EUR
			assert.Error(t, c.UnmarshalText(input))
			assert.Equal(t, EUR, c, "the receiver must not change on failure")
		})
	}
}

// TestCurrencyUnmarshalTextDoesNotRetainInput checks the decoder copies what
// it needs: encoders are free to reuse the slice they hand over.
func TestCurrencyUnmarshalTextDoesNotRetainInput(t *testing.T) {
	text := []byte("EUR")

	var c Currency
	require.NoError(t, c.UnmarshalText(text))

	copy(text, "XXX")
	assert.Equal(t, EUR, c)
}
