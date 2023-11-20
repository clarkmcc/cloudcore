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
		"lastHeartbeatTimestamp": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"online": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Resolve: func(p graphql.ResolveParams) (any, error) {
				h := p.Source.(types.Host)
				// todo: make this configurable
				if h.LastHeartbeatTimestamp.After(time.Now().Add(-time.Minute)) && h.Online {
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
		"events": &graphql.Field{
			Type:        graphql.NewList(agentEvent),
			Description: "The events for the host and its agent.",
			Resolve: wrapper[types.Host](func(rctx resolveContext[types.Host]) ([]types.AgentEventLog, error) {
				return rctx.db.GetEventLogsByHost(rctx, rctx.source.ID, 15)
			}),
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
	Resolve: wrapper[any](func(rctx resolveContext[any]) ([]types.Host, error) {
		projectID := rctx.getStringArg("projectId")
		err := rctx.canAccessProject(projectID)
		if err != nil {
			return nil, err
		}
		return rctx.db.ListProjectHosts(rctx, projectID)
	}),
}

var agentEvent = graphql.NewObject(graphql.ObjectConfig{
	Name: "AgentEvent",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"createdAt": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"agentId": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"hostId": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"type": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"message": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})

var hostDetails = &graphql.Field{
	Type: hostType,
	Args: graphql.FieldConfigArgument{
		"hostId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"projectId": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: wrapper[any](func(rctx resolveContext[any]) (types.Host, error) {
		projectID := rctx.getStringArg("projectId")
		err := rctx.canAccessProject(projectID)
		if err != nil {
			return types.Host{}, err
		}
		return rctx.db.GetHost(rctx, rctx.getStringArg("hostId"), projectID)
	}),
}
