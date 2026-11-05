package domain

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

var ErrInvalidAmount = errors.New("invalid amount")

type Money struct {
	amount   decimal.Decimal
	currency Currency
}

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

func (m Money) Amount() decimal.Decimal {
	return m.amount
}

func (m Money) Currency() Currency {
	return m.currency
}
