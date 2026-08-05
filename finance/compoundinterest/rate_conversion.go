package compoundinterest

import "github.com/yeferson59/gofinance/v2/decimal"

// periodicEffective returns the effective periodic rate implied by rt — the
// rate actually applied at the end of each compounding period — together with
// the number of compounding periods per year.
//
// Every conversion in this file is derived from this pair, so that a rate
// declared in any of the five TypeRate forms converts to any of the others.
// The two anticipated (discount) forms are folded in through i = d/(1-d), the
// inverse of d = i/(1+i).
//
// It returns ErrInvalidTypeRate for an unknown rate type and
// ErrInvalidAnticipatedRate for an anticipated rate of 100% or more.
func (rt RateInterest) periodicEffective() (decimal.Decimal, decimal.Decimal, error) {
	periodsPerYear, err := rt.compoundingFrequency.PeriodsPerYear()
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, err
	}

	switch rt.typeRate {
	case RateEffectyPeriodic:
		return rt.value, periodsPerYear, nil

	case RateEffectyNominal:
		periodic, err := rt.value.Div(periodsPerYear)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		return periodic, periodsPerYear, nil

	case RateEffectyAnnually:
		exponent, err := decimal.One.Div(periodsPerYear)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		growth, err := rt.value.Add(decimal.One).Pow(exponent)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		return growth.Sub(decimal.One), periodsPerYear, nil

	case RateAnticipateEffectyPeriodic:
		periodic, err := anticipatedToEffective(rt.value)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		return periodic, periodsPerYear, nil

	case RateAnticipateEffectyNominal:
		discount, err := rt.value.Div(periodsPerYear)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		periodic, err := anticipatedToEffective(discount)
		if err != nil {
			return decimal.Decimal{}, decimal.Decimal{}, err
		}

		return periodic, periodsPerYear, nil

	default:
		return decimal.Decimal{}, decimal.Decimal{}, ErrInvalidTypeRate
	}
}

// anticipatedToEffective converts an anticipated (discount) periodic rate d
// into the equivalent effective periodic rate:
//
//	i = d / (1 - d)
//
// It returns ErrInvalidAnticipatedRate when d ≥ 1, which leaves nothing to
// discount from.
func anticipatedToEffective(discount decimal.Decimal) (decimal.Decimal, error) {
	remainder := decimal.One.Sub(discount)
	if !remainder.IsPos() {
		return decimal.Decimal{}, ErrInvalidAnticipatedRate
	}

	return discount.Div(remainder)
}

// effectiveToAnticipated converts an effective periodic rate i into the
// equivalent anticipated (discount) periodic rate:
//
//	d = i / (1 + i)
func effectiveToAnticipated(periodic decimal.Decimal) (decimal.Decimal, error) {
	return periodic.Div(decimal.One.Add(periodic))
}

// RatePeriodic converts the current interest rate to a periodic (per-period) rate:
// the rate actually applied to each compounding period. It accepts any of the five
// rate types, including the anticipated (discount) ones.
//
// Returns:
//   - The periodic interest rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
//
// Example:
//
//	rate, _ := NewRateInterest(0.12, Monthly, RateEffectyNominal)
//	periodicRate, _ := rate.RatePeriodic()  // Converts 12% nominal to ~1% monthly
func (rt RateInterest) RatePeriodic() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyPeriodic {
		return rt.value, nil
	}

	periodic, _, err := rt.periodicEffective()

	return periodic, err
}

// RateNominal converts the current interest rate to a nominal rate.
// A nominal rate is the stated rate before considering compounding frequency:
// the periodic rate times the number of periods per year. It accepts any of the
// five rate types, including the anticipated (discount) ones.
//
// Returns:
//   - The nominal interest rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
func (rt RateInterest) RateNominal() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyNominal {
		return rt.value, nil
	}

	periodic, periodsPerYear, err := rt.periodicEffective()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return periodic.Mul(periodsPerYear), nil
}

// RateEffectyAnnually converts the current interest rate to an effective annual rate.
// The effective annual rate accounts for compounding and represents the true yearly
// return. It accepts any of the five rate types, including the anticipated
// (discount) ones: the effective annual rate is the same quantity whichever form
// the rate was quoted in.
//
// Returns:
//   - The effective annual interest rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
func (rt RateInterest) RateEffectyAnnually() (decimal.Decimal, error) {
	if rt.typeRate == RateEffectyAnnually {
		return rt.value, nil
	}

	periodic, periodsPerYear, err := rt.periodicEffective()
	if err != nil {
		return decimal.Decimal{}, err
	}

	growth, err := decimal.One.Add(periodic).Pow(periodsPerYear)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return growth.Sub(decimal.One), nil
}

// RatePeriodicToPeriodic converts a periodic rate from one compounding frequency to another.
// For example, converts a monthly rate to a quarterly rate:
//
//	i' = (1+i)^(m/m') - 1
//
// Parameters:
//   - newCompoundingFrequency: The target compounding frequency
//
// Returns:
//   - The converted periodic rate for the new frequency as a Decimal
//   - An error if the conversion fails
//
// Example:
//
//	rate, _ := NewRateInterest(0.01, Monthly, RateEffectyPeriodic)
//	quarterlyRate, _ := rate.RatePeriodicToPeriodic(Quarterly)
func (rt RateInterest) RatePeriodicToPeriodic(newCompoundingFrequency CompoundingFrequency) (decimal.Decimal, error) {
	currentPeriodicRate, currentPeriodsPerYear, err := rt.periodicEffective()
	if err != nil {
		return decimal.Decimal{}, err
	}

	newPeriodsPerYear, err := newCompoundingFrequency.PeriodsPerYear()
	if err != nil {
		return decimal.Decimal{}, err
	}

	exponent, err := currentPeriodsPerYear.Div(newPeriodsPerYear)
	if err != nil {
		return decimal.Decimal{}, err
	}

	growth, err := currentPeriodicRate.Add(decimal.One).Pow(exponent)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return growth.Sub(decimal.One), nil
}

// RateNominalToNominal converts a nominal rate from one compounding frequency to another.
// For example, converts a monthly nominal rate to a quarterly nominal rate.
//
// Parameters:
//   - newCompoundingFrequency: The target compounding frequency
//
// Returns:
//   - The converted nominal rate for the new frequency as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) RateNominalToNominal(newCompoundingFrequency CompoundingFrequency) (decimal.Decimal, error) {
	newPeriodicRate, err := rt.RatePeriodicToPeriodic(newCompoundingFrequency)
	if err != nil {
		return decimal.Decimal{}, err
	}

	newRateInterest, err := NewRateInterest(newPeriodicRate, newCompoundingFrequency, RateEffectyPeriodic)
	if err != nil {
		return decimal.Decimal{}, err
	}

	return newRateInterest.RateNominal()
}

// RateAnticipateEffectyAnnually converts the current rate to an effective annual rate.
//
// An effective annual rate is a single quantity regardless of whether the rate was
// quoted in ordinary or anticipated (discount) form, so this is equivalent to
// RateEffectyAnnually and is kept for source compatibility.
//
// Returns:
//   - The effective annual rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
func (rt RateInterest) RateAnticipateEffectyAnnually() (decimal.Decimal, error) {
	return rt.RateEffectyAnnually()
}

// RateAnticipateNominal converts the current rate to an anticipated (discount) nominal
// rate: the anticipated periodic rate times the number of periods per year. It accepts
// any of the five rate types.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
func (rt RateInterest) RateAnticipateNominal() (decimal.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyNominal {
		return rt.value, nil
	}

	discount, err := rt.RateAnticipatePeriodic()
	if err != nil {
		return decimal.Decimal{}, err
	}

	periodsPerYear, err := rt.compoundingFrequency.PeriodsPerYear()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return discount.Mul(periodsPerYear), nil
}

// RateAnticipatePeriodic converts the current rate to an anticipated (discount)
// periodic rate, d = i/(1+i): the rate charged at the beginning of each period that
// is equivalent to the ordinary rate charged at the end. It accepts any of the five
// rate types.
//
// Returns:
//   - The anticipated periodic rate as a Decimal
//   - ErrInvalidTypeRate for an unknown rate type, ErrInvalidAnticipatedRate for
//     an anticipated rate of 100% or more, or an error if the conversion fails
func (rt RateInterest) RateAnticipatePeriodic() (decimal.Decimal, error) {
	if rt.typeRate == RateAnticipateEffectyPeriodic {
		return rt.value, nil
	}

	periodic, _, err := rt.periodicEffective()
	if err != nil {
		return decimal.Decimal{}, err
	}

	return effectiveToAnticipated(periodic)
}

// ToAnticipateNominal converts an effective or periodic rate to an anticipated nominal
// rate. This is useful for discount instruments like Treasury bills.
//
// Equivalent to RateAnticipateNominal, which now accepts every rate type.
//
// Returns:
//   - The anticipated nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipateNominal() (decimal.Decimal, error) {
	return rt.RateAnticipateNominal()
}

// ToAnticipatePeriodic converts an effective or periodic rate to an anticipated periodic rate.
//
// Equivalent to RateAnticipatePeriodic, which now accepts every rate type.
//
// Returns:
//   - The anticipated periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToAnticipatePeriodic() (decimal.Decimal, error) {
	return rt.RateAnticipatePeriodic()
}

// ToNominal converts an anticipated rate to a nominal rate.
//
// Equivalent to RateNominal, which now accepts every rate type.
//
// Returns:
//   - The nominal rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToNominal() (decimal.Decimal, error) {
	return rt.RateNominal()
}

// ToPeriodic converts an anticipated rate to a periodic rate.
//
// Equivalent to RatePeriodic, which now accepts every rate type.
//
// Returns:
//   - The periodic rate as a Decimal
//   - An error if the conversion fails
func (rt RateInterest) ToPeriodic() (decimal.Decimal, error) {
	return rt.RatePeriodic()
}
