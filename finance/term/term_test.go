package term

import (
	"errors"
	"testing"

	"github.com/yeferson59/gofinance/decimal"
)

func TestUnitValid(t *testing.T) {
	for _, u := range []Unit{Days, Weeks, Months, Years} {
		if !u.Valid() {
			t.Errorf("expected %q to be valid", u)
		}
	}

	if Unit("fortnights").Valid() {
		t.Error("expected an unknown unit to be invalid")
	}
}

func TestFrequencyPeriodsPerYear(t *testing.T) {
	cases := map[Frequency]int64{
		Daily:        365,
		Monthly:      12,
		Bimonthly:    6,
		Quarterly:    4,
		FourMonthly:  3,
		SemiAnnually: 2,
		Annually:     1,
	}

	for freq, want := range cases {
		got, err := freq.PeriodsPerYear()
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", freq, err)
		}
		if !got.Equal(decimal.MustFromInt64(want, 0)) {
			t.Errorf("%q: expected %d periods per year, got %s", freq, want, got)
		}
	}
}

func TestFrequencyMonthsPerPeriod(t *testing.T) {
	cases := map[Frequency]int64{
		Monthly:      1,
		Bimonthly:    2,
		Quarterly:    3,
		FourMonthly:  4,
		SemiAnnually: 6,
		Annually:     12,
	}

	for freq, want := range cases {
		got, err := freq.MonthsPerPeriod()
		if err != nil {
			t.Fatalf("%q: unexpected error: %v", freq, err)
		}
		if !got.Equal(decimal.MustFromInt64(want, 0)) {
			t.Errorf("%q: expected %d months per period, got %s", freq, want, got)
		}
	}
}

func TestFrequencyInvalid(t *testing.T) {
	bad := Frequency("weekly-ish")

	if bad.Valid() {
		t.Error("expected an unknown frequency to be invalid")
	}

	if _, err := bad.PeriodsPerYear(); !errors.Is(err, ErrInvalidFrequency) {
		t.Errorf("expected ErrInvalidFrequency, got %v", err)
	}

	if _, err := bad.MonthsPerPeriod(); !errors.Is(err, ErrInvalidFrequency) {
		t.Errorf("expected ErrInvalidFrequency, got %v", err)
	}
}
