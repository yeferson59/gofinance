package money

import (
	"bytes"
	"encoding"
	"encoding/gob"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.BinaryMarshaler   = Currency(0)
	_ encoding.BinaryAppender    = Currency(0)
	_ encoding.BinaryUnmarshaler = (*Currency)(nil)

	_ encoding.BinaryMarshaler   = Money{}
	_ encoding.BinaryAppender    = Money{}
	_ encoding.BinaryUnmarshaler = (*Money)(nil)
)

// TestCurrencyBinaryLayout pins the encoding byte for byte. A binary format is
// persisted, so a change here is a change to data already written.
func TestCurrencyBinaryLayout(t *testing.T) {
	cases := map[Currency]string{
		USD: "01555344", // version, 'U', 'S', 'D'
		EUR: "01455552",
		JPY: "014a5059",
		XXX: "01585858",
	}

	for currency, want := range cases {
		encoded, err := currency.MarshalBinary()
		require.NoError(t, err)
		assert.Equal(t, want, hex.EncodeToString(encoded), currency)
		assert.Len(t, encoded, currencyBinaryLen)
	}
}

// TestMoneyBinaryLayout pins the amount's encoding: version, ISO code, then
// the decimal package's own bytes.
func TestMoneyBinaryLayout(t *testing.T) {
	encoded, err := MustMoneyFromString("99.95", JPY).MarshalBinary()
	require.NoError(t, err)

	assert.Equal(t, "014a505901020000000000000000000000000000270b", hex.EncodeToString(encoded))
	assert.Len(t, encoded, moneyBinaryPrefixLen+18)
}

// TestCurrencyBinaryRoundTripAllCurrencies checks every declared currency
// survives, which is the whole reason the code is written instead of the
// integer behind the constant.
func TestCurrencyBinaryRoundTripAllCurrencies(t *testing.T) {
	for currency := range currencyCode {
		encoded, err := currency.MarshalBinary()
		require.NoError(t, err)

		var decoded Currency
		require.NoError(t, decoded.UnmarshalBinary(encoded))
		assert.Equal(t, currency, decoded)
	}
}

// TestMoneyBinaryRoundTrip checks the amount and its currency both survive.
func TestMoneyBinaryRoundTrip(t *testing.T) {
	amounts := []Money{
		{},
		MoneyZero,
		MoneyOne,
		MustMoneyFromString("-1234.5678", EUR),
		MustMoneyFromString("0.0000000000000000001", BHD),
		MustMoneyFromString("99999999999999999999999999999999999999", JPY),
	}

	for _, amount := range amounts {
		encoded, err := amount.MarshalBinary()
		require.NoError(t, err, amount)

		var decoded Money
		require.NoError(t, decoded.UnmarshalBinary(encoded), amount)

		assert.True(t, decoded.Equal(amount), "%v became %v", amount, decoded)
		assert.Equal(t, amount.GetCurrency(), decoded.GetCurrency())
	}
}

// TestBinaryMarshalUnknownCurrency checks an amount whose currency cannot be
// named is refused, as it is by every other encoder here.
func TestBinaryMarshalUnknownCurrency(t *testing.T) {
	_, err := Currency(255).MarshalBinary()
	assert.Error(t, err)

	_, err = Money{currency: Currency(255)}.MarshalBinary()
	assert.Error(t, err)

	// A failure must not truncate a buffer the caller had already filled.
	buf := []byte{0xff}

	buf, err = Currency(255).AppendBinary(buf)
	assert.Error(t, err)
	assert.Equal(t, []byte{0xff}, buf)
}

// TestBinaryUnmarshalRejects checks nothing the layout does not define is
// waved through.
func TestBinaryUnmarshalRejects(t *testing.T) {
	currency, err := USD.MarshalBinary()
	require.NoError(t, err)

	amount, err := MustMoneyFromString("1.5", USD).MarshalBinary()
	require.NoError(t, err)

	craft := func(base []byte, mutate func([]byte)) []byte {
		out := append([]byte(nil), base...)
		mutate(out)

		return out
	}

	t.Run("currency", func(t *testing.T) {
		cases := map[string][]byte{
			"empty":           {},
			"nil":             nil,
			"truncated":       currency[:3],
			"padded":          append(append([]byte(nil), currency...), 0),
			"unknown version": craft(currency, func(b []byte) { b[0] = 2 }),
			"unknown code":    craft(currency, func(b []byte) { copy(b[1:], "ZZZ") }),
			"lowercase code":  craft(currency, func(b []byte) { copy(b[1:], "usd") }),
		}

		for name, input := range cases {
			t.Run(name, func(t *testing.T) {
				c := EUR
				assert.Error(t, c.UnmarshalBinary(input))
				assert.Equal(t, EUR, c, "the receiver must not change on failure")
			})
		}
	})

	t.Run("money", func(t *testing.T) {
		cases := map[string][]byte{
			"empty":            {},
			"nil":              nil,
			"prefix only":      amount[:moneyBinaryPrefixLen],
			"truncated amount": amount[:len(amount)-1],
			"padded":           append(append([]byte(nil), amount...), 0),
			"unknown version":  craft(amount, func(b []byte) { b[0] = 2 }),
			"unknown code":     craft(amount, func(b []byte) { copy(b[1:], "ZZZ") }),
			"bad decimal":      craft(amount, func(b []byte) { b[moneyBinaryPrefixLen] = 2 }),
		}

		for name, input := range cases {
			t.Run(name, func(t *testing.T) {
				m := MustMoneyFromString("99", EUR)
				assert.Error(t, m.UnmarshalBinary(input))
				assert.True(t, m.Equal(MustMoneyFromString("99", EUR)), "the receiver must not change on failure")
				assert.Equal(t, EUR, m.GetCurrency())
			})
		}
	})
}

// TestGobRoundTrip is the reason the binary pair exists. Money and Decimal
// have only unexported fields, so gob could not encode them at all — it
// consults GobEncoder and BinaryMarshaler, and never falls back to
// MarshalText. A Currency did encode, but as the integer behind the constant,
// which stops meaning the same thing the moment the list is renumbered.
func TestGobRoundTrip(t *testing.T) {
	type portfolio struct {
		Name     string
		Balance  Money
		Reported Currency
		Weights  map[Currency]Money
	}

	original := portfolio{
		Name:     "main",
		Balance:  MustMoneyFromString("1234.56", EUR),
		Reported: JPY,
		Weights: map[Currency]Money{
			USD: MustMoneyFromString("0.75", USD),
			BHD: MustMoneyFromString("-0.001", BHD),
		},
	}

	var buf bytes.Buffer
	require.NoError(t, gob.NewEncoder(&buf).Encode(original))

	var decoded portfolio
	require.NoError(t, gob.NewDecoder(&buf).Decode(&decoded))

	assert.Equal(t, original.Name, decoded.Name)
	assert.Equal(t, original.Reported, decoded.Reported)
	assert.True(t, decoded.Balance.Equal(original.Balance))
	assert.Equal(t, EUR, decoded.Balance.GetCurrency())

	require.Len(t, decoded.Weights, 2)

	for currency, want := range original.Weights {
		got, ok := decoded.Weights[currency]
		require.True(t, ok, currency)
		assert.True(t, got.Equal(want), "%v: %v became %v", currency, want, got)
		assert.Equal(t, currency, got.GetCurrency())
	}
}

// TestCurrencyBinaryIsNotTheEnumInteger guards the decision that makes the
// format survive a renumbering: the bytes must carry the ISO code, so a build
// that assigns USD a different number still reads it as USD.
func TestCurrencyBinaryIsNotTheEnumInteger(t *testing.T) {
	encoded, err := USD.MarshalBinary()
	require.NoError(t, err)

	assert.Equal(t, "USD", string(encoded[1:]))
	assert.NotContains(t, encoded[1:], byte(USD), "the enum's number must not appear")
}
