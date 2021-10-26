package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New returns suggared logger in production state
func New(service string) (*zap.SugaredLogger, error) {
	// Init production level logger
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]interface{}{
		"service": service,
	}

	log, err := config.Build()
	if err != nil {
		return nil, err
	}

	return log.Sugar(), nil
}
