package appbackend

import (
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/graphql-go/graphql"
)

var ensureUser = &graphql.Field{
	Type: graphql.NewList(projectType),
	Args: graphql.FieldConfigArgument{},
	Resolve: func(p graphql.ResolveParams) (any, error) {
		db := db(p.Context)
		subject := middleware.SubjectFromContext(p.Context)
		return db.UpsertUser(p.Context, subject)
	},
}
