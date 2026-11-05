# Kinfolio — Decisions

## D-01 Language: Go
- Chosen to learn a language that thinks differently from C#: simple, explicit, composition over inheritance.
- Alternatives considered: Kotlin (best PDF support, closest to C#), F# (stays on .NET), Rust, Java.

## D-02 PDF text extraction: `pdftotext -layout` (Poppler)
- Native Go PDF libraries are either basic (`ledongthuc/pdf`), not for extraction (`pdfcpu`),
  commercial (`unipdf`) or AGPL + cgo (`go-fitz`).
- `pdftotext -layout` keeps column alignment, so tables can be parsed line by line.
- Called via `os/exec`. Requires `brew install poppler`.
- Hidden behind a `TextExtractor` interface so it can be replaced.
- **To verify:** run against DBS, Moomoo and Tiger sample PDFs before building the parsers.
- Prefer CSV over PDF where a source offers both (Moomoo, Tiger).

## D-03 Money: `shopspring/decimal`
- Never use `float64` for money.
- Amounts always carry a currency (`Money{Amount, Currency}`); mixing currencies without
  an explicit FX conversion is an error.

## D-04 Modelling without classes
- Importers: one `Importer` interface; Moomoo, Tiger and DBS each implement it.
- Transaction kinds (trade, dividend, fee, interest, transfer, option event): an interface
  plus a type switch. Go does not check that a switch covers every kind, so add tests
  to make sure each kind is handled.

## D-05 Application type: CLI
- A single command-line binary, `kinfolio`.

## D-06 Project layout (proposed)

```
cmd/kinfolio/          main package (CLI entry point)
internal/domain/       Money, Currency, Owner, Transaction, Lot, Position
internal/importer/     Importer interface, TextExtractor
  moomoo/ tiger/ dbs/  one package per source
internal/validate/     statement integrity checks (UC-02), transfer reconciliation (UC-03)
internal/ledger/       per-owner lots and cash (UC-06..08)
internal/report/       statement generation (UC-09, UC-10)
docs/
testdata/              anonymised sample statements
```

`internal/` is enforced by the Go compiler: these packages cannot be imported from outside the module.

## D-07 Money zero value: not valid, not guarded
- `Money{}` (currency `""`) is not a valid `Money`; always create one with `NewMoney`. This is
  stated in the doc comment, not enforced in code (YAGNI).
- Go can't hide a struct's zero value (like `default(T)` for a C# struct). `Money{}` can still
  appear through ignored errors, uninitialised struct fields or `var m Money`.
- Current behaviour is accepted: `Money{}.Add(x)` with a real currency returns
  `ErrCurrencyMismatch`, so it fails loudly. `Money{}.Add(Money{})` returns `Money{}` without
  error; this is harmless and deliberately **not** pinned by a test.
- Revisit when summing is needed (e.g. `var total Money` in a loop): add a constructor such as
  `Zero(currency)` or a `Sum(currency, ...Money)` function, driven by that real use.
