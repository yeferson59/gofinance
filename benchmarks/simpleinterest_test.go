package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/finance/simpleinterest"
)

func BenchmarkNewPeriod(b *testing.B) {
	testcases := []struct {
		value float64
		time  simpleinterest.Periods
	}{
		{
			1,
			simpleinterest.Days,
		},
		{
			2,
			simpleinterest.Weeks,
		},
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			_ = simpleinterest.NewPeriod(testcase.value, testcase.time)
		}
	}
}

func BenchmarkNewSimpleInterest(b *testing.B) {
	testcases := []simpleinterest.Period{
		simpleinterest.NewPeriod(1, simpleinterest.Days),
		simpleinterest.NewPeriod(2, simpleinterest.Weeks),
		simpleinterest.NewPeriod(0.5, simpleinterest.Years),
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			_ = simpleinterest.New(0, 1_000, 0, 0.5, testcase)
		}
	}
}

func BenchmarkSimpleInterest(b *testing.B) {
	testcases := []simpleinterest.SimpleInterest{
		simpleinterest.New(0, 1_000, 200, 0.9, simpleinterest.NewPeriod(1, simpleinterest.Days)),
		simpleinterest.New(0, 1_500, 100, 0.10, simpleinterest.NewPeriod(2.9, simpleinterest.Months)),
		simpleinterest.New(0, 1_200, 20, 0.6, simpleinterest.NewPeriod(0.9, simpleinterest.Years)),
		simpleinterest.New(0, 1_000, 250, 0.1, simpleinterest.NewPeriod(1.5, simpleinterest.Weeks)),
		simpleinterest.New(0, 1_000, 80, 0.25, simpleinterest.NewPeriod(2.5, simpleinterest.Days)),
	}

	b.Run("future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.Future()
			}
		}
	})

	b.Run("future with rate interest", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.FutureWithRateInterest()
			}
		}
	})

	b.Run("present", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.Present()
			}
		}
	})

	b.Run("present with future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.PresentWithFuture()
			}
		}
	})

	b.Run("interest", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.Interest()
			}
		}
	})

	b.Run("rate interest", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateInterest()
			}
		}
	})

	b.Run("rate interest with present and future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateInterestWithPresentAndFuture()
			}
		}
	})

	b.Run("periods", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.Periods()
			}
		}
	})

	b.Run("periods with present and future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.PeriodsWithPresentAndFuture()
			}
		}
	})
}
