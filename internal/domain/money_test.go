package domain_test

import (
	"testing"

	"github.com/StevenYu0124/kinfolio/internal/domain"
	"github.com/shopspring/decimal"
)

func TestNewMoney_ValidInput(t *testing.T) {
	m, err := domain.NewMoney("10.5", domain.SGD)
	if err != nil {
		t.Fatalf("NewMoney() error = %v, want nil", err)
	}

	wantAmount := decimal.RequireFromString("10.5")
	if !m.Amount().Equal(wantAmount) {
		t.Errorf("Amount() = %v, want %v", m.Amount(), wantAmount)
	}

	wantCurrency := domain.SGD
	if m.Currency() != wantCurrency {
		t.Errorf("Currency() = %q, want %q", m.Currency(), wantCurrency)
	}
}

func TestNewMoney_InvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		currency domain.Currency
	}{
		{"non-numeric amount", "abc", domain.SGD},
		{"empty amount", "", domain.SGD},
		{"unknown currency", "12.34", "ABC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.NewMoney(tt.amount, tt.currency)
			if err == nil {
				t.Errorf("NewMoney(%q, %q) error = nil, want an error", tt.amount, tt.currency)
			}
		})
	}
}
