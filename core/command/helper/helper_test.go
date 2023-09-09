package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveAddr(t *testing.T) {
	tests := []struct {
		address    string
		defaultIP  IPBinding
		expectedIP string
		err        bool
	}{
		{"localhost:8080", LocalHostBinding, "127.0.0.1", false},
		{":8080", LocalHostBinding, "127.0.0.1", false},
		{"invalid", LocalHostBinding, "", true},
	}

	for _, test := range tests {
		addr, err := ResolveAddr(test.address, test.defaultIP)
		if test.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expectedIP, addr.IP.String())
		}
	}
}
