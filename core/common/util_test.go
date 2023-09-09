package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInt64ToBytes tests converting int64 to Bytes
func TestInt64ToBytes(t *testing.T) {
	tests := []struct {
		input  int64
		output []byte
	}{
		{input: 0, output: []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{input: 1, output: []byte{1, 0, 0, 0, 0, 0, 0, 0}},
		{input: -1, output: []byte{255, 255, 255, 255, 255, 255, 255, 255}},
		// Add more test values as needed
	}

	for _, test := range tests {
		result := Int64ToBytes(test.input)
		assert.Equal(t, test.output, result, "For input %d", test.input)
	}
}

// TestInt64ToBytes tests converting Bytes to int64
func TestBytesToInt64(t *testing.T) {
	tests := []struct {
		input  []byte
		output int64
	}{
		{input: []byte{0, 0, 0, 0, 0, 0, 0, 0}, output: 0},
		{input: []byte{1, 0, 0, 0, 0, 0, 0, 0}, output: 1},
		{input: []byte{255, 255, 255, 255, 255, 255, 255, 255}, output: -1},
		// Add more test values as needed
	}

	for _, test := range tests {
		result := BytesToInt64(test.input)
		assert.Equal(t, test.output, result, "For input %v", test.input)
	}
}

// TestConversionIntegrity tests integrity between TestInt64ToBytes TestBytesToInt64
func TestConversionIntegrity(t *testing.T) {
	tests := []int64{
		0,
		1,
		-1,
		1234567890,
		-1234567890,
		// Add more test values as needed
	}

	for _, test := range tests {
		bytes := Int64ToBytes(test)
		result := BytesToInt64(bytes)
		assert.Equal(t, test, result, "For input %d", test)
	}
}
