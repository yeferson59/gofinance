package decimal

import (
	"encoding"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.TextMarshaler   = Decimal{}
	_ encoding.TextAppender    = Decimal{}
	_ encoding.TextUnmarshaler = (*Decimal)(nil)
)

// TestDecimalMarshalText checks the text form is the same plain notation
// String produces: no quotes, no exponent, no thousands separators.
func TestDecimalMarshalText(t *testing.T) {
	cases := []string{"0", "1", "-1", "1.5", "-0.001", "123.45", "0.0000000000000000001"}

	for _, want := range cases {
		text, err := MustFromString(want).MarshalText()
		require.NoError(t, err)
		assert.Equal(t, want, string(text))
	}
}

// TestDecimalMarshalTextZeroValue checks the zero Decimal encodes to something
// readable rather than empty text, which UnmarshalText would reject.
func TestDecimalMarshalTextZeroValue(t *testing.T) {
	var d Decimal

	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "0", string(text))

	var decoded Decimal
	require.NoError(t, decoded.UnmarshalText(text))
	assert.True(t, decoded.IsZero())
}

// TestDecimalMarshalTextAgreesWithString checks the two never drift apart, so
// a value logged with %v and one encoded to YAML read identically.
func TestDecimalMarshalTextAgreesWithString(t *testing.T) {
	for _, input := range []string{"0", "-42", "3.14159", "999999.999999"} {
		d := MustFromString(input)

		text, err := d.MarshalText()
		require.NoError(t, err)
		assert.Equal(t, d.String(), string(text))
	}
}

// TestDecimalAppendText checks AppendText writes into a caller's buffer rather
// than allocating a fresh one per value.
func TestDecimalAppendText(t *testing.T) {
	buf := []byte("total=")

	buf, err := MustFromString("12.50").AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "total=12.5", string(buf))

	buf = append(buf, ' ')

	buf, err = MustFromString("-0.25").AppendText(buf)
	require.NoError(t, err)
	assert.Equal(t, "total=12.5 -0.25", string(buf))
}

// TestDecimalUnmarshalText checks the accepted forms, including the ones only
// a hand-written document produces: a leading plus, and a value written with
// trailing zeros that normalise away.
func TestDecimalUnmarshalText(t *testing.T) {
	cases := map[string]string{
		"0":      "0",
		"-0":     "0",
		"+7":     "7",
		"1.500":  "1.5",
		"-1.500": "-1.5",
		"0.10":   "0.1",
		"42":     "42",
	}

	for input, want := range cases {
		var d Decimal
		require.NoError(t, d.UnmarshalText([]byte(input)), input)
		assert.Equal(t, want, d.String(), input)
	}
}

// TestDecimalUnmarshalTextRoundTrip checks a value survives the encode/decode
// cycle exactly, which is the contract the whole text pair rests on.
func TestDecimalUnmarshalTextRoundTrip(t *testing.T) {
	for _, input := range []string{"0", "1", "-1", "0.5", "-1234.5678", "1.0000000000000000001"} {
		original := MustFromString(input)

		text, err := original.MarshalText()
		require.NoError(t, err)

		var decoded Decimal
		require.NoError(t, decoded.UnmarshalText(text))
		assert.True(t, original.Equal(decoded), "%s decoded to %s", original, decoded)
	}
}

// TestDecimalUnmarshalTextInvalid checks the decoder holds the same line the
// rest of the package holds. Exponent notation matters most: it is valid JSON
// and valid YAML, so rejecting it here is what stops "1e2" from arriving as a
// number this package cannot represent exactly.
func TestDecimalUnmarshalTextInvalid(t *testing.T) {
	cases := map[string][]byte{
		"empty":            []byte(""),
		"nil":              nil,
		"space":            []byte(" 1"),
		"trailing space":   []byte("1 "),
		"exponent":         []byte("1e2"),
		"exponent capital": []byte("1E2"),
		"letters":          []byte("abc"),
		"currency":         []byte("$1.00"),
		"separators":       []byte("1,000"),
		"two dots":         []byte("1.2.3"),
		"bare dot":         []byte("."),
		"sign only":        []byte("-"),
		"trailing dot":     []byte("1."),
		"too precise":      []byte("0.00000000000000000001"),
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			d := MustFromString("99")
			assert.Error(t, d.UnmarshalText(input))
			assert.Equal(t, "99", d.String(), "the receiver must not change on failure")
		})
	}
}

// TestDecimalUnmarshalTextDoesNotRetainInput checks the decoder copies what it
// needs: encoders are free to reuse the slice they hand over.
func TestDecimalUnmarshalTextDoesNotRetainInput(t *testing.T) {
	text := []byte("1.25")

	var d Decimal
	require.NoError(t, d.UnmarshalText(text))

	copy(text, "9.99")
	assert.Equal(t, "1.25", d.String())
}
