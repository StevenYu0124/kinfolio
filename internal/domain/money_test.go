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
