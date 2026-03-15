package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func BenchmarkNewAnnuity(b *testing.B) {
	testcases := []struct {
		value        float64
		present      float64
		future       float64
		period       compositeinterest.Period
		rateInterest compositeinterest.RateInterest
	}{
		{
			0,
			1_000,
			0,
			compositeinterest.Period{},
			compositeinterest.RateInterest{},
		},
		{
			0,
			2_000,
			0,
			compositeinterest.Period{},
			compositeinterest.RateInterest{},
		},
	}

	for _, testcase := range testcases {
		b.ReportAllocs()
		b.StartTimer()
		for b.Loop() {
			valueMoney, _ := money.New(int64(testcase.value*100), 2, money.USD)
			presentMoney, _ := money.New(int64(testcase.present*100), 2, money.USD)
			futureMoney, _ := money.New(int64(testcase.future*100), 2, money.USD)
			_, _ = annuities.New(valueMoney, presentMoney, futureMoney, testcase.period, testcase.rateInterest)
		}
	}
}
