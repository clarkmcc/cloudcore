package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func New(level string, debug bool) *zap.Logger {
	l := zapcore.InfoLevel
	if v, err := zapcore.ParseLevel(level); err == nil {
		l = v
	}
	cfg := zap.NewProductionEncoderConfig()
	var encoder zapcore.Encoder
	if debug {
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder
		encoder = zapcore.NewConsoleEncoder(cfg)
	} else {
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	return zap.New(zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l))
}
