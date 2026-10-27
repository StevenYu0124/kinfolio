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

Done (2026-10-28): `internal/domain/money.go` with `NewMoney(amount, currency string) (Money, error)`,
unexported fields, `Amount()`/`Currency()` getters; rejects non-numeric amount and empty currency.
Tests: one happy-path test plus a table-driven `TestNewMoney_InvalidInput`.

Next session, start with:
1. **Currency as an enum-like type** (user's choice). Go has no `enum`; discuss options
   (`type Currency string` + typed constants vs `int` + `iota`, plus a `ParseCurrency` factory) and
   drive it with TDD. This likely replaces the empty-currency check and the open "`"  "` currency" question.

Smaller leftovers:
- Commit at each green (habit to keep).
- Doc comments on exported names: user asked Claude to write these once the API settles.
- Then: `Add` (same currency → sum, different currency → error), per D-03.
