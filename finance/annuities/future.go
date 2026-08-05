package annuities

import (
	"errors"

	"github.com/yeferson59/gofinance/v2/decimal"
	"github.com/yeferson59/gofinance/v2/finance/compoundinterest"
	"github.com/yeferson59/gofinance/v2/money"
)

// Future returns the future value of the ordinary annuity: the periodic
// payments (Value) alone, each made at the end of its period, grown to the
// end of the term:
//
//	FV = PMT × ((1+i)ⁿ − 1) / i
//
// It never substitutes the principal's growth or a pre-set future value —
// use PrincipalFuture for the compounded principal alone, or
// FutureWithContributions for principal plus payments.
func (a Annuity) Future() (money.Money, error) {
	return a.contributionsFuture()
}

// contributionsFuture computes the future value of the periodic contributions
// (the annuity's payment amount) alone, assuming each one is made at the end
// of its period (ordinary annuity). It backs Future and is combined with
// PrincipalFuture in FutureWithContributions.
func (a Annuity) contributionsFuture() (money.Money, error) {
	periods, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	// With no interest the contributions never grow, so the future value is
	// just their sum. The general formula divides by the rate, so the limit
	// is returned directly.
	if rateInterest.IsZero() {
		return money.FromDecimal(a.value.ToDecimal().Mul(periods), a.currency), nil
	}

	growthPower, err := decimal.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	result, err := growthPower.Sub(decimal.One).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(a.value.ToDecimal().Mul(result), a.currency), nil
}

// AnticipateFuture is like Future, but assumes each payment is made at the
// beginning of its period (annuity due), so it also earns interest during its
// own first period: FV_due = FV_ordinary × (1+i). Like Future, it never
// substitutes the principal's growth or a pre-set future value.
func (a Annuity) AnticipateFuture() (money.Money, error) {
	return a.contributionsAnticipateFuture()
}

// contributionsAnticipateFuture is like contributionsFuture, but assumes each
// contribution is made at the beginning of its period (annuity due), so it
// also earns interest during its own first period: FV_due = FV_ordinary × (1+i).
func (a Annuity) contributionsAnticipateFuture() (money.Money, error) {
	ordinary, err := a.contributionsFuture()
	if err != nil {
		return money.Money{}, err
	}

	_, rateInterest, err := a.compoundInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	return money.FromDecimal(ordinary.ToDecimal().Mul(decimal.One.Add(rateInterest)), a.currency), nil
}

// PrincipalFuture returns the future value of the initial principal (Present)
// alone, compounded over the annuity's periods and rate — or the pre-set
// future value if one was configured:
//
//	FV = PV × (1+i)ⁿ
//
// If no principal was configured, it returns zero instead of an error, so it
// can be added freely to the payments' future value (which is exactly what
// FutureWithContributions does).
func (a Annuity) PrincipalFuture() (money.Money, error) {
	principal, err := a.compoundInterest.Future()
	if errors.Is(err, compoundinterest.ErrInvalidOperation) {
		return money.MustMoneyFromFloat64(0, a.currency), nil
	}
	if err != nil {
		return money.Money{}, err
	}

	return principal, nil
}

// FutureWithContributions returns the future value of a recurring investment:
// the initial principal (Present) compounded over the term, plus the future
// value of equal periodic contributions (Value) made at the end of each
// period (ordinary annuity).
//
// Use this to model an investment plan such as "I have $1,000 today and add
// $100 every month" — Present is the $1,000, Value is the $100 contribution.
// If no principal was configured, only the contributions' growth is
// returned.
//
// Example:
//
//	future, err := NewAnnuity().
//	    Present(1000, money.USD).
//	    Value(100, money.USD).
//	    AnnualRate(0.06).
//	    Periods(12).
//	    Monthly().
//	    MustBuild().
//	    FutureWithContributions()
func (a Annuity) FutureWithContributions() (money.Money, error) {
	contributions, err := a.contributionsFuture()
	if err != nil {
		return money.Money{}, err
	}

	principal, err := a.PrincipalFuture()
	if err != nil {
		return money.Money{}, err
	}

	// TryAdd rather than Add: the two halves are built from the annuity's
	// resolved currency, so they agree, but a function that returns an error
	// must report a mismatch rather than panic if that ever stops holding.
	return contributions.TryAdd(principal)
}

// AnticipateFutureWithContributions is like FutureWithContributions, but
// assumes each contribution is made at the beginning of its period (annuity
// due), so it also earns interest during its own first period.
func (a Annuity) AnticipateFutureWithContributions() (money.Money, error) {
	contributions, err := a.contributionsAnticipateFuture()
	if err != nil {
		return money.Money{}, err
	}

	principal, err := a.PrincipalFuture()
	if err != nil {
		return money.Money{}, err
	}

	// TryAdd rather than Add: the two halves are built from the annuity's
	// resolved currency, so they agree, but a function that returns an error
	// must report a mismatch rather than panic if that ever stops holding.
	return contributions.TryAdd(principal)
}
