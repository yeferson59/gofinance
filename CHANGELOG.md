# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- **Performance**: Eliminated duplicate `math.Pow()` call in annuities Present() calculation
- **Performance**: Optimized Period struct lookup from O(7) sequential checks to O(1) switch-based lookup
- **Performance**: Cached repeated exponential calculations in rate conversion methods

### Fixed
- Improved error handling consistency in NewPeriod() - now returns error for invalid compounding frequency

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
