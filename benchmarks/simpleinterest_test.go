package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/finance/simpleinterest"
)

func mustDecimal(f float64) simpleinterest.Decimal {
	d, _ := simpleinterest.NewFromFloat64(f)
	return d
}

func mustDecimalInt(i int64) simpleinterest.Decimal {
	d, _ := simpleinterest.NewFromInt64(i, 0)
	return d
}

func BenchmarkNewPeriod(b *testing.B) {
	testcases := []struct {
		value simpleinterest.Decimal
		time  simpleinterest.Periods
	}{
		{
			mustDecimal(1),
			simpleinterest.Days,
		},
		{
			mustDecimal(2),
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
		simpleinterest.NewPeriod(mustDecimal(1), simpleinterest.Days),
		simpleinterest.NewPeriod(mustDecimal(2), simpleinterest.Weeks),
		simpleinterest.NewPeriod(mustDecimal(0.5), simpleinterest.Years),
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			_ = simpleinterest.New(
				mustDecimalInt(0),
				mustDecimalInt(1_000),
				mustDecimalInt(0),
				mustDecimal(0.5),
				testcase,
			)
		}
	}
}

func BenchmarkSimpleInterest(b *testing.B) {
	testcases := []simpleinterest.SimpleInterest{
		simpleinterest.New(
			mustDecimalInt(0),
			mustDecimalInt(1_000),
			mustDecimalInt(200),
			mustDecimal(0.9),
			simpleinterest.NewPeriod(mustDecimal(1), simpleinterest.Days),
		),
		simpleinterest.New(
			mustDecimalInt(0),
			mustDecimalInt(1_500),
			mustDecimalInt(100),
			mustDecimal(0.10),
			simpleinterest.NewPeriod(mustDecimal(2.9), simpleinterest.Months),
		),
		simpleinterest.New(
			mustDecimalInt(0),
			mustDecimalInt(1_200),
			mustDecimalInt(20),
			mustDecimal(0.6),
			simpleinterest.NewPeriod(mustDecimal(0.9), simpleinterest.Years),
		),
		simpleinterest.New(
			mustDecimalInt(0),
			mustDecimalInt(1_000),
			mustDecimalInt(250),
			mustDecimal(0.1),
			simpleinterest.NewPeriod(mustDecimal(1.5), simpleinterest.Weeks),
		),
		simpleinterest.New(
			mustDecimalInt(0),
			mustDecimalInt(1_000),
			mustDecimalInt(80),
			mustDecimal(0.25),
			simpleinterest.NewPeriod(mustDecimal(2.5), simpleinterest.Days),
		),
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
