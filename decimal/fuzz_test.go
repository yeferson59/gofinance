package decimal

import (
	"bytes"
	json "encoding/json/v2"
	"math/big"
	"strings"
	"testing"
)

// The fuzz targets below all follow the same shape: drive an entry point with
// arbitrary input and assert a property that must hold whatever comes back.
// None of them assert a specific value, so they stay useful as the engine
// changes.
//
// math/big is the oracle where one is needed. It is arbitrary-precision and
// part of the standard library, so a disagreement points at this package.

// bigFromDecimal converts a Decimal into an exact big.Rat, the common ground
// for comparing against math/big.
func bigFromDecimal(t *testing.T, d Decimal) *big.Rat {
	t.Helper()

	rat, ok := new(big.Rat).SetString(d.String())
	if !ok {
		t.Fatalf("could not read %q as a rational", d.String())
	}

	return rat
}

// keptPrecision returns how many digits a Decimal actually carries after the
// point, read from its printed form. Decimal exposes no precision accessor, and
// the tolerance checks below need to know how much the engine kept.
func keptPrecision(d Decimal) int {
	printed := d.String()

	point := strings.IndexByte(printed, '.')
	if point < 0 {
		return 0
	}

	return len(printed) - point - 1
}

// FuzzNewFromString checks the parser never panics and that whatever it
// accepts round-trips through String back to the same value.
func FuzzNewFromString(f *testing.F) {
	seeds := []string{
		"0", "-0", "1", "-1", "0.1", "-0.1", "3.14159",
		"1234567890123456789", "-1234567890123456789",
		"0.0000000000000000001", "1e10", "1E-5", "+5", ".5", "5.",
		"", " ", "abc", "1.2.3", "--1", "1-", "0x10", "NaN", "Inf",
		"99999999999999999999999999999999999999999",
		"0.00000000000000000000000000000000000001",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		parsed, err := NewFromString(input)
		if err != nil {
			// Rejecting an input is always allowed; it must simply not panic.
			return
		}

		// Anything accepted must print back to something that parses again to
		// the same value.
		printed := parsed.String()

		reparsed, err := NewFromString(printed)
		if err != nil {
			t.Fatalf("NewFromString(%q) printed %q, which no longer parses: %v", input, printed, err)
		}

		if !reparsed.Equal(parsed) {
			t.Fatalf("NewFromString(%q) = %v, printed %q, reparsed to %v", input, parsed, printed, reparsed)
		}
	})
}

// FuzzStringIsExact checks String prints the exact value rather than an
// approximation: math/big must read the printed text back to the same
// rational.
func FuzzStringIsExact(f *testing.F) {
	f.Add(int64(0), uint8(0))
	f.Add(int64(1), uint8(0))
	f.Add(int64(-1), uint8(5))
	f.Add(int64(1234567890), uint8(9))
	f.Add(int64(-999999999999999999), uint8(19))

	f.Fuzz(func(t *testing.T, coefficient int64, precision uint8) {
		value, err := NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		printed := value.String()

		// The printed form must be a plain decimal literal, not scientific
		// notation, so downstream parsers and SQL drivers can read it.
		if strings.ContainsAny(printed, "eE") {
			t.Fatalf("String() produced scientific notation: %q", printed)
		}

		// coefficient / 10^precision, computed exactly.
		expected := new(big.Rat).SetFrac(
			big.NewInt(coefficient),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(precision)), nil),
		)

		if got := bigFromDecimal(t, value); got.Cmp(expected) != 0 {
			t.Fatalf("NewFromInt64(%d, %d).String() = %q, want %s",
				coefficient, precision, printed, expected.FloatString(int(precision)))
		}
	})
}

// FuzzAddAgainstBig cross-checks addition with math/big. Whenever TryAdd
// succeeds its result must be exact, since decimal addition needs no rounding
// while the result fits.
func FuzzAddAgainstBig(f *testing.F) {
	f.Add(int64(1), uint8(0), int64(2), uint8(0))
	f.Add(int64(-5), uint8(2), int64(5), uint8(2))
	f.Add(int64(1), uint8(19), int64(1), uint8(0))

	f.Fuzz(func(t *testing.T, leftCoef int64, leftPrec uint8, rightCoef int64, rightPrec uint8) {
		left, err := NewFromInt64(leftCoef, leftPrec)
		if err != nil {
			return
		}

		right, err := NewFromInt64(rightCoef, rightPrec)
		if err != nil {
			return
		}

		sum, err := left.TryAdd(right)
		if err != nil {
			// Overflow is a legitimate outcome; it must be reported, not
			// panicked, which the absence of a panic here already proves.
			return
		}

		expected := new(big.Rat).Add(bigFromDecimal(t, left), bigFromDecimal(t, right))

		if got := bigFromDecimal(t, sum); got.Cmp(expected) != 0 {
			t.Fatalf("%v + %v = %v, want %s", left, right, sum, expected.FloatString(40))
		}
	})
}

// FuzzMulAgainstBig cross-checks multiplication with math/big. A product can
// need more digits than the engine keeps, so the check allows rounding but
// bounds it: the result must be within one unit of the last place it kept.
func FuzzMulAgainstBig(f *testing.F) {
	f.Add(int64(3), uint8(0), int64(4), uint8(0))
	f.Add(int64(-25), uint8(2), int64(4), uint8(1))
	f.Add(int64(1), uint8(10), int64(1), uint8(10))

	f.Fuzz(func(t *testing.T, leftCoef int64, leftPrec uint8, rightCoef int64, rightPrec uint8) {
		left, err := NewFromInt64(leftCoef, leftPrec)
		if err != nil {
			return
		}

		right, err := NewFromInt64(rightCoef, rightPrec)
		if err != nil {
			return
		}

		product, err := left.TryMul(right)
		if err != nil {
			return
		}

		expected := new(big.Rat).Mul(bigFromDecimal(t, left), bigFromDecimal(t, right))

		// Tolerance: one unit in the last place the result actually carries.
		got := bigFromDecimal(t, product)

		difference := new(big.Rat).Sub(got, expected)
		difference.Abs(difference)

		tolerance := new(big.Rat).SetFrac(
			big.NewInt(1),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(keptPrecision(product))), nil),
		)

		if difference.Cmp(tolerance) > 0 {
			t.Fatalf("%v × %v = %v, want %s (off by %s)",
				left, right, product, expected.FloatString(40), difference.FloatString(40))
		}
	})
}

// FuzzDivMulRoundTrip checks that dividing and multiplying back stays close to
// the original: (a/b)×b ≈ a, within the precision the quotient kept.
func FuzzDivMulRoundTrip(f *testing.F) {
	f.Add(int64(1), int64(3))
	f.Add(int64(10), int64(4))
	f.Add(int64(-7), int64(11))
	f.Add(int64(1), int64(1))

	f.Fuzz(func(t *testing.T, numerator, denominator int64) {
		if denominator == 0 {
			return
		}

		left, err := NewFromInt64(numerator, 0)
		if err != nil {
			return
		}

		right, err := NewFromInt64(denominator, 0)
		if err != nil {
			return
		}

		quotient, err := left.Div(right)
		if err != nil {
			return
		}

		recovered, err := quotient.TryMul(right)
		if err != nil {
			return
		}

		// The quotient is rounded to at most 19 places, so multiplying back
		// can miss by that rounding scaled by the divisor.
		difference := new(big.Rat).Sub(bigFromDecimal(t, recovered), bigFromDecimal(t, left))
		difference.Abs(difference)

		divisor := new(big.Rat).Abs(bigFromDecimal(t, right))
		tolerance := new(big.Rat).SetFrac(big.NewInt(1), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		tolerance.Mul(tolerance, divisor)

		if difference.Cmp(tolerance) > 0 {
			t.Fatalf("(%v / %v) × %v = %v, want %v (off by %s)",
				left, right, right, recovered, left, difference.FloatString(40))
		}
	})
}

// FuzzCompareIsConsistent checks the comparison operators agree with each
// other and with subtraction's sign, on arbitrary pairs.
func FuzzCompareIsConsistent(f *testing.F) {
	f.Add(int64(1), uint8(0), int64(1), uint8(0))
	f.Add(int64(1), uint8(0), int64(10), uint8(1))
	f.Add(int64(-1), uint8(3), int64(1), uint8(3))

	f.Fuzz(func(t *testing.T, leftCoef int64, leftPrec uint8, rightCoef int64, rightPrec uint8) {
		left, err := NewFromInt64(leftCoef, leftPrec)
		if err != nil {
			return
		}

		right, err := NewFromInt64(rightCoef, rightPrec)
		if err != nil {
			return
		}

		equal := left.Equal(right)
		less := left.LessThan(right)
		greater := left.GreaterThan(right)

		// Exactly one of the three must hold.
		count := 0
		for _, held := range []bool{equal, less, greater} {
			if held {
				count++
			}
		}

		if count != 1 {
			t.Fatalf("%v vs %v: equal=%v less=%v greater=%v", left, right, equal, less, greater)
		}

		if left.LessThanOrEqual(right) != (less || equal) {
			t.Fatalf("%v vs %v: LessThanOrEqual disagrees with LessThan/Equal", left, right)
		}

		if left.GreaterThanOrEqual(right) != (greater || equal) {
			t.Fatalf("%v vs %v: GreaterThanOrEqual disagrees with GreaterThan/Equal", left, right)
		}

		// The comparison must agree with the sign of the difference.
		if difference, err := left.TrySub(right); err == nil {
			switch {
			case equal && !difference.IsZero():
				t.Fatalf("%v == %v but the difference is %v", left, right, difference)
			case less && !difference.IsNeg():
				t.Fatalf("%v < %v but the difference is %v", left, right, difference)
			case greater && !difference.IsPos():
				t.Fatalf("%v > %v but the difference is %v", left, right, difference)
			}
		}
	})
}

// FuzzRoundBankStaysWithinAUnit checks banker's rounding moves a value by less
// than one unit of the target precision, and never changes its sign.
func FuzzRoundBankStaysWithinAUnit(f *testing.F) {
	f.Add(int64(125), uint8(2), uint8(1))
	f.Add(int64(-125), uint8(2), uint8(1))
	f.Add(int64(1), uint8(19), uint8(0))

	f.Fuzz(func(t *testing.T, coefficient int64, precision, target uint8) {
		value, err := NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		rounded := value.RoundBank(target)

		difference := new(big.Rat).Sub(bigFromDecimal(t, rounded), bigFromDecimal(t, value))
		difference.Abs(difference)

		unit := new(big.Rat).SetFrac(
			big.NewInt(1),
			new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(target)), nil),
		)

		// Rounding moves by at most half a unit, so a whole unit is a safe
		// bound that also tolerates the engine's own precision cap.
		if difference.Cmp(unit) > 0 {
			t.Fatalf("RoundBank(%d) moved %v to %v, by %s", target, value, rounded, difference.FloatString(40))
		}

		if !rounded.IsZero() && rounded.Sign() != value.Sign() {
			t.Fatalf("RoundBank(%d) changed the sign of %v to %v", target, value, rounded)
		}
	})
}

// FuzzAbsNegAreInverses checks the two sign operations behave as their names
// promise on arbitrary values.
func FuzzAbsNegAreInverses(f *testing.F) {
	f.Add(int64(5), uint8(0))
	f.Add(int64(-5), uint8(3))
	f.Add(int64(0), uint8(0))

	f.Fuzz(func(t *testing.T, coefficient int64, precision uint8) {
		value, err := NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		if absolute := value.Abs(); absolute.IsNeg() {
			t.Fatalf("Abs(%v) = %v is negative", value, absolute)
		}

		if negated := value.Neg().Neg(); !negated.Equal(value) {
			t.Fatalf("Neg(Neg(%v)) = %v", value, negated)
		}

		if value.IsZero() {
			return
		}

		if value.Neg().Sign() == value.Sign() {
			t.Fatalf("Neg(%v) kept the sign", value)
		}
	})
}

// FuzzUnmarshalText checks the text decoder never panics on arbitrary bytes,
// and that whatever it accepts encodes back to text that decodes to the same
// value. Round-tripping is the contract every text-based encoder relies on.
func FuzzUnmarshalText(f *testing.F) {
	seeds := []string{
		"0", "-0", "+0", "1", "-1", "0.1", "-0.1", "3.14159",
		"1234567890123456789", "-1234567890123456789",
		"0.0000000000000000001", "1.500", "42",
		"", " ", "1 ", "1e2", "1E-5", ".5", "5.", "abc", "$1.00",
		"1,000", "1.2.3", "--1", "0x10", "NaN", "Inf",
		"99999999999999999999999999999999999999999",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		var decoded Decimal
		if err := decoded.UnmarshalText(input); err != nil {
			// Rejecting an input is always allowed; it must simply not panic.
			return
		}

		encoded, err := decoded.MarshalText()
		if err != nil {
			t.Fatalf("UnmarshalText(%q) accepted, but MarshalText then failed: %v", input, err)
		}

		// The encoder must never emit notation its own decoder rejects.
		if strings.ContainsAny(string(encoded), "eE") {
			t.Fatalf("MarshalText produced scientific notation: %q", encoded)
		}

		var again Decimal
		if err := again.UnmarshalText(encoded); err != nil {
			t.Fatalf("MarshalText produced %q, which no longer decodes: %v", encoded, err)
		}

		if !again.Equal(decoded) {
			t.Fatalf("UnmarshalText(%q) = %v, encoded %q, decoded back to %v",
				input, decoded, encoded, again)
		}
	})
}

// FuzzUnmarshalJSON checks the JSON decoder never panics on arbitrary bytes,
// that the v1 and v2 entry points agree on every input, and that whatever is
// accepted round-trips through MarshalJSON to the same value.
func FuzzUnmarshalJSON(f *testing.F) {
	seeds := []string{
		"0", "-0", "1", "-1", "1.5", "-0.001", "1234567890123456789",
		"0.0000000000000000001", `"9.99"`, `"-0.5"`, `"1.500"`, `"1.5"`,
		"1e2", `"1e2"`, `""`, `" 1"`, `"abc"`, "true", "null", "[1]", `{"a":1}`,
		"not-json", "", "0.00000000000000000001",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		var viaV1 Decimal
		errV1 := viaV1.UnmarshalJSON(input)

		var viaV2 Decimal
		errV2 := json.Unmarshal(input, &viaV2)

		// The two entry points must accept and reject the same documents;
		// MarshalerTo's contract is that they behave equivalently under
		// default options.
		if (errV1 == nil) != (errV2 == nil) {
			t.Fatalf("UnmarshalJSON(%q) = %v, but json.Unmarshal = %v", input, errV1, errV2)
		}

		if errV1 != nil {
			// Rejecting an input is always allowed; it must simply not panic.
			return
		}

		if !viaV1.Equal(viaV2) {
			t.Fatalf("UnmarshalJSON(%q) = %v, but json.Unmarshal = %v", input, viaV1, viaV2)
		}

		encoded, err := viaV1.MarshalJSON()
		if err != nil {
			t.Fatalf("UnmarshalJSON(%q) accepted, but MarshalJSON then failed: %v", input, err)
		}

		var again Decimal
		if err := again.UnmarshalJSON(encoded); err != nil {
			t.Fatalf("MarshalJSON produced %q, which no longer decodes: %v", encoded, err)
		}

		if !again.Equal(viaV1) {
			t.Fatalf("UnmarshalJSON(%q) = %v, encoded %q, decoded back to %v",
				input, viaV1, encoded, again)
		}
	})
}

// FuzzJSONMapKeyRoundTrip checks a Decimal used as a map key survives the
// round trip that MarshalJSON alone could not even start.
//
// The decoded key is compared with Equal rather than ==. A Decimal is
// comparable, so Go lets it be a map key, but == compares the stored
// representation while Equal compares the value: the encoded form is String's,
// which trims trailing zeros, so NewFromInt64(150, 2) comes back as coefficient
// 15 scale 1. Numerically the same number, a different struct.
func FuzzJSONMapKeyRoundTrip(f *testing.F) {
	f.Add(int64(0), uint8(0))
	f.Add(int64(1), uint8(0))
	f.Add(int64(-1), uint8(5))
	f.Add(int64(1234567890), uint8(9))
	f.Add(int64(-999999999999999999), uint8(19))

	f.Fuzz(func(t *testing.T, coefficient int64, precision uint8) {
		key, err := NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		encoded, err := json.Marshal(map[Decimal]int{key: 1})
		if err != nil {
			t.Fatalf("marshaling %v as a map key: %v", key, err)
		}

		var decoded map[Decimal]int
		if err := json.Unmarshal(encoded, &decoded); err != nil {
			t.Fatalf("%q no longer decodes: %v", encoded, err)
		}

		if len(decoded) != 1 {
			t.Fatalf("%v encoded to %q, which decoded to %d keys", key, encoded, len(decoded))
		}

		for got, count := range decoded {
			if !got.Equal(key) || count != 1 {
				t.Fatalf("%v encoded to %q, which decoded to %v", key, encoded, decoded)
			}
		}
	})
}

// FuzzUnmarshalBinary checks the binary decoder never panics on arbitrary
// bytes, and that whatever it accepts encodes back to exactly the same bytes.
// A binary format is persisted, so byte-for-byte stability is the property
// that matters: two writers of the same value must produce the same blob.
func FuzzUnmarshalBinary(f *testing.F) {
	for _, seed := range []string{"0", "1", "-1", "1.5", "-0.001", "1234567890123456789"} {
		encoded, err := MustFromString(seed).MarshalBinary()
		if err != nil {
			f.Fatal(err)
		}

		f.Add(encoded)
	}

	f.Add([]byte(nil))
	f.Add([]byte{binaryVersion})
	f.Add(make([]byte, binaryLen))

	f.Fuzz(func(t *testing.T, input []byte) {
		var decoded Decimal
		if err := decoded.UnmarshalBinary(input); err != nil {
			// Rejecting an input is always allowed; it must simply not panic.
			return
		}

		encoded, err := decoded.MarshalBinary()
		if err != nil {
			t.Fatalf("UnmarshalBinary(%x) accepted, but MarshalBinary then failed: %v", input, err)
		}

		if !bytes.Equal(encoded, input) {
			t.Fatalf("UnmarshalBinary(%x) re-encoded as %x", input, encoded)
		}

		var again Decimal
		if err := again.UnmarshalBinary(encoded); err != nil {
			t.Fatalf("re-encoded %x, which no longer decodes: %v", encoded, err)
		}

		if again != decoded {
			t.Fatalf("round trip of %x changed %v into %v", input, decoded, again)
		}

		// The binary and text forms must describe the same number, however
		// differently they store it.
		text, err := decoded.MarshalText()
		if err != nil {
			t.Fatalf("decoded %x but could not write it as text: %v", input, err)
		}

		var viaText Decimal
		if err := viaText.UnmarshalText(text); err != nil {
			t.Fatalf("decoded %x, whose text %q does not parse: %v", input, text, err)
		}

		if !viaText.Equal(decoded) {
			t.Fatalf("%x is %v as binary but %v as text %q", input, decoded, viaText, text)
		}
	})
}
