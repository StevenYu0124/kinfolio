# Kinfolio — Use Cases

Kinfolio turns broker and bank statements into per-owner statements for a family whose
investments are traded through one person's accounts.

## Background

- The manager trades on behalf of family members (e.g. Dad calls, manager places the order).
  All trades live in the manager's own broker accounts, mixed with the manager's own trades.
- Brokers report P&L **excluding fees and margin interest**. Kinfolio reports **net** P&L per owner.
- Sources: **Moomoo** (CSV/PDF), **Tiger** (CSV/PDF), **DBS** (PDF, cash flow only).
- Products: **stocks and options**.

## Actors

- **Manager** — imports statements, tags transactions to owners, generates statements.
- **Owner** — the manager or a family member. Receives statements; does not use the app.

## Principles

- Every transaction, holding and cash movement belongs to exactly **one owner**.
  Untagged items belong to the manager (default owner).
- Imports are all-or-nothing: if validation fails, nothing is saved and a clear report is produced.
- Base currency is chosen at statement generation time.
- Statement period follows the input period.

## A. Input & Validation

### UC-01 Import statement
Detect source, account and period; parse into the normalized model (transactions, positions, cash balances).
- Validate: supported source, known account, readable period, not a duplicate (file hash).

### UC-02 Validate statement integrity
- Cash: opening + movements = closing, per currency.
- Positions: opening + buys − sells ± option events = closing.
- Printed totals = sum of lines.
- Periods are continuous per account (no gaps or overlaps).
- Periods across sources that don't align are reported, not guessed.

### UC-03 Reconcile transfers across sources
Match DBS transfers with broker deposits/withdrawals by amount, date window, currency and FX tolerance.
Unmatched transfers fail validation.

## B. Ownership

### UC-04 Manage owners
Add or edit owners. The manager is the default owner.

### UC-05 Tag to owner
Assign an owner to: trades, cash contributions (e.g. "this DBS transfer is Dad's money"),
fees, and margin interest.
- Validate: every tag refers to an existing item; each item has at most one owner.

### UC-06 Track positions per owner (lots)
- Buys open lots for the owner; sells close that owner's lots (FIFO).
- Options: open, close, expire, assignment and exercise. Resulting stock lots go to the option's owner.
- Validate: an owner cannot sell more than they hold;
  sum of owners' positions = broker positions.

### UC-07 Track cash per owner
Contributions, trade settlements, fees, interest, dividends and option premiums.
- Validate: sum of owners' cash = broker cash.

### UC-08 Attribute income & costs
- Trade fees follow the trade's owner.
- Dividends follow the owner holding the shares on the ex-date.
- Margin interest and account-level fees are tagged explicitly (UC-05); untagged goes to the manager.

## C. Output

### UC-09 Owner statement
For one owner, one period, a chosen base currency:
- Opening/closing positions and cash, trades in the period.
- Realised and unrealised P&L, with price and FX effects shown separately.
- Dividends, option premiums, fees, interest.
- **Net P&L after fees and interest.**
- Format: PDF first.

FX: cost uses the trade-date rate; valuation uses the closing rate. Rates come from statements where
available, otherwise from an external or user-provided rate source.

### UC-10 Consolidated statement
All owners side by side, plus the reconciliation checks from UC-06 and UC-07.

## Deferred

- **Split allocation** — one transaction shared across owners by percentage. Build when first needed.
- **Exception review & manual adjustments** — audit-trailed corrections.
- **Period close / lock.**
