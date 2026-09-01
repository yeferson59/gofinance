package money

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ json.Marshaler       = Currency(0)
	_ json.MarshalerTo     = Currency(0)
	_ json.Unmarshaler     = (*Currency)(nil)
	_ json.UnmarshalerFrom = (*Currency)(nil)

	_ json.Marshaler       = Money{}
	_ json.MarshalerTo     = Money{}
	_ json.Unmarshaler     = (*Money)(nil)
	_ json.UnmarshalerFrom = (*Money)(nil)
)

// TestCurrencyJSONWireFormat checks the encoding has not changed: the ISO code
// in a JSON string, from both entry points.
func TestCurrencyJSONWireFormat(t *testing.T) {
	for _, currency := range []Currency{USD, EUR, JPY, XXX} {
		v1, err := currency.MarshalJSON()
		require.NoError(t, err)

		v2, err := json.Marshal(currency)
		require.NoError(t, err)

		assert.Equal(t, `"`+currency.String()+`"`, string(v1))
		assert.Equal(t, string(v1), string(v2), "the two entry points must agree")
	}
}

// TestCurrencyMarshalJSONUnknown checks an unrecognised currency is now an
// error. It used to encode as "Currency(255)", a document no decoder could
// read back.
func TestCurrencyMarshalJSONUnknown(t *testing.T) {
	_, err := Currency(255).MarshalJSON()
	assert.Error(t, err)

	_, err = json.Marshal(Currency(255))
	assert.Error(t, err)
}

// TestCurrencyUnmarshalJSONForms covers both accepted forms and checks the v1
// and v2 paths agree on each.
func TestCurrencyUnmarshalJSONForms(t *testing.T) {
	cases := map[string]Currency{
		`"USD"`:   USD,
		`"usd"`:   USD,
		`" EUR "`: EUR,
		`"XXX"`:   XXX,
		`1`:       XTS,
		`143`:     USD,
	}

	for input, want := range cases {
		var v1 Currency
		require.NoError(t, v1.UnmarshalJSON([]byte(input)), input)
		assert.Equal(t, want, v1, input)

		var v2 Currency
		require.NoError(t, json.Unmarshal([]byte(input), &v2), input)
		assert.Equal(t, want, v2, input)
	}
}

// TestCurrencyUnmarshalJSONRejects covers inputs the decoder used to accept
// silently. A kind it had no case for left the receiver as XXX with no error,
// and a number wider than the type was truncated into range: 256 decoded as
// currency 0, which is XXX again.
func TestCurrencyUnmarshalJSONRejects(t *testing.T) {
	cases := map[string]string{
		"boolean":          "true",
		"null":             "null",
		"array":            "[1]",
		"object":           `{"code":"USD"}`,
		"out of range":     "256",
		"far out of range": "99999",
		"negative":         "-1",
		"fractional":       "1.5",
		"unknown code":     `"ZZZ"`,
		"wrong length":     `"US"`,
		"empty string":     `""`,
		"trailing content": `"USD" "EUR"`,
		"trailing garbage": `"USD"x`,
		"two numbers":      `1 1`,
		"malformed":        "not-json",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			v1 := EUR
			assert.Error(t, v1.UnmarshalJSON([]byte(input)))
			assert.Equal(t, EUR, v1, "the receiver must not change on failure")

			v2 := EUR
			assert.Error(t, json.Unmarshal([]byte(input), &v2))
		})
	}
}

// TestCurrencyAsJSONMapKey checks the streaming encoder did not cost the map
// key support the string form already gave.
func TestCurrencyAsJSONMapKey(t *testing.T) {
	encoded, err := json.Marshal(map[Currency]int{USD: 1})
	require.NoError(t, err)
	assert.Equal(t, `{"USD":1}`, string(encoded))

	var decoded map[Currency]int
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, map[Currency]int{USD: 1}, decoded)
}

// TestMoneyJSONWireFormat checks the encoding has not changed: the value as a
// JSON string so it survives readers that parse numbers as float64, plus the
// ISO code.
func TestMoneyJSONWireFormat(t *testing.T) {
	amount := MustMoneyFromString("1234.56", EUR)

	v1, err := amount.MarshalJSON()
	require.NoError(t, err)
	assert.Equal(t, `{"value":"1234.56","currency":"EUR"}`, string(v1))

	v2, err := json.Marshal(amount)
	require.NoError(t, err)
	assert.Equal(t, string(v1), string(v2), "the two entry points must agree")
}

// TestMoneyMarshalJSONUnsetCurrency checks an amount whose currency cannot be
// named is refused rather than encoded into a document nothing can read.
func TestMoneyMarshalJSONUnsetCurrency(t *testing.T) {
	amount := Money{value: MoneyOne.value, currency: Currency(255)}

	_, err := amount.MarshalJSON()
	assert.Error(t, err)

	_, err = json.Marshal(amount)
	assert.Error(t, err)
}

// TestMoneyUnmarshalJSONForms covers every accepted shape and checks the v1
// and v2 paths agree on each.
func TestMoneyUnmarshalJSONForms(t *testing.T) {
	cases := map[string]Money{
		`{"value":"1234.56","currency":"EUR"}`: MustMoneyFromString("1234.56", EUR),
		`{"value":1234.56,"currency":"EUR"}`:   MustMoneyFromString("1234.56", EUR),
		`{"value":"1","currency":""}`:          MustMoneyFromString("1", USD),
		`{"value":"1"}`:                        MustMoneyFromString("1", USD),
		`{"value":"1","currency":143}`:         MustMoneyFromString("1", USD),
		`{"currency":"JPY","value":"1"}`:       MustMoneyFromString("1", JPY),
		`{"value":"1","note":{"a":[1,2]}}`:     MustMoneyFromString("1", USD),
		`123.45`:                               MustMoneyFromString("123.45", USD),
		`"123.45"`:                             MustMoneyFromString("123.45", USD),
	}

	for input, want := range cases {
		var v1 Money
		require.NoError(t, v1.UnmarshalJSON([]byte(input)), input)
		assert.True(t, v1.Equal(want), "%s decoded to %v %v", input, v1, v1.GetCurrency())

		var v2 Money
		require.NoError(t, json.Unmarshal([]byte(input), &v2), input)
		assert.True(t, v2.Equal(want), "%s decoded to %v %v", input, v2, v2.GetCurrency())
	}
}

// TestMoneyUnmarshalJSONRejects covers what must not decode. An object with no
// amount matters most: the currency alone describes no money, and defaulting
// it to zero would turn a truncated document into a valid one.
func TestMoneyUnmarshalJSONRejects(t *testing.T) {
	cases := map[string]string{
		"no value member":  `{"currency":"USD"}`,
		"empty object":     `{}`,
		"null value":       `{"value":null}`,
		"unknown currency": `{"value":"1","currency":"ZZZ"}`,
		"bad currency":     `{"value":"1","currency":true}`,
		"letters":          `{"value":"abc","currency":"USD"}`,
		"exponent":         `{"value":"1e2","currency":"USD"}`,
		"bare exponent":    `1e2`,
		"boolean":          `true`,
		"array":            `[1]`,
		"trailing content": `{"value":"1"} 2`,
		"duplicate member": `{"value":"1","value":"2"}`,
		"malformed":        `not-json`,
		"empty":            ``,
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			v1 := MustMoneyFromString("99", EUR)
			assert.Error(t, v1.UnmarshalJSON([]byte(input)))
			assert.True(t, v1.Equal(MustMoneyFromString("99", EUR)), "the receiver must not change on failure")

			var v2 Money
			assert.Error(t, json.Unmarshal([]byte(input), &v2))
		})
	}
}

// TestMoneyJSONRoundTripAllCurrencies checks every declared currency survives
// an encode/decode cycle with its amount.
func TestMoneyJSONRoundTripAllCurrencies(t *testing.T) {
	for currency := range currencyCode {
		amount := NewFromDecimal(MustMoneyFromString("-1234.5678", USD).GetDecimal(), currency)

		encoded, err := json.Marshal(amount)
		require.NoError(t, err)

		var decoded Money
		require.NoError(t, json.Unmarshal(encoded, &decoded))

		assert.True(t, decoded.Equal(amount), "%v round-tripped as %s", currency, encoded)
		assert.Equal(t, currency, decoded.GetCurrency())
	}
}

// TestMoneyMarshalJSONToComposes checks the streaming encoder writes into a
// surrounding document instead of replacing it, and that the amount is one
// value from the encoder's point of view.
func TestMoneyMarshalJSONToComposes(t *testing.T) {
	var out []byte

	enc := jsontext.NewEncoder(sliceWriter{&out})

	require.NoError(t, enc.WriteToken(jsontext.BeginArray))
	require.NoError(t, MustMoneyFromString("1.50", USD).MarshalJSONTo(enc))
	require.NoError(t, USD.MarshalJSONTo(enc))
	require.NoError(t, enc.WriteToken(jsontext.EndArray))

	assert.JSONEq(t, `[{"value":"1.5","currency":"USD"},"USD"]`, string(out))
}

// TestMoneyUnmarshalJSONFromReadsOneValue checks a rejected value is consumed
// whole, as the UnmarshalerFrom contract requires: stopping partway would
// strand the decoder inside it.
func TestMoneyUnmarshalJSONFromReadsOneValue(t *testing.T) {
	dec := jsontext.NewDecoder(bytesReader(`[[1,2],{"value":"3.5","currency":"JPY"}]`))

	_, err := dec.ReadToken() // '['
	require.NoError(t, err)

	var m Money
	assert.Error(t, m.UnmarshalJSONFrom(dec), "an array is not an amount")

	require.NoError(t, m.UnmarshalJSONFrom(dec))
	assert.True(t, m.Equal(MustMoneyFromString("3.5", JPY)))
	assert.Equal(t, JPY, m.GetCurrency())
}

// TestCurrencyUnmarshalJSONFromReadsOneValue checks the same for a currency:
// the value it refuses must still be consumed.
func TestCurrencyUnmarshalJSONFromReadsOneValue(t *testing.T) {
	dec := jsontext.NewDecoder(bytesReader(`[{"a":[1]},"EUR"]`))

	_, err := dec.ReadToken() // '['
	require.NoError(t, err)

	var c Currency
	assert.Error(t, c.UnmarshalJSONFrom(dec), "an object is not a currency")

	require.NoError(t, c.UnmarshalJSONFrom(dec))
	assert.Equal(t, EUR, c)
}

// TestJSONWhitespaceIsJSONsOwn checks the trailing-content check trims what
// JSON calls whitespace and nothing else. bytes.TrimSpace would also swallow a
// vertical tab or a form feed, which JSON does not allow — found by the fuzz
// target that cross-checks the two decoding paths.
func TestJSONWhitespaceIsJSONsOwn(t *testing.T) {
	accepted := []string{"\"USD\" ", "\"USD\"\t", "\"USD\"\n", "\"USD\"\r\n", " \"USD\" "}

	for _, input := range accepted {
		var c Currency
		assert.NoError(t, c.UnmarshalJSON([]byte(input)), "%q", input)
		assert.Equal(t, USD, c, "%q", input)
	}

	// A vertical tab, a form feed and a non-breaking space are all whitespace
	// to bytes.TrimSpace and content to JSON.
	rejected := []string{"\"USD\"\v", "\"USD\"\f", "\"USD\"\x00", "\"USD\"\u00a0"}

	for _, input := range rejected {
		var c Currency
		assert.Error(t, c.UnmarshalJSON([]byte(input)), "%q", input)

		var m Money
		assert.Error(t, m.UnmarshalJSON([]byte(`{"value":"1"}`+input[5:])), "%q", input)
	}
}
