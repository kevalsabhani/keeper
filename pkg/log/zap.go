package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a Zap logger configured for the given environment.
// In production it uses JSON encoding with ISO8601 timestamps;
// in all other environments it uses colored console output for readability.
func New(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return config.Build()
}
