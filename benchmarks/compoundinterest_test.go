package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

func BenchmarkPeriod(b *testing.B) {
	testcases := []float64{
		6,
		12,
		18,
		24,
		36,
		48,
		0,
	}

	for _, value := range testcases {
		b.StartTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, _ = compoundinterest.NewPeriod(decimal.MustFromFloat64(value), compoundinterest.Monthly)
		}
	}
}

func BenchmarkNewRateInterest(b *testing.B) {
	testcases := []struct {
		value                float64
		compoundingFrequency compoundinterest.CompoundingFrequency
		typeRate             compoundinterest.TypeRate
	}{
		{
			0.10,
			compoundinterest.Annually,
			compoundinterest.RateEffectyAnnually,
		},
		{
			0.24,
			compoundinterest.Bimonthly,
			compoundinterest.RateAnticipateEffectyNominal,
		},
		{
			0.25,
			compoundinterest.Annually,
			compoundinterest.RateAnticipateEffectyPeriodic,
		},
		{
			0.244,
			compoundinterest.Daily,
			compoundinterest.RateAnticipateEffectyNominal,
		},
	}

	for _, testcase := range testcases {
		b.StartTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, _ = compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.value), testcase.compoundingFrequency, testcase.typeRate)
		}
	}
}

func BenchmarkRateInterest(b *testing.B) {
	rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.25), compoundinterest.Bimonthly, compoundinterest.RateEffectyAnnually)
	rateTwo, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(0.18), compoundinterest.Monthly, compoundinterest.RateEffectyAnnually)

	testcases := []compoundinterest.RateInterest{
		rate,
		rateTwo,
	}

	b.Run("rate periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RatePeriodic()
			}
		}
	})

	b.Run("rate nominal", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateNominal()
			}
		}
	})

	b.Run("rate effecty annually", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateEffectyAnnually()
			}
		}
	})

	b.Run("rate anticipated periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateAnticipatePeriodic()
			}
		}
	})

	b.Run("rate anticipated nominal", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateAnticipateNominal()
			}
		}
	})

	b.Run("rate anticipated effecty annually", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateAnticipateEffectyAnnually()
			}
		}
	})

	b.Run("rate nominal to nominal", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RateNominalToNominal(compoundinterest.Annually)
			}
		}
	})

	b.Run("rate periodic to periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RatePeriodicToPeriodic(compoundinterest.Bimonthly)
			}
		}
	})

	b.Run("to nominal", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.ToNominal()
			}
		}
	})

	b.Run("to periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.ToPeriodic()
			}
		}
	})

	b.Run("to anticipated nominal", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.ToAnticipateNominal()
			}
		}
	})

	b.Run("to anticipated periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.ToAnticipatePeriodic()
			}
		}
	})
}

func BenchmarkNewCompoundInterest(b *testing.B) {
	testcases := []struct {
		numberPeriod         float64
		valueInterest        float64
		compoundingFrequency compoundinterest.CompoundingFrequency
		typeRate             compoundinterest.TypeRate
		present              float64
		future               float64
	}{
		{
			2,
			0.24,
			compoundinterest.Monthly,
			compoundinterest.RateEffectyNominal,
			1_000,
			0,
		},
		{
			2.5,
			0.24,
			compoundinterest.Monthly,
			compoundinterest.RateEffectyNominal,
			0,
			2_500,
		},
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
			rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
			presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
			futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
			_, _ = compoundinterest.New(presentMoney, futureMoney, rate, period)
		}
	}
}

func BenchmarkCompoundInterest(b *testing.B) {
	testcases := []struct {
		numberPeriod         float64
		valueInterest        float64
		compoundingFrequency compoundinterest.CompoundingFrequency
		typeRate             compoundinterest.TypeRate
		present              float64
		future               float64
	}{
		{
			2,
			0.25,
			compoundinterest.Monthly,
			compoundinterest.RateEffectyNominal,
			1_000,
			0,
		},
		{
			2.5,
			0.24,
			compoundinterest.Monthly,
			compoundinterest.RateEffectyNominal,
			0,
			2_500,
		},
	}

	b.Run("future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compoundinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Future()
			}
		}
	})

	b.Run("present", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compoundinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Present()
			}
		}
	})

	b.Run("get equals rate interest periods", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compoundinterest.New(presentMoney, futureMoney, rate, period)

				_, _, _ = ci.GetEqualsRateInterestPeriods()
			}
		}
	})

	b.Run("interest", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compoundinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Interest()
			}
		}
	})

	b.Run("periods", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compoundinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compoundinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compoundinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Periods()
			}
		}
	})
}
