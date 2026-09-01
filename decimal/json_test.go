package decimal

import (
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	_ json.Marshaler       = Decimal{}
	_ json.MarshalerTo     = Decimal{}
	_ json.Unmarshaler     = (*Decimal)(nil)
	_ json.UnmarshalerFrom = (*Decimal)(nil)
)

// TestDecimalJSONValueIsBareNumber checks the wire format has not changed: a
// Decimal in a value position is still a bare JSON number, carrying every
// digit rather than being rounded through float64.
func TestDecimalJSONValueIsBareNumber(t *testing.T) {
	type wrapper struct {
		D Decimal `json:"d"`
	}

	cases := map[string]string{
		"0":                     `{"d":0}`,
		"1.5":                   `{"d":1.5}`,
		"-0.001":                `{"d":-0.001}`,
		"1234567890123456789":   `{"d":1234567890123456789}`,
		"0.0000000000000000001": `{"d":0.0000000000000000001}`,
	}

	for input, want := range cases {
		encoded, err := json.Marshal(wrapper{MustFromString(input)})
		require.NoError(t, err)
		assert.Equal(t, want, string(encoded), input)
	}
}

// TestDecimalMarshalJSONMatchesMarshalJSONTo checks the v1 and v2 entry points
// agree under default options, which the MarshalerTo contract requires.
func TestDecimalMarshalJSONMatchesMarshalJSONTo(t *testing.T) {
	for _, input := range []string{"0", "-42", "3.14159", "999999.999999", "-0.0000000000000000001"} {
		d := MustFromString(input)

		v1, err := d.MarshalJSON()
		require.NoError(t, err)

		v2, err := json.Marshal(d)
		require.NoError(t, err)

		assert.Equal(t, string(v1), string(v2), input)
		assert.Equal(t, input, string(v1), input)
	}
}

// TestDecimalAsJSONMapKey covers the gap that made a Decimal unusable as a map
// key: a JSON object name must be a string, and the bare number MarshalJSON
// returns is rejected outright by the encoder.
func TestDecimalAsJSONMapKey(t *testing.T) {
	encoded, err := json.Marshal(map[Decimal]int{MustFromString("1.5"): 1})
	require.NoError(t, err)
	assert.Equal(t, `{"1.5":1}`, string(encoded))

	var decoded map[Decimal]int
	require.NoError(t, json.Unmarshal([]byte(`{"2.25":3}`), &decoded))
	assert.Equal(t, map[Decimal]int{MustFromString("2.25"): 3}, decoded)
}

// TestDecimalAsNestedJSONMapKey checks the key detection reads the encoder's
// current position rather than assuming a top-level object.
func TestDecimalAsNestedJSONMapKey(t *testing.T) {
	value := map[string]map[Decimal]Decimal{
		"rates": {MustFromString("1.5"): MustFromString("-0.25")},
	}

	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	assert.Equal(t, `{"rates":{"1.5":-0.25}}`, string(encoded))

	var decoded map[string]map[Decimal]Decimal
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, value, decoded)
}

// TestDecimalJSONStringTag checks the `,string` tag, which producers use so a
// decimal survives readers that parse every JSON number as a float64.
func TestDecimalJSONStringTag(t *testing.T) {
	type wrapper struct {
		D Decimal `json:"d,string"`
	}

	encoded, err := json.Marshal(wrapper{MustFromString("1.5")})
	require.NoError(t, err)
	assert.Equal(t, `{"d":"1.5"}`, string(encoded))

	var decoded wrapper
	require.NoError(t, json.Unmarshal(encoded, &decoded))
	assert.Equal(t, "1.5", decoded.D.String())
}

// TestDecimalJSONInArrayAndTopLevel checks the positions that must stay bare
// numbers are not caught by the object-name detection.
func TestDecimalJSONInArrayAndTopLevel(t *testing.T) {
	encoded, err := json.Marshal([]Decimal{MustFromString("1.5"), MustFromString("-0.25")})
	require.NoError(t, err)
	assert.Equal(t, `[1.5,-0.25]`, string(encoded))

	encoded, err = json.Marshal(MustFromString("1.5"))
	require.NoError(t, err)
	assert.Equal(t, `1.5`, string(encoded))

	encoded, err = json.Marshal(map[string]Decimal{"a": MustFromString("1.5")})
	require.NoError(t, err)
	assert.Equal(t, `{"a":1.5}`, string(encoded))
}

// TestDecimalUnmarshalJSONAcceptsQuotedNumber checks the string form decodes
// even without the `,string` tag: it is what map keys produce, and what a
// producer sends to protect precision.
func TestDecimalUnmarshalJSONAcceptsQuotedNumber(t *testing.T) {
	cases := map[string]string{
		`"9.99"`:     "9.99",
		`"-0.5"`:     "-0.5",
		`"1.500"`:    "1.5",
		`"\u0031.5"`: "1.5", // escapes are legal JSON, if pointless in a number
		"9.99":       "9.99",
		"-0.5":       "-0.5",
		"1234567890": "1234567890",
	}

	for input, want := range cases {
		var d Decimal
		require.NoError(t, d.UnmarshalJSON([]byte(input)), input)
		assert.Equal(t, want, d.String(), input)

		var viaV2 Decimal
		require.NoError(t, json.Unmarshal([]byte(input), &viaV2), input)
		assert.Equal(t, want, viaV2.String(), input)
	}
}

// TestDecimalUnmarshalJSONRejects checks a quoted value is held to the same
// standard as an unquoted one, and that composites are refused rather than
// partially consumed.
func TestDecimalUnmarshalJSONRejects(t *testing.T) {
	cases := map[string]string{
		"exponent":        "1e2",
		"quoted exponent": `"1e2"`,
		"quoted letters":  `"abc"`,
		"quoted empty":    `""`,
		"quoted space":    `" 1"`,
		"quoted currency": `"$1.00"`,
		"boolean":         "true",
		"null":            "null",
		"object":          `{"a":1}`,
		"array":           "[1]",
		"malformed":       "not-json",
		"too precise":     "0.00000000000000000001",
	}

	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			d := MustFromString("99")
			assert.Error(t, d.UnmarshalJSON([]byte(input)))
			assert.Equal(t, "99", d.String(), "the receiver must not change on failure")

			viaV2 := MustFromString("99")
			assert.Error(t, json.Unmarshal([]byte(input), &viaV2))
			assert.Equal(t, "99", viaV2.String(), "the receiver must not change on failure")
		})
	}
}

// TestDecimalUnmarshalJSONFromReadsOneValue checks the decoder is left on the
// next value after a rejected composite, which the UnmarshalerFrom contract
// requires: reading a single token would strand it inside the object.
func TestDecimalUnmarshalJSONFromReadsOneValue(t *testing.T) {
	dec := jsontext.NewDecoder(bytesReader(`[{"a":1},2.5]`))

	_, err := dec.ReadToken() // '['
	require.NoError(t, err)

	var d Decimal
	assert.Error(t, d.UnmarshalJSONFrom(dec), "an object is not a decimal")

	// The rejected object must have been consumed whole, leaving the next
	// element readable.
	require.NoError(t, d.UnmarshalJSONFrom(dec))
	assert.Equal(t, "2.5", d.String())
}

// TestDecimalMarshalJSONToWritesIntoCallerBuffer checks the streaming encoder
// composes with a surrounding document rather than replacing it.
func TestDecimalMarshalJSONToWritesIntoCallerBuffer(t *testing.T) {
	var out []byte

	enc := jsontext.NewEncoder(sliceWriter{&out})

	require.NoError(t, enc.WriteToken(jsontext.BeginArray))
	require.NoError(t, MustFromString("1.5").MarshalJSONTo(enc))
	require.NoError(t, MustFromString("-0.25").MarshalJSONTo(enc))
	require.NoError(t, enc.WriteToken(jsontext.EndArray))

	assert.JSONEq(t, `[1.5,-0.25]`, string(out))
}

// TestDecimalJSONMapKeyNormalisesRepresentation records a sharp edge that
// predates the map-key support but only becomes visible through it: Decimal is
// comparable, so Go accepts it as a map key, but == compares the stored
// representation while Equal compares the value. The encoded key is String's
// form, which trims trailing zeros, so a round trip can return a key that is
// numerically equal to the original and yet not the same map key.
func TestDecimalJSONMapKeyNormalisesRepresentation(t *testing.T) {
	padded := MustFromInt64(150, 2) // 1.50
	trimmed := MustFromString("1.5")

	assert.True(t, padded.Equal(trimmed), "the same number")
	assert.NotEqual(t, padded, trimmed, "a different representation")

	encoded, err := json.Marshal(map[Decimal]int{padded: 1})
	require.NoError(t, err)
	assert.Equal(t, `{"1.5":1}`, string(encoded))

	var decoded map[Decimal]int
	require.NoError(t, json.Unmarshal(encoded, &decoded))

	_, found := decoded[padded]
	assert.False(t, found, "== does not find it: the trailing zero is gone")

	assert.Equal(t, 1, decoded[trimmed], "the trimmed key is what came back")
}

// TestDecimalUnmarshalJSONRejectsTrailingContent covers a case json.Unmarshal
// never produces but a direct caller can: "00" is two JSON numbers, and
// reading only the first one silently decoded it as zero. The fuzz corpus
// entry that found it is kept alongside this test.
func TestDecimalUnmarshalJSONRejectsTrailingContent(t *testing.T) {
	for _, input := range []string{"00", "1 2", "1.5 2", "0.5false", `"1" "2"`, "1,2"} {
		d := MustFromString("99")
		assert.Error(t, d.UnmarshalJSON([]byte(input)), input)
		assert.Equal(t, "99", d.String(), input)
	}
}

// TestDecimalUnmarshalJSONTrailingZeros guards the reordering that check
// forced: jsontext hands back a value aliasing the decoder's buffer, which the
// next read overwrites, so the value must be parsed before anything else is
// read. Reading first turned "5.00" into "#.00" and rejected it.
func TestDecimalUnmarshalJSONTrailingZeros(t *testing.T) {
	cases := map[string]string{
		"5.00":     "5",
		"9.99":     "9.99",
		"0.10":     "0.1",
		"-1.500":   "-1.5",
		`"5.00"`:   "5",
		"100.0000": "100",
	}

	for input, want := range cases {
		var d Decimal
		require.NoError(t, d.UnmarshalJSON([]byte(input)), input)
		assert.Equal(t, want, d.String(), input)
	}
}

// TestDecimalJSONWhitespaceIsJSONsOwn checks the trailing-content check trims
// what JSON calls whitespace and nothing else: a vertical tab or a form feed
// is content, not spacing, however much bytes.TrimSpace disagrees.
func TestDecimalJSONWhitespaceIsJSONsOwn(t *testing.T) {
	for _, input := range []string{"1.5 ", "1.5\t", "1.5\n", "1.5\r\n", " 1.5 "} {
		var d Decimal
		assert.NoError(t, d.UnmarshalJSON([]byte(input)), "%q", input)
		assert.Equal(t, "1.5", d.String(), "%q", input)
	}

	for _, input := range []string{"1.5\v", "1.5\f", "1.5\x00", "1.5\u00a0"} {
		var d Decimal
		assert.Error(t, d.UnmarshalJSON([]byte(input)), "%q", input)
	}
}
