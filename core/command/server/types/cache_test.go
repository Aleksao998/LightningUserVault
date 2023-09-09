package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertStringToCacheType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected CacheType
		err      bool
	}{
		{"MEMCACHE", MEMCACHE, false},
		{"memcache", MEMCACHE, false},
		{"MeMcAcHe", MEMCACHE, false},
		{"INVALID", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ConvertStringToCacheType(tt.input)
			if tt.err {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tt.input)
				assert.Equal(t, tt.expected, got, "Expected %s but got %s for input: %s", tt.expected, got, tt.input)
			}
		})
	}
}
