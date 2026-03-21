package benchmarks

import (
	"os"
	"testing"

	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func BenchmarkNewAnnuity(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.05), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuities.New(
			money.MustMoneyFromFloat64(100, money.USD),
			money.MustMoneyFromFloat64(1000, money.USD),
			money.MustMoneyFromFloat64(0, money.USD),
			period,
			rate,
		)
	}
}

func BenchmarkAnnuityPresent(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.05), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	annuity, _ := annuities.New(
		money.MustMoneyFromFloat64(100, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		period,
		rate,
	)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuity.Present()
	}
}

func BenchmarkAnnuityFuture(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(12), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.05), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	annuity, _ := annuities.New(
		money.MustMoneyFromFloat64(100, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD),
		period,
		rate,
	)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuity.Future()
	}
}

func BenchmarkAnnuityPaymentFromPresentValue(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(360), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.06), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	annuity, _ := annuities.New(
		money.MustMoneyFromFloat64(0, money.USD),
		money.MustMoneyFromFloat64(300000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		period,
		rate,
	)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuity.PaymentFromPresentValue()
	}
}

func BenchmarkAnnuityPaymentFromFutureValue(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(120), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.08), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	annuity, _ := annuities.New(
		money.MustMoneyFromFloat64(0, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		money.MustMoneyFromFloat64(50000, money.USD),
		period,
		rate,
	)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuity.PaymentFromFutureValue()
	}
}

func BenchmarkBuildSchedule(b *testing.B) {
	pv := money.MustMoneyFromFloat64(200000, money.USD)
	rate := money.MustMoneyFromFloat64(0.005, money.USD)
	payment := money.MustMoneyFromFloat64(1074, money.USD)

	testcases := []struct {
		name string
		nper int
	}{
		{"12_months", 12},
		{"120_months", 120},
		{"360_months", 360},
		{"600_months", 600},
	}

	for _, tc := range testcases {
		b.Run(tc.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = annuities.BuildSchedule(pv, rate, payment, tc.nper)
			}
		})
	}
}

func BenchmarkWriteCSV(b *testing.B) {
	schedule := annuities.BuildSchedule(
		money.MustMoneyFromFloat64(200000, money.USD),
		money.MustMoneyFromFloat64(0.005, money.USD),
		money.MustMoneyFromFloat64(1074, money.USD),
		360,
	)

	headers := []string{"Period", "Balance", "Payment", "Interest", "SumInterest", "Principal"}

	tmpFile := "/tmp/benchmark_test.csv"

	testcases := []struct {
		name string
		rows int
	}{
		{"360_rows", 360},
		{"1200_rows", 1200},
		{"3600_rows", 3600},
	}

	for _, tc := range testcases {
		b.Run(tc.name, func(b *testing.B) {
			extendedRows := make([]annuities.Schedule, tc.rows)
			copy(extendedRows, schedule)
			for i := len(schedule); i < tc.rows; i++ {
				extendedRows[i] = schedule[i%len(schedule)]
			}

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				_ = annuities.WriteCSV(tmpFile, headers, extendedRows)
			}
		})
	}

	os.Remove(tmpFile)
}

func BenchmarkAnnuityAllMethods(b *testing.B) {
	period, _ := compositeinterest.NewPeriod(money.MustFromFloat64(120), compositeinterest.Monthly)
	rate, _ := compositeinterest.NewRateInterest(money.MustFromFloat64(0.05), compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
	annuity, _ := annuities.New(
		money.MustMoneyFromFloat64(100, money.USD),
		money.MustMoneyFromFloat64(1000, money.USD),
		money.MustMoneyFromFloat64(0, money.USD),
		period,
		rate,
	)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		_, _ = annuity.Present()
		_, _ = annuity.Future()
		_, _ = annuity.PaymentFromPresentValue()
		_, _ = annuity.PaymentFromFutureValue()
		_ = annuities.BuildSchedule(
			money.MustMoneyFromFloat64(1000, money.USD),
			money.MustMoneyFromFloat64(0.05/12, money.USD),
			money.MustMoneyFromFloat64(100, money.USD),
			120,
		)
	}
}
