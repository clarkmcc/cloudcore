package events

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/client"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/fx"
)

func NewLifecycleNotifications(lc fx.Lifecycle, client *client.Client) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Notify(ctx, &rpc.ClientNotification{
				Type: rpc.ClientNotification_AGENT_STARTUP,
			})
		},
		OnStop: func(context.Context) error {
			return client.Notify(context.Background(), &rpc.ClientNotification{
				Type: rpc.ClientNotification_AGENT_SHUTDOWN,
			})
		},
	})
}
