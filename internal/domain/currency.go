package domain

import (
	"errors"
	"fmt"
)

// Currency is an ISO 4217 currency code. Only the codes declared as constants below are
// supported; use ParseCurrency to turn text input into a Currency.
type Currency string

// Supported currencies.
const (
	SGD Currency = "SGD"
	HKD Currency = "HKD"
)

// ErrUnknownCurrency is returned for a well-formed currency code that is not supported.
var ErrUnknownCurrency = errors.New("unknown currency")

// ParseCurrency returns the Currency for the code s, such as "SGD". Codes are case-sensitive.
// It returns an error wrapping ErrUnknownCurrency if s is not a supported currency.
func ParseCurrency(s string) (Currency, error) {
	c := Currency(s)
	if err := c.validate(); err != nil {
		return "", err
	}
	return c, nil
}

func (c Currency) validate() error {
	switch c {
	case SGD, HKD:
		return nil
	default:
		return fmt.Errorf("%w %q", ErrUnknownCurrency, c)
	}
}
