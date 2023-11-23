package agent

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewLifecycleNotifications(lc fx.Lifecycle, client *Client, logger *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return logError(logger, client.Notify(ctx, &rpc.ClientNotification{
				Type: rpc.ClientNotification_AGENT_STARTUP,
			}))
		},
		OnStop: func(context.Context) error {
			return logError(logger, client.Notify(context.Background(), &rpc.ClientNotification{
				Type: rpc.ClientNotification_AGENT_SHUTDOWN,
			}))
		},
	})
}

func logError(logger *zap.Logger, err error) error {
	if err != nil {
		logger.Error("failed to send lifecycle notification", zap.Error(err))
	}
	return nil
}
