package registeredtasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/sysinfo"
	"github.com/clarkmcc/cloudcore/internal/tasks"
	"time"
)

func init() {
	tasks.DefaultRegistry.Register(&tasks.Task{
		Name:     "metadata",
		Schedule: tasks.Interval(10 * time.Second),
		Action: func(ctx context.Context) error {
			md, err := sysinfo.BuildSystemMetadata(ctx, tasks.AgentDB(ctx), tasks.Logger(ctx))
			if err != nil {
				return err
			}
			return tasks.Client(ctx).UploadMetadata(ctx, md)
		},
	})
}
