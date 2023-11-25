package appbackend

import (
	"database/sql"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/spf13/cast"
)

var schemaConfig = graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"hosts":          hostList,
			"host":           hostDetails,
			"projectMetrics": projectMetrics,
			"packages":       listPackages,
		},
	}),
	Mutation: graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"ensureUser":              ensureUser,
			"projectCreate":           projectCreate,
			"buildDeployAgentCommand": buildDeployAgentCommand,
			"hostGroupCreate":         newHostGroup,
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

var nullStringType = graphql.NewScalar(graphql.ScalarConfig{
	Name: "NullString",
	Serialize: func(value any) any {
		switch v := value.(type) {
		case sql.NullString:
			return v.String
		case string:
			return v
		default:
			return cast.ToString(v)
		}
	},
	ParseValue: func(value any) any {
		switch v := value.(type) {
		case string:
			return sql.NullString{String: v, Valid: true}
		case nil:
			return sql.NullString{Valid: false}
		case *string:
			if v == nil {
				return sql.NullString{Valid: false}
			}
			return sql.NullString{String: *v, Valid: true}
		default:
			return sql.NullString{String: cast.ToString(v), Valid: true}
		}
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		return valueAST.GetValue()
	},
})

var nullInt64Type = graphql.NewScalar(graphql.ScalarConfig{
	Name: "NullInt64",
	Serialize: func(value any) any {
		switch v := value.(type) {
		case sql.NullInt64:
			return v.Int64
		case int64:
			return v
		default:
			return cast.ToInt64(v)
		}
	},
	ParseValue: func(value any) any {
		switch v := value.(type) {
		case int64:
			return sql.NullInt64{Int64: v, Valid: true}
		case nil:
			return sql.NullInt64{Valid: false}
		case *int64:
			if v == nil {
				return sql.NullInt64{Valid: false}
			}
			return sql.NullInt64{Int64: *v, Valid: true}
		default:
			return sql.NullInt64{Int64: cast.ToInt64(v), Valid: true}
		}
	},
	ParseLiteral: func(valueAST ast.Value) any {
		return valueAST.GetValue()
	},
})
