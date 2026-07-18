# Composite Interest Package

Package for performing compound interest calculations in Go.

## Description

The `compoundinterest` package provides a complete solution for compound interest financial calculations. It allows you to calculate present values, future values, interest rates, and periods, as well as perform conversions between different types of rates.

## Features

- ✅ **Future Value** (FV) Calculation
- ✅ **Present Value** (PV) Calculation
- ✅ **Interest Rate** Calculation
- ✅ **Number of Periods** Calculation
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

- ✅ **Rate conversion** between all types
- ✅ Test Coverage: **86.5%**
- ✅ **105 comprehensive tests** using Testify

## Installation

```bash
go get github.com/yeferson59/gofinance
```

## Basic Usage

### 1. Calculate Future Value

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/compoundinterest"
)

func main() {
  // Create a monthly interest rate of 1% (periodic)
  rate, _ := compoundinterest.NewRateInterest(0.01, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)

  // Create a period of 12 months
  period, _ := compoundinterest.NewPeriod(12, compoundinterest.Monthly)

  // Create a compound interest object
  // Parameters: (present=1000, future=0, rateInterest, period)
  ci, _ := compoundinterest.New(1000, 0, rate, period)

  // Calculate future value
  future, _ := ci.Future()
  fmt.Printf("Future value: $%.2f\n", future)
  // Output: Future value: $1126.83
}
```

### 2. Calculate Present Value

```go
// Determine how much money to invest today to get $1000 in 12 months
rate, _ := compoundinterest.NewRateInterest(0.01, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
period, _ := compoundinterest.NewPeriod(12, compoundinterest.Monthly)

ci, _ := compoundinterest.New(0, 1000, rate, period)
present, _ := ci.Present()
fmt.Printf("Present value: $%.2f\n", present)
// Output: Present value: $887.45
```

### 3. Calculate Interest Rate

```go
// Find the interest rate between two values
rate, _ := compoundinterest.NewRateInterest(0.01, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
period, _ := compoundinterest.NewPeriod(12, compoundinterest.Monthly)

ci, _ := compoundinterest.New(1000, 1126.83, rate, period)
interest, _ := ci.Interest()
fmt.Printf("Periodic rate: %.4f (%.2f%%)\n", interest, interest*100)
// Output: Periodic rate: 0.0100 (1.00%)
```

### 4. Calculate Periods

```go
// Determine how many periods are needed for money to grow from $1000 to $2000
rate, _ := compoundinterest.NewRateInterest(0.05, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
period, _ := compoundinterest.NewPeriod(1, compoundinterest.Monthly) // initial value

ci, _ := compoundinterest.New(1000, 2000, rate, period)
periods, _ := ci.Periods()
fmt.Printf("Number of periods: %.2f months\n", periods)
```

### 5. Convert Rates

```go
// Convert an annual nominal rate of 12% to monthly periodic rate
rate, _ := compoundinterest.NewRateInterest(0.12, compoundinterest.Monthly, compoundinterest.RateEffectyNominal)

periodic, _ := rate.RatePeriodic()
fmt.Printf("Periodic rate: %.4f\n", periodic)
// Output: Periodic rate: 0.0100

annual, _ := rate.RateEffectyAnnually()
fmt.Printf("Effective annual rate: %.4f\n", annual)
// Output: Effective annual rate: 0.1268
```

## Mathematical Formulas

### Future Value

```
FV = PV × (1 + i)^n

where:
  FV = Future Value
  PV = Present Value
  i  = Periodic rate
  n  = Number of periods
```

### Present Value

```
PV = FV / (1 + i)^n
```

### Interest Rate

```
i = (FV / PV)^(1/n) - 1
```

### Number of Periods

```
n = ln(FV / PV) / ln(1 + i)
```

### Rate Conversions

**Nominal to Periodic:**

```
i_periodic = j / m

where:
  j = nominal rate
  m = number of periods per year
```

**Periodic to Effective Annual:**

```
i_a = (1 + i_periodic)^m - 1
```

**Nominal to Effective Annual:**

```
i_a = (1 + j/m)^m - 1
```

## Main Types

### CompoundingFrequency

Defines the frequency at which interest compounds:

- `Daily`: Daily (365)
- `Monthly`: Monthly (12)
- `Bimonthly`: Bimonthly (6)
- `Quarterly`: Quarterly (4)
- `FourMonthly`: Every four months (3)
- `SemiAnnually`: Semi-annually (2)
- `Annually`: Annually (1)

### TypeRate

Defines the type of rate:

- `RateEffectyPeriodic`: Ordinary periodic rate
- `RateEffectyNominal`: Ordinary nominal rate
- `RateEffectyAnnually`: Ordinary effective annual rate
- `RateAnticipateEffectyPeriodic`: Anticipated periodic rate
- `RateAnticipateEffectyNominal`: Anticipated nominal rate
- `RateAnticipateEffectyAnnually`: Anticipated effective annual rate

### Period

Represents a specific period:

```go
period, err := NewPeriod(12, Monthly)
```

### RateInterest

Represents a rate with its frequency and type:

```go
rate, err := NewRateInterest(0.12, Monthly, RateEffectyNominal)
```

### CompoundInterest

Contains all parameters for a calculation:

```go
ci, err := New(1000, 1636.36, rate, period)
```

## Main Methods

| Method       | Description                  | Use                      |
| ------------ | ---------------------------- | ------------------------ |
| `Future()`   | Calculates future value      | When you know: PV, i, n  |
| `Present()`  | Calculates present value     | When you know: FV, i, n  |
| `Interest()` | Calculates interest rate     | When you know: PV, FV, n |
| `Periods()`  | Calculates number of periods | When you know: PV, FV, i |

## Rate Conversion Methods

| Method                            | Converts to                      |
| --------------------------------- | -------------------------------- |
| `RatePeriodic()`                  | Periodic Rate                    |
| `RateNominal()`                   | Nominal Rate                     |
| `RateEffectyAnnually()`           | Effective Annual Rate            |
| `RatePeriodicToPeriodic()`        | Periodic to different frequency  |
| `RateNominalToNominal()`          | Nominal to different frequency   |
| `RateAnticipateEffectyAnnually()` | Anticipated Effective Annual     |
| `RateAnticipateNominal()`         | Anticipated Nominal              |
| `RateAnticipatePeriodic()`        | Anticipated Periodic             |
| `ToAnticipateNominal()`           | Ordinary to Anticipated Nominal  |
| `ToAnticipatePeriodic()`          | Ordinary to Anticipated Periodic |
| `ToNominal()`                     | Anticipated to Ordinary Nominal  |
| `ToPeriodic()`                    | Anticipated to Ordinary Periodic |

## Advanced Examples

### Converting Between Frequencies

```go
// Convert a monthly rate of 1% to its equivalent quarterly rate
rate, _ := compoundinterest.NewRateInterest(0.01, compoundinterest.Monthly, compoundinterest.RateEffectyPeriodic)
quarterly, _ := rate.RatePeriodicToPeriodic(compoundinterest.Quarterly)
fmt.Printf("Equivalent quarterly rate: %.4f\n", quarterly)
```

### Anticipated Rates

```go
// Convert an ordinary rate to its equivalent anticipated rate
rate, _ := compoundinterest.NewRateInterest(0.10, compoundinterest.Monthly, compoundinterest.RateEffectyAnnually)
anticipatedNominal, _ := rate.ToAnticipateNominal()
fmt.Printf("Anticipated nominal rate: %.4f\n", anticipatedNominal)
```

## Testing

The package includes a comprehensive test suite of 105 tests with 86.5% coverage:

```bash
# Run all tests
go test ./finance/compoundinterest -v

# Run tests with coverage
go test ./finance/compoundinterest -cover

# Run a specific test
go test -run TestFutureWithPeriodicRate ./finance/compoundinterest
```

### Test Files

- `future_test.go` - 10 tests for future value calculations
- `present_test.go` - 12 tests for present value calculations
- `rate_interest_test.go` - 13 tests for rate calculations
- `periods_test.go` - 16 tests for period calculations
- `rate_conversion_test.go` - 31 tests for rate conversions
- `utils_test.go` - 23 tests for utility functions

## Use Cases

1. **Investment Calculations**: Determine how much an investment will grow
2. **Loans**: Calculate payments or interest rates
3. **Financial Analysis**: Compare different investment options
4. **Valuation**: Calculate present value of future cash flows
5. **Discounts**: Work with anticipated rates and discounts

## Important Notes

- All rates must be entered as decimals (0.05 for 5%)
- Periods must be positive numbers
- The compounding frequency must match the unit of periods
- For accurate results, it's important to use the correct rate type

## Contributing

Contributions are welcome. Please ensure:

1. Write tests for new features
2. Maintain test coverage above 85%
3. Document the code properly
4. Follow the existing code style

## License

MIT License

## Author

Yeferson Toloza
