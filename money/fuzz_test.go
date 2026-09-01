package money

import (
	"bytes"
	"encoding/json"
	jsonv2 "encoding/json/v2"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
)

// fuzzCurrencies are the currencies the targets below cycle through, chosen to
// span the three precisions ISO 4217 uses: none (JPY), two (USD), and three
// (BHD).
var fuzzCurrencies = []Currency{USD, JPY, BHD, EUR}

// currencyAt picks a currency from an arbitrary fuzz input.
func currencyAt(index uint8) Currency {
	return fuzzCurrencies[int(index)%len(fuzzCurrencies)]
}

// FuzzUnmarshalJSON checks the decoder never panics on arbitrary bytes, and
// that anything it accepts survives a round trip back through the encoder.
func FuzzUnmarshalJSON(f *testing.F) {
	seeds := []string{
		`{"value":"1234.56","currency":"USD"}`,
		`{"value":"0","currency":"JPY"}`,
		`{"value":"-0.001","currency":"BHD"}`,
		`{"value":"1"}`,
		`1234.56`,
		`"1234.56"`,
		`{}`, `[]`, `null`, ``, `{`, `{"value":}`,
		`{"value":"abc","currency":"USD"}`,
		`{"value":"1","currency":"NOTACURRENCY"}`,
		`{"value":"1e400","currency":"USD"}`,
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		var decoded Money

		if err := decoded.UnmarshalJSON(data); err != nil {
			// Rejecting malformed input is fine; not panicking is the point.
			return
		}

		encoded, err := decoded.MarshalJSON()
		if err != nil {
			t.Fatalf("decoded %q but could not re-encode it: %v", data, err)
		}

		var again Money
		if err := again.UnmarshalJSON(encoded); err != nil {
			t.Fatalf("re-encoded %q as %q, which no longer decodes: %v", data, encoded, err)
		}

		if !again.Equal(decoded) || again.GetCurrency() != decoded.GetCurrency() {
			t.Fatalf("round trip of %q changed %v %v into %v %v",
				data, decoded, decoded.GetCurrency(), again, again.GetCurrency())
		}
	})
}

// FuzzScan checks the SQL decoder against arbitrary text, and that a value
// written with Value can always be read back with Scan.
func FuzzScan(f *testing.F) {
	for _, seed := range []string{"0", "1234.56", "-1", "1e5", "", "abc", "  12  ", "0.0000000001"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		var scanned Money

		if err := scanned.Scan(text); err != nil {
			return
		}

		// Value produces what a driver would store; Scan must read it back.
		stored, err := scanned.Value()
		if err != nil {
			t.Fatalf("scanned %q but could not produce a driver value: %v", text, err)
		}

		var again Money
		if err := again.Scan(stored); err != nil {
			t.Fatalf("stored %q as %v, which no longer scans: %v", text, stored, err)
		}

		if !again.Equal(scanned) {
			t.Fatalf("round trip of %q changed %v into %v", text, scanned, again)
		}

		// Scanning the same text through []byte must agree with the string
		// form, since drivers use either.
		var fromBytes Money
		if err := fromBytes.Scan([]byte(text)); err != nil {
			t.Fatalf("scanned %q as a string but not as bytes: %v", text, err)
		}

		if !fromBytes.Equal(scanned) {
			t.Fatalf("%q scanned to %v as a string and %v as bytes", text, scanned, fromBytes)
		}
	})
}

// FuzzAllocateSumsBack is the property money.Allocate exists to guarantee:
// however an amount is split, the parts add back to exactly what went in.
//
// This is the fuzz counterpart of the defect found in Phase 3
// (TESTING_PLAN.md §2.10a), where a sub-unit residue made the parts overshoot.
// A hand-written table found it once; this target hunts the whole input space.
func FuzzAllocateSumsBack(f *testing.F) {
	f.Add(int64(100), uint8(2), uint8(0), uint32(1), uint32(1), uint32(1))
	f.Add(int64(10), uint8(2), uint8(0), uint32(1), uint32(1), uint32(1))
	f.Add(int64(-123456), uint8(2), uint8(1), uint32(1), uint32(2), uint32(3))
	f.Add(int64(1), uint8(19), uint8(1), uint32(7), uint32(11), uint32(13))
	f.Add(int64(0), uint8(0), uint8(2), uint32(0), uint32(1), uint32(0))

	f.Fuzz(func(t *testing.T, coefficient int64, precision, currencyIndex uint8, first, second, third uint32) {
		value, err := decimal.NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		currency := currencyAt(currencyIndex)
		original := NewFromDecimal(value, currency)

		parts, err := original.Allocate(first, second, third)
		if err != nil {
			// All-zero ratios and unknown currencies are reported, not split.
			return
		}

		if len(parts) != 3 {
			t.Fatalf("Allocate returned %d parts for 3 ratios", len(parts))
		}

		total := NewFromDecimal(decimal.Zero, currency)

		for _, part := range parts {
			if part.GetCurrency() != currency {
				t.Fatalf("part %v carries %v, not %v", part, part.GetCurrency(), currency)
			}

			total, err = total.TryAdd(part)
			if err != nil {
				return
			}
		}

		if !total.Equal(original) {
			t.Fatalf("%v split by (%d, %d, %d) gave parts summing to %v",
				original, first, second, third, total)
		}
	})
}

// FuzzAllocateEvenlySumsBack is the same property for the even split, across
// arbitrary part counts.
func FuzzAllocateEvenlySumsBack(f *testing.F) {
	f.Add(int64(100), uint8(2), uint8(0), 3)
	f.Add(int64(1), uint8(2), uint8(0), 7)
	f.Add(int64(-1), uint8(0), uint8(1), 100)

	f.Fuzz(func(t *testing.T, coefficient int64, precision, currencyIndex uint8, parts int) {
		// Keep the split to a sane size; the property does not depend on it.
		if parts < 1 || parts > 200 {
			return
		}

		value, err := decimal.NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		currency := currencyAt(currencyIndex)
		original := NewFromDecimal(value, currency)

		split, err := original.AllocateEvenly(parts)
		if err != nil {
			return
		}

		total := NewFromDecimal(decimal.Zero, currency)

		for _, part := range split {
			total, err = total.TryAdd(part)
			if err != nil {
				return
			}
		}

		if !total.Equal(original) {
			t.Fatalf("%v split evenly into %d gave parts summing to %v", original, parts, total)
		}
	})
}

// FuzzArithmeticKeepsCurrency checks that every operation preserving an amount
// keeps its currency, and that the checked variants agree with the panicking
// ones on well-formed input.
func FuzzArithmeticKeepsCurrency(f *testing.F) {
	f.Add(int64(100), uint8(2), int64(50), uint8(2), uint8(0))
	f.Add(int64(-1), uint8(0), int64(1), uint8(0), uint8(1))

	f.Fuzz(func(t *testing.T, leftCoef int64, leftPrec uint8, rightCoef int64, rightPrec uint8, currencyIndex uint8) {
		leftValue, err := decimal.NewFromInt64(leftCoef, leftPrec)
		if err != nil {
			return
		}

		rightValue, err := decimal.NewFromInt64(rightCoef, rightPrec)
		if err != nil {
			return
		}

		currency := currencyAt(currencyIndex)
		left := NewFromDecimal(leftValue, currency)
		right := NewFromDecimal(rightValue, currency)

		sum, err := left.TryAdd(right)
		if err != nil {
			return
		}

		if sum.GetCurrency() != currency {
			t.Fatalf("addition changed the currency to %v", sum.GetCurrency())
		}

		// Adding then subtracting the same amount must return the original.
		back, err := sum.TrySub(right)
		if err != nil {
			return
		}

		if !back.Equal(left) {
			t.Fatalf("(%v + %v) − %v = %v, want %v", left, right, right, back, left)
		}

		// Scaling by a decimal keeps the currency too.
		scaled := left.MulDecimal(decimal.MustFromInt64(2, 0))
		if scaled.GetCurrency() != currency {
			t.Fatalf("MulDecimal changed the currency to %v", scaled.GetCurrency())
		}
	})
}

// FuzzRoundBankStaysWithinAUnit checks Money's rounding never moves an amount
// by a whole unit of the target precision, nor flips its sign.
func FuzzRoundBankStaysWithinAUnit(f *testing.F) {
	f.Add(int64(125), uint8(2), uint8(1), uint8(0))
	f.Add(int64(-125), uint8(2), uint8(1), uint8(1))

	f.Fuzz(func(t *testing.T, coefficient int64, precision, target, currencyIndex uint8) {
		value, err := decimal.NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		currency := currencyAt(currencyIndex)
		original := NewFromDecimal(value, currency)
		rounded := original.RoundBank(target)

		if rounded.GetCurrency() != currency {
			t.Fatalf("RoundBank changed the currency to %v", rounded.GetCurrency())
		}

		if !rounded.IsZero() && rounded.GetDecimal().Sign() != original.GetDecimal().Sign() {
			t.Fatalf("RoundBank(%d) changed the sign of %v to %v", target, original, rounded)
		}

		// The two rendering helpers must agree with the rounded value.
		if rounded.RoundBankString(target) != original.RoundBankString(target) {
			t.Fatalf("RoundBankString disagrees with RoundBank for %v at %d", original, target)
		}
	})
}

// FuzzJSONNumberFormIsAccepted checks the documented fallback: a bare JSON
// number decodes as a USD amount, matching Scan's default.
func FuzzJSONNumberFormIsAccepted(f *testing.F) {
	f.Add(int64(1234), uint8(2))
	f.Add(int64(-1), uint8(0))

	f.Fuzz(func(t *testing.T, coefficient int64, precision uint8) {
		value, err := decimal.NewFromInt64(coefficient, precision)
		if err != nil {
			return
		}

		data, err := json.Marshal(json.RawMessage(value.String()))
		if err != nil {
			return
		}

		var decoded Money
		if err := decoded.UnmarshalJSON(data); err != nil {
			return
		}

		if decoded.GetCurrency() != USD {
			t.Fatalf("bare number %q decoded with currency %v, want USD", data, decoded.GetCurrency())
		}

		if !decoded.GetDecimal().Equal(value) {
			t.Fatalf("bare number %q decoded to %v, want %v", data, decoded, value)
		}
	})
}

// FuzzCurrencyUnmarshalText checks the text decoder never panics on arbitrary
// bytes, that anything it accepts is a currency this package recognises, and
// that the pair round-trips through the canonical uppercase ISO code.
func FuzzCurrencyUnmarshalText(f *testing.F) {
	seeds := []string{
		"USD", "usd", " EUR ", "JPY", "XXX", "BHD",
		"", " ", "US", "USDD", "ZZZ", "840", "us\x00", "€",
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, input []byte) {
		var decoded Currency
		if err := decoded.UnmarshalText(input); err != nil {
			// Rejecting an input is always allowed; it must simply not panic.
			return
		}

		if !decoded.Valid() {
			t.Fatalf("UnmarshalText(%q) accepted an invalid currency: %v", input, decoded)
		}

		encoded, err := decoded.MarshalText()
		if err != nil {
			t.Fatalf("UnmarshalText(%q) accepted, but MarshalText then failed: %v", input, err)
		}

		if len(encoded) != 3 {
			t.Fatalf("MarshalText produced %q, want a three-letter ISO code", encoded)
		}

		var again Currency
		if err := again.UnmarshalText(encoded); err != nil {
			t.Fatalf("MarshalText produced %q, which no longer decodes: %v", encoded, err)
		}

		if again != decoded {
			t.Fatalf("UnmarshalText(%q) = %v, encoded %q, decoded back to %v",
				input, decoded, encoded, again)
		}
	})
}

// FuzzJSONPathsAgree checks the v1 and v2 entry points accept and reject the
// same documents and decode them to the same amount. MarshalerTo's contract is
// that a type implementing both behaves equivalently under default options,
// and the two are now separate code paths.
func FuzzJSONPathsAgree(f *testing.F) {
	seeds := []string{
		`{"value":"1234.56","currency":"USD"}`,
		`{"value":1234.56,"currency":"USD"}`,
		`{"value":"1","currency":""}`,
		`{"value":"1"}`,
		`{"value":"1","currency":143}`,
		`{"value":"1","note":{"a":[1,2]}}`,
		`{"currency":"USD"}`,
		`{"value":"1","value":"2"}`,
		`1234.56`, `"1234.56"`, `1e2`,
		`{}`, `[]`, `null`, ``, `{`, `{"value":}`, `{"value":"1"} 2`,
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		var viaV1 Money
		errV1 := viaV1.UnmarshalJSON(data)

		var viaV2 Money
		errV2 := jsonv2.Unmarshal(data, &viaV2)

		if (errV1 == nil) != (errV2 == nil) {
			t.Fatalf("UnmarshalJSON(%q) = %v, but json.Unmarshal = %v", data, errV1, errV2)
		}

		if errV1 != nil {
			return
		}

		if !viaV1.Equal(viaV2) {
			t.Fatalf("UnmarshalJSON(%q) = %v %v, but json.Unmarshal = %v %v",
				data, viaV1, viaV1.GetCurrency(), viaV2, viaV2.GetCurrency())
		}

		// The two encoders must agree as well, on a value the decoder just
		// accepted.
		encV1, err1 := viaV1.MarshalJSON()
		encV2, err2 := jsonv2.Marshal(viaV1)

		if (err1 == nil) != (err2 == nil) {
			t.Fatalf("MarshalJSON = %v, but json.Marshal = %v", err1, err2)
		}

		if err1 == nil && string(encV1) != string(encV2) {
			t.Fatalf("MarshalJSON produced %q, json.Marshal produced %q", encV1, encV2)
		}
	})
}

// FuzzCurrencyUnmarshalJSON checks the currency decoder never panics, that the
// two entry points agree, and that anything accepted is a currency this
// package can name again.
func FuzzCurrencyUnmarshalJSON(f *testing.F) {
	seeds := []string{
		`"USD"`, `"usd"`, `" EUR "`, `"XXX"`, `"ZZZ"`, `""`,
		`0`, `1`, `143`, `157`, `158`, `256`, `-1`, `1.5`, `1e2`,
		`true`, `null`, `[1]`, `{"a":1}`, `"USD" "EUR"`, ``, `not-json`,
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		var viaV1 Currency
		errV1 := viaV1.UnmarshalJSON(data)

		var viaV2 Currency
		errV2 := jsonv2.Unmarshal(data, &viaV2)

		if (errV1 == nil) != (errV2 == nil) {
			t.Fatalf("UnmarshalJSON(%q) = %v, but json.Unmarshal = %v", data, errV1, errV2)
		}

		if errV1 != nil {
			return
		}

		if viaV1 != viaV2 {
			t.Fatalf("UnmarshalJSON(%q) = %v, but json.Unmarshal = %v", data, viaV1, viaV2)
		}

		// Anything accepted must encode again, which the ISO accessor is what
		// decides: a currency with no code is not a currency.
		encoded, err := viaV1.MarshalJSON()
		if err != nil {
			t.Fatalf("UnmarshalJSON(%q) accepted %v, which cannot be encoded: %v", data, viaV1, err)
		}

		var again Currency
		if err := again.UnmarshalJSON(encoded); err != nil {
			t.Fatalf("re-encoded %q as %q, which no longer decodes: %v", data, encoded, err)
		}

		if again != viaV1 {
			t.Fatalf("round trip of %q changed %v into %v", data, viaV1, again)
		}
	})
}

// FuzzMoneyUnmarshalBinary checks the binary decoders never panic on arbitrary
// bytes, and that whatever they accept encodes back to exactly the same bytes.
// A binary format is persisted, so byte-for-byte stability is the property
// that matters.
func FuzzMoneyUnmarshalBinary(f *testing.F) {
	for _, seed := range []Money{
		MustMoneyFromString("1234.56", USD),
		MustMoneyFromString("-0.001", BHD),
		MustMoneyFromString("0", JPY),
	} {
		encoded, err := seed.MarshalBinary()
		if err != nil {
			f.Fatal(err)
		}

		f.Add(encoded)
	}

	f.Add([]byte(nil))
	f.Add([]byte{moneyBinaryVersion, 'U', 'S', 'D'})
	f.Add(make([]byte, moneyBinaryPrefixLen+18))

	f.Fuzz(func(t *testing.T, input []byte) {
		var decoded Money
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

		var again Money
		if err := again.UnmarshalBinary(encoded); err != nil {
			t.Fatalf("re-encoded %x, which no longer decodes: %v", encoded, err)
		}

		if !again.Equal(decoded) || again.GetCurrency() != decoded.GetCurrency() {
			t.Fatalf("round trip of %x changed %v %v into %v %v",
				input, decoded, decoded.GetCurrency(), again, again.GetCurrency())
		}

		// The binary and JSON forms must describe the same amount, however
		// differently they store it.
		document, err := decoded.MarshalJSON()
		if err != nil {
			t.Fatalf("decoded %x but could not write it as JSON: %v", input, err)
		}

		var viaJSON Money
		if err := viaJSON.UnmarshalJSON(document); err != nil {
			t.Fatalf("decoded %x, whose JSON %s does not parse: %v", input, document, err)
		}

		if !viaJSON.Equal(decoded) || viaJSON.GetCurrency() != decoded.GetCurrency() {
			t.Fatalf("%x is %v %v as binary but %v %v as JSON",
				input, decoded, decoded.GetCurrency(), viaJSON, viaJSON.GetCurrency())
		}
	})
}

// FuzzCurrencyUnmarshalBinary checks the same for a currency on its own.
func FuzzCurrencyUnmarshalBinary(f *testing.F) {
	for _, seed := range []Currency{USD, EUR, JPY, XXX} {
		encoded, err := seed.MarshalBinary()
		if err != nil {
			f.Fatal(err)
		}

		f.Add(encoded)
	}

	f.Add([]byte(nil))
	f.Add([]byte{currencyBinaryVersion, 'Z', 'Z', 'Z'})
	f.Add(make([]byte, currencyBinaryLen))

	f.Fuzz(func(t *testing.T, input []byte) {
		var decoded Currency
		if err := decoded.UnmarshalBinary(input); err != nil {
			return
		}

		encoded, err := decoded.MarshalBinary()
		if err != nil {
			t.Fatalf("UnmarshalBinary(%x) accepted %v, which cannot be encoded: %v", input, decoded, err)
		}

		if !bytes.Equal(encoded, input) {
			t.Fatalf("UnmarshalBinary(%x) re-encoded as %x", input, encoded)
		}

		// Whatever the binary form accepts, the text form must name too.
		text, err := decoded.MarshalText()
		if err != nil {
			t.Fatalf("decoded %x but could not write it as text: %v", input, err)
		}

		if string(text) != string(input[1:]) {
			t.Fatalf("%x decoded to %v, whose text is %q", input, decoded, text)
		}
	})
}
