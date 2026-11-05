package domain_test

import (
	"errors"
	"testing"

	"github.com/StevenYu0124/kinfolio/internal/domain"
	"github.com/shopspring/decimal"
)

func TestNewMoney_ValidInput(t *testing.T) {
	amount, currency := "10.5", domain.SGD
	m, err := domain.NewMoney(amount, currency)
	if err != nil {
		t.Fatalf("NewMoney(%q, %q) error = %v, want nil", amount, currency, err)
	}

	wantAmount := decimal.RequireFromString(amount)
	if !m.Amount().Equal(wantAmount) {
		t.Errorf("Amount() = %v, want %v", m.Amount(), wantAmount)
	}

	wantCurrency := currency
	if m.Currency() != wantCurrency {
		t.Errorf("Currency() = %q, want %q", m.Currency(), wantCurrency)
	}
}

func TestNewMoney_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		currency domain.Currency
		wantErr  error
	}{
		{"non-numeric amount", "abc", domain.SGD, domain.ErrInvalidAmount},
		{"empty amount", "", domain.SGD, domain.ErrInvalidAmount},
		{"unknown currency", "12.34", "ABC", domain.ErrUnknownCurrency},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewMoney(tt.amount, tt.currency)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewMoney(%q, %q) error = %v, want %v", tt.amount, tt.currency, err, tt.wantErr)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	currency := domain.SGD
	money1 := mustMoney(t, "12", currency)
	money2 := mustMoney(t, "23", currency)
	sum, err := money1.Add(money2)
	if err != nil {
		t.Fatalf("%v.Add(%v) error = %v, want nil", money1.Amount(), money2.Amount(), err)
	}

	wantAmount := decimal.RequireFromString("35")
	if !sum.Amount().Equal(wantAmount) {
		t.Errorf("%v.Add(%v) Amount() = %v, want %v", money1.Amount(), money2.Amount(), sum.Amount(), wantAmount)
	}

	wantCurrency := currency
	if sum.Currency() != wantCurrency {
		t.Errorf("%v.Add(%v) Currency() = %q, want %q", money1.Currency(), money2.Currency(), sum.Currency(), wantCurrency)
	}
}

func TestMoney_Add_CurrencyMismatch(t *testing.T) {
	money1 := mustMoney(t, "12", domain.SGD)
	money2 := mustMoney(t, "23", domain.HKD)
	_, err := money1.Add(money2)
	wantErr := domain.ErrCurrencyMismatch
	if !errors.Is(err, wantErr) {
		t.Errorf("%v.Add(%v) error = %v, want %v", money1.Currency(), money2.Currency(), err, wantErr)
	}
}

func mustMoney(t *testing.T, amount string, c domain.Currency) domain.Money {
	t.Helper()
	m, err := domain.NewMoney(amount, c)
	if err != nil {
		t.Fatalf("NewMoney(%q, %q) error = %v, want nil", amount, c, err)
	}
	return m
}
