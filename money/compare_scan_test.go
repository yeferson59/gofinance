package money

import (
	"errors"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
)

// Ordering two amounts in different currencies has no answer without an
// exchange rate. TryCmp reports the mismatch instead of comparing the bare
// numbers.
func TestTryCmpRejectsCurrencyMismatch(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)
	cop := MustMoneyFromFloat64(5, COP)

	if _, err := usd.TryCmp(cop); !errors.Is(err, ErrCurrencyMismatch) {
		t.Fatalf("TryCmp error = %v, want ErrCurrencyMismatch", err)
	}

	cmp, err := usd.TryCmp(MustMoneyFromFloat64(20, USD))
	if err != nil {
		t.Fatalf("TryCmp same currency error = %v", err)
	}

	if cmp != -1 {
		t.Errorf("TryCmp(10 USD, 20 USD) = %d, want -1", cmp)
	}
}

// The panicking comparisons follow Add and Sub: a currency mismatch panics
// with ErrCurrencyMismatch rather than returning a meaningless answer.
func TestComparisonsPanicOnCurrencyMismatch(t *testing.T) {
	usd := MustMoneyFromFloat64(10, USD)
	cop := MustMoneyFromFloat64(5, COP)

	comparisons := map[string]func(){
		"Cmp":                func() { usd.Cmp(cop) },
		"GreaterThan":        func() { usd.GreaterThan(cop) },
		"GreaterThanOrEqual": func() { usd.GreaterThanOrEqual(cop) },
		"LessThan":           func() { usd.LessThan(cop) },
		"LessThanOrEqual":    func() { usd.LessThanOrEqual(cop) },
	}

	for name, compare := range comparisons {
		t.Run(name, func(t *testing.T) {
			defer func() {
				r := recover()
				err, ok := r.(error)
				if !ok || !errors.Is(err, ErrCurrencyMismatch) {
					t.Fatalf("recovered %v, want a panic with ErrCurrencyMismatch", r)
				}
			}()

			compare()
		})
	}
}

func TestComparisonsWithinOneCurrency(t *testing.T) {
	ten := MustMoneyFromFloat64(10, COP)
	five := MustMoneyFromFloat64(5, COP)

	if !ten.GreaterThan(five) || !ten.GreaterThanOrEqual(five) || ten.LessThan(five) || ten.LessThanOrEqual(five) {
		t.Error("10 COP vs 5 COP ordered wrongly")
	}

	if ten.Cmp(ten) != 0 || !ten.GreaterThanOrEqual(ten) || !ten.LessThanOrEqual(ten) {
		t.Error("10 COP vs itself ordered wrongly")
	}
}

// A database column holding only the amount carries no currency. Scan keeps
// the currency the receiver already has, so presetting it is enough to read
// a COP column back as COP.
func TestScanKeepsReceiverCurrency(t *testing.T) {
	m := NewFromDecimal(decimal.Zero, COP)

	if err := m.Scan("1500.50"); err != nil {
		t.Fatalf("Scan error = %v", err)
	}

	if m.GetCurrency() != COP {
		t.Errorf("currency = %v, want COP", m.GetCurrency())
	}

	if m.String() != "1500.5" {
		t.Errorf("value = %s, want 1500.5", m.String())
	}
}

// A receiver with no currency set still defaults to USD, as before.
func TestScanWithoutCurrencyDefaultsToUSD(t *testing.T) {
	var m Money

	if err := m.Scan(int64(42)); err != nil {
		t.Fatalf("Scan error = %v", err)
	}

	if m.GetCurrency() != USD {
		t.Errorf("currency = %v, want USD", m.GetCurrency())
	}
}

// Text written the way StringMoney writes it ("COP 1500.50") names its own
// currency, which wins over the receiver's.
func TestScanReadsCurrencyPrefix(t *testing.T) {
	for _, src := range []any{"COP 1500.50", []byte("COP 1500.50")} {
		m := NewFromDecimal(decimal.Zero, EUR)

		if err := m.Scan(src); err != nil {
			t.Fatalf("Scan(%#v) error = %v", src, err)
		}

		if m.GetCurrency() != COP || m.String() != "1500.5" {
			t.Errorf("Scan(%#v) = %s %v, want 1500.5 COP", src, m.String(), m.GetCurrency())
		}
	}
}

func TestScanRejectsUnknownCurrencyPrefix(t *testing.T) {
	m := NewFromDecimal(decimal.One, EUR)

	if err := m.Scan("ZZZ 10"); err == nil {
		t.Fatal("Scan(\"ZZZ 10\") returned no error")
	}

	if m.GetCurrency() != EUR || !m.GetDecimal().Equal(decimal.One) {
		t.Errorf("a failed Scan changed the receiver to %s %v", m.String(), m.GetCurrency())
	}
}

// StringMoney output scans back to the same amount and currency.
func TestStringMoneyRoundTripsThroughScan(t *testing.T) {
	original := MustMoneyFromString("-1234.56", COP)

	text, err := original.StringMoney()
	if err != nil {
		t.Fatal(err)
	}

	var scanned Money
	if err := scanned.Scan(text); err != nil {
		t.Fatalf("Scan(%q) error = %v", text, err)
	}

	if !scanned.Equal(original) {
		t.Errorf("Scan(%q) = %s %v, want %s %v", text, scanned, scanned.GetCurrency(), original, original.GetCurrency())
	}
}

// Format and StringMoney show the currency's own precision, rounding (half
// away from zero) when the amount carries more digits.
func TestFormatRoundsToCurrencyPrecision(t *testing.T) {
	tests := []struct {
		amount     string
		currency   Currency
		format     string
		stringForm string
	}{
		{"10.005", USD, "$10.01", "USD 10.01"},
		{"10.004", USD, "$10.00", "USD 10.00"},
		{"1500.5", JPY, "¥1501", "JPY 1501"},
		{"1.2345", KWD, "KWD1.235", "KWD 1.235"},
	}

	for _, tt := range tests {
		m := MustMoneyFromString(tt.amount, tt.currency)

		if got, err := m.Format(); err != nil || got != tt.format {
			t.Errorf("Format(%s %v) = %q, %v, want %q", tt.amount, tt.currency, got, err, tt.format)
		}

		if got, err := m.StringMoney(); err != nil || got != tt.stringForm {
			t.Errorf("StringMoney(%s %v) = %q, %v, want %q", tt.amount, tt.currency, got, err, tt.stringForm)
		}
	}
}

// The minus sign goes before the symbol, as amounts are written.
func TestFormatPutsSignBeforeSymbol(t *testing.T) {
	tests := []struct {
		amount string
		want   string
	}{
		{"-10", "-$10.00"},
		{"-0.004", "$0.00"},
	}

	for _, tt := range tests {
		got, err := MustMoneyFromString(tt.amount, USD).Format()
		if err != nil || got != tt.want {
			t.Errorf("Format(%s USD) = %q, %v, want %q", tt.amount, got, err, tt.want)
		}
	}
}

// Currency.Scan used to accept any type it did not recognise — NULL
// included — and silently set the currency to XXX.
func TestCurrencyScanRejectsUnsupportedTypes(t *testing.T) {
	for _, src := range []any{nil, int64(32), 3.5, true} {
		c := COP

		if err := c.Scan(src); err == nil {
			t.Errorf("Scan(%#v) returned no error", src)
		}

		if c != COP {
			t.Errorf("Scan(%#v) changed the currency to %v", src, c)
		}
	}
}

// Currencies introduced since the table was written are recognised, with
// their ISO 4217 minor unit.
func TestNewISOCurrencies(t *testing.T) {
	for _, code := range []string{"SLE", "ZWG", "XCG"} {
		c, err := GetCurrencyFromISOCode(code)
		if err != nil {
			t.Errorf("GetCurrencyFromISOCode(%q) error = %v", code, err)
			continue
		}

		if !c.Valid() || c.String() != code {
			t.Errorf("%s: Valid() = %v, String() = %q", code, c.Valid(), c.String())
		}

		if prec, err := c.GetCurrencyPrecisionCode(); err != nil || prec != 2 {
			t.Errorf("%s precision = %d, %v, want 2", code, prec, err)
		}
	}

	if SLE.String() != "SLE" || ZWG.String() != "ZWG" || XCG.String() != "XCG" {
		t.Error("new currency constants do not name their codes")
	}
}

// The codes they replace stay readable, so stored historical amounts still
// decode.
func TestReplacedISOCurrenciesStillDecode(t *testing.T) {
	for _, code := range []string{"HRK", "SLL", "ZWL", "ANG"} {
		if _, err := GetCurrencyFromISOCode(code); err != nil {
			t.Errorf("GetCurrencyFromISOCode(%q) error = %v", code, err)
		}
	}
}

// Valid follows the currency table rather than a hard-coded upper bound.
func TestValidMatchesCurrencyTable(t *testing.T) {
	for c := range currencyCode {
		if !c.Valid() {
			t.Errorf("%v is in the table but not Valid", c)
		}
	}

	if Currency(255).Valid() {
		t.Error("Currency(255) reported Valid")
	}
}
