package tvm_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/tvm"
)

func ExampleConfig_SolvePMT() {
	// A $300,000 mortgage at 0.5% a month over 30 years. The payment comes
	// back negative: money leaving the borrower.
	payment, err := tvm.NewTVM().PV(300000).Rate(0.005).N(360).SolvePMT()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", payment.InexactFloat64())
	// Output: -1798.65
}

func ExampleConfig_SolveFV() {
	// Saving $100 a month at 1% a month for a year.
	future, err := tvm.NewTVM().PMT(-100).Rate(0.01).N(12).SolveFV()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f\n", future.InexactFloat64())
	// Output: 1268.25
}

func ExampleConfig_SolveN() {
	// How long until $100 monthly deposits at 1% reach $5,000?
	periods, err := tvm.NewTVM().PMT(-100).Rate(0.01).FV(5000).SolveN()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f periods\n", periods.InexactFloat64())
	// Output: 40.75 periods
}

func ExampleConfig_SolveRate() {
	// Recover the rate implied by a known loan and payment.
	rate, err := tvm.NewTVM().PV(300000).PMT(-1798.65).N(360).SolveRate()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.4f per period\n", rate.InexactFloat64())
	// Output: 0.0050 per period
}

func ExampleConfig_Due() {
	// Paying at the start of each period rather than the end services less
	// interest, so the payment is smaller.
	ordinary, err := tvm.NewTVM().PV(100000).Rate(0.01).N(60).SolvePMT()
	if err != nil {
		panic(err)
	}

	due, err := tvm.NewTVM().PV(100000).Rate(0.01).N(60).Due().SolvePMT()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ordinary %.2f, due %.2f\n", ordinary.InexactFloat64(), due.InexactFloat64())
	// Output: ordinary -2224.44, due -2202.42
}
