package appbackend

import (
	"context"
	"encoding/json"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/graphql-go/graphql"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

var _ http.Handler = &Server{}

type Server struct {
	schema graphql.Schema
	logger *zap.Logger
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		s.logger.Warn("failed to decode request", zap.Error(err))
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	res := graphql.Do(graphql.Params{
		Schema:         s.schema,
		RequestString:  req.Query,
		VariableValues: req.Variables,
		OperationName:  req.Operation,
		Context:        r.Context(),
	})
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		s.logger.Warn("failed to encode response", zap.Error(err))
		_, _ = w.Write([]byte(err.Error()))
		return
	}
}

func New(lc fx.Lifecycle, config *config.Config, logger *zap.Logger) (*Server, error) {
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "RootQuery",
			Fields: graphql.Fields{
				"ping": &graphql.Field{
					Type: graphql.String,
					Resolve: func(p graphql.ResolveParams) (any, error) {
						return "pong", nil
					},
				},
			},
		}),
		//Mutation: graphql.NewObject(graphql.ObjectConfig{
		//	Name:   "RootMutation",
		//	Fields: graphql.Fields{},
		//}),
		//Subscription: graphql.NewObject(graphql.ObjectConfig{
		//	Name:   "RootSubscription",
		//	Fields: graphql.Fields{},
		//}),
	})
	if err != nil {
		return nil, err
	}
	srv := Server{
		schema: s,
		logger: logger.Named("app-backend"),
	}
	var listener net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.AppServer.Port))
			if err != nil {
				return err
			}
			srv.logger.Info("listening on port", zap.Int("port", config.AppServer.Port))
			go http.Serve(listener, &srv)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return listener.Close()
		},
	})
	return &srv, nil
}

type request struct {
	Query     string         `json:"query"`
	Operation string         `json:"operationName"`
	Variables map[string]any `json:"variables"`
}
