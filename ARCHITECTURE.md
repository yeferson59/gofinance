# Architecture

This document describes the layering rules of the repository, the reasoning
behind the current structure, and the follow-up work that is recommended but
not yet executed.

## Modules

The repository hosts three Go modules:

| Module | Path | External runtime deps |
| ------ | ---- | --------------------- |
| `github.com/yeferson59/gofinance/v2` | repo root | none (testify is test-only) |
| `github.com/yeferson59/gofinance/charts` | `charts/` | go-echarts |
| `github.com/yeferson59/gofinance/examples` | `examples/` | go-echarts (via `charts`) |

Splitting `charts` and `examples` out keeps the library module dependency-free:
consumers who only need financial math never download a charting stack. The
submodules use `replace` directives pointing at the repo root during local
development; external consumers resolve the published version from the
`require` line.

## Layering

```
decimal        fixed-point numeric engine (stdlib only)
   ▲
money          Currency + Money (an amount bound to a currency)
   ▲
finance/*      domain math: interest, annuities, bonds, TVM, …
   │           (shared time vocabulary lives in finance/term)
   ▲
charts         optional visualization (separate module)
```

The rules, from the bottom up:

1. **`decimal` is the only numeric API.** Rates, factors, ratios, periods and
   any other currency-less quantity are `decimal.Decimal`. No package may wrap
   or mirror the decimal API.
2. **`money` adds exactly one concept: currency.** `Money` is a
   `decimal.Decimal` plus a `Currency`, with currency-aware operations
   (checked add/sub, allocation, conversion, formatting, JSON/SQL codecs).
   `money.FromDecimal(d, currency)` is the bridge from a computed decimal to a
   monetary amount. `money.Decimal` remains only as a deprecated alias of
   `decimal.Decimal` for source compatibility.
3. **`finance/*` computes on `decimal` and carries amounts as `money.Money`.**
   A package that involves no amounts at all (`tvm`, `daycount`) must not
   import `money`. When a rate multiplies an amount, use
   `Money.MulDecimal`/`DivDecimal` — never attach a placeholder currency to a
   rate just to make the types line up.
4. **Domain packages perform no I/O.** Rendering and persistence live at the
   edges: `charts` renders schedules, and CSV export is written against
   `io.Writer` (`annuities.WriteCSVTo`), with the file-on-disk variant kept
   only as a thin convenience wrapper.

## History: what the 2026-07 restructuring changed

The previous structure had four structural problems, all fixed in this pass:

- `money` contained a 335-line hand-written mirror of the whole
  `decimal.Decimal` API. Every finance package depended on `money.Decimal`
  (407 references) and `decimal` was almost unused outside its own package;
  each new decimal method had to be added twice. The mirror is now a type
  alias plus deprecated constructor forwarders.
- Because rates were `money.Decimal`, amount×rate math was expressed as
  `amount.Mul(rate.ToMoney())`, silently attaching USD to a rate.
  `MulDecimal`/`DivDecimal` express the dimensionally correct operation.
- `finance/charts` lived inside the library module, so every consumer pulled
  go-echarts. It is now the separate `charts` module.
- `annuities` wrote CSV files directly to disk, and a generated
  `amortizacion.csv` was committed at the repo root.

## Executed follow-ups (second pass)

- **`finance/compositeinterest` renamed to `finance/compoundinterest`.**
  "Composite interest" is a non-standard term. The types followed suit
  (`CompoundInterest`, `CompoundConfig`, `NewCompound`); no forwarding
  package is kept — this is part of the breaking release described below.
- **`Money.Mul(Money)`, `Money.Div(Money)` and `Money.MustDiv` are
  deprecated.** Multiplying two monetary amounts is dimensionally
  meaningless (money²), and dividing them yields a currency-less ratio.
  Use `MulDecimal`/`DivDecimal` to scale amounts and
  `ToDecimal()` + `decimal.Div` for ratios. `Add`/`Sub` now document that
  they don't currency-check and point to `SafeAdd`/`SafeSub`.
- **Shared time vocabulary in `finance/term`.** `term.Unit`
  (days/weeks/months/years) and `term.Frequency` (daily … annually, with
  `PeriodsPerYear`/`MonthsPerPeriod`) are the one set of types;
  `simpleinterest.Periods` and `compoundinterest.CompoundingFrequency` are
  aliases of them. The confusing `QuarterlyOne`/`QuarterlyTwo` constants
  became `Quarterly` (4/yr) and `FourMonthly` (3/yr), with the old names
  kept as deprecated aliases. Day-count conventions answer a different
  question (date range → year fraction) and deliberately stay in
  `finance/daycount`.

## Executed follow-ups (third pass)

- **The module path is now `github.com/yeferson59/gofinance/v2`.** The
  accumulated breaking changes (`ToMoney` removal, the `charts` module
  split, the `compoundinterest` rename, the `Add`/`Sub` redesign below)
  ride a proper major release. Release-time actions left for the
  maintainer: tag `v2.0.0` on the root module and `charts/v1.0.0` for the
  charts module (its own path carries no `/vN` suffix, so its tags stay in
  the v0/v1 series while requiring the library's `/v2`).
- **`Money.Add`/`Sub` are currency-checked by default.** They now panic on
  a currency mismatch, exactly as they already panicked on overflow —
  mirroring the decimal engine's `Add`/`TryAdd` split. The error-returning
  forms are `TryAdd`/`TrySub` (mismatch → `ErrCurrencyMismatch`);
  `SafeAdd`/`SafeSub` survive as deprecated aliases of them.
