package registeredtasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/tasks"
	"time"
)

func init() {
	tasks.DefaultRegistry.Register(&tasks.Task{
		Name:     "heartbeat",
		Schedule: tasks.Interval(10 * time.Second),
		Action: func(ctx context.Context) error {
			return tasks.Client(ctx).Notify(ctx, &rpc.ClientNotification{
				Type: rpc.ClientNotification_HEARTBEAT,
			})
		},
	})
}
