package loans_test

import (
	"fmt"

	"github.com/yeferson59/gofinance/v2/finance/loans"
	"github.com/yeferson59/gofinance/v2/money"
)

func ExampleConfig_Payment() {
	payment, err := loans.NewLoan().
		Principal(300000, money.USD).
		AnnualRate(0.06).
		Years(30).
		Monthly().
		Payment()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%.2f a month\n", payment.InexactFloat64())
	// Output: 1798.65 a month
}

func ExampleConfig_APR() {
	// Fees reduce the cash actually received without reducing the payments,
	// so the APR comes out above the note rate.
	loan := loans.NewLoan().
		Principal(300000, money.USD).
		AnnualRate(0.06).
		Years(30).
		Monthly()

	withoutFees, err := loan.APR()
	if err != nil {
		panic(err)
	}

	withFees, err := loan.Fees(6000).APR()
	if err != nil {
		panic(err)
	}

	fmt.Printf("note rate %.4f, APR with fees %.4f\n",
		withoutFees.InexactFloat64(), withFees.InexactFloat64())
	// Output: note rate 0.0600, APR with fees 0.0619
}

func ExampleConfig_Savings() {
	// Paying an extra 200 a month retires the loan sooner and saves interest.
	savings, err := loans.NewLoan().
		Principal(300000, money.USD).
		AnnualRate(0.06).
		Years(30).
		Monthly().
		ExtraPayment(200).
		Savings()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d payments avoided, %.2f saved\n",
		savings.PeriodsSaved, savings.InterestSaved.InexactFloat64())
	// Output: 81 payments avoided, 91173.43 saved
}

func ExampleConfig_Quarterly() {
	// The payment frequency changes both the number of payments and the
	// periodic rate derived from the annual one.
	quarterly, err := loans.NewLoan().
		Principal(50000, money.USD).
		AnnualRate(0.08).
		Years(5).
		Quarterly().
		Payment()
	if err != nil {
		panic(err)
	}

	semiAnnually, err := loans.NewLoan().
		Principal(50000, money.USD).
		AnnualRate(0.08).
		Years(5).
		SemiAnnually().
		Payment()
	if err != nil {
		panic(err)
	}

	annually, err := loans.NewLoan().
		Principal(50000, money.USD).
		AnnualRate(0.08).
		Years(5).
		Annually().
		Payment()
	if err != nil {
		panic(err)
	}

	fmt.Printf("quarterly %.2f, semiannually %.2f, annually %.2f\n",
		quarterly.InexactFloat64(), semiAnnually.InexactFloat64(), annually.InexactFloat64())
	// Output: quarterly 3057.84, semiannually 6164.55, annually 12522.82
}
