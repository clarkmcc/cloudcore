package appbackend

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/graphql-go/graphql"
	"github.com/spf13/cast"
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

type resolveContext[S any] struct {
	context.Context
	db     database.Database
	params graphql.ResolveParams
	logger *zap.Logger
	source S
}

func (r *resolveContext[S]) getStringArg(name string) string {
	return cast.ToString(r.params.Args[name])
}

func (r *resolveContext[S]) canAccessProject(projectId string) error {
	sub := middleware.SubjectFromContext(r)
	if sub == "" {
		return fmt.Errorf("no subject in context")
	}
	ok, err := r.db.CanAccessProject(r, sub, projectId)
	if err != nil {
		return fmt.Errorf("checking if user can access project: %w", err)
	}
	if !ok {
		return fmt.Errorf("Not authorized to access project")
	}
	return nil
}

type resolverFunc[T any, S any] func(rctx resolveContext[S]) (T, error)

func wrapper[S any, T any](fn resolverFunc[T, S]) func(params graphql.ResolveParams) (any, error) {
	return func(params graphql.ResolveParams) (any, error) {
		ctx := params.Context
		source, ok := params.Source.(S)
		if !ok {
			var s S
			return nil, fmt.Errorf("invalid source type: %T, expected %T", params.Source, s)
		}
		return fn(resolveContext[S]{
			Context: ctx,
			db:      db(ctx),
			logger:  logger(ctx),
			source:  source,
			params:  params,
		})
	}
}
