package client

import (
	"go.uber.org/zap"
	"time"
)

func keepalive(logger *zap.Logger, do func() error) {
	for {
		err := do()
		if err != nil {
			logger.Error("retrying after error", zap.Error(err))
		}
		time.Sleep(time.Second)
	}
}
