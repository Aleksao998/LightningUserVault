package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestConvertStringToLogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected zapcore.Level
		err      bool
	}{
		{"ERROR", zapcore.ErrorLevel, false},
		{"error", zapcore.ErrorLevel, false},
		{"ErRoR", zapcore.ErrorLevel, false},
		{"INFO", zapcore.InfoLevel, false},
		{"info", zapcore.InfoLevel, false},
		{"DEBUG", zapcore.DebugLevel, false},
		{"debug", zapcore.DebugLevel, false},
		{"WARN", zapcore.WarnLevel, false},
		{"warn", zapcore.WarnLevel, false},
		{"INVALID", zapcore.InfoLevel, true}, // Defaulting to InfoLevel for invalid input
		{"", zapcore.InfoLevel, true},        // Defaulting to InfoLevel for empty input
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			got, err := ConvertStringToLogLevel(tt.input)
			if tt.err {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", tt.input)
				assert.Equal(t, tt.expected, got, "Expected %v but got %v for input: %s", tt.expected, got, tt.input)
			}
		})
	}
}
