package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/decimal"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
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
			_, _ = compositeinterest.NewPeriod(decimal.MustFromFloat64(value), compositeinterest.Monthly)
		}
	}
}

func BenchmarkNewRateInterest(b *testing.B) {
	testcases := []struct {
		value                float64
		compoundingFrequency compositeinterest.CompoundingFrequency
		typeRate             compositeinterest.TypeRate
	}{
		{
			0.10,
			compositeinterest.Annually,
			compositeinterest.RateEffectyAnnually,
		},
		{
			0.24,
			compositeinterest.Bimonthly,
			compositeinterest.RateAnticipateEffectyNominal,
		},
		{
			0.25,
			compositeinterest.Annually,
			compositeinterest.RateAnticipateEffectyPeriodic,
		},
		{
			0.244,
			compositeinterest.Daily,
			compositeinterest.RateAnticipateEffectyNominal,
		},
	}

	for _, testcase := range testcases {
		b.StartTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, _ = compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.value), testcase.compoundingFrequency, testcase.typeRate)
		}
	}
}

func BenchmarkRateInterest(b *testing.B) {
	rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(0.25), compositeinterest.Bimonthly, compositeinterest.RateEffectyAnnually)
	rateTwo, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(0.18), compositeinterest.Monthly, compositeinterest.RateEffectyAnnually)

	testcases := []compositeinterest.RateInterest{
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
				_, _ = testcase.RateNominalToNominal(compositeinterest.Annually)
			}
		}
	})

	b.Run("rate periodic to periodic", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				_, _ = testcase.RatePeriodicToPeriodic(compositeinterest.Bimonthly)
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

func BenchmarkNewCompositeInterest(b *testing.B) {
	testcases := []struct {
		numberPeriod         float64
		valueInterest        float64
		compoundingFrequency compositeinterest.CompoundingFrequency
		typeRate             compositeinterest.TypeRate
		present              float64
		future               float64
	}{
		{
			2,
			0.24,
			compositeinterest.Monthly,
			compositeinterest.RateEffectyNominal,
			1_000,
			0,
		},
		{
			2.5,
			0.24,
			compositeinterest.Monthly,
			compositeinterest.RateEffectyNominal,
			0,
			2_500,
		},
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
			rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
			presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
			futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
			_, _ = compositeinterest.New(presentMoney, futureMoney, rate, period)
		}
	}
}

func BenchmarkCompositeInterest(b *testing.B) {
	testcases := []struct {
		numberPeriod         float64
		valueInterest        float64
		compoundingFrequency compositeinterest.CompoundingFrequency
		typeRate             compositeinterest.TypeRate
		present              float64
		future               float64
	}{
		{
			2,
			0.25,
			compositeinterest.Monthly,
			compositeinterest.RateEffectyNominal,
			1_000,
			0,
		},
		{
			2.5,
			0.24,
			compositeinterest.Monthly,
			compositeinterest.RateEffectyNominal,
			0,
			2_500,
		},
	}

	b.Run("future", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compositeinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Future()
			}
		}
	})

	b.Run("present", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compositeinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Present()
			}
		}
	})

	b.Run("get equals rate interest periods", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compositeinterest.New(presentMoney, futureMoney, rate, period)

				_, _, _ = ci.GetEqualsRateInterestPeriods()
			}
		}
	})

	b.Run("interest", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compositeinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Interest()
			}
		}
	})

	b.Run("periods", func(b *testing.B) {
		for _, testcase := range testcases {
			b.ReportAllocs()
			b.StartTimer()
			for b.Loop() {
				period, _ := compositeinterest.NewPeriod(decimal.MustFromFloat64(testcase.numberPeriod), testcase.compoundingFrequency)
				rate, _ := compositeinterest.NewRateInterest(decimal.MustFromFloat64(testcase.valueInterest), testcase.compoundingFrequency, testcase.typeRate)
				presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
				futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
				ci, _ := compositeinterest.New(presentMoney, futureMoney, rate, period)

				_, _ = ci.Periods()
			}
		}
	})
}
