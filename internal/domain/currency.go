package domain

import (
	"errors"
	"fmt"
)

type Currency string

const (
	SGD Currency = "SGD"
	HKD Currency = "HKD"
)

var ErrUnknownCurrency = errors.New("unknown currency")

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
