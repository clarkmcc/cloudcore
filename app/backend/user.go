package appbackend

import (
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
	"github.com/graphql-go/graphql"
)

var ensureUser = &graphql.Field{
	Type: graphql.NewList(projectType),
	Args: graphql.FieldConfigArgument{},
	Resolve: wrapper[any](func(rctx resolveContext[any]) ([]types.Project, error) {
		return rctx.db.UpsertUser(rctx, middleware.SubjectFromContext(rctx))
	}),
}
