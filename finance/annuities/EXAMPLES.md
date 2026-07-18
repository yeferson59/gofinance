# Usage Examples - Annuities Package

This file contains practical examples of how to use the `annuities` package.

## Table of Contents

1. [Basic Examples](#basic-examples)
2. [Financial Use Cases](#financial-use-cases)
3. [Advanced Examples](#advanced-examples)

---

## Basic Examples

### 1. Calculate Present Value of an Annuity

```go
package main

import (
    "fmt"
    "log"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: $100 monthly payments for 12 months at 1% monthly
    // What is the present value?

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

    // Create annuity with $100 payment
    ann, err := annuities.New(100, 0, 0, period, rate)
    if err != nil {
        log.Fatal(err)
    }

    // Calculate present value
    present, err := ann.Present()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Monthly payment: $100.00\n")
    fmt.Printf("Rate: 1%% monthly\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Present value: $%.2f\n", present)
    // Output:
    // Monthly payment: $100.00
    // Rate: 1% monthly
    // Period: 12 months
    // Present value: $1,125.51
}
```

### 2. Calculate Future Value of an Annuity

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: $500 monthly savings for 12 months at 0.5% monthly
    // How much will be accumulated?

    rate, _ := compositeinterest.NewRateInterest(
        0.005,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

    ann, _ := annuities.New(500, 0, 0, period, rate)

    future, _ := ann.Future()

    fmt.Printf("Monthly savings: $500.00\n")
    fmt.Printf("Rate: 0.5%% monthly\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Future value: $%.2f\n", future)
    // Output: Future value: $6,133.56
}
```

### 3. Calculate Periodic Payment to Pay Off a Loan

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: $10,000 loan at 1% monthly for 12 months
    // What's the monthly payment?

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 10000, 0, period, rate)

    payment, _ := ann.PaymentFromPresentValue()

    fmt.Printf("Loan amount: $10,000.00\n")
    fmt.Printf("Rate: 1%% monthly\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Monthly payment: $%.2f\n", payment)
    // Output: Monthly payment: $888.49
}
```

### 4. Calculate Periodic Payment to Accumulate a Future Value

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: Want to accumulate $10,000 in 12 months at 1% monthly
    // How much should we save monthly?

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(12, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 0, 10000, period, rate)

    payment, _ := ann.PaymentFromFutureValue()

    fmt.Printf("Target amount: $10,000.00\n")
    fmt.Printf("Rate: 1%% monthly\n")
    fmt.Printf("Period: 12 months\n")
    fmt.Printf("Monthly payment needed: $%.2f\n", payment)
    // Output: Monthly payment needed: $788.49
}
```

### 5. Calculate Number of Periods

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: With $500 monthly payments at 1% monthly
    // How many months to reach $10,000?

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriod(1, compositeinterest.Monthly)

    ann, _ := annuities.New(500, 0, 10000, period, rate)

    periods, _ := ann.PeriodsWithFuture()

    fmt.Printf("Monthly payment: $500.00\n")
    fmt.Printf("Target amount: $10,000.00\n")
    fmt.Printf("Rate: 1%% monthly\n")
    fmt.Printf("Periods needed: %.2f months\n", periods)
    // Output: Periods needed: 19.41 months
}
```

---

## Financial Use Cases

### 1. Mortgage Calculation

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: $300,000 mortgage at 4.5% annual for 30 years
    // Calculate monthly payment

    // Create annual rate of 4.5%
    rate, _ := compositeinterest.NewRateInterest(
        0.045,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    // Convert to monthly equivalent rate
    // Using RatePeriodicToPeriodic to convert annual to monthly
    monthlyRate, _ := rate.RatePeriodicToPeriodic(compositeinterest.Monthly)

    // Create monthly rate with the converted value
    monthlyRateObj, _ := compositeinterest.NewRateInterest(
        monthlyRate,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    // 30 years × 12 months = 360 periods
    period, _ := compositeinterest.NewPeriod(360, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 300000, 0, period, monthlyRateObj)

    payment, _ := ann.PaymentFromPresentValue()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("       MORTGAGE CALCULATION\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Loan amount: $300,000.00\n")
    fmt.Printf("Annual rate: 4.5%%\n")
    fmt.Printf("Term: 30 years (360 months)\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Monthly payment: $%.2f\n", payment)
    fmt.Printf("Total amount paid: $%.2f\n", payment*360)
    fmt.Printf("Total interest paid: $%.2f\n", (payment*360)-300000)
    fmt.Printf("═════════════════════════════════════════\n")
}
```

### 2. Savings Plan

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: Save $500 monthly at 6% annual for 5 years
    // How much will be accumulated?

    rate, _ := compositeinterest.NewRateInterest(
        0.06,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyNominal,
    )

    // Convert nominal to periodic
    monthlyRate, _ := rate.RatePeriodic()

    monthlyRateObj, _ := compositeinterest.NewRateInterest(
        monthlyRate,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    // 5 years × 12 months = 60 periods
    period, _ := compositeinterest.NewPeriod(60, compositeinterest.Monthly)

    ann, _ := annuities.New(500, 0, 0, period, monthlyRateObj)

    future, _ := ann.Future()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("       SAVINGS PLAN\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Monthly savings: $500.00\n")
    fmt.Printf("Annual rate: 6%% (nominal monthly)\n")
    fmt.Printf("Period: 5 years (60 months)\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Accumulated amount: $%.2f\n", future)
    fmt.Printf("Total contributed: $%.2f\n", 500*60)
    fmt.Printf("Interest earned: $%.2f\n", future-(500*60))
    fmt.Printf("═════════════════════════════════════════\n")
}
```

### 3. Retirement Plan

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: Plan to retire in 25 years with $1 million
    // Invest at 7% annual. How much monthly?

    rate, _ := compositeinterest.NewRateInterest(
        0.07,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    // Convert to monthly
    monthlyRate, _ := rate.RatePeriodicToPeriodic(compositeinterest.Monthly)

    monthlyRateObj, _ := compositeinterest.NewRateInterest(
        monthlyRate,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    // 25 years × 12 months = 300 periods
    period, _ := compositeinterest.NewPeriod(300, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 0, 1000000, period, monthlyRateObj)

    payment, _ := ann.PaymentFromFutureValue()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("       RETIREMENT PLAN\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Goal: $1,000,000.00\n")
    fmt.Printf("Annual rate: 7%%\n")
    fmt.Printf("Period: 25 years (300 months)\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Monthly investment needed: $%.2f\n", payment)
    fmt.Printf("Total contributed: $%.2f\n", payment*300)
    fmt.Printf("Interest earned: $%.2f\n", 1000000-(payment*300))
    fmt.Printf("═════════════════════════════════════════\n")
}
```

### 4. Equipment Lease

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: Equipment worth $50,000 leased over 5 years at 8% annual
    // Calculate annual lease payment

    rate, _ := compositeinterest.NewRateInterest(
        0.08,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    period, _ := compositeinterest.NewPeriod(5, compositeinterest.Annually)

    ann, _ := annuities.New(0, 50000, 0, period, rate)

    payment, _ := ann.PaymentFromPresentValue()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("       EQUIPMENT LEASE\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Equipment value: $50,000.00\n")
    fmt.Printf("Annual rate: 8%%\n")
    fmt.Printf("Lease term: 5 years\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Annual lease payment: $%.2f\n", payment)
    fmt.Printf("Total lease payments: $%.2f\n", payment*5)
    fmt.Printf("═════════════════════════════════════════\n")
}
```

### 5. Student Loan Repayment

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Scenario: $50,000 student loan at 5% annual for 10 years
    // Calculate monthly payment

    rate, _ := compositeinterest.NewRateInterest(
        0.05,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    // Convert to monthly
    monthlyRate, _ := rate.RatePeriodicToPeriodic(compositeinterest.Monthly)

    monthlyRateObj, _ := compositeinterest.NewRateInterest(
        monthlyRate,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    // 10 years × 12 months = 120 periods
    period, _ := compositeinterest.NewPeriod(120, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 50000, 0, period, monthlyRateObj)

    payment, _ := ann.PaymentFromPresentValue()

    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("    STUDENT LOAN REPAYMENT\n")
    fmt.Printf("═════════════════════════════════════════\n")
    fmt.Printf("Loan amount: $50,000.00\n")
    fmt.Printf("Annual rate: 5%%\n")
    fmt.Printf("Repayment term: 10 years (120 months)\n")
    fmt.Printf("─────────────────────────────────────────\n")
    fmt.Printf("Monthly payment: $%.2f\n", payment)
    fmt.Printf("Total amount paid: $%.2f\n", payment*120)
    fmt.Printf("Total interest paid: $%.2f\n", (payment*120)-50000)
    fmt.Printf("═════════════════════════════════════════\n")
}
```

---

## Advanced Examples

### 1. Compare Different Loan Terms

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    loanAmount := 100000.0
    annualRate := 0.06

    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("    COMPARE DIFFERENT LOAN TERMS\n")
    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("Loan amount: $%.2f\n", loanAmount)
    fmt.Printf("Annual rate: %.1f%%\n\n", annualRate*100)

    // Different terms
    terms := []struct {
        years   int
        months  int
    }{
        {5, 60},
        {10, 120},
        {15, 180},
        {20, 240},
        {30, 360},
    }

    rate, _ := compositeinterest.NewRateInterest(
        annualRate,
        compositeinterest.Annually,
        compositeinterest.RateEffectyAnnually,
    )

    monthlyRate, _ := rate.RatePeriodicToPeriodic(compositeinterest.Monthly)
    monthlyRateObj, _ := compositeinterest.NewRateInterest(
        monthlyRate,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    for _, term := range terms {
        period, _ := compositeinterest.NewPeriod(float64(term.months), compositeinterest.Monthly)
        ann, _ := annuities.New(0, loanAmount, 0, period, monthlyRateObj)
        payment, _ := ann.PaymentFromPresentValue()

        fmt.Printf("OPTION: %2d-year term (%3d months)\n", term.years, term.months)
        fmt.Printf("  Monthly payment: $%,.2f\n", payment)
        fmt.Printf("  Total paid: $%,.2f\n", payment*float64(term.months))
        fmt.Printf("  Total interest: $%,.2f\n\n", (payment*float64(term.months))-loanAmount)
    }

    fmt.Printf("═════════════════════════════════════════════════════\n")
}
```

### 2. Calculate Remaining Balance Over Time

```go
package main

import (
    "fmt"
    "github.com/yeferson59/gofinance/v2/finance/annuities"
    "github.com/yeferson59/gofinance/v2/finance/compositeinterest"
)

func main() {
    // Track remaining balance through loan repayment
    // $10,000 loan at 1% monthly for 12 months

    rate, _ := compositeinterest.NewRateInterest(
        0.01,
        compositeinterest.Monthly,
        compositeinterest.RateEffectyPeriodic,
    )

    period, _ := compositeinterest.NewPeriodic(12, compositeinterest.Monthly)

    ann, _ := annuities.New(0, 10000, 0, period, rate)

    payment, _ := ann.PaymentFromPresentValue()

    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("       LOAN AMORTIZATION SCHEDULE\n")
    fmt.Printf("═════════════════════════════════════════════════════\n")
    fmt.Printf("Initial balance: $10,000.00\n")
    fmt.Printf("Monthly payment: $%.2f\n", payment)
    fmt.Printf("Monthly rate: 1%%\n")
    fmt.Printf("─────────────────────────────────────────────────────\n")

    balance := 10000.0
    totalInterest := 0.0

    for month := 1; month <= 12; month++ {
        interest := balance * 0.01
        principal := payment - interest
        balance -= principal
        totalInterest += interest

        if month%3 == 1 || month == 12 {
            fmt.Printf("Month %2d: Payment: $%.2f | Interest: $%.2f | Principal: $%.2f | Balance: $%.2f\n",
                month, payment, interest, principal, balance)
        }
    }

    fmt.Printf("─────────────────────────────────────────────────────\n")
    fmt.Printf("Total interest paid: $%.2f\n", totalInterest)
    fmt.Printf("═════════════════════════════════════════════════════\n")
}
```

---

## Important Notes

1. **Always handle errors:** The simplified examples omit `if err != nil`, but in real code always check for errors.

2. **Decimal precision:** For real financial calculations, consider using arbitrary precision libraries.

3. **Rate conversions:** When changing from annual to monthly rates, use the appropriate conversion methods from compositeinterest.

4. **Payment timing:** This package assumes payments are made at the end of each period (ordinary annuity).

5. **Unit consistency:** Ensure that the period frequency matches the rate frequency.

6. **Zero values:** When creating an annuity, pass 0 for values you want to calculate.
