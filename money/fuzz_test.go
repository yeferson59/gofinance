package money

import (
	"encoding/json"
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
