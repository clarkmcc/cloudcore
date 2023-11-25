package appbackend

import (
	"context"
	"fmt"
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/pkg/packages"
	"github.com/graphql-go/graphql"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type contextKeyLogger struct{}
type contextKeyDatabase struct{}
type contextKeyPackages struct{}

func (s *Server) graphqlContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, contextKeyLogger{}, s.logger)
	ctx = context.WithValue(ctx, contextKeyDatabase{}, s.database)
	ctx = context.WithValue(ctx, contextKeyPackages{}, s.packages)
	return ctx
}

type resolveContext[S any] struct {
	context.Context
	db       database.Database
	params   graphql.ResolveParams
	logger   *zap.Logger
	packages packages.Provider
	source   S
}

func (r *resolveContext[S]) getStringArg(name string) string {
	return cast.ToString(r.params.Args[name])
}

func (r *resolveContext[S]) getBoolArg(name string) bool {
	return cast.ToBool(r.params.Args[name])
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
			Context:  ctx,
			db:       ctx.Value(contextKeyDatabase{}).(database.Database),
			packages: ctx.Value(contextKeyPackages{}).(packages.Provider),
			logger:   ctx.Value(contextKeyLogger{}).(*zap.Logger),
			source:   source,
			params:   params,
		})
	}
}
