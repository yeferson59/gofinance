# 🧭 GoFinance — Feature Candidates

A curated, prioritized list of features that would be **interesting and valuable**
additions to GoFinance. These are *candidates*, not commitments — a roadmap to discuss,
refine, and pull from.

Each entry notes what existing building blocks it can reuse and where it fits in the
current architecture, so nothing gets reinvented.

## How to read this list

- **Priority** — how much value it adds relative to effort:
  - 🟢 **High** — high-value, low-friction, or unblocks other features
  - 🟡 **Medium** — clearly useful, moderate scope
  - ⚪ **Low** — nice-to-have or niche
- **Effort** — rough size: **S** (a file or two), **M** (a small package), **L** (a full package with edge cases)
- **✅** marks a feature that is already implemented and shipped.

### Guiding principles for every feature

- Build on `decimal.Decimal` / `money.Money` — **never** an internal `float64`. This is
  the library's core differentiator (19 exact fractional digits).
- Mirror the existing **fluent builder** style (`compositeinterest.NewComposite()`,
  `annuities.NewAnnuity()`, `compositeinterest.NewRateConversion()`).
- Follow the **`Try*` (returns `error`) / `Must*` (panics)** convention already used
  across `decimal` and the finance packages.
- Keep core packages **zero external dependencies** and **immutable / goroutine-safe**.

---

## 1. 📈 Investment instruments

Proposed home: a new `finance/investment` package (plus `finance/bonds`).

| # | Feature | What it does | Reuses | Priority | Effort |
|---|---------|--------------|--------|----------|--------|
| 1.1 | ✅ **NPV / VAN** | Net present value of a series of periodic cash flows at a discount rate | `decimal.Pow`, `decimal.Div`, discounting logic from `compositeinterest` | 🟢 High | S |
| 1.2 | ✅ **IRR / TIR** | Internal rate of return — solve for the rate where NPV = 0 | Newton–Raphson / bisection over 1.1, the root-solver pattern in `finance/*/root.go` | 🟢 High | M |
| 1.3 | ✅ **XNPV / XIRR** | NPV/IRR for *irregular, date-based* cash flows | 1.1, 1.2 + day-count conventions (§3.1) | 🟡 Medium | M |
| 1.4 | ✅ **Perpetuities** | Present value of a level and a *growing* perpetuity (Gordon model) | `decimal.Div`, `money.Money` | 🟡 Medium | S |
| 1.5 | ✅ **Bond valuation** | Clean price, yield-to-maturity, accrued interest, Macaulay & modified duration, convexity | `decimal` math, discounting, day-count (§3.1) → `finance/bonds` | 🟡 Medium | L |

> **Why NPV/IRR first:** they are the cornerstone of capital budgeting and the single
> most-requested finance primitives missing today. Both are pure decimal math over a
> cash-flow slice — a natural, high-value fit for the existing engine.

---

## 2. 🏦 Loans & accounting

Proposed home: extend `finance/annuities`; add `finance/depreciation`.

| # | Feature | What it does | Reuses | Priority | Effort |
|---|---------|--------------|--------|----------|--------|
| 3.1 | ✅ **Day-count conventions** | 30/360, Actual/365, Actual/360, Actual/Actual — precise interest by calendar dates | shared utility → **unblocks** XIRR (1.3), bonds (1.5), date-based interest | 🟢 High | M |
| 3.2 | **APR / effective annual rate** | Effective APR from a nominal rate plus fees/points | `compositeinterest` rate conversions (`NewRateConversion()`) | 🟡 Medium | S |
| 3.3 | **Extra payments / early payoff** | Interest saved and term shortened when overpaying a loan | extends `annuities.BuildSchedule` (`finance/annuities/utils.go`) | 🟡 Medium | M |
| 3.4 | **Refinance comparator** | Break-even point between two loan offers | 3.2, 3.3, NPV (1.1) | ⚪ Low | S |
| 3.5 | ✅ **Depreciation** | Straight-line, declining balance (single/double), sum-of-years'-digits, MACRS | `decimal` arithmetic → `finance/depreciation` | 🟡 Medium | M |

> **Why day-count first in this section:** it's a small shared utility that several
> other features (XIRR, bonds, accrued interest) depend on. Building it once removes a
> blocker from three places.

---

## 3. 📊 Metrics & returns

Proposed home: a new `finance/returns` package.

| # | Feature | What it does | Reuses | Priority | Effort |
|---|---------|--------------|--------|----------|--------|
| 4.1 | ✅ **CAGR / annualized return** | Compound annual growth rate from start value, end value, and horizon | `decimal.Pow` (fractional exponent already supported) | 🟢 High | S |
| 4.2 | ✅ **ROI / holding-period return** | Simple return over a period, with/without income | `money.Money`, `decimal.Div` | 🟢 High | S |
| 4.3 | ✅ **Inflation adjustment** | Real vs nominal values; future/present purchasing power | `decimal.Pow` | 🟡 Medium | S |
| 4.4 | **TWR vs MWR** | Time-weighted vs money-weighted returns for a portfolio | 4.1, IRR (1.2) | ⚪ Low | M |
| 4.5 | **Volatility & Sharpe ratio** | Std. deviation of returns and risk-adjusted return | `decimal.Sqrt` (already available, correctly rounded) | ⚪ Low | M |

> **Why CAGR/ROI first:** tiny, universally useful, and they showcase the fractional
> `decimal.Pow` and `decimal.Sqrt` kernels the library already ships.

---

## 4. 🛠️ Tooling & developer experience

| # | Feature | What it does | Reuses | Priority | Effort |
|---|---------|--------------|--------|----------|--------|
| 5.1 | ✅ **General TVM solver** | Solve any one of N, I/Y, PV, PMT, FV given the other four (HP-12C style) | unifies logic already split across `compositeinterest` and `annuities` | 🟢 High | M |
| 5.2 | **CLI (`cmd/gofinance`)** | Subcommands for NPV/IRR/payment/amortization with tabular and `--json` output | all finance packages; already on the README roadmap | 🟡 Medium | M |
| 5.3 | **Schedule export API** | First-class CSV/JSON export of amortization/investment schedules | `annuities.BuildSchedule`; formalizes today's ad-hoc `amortizacion.csv` | ⚪ Low | S |
| 5.4 | **More runnable examples** | Recipe-style examples per instrument in `examples/` | existing `examples/main.go` | ⚪ Low | S |
| 5.5 | **Web API wrapper** *(optional, non-core)* | Thin HTTP layer over the calculators | all finance packages; on the roadmap | ⚪ Low | L |

> **Why the TVM solver stands out:** the pieces already exist but are scattered. A single
> solver gives users the familiar financial-calculator mental model and removes
> duplicated formula code between the compound-interest and annuities packages.

---

## 5. ⭐ Recommended quick wins (start here)

The best value-to-effort ratio, roughly in order — **all shipped ✅**:

1. ✅ **CAGR / ROI** (4.1, 4.2) — trivial, high-demand, shows off the decimal engine.
2. ✅ **NPV** (1.1) — cornerstone primitive, pure decimal math over a cash-flow slice.
3. ✅ **Day-count conventions** (3.1) — small, and unblocks XIRR + bonds.
4. ✅ **IRR** (1.2) — the natural companion to NPV via a root solver.
5. ✅ **General TVM solver** (5.1) — unifies existing logic behind a familiar API.

The follow-up batch is also shipped: ✅ XNPV/XIRR (1.3), ✅ perpetuities (1.4),
✅ inflation adjustment (4.3), ✅ depreciation (3.5), and ✅ bond valuation (1.5).
Remaining candidates worth picking up next: APR/effective rate (3.2), extra
payments / early payoff (3.3), TWR/MWR (4.4), volatility & Sharpe (4.5), and the
tooling items (CLI, schedule export).

---

## 6. 🧩 Design notes

- **Error style:** expose a `Try*` variant returning `error` for every operation that
  can fail (bad input, no convergence, currency mismatch) plus a `Must*` panicking
  variant, matching `decimal.TryAdd`/`Add` and `annuities.MustPayment`.
- **Convergence:** iterative solvers (IRR, YTM, TVM-rate) should accept an optional
  tolerance/max-iterations and return a typed error (e.g. `ErrNoConvergence`) rather
  than silently returning zero.
- **Money vs Decimal:** cash amounts stay `money.Money` (currency-checked); rates,
  factors, and ratios stay `decimal.Decimal` — consistent with how
  `annuities.BuildSchedule` already separates them.
- **Zero dependencies:** keep new core packages stdlib-only; anything needing
  `go-echarts` or HTTP frameworks lives in separate modules (like `charts`
  today).
- **Changelog:** record each accepted feature under `[Unreleased]` in `CHANGELOG.md`,
  following the existing Keep a Changelog format.

---

*This document is a living roadmap. Items can be promoted, deferred, or dropped as
priorities change — treat priorities and effort estimates as starting points for
discussion, not fixed scope.*
