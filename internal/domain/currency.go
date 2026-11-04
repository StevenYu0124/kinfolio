package domain

import "fmt"

type Currency string

const (
	SGD Currency = "SGD"
	HKD Currency = "HKD"
)

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
		return fmt.Errorf("unknown currency %q", c)
	}
}
