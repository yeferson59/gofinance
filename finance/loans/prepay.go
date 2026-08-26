package loans

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/annuities"
	"github.com/yeferson59/gofinance/v2/money"
)

// Payoff describes how a loan is actually retired: how many payments it takes,
// what they add up to, and the period-by-period amortization behind them. The
// rows are annuities.Schedule values, so they can be exported with
// annuities.WriteCSVTo or plotted like any other schedule.
type Payoff struct {
	// Periods is the number of payments actually made, which is shorter than
	// the scheduled term whenever an extra payment is applied.
	Periods int

	// Payment is the scheduled level payment at the note rate.
	Payment money.Money

	// ExtraPayment is the amount added to every scheduled payment.
	ExtraPayment money.Money

	// FinalPayment is the last payment, which only settles the outstanding
	// balance and is therefore usually smaller than Payment + ExtraPayment.
	FinalPayment money.Money

	// TotalPaid is the sum of every payment made.
	TotalPaid money.Money

	// TotalInterest is the part of TotalPaid that went to interest.
	TotalInterest money.Money

	// Schedule is the amortization table, opening with a period-0 row that
	// carries the original balance.
	Schedule []annuities.Schedule
}

// Savings compares retiring the loan on schedule against retiring it with the
// configured extra payment.
type Savings struct {
	// Scheduled is the payoff with the contractual payment alone.
	Scheduled Payoff

	// Accelerated is the payoff with the extra payment applied every period.
	Accelerated Payoff

	// PeriodsSaved is how many payments the borrower avoids.
	PeriodsSaved int

	// InterestSaved is the interest never charged because the balance is
	// retired sooner.
	InterestSaved money.Money
}

// Payoff amortizes the loan with the configured extra payment (zero unless
// ExtraPayment was set), stopping as soon as the balance reaches zero. Paying
// more than the schedule demands retires the principal early, so the last
// payment only settles what is left.
//
// It returns the term errors of PeriodicRate and ErrNegativeExtra if the extra
// payment is negative.
func (l Config) Payoff() (Payoff, error) {
	rate, n, err := l.terms()
	if err != nil {
		return Payoff{}, err
	}

	payment, err := levelPayment(l.principal, rate, n)
	if err != nil {
		return Payoff{}, err
	}

	extra := money.NewFromDecimal(l.extra, l.principal.GetCurrency())

	return amortize(l.principal, rate, payment, extra, n)
}

// MustPayoff is like Payoff but panics on error.
func (l Config) MustPayoff() Payoff {
	p, err := l.Payoff()
	if err != nil {
		panic(err)
	}

	return p
}

// Savings measures what the configured extra payment buys: how many payments
// it removes from the term and how much interest it never accrues. Both
// payoffs are amortized from the same terms, so the two figures are directly
// comparable. With no extra payment configured the savings are zero.
//
// It returns the same errors as Payoff.
func (l Config) Savings() (Savings, error) {
	accelerated, err := l.Payoff()
	if err != nil {
		return Savings{}, err
	}

	base := l
	base.extra = decimal.Zero

	scheduled, err := base.Payoff()
	if err != nil {
		return Savings{}, err
	}

	return Savings{
		Scheduled:     scheduled,
		Accelerated:   accelerated,
		PeriodsSaved:  scheduled.Periods - accelerated.Periods,
		InterestSaved: scheduled.TotalInterest.Sub(accelerated.TotalInterest),
	}, nil
}

// MustSavings is like Savings but panics on error.
func (l Config) MustSavings() Savings {
	s, err := l.Savings()
	if err != nil {
		panic(err)
	}

	return s
}

// amortize runs the balance down from pv, charging rate on the outstanding
// balance and applying payment + extra each period, for at most maxPeriods
// periods. Whenever the amount due would overshoot the remaining balance —
// which is what happens on the final period, early or not — only the payoff
// amount is charged.
func amortize(pv money.Money, rate decimal.Decimal, payment, extra money.Money, maxPeriods int) (Payoff, error) {
	if extra.IsNegative() {
		return Payoff{}, ErrNegativeExtra
	}

	currency := pv.GetCurrency()
	zero := money.NewFromDecimal(decimal.Zero, currency)
	due := payment.Add(extra)

	balance, sumInterest, totalPaid, final := pv, zero, zero, zero

	rows := make([]annuities.Schedule, 0, maxPeriods+1)
	rows = append(rows, annuities.Schedule{
		Period:      decimal.Zero,
		Balance:     pv,
		Payment:     zero,
		Interest:    zero,
		SumInterest: zero,
		Principal:   zero,
	})

	periods := 0

	for p := 1; p <= maxPeriods; p++ {
		interest := balance.MulDecimal(rate)
		payoffAmount := balance.Add(interest)

		// The scheduled amount is capped by what is still owed, and the last
		// contractual period always settles the balance in full so no rounding
		// residue survives the term.
		paid := due
		if paid.GreaterThanOrEqual(payoffAmount) || p == maxPeriods {
			paid = payoffAmount
		}

		principal := paid.Sub(interest)
		balance = balance.Sub(principal)
		sumInterest = sumInterest.Add(interest)
		totalPaid = totalPaid.Add(paid)

		index, err := decimal.NewFromInt64(int64(p), 0)
		if err != nil {
			return Payoff{}, err
		}

		rows = append(rows, annuities.Schedule{
			Period:      index,
			Balance:     balance,
			Payment:     paid,
			Interest:    interest,
			SumInterest: sumInterest,
			Principal:   principal,
		})

		periods, final = p, paid

		if !balance.IsPositive() {
			break
		}
	}

	return Payoff{
		Periods:       periods,
		Payment:       payment,
		ExtraPayment:  extra,
		FinalPayment:  final,
		TotalPaid:     totalPaid,
		TotalInterest: sumInterest,
		Schedule:      rows,
	}, nil
}
