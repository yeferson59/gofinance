// Package loans looks at an amortizing loan from the borrower's side: the level
// payment, the true annual cost once up-front fees are taken into account
// (APR), what overpaying every period saves, and whether a refinance offer is
// worth its closing costs.
//
// A loan is described by a principal, a periodic (or nominal annual) interest
// rate, a payment frequency, and either a number of payments or a number of
// years. Fees and extra payments are expressed as plain amounts in the
// principal's currency. Every calculation runs on the decimal engine, so the
// closed-form results are exact to the engine's precision and the APR search
// bisects against exact residuals.
//
// Example — a $250,000 mortgage at 6.5% nominal annual over 30 years, with
// $3,500 of up-front fees:
//
//	loan := loans.NewLoan().
//	    Principal(250000, money.USD).
//	    AnnualRate(0.065).
//	    Years(30).
//	    Monthly().
//	    Fees(3500)
//
//	payment := loan.MustPayment() // ≈ 1580.17 per month
//	apr := loan.MustAPR()         // ≈ 0.0664 — above the 6.5% note rate
package loans

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/money"
)

var (
	// ErrNonPositivePrincipal is returned when the amount borrowed is zero or
	// negative.
	ErrNonPositivePrincipal = errors.New("loans: principal must be positive")

	// ErrInvalidPeriods is returned when neither Periods nor Years resolves to
	// a positive number of payments.
	ErrInvalidPeriods = errors.New("loans: number of periods must be positive")

	// ErrInvalidFrequency is returned when the number of payments per year is
	// not at least one.
	ErrInvalidFrequency = errors.New("loans: payments per year must be at least 1")

	// ErrInvalidRate is returned when a rate is not greater than −1, which
	// makes the growth factor (1+rate) zero or negative.
	ErrInvalidRate = errors.New("loans: rate must be greater than -1")

	// ErrNegativeFees is returned when the up-front fees are negative.
	ErrNegativeFees = errors.New("loans: fees cannot be negative")

	// ErrFeesExceedPrincipal is returned when the fees swallow the whole
	// principal, leaving the borrower with nothing to finance.
	ErrFeesExceedPrincipal = errors.New("loans: fees must be smaller than the principal")

	// ErrNegativeExtra is returned when the extra payment applied on top of the
	// scheduled one is negative.
	ErrNegativeExtra = errors.New("loans: extra payment cannot be negative")

	// ErrFrequencyMismatch is returned by Compare when the two loans pay on
	// different schedules, which makes a period-by-period comparison
	// meaningless.
	ErrFrequencyMismatch = errors.New("loans: loans must share the same payments per year")

	// ErrNoBreakEven is returned by Compare when the offer's savings never add
	// up to its closing costs.
	ErrNoBreakEven = errors.New("loans: the offer never recovers its closing costs")

	// ErrNoConvergence is returned by the APR solver when no rate reproduces
	// the loan's net proceeds.
	ErrNoConvergence = errors.New("loans: rate did not converge")
)

// Config is a fluent builder describing an amortizing loan. Create one with
// NewLoan, set its terms, then call Payment, APR, Payoff, Savings, or pass it
// to Compare.
type Config struct {
	principal money.Money
	rate      decimal.Decimal
	annual    bool
	periods   int
	years     int
	freq      int
	fees      decimal.Decimal
	extra     decimal.Decimal
}

// NewLoan returns a Config defaulting to monthly payments (12 per year) with
// no fees and no extra payment.
func NewLoan() Config {
	return Config{
		rate:  decimal.Zero,
		freq:  12,
		fees:  decimal.Zero,
		extra: decimal.Zero,
	}
}

// Principal sets the amount borrowed and the currency every other amount of
// the loan is expressed in.
func (l Config) Principal(amount float64, currency money.Currency) Config {
	l.principal = money.MustMoneyFromFloat64(amount, currency)
	return l
}

// Rate sets the interest rate charged each payment period as a fraction (e.g.
// 0.005 for 0.5% per month). It replaces any rate set with AnnualRate.
func (l Config) Rate(rate float64) Config {
	l.rate = decimal.MustFromFloat64(rate)
	l.annual = false

	return l
}

// AnnualRate sets the nominal annual rate as a fraction (e.g. 0.065 for 6.5%).
// It is divided by the payments per year to obtain the periodic rate, whatever
// order the builder methods are called in. It replaces any rate set with Rate.
func (l Config) AnnualRate(rate float64) Config {
	l.rate = decimal.MustFromFloat64(rate)
	l.annual = true

	return l
}

// Periods sets the total number of payments. It takes precedence over Years.
func (l Config) Periods(n int) Config {
	l.periods = n
	return l
}

// Years sets the term in years; the number of payments is that many times the
// payments per year, whatever order the builder methods are called in.
func (l Config) Years(n int) Config {
	l.years = n
	return l
}

// PaymentsPerYear sets how many payments are made each year (12 by default).
func (l Config) PaymentsPerYear(n int) Config {
	l.freq = n
	return l
}

// Monthly sets twelve payments per year. This is the default.
func (l Config) Monthly() Config {
	l.freq = 12
	return l
}

// Quarterly sets four payments per year.
func (l Config) Quarterly() Config {
	l.freq = 4
	return l
}

// SemiAnnually sets two payments per year.
func (l Config) SemiAnnually() Config {
	l.freq = 2
	return l
}

// Annually sets one payment per year.
func (l Config) Annually() Config {
	l.freq = 1
	return l
}

// Fees sets the up-front fees, points, and other finance charges the borrower
// pays to obtain the loan, in the principal's currency. They are deducted from
// the cash actually received, which is what makes the APR exceed the note rate.
func (l Config) Fees(amount float64) Config {
	l.fees = decimal.MustFromFloat64(amount)
	return l
}

// ExtraPayment sets an additional amount paid on top of the scheduled payment
// every period, in the principal's currency. It is what Payoff amortizes and
// what Savings measures against the scheduled payment.
func (l Config) ExtraPayment(amount float64) Config {
	l.extra = decimal.MustFromFloat64(amount)
	return l
}

// PeriodicRate returns the interest rate charged each payment period, dividing
// a nominal annual rate by the payments per year when one was set.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, ErrNonPositivePrincipal,
// or ErrInvalidRate on invalid terms.
func (l Config) PeriodicRate() (decimal.Decimal, error) {
	rate, _, err := l.terms()

	return rate, err
}

// MustPeriodicRate is like PeriodicRate but panics on error.
func (l Config) MustPeriodicRate() decimal.Decimal {
	d, err := l.PeriodicRate()
	if err != nil {
		panic(err)
	}

	return d
}

// NumberOfPayments returns how many payments the loan schedules: the explicit
// Periods when set, otherwise the years times the payments per year.
//
// It returns ErrInvalidFrequency, ErrInvalidPeriods, ErrNonPositivePrincipal,
// or ErrInvalidRate on invalid terms.
func (l Config) NumberOfPayments() (int, error) {
	_, n, err := l.terms()

	return n, err
}

// terms validates the configuration and returns the periodic rate together
// with the total number of payments.
func (l Config) terms() (decimal.Decimal, int, error) {
	if l.freq < 1 {
		return decimal.Decimal{}, 0, ErrInvalidFrequency
	}

	n := l.periods
	if n == 0 {
		n = l.years * l.freq
	}

	if n <= 0 {
		return decimal.Decimal{}, 0, ErrInvalidPeriods
	}

	if !l.principal.IsPositive() {
		return decimal.Decimal{}, 0, ErrNonPositivePrincipal
	}

	rate := l.rate

	if l.annual {
		periodsPerYear, err := decimal.NewFromInt64(int64(l.freq), 0)
		if err != nil {
			return decimal.Decimal{}, 0, err
		}

		rate, err = l.rate.Div(periodsPerYear)
		if err != nil {
			return decimal.Decimal{}, 0, err
		}
	}

	if !decimal.One.Add(rate).IsPos() {
		return decimal.Decimal{}, 0, ErrInvalidRate
	}

	return rate, n, nil
}

// periodsPerYear returns the payment frequency as a decimal, for annualizing
// periodic rates.
func (l Config) periodsPerYear() (decimal.Decimal, error) {
	if l.freq < 1 {
		return decimal.Decimal{}, ErrInvalidFrequency
	}

	return decimal.NewFromInt64(int64(l.freq), 0)
}

// annuityFactor returns the present value of one unit paid at the end of each
// of n periods, discounted at the periodic rate:
//
//	a(i, n) = (1 − (1+i)^−n) / i    (n when i is zero)
func annuityFactor(rate decimal.Decimal, n int) (decimal.Decimal, error) {
	periods, err := decimal.NewFromInt64(int64(n), 0)
	if err != nil {
		return decimal.Decimal{}, err
	}

	if rate.IsZero() {
		return periods, nil
	}

	growth := decimal.One.Add(rate)
	if !growth.IsPos() {
		return decimal.Decimal{}, ErrInvalidRate
	}

	pow, err := growth.Pow(periods)
	if err != nil {
		return decimal.Decimal{}, err
	}

	discount, err := decimal.One.Div(pow)
	if err != nil {
		return decimal.Decimal{}, err
	}

	numerator, err := decimal.One.TrySub(discount)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return numerator.Div(rate)
}

// levelPayment returns the payment that amortizes pv over n periods at the
// given periodic rate: PMT = PV / a(i, n).
func levelPayment(pv money.Money, rate decimal.Decimal, n int) (money.Money, error) {
	factor, err := annuityFactor(rate, n)
	if err != nil {
		return money.Money{}, err
	}

	return pv.DivDecimal(factor)
}
