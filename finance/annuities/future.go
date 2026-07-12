package annuities

import (
	"errors"

	"github.com/yeferson59/gofinance/finance/compositeinterest"
	"github.com/yeferson59/gofinance/money"
)

func (a Annuity) Future() (money.Money, error) {
	if future, err := a.compositeInterest.Future(); err == nil && !future.IsZero() {
		return future, nil
	}

	return a.contributionsFuture()
}

// contributionsFuture returns the future value of the periodic contributions
// (the annuity's payment amount) alone, assuming each one is made at the end
// of its period (ordinary annuity). Unlike Future, it never substitutes a
// pre-set or principal-derived future value, so it can be combined with
// principalFuture in FutureWithContributions.
func (a Annuity) contributionsFuture() (money.Money, error) {
	periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	growthPower, err := money.One.Add(rateInterest).Pow(periods)
	if err != nil {
		return money.Money{}, err
	}

	result, err := growthPower.Sub(money.One).Div(rateInterest)
	if err != nil {
		return money.Money{}, err
	}

	return a.value.ToDecimal().Mul(result).ToMoney(a.value.Currency()), nil
}

func (a Annuity) AnticipateFuture() (money.Money, error) {
	if future, err := a.compositeInterest.Future(); err == nil && !future.IsZero() {
		return future, nil
	}

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

	_, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
	if err != nil {
		return money.Money{}, err
	}

	return ordinary.ToDecimal().Mul(money.One.Add(rateInterest)).ToMoney(ordinary.Currency()), nil
}

// principalFuture returns the future value of the initial principal (Present),
// compounded over the annuity's periods and rate: PV × (1+i)^n. If no
// principal was configured, it returns zero instead of an error, so it can be
// added freely to contributionsFuture/contributionsAnticipateFuture.
func (a Annuity) principalFuture() (money.Money, error) {
	principal, err := a.compositeInterest.Future()
	if errors.Is(err, compositeinterest.ErrInvalidOperation) {
		return money.MustMoneyFromFloat64(0, a.value.Currency()), nil
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

	principal, err := a.principalFuture()
	if err != nil {
		return money.Money{}, err
	}

	return contributions.Add(principal), nil
}

// AnticipateFutureWithContributions is like FutureWithContributions, but
// assumes each contribution is made at the beginning of its period (annuity
// due), so it also earns interest during its own first period.
func (a Annuity) AnticipateFutureWithContributions() (money.Money, error) {
	contributions, err := a.contributionsAnticipateFuture()
	if err != nil {
		return money.Money{}, err
	}

	principal, err := a.principalFuture()
	if err != nil {
		return money.Money{}, err
	}

	return contributions.Add(principal), nil
}
