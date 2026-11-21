package lib

import (
	"math"
	"testing"
	"time"
)

func TestIsValidHTTPTimeoutValue(t *testing.T) {
	maxValidSeconds := math.MaxInt64 / int64(time.Second)

	tests := []struct {
		name    string
		seconds int64
		want    bool
	}{
		{"negative value is invalid", -1, false},
		// 0 is a valid value: it disables the timeout entirely (no timeout),
		// rather than making requests fail instantly.
		{"zero is valid and disables the timeout", 0, true},
		{"typical default value is valid", 60, true},
		{"largest representable value is valid", maxValidSeconds, true},
		{"value overflowing time.Duration is invalid", maxValidSeconds + 1, false},
		{"math.MaxInt64 is invalid (overflows time.Duration)", math.MaxInt64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidHTTPTimeoutValue(tt.seconds); got != tt.want {
				t.Errorf("IsValidHTTPTimeoutValue(%d) = %v, want %v", tt.seconds, got, tt.want)
			}
		})
	}
}
