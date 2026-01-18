# Annuities Package

Package for performing annuity calculations in Go.

## Description

The `annuities` package provides a complete solution for annuity financial calculations. It allows you to calculate payments, present values, future values, and periods for annuities with compound interest. An annuity is a series of equal payments made at regular intervals.

## Features

- ✅ **Payment Calculation from Present Value** - Calculate the periodic payment needed to amortize a present value
- ✅ **Payment Calculation from Future Value** - Calculate the periodic payment needed to accumulate a future value
- ✅ **Present Value Calculation** - Calculate the present value of an annuity
- ✅ **Future Value Calculation** - Calculate the future value of an annuity
- ✅ **Number of Periods Calculation** - Calculate how many periods are needed
- ✅ Support for multiple **compounding frequencies**:
  - Daily (365 periods/year)
  - Monthly (12 periods/year)
  - Bimonthly (6 periods/year)
  - Quarterly (4 or 3 periods/year)
  - Semi-annually (2 periods/year)
  - Annually (1 period/year)

- ✅ Support for multiple **types of rates**:
  - **Ordinary Rates** (charged at the end of the period):
    - Periodic Rate (i)
    - Nominal Annual Rate (j)
    - Effective Annual Rate (i_a)
  - **Anticipated Rates** (charged at the beginning of the period - discounts)
    - Anticipated Periodic Rate
    - Anticipated Nominal Rate
    - Anticipated Effective Annual Rate

- ✅ **Test Coverage: 75%** with comprehensive tests

## Installation

```bash
go get github.com/yeferson59/gofinance
```

## Basic Usage

### 1. Calculate Periodic Payment from Present Value

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/annuities"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
  // Create a monthly interest rate of 1% (periodic)
  rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)

  // Create a period of 12 months
  period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

  // Create an annuity object with monthly payment of $100
  // Parameters: (payment=100, present=0, future=0, period, rateInterest)
  ann, _ := annuities.New(100, 0, 0, period, rate)

  // Calculate present value
  present, _ := ann.Present()
  fmt.Printf("Present value: $%.2f\n", present)
  // Output: Present value: $1,125.51
}
```

### 2. Calculate Periodic Payment to Amortize a Loan

```go
// Calculate monthly payment for a $10,000 loan at 1% monthly for 12 months
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

ann, _ := annuities.New(0, 10000, 0, period, rate)
payment, _ := ann.PaymentFromPresentValue()
fmt.Printf("Monthly payment: $%.2f\n", payment)
// Output: Monthly payment: $888.49
```

### 3. Calculate Future Value of an Annuity

```go
// Calculate future value of $100 monthly payments at 1% monthly for 12 months
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

ann, _ := annuities.New(100, 0, 0, period, rate)
future, _ := ann.Future()
fmt.Printf("Future value: $%.2f\n", future)
// Output: Future value: $1,268.25
```

### 4. Calculate Payment to Accumulate a Future Value

```go
// Calculate monthly payment needed to accumulate $10,000 in 12 months at 1% monthly
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

ann, _ := annuities.New(0, 0, 10000, period, rate)
payment, _ := ann.PaymentFromFutureValue()
fmt.Printf("Monthly payment: $%.2f\n", payment)
// Output: Monthly payment: $788.49
```

### 5. Calculate Number of Periods

```go
// Calculate how many months are needed for $100 payments to reach present value
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(1, compositeinterest.Monthly)

ann, _ := annuities.New(100, 5000, 0, period, rate)
periods, _ := ann.PeriodsWithPresent()
fmt.Printf("Number of periods: %.2f months\n", periods)
```

## Mathematical Formulas

### Payment from Present Value

```
PMT = PV × [i(1 + i)^n] / [(1 + i)^n - 1]

where:
  PMT = Periodic payment
  PV = Present value
  i = Periodic rate
  n = Number of periods
```

### Present Value of an Annuity

```
PV = PMT × [(1 + i)^n - 1] / [i(1 + i)^n]
```

### Future Value of an Annuity

```
FV = PMT × [((1 + i)^n - 1) / i]
```

### Payment from Future Value

```
PMT = FV × [i / ((1 + i)^n - 1)]
```

### Number of Periods (from Present Value)

```
n = ln(PMT / (PMT - PV×i)) / ln(1 + i)
```

### Number of Periods (from Future Value)

```
n = ln((FV×i + PMT) / PMT) / ln(1 + i)
```

## Main Types

### Annuity

Represents an annuity with periodic payments:

```go
ann, err := annuities.New(payment, present, future, period, rateInterest)
```

## Main Methods

| Method                    | Description                              | Use                                      |
| ------------------------- | ---------------------------------------- | ---------------------------------------- |
| `PaymentFromPresentValue()` | Calculates periodic payment              | When you know: PV, i, n                  |
| `PaymentFromFutureValue()` | Calculates payment to reach FV           | When you know: FV, i, n                  |
| `Present()`               | Calculates present value                 | When you know: PMT, i, n                 |
| `Future()`                | Calculates future value                  | When you know: PMT, i, n                 |
| `PeriodsWithPresent()`    | Calculates periods with present value    | When you know: PMT, PV, i                |
| `PeriodsWithFuture()`     | Calculates periods with future value     | When you know: PMT, FV, i                |

## Advanced Examples

### Mortgage Calculation

```go
// Calculate monthly payment for a $300,000 mortgage at 4.5% annual for 30 years
rate, _ := compositeinterest.NewRateInterest(
    0.045,
    compositeinterest.Annually,
    compositeinterest.RateEffectyAnnually,
)

// Convert to monthly equivalent
monthlyRate, _ := rate.RatePeriodicToPeriodic(compositeinterest.Monthly)
rate, _ = compositeinterest.NewRateInterest(
    monthlyRate,
    compositeinterest.Monthly,
    compositeinterest.RateEffectyPeriodic,
)

period, _ := compositeinterest.NewPeriod(360, compositeinterest.Monthly) // 30 years × 12 months

ann, _ := annuities.New(0, 300000, 0, period, rate)
payment, _ := ann.PaymentFromPresentValue()
fmt.Printf("Monthly payment: $%.2f\n", payment)
```

### Savings Plan

```go
// How much will be accumulated if saving $500/month at 6% annual for 5 years?
rate, _ := compositeinterest.NewRateInterest(0.06, compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
monthlyRate, _ := rate.RatePeriodic()
rate, _ = compositeinterest.NewRateInterest(monthlyRate, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)

period, _ := compositeinterest.NewPeriod(60, compositeinterest.Monthly) // 5 years × 12 months

ann, _ := annuities.New(500, 0, 0, period, rate)
future, _ := ann.Future()
fmt.Printf("Accumulated amount: $%.2f\n", future)
```

### Lease Payment Calculation

```go
// Calculate lease payment for equipment worth $50,000 over 5 years at 8% annual
rate, _ := compositeinterest.NewRateInterest(0.08, compositeinterest.Annually, compositeinterest.RateEffectyAnnually)
period, _ := compositeinterest.NewPeriod(5, compositeinterest.Annually)

ann, _ := annuities.New(0, 50000, 0, period, rate)
payment, _ := ann.PaymentFromPresentValue()
fmt.Printf("Annual lease payment: $%.2f\n", payment)
```

## Testing

The package includes comprehensive tests with 75% coverage:

```bash
# Run all tests
go test ./finance/annuities -v

# Run tests with coverage
go test ./finance/annuities -cover

# Run a specific test
go test -run TestPresentValue ./finance/annuities
```

### Test Files

- `present_test.go` - Tests for present value calculations

## Use Cases

1. **Loan Amortization**: Calculate monthly payments for mortgages or personal loans
2. **Savings Plans**: Determine how much you'll save with regular deposits
3. **Investment Returns**: Calculate future value of regular investments
4. **Lease Calculations**: Determine periodic lease payments
5. **Retirement Planning**: Calculate payment needed to accumulate retirement funds
6. **Insurance**: Calculate annuity benefits

## Important Notes

- All rates must be entered as decimals (0.05 for 5%)
- Periods must be positive numbers
- The compounding frequency must match the unit of periods
- For accurate results, ensure consistency between payment frequency and rate frequency
- When payment is known (not 0), it's used to calculate PV or FV

## Contributing

Contributions are welcome. Please ensure:

1. Write tests for new features
2. Maintain test coverage above 75%
3. Document the code properly
4. Follow the existing code style

## License

MIT License

## Author

Yeferson Toloza
