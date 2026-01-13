# Development Guide - Composite Interest Package

This guide provides detailed information for developers working with the `compositeinterest` package.

## Project Structure

```
compositeinterest/
├── README.md                      # Main package documentation
├── DEVELOPMENT.md                 # This guide
├── EXAMPLES.md                    # Usage examples
├── root.go                        # Definition of main types
├── consts.go                      # Constants for frequencies and rate types
├── data.go                        # Mapping of compounding frequencies
├── future.go                      # Future value calculation
├── present.go                     # Present value calculation
├── rate_interest.go               # Interest rate calculation
├── periods.go                     # Period calculation
├── rate_conversion.go             # Conversions between rate types
├── utils.go                       # Utility functions
├── future_test.go                 # Tests for Future()
├── present_test.go                # Tests for Present()
├── rate_interest_test.go          # Tests for Interest()
├── periods_test.go                # Tests for Periods()
├── rate_conversion_test.go        # Tests for conversions
└── utils_test.go                  # Tests for utility functions
```

## Key Financial Concepts

### 1. Periodic Rate (i)

The interest rate applied in each compounding period.

**Example:**

- 12% nominal monthly = 1% periodic (12% / 12)
- 12% nominal quarterly = 3% periodic (12% / 4)

**Usage in formulas:**

```
FV = PV × (1 + i)^n
```

### 2. Nominal Rate (j)

Annual rate divided by the compounding frequency.
Does not account for the effects of multiple compounding.

**Conversion:**

```
i = j / m  (where m is the frequency per year)
```

### 3. Effective Annual Rate (i_a)

The actual rate earned in one year considering compounding.

**Conversion from nominal:**

```
i_a = (1 + j/m)^m - 1
```

**Conversion from periodic:**

```
i_a = (1 + i)^m - 1
```

### 4. Anticipated Rates

Charged at the beginning of the period (discounts).
More common in bills of exchange and securities.

**Relationship with ordinary rates:**

```
d_a = i_a / (1 + i_a)  (where d_a is annual anticipated rate)
```

## Typical Workflow

### 1. Create a Rate

```go
rate, err := compositeinterest.NewRateInterest(
    0.12,                                    // value
    compositeinterest.Monthly,               // frequency
    compositeinterest.RateEffectyNominal,    // type
)
```

### 2. Create a Period

```go
period, err := compositeinterest.NewPeriod(
    12,                              // number of periods
    compositeinterest.Monthly,       // frequency
)
```

### 3. Create a CompositeInterest Object

```go
ci, err := compositeinterest.New(
    1000,        // present (present value)
    0,           // future (0 if to be calculated)
    rate,        // interest rate
    period,      // period
)
```

### 4. Calculate the Desired Value

```go
future, err := ci.Future()
// or
present, err := ci.Present()
// or
interest, err := ci.Interest()
// or
periods, err := ci.Periods()
```

## Implementation Details

### Internal Rate Conversion

The `getEqualsRateInterestPeriods()` method is crucial:

1. **Converts rate to periodic:** Any rate type is converted to periodic
2. **Adjusts periods:** If the period frequency doesn't match the rate frequency, the number of periods is multiplied

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

1. **getPeriod():** If no period is set
2. **getCompoundingFrequency():** If the frequency doesn't exist
3. **Logarithmic calculations:** In `Periods()` if division by zero occurs

**Best practice:**

```go
future, err := ci.Future()
if err != nil {
    log.Fatal("Error calculating future value:", err)
}
```

## Testing

### Running Tests

```bash
# All tests
go test ./finance/compositeinterest -v

# Tests from a specific file
go test -run TestFuture ./finance/compositeinterest

# With coverage
go test ./finance/compositeinterest -cover

# Detailed coverage
go test ./finance/compositeinterest -coverprofile=coverage.out
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

    // Act: Execute the operation
    ci, err := compositeinterest.New(1000, 0, rate, period)
    require.NoError(t, err)

    result, err := ci.Future()

    // Assert: Verify results
    require.NoError(t, err)
    assert.InDelta(t, 1126.83, result, 0.01)
}
```

**Using Testify:**

- `require.NoError()` for fatal errors
- `assert.Equal()` for exact values
- `assert.InDelta()` for values with tolerance
- `assert.True()` for boolean conditions

## Common Use Cases

### 1. Loan with Monthly Nominal Rate

```go
// $10,000 loan with 2% nominal monthly for 24 months
rate, _ := compositeinterest.NewRateInterest(0.24, compositeinterest.Monthly, compositeinterest.RateEffectyNominal)
period, _ := compositeinterest.NewPeriod(24, compositeinterest.Monthly)
ci, _ := compositeinterest.New(10000, 0, rate, period)
amount, _ := ci.Future()
```

### 2. Investment with Effective Annual Rate

```go
// $5,000 investment with 8% effective annual for 3 years
rate, _ := compositeinterest.NewRateInterest(0.08, compositeinterest.Annually, compositeinterest.RateEffectyAnnually)
period, _ := compositeinterest.NewPeriod(3, compositeinterest.Annually)
ci, _ := compositeinterest.New(5000, 0, rate, period)
future, _ := ci.Future()
```

### 3. Discount Bill of Exchange

```go
// $50,000 bill discounted at 10% annual anticipated rate
rate, _ := compositeinterest.NewRateInterest(0.10, compositeinterest.Annually, compositeinterest.RateAnticipateEffectyAnnually)
period, _ := compositeinterest.NewPeriod(1, compositeinterest.Annually)
ci, _ := compositeinterest.New(0, 50000, rate, period)
present, _ := ci.Present()
```

## Performance and Optimization

### Optimization Points

1. **Reuse objects:** If doing multiple calculations with the same rate/period, create the structure once
2. **Cache conversions:** If converting the same rate multiple times, consider caching the result
3. **Avoid unnecessary creations:** Reuse `RateInterest` and `Period` in loops

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
rate, _ := NewRateInterest(0.12, Monthly, RateEffectyNominal)
periodic, _ := rate.RatePeriodic()
log.Printf("Periodic rate: %.6f\n", periodic)
```

2. **Verify conversions:**

```go
original := 0.01
periodic, _ := rate.RatePeriodic()
nominal, _ := rate.RateNominal()
fmt.Printf("Original: %.6f\nPeriodic: %.6f\nNominal: %.6f\n", original, periodic, nominal)
```

3. **Use test helper for debugging:**

```go
func debugRate(t *testing.T, rate *RateInterest) {
    periodic, _ := rate.RatePeriodic()
    nominal, _ := rate.RateNominal()
    annual, _ := rate.RateEffectyAnnually()
    t.Logf("Periodic: %.6f, Nominal: %.6f, Annual: %.6f", periodic, nominal, annual)
}
```

## Common Changes and How to Make Them

### Adding a New Frequency

1. **Add constant in `consts.go`:**

```go
const Weekly CompoundingFrequency = "weekly"
```

2. **Update `data.go`:**

```go
var countCompoundingFrequency = map[CompoundingFrequency]float64{
    // ...
    Weekly: 52,
    // ...
}
```

3. **Update `root.go` (Period struct if necessary)**

4. **Write tests:**

```go
func TestGetCompoundingFrequencyWeekly(t *testing.T) {
    value, err := getCompoundingFrequency(Weekly)
    require.NoError(t, err)
    assert.Equal(t, 52.0, value)
}
```

### Adding a New Conversion Function

1. **Add in `rate_conversion.go`:**

```go
func (rt *RateInterest) NewConversion() (float64, error) {
    // implementation
    return result, nil
}
```

2. **Write comprehensive tests**
3. **Document in comments**
4. **Update this guide file**

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

**Q: What's the difference between nominal and periodic rate?**
A: The nominal rate is annual (j) and divided by frequency to get periodic (i = j/m).

**Q: When should I use anticipated rate?**
A: In bill discounts, documentary credit operations, and other instruments where interest is charged at the beginning.

**Q: How do I know what rate type I have?**
A: If given an annual percentage, it's likely nominal. If given with a specific frequency, it's periodic.

**Q: Why do I get an error in Periods()?**
A: Probably because present value = future value (no gain), or the period isn't set.

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
- Test coverage > 85%
- Code formatted with `gofmt`

## Support

For questions or bug reports, please open an issue in the repository.
