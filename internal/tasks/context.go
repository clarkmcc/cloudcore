package tasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/client"
	"go.uber.org/zap"
)

const contextKeyLogger = "logger"
const contextKeyClient = "client"

func (e *Executor) executionContext(ctx context.Context, task *Task) context.Context {
	ctx = context.WithValue(ctx, contextKeyLogger, e.logger.Named(task.Name))
	ctx = context.WithValue(ctx, contextKeyClient, e.client)
	return e.tomb.Context(ctx)
}

func GetLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(contextKeyLogger).(*zap.Logger)
}

func GetClient(ctx context.Context) *client.Client {
	return ctx.Value(contextKeyClient).(*client.Client)
}
