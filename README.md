# 💰 GoFinance

A robust, type-safe Go library for financial calculations and money management. **GoFinance** provides comprehensive tools for handling complex financial mathematics including simple interest, compound interest, annuities, and precise monetary operations.

---

## ✨ Features

### 💵 **Money Management**

- Precise decimal-based money handling using `udecimal` for accuracy
- Multi-currency support with ISO code standards
- Thread-safe operations with built-in mutex protection
- Clean string representation of monetary values

### 📊 **Simple Interest Calculations**

- Calculate future and present values
- Determine interest rates and periods
- Support for multiple time units:
  - Days, Weeks, Months, Years
- Comprehensive interest computation

### 🔄 **Compound Interest Calculations**

- Complex interest rate calculations with flexible compounding frequencies
- Support for multiple compounding periods:
  - Daily, Monthly, Bimonthly
  - Quarterly (two variations), Semi-annually, Annually
- Rate conversion and period calculations
- Advanced financial data structures

### 📈 **Annuities**

- Annuity calculations for financial planning
- Support for various annuity types and scenarios

---

## 🚀 Getting Started

### Prerequisites

- Go 1.26.1 or higher
- Basic knowledge of financial mathematics

### Installation

```bash
go get github.com/yeferson59/gofinance
```

### Quick Example: Money Management

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/money"
)

func main() {
    // Create money with USD currency
    m, err := money.New(10000, 2, money.USD)
    if err != nil {
        panic(err)
    }

    // Get string representation
    str, err := m.String()
    if err != nil {
        panic(err)
    }

    fmt.Println(str) // Output: 100.00 USD
}
```

### Quick Example: Simple Interest

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/simpleinterest"
)

func main() {
    // Create a period of 2 years
    period := simpleinterest.NewPeriod(2, simpleinterest.Years)

    // Create simple interest calculation
    si := simpleinterest.New(
        future,
        present,
        interest,
        rateInterest,
        period,
    )

    periods, err := si.GetPeriods()
    if err != nil {
        panic(err)
    }

    fmt.Println("Periods:", periods)
}
```

### Quick Example: Compound Interest

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Create monthly compounding periods
    period, err := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
    if err != nil {
        panic(err)
    }

    // Create interest rate with monthly compounding
    rateInterest, err := compositeinterest.NewRateInterest(
        0.08,
        compositeinterest.Monthly,
        compositeinterest.Nominal,
    )
    if err != nil {
        panic(err)
    }

    // Create compound interest calculation
    ci, err := compositeinterest.New(
        1000,    // present value
        1500,    // future value
        rateInterest,
        period,
    )
    if err != nil {
        panic(err)
    }

    fmt.Println("Compound interest:", ci)
}
```

---

## 📦 Project Structure

```
gofinance/
├── money/                          # Money and currency handling
│   ├── money.go                   # Core Money struct
│   ├── currency.go                # Currency definitions
│   └── root.go                    # Root package definitions
│
├── finance/                        # Financial calculations
│   ├── simpleinterest/            # Simple interest calculations
│   │   ├── interest.go            # Interest computation
│   │   ├── future.go              # Future value calculations
│   │   ├── present.go             # Present value calculations
│   │   ├── rate_interest.go       # Interest rate calculations
│   │   ├── periods.go             # Period definitions
│   │   └── *_test.go              # Comprehensive tests
│   │
│   ├── compositeinterest/         # Compound interest calculations
│   │   ├── root.go                # Core structures
│   │   ├── future.go              # Future value calculations
│   │   ├── present.go             # Present value calculations
│   │   ├── rate_interest.go       # Rate conversion
│   │   ├── rate_conversion.go     # Period conversions
│   │   ├── periods.go             # Compounding frequencies
│   │   ├── utils.go               # Utility functions
│   │   └── data.go                # Data structures
│   │
│   └── annuities/               # Annuity calculations
│       └── root.go                # Annuity definitions
│
├── go.mod                          # Module definition
├── go.sum                          # Dependencies hash
├── Makefile                        # Development tasks
├── .golangci.yaml                 # Linting configuration
└── README.md                       # This file
```

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

### Annuities

An annuity is a series of equal payments made at regular intervals. Useful for loans, pensions, and investments.

---

## 🔐 Thread Safety

The Money struct uses `sync.RWMutex` to ensure thread-safe operations. You can safely use Money objects in concurrent applications.

```go
var m *money.Money
var wg sync.WaitGroup

for i := 0; i < 100; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        str, _ := m.String() // Safe concurrent access
    }()
}
wg.Wait()
```

---

## 📖 Dependencies

- **github.com/quagmt/udecimal** - Decimal arithmetic for precise monetary calculations
- **github.com/stretchr/testify** - Testing utilities (dev dependency)

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
- Check existing documentation
- Review the test files for usage examples

---

## 🎯 Roadmap

- [ ] Additional financial instruments support
- [ ] Performance optimizations
- [ ] Extended documentation with more examples
- [ ] CLI tools for quick calculations
- [ ] Web API wrapper

---

**Happy calculating! 🧮**
