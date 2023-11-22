package appbackend

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
	"github.com/graphql-go/graphql"
)

var dashboardMetricsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ProjectMetrics",
	Fields: graphql.Fields{
		"totalHosts": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"totalAgents": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"onlineHosts": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"offlineHosts": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Int),
		},
		"hostsByOsName": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "OsNameCount",
				Fields: graphql.Fields{
					"osName": &graphql.Field{
						Type: graphql.NewNonNull(graphql.String),
					},
					"count": &graphql.Field{
						Type: graphql.NewNonNull(graphql.Int),
					},
				},
			})),
		},
	},
})

var projectMetrics = &graphql.Field{
	Type: dashboardMetricsType,
	Args: graphql.FieldConfigArgument{
		"projectId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: wrapper[any](func(rctx resolveContext[any]) (*types.ProjectMetrics, error) {
		projectID := rctx.getStringArg("projectId")
		err := rctx.canAccessProject(projectID)
		if err != nil {
			return nil, err
		}
		return rctx.db.GetProjectMetrics(rctx, projectID)
	}),
}
