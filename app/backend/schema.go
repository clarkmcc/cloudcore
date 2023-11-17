package appbackend

import "github.com/graphql-go/graphql"

var schemaConfig = graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"ping": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (any, error) {
					return "pong", nil
				},
			},
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"ensureUser": ensureUser,
		},
	}),
}

var statusType = graphql.NewEnum(graphql.EnumConfig{
	Name: "Status",
	Values: graphql.EnumValueConfigMap{
		"active": &graphql.EnumValueConfig{
			Value: "active",
		},
		"deleted": &graphql.EnumValueConfig{
			Value: "deleted",
		},
	},
})
