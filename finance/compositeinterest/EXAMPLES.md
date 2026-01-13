# Usage Examples - Composite Interest Package

This file contains practical examples of how to use the `compositeinterest` package.

## Table of Contents

1. [Basic Examples](#basic-examples)
2. [Rate Conversion](#rate-conversion)
3. [Financial Use Cases](#financial-use-cases)
4. [Advanced Examples](#advanced-examples)

---

## Basic Examples

### 1. Simple Future Value Calculation

```go
package main

import (
    "fmt"
    "log"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: $1,000 invested at 1% monthly for 12 months

    // Create rate: 1% periodic monthly
    rate, err := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Create period: 12 months
    period, err := compositeinterest.NewPeriod(12, compositeinterest.Monthly)
    if err != nil {
        log.Fatal(err)
    }

    // Create CompositeInterest object
    // Present: 1000, Future: 0 (to calculate)
    ci, err := compositeinterest.New(1000, 0, rate, period)
    if err != nil {
        log.Fatal(err)
    }

    // Calculate future value
    future, err := ci.Future()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Initial capital: $1,000.00\n")
    fmt.Printf("Rate: 1%% monthly\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Future value: $%.2f\n", future)
    // Output:
    // Initial capital: $1,000.00
    // Rate: 1% monthly
    // Period: 12 months
    // Future value: $1,126.83
}
```

### 2. Calculate Present Value

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: I want to have $10,000 in 2 years
    // How much should I invest today at 5% annual?

    rate, _ := compositeinterest.NewRateInterest(
        0.05,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    period, _ := compositeinterest.NewPeriod(2, compositeinterest.Annually)

    // Present: 0 (to calculate), Future: 10000
    ci, _ := compositeinterest.New(0, 10000, rate, period)

    present, _ := ci.Present()

    fmt.Printf("Desired future value: $10,000.00\n")
    fmt.Printf("Rate: 5%% annual\n")
    fmt.Printf("Period: 2 years\n")
    fmt.Printf("Should invest today: $%.2f\n", present)
    // Output: Should invest today: $9,070.29
}
```

### 3. Calculate Interest Rate

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: I invested $1,000 and after 12 months I have $1,126.83
    // What's the monthly rate?

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

    ci, _ := compositeinterest.New(1000, 1126.83, rate, period)

    interest, _ := ci.Interest()

    fmt.Printf("Initial capital: $1,000.00\n")
    fmt.Printf("Final capital: $1,126.83\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Periodic rate: %.4f (%.2f%%)\n", interest, interest*100)
    // Output: Periodic rate: 0.0100 (1.00%)
}
```

### 4. Calculate Number of Periods

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: With $1,000 and 5% periodic, how many months to reach $1,500?

    rate, _ := compositeinterest.NewRateInterest(
        0.05,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(1, compositeinterest.Monthly)

    ci, _ := compositeinterest.New(1000, 1500, rate, period)

    periods, _ := ci.Periods()

    fmt.Printf("Initial capital: $1,000.00\n")
    fmt.Printf("Target capital: $1,500.00\n")
    fmt.Printf("Rate: 5%% periodic\n")
    fmt.Printf("Periods needed: %.2f months\n", periods)
    // Output: Periods needed: 8.11 months
}
```

---

## Rate Conversion

### 1. Convert Nominal to Periodic

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: Annual nominal rate of 12% compounded monthly
    // What's the monthly periodic rate?

    rate, _ := compositeinterest.NewRateInterest(
        0.12,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyNominal,
    )

    // Convert to periodic
    periodic, _ := rate.RatePeriodic()

    fmt.Printf("Annual nominal rate: 12%%\n")
    fmt.Printf("Frequency: Monthly\n")
    fmt.Printf("Periodic rate: %.4f (%.2f%%)\n", periodic, periodic*100)
    // Output: Periodic rate: 0.0100 (1.00%)
}
```

### 2. Convert to Effective Annual Rate

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: What's the effective annual rate equivalent
    // of a 12% nominal rate compounded monthly?

    rate, _ := compositeinterest.NewRateInterest(
        0.12,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyNominal,
    )

    annual, _ := rate.RateEffectyAnnually()

    fmt.Printf("Annual nominal rate: 12%%\n")
    fmt.Printf("Effective annual rate: %.4f (%.2f%%)\n", annual, annual*100)
    // Output: Effective annual rate: 0.1268 (12.68%)
}
```

### 3. Convert Between Frequencies

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: I have a monthly rate of 1%
    // What's its equivalent quarterly rate?

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    // Convert to equivalent quarterly rate
    quarterly, _ := rate.RatePeriodicToPeriodic(compositeinterest.QuarterlyOne)

    fmt.Printf("Monthly rate: 1%%\n")
    fmt.Printf("Equivalent quarterly rate: %.4f (%.2f%%)\n", quarterly, quarterly*100)
    // Output: Equivalent quarterly rate: 0.0303 (3.03%)
}
```

### 4. Convert to Anticipated Rate

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: I have an ordinary rate of 10% annual
    // What's its equivalent anticipated (discount) rate?

    rate, _ := compositeinterest.NewRateInterest(
        0.10,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    // Convert to anticipated nominal rate
    anticipated, _ := rate.ToAnticipateNominal()

    fmt.Printf("Ordinary rate: 10%% annual\n")
    fmt.Printf("Equivalent anticipated rate: %.4f (%.2f%%)\n", anticipated, anticipated*100)
    // Output: Equivalent anticipated rate: 0.0909 (9.09%)

    // Verify reverse conversion
    rate2, _ := compositeinterest.NewRateInterest(
        anticipated,
        compositeinterest.Annually,
        compositeinterest.RateAnticipateEffectyNominal,
    )
    back, _ := rate2.ToNominal()
    fmt.Printf("Reverse conversion: %.4f (matches 10%%)\n", back)
}
```

---

## Financial Use Cases

### 1. Mortgage Simulation

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: $300,000 mortgage at 4.5% annual for 30 years
    // (simplified, without periodic payments)

    rate, _ := compositeinterest.NewRateInterest(
        0.045,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    period, _ := compositeinterest.NewPeriod(30, compositeinterest.Annually)

    ci, _ := compositeinterest.New(300000, 0, rate, period)

    totalDebt, _ := ci.Future()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("       MORTGAGE SIMULATION\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Initial capital: $%,.2f\n", 300000.0)
    fmt.Printf("Annual rate: 4.5%%\n")
    fmt.Printf("Term: 30 years\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Total debt at end: $%,.2f\n", totalDebt)
    fmt.Printf("Total interest paid: $%,.2f\n", totalDebt-300000)
    fmt.Printf("═════════════════════════════════════════\n")
}
```

### 2. Compare Investment Options

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    capital := 10000.0
    years := 5

    // Option 1: 6% nominal monthly
    rate1, _ := compositeinterest.NewRateInterest(
        0.06,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyNominal,
    )
    period1, _ := compositeinterest.NewPeriod(float64(years*12), compositeinterest.Monthly)
    ci1, _ := compositeinterest.New(capital, 0, rate1, period1)
    future1, _ := ci1.Future()

    // Option 2: 6.17% effective annual
    rate2, _ := compositeinterest.NewRateInterest(
        0.0617,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )
    period2, _ := compositeinterest.NewPeriod(float64(years), compositeinterest.Annually)
    ci2, _ := compositeinterest.New(capital, 0, rate2, period2)
    future2, _ := ci2.Future()

    // Option 3: 0.5% periodic monthly
    rate3, _ := compositeinterest.NewRateInterest(
        0.005,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )
    period3, _ := compositeinterest.NewPeriod(float64(years*12), compositeinterest.Monthly)
    ci3, _ := compositeinterest.New(capital, 0, rate3, period3)
    future3, _ := ci3.Future()

    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("     INVESTMENT OPTIONS COMPARISON\n")
    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("Initial capital: $%,.2f\n", capital)
    fmt.Printf("Period: %d years\n\n", years)

    fmt.Printf("OPTION 1: 6%% nominal monthly\n")
    fmt.Printf("  Future value: $%,.2f\n", future1)
    fmt.Printf("  Gain: $%,.2f\n\n", future1-capital)

    fmt.Printf("OPTION 2: 6.17%% effective annual\n")
    fmt.Printf("  Future value: $%,.2f\n", future2)
    fmt.Printf("  Gain: $%,.2f\n\n", future2-capital)

    fmt.Printf("OPTION 3: 0.5%% periodic monthly\n")
    fmt.Printf("  Future value: $%,.2f\n", future3)
    fmt.Printf("  Gain: $%,.2f\n", future3-capital)

    fmt.Printf("─────────────────────────────────────────────────────\n")
    if future1 > future2 && future1 > future3 {
        fmt.Printf("✓ BEST OPTION: #1\n")
    } else if future2 > future3 {
        fmt.Printf("✓ BEST OPTION: #2\n")
    } else {
        fmt.Printf("✓ BEST OPTION: #3\n")
    }
    fmt.Printf("═════════════════════════════════════════════════════\n")
}
```

### 3. Bill of Exchange Discount

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Scenario: A company has a $50,000 bill of exchange
    // due in 6 months. A bank discounts it at 10% annual anticipated rate.
    // How much does the company receive?

    rate, _ := compositeinterest.NewRateInterest(
        0.10,
        compositeinterest.Annually,
        compositeinterest.RateAnticipateEffectyAnnually,
    )

    // 6 months is 0.5 years
    period, _ := compositeinterest.NewPeriod(0.5, compositeinterest.Annually)

    ci, _ := compositeinterest.New(0, 50000, rate, period)

    received, _ := ci.Present()
    discount := 50000 - received

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("    BILL OF EXCHANGE DISCOUNT\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Face value: $50,000.00\n")
    fmt.Printf("Discount rate: 10%% annual (anticipated)\n")
    fmt.Printf("Term: 6 months\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Amount received: $%.2f\n", received)
    fmt.Printf("Discount: $%.2f\n", discount)
    fmt.Printf("═════════════════════════════════════════\n")
}
```

---

## Advanced Examples

### 1. Calculate Equivalent Rates for All Frequencies

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/finance/compositeinterest"
)

func main() {
    // Find equivalent rates for all frequencies

    baseRate, _ := compositeinterest.NewRateInterest(
        0.12,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyNominal,
    )

    frequencies := []struct {
        name  string
        freq  compositeinterest.CompoundingFrequency
    }{
        {"Daily", compositeinterest.Daily},
        {"Monthly", compositeinterest.Monthly},
        {"Bimonthly", compositeinterest.Bimonthly},
        {"Quarterly (4x)", compositeinterest.QuarterlyOne},
        {"Quarterly (3x)", compositeinterest.QuarterlyTwo},
        {"Semi-annually", compositeinterest.SemiAnnually},
        {"Annually", compositeinterest.Annually},
    }

    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("  EQUIVALENT NOMINAL RATES TO 12%% MONTHLY\n")
    fmt.Printf("═════════════════════════════════════════════════════\n")

    for _, f := range frequencies {
        equivalent, _ := baseRate.RateNominalToNominal(f.freq)
        fmt.Printf("%-20s: %.4f (%.2f%%)\n", f.name, equivalent, equivalent*100)
    }

    fmt.Printf("═════════════════════════════════════════════════════\n")
}
```

---

## Important Notes

1. **Always handle errors:** The simplified examples omit `if err != nil`, but in real code always check for errors.

2. **Decimal precision:** For real financial calculations, consider using arbitrary precision libraries.

3. **Fractional periods:** Some calculations may result in fractional periods (e.g., 8.11 months).

4. **Input validation:** Always validate that rates are positive and periods are valid.

5. **Unit consistency:** Ensure the period units match the compounding frequency.
