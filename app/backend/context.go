package appbackend

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"go.uber.org/zap"
)

const contextKeyLogger = "logger"
const contextKeyDatabase = "database"

func (s *Server) graphqlContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, contextKeyLogger, s.logger)
	ctx = context.WithValue(ctx, contextKeyDatabase, s.database)
	return ctx
}

func logger(ctx context.Context) *zap.Logger {
	return ctx.Value(contextKeyLogger).(*zap.Logger)
}

func db(ctx context.Context) database.Database {
	return ctx.Value(contextKeyDatabase).(database.Database)
}
