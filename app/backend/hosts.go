package appbackend

import "github.com/graphql-go/graphql"

var hostType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Host",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
		"created_at": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"updated_at": &graphql.Field{
			Type: graphql.NewNonNull(graphql.DateTime),
		},
		"status": &graphql.Field{
			Type: graphql.NewNonNull(statusType),
		},
		"identifier": &graphql.Field{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "An identifier for the host as determined by the agent. This is usually extracted from the host somehow (i.e. a Host ID).",
		},
		"hostname": &graphql.Field{
			Type:        graphql.String,
			Description: "The hostname of the host.",
		},
		"public_ip_address": &graphql.Field{
			Type:        graphql.String,
			Description: "The public IP address of the host.",
		},
		"private_ip_address": &graphql.Field{
			Type:        graphql.String,
			Description: "The private IP address of the host.",
		},
		"os_name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the operating system (i.e. darwin).",
		},
		"os_family": &graphql.Field{
			Type:        graphql.String,
			Description: "The family of the operating system (i.e. Standalone Workstation).",
		},
		"os_version": &graphql.Field{
			Type:        graphql.String,
			Description: "The version of the operating system (i.e. 14.0).",
		},
		"kernel_architecture": &graphql.Field{
			Type:        graphql.String,
			Description: "The architecture of the kernel (i.e. arm64).",
		},
		"kernel_version": &graphql.Field{
			Type:        graphql.String,
			Description: "The version of the kernel (i.e. 23.0.0).",
		},
		"cpu_model": &graphql.Field{
			Type:        graphql.String,
			Description: "The model of the CPU (i.e. Apple M1 Max).",
		},
		"cpu_cores": &graphql.Field{
			Type:        graphql.Int,
			Description: "The number of CPU cores (i.e. 10).",
		},
	},
})

var hostList = &graphql.Field{
	Type:    graphql.NewList(hostType),
	Args:    graphql.FieldConfigArgument{},
	Resolve: nil,
}
