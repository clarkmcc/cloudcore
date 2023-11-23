package tasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agent"
	"go.uber.org/zap"
)

type contextKeyLogger struct{}
type contextKeyClient struct{}
type contextKeyDatabase struct{}

// executionContext builds a context.Context that contains all the dependencies
// that a task may want to extract.
func (e *Executor) executionContext(ctx context.Context, task *Task) context.Context {
	ctx = context.WithValue(ctx, contextKeyLogger{}, e.logger.Named(task.Name))
	ctx = context.WithValue(ctx, contextKeyClient{}, e.client)
	ctx = context.WithValue(ctx, contextKeyDatabase{}, e.db)
	return e.tomb.Context(ctx)
}

// Logger returns the task's logger from the context
func Logger(ctx context.Context) *zap.Logger {
	return ctx.Value(contextKeyLogger{}).(*zap.Logger)
}

// Client returns the gRPC client from the context
func Client(ctx context.Context) *agent.Client {
	return ctx.Value(contextKeyClient{}).(*agent.Client)
}

// Database returns the agent database available in the task context
func Database(ctx context.Context) agent.Database {
	return ctx.Value(contextKeyDatabase{}).(agent.Database)
}
