package detector

import "testing"

func TestIsMixedCaps(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"mixedCaps", true},
		{"MixedCaps", true},
		{"snake_case", false},
		{"Snake_Case", false},
		{"userID", true},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := IsMixedCaps(tt.in); got != tt.want {
				t.Errorf("got %v want %v", got, tt.want)
			}
		})
	}
}
