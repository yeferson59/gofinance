package decimal

import (
	"encoding"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ encoding.BinaryMarshaler   = Decimal{}
	_ encoding.BinaryAppender    = Decimal{}
	_ encoding.BinaryUnmarshaler = (*Decimal)(nil)
)

// TestDecimalBinaryLayout pins the encoding byte for byte. A binary format is
// persisted, so a change here is a change to data already written: this test
// is what makes that change deliberate rather than accidental.
func TestDecimalBinaryLayout(t *testing.T) {
	cases := map[string]string{
		// version, header (sign<<7 | scale), coefficient hi, coefficient lo.
		"0":         "010000000000000000000000000000000000",
		"1":         "010000000000000000000000000000000001",
		"-1":        "018000000000000000000000000000000001",
		"1234.5678": "010400000000000000000000000000bc614e",
		"-1.50":     "018200000000000000000000000000000096",
	}

	for input, want := range cases {
		encoded, err := MustFromString(input).MarshalBinary()
		require.NoError(t, err)
		assert.Equal(t, want, hex.EncodeToString(encoded), input)
		assert.Len(t, encoded, binaryLen, input)
	}
}

// TestDecimalBinaryPreservesRepresentation checks what binary offers that text
// does not: the trailing zero survives. MarshalText writes the value, this
// writes the number as it is stored.
func TestDecimalBinaryPreservesRepresentation(t *testing.T) {
	padded := MustFromInt64(150, 2) // 1.50

	encoded, err := padded.MarshalBinary()
	require.NoError(t, err)

	var decoded Decimal
	require.NoError(t, decoded.UnmarshalBinary(encoded))

	assert.Equal(t, padded, decoded, "the same representation, not merely the same value")

	// The text form normalises instead, which is the right answer for a
	// document and the wrong one for a faithful copy.
	text, err := padded.MarshalText()
	require.NoError(t, err)

	var viaText Decimal
	require.NoError(t, viaText.UnmarshalText(text))

	assert.True(t, viaText.Equal(padded))
	assert.NotEqual(t, padded, viaText)
}

// TestDecimalBinaryRoundTrip covers the range the layout has to carry,
// including the widest coefficient and the deepest scale.
func TestDecimalBinaryRoundTrip(t *testing.T) {
	values := []Decimal{
		{},
		Zero,
		One,
		MustFromString("-1"),
		MustFromString("0.0000000000000000001"),
		MustFromString("-0.0000000000000000001"),
		MustFromString("1234567890123456789"),
		MustFromString("-99999999999999999999999999999999999999"),
		MustFromInt64(150, 2),
	}

	for _, value := range values {
		encoded, err := value.MarshalBinary()
		require.NoError(t, err, value)

		var decoded Decimal
		require.NoError(t, decoded.UnmarshalBinary(encoded), value)
		assert.Equal(t, value, decoded, "%v did not survive the round trip", value)
	}
}

// TestDecimalAppendBinary checks AppendBinary writes into a caller's buffer
// rather than allocating one per value.
func TestDecimalAppendBinary(t *testing.T) {
	buf := []byte{0xff}

	buf, err := MustFromString("1").AppendBinary(buf)
	require.NoError(t, err)

	buf, err = MustFromString("-1").AppendBinary(buf)
	require.NoError(t, err)

	require.Len(t, buf, 1+2*binaryLen)
	assert.Equal(t, byte(0xff), buf[0], "the caller's byte must survive")

	var first, second Decimal
	require.NoError(t, first.UnmarshalBinary(buf[1:1+binaryLen]))
	require.NoError(t, second.UnmarshalBinary(buf[1+binaryLen:]))

	assert.Equal(t, "1", first.String())
	assert.Equal(t, "-1", second.String())
}

// TestDecimalUnmarshalBinaryRejects checks nothing the layout does not define
// is waved through. Binary input is machine-written, so an oddity in it means
// the reader and the writer disagree, not that the input was sloppy.
func TestDecimalUnmarshalBinaryRejects(t *testing.T) {
	valid, err := MustFromString("1.5").MarshalBinary()
	require.NoError(t, err)

	truncated := valid[:len(valid)-1]

	padded := append(append([]byte(nil), valid...), 0)

	wrongVersion := append([]byte(nil), valid...)
	wrongVersion[0] = 2

	reservedBits := append([]byte(nil), valid...)
	reservedBits[1] |= 0x40

	scaleTooDeep := append([]byte(nil), valid...)
	scaleTooDeep[1] = 20

	cases := map[string][]byte{
		"empty":           {},
		"nil":             nil,
		"truncated":       truncated,
		"padded":          padded,
		"unknown version": wrongVersion,
		"reserved bits":   reservedBits,
		"scale too deep":  scaleTooDeep,
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			d := MustFromString("99")
			assert.Error(t, d.UnmarshalBinary(input))
			assert.Equal(t, "99", d.String(), "the receiver must not change on failure")
		})
	}
}

// TestDecimalUnmarshalBinaryDoesNotRetainInput checks the decoder copies what
// it needs: callers are free to reuse the buffer they hand over.
func TestDecimalUnmarshalBinaryDoesNotRetainInput(t *testing.T) {
	encoded, err := MustFromString("1.25").MarshalBinary()
	require.NoError(t, err)

	var d Decimal
	require.NoError(t, d.UnmarshalBinary(encoded))

	for i := range encoded {
		encoded[i] = 0xff
	}

	assert.Equal(t, "1.25", d.String())
}

// TestDecimalBinaryZeroIsCanonical checks zero has exactly one encoding, so
// equal values encode to equal bytes and a blob can be compared, deduplicated
// or used as a cache key. A sign or a scale attached to a zero coefficient is
// rejected rather than quietly collapsed — found by the round-trip fuzz
// target, which caught such a blob re-encoding as different bytes.
func TestDecimalBinaryZeroIsCanonical(t *testing.T) {
	encoded, err := Zero.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, make([]byte, binaryLen-1), encoded[1:], "canonical zero is all zero bytes")

	crafted := func(header byte) []byte {
		return append([]byte{binaryVersion, header}, make([]byte, 16)...)
	}

	var d Decimal

	assert.Error(t, d.UnmarshalBinary(crafted(binaryNegBit)), "negative zero")
	assert.Error(t, d.UnmarshalBinary(crafted(5)), "zero with a scale")
	assert.Error(t, d.UnmarshalBinary(crafted(binaryNegBit|5)), "both")

	require.NoError(t, d.UnmarshalBinary(crafted(0)))
	assert.Equal(t, Zero, d)
}
