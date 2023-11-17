package appbackend

import "github.com/graphql-go/graphql"

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
