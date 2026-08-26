package loans

import (
	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/investment"
	"github.com/yeferson59/gofinance/v2/money"
)

// Comparison weighs a refinance offer against the loan currently being paid.
// Every amount is expressed in the loans' shared currency.
type Comparison struct {
	// CurrentPayment and OfferPayment are the scheduled level payments of each
	// loan.
	CurrentPayment money.Money
	OfferPayment   money.Money

	// PaymentSavings is what the offer saves each period, positive when the
	// offer is the cheaper one.
	PaymentSavings money.Money

	// ClosingCosts is the offer's up-front fees, treated as paid at closing.
	ClosingCosts money.Money

	// BreakEvenPeriods is the number of periods of accumulated savings needed
	// to recover the closing costs — the point past which refinancing is ahead
	// in nominal terms. It is zero when the offer costs nothing to take.
	BreakEvenPeriods int

	// CurrentTotalInterest and OfferTotalInterest are the interest each loan
	// accrues over its whole life.
	CurrentTotalInterest money.Money
	OfferTotalInterest   money.Money

	// InterestSaved is how much less interest the offer accrues, before
	// closing costs.
	InterestSaved money.Money

	// NetPresentValue discounts the whole stream of payment differences at the
	// offer's periodic rate and subtracts the closing costs paid today. It is
	// the figure to judge the offer by when the two terms differ: positive
	// means the offer is worth taking.
	NetPresentValue money.Money
}

// Compare evaluates a refinance offer against the loan currently outstanding.
// Model the current loan with its *remaining* balance and *remaining* number of
// payments, and the offer with what the new lender would advance and its own
// term and fees.
//
// The comparison runs off both amortization schedules, so a different term, a
// configured extra payment, or a smaller final payment on either side is
// accounted for period by period. The offer's fees are treated as closing costs
// paid today rather than financed.
//
// It returns money.ErrCurrencyMismatch when the loans are in different
// currencies, ErrFrequencyMismatch when they pay on different schedules,
// ErrNoBreakEven when the savings never recover the closing costs, and any term
// error of either loan.
func Compare(current, offer Config) (Comparison, error) {
	if current.principal.GetCurrency() != offer.principal.GetCurrency() {
		return Comparison{}, money.ErrCurrencyMismatch
	}

	if current.freq != offer.freq {
		return Comparison{}, ErrFrequencyMismatch
	}

	currentPayoff, err := current.Payoff()
	if err != nil {
		return Comparison{}, err
	}

	offerPayoff, err := offer.Payoff()
	if err != nil {
		return Comparison{}, err
	}

	// Validates the fees and reports them as the cost of taking the offer.
	if _, err = offer.NetProceeds(); err != nil {
		return Comparison{}, err
	}

	currency := offer.principal.GetCurrency()
	closingCosts := money.NewFromDecimal(offer.fees, currency)

	savings := paymentDifferences(currentPayoff, offerPayoff, currency)

	breakEven, err := breakEvenPeriod(savings, closingCosts)
	if err != nil {
		return Comparison{}, err
	}

	rate, _, err := offer.terms()
	if err != nil {
		return Comparison{}, err
	}

	npv, err := investment.NPV(rate, append([]money.Money{closingCosts.Neg()}, savings...))
	if err != nil {
		return Comparison{}, err
	}

	return Comparison{
		CurrentPayment:       currentPayoff.Payment,
		OfferPayment:         offerPayoff.Payment,
		PaymentSavings:       currentPayoff.Payment.Sub(offerPayoff.Payment),
		ClosingCosts:         closingCosts,
		BreakEvenPeriods:     breakEven,
		CurrentTotalInterest: currentPayoff.TotalInterest,
		OfferTotalInterest:   offerPayoff.TotalInterest,
		InterestSaved:        currentPayoff.TotalInterest.Sub(offerPayoff.TotalInterest),
		NetPresentValue:      npv,
	}, nil
}

// MustCompare is like Compare but panics on error.
func MustCompare(current, offer Config) Comparison {
	c, err := Compare(current, offer)
	if err != nil {
		panic(err)
	}

	return c
}

// paymentDifferences returns, for every period until the longer of the two
// loans is retired, what the current loan costs minus what the offer costs. A
// loan that is already paid off contributes nothing from then on, which is
// exactly how a shorter term shows up as a cost later in the horizon.
func paymentDifferences(current, offer Payoff, currency money.Currency) []money.Money {
	horizon := max(current.Periods, offer.Periods)
	zero := money.NewFromDecimal(decimal.Zero, currency)

	differences := make([]money.Money, 0, horizon)

	for p := 1; p <= horizon; p++ {
		differences = append(differences, paymentAt(current, p, zero).Sub(paymentAt(offer, p, zero)))
	}

	return differences
}

// paymentAt returns the payment made in period p, or zero once the loan is
// retired. Schedule index p holds period p because index 0 is the opening row.
func paymentAt(p Payoff, period int, zero money.Money) money.Money {
	if period > p.Periods {
		return zero
	}

	return p.Schedule[period].Payment
}

// breakEvenPeriod returns the first period at which the accumulated savings
// cover the closing costs, or ErrNoBreakEven if they never do. Costs of zero
// break even immediately, before any payment is made.
func breakEvenPeriod(savings []money.Money, closingCosts money.Money) (int, error) {
	accumulated := money.NewFromDecimal(decimal.Zero, closingCosts.GetCurrency())

	if accumulated.GreaterThanOrEqual(closingCosts) {
		return 0, nil
	}

	for p, saving := range savings {
		accumulated = accumulated.Add(saving)

		if accumulated.GreaterThanOrEqual(closingCosts) {
			return p + 1, nil
		}
	}

	return 0, ErrNoBreakEven
}
