package types

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// ConvertStringToLogLevel converts a string to its corresponding zapcore.Level
func ConvertStringToLogLevel(s string) (zapcore.Level, error) {
	switch strings.ToUpper(s) {
	case "ERROR":
		return zapcore.ErrorLevel, nil
	case "INFO":
		return zapcore.InfoLevel, nil
	case "DEBUG":
		return zapcore.DebugLevel, nil
	case "WARN":
		return zapcore.WarnLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s", s) // Defaulting to InfoLevel for invalid input
	}
}
