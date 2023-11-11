package tasks

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/client"
	"go.uber.org/zap"
)

const contextKeyLogger = "logger"
const contextKeyClient = "client"
const contextKeyAgentDB = "agentdb"

func (e *Executor) executionContext(ctx context.Context, task *Task) context.Context {
	ctx = context.WithValue(ctx, contextKeyLogger, e.logger.Named(task.Name))
	ctx = context.WithValue(ctx, contextKeyClient, e.client)
	ctx = context.WithValue(ctx, contextKeyAgentDB, e.db)
	return e.tomb.Context(ctx)
}

func Logger(ctx context.Context) *zap.Logger {
	return ctx.Value(contextKeyLogger).(*zap.Logger)
}

func Client(ctx context.Context) *client.Client {
	return ctx.Value(contextKeyClient).(*client.Client)
}

func AgentDB(ctx context.Context) agentdb.AgentDB {
	return ctx.Value(contextKeyAgentDB).(agentdb.AgentDB)
}
