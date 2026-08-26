package annuities

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

func TestBuildScheduleAmortizesToZero(t *testing.T) {
	pv := money.MustMoneyFromFloat64(200000, money.USD)
	rate := decimal.MustFromFloat64(0.005)
	nper := decimal.MustFromFloat64(360)

	payment := NewAnnuity().
		Present(200000, money.USD).
		Rate(0.005).
		Periods(360).
		Monthly().
		MustPayment()

	rows, err := BuildSchedule(pv, rate, payment, nper)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rows) != 361 {
		t.Fatalf("expected 361 rows (period 0..360), got %d", len(rows))
	}

	first := rows[0]
	if !first.Balance.Equal(pv) {
		t.Errorf("expected period 0 balance to equal pv %s, got %s", pv.String(), first.Balance.String())
	}
	if first.Payment.GetCurrency() != pv.GetCurrency() {
		t.Errorf("expected period 0 payment currency to match pv currency %v, got %v", pv.GetCurrency(), first.Payment.GetCurrency())
	}

	last := rows[len(rows)-1]
	// The balance should be fully amortized (allowing for the last
	// payment's rounding to the currency's smallest unit).
	if last.Balance.Abs().InexactFloat64() > 1 {
		t.Errorf("expected balance to amortize to ~0, got %s", last.Balance.String())
	}
}

func TestBuildScheduleCurrencyMismatch(t *testing.T) {
	pv := money.MustMoneyFromFloat64(1000, money.USD)
	payment := money.MustMoneyFromFloat64(100, money.EUR)

	_, err := BuildSchedule(pv, decimal.MustFromFloat64(0.01), payment, decimal.MustFromFloat64(12))
	if !errors.Is(err, money.ErrCurrencyMismatch) {
		t.Errorf("expected ErrCurrencyMismatch, got %v", err)
	}
}

func TestBuildScheduleInvalidPeriods(t *testing.T) {
	pv := money.MustMoneyFromFloat64(1000, money.USD)
	payment := money.MustMoneyFromFloat64(100, money.USD)

	tests := []decimal.Decimal{
		decimal.MustFromFloat64(0),
		decimal.MustFromFloat64(-12),
	}

	for _, nper := range tests {
		if _, err := BuildSchedule(pv, decimal.MustFromFloat64(0.01), payment, nper); !errors.Is(err, ErrInvalidPeriods) {
			t.Errorf("expected ErrInvalidPeriods for nper=%s, got %v", nper.String(), err)
		}
	}
}

func TestBuildScheduleNonUSDCurrency(t *testing.T) {
	pv := money.MustMoneyFromFloat64(1000, money.JPY)
	payment := money.MustMoneyFromFloat64(90, money.JPY)

	rows, err := BuildSchedule(pv, decimal.MustFromFloat64(0.01), payment, decimal.MustFromFloat64(12))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	for i, r := range rows {
		if r.Balance.GetCurrency() != money.JPY {
			t.Errorf("row %d: expected balance currency JPY, got %v", i, r.Balance.GetCurrency())
		}
		if r.SumInterest.GetCurrency() != money.JPY {
			t.Errorf("row %d: expected sum-interest currency JPY, got %v", i, r.SumInterest.GetCurrency())
		}
	}
}

func TestWriteCSVRoundsToTargetCurrencyPrecision(t *testing.T) {
	pv := money.MustMoneyFromFloat64(1000, money.JPY)
	payment := money.MustMoneyFromFloat64(90, money.JPY)

	rows, err := BuildSchedule(pv, decimal.MustFromFloat64(0.01), payment, decimal.MustFromFloat64(12))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	path := filepath.Join(t.TempDir(), "schedule.csv")
	headers := []string{"Period", "Balance", "Payment", "Interest", "SumInterest", "Principal"}

	if err := WriteCSV(path, headers, rows); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	content := string(data)
	if len(content) == 0 {
		t.Fatal("expected non-empty CSV output")
	}
	// JPY has zero decimal places, so no row should contain a decimal point.
	for i, line := range strings.Split(content, "\n") {
		if i == 0 || line == "" {
			continue // header or trailing newline
		}
		if strings.Contains(line, ".") {
			t.Errorf("expected no decimal point in JPY row, got %q", line)
		}
	}
}
