# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Breaking**: `money.Decimal` is now a type alias of `decimal.Decimal` instead of a wrapper type. The full arithmetic API comes directly from the `decimal` package; `money`'s `NewFrom*`/`MustFrom*` decimal constructors and `money.Zero`/`money.One` remain as deprecated forwarders. `Decimal.ToMoney` was removed — use the new `money.FromDecimal(d, currency)`
- **Breaking**: the chart rendering package moved from `finance/charts` to the separate `charts` Go module (`github.com/yeferson59/gofinance/charts`); the root library module now has zero external runtime dependencies. `examples/` is also its own module
- All `finance/*` packages now use `decimal.Decimal` directly for rates, factors and periods; `finance/tvm` and `finance/daycount` no longer depend on `money`
- `Money.Scan` now defaults the currency to USD (matching `UnmarshalJSON`) instead of leaving it unset
- Repository layering rules are documented in `ARCHITECTURE.md`

### Added
- `money.FromDecimal(d, currency)` to turn a computed decimal into a monetary amount
- Dimensionally-correct `Money.MulDecimal`, `Money.DivDecimal` and `Money.MustDivDecimal` for amount×rate / amount÷rate math without attaching a placeholder currency
- `annuities.WriteCSVTo(io.Writer, ...)` so schedule CSV export is destination-agnostic; `WriteCSV` remains as a file convenience wrapper
- `finance/investment`: date-based `XNPV`/`XIRR` for irregular cash flows (`DatedCashFlow`, Actual/365 basis) and `Perpetuity`/`GrowingPerpetuity` (Gordon model), all with `Must*` variants and typed errors (`ErrDatesBeforeBase`, `ErrNonPositiveRate`, `ErrRateBelowGrowth`)
- New `finance/bonds` package: fluent `NewBond()` builder with clean `Price` from yield, `YTM` from price (bracketed bisection), `MacaulayDuration`/`ModifiedDuration`/`Convexity`, `CouponPayment`, and `AccruedInterest` (reusing `finance/daycount`)
- New `finance/depreciation` package: `StraightLine`, `DecliningBalance`, `DoubleDecliningBalance` (with straight-line switchover), `SumOfYearsDigits`, and `MACRS` (GDS half-year tables for 3/5/7/10/15/20-year recovery), each returning a year-by-year `Schedule`
- `finance/returns`: inflation adjustment `RealValue`, `NominalValue`, and Fisher `RealRate` (with `Must*` variants and `ErrInvalidInflationRate`)
- Runnable examples for the new instruments in `examples/main.go`
- New `finance/returns` package: `CAGR`, `ROI`, `HoldingPeriodReturn`, and `Annualized` return metrics (with `Must*` variants), computed on the decimal engine
- New `finance/investment` package: `NPV` of a periodic cash-flow stream at a discount rate, and `IRR` via Newton–Raphson with a bracketed-bisection fallback (with `Must*` variants); typed errors `ErrNoCashFlows`, `ErrInvalidRate`, `ErrNoSignChange`, `ErrNoConvergence`
- New `finance/tvm` package: a general time-value-of-money solver (`SolveN`, `SolveRate`, `SolvePV`, `SolvePMT`, `SolveFV`, plus `Must*`) with a fluent `NewTVM()` builder and ordinary/annuity-due timing
- New `finance/daycount` package: `Days` and `YearFraction` under the 30/360, Actual/360, Actual/365 (Fixed), and Actual/Actual (ISDA) conventions
- Runnable examples for the new packages in `examples/main.go`, and a `FEATURES.md` roadmap of candidate features
- Fluent rate conversion builder `compositeinterest.NewRateConversion()` for converting rates between periodic, nominal, effective annual, and anticipated types, and across compounding frequencies
- `money.SafeAdd`/`money.SafeSub` currency-checked arithmetic returning `ErrCurrencyMismatch` on mixed currencies
- `simpleinterest` builder `PresentValue()` convenience method
- Manual Release GitHub Actions workflow (`workflow_dispatch`) to create a version tag and GitHub Release with a semver bump or explicit version
- New standalone `decimal` package (`github.com/yeferson59/gofinance/decimal`): the fixed-point decimal engine that used to live inside `money` is now independently importable, with error-returning `TryAdd`/`TrySub`/`TryMul` alongside the existing panicking `Add`/`Sub`/`Mul`
- `decimal.ErrLogNonPositive` and `decimal.ErrPowNegBase` errors for the domain violations of the new native `Ln`/`Log*`/`Pow` implementations
- `Decimal.Sqrt`/`MustSqrt` (in both `decimal` and `money`): direct square root via Newton's integer iteration on the exact 256-bit radicand — always correctly rounded to 19 fractional digits (exact for perfect squares), ~0.4µs and zero allocations; returns the new `decimal.ErrSqrtNegative` for negative input

### Changed
- **Precision**: `Decimal.Pow`, `Ln`, `Log`, `Log2`, and `Log10` are now computed natively on the 128-bit decimal engine (120-bit binary fixed-point internals, zero allocations) instead of round-tripping through `float64`/`math.Pow`, so results are accurate to the full 19-digit precision (previously ~15-16 digits). Integer exponents use exact binary exponentiation on 38-significant-digit intermediates: powers whose exact value fits in 38 digits (e.g. `1.05^12`) now come out exact
- **Performance**: 128-bit-divisor division (`u256/u128` and `u128/u128`) replaced bit-by-bit binary long division with normalized 2-word Knuth (algorithm D) steps built on `bits.Div64` — `Ln`/`Log10` run at ~0.2µs and `Pow` at ~0.5µs on Apple M1 (11-13× faster than the first native version), and `Decimal.Div` by divisors wider than 64 bits speeds up as a side effect
- **Behavior**: `Pow` now defines `0^0 = 1`, returns `ErrDivideByZero` for `0^negative`, `ErrPowNegBase` for a negative base with a fractional exponent, `ErrOverflow` past the representable range, and rounds to zero below it; `Ln`/`Log*` now return `ErrLogNonPositive` for non-positive input instead of a formatting error
- **Performance**: Eliminated duplicate `math.Pow()` call in annuities Present() calculation
- **Performance**: Optimized Period struct lookup from O(7) sequential checks to O(1) switch-based lookup
- **Performance**: Cached repeated exponential calculations in rate conversion methods
- **Breaking**: `simpleinterest.Periods()` and `PeriodsWithPresentAndFuture()` now return `money.Decimal` instead of `money.Money`, matching `compositeinterest.Periods()`
- **Breaking**: `SimpleConfig.Build()` and `CompositeConfig.Build()` now use value receivers so documented fluent chains compile
- **Breaking**: `money`'s decimal engine moved into the new `decimal` package; `money.ErrOverflow`, `ErrDivideByZero`, `ErrEmptyString`, `ErrInvalidFormat`, `ErrPrecOutOfRange`, and `ErrIntPartOverflow` no longer exist — use `decimal.ErrOverflow`, etc. instead. `money.Decimal` and `money.Money` keep their existing public methods, now implemented as a wrapper around `decimal.Decimal`
- `compositeinterest.NewComposite()` now defaults the rate type to `RateEffectyPeriodic`, matching `annuities.NewAnnuity()`
- `money.MoneyZero`/`money.MoneyOne` now use the `USD` constant instead of the magic number 143
- `compositeinterest.Bimonthly` constant value corrected from `"Bimontly"` to `"bimonthly"`

### Fixed
- Improved error handling consistency in NewPeriod() - now returns error for invalid compounding frequency
- Fixed MMK and MNT ISO currency code typos ("MNK", "MINT")
- `GetCurrencyPrecisionCode()` now covers all currencies per ISO 4217 minor units (default 2) instead of only USD/COP
- `GetEqualsRateInterestPeriods()` now propagates errors instead of silently returning zero values
- `annuities.Years()` now honors the configured payment frequency instead of always assuming monthly
- `annuities.AnnualRate()` now uses divisor 3 for `QuarterlyTwo` instead of 4
- `Annuity.Future()` now falls back to the payment-based formula when the compound-interest future value is not derivable

## [1.1.0] - 2026-01-06

### Added
- Comprehensive root calculation tests for more complex operations
- Data validation utilities for financial calculations

### Fixed
- Validate zero values in Periods to prevent invalid calculations
- Fixed equals periods with different compounding frequencies to include days period
- Corrected operation logic for periods with different compounding frequencies
- Improved code coverage with additional test flags

## [1.0.0] - 2026-01-04

### Added

#### Core Functionality
- **Compound Interest Calculations**: Complete package for financial calculations
  - Future Value calculation using formula: FV = PV × (1 + i)^n
  - Present Value calculation using formula: PV = FV / (1 + i)^n
  - Interest Rate calculation using formula: i = (FV/PV)^(1/n) - 1
  - Number of Periods calculation using formula: n = ln(FV/PV) / ln(1 + i)

#### Compounding Frequencies Support
- Daily (365 periods per year)
- Monthly (12 periods per year)
- Bimonthly (6 periods per year)
- Quarterly (4 or 3 periods per year)
- Semi-annually (2 periods per year)
- Annually (1 period per year)

#### Interest Rate Types Support
**Ordinary Rates (charged at period end):**
- Periodic Rate
- Nominal Annual Rate
- Effective Annual Rate

**Anticipated Rates (charged at period beginning - discounts):**
- Anticipated Periodic Rate
- Anticipated Nominal Rate
- Anticipated Effective Annual Rate

#### Rate Conversion Methods
- `RatePeriodic()` - Convert to periodic rate
- `RateNominal()` - Convert to nominal rate
- `RateEffectyAnnually()` - Convert to effective annual rate
- `RatePeriodicToPeriodic()` - Convert between different frequencies
- `RateNominalToNominal()` - Convert nominal rates between frequencies
- `RateAnticipateEffectyAnnually()` - Convert to anticipated effective annual
- `RateAnticipateNominal()` - Convert to anticipated nominal
- `RateAnticipatePeriodic()` - Convert to anticipated periodic
- `ToAnticipateNominal()` - Convert ordinary to anticipated nominal
- `ToAnticipatePeriodic()` - Convert ordinary to anticipated periodic
- `ToNominal()` - Convert anticipated to ordinary nominal
- `ToPeriodic()` - Convert anticipated to ordinary periodic

#### Type Definitions
- `CompoundingFrequency` - Enum for compounding frequencies
- `TypeRate` - Enum for rate types
- `Period` - Structure for period specifications
- `RateInterest` - Structure for rate specifications
- `CompositeInterest` - Main structure for compound interest calculations

#### Documentation
- Complete README.md with installation and usage guide
- DEVELOPMENT.md with comprehensive developer guide
- EXAMPLES.md with practical usage examples
- Inline code documentation for all public functions and types

#### Testing
- 105 comprehensive unit tests
- 86.5% code coverage
- Test files:
  - `future_test.go` - 10 tests for Future() calculations
  - `present_test.go` - 12 tests for Present() calculations
  - `rate_interest_test.go` - 13 tests for Interest() calculations
  - `periods_test.go` - 16 tests for Periods() calculations
  - `rate_conversion_test.go` - 31 tests for rate conversions
  - `utils_test.go` - 23 tests for utility functions

### Technical Details

#### Package Structure
```
finance/compositeinterest/
├── README.md              # Package documentation
├── DEVELOPMENT.md         # Developer guide
├── EXAMPLES.md           # Usage examples
├── root.go               # Main types and functions
├── consts.go             # Constants
├── data.go               # Data structures
├── future.go             # Future value calculation
├── present.go            # Present value calculation
├── rate_interest.go      # Interest rate calculation
├── periods.go            # Periods calculation
├── rate_conversion.go    # Rate conversion methods
├── utils.go              # Utility functions
├── *_test.go             # Test files
```

#### Key Features
- ✅ No external dependencies (only Go standard library)
- ✅ Full support for multiple rate types and frequencies
- ✅ Comprehensive error handling
- ✅ Efficient mathematical calculations
- ✅ Well-documented code with examples
- ✅ 100% English documentation
- ✅ Production-ready code quality

### Dependencies
- Go 1.18 or higher
- Testing: `github.com/stretchr/testify` v1.x (test dependency only)

### Known Limitations
None

### Performance
- All calculations are performed using standard Go math library
- No external dependencies for core functionality
- Average calculation time: < 1 microsecond per operation

### Backward Compatibility
- Initial release, no previous versions to compare

### Contributors
- Yeferson Toloza (Initial development)

---

## Release Information

**Release Date:** January 4, 2026
**Tag:** v1.0.0
**Status:** Stable

This is the first stable release of the gofinance composite interest package, featuring a complete implementation of compound interest calculations with support for multiple rate types and compounding frequencies.
