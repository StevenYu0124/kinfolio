# Kinfolio

Per-owner investment statements for a family, built from Moomoo, Tiger and DBS statements.
See `docs/use-cases.md` and `docs/decisions.md`.

## Working agreement

- This is a learning project: the user is learning Go (background: C#, prefers OOP and clean code).
- **Do not write or edit code, and do not run commands**, unless the user explicitly asks for help
  with a specific task. Explain, review, suggest and point to docs instead.
- Documentation (files under `docs/`, this file) may be written and edited freely.
- When explaining Go concepts, relate them to C# equivalents where helpful.

## TDD practice

The user is practising test-driven development. TDD is **optional per module**: use it where
the code is pure logic with no I/O, starting with the `Money` type in `internal/domain` (see D-03 in
`docs/decisions.md`). For I/O-heavy code (e.g. DB repositories, importers reading files), tests may
be written after the code, typically as integration tests. Check which approach the user wants
when starting a new module.

- Follow red → green → refactor in small steps: one failing test, the least code to pass it,
  then clean up with the tests green.
- The user writes both the tests and the implementation. Claude acts as a coach:
  - suggest the *next* test case to write (one at a time), not the whole list up front;
  - review tests for clarity and whether they really fail for the right reason;
  - call out when the implementation does more than the current tests require;
  - point out refactoring opportunities only once the bar is green.
- Go testing conventions: tests live in `xxx_test.go` beside the code, use the standard
  `testing` package, prefer table-driven tests with `t.Run` subtests, and run with
  `go test ./...`. Relate to xUnit/NUnit (`[Theory]`/`[InlineData]`) where helpful.
- Running `go test` is still a command: only run it when the user asks.

## Progress & next session

Done (2026-10-28): `internal/domain/money.go` with `NewMoney`, unexported fields, `Amount()`/`Currency()`
getters; rejects non-numeric amount.

Done (2026-11-05): `internal/domain/currency.go`, `type Currency string` with typed constants `SGD`, `HKD`
(option A, chosen over `iota` and struct-based enums). `ParseCurrency(s string) (Currency, error)` for
text input; unexported `(c Currency) validate()` holds the single list of supported codes.
`NewMoney(amount string, currency Currency)` now takes the type and delegates to `validate()`.
Decisions and conventions agreed this session:
- Rules are tested once: currency rules in `currency_test.go`; `money_test.go` has one
  "unknown currency" row that only proves delegation.
- Untyped literals (`"XYZ"`, `""`) still convert to `Currency`, so runtime validation stays. Prefer
  `domain.SGD` over `"SGD"` in code and tests; don't add tests for language behaviour.
- API (signature) changes are done as one compiler-guided refactor, green → green; use parallel change
  (expand → migrate → contract) once there are many callers.

Done (2026-11-05, second session): sentinel errors. `ErrUnknownCurrency` (`currency.go`, returned by
`validate()`) and `ErrInvalidAmount` (`money.go`, returned by `NewMoney`). Tests assert with `errors.Is`;
`TestNewMoney_InvalidInput` has a `wantErr` column (a real sentinel on every row, never `nil` = "any
error"). `NewMoney` returns `validate()`'s error unwrapped. Test failure messages show the call's
arguments, taken from locals/table fields so they can't drift from the call.

Error conventions (supersede the earlier "each layer prefixes where it is, always `%w`"):
- Callers distinguish errors by **sentinel identity** (`errors.Is`), never by message text. Tests don't
  pin message text.
- Sentinels: `var ErrXxx = errors.New("...")`, lowercase, no trailing punctuation, name the *kind*.
  "unknown" = well-formed but unsupported; "invalid" = malformed.
- Message shape: `"%w %q"` → `<sentinel> "<input>"`; add `": %v"` for a lower-level cause
  (`invalid amount "abc": can't convert abc to decimal`).
- `%w` only for errors callers may depend on (own sentinels, own layers); `%v` for third-party
  implementation details (e.g. `decimal`'s errors) so they don't become part of the API.
- Add a prefix only when it adds context the inner error lacks; otherwise return `err` as is.
- `errors.As` / custom error types (≈ C# `catch (SomeException ex)`) only when a caller needs data
  from the error; sentinels (≈ a shared well-known instance) are enough so far.

Next session, start with:
1. **`Add`** (same currency → sum, different currency → error), per D-03, via TDD. Consider what `Add`
   does with the zero value `Money{}` (currency `""`); a sentinel like `ErrCurrencyMismatch` may fit here.

Smaller leftovers:
- Commit at each green (habit to keep; this session piled up several cycles uncommitted).
- Doc comments on exported names (now incl. the two sentinels) plus an `ExampleNewMoney` function:
  user asked Claude to write them once the API settles; it's close to settled now.
