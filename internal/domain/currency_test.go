package domain_test

import (
	"errors"
	"testing"

	"github.com/StevenYu0124/kinfolio/internal/domain"
)

func TestParseCurrency_ValidCode(t *testing.T) {
	tests := []struct {
		input string
		want  domain.Currency
	}{
		{"SGD", domain.SGD},
		{"HKD", domain.HKD},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := domain.ParseCurrency(tt.input)
			if err != nil {
				t.Fatalf("ParseCurrency(%q) error = %v, want nil", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseCurrency(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseCurrency_InvalidCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"unknown code", "ABC"},
	}

	wantErr := domain.ErrUnknownCurrency
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := domain.ParseCurrency(tt.input)
			if !errors.Is(err, wantErr) {
				t.Errorf("ParseCurrency(%q) error = %v, want %v", tt.input, err, wantErr)
			}
		})
	}
}
