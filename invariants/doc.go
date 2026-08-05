// Package invariants holds the cross-package tests: properties that must hold
// between packages rather than inside one of them.
//
// They live here, in a package of their own, because each check imports two or
// more of the finance packages and belongs to none. The properties fall into
// three groups:
//
//   - Agreement. Two packages that model the same thing must produce the same
//     number. A loan payment is the same whether it comes from finance/tvm,
//     finance/annuities or finance/loans; a bond's price is the present value
//     of its own cash flows; a gradient series with no gradient is an ordinary
//     annuity.
//
//   - Closure. A schedule must add up: an amortization table ends at a zero
//     balance and its principal sums to the amount borrowed, an allocation
//     sums back to the amount split.
//
//   - Robustness. No function that returns an error may panic on a valid but
//     extreme input. Several such panics were found in the solvers
//     (TESTING_PLAN.md §2.8), all from the same pattern — the panicking
//     decimal helpers used inside error-returning code — so the sweep here
//     looks for it across the whole library rather than package by package.
//
// The package contains no non-test code.
package invariants
