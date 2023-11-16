package logger

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcored/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(config *config.Logging) *zap.Logger {
	level := zap.InfoLevel
	if v, err := zapcore.ParseLevel(config.Level); err == nil {
		level = v
	}
	cfg := zap.NewProductionEncoderConfig()
	var encoder zapcore.Encoder
	if config.Debug {
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
	} else {
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level))
}
