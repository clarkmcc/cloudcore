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
			db := tasks.AgentDB(ctx)
			md, err := sysinfo.BuildSystemMetadata(ctx, db, tasks.Logger(ctx))
			if err != nil {
				return err
			}
			res, err := tasks.Client(ctx).UploadMetadata(ctx, md)
			if err != nil {
				return err
			}
			return db.SaveAgentID(ctx, res.GetAgentId())
		},
	})
}
