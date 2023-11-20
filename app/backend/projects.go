package appbackend

import (
	"fmt"
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/graphql-go/graphql"
)

var projectType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Project",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"created_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"updated_at": &graphql.Field{
			Type: graphql.DateTime,
		},
		"status": &graphql.Field{
			Type: statusType,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var projectCreate = &graphql.Field{
	Type: graphql.NewObject(graphql.ObjectConfig{
		Name: "ProjectCreate",
		Fields: graphql.Fields{
			"project": &graphql.Field{
				Type: projectType,
			},
			"allProjects": &graphql.Field{
				Type: graphql.NewList(projectType),
			},
		},
	}),
	Args: graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"description": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	},
	Resolve: wrapper[any](func(rctx resolveContext[any]) (map[string]any, error) {
		sub := middleware.SubjectFromContext(rctx)
		project, err := rctx.db.CreateProject(rctx, sub,
			rctx.getStringArg("name"),
			rctx.getStringArg("description"))
		if err != nil {
			return nil, fmt.Errorf("creating new project: %w", err)
		}
		projects, err := rctx.db.GetUserProjects(rctx, sub)
		if err != nil {
			return nil, fmt.Errorf("getting user projects: %w", err)
		}
		return map[string]any{
			"project":     project,
			"allProjects": projects,
		}, nil
	}),
}
