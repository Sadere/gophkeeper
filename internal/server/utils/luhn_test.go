package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckLuhn(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "invalid number",
			number: "4147203059780942",
			want:   false,
		},
		{
			name:   "valid number",
			number: "5062821234567892",
			want:   true,
		},
		{
			name:   "incorrect input number",
			number: "123456abcd",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckLuhn(tt.number)

			if tt.want {
				assert.Truef(t, result, "CheckLuhn() must be true for number %s", tt.number)
			} else {
				assert.Falsef(t, result, "CheckLuhn() must be false for number %s", tt.number)
			}
		})
	}
}
