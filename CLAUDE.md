# Kinfolio

Per-owner investment statements for a family, built from Moomoo, Tiger and DBS statements.
See `docs/use-cases.md` and `docs/decisions.md`.

## Working agreement

- This is a learning project: the user is learning Go (background: C#, prefers OOP and clean code).
- **Do not write or edit code, and do not run commands**, unless the user explicitly asks for help
  with a specific task. Explain, review, suggest and point to docs instead.
  - Exception: `go test` (any package or flags, e.g. `go test ./...`, `-run`, `-v`, `-count=1`)
    may be run without asking, e.g. to check a red/green state while reviewing.
  - Exception: Claude may `git commit` without asking (e.g. at a green, or after doc edits).
    Commit messages follow the existing style: short, imperative (`Add Money.Add with ...`).
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
- Claude may run `go test` without asking (see the exception under "Working agreement").

## Conventions

Agreed while building `internal/domain`. Progress lives in the git log; this section keeps the
*why*/*how* that commit messages don't.

Design:
- `Currency` is `type Currency string` with typed constants (`SGD`, `HKD`), chosen over `iota` and
  struct-based enums. `ParseCurrency` is for text input; unexported `validate()` holds the single list
  of supported codes. Untyped literals (`"XYZ"`) still convert to `Currency`, so runtime validation
  stays. Prefer `domain.SGD` over `"SGD"` in code and tests.
- Methods on `Money` use value receivers and build the struct directly (no round trip through
  `NewMoney`). `String()` → `"12 SGD"`, amount as `decimal` prints it; display formatting belongs in
  `report`. Inside `String()`, format the fields, never the receiver (infinite recursion).
- API (signature) changes are done as one compiler-guided refactor, green → green; use parallel change
  (expand → migrate → contract) once there are many callers.

Errors:
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

Tests:
- Rules are tested once, where they live (currency rules in `currency_test.go`); callers get one row
  that only proves delegation. Don't add tests for language behaviour.
- Error tests assert with `errors.Is` against a `wantErr` column holding a real sentinel on every row
  (never `nil` = "any error").
- Table-driven tests have a `name` column used in `t.Run(tt.name, ...)`; use `got`/`want`.
- Method tests are named `TestType_Method[_Case]` (e.g. `TestMoney_Add_CurrencyMismatch`).
- Test setup uses the helper `mustMoney(t, amount, currency)` with `t.Helper()`; never drop setup errors.
- Failure messages show the call's arguments, taken from locals/table fields so they can't drift from
  the call. Print `Money` with `%v` (`12 SGD.Add(23 SGD) ...`); in `TestMoney_String` itself, print the
  raw inputs instead (formatting the value would be circular).
- A test that is green from the start is fine when it documents a decision (e.g. "trailing zero dropped").
- Examples don't print third-party error text in `// Output:` (it would pin `decimal`'s messages).
- Run tests with `./` paths (`go test ./internal/domain`, `go test ./...`); `internal/domain` alone is
  treated as a standard-library path. `(cached)` means unchanged since a passing run; `-count=1` forces.

## Next session

Start with:
1. Decide what comes next in the domain (e.g. `Sub`/`Neg`, or `Owner`/`Transaction` per D-06), or a
   summing helper (`Zero(currency)` / `Sum`) if a real use appears (D-07).

Smaller leftovers:
- Commit at each green, *before* starting the next red (still slipping: two cycles piled up again).
- Examples use `a, _ := domain.NewMoney(...)` for brevity; the user may want them to check errors
  (a good first example for the user to write).
