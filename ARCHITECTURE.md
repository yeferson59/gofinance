# Architecture

This document describes the layering rules of the repository, the reasoning
behind the current structure, and the follow-up work that is recommended but
not yet executed.

## Modules

The repository hosts three Go modules:

| Module | Path | External runtime deps |
| ------ | ---- | --------------------- |
| `github.com/yeferson59/gofinance` | repo root | none (testify is test-only) |
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

## Recommended follow-ups (not yet executed)

- **Rename `finance/compositeinterest` → `finance/compoundinterest`.**
  "Composite interest" is a non-standard term; the rename touches README,
  examples, benchmarks and all import paths, so it should ride a major
  version bump.
- **Deprecate `Money.Mul(Money)` and `Money.Div(Money)`.** Multiplying two
  monetary amounts is dimensionally meaningless (money²); the surviving
  cases are ratios, better expressed via `ToDecimal()` or
  `MulDecimal`/`DivDecimal`. Same for `Money.Add`/`Sub` silently ignoring a
  currency mismatch — consider making the checked (`Safe*`) behavior the
  default in a major release.
- **Unify time/frequency vocabulary.** `simpleinterest.Periods`,
  `compositeinterest.CompoundingFrequency` and `daycount.Convention` model
  overlapping concepts with different types; a shared `finance/term` (or
  similar) package would let them interoperate.
- **Version note.** Removing `Decimal.ToMoney` and moving the `charts`
  import path are breaking changes relative to v1.4.2. Under strict semver
  they call for a v2 (`/v2` module path); alternatively document them
  loudly in the next minor release notes, since the rest of the surface is
  source-compatible through the alias and forwarders.
