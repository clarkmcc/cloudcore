package registeredtasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/tasks"
)

func init() {
	tasks.DefaultRegistry.Register(&tasks.Task{
		Name:     "metadata",
		Schedule: tasks.RunOnceNow{},
		Action: func(ctx context.Context) error {
			return tasks.GetClient(ctx).UploadMetadata(ctx, &rpc.SystemMetadata{})
		},
	})
}
