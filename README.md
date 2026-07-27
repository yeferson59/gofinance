[![CI](https://github.com/yeferson59/gofinance/actions/workflows/ci.yml/badge.svg)](https://github.com/yeferson59/gofinance/actions/workflows/ci.yml)
[![codecov](https://codecov.io/github/yeferson59/gofinance/graph/badge.svg?token=OQAWSDRWEZ)](https://codecov.io/github/yeferson59/gofinance)

# 💰 GoFinance

A robust, type-safe Go library for financial calculations and money management. **GoFinance** provides comprehensive tools for handling complex financial mathematics including simple interest, compound interest, annuities, and precise monetary operations — all built on its own zero-dependency, allocation-free decimal engine.

---

## ✨ Features

### 🔢 **Decimal Engine** (`decimal`)

- Fixed-point decimal type backed by a 128-bit coefficient: up to 19 digits after the decimal point, no `float64` rounding surprises
- Zero external dependencies and zero heap allocations for arithmetic
- Native mathematical functions computed directly on the decimal representation:
  - `Pow` — exact binary exponentiation for integer exponents (e.g. `(1+i)^n`), 120-bit fixed-point `exp(e·ln x)` for fractional ones
  - `Sqrt` — Newton's integer iteration on the exact 256-bit radicand, always correctly rounded (exact for perfect squares)
  - `Ln`, `Log`, `Log2`, `Log10` — accurate to the full 19-digit precision
- Banker's, half-away, half-toward, truncating and away-from-zero rounding modes
- JSON marshaling/unmarshaling and exact string round-tripping

### 💵 **Money Management** (`money`)

- Precise monetary amounts built on the `decimal` engine
- Multi-currency support with ISO 4217 codes, symbols, and per-currency precision
- Currency-checked arithmetic: `Add`/`Sub` panic on a currency mismatch (like the decimal engine does on overflow); `TryAdd`/`TrySub` return `ErrCurrencyMismatch` instead
- Penny-exact allocation: split amounts by ratios or evenly without losing a cent (`Allocate`, `AllocateEvenly`)
- Currency conversion with explicit exchange rates (`Convert`, `ConvertFloat64`)

### 📊 **Simple Interest** (`finance/simpleinterest`)

- Future value, present value, interest, rate, and period calculations
- Fluent builder API with day/week/month/year time units

### 🔄 **Compound Interest** (`finance/compoundinterest`)

- Future/present value and period calculations with flexible compounding frequencies (daily, monthly, bimonthly, quarterly, semi-annual, annual)
- Rate conversions between periodic, nominal, effective annual, and anticipated rate types

### 📈 **Annuities** (`finance/annuities`)

- Payment, present/future value, rate, and period calculations
- Full amortization schedule generation (`BuildSchedule`)
- Optional chart rendering of schedules via the separate `charts` module (go-echarts)

### 📊 **Returns** (`finance/returns`)

- Compound annual growth rate (`CAGR`), simple return on investment (`ROI`), holding-period return, and annualization of a cumulative return — all on the decimal engine
- Inflation adjustment: real vs nominal value and the Fisher real rate (`RealValue`, `NominalValue`, `RealRate`)
- Portfolio performance: time-weighted return over valuation subperiods (`TimeWeightedReturn`, `ChainReturns`) vs money-weighted return (`MoneyWeightedReturn`)
- Risk metrics: sample and population variance/volatility, square-root-of-time annualization, and the Sharpe ratio (`Volatility`, `AnnualizedVolatility`, `SharpeRatio`, `AnnualizedSharpeRatio`)

### 🏠 **Loans** (`finance/loans`)

- Level payment, effective annual rate, and the fee-inclusive `APR`/`EffectiveAPR` solved from the loan's net proceeds
- Extra payments and early payoff: how many payments are avoided and how much interest is never charged (`Payoff`, `Savings`)
- Refinance comparator: payment delta, break-even period on closing costs, and the net present value of the switch (`Compare`)

### 💹 **Investment** (`finance/investment`)

- Net present value (`NPV`) of a periodic cash-flow stream at a discount rate
- Internal rate of return (`IRR`) via Newton–Raphson with a bracketed-bisection fallback
- Date-based `XNPV`/`XIRR` for irregular cash flows (Actual/365 basis)
- Level and growing (`Perpetuity`, `GrowingPerpetuity`, Gordon model) perpetuities

### 🏦 **Bonds** (`finance/bonds`)

- Clean price from yield, yield to maturity from price (`Price`, `YTM`)
- Macaulay and modified duration, convexity, and accrued interest

### 🏭 **Depreciation** (`finance/depreciation`)

- Straight-line, declining balance, double declining (with straight-line switch), sum-of-years'-digits, and MACRS (GDS) schedules

### 🧮 **Time Value of Money** (`finance/tvm`)

- General financial-calculator solver: set any four of N, rate, PV, PMT, FV and solve for the fifth, with ordinary-annuity and annuity-due timing

### 📆 **Day-Count Conventions** (`finance/daycount`)

- Year fractions and day counts under 30/360, Actual/360, Actual/365 (Fixed), and Actual/Actual (ISDA)

### 🗓️ **Shared Time Vocabulary** (`finance/term`)

- One set of types for time units (`term.Unit`) and compounding/payment cadences (`term.Frequency`, with `PeriodsPerYear` and `MonthsPerPeriod`), shared by the interest packages

See [`FEATURES.md`](FEATURES.md) for the full roadmap of proposed and in-progress features.

---

## 🚀 Getting Started

### Prerequisites

- Go 1.26.5 or higher

### Installation

```bash
go get github.com/yeferson59/gofinance/v2
```

### Quick Example: Money Management

```go
package main

import (
    "fmt"

    "github.com/yeferson59/gofinance/v2/money"
)

func main() {
    // $100.00 from an integer amount (10000) and a precision (2)
    m, err := money.New(10000, 2, money.USD)
    if err != nil {
        panic(err)
    }

    str, err := m.StringMoney()
    if err != nil {
        panic(err)
    }
    fmt.Println(str) // USD 100.00

    // Split $100.00 in the ratio 1:1:1 without losing a cent
    parts, err := m.Allocate(1, 1, 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(parts[0].StringFixed(2), parts[1].StringFixed(2), parts[2].StringFixed(2))
    // 33.34 33.33 33.33
}
```

### Quick Example: Decimal Math

```go
package main

import (
    "fmt"

    "github.com/yeferson59/gofinance/v2/decimal"
)

func main() {
    // Exact integer power: 30-year monthly compounding factor
    growth := decimal.MustFromString("1.005").
        MustPow(decimal.MustFromString("360"))
    fmt.Println(growth) // 6.0225752122632161841

    // Correctly rounded square root
    fmt.Println(decimal.MustFromString("2").MustSqrt())
    // 1.4142135623730950488

    // Natural logarithm at full 19-digit precision
    fmt.Println(decimal.MustFromString("2").MustLn())
    // 0.6931471805599453094
}
```

### Quick Example: Compound Interest

```go
package main

import (
    "fmt"

    "github.com/yeferson59/gofinance/v2/finance/compoundinterest"
    "github.com/yeferson59/gofinance/v2/money"
)

func main() {
    ci := compoundinterest.NewCompound().
        Present(1000, money.USD).
        Rate(0.05).
        Periods(12).
        Monthly().
        RateType(compoundinterest.RateEffectyPeriodic).
        MustBuild()

    future, err := ci.Future()
    if err != nil {
        panic(err)
    }
    fmt.Println("Future value:", future.StringFixed(2))
}
```

### Quick Example: Annuities

```go
package main

import (
    "fmt"

    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/money"
)

func main() {
    // Monthly payment for a $300,000 loan at 6% over 360 months
    payment := annuities.NewAnnuity().
        Present(300000, money.USD).
        AnnualRate(0.06).
        Periods(360).
        Monthly().
        MustPayment()

    fmt.Println("Monthly payment:", payment.StringFixed(2))
}
```

More runnable examples live in [`examples/main.go`](examples/main.go).

---

## 🎯 Precision & Performance

The `decimal` engine performs every operation on a 128-bit integer coefficient, so results carry up to 19 exact fractional digits — where `float64`-based math starts drifting after ~15-16 significant digits:

| Expression | `float64`            | GoFinance `decimal`             |
| ---------- | -------------------- | ------------------------------- |
| `1.005^12` | `1.0616778118644976` | `1.0616778118644995688` (exact) |
| `ln(2)`    | `0.6931471805599453` | `0.6931471805599453094`         |
| `sqrt(2)`  | `1.4142135623730951` | `1.4142135623730950488`         |

Internally, `Ln`/`Log*`/`Pow` run on 120-bit binary fixed-point kernels (~36 internal decimal digits), integer `Pow` uses exact 38-significant-digit binary exponentiation, and `Sqrt` computes the integer square root of the exact 256-bit radicand — so it is _always_ correctly rounded.

Indicative benchmarks (Apple M1, `make bench`):

| Operation             | Time        | Allocations |
| --------------------- | ----------- | ----------- |
| `Add` / `Mul` / `Cmp` | ~20 ns      | 0           |
| `Div`                 | ~60 ns      | 0           |
| `Parse`               | ~70 ns      | 0           |
| `Sqrt`                | ~0.2-0.4 µs | 0           |
| `Ln` / `Log10`        | ~0.2-0.7 µs | 0           |
| `Pow`                 | ~0.5-0.8 µs | 0           |

---

## 📦 Project Structure

```
gofinance/
├── decimal/                        # Fixed-point decimal engine (stdlib only)
│   ├── decimal128.go              # Core 128-bit coefficient arithmetic
│   ├── math.go                    # Native Pow, Sqrt, Ln, Log, Log2, Log10
│   ├── u128.go / u256.go          # 128/256-bit integer primitives
│   └── root.go                    # Public Decimal API
│
├── money/                          # Money and currency handling
│   ├── money.go / root.go         # Core Money type and operations
│   ├── currency.go                # ISO 4217 currencies, symbols, precision
│   ├── allocate.go                # Penny-exact allocation
│   ├── convert.go                 # Currency conversion
│   └── decimal.go                 # Deprecated aliases into decimal + FromDecimal
│
├── finance/
│   ├── simpleinterest/            # Simple interest (fluent builder)
│   ├── compoundinterest/         # Compound interest and rate conversions
│   ├── annuities/                 # Annuities and amortization schedules
│   ├── returns/                   # CAGR, ROI, inflation, TWR/MWR, volatility
│   ├── investment/                # NPV/IRR, XNPV/XIRR, perpetuities
│   ├── loans/                     # APR with fees, early payoff, refinancing
│   ├── bonds/                     # Bond pricing, YTM, duration, convexity
│   ├── depreciation/              # Straight-line, declining balance, MACRS
│   ├── tvm/                       # General time-value-of-money solver
│   ├── daycount/                  # Day-count conventions (30/360, Actual/*, …)
│   └── term/                      # Shared time vocabulary (Unit, Frequency)
│
├── charts/                         # Separate module: chart rendering (go-echarts)
├── benchmarks/                     # Cross-package benchmark suites
├── examples/                       # Separate module: runnable usage examples
├── ARCHITECTURE.md                 # Layering rules and design decisions
├── Makefile                        # Development tasks
└── .golangci.yaml                  # Linting configuration
```

The repository hosts three Go modules: the root library module (zero
runtime dependencies), `charts` (opt-in go-echarts rendering), and
`examples` (a runnable showcase). See [`ARCHITECTURE.md`](ARCHITECTURE.md)
for the layering rules between `decimal`, `money`, and `finance/*`.

---

## 🛠️ Development

### Setup

```bash
# Clone the repository
git clone https://github.com/yeferson59/gofinance.git
cd gofinance

# Install dependencies
go mod download
```

### Available Commands

```bash
# Run linting
make lint

# Run tests with coverage report
make test

# Format code and run imports
make fmt

# Run benchmarks
make bench
```

### Running Tests

```bash
go test -v ./...
```

### Code Coverage

```bash
go test -coverpkg=./... -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 📚 Key Concepts

### Simple Interest

Simple interest is calculated on the principal amount only. The formula is:

- **Interest = Principal × Rate × Time**

Use this when interest is not compounded.

### Compound Interest

Compound interest is calculated on the principal plus accumulated interest. It's more common in real-world scenarios:

- **A = P(1 + r/n)^(nt)**

Where:

- A = Final amount
- P = Principal
- r = Annual rate
- n = Compounding frequency
- t = Time in years

GoFinance evaluates the `(1 + r/n)^(nt)` factor with exact binary exponentiation on the decimal engine, so compounding factors don't accumulate floating-point drift.

### Annuities

An annuity is a series of equal payments made at regular intervals. Useful for loans, pensions, and investments.

---

## 🔐 Concurrency

`Decimal` and `Money` are immutable value types: every operation returns a new value and never mutates its receiver. They can be shared freely across goroutines without locks.

---

## 📖 Dependencies

The library module (`decimal`, `money`, `finance/...`) depends only on the Go standard library at runtime.

- **github.com/go-echarts/go-echarts/v2** — used exclusively by the optional `charts` module
- **github.com/stretchr/testify** — testing utilities (dev dependency)

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is open source and available under the terms specified in the repository.

---

## 👨‍💻 Author

**Yeferson Toloza**

- GitHub: [@yeferson59](https://github.com/yeferson59)
- Project: [GoFinance](https://github.com/yeferson59/gofinance)

---

## 📞 Support

If you have questions or issues:

- Open an issue on GitHub
- Check existing documentation ([`CHANGELOG.md`](CHANGELOG.md), [`PERFORMANCE.md`](PERFORMANCE.md))
- Review the test files for usage examples

---

## 🎯 Roadmap

See [`FEATURES.md`](FEATURES.md) for the full, prioritized list.

- [x] Additional financial instruments (NPV/IRR, returns, TVM solver, day-count conventions)
- [x] Bond valuation (YTM, duration, convexity) and XNPV/XIRR
- [x] Depreciation schedules (straight-line, declining balance, MACRS)
- [x] Performance optimizations (native decimal math, Knuth division, zero-allocation kernels)
- [ ] Extended documentation with more examples
- [ ] CLI tools for quick calculations
- [ ] Web API wrapper

---

**Happy calculating! 🧮**
