package helper

import (
	"errors"
	"testing"
)

func TestSanitizeRaidName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		shouldErr bool
	}{
		{name: "valid simple", input: "dogFood", want: "dogFood"},
		{name: "valid mixed chars", input: "dog-food_1.test", want: "dog-food_1.test"},
		{name: "trim spaces", input: "  dogFood  ", want: "dogFood"},
		{name: "invalid empty", input: "", shouldErr: true},
		{name: "invalid spaces only", input: "   ", shouldErr: true},
		{name: "invalid internal space", input: "dog Food", shouldErr: true},
		{name: "invalid slash", input: "dog/food", shouldErr: true},
		{name: "invalid leading dash", input: "-dogFood", shouldErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SanitizeRaidName(tt.input)
			if tt.shouldErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tt.input)
				}
				if !errors.Is(err, ErrInvalidRaidName) {
					t.Fatalf("expected ErrInvalidRaidName, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
