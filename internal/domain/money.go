package domain

import (
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
)

type Money struct {
	amount   decimal.Decimal
	currency string
}

func NewMoney(amount, currency string) (Money, error) {
	if currency == "" {
		return Money{}, errors.New("currency must not be empty")
	}

	a, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, fmt.Errorf("invalid amount %q: %w", amount, err)
	}

	return Money{amount: a, currency: currency}, nil
}

func (m Money) Amount() decimal.Decimal {
	return m.amount
}

func (m Money) Currency() string {
	return m.currency
}
