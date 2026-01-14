package benchmarks

import (
	"testing"

	"github.com/yeferson59/gofinance/finance/annuities"
	"github.com/yeferson59/gofinance/finance/compositeinterest"
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
			_, _ = annuities.New(testcase.value, testcase.present, testcase.future, testcase.period, testcase.rateInterest)
		}
	}
}
