package appbackend

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database/types"
	"github.com/graphql-go/graphql"
	"time"
)

var hostType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Host",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"updatedAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"online": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Resolve: func(p graphql.ResolveParams) (any, error) {
				h := p.Source.(types.Host)
				if h.LastHeartbeatTimestamp.After(time.Now().Add(-time.Hour)) && h.Online {
					return true, nil
				}
				return false, nil
			},
		},
		"identifier": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "An identifier for the host as determined by the agent. This is usually extracted from the host somehow (i.e. a Host ID).",
		},
		"hostname": &graphql.Field{
			Type:        nullStringType,
			Description: "The hostname of the host.",
		},
		"publicIpAddress": &graphql.Field{
			Type:        nullStringType,
			Description: "The public IP address of the host.",
		},
		"privateIpAddress": &graphql.Field{
			Type:        nullStringType,
			Description: "The private IP address of the host.",
		},
		"osName": &graphql.Field{
			Type:        nullStringType,
			Description: "The name of the operating system (i.e. darwin).",
		},
		"osFamily": &graphql.Field{
			Type:        nullStringType,
			Description: "The family of the operating system (i.e. Standalone Workstation).",
		},
		"osVersion": &graphql.Field{
			Type:        nullStringType,
			Description: "The version of the operating system (i.e. 14.0).",
		},
		"kernelArchitecture": &graphql.Field{
			Type:        nullStringType,
			Description: "The architecture of the kernel (i.e. arm64).",
		},
		"kernelVersion": &graphql.Field{
			Type:        nullStringType,
			Description: "The version of the kernel (i.e. 23.0.0).",
		},
		"cpuModel": &graphql.Field{
			Type:        nullStringType,
			Description: "The model of the CPU (i.e. Apple M1 Max).",
		},
		"cpuCores": &graphql.Field{
			Type:        nullInt64Type,
			Description: "The number of CPU cores (i.e. 10).",
		},
	},
})

var hostList = &graphql.Field{
	Type: graphql.NewList(hostType),
	Args: graphql.FieldConfigArgument{
		"projectId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: wrapper(func(rctx resolveContext) ([]types.Host, error) {
		projectID := rctx.getStringArg("projectId")
		err := rctx.canAccessProject(projectID)
		if err != nil {
			return nil, err
		}
		return rctx.db.ListProjectHosts(rctx, projectID)
	}),
}
