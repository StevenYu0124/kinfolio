package domain

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var (
	// ErrInvalidAmount is returned by NewMoney when the amount is not a valid decimal number.
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrCurrencyMismatch is returned when an operation combines Money in different currencies.
	// Converting between currencies always has to be done explicitly.
	ErrCurrencyMismatch = errors.New("currency mismatch")
)

// Money is an exact decimal amount in a single currency. It is an immutable value:
// operations return a new Money and never change the receiver.
//
// The zero value is not a valid Money; always create one with NewMoney.
type Money struct {
	amount   decimal.Decimal
	currency Currency
}

// NewMoney returns a Money for the decimal amount, such as "10.50" or "-3", in the given currency.
// It returns an error wrapping ErrUnknownCurrency if the currency is not supported, or
// ErrInvalidAmount if the amount is not a decimal number.
func NewMoney(amount string, currency Currency) (Money, error) {
	if err := currency.validate(); err != nil {
		return Money{}, err
	}

	a, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("%w %q: %v", ErrInvalidAmount, amount, err)
	}

	return Money{amount: a, currency: currency}, nil
}

// Amount returns the decimal amount.
func (m Money) Amount() decimal.Decimal {
	return m.amount
}

// Currency returns the currency of the amount.
func (m Money) Currency() Currency {
	return m.currency
}

// Add returns the sum of m and other. Both must be in the same currency;
// otherwise Add returns an error wrapping ErrCurrencyMismatch.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, fmt.Errorf("%w: %s and %s", ErrCurrencyMismatch, m.currency, other.currency)
	}
	return Money{amount: m.amount.Add(other.amount), currency: m.currency}, nil
}

// String returns the amount and currency code, such as "10.5 SGD". It is meant for logs
// and debugging, not for display: the amount is not padded to a fixed number of decimals.
func (m Money) String() string {
	return fmt.Sprintf("%s %s", m.amount, m.currency)
}
