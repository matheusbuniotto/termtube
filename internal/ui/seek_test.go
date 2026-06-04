package ui

import "testing"

func TestParseSeek(t *testing.T) {
	tests := []struct {
		in   string
		want float64
	}{
		{"1:30", 90},
		{"1 30", 90},
		{"90", 90},
		{":jump 2:05", 125},
		{"jump 0 45", 45},
		{"1m30s", 90},
	}
	for _, tc := range tests {
		got, err := ParseSeek(tc.in)
		if err != nil {
			t.Fatalf("%q: %v", tc.in, err)
		}
		if got != tc.want {
			t.Fatalf("%q: got %v want %v", tc.in, got, tc.want)
		}
	}
}
