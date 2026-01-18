# Development Guide - Annuities Package

This guide provides detailed information for developers working with the `annuities` package.

## Project Structure

```
annuities/
├── README.md                      # Main package documentation
├── DEVELOPMENT.md                 # This guide
├── EXAMPLES.md                    # Usage examples
├── root.go                        # Definition of main types
├── future.go                      # Future value calculation
├── present.go                     # Present value calculation
├── periods.go                     # Period calculation
├── present_test.go                # Tests for annuity calculations
```

## Key Financial Concepts

### 1. Ordinary Annuity

An annuity where payments are made at the end of each period.

**Formula for Present Value:**

```
PV = PMT × [(1 + i)^n - 1] / [i(1 + i)^n]
```

### 2. Annuity Periodic Payment

The regular payment made each period to amortize a present value or accumulate a future value.

**From Present Value:**

```
PMT = PV × [i(1 + i)^n] / [(1 + i)^n - 1]
```

**From Future Value:**

```
PMT = FV × [i / ((1 + i)^n - 1)]
```

### 3. Future Value of Annuity

The accumulated value of periodic payments with compound interest.

```
FV = PMT × [((1 + i)^n - 1) / i]
```

### 4. Periodic Rate (i)

The interest rate applied in each compounding period.

**Example:**

- 12% nominal monthly = 1% periodic (12% / 12)
- 12% nominal quarterly = 3% periodic (12% / 4)

## Typical Workflow

### 1. Create a Rate

```go
rate, err := compositeinterest.NewRateInterest(
    0.01,                                    // value
    compositeinterest.Monthly,               // frequency
    compositeinterest.RateEffectyPeriodic,   // type
)
```

### 2. Create a Period

```go
period, err := compositeinterest.NewPeriod(
    12,                              // number of periods
    compositeinterest.Monthly,       // frequency
)
```

### 3. Create an Annuity Object

```go
ann, err := annuities.New(
    100,         // payment (periodic payment amount)
    0,           // present (0 if not used)
    0,           // future (0 if not used)
    period,      // period
    rate,        // interest rate
)
```

### 4. Calculate the Desired Value

```go
present, err := ann.Present()
// or
future, err := ann.Future()
// or
payment, err := ann.PaymentFromPresentValue()
// or
payment, err := ann.PaymentFromFutureValue()
// or
periods, err := ann.PeriodsWithPresent()
// or
periods, err := ann.PeriodsWithFuture()
```

## Implementation Details

### Annuity Type

```go
type Annuity struct {
    value             float64
    compositeInterest compositeinterest.CompositeInterest
}
```

- `value`: The periodic payment amount
- `compositeInterest`: Underlying composite interest calculation

### Using CompositeInterest

The annuity package leverages the `compositeinterest` package:

1. **GetEqualsRateInterestPeriods():** Converts rate to periodic and adjusts periods to match rate frequency

**Example:**

```
If rate is: 12% nominal monthly
Result periodic: 1% monthly

If period is: 12 months and rate is monthly
Adjusted period: 12 (no change, frequencies match)

If period is: 1 year and rate is monthly
Adjusted period: 12 (1 year × 12 months/year)
```

### Error Handling

Errors can occur in:

1. **Creating rate or period:** If invalid parameters
2. **Accessing composite interest methods:** If underlying data is invalid
3. **Logarithmic calculations:** In `PeriodsWithPresent()` or `PeriodsWithFuture()` if division by zero occurs

**Best practice:**

```go
present, err := ann.Present()
if err != nil {
    log.Fatal("Error calculating present value:", err)
}
```

## Testing

### Running Tests

```bash
# All tests
go test ./finance/annuities -v

# Tests from a specific file
go test -run TestPresent ./finance/annuities

# With coverage
go test ./finance/annuities -cover

# Detailed coverage
go test ./finance/annuities -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Writing New Tests

**Recommended structure:**

```go
func TestNewFunctionality(t *testing.T) {
    // Arrange: Prepare data
    rate, err := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
    require.NoError(t, err)

    period, err := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
    require.NoError(t, err)

    ann, err := annuities.New(100, 0, 0, period, rate)
    require.NoError(t, err)

    // Act: Execute the operation
    result, err := ann.Present()

    // Assert: Verify results
    require.NoError(t, err)
    assert.InDelta(t, 1125.51, result, 0.01)
}
```

**Using Testify:**

- `require.NoError()` for fatal errors
- `assert.Equal()` for exact values
- `assert.InDelta()` for values with tolerance
- `assert.True()` for boolean conditions

## Common Use Cases

### 1. Loan Amortization

```go
// $10,000 loan with 1% monthly for 12 months
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
ann, _ := annuities.New(0, 10000, 0, period, rate)
payment, _ := ann.PaymentFromPresentValue()
```

### 2. Savings Accumulation

```go
// $500 monthly savings at 0.5% monthly for 12 months
rate, _ := compositeinterest.NewRateInterest(0.005, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
ann, _ := annuities.New(500, 0, 0, period, rate)
future, _ := ann.Future()
```

### 3. Calculate Payment to Reach Goal

```go
// $100,000 goal in 5 years at 6% annual
rate, _ := compositeinterest.NewRateInterest(0.06, compositeinterest.Annually, compositeinterest.RateEffectyAnnually)
period, _ := compositeinterest.NewPeriod(5, compositeinterest.Annually)
ann, _ := annuities.New(0, 0, 100000, period, rate)
payment, _ := ann.PaymentFromFutureValue()
```

## Performance and Optimization

### Optimization Points

1. **Reuse objects:** If doing multiple calculations with the same rate/period, create the annuity once
2. **Cache conversions:** If converting the same rate multiple times, consider caching the result
3. **Avoid unnecessary creations:** Reuse rate and period objects in loops

### Efficient Example

```go
// ❌ Inefficient: Creates rate 1000 times
for i := 0; i < 1000; i++ {
    rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
    // ...
}

// ✅ Efficient: Creates once
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
for i := 0; i < 1000; i++ {
    // Reuse rate
}
```

## Debugging

### Useful Techniques

1. **Print intermediate values:**

```go
rate, _ := compositeinterest.NewRateInterest(0.01, compositeinterest.Monthly, compositeinterest.RateEffectyPeriodic)
period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
ann, _ := annuities.New(100, 0, 0, period, rate)

periods, rateInterest, _ := ann.compositeInterest.GetEqualsRateInterestPeriods()
log.Printf("Periods: %.2f, Rate: %.6f\n", periods, rateInterest)
```

2. **Verify calculations manually:**

```go
// After calculating, verify using the formula
// PV = PMT × [(1 + i)^n - 1] / [i(1 + i)^n]
i := 0.01
n := 12.0
pmt := 100.0
pow := math.Pow(1+i, n)
expectedPV := pmt * ((pow - 1) / (i * pow))
```

## Common Changes and How to Make Them

### Adding a New Method

1. **Create new file** (e.g., `rate.go`)
2. **Implement the method:**

```go
func (a Annuity) NewCalculation() (float64, error) {
    periods, rateInterest, err := a.compositeInterest.GetEqualsRateInterestPeriods()
    if err != nil {
        return 0, err
    }
    
    // Your calculation here
    result := calculateSomething(a.value, periods, rateInterest)
    
    return result, nil
}
```

3. **Write tests** in `_test.go` file
4. **Update documentation** in README.md and EXAMPLES.md

### Extending Annuity Type

If you need to add fields to the `Annuity` struct:

1. **Update the struct** in `root.go`
2. **Update `New()` function** to initialize new fields
3. **Update tests** that create Annuity objects

## Additional Resources

### Recommended Books

- "Financial Mathematics" - Héctor Vidaurri
- "Financial Engineering and Computation" - Phelim P. Boyle
- "Introduction to Financial Mathematics" - Chris Kenyon

### Useful Links

- [Go Documentation](https://golang.org/doc/)
- [Testify README](https://github.com/stretchr/testify)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Frequently Asked Questions (FAQ)

**Q: What's the difference between `PaymentFromPresentValue()` and `PaymentFromFutureValue()`?**
A: `PaymentFromPresentValue()` calculates the payment needed to pay off a present value (loan). `PaymentFromFutureValue()` calculates the payment needed to accumulate a future value (savings).

**Q: How do I know which method to use?**
A:
- If you have an initial debt/principal → Use `PaymentFromPresentValue()`
- If you have a goal amount → Use `PaymentFromFutureValue()`
- If you know the payment → Use `Present()` or `Future()` to calculate values

**Q: What's the "value" field in the Annuity struct?**
A: The "value" field is the periodic payment amount. This is what gets paid/received each period.

**Q: Why do I get an error in PeriodsWithPresent() or PeriodsWithFuture()?**
A: Probably because the payment equals the interest (no progress toward goal), or when payment is less than the interest accruing.

**Q: Can I use anticipated rates with annuities?**
A: Yes! Create a rate with anticipated rate type and pass it to the annuity. The underlying `compositeinterest` package handles the conversions.

## Contributing

To contribute to development:

1. Fork the repository
2. Create a branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add your-feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request

**Requirements for PR:**

- Tests for new features
- Updated documentation
- Test coverage maintained above 75%
- Code formatted with `gofmt`

## Support

For questions or bug reports, please open an issue in the repository.
