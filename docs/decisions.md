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

## D-04 Importers: one `Importer` interface, one implementation per source (proposed)
- Moomoo, Tiger and DBS each implement `Importer`, in their own package (see D-06).
- Adding a source means adding a package; the rest of the app only sees the interface.

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

## D-08 Transaction kinds: interface + type switch (proposed)
- Kinds: trade, dividend, fee, interest, transfer, option event. Each is its own struct; all
  satisfy a `Transaction` interface. Code that handles them (e.g. `ledger`) uses a type switch,
  like C# pattern matching (`switch (tx) { case Trade t: ... }`).
- Why: keeps per-kind logic in the package that owns it (`ledger`, `report`), so `domain` holds
  data only and does not depend on them.
- Alternatives considered:
  - One struct with a `Kind` field and optional fields: simple and close to CSV rows, but most
    fields are empty for most kinds and nothing stops invalid combinations.
  - Interface with behaviour (e.g. `Apply(l *Ledger)` per kind, like C# virtual methods): the
    compiler forces every kind to implement it, but ledger logic moves into `domain` and
    `domain` depends on `ledger`.
- Cost: Go does not check that a switch covers every kind, so add tests to make sure each
  kind is handled.

## D-09 Currency: `type Currency string` with typed constants
- `type Currency string` with constants `SGD`, `HKD`; codes read and print as themselves.
- Alternatives considered: `iota` integer enum (needs a lookup table to/from text), struct-based
  enum (unexported field, closest to a closed C# enum class, but more code).
- Untyped literals (`"XYZ"`, `""`) still convert to `Currency`, so values are validated at runtime
  (`ParseCurrency` for text input, `NewMoney` for constructed values).

## D-10 Money is an immutable value type
- Methods use value receivers and return a new `Money` (like a C# `readonly struct`).
- Methods that combine valid `Money` values build the result directly, not via `NewMoney`: the
  inputs are already valid.
- `String()` (`"12 SGD"`) is for logs, errors and test messages; amount as `decimal` prints it.
  Display formatting (padding, thousands separators) belongs in `report`.

## D-11 Errors: sentinel errors, matched with `errors.Is`
- Callers distinguish errors by **sentinel identity** (`errors.Is`), never by message text. Tests
  don't pin message text.
- Sentinels: `var ErrXxx = errors.New("...")`, lowercase, no trailing punctuation, name the *kind*.
  "unknown" = well-formed but unsupported; "invalid" = malformed.
- Message shape: `"%w %q"` → `<sentinel> "<input>"`; add `": %v"` for a lower-level cause
  (`invalid amount "abc": can't convert abc to decimal`).
- `%w` only for errors callers may depend on (own sentinels, own layers); `%v` for third-party
  implementation details (e.g. `decimal`'s errors) so they don't become part of the API.
- Add a prefix only when it adds context the inner error lacks; otherwise return `err` as is.
- Custom error types with `errors.As` (≈ C# `catch (SomeException ex)`) only when a caller needs
  data from the error; sentinels (≈ a shared well-known instance) are enough so far.
- Supersedes an earlier idea: "each layer prefixes where it is, always `%w`".
