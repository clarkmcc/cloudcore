package appbackend

import (
	"context"
	"encoding/json"
	"github.com/clarkmcc/cloudcore/app/backend/middleware"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/pkg/packages"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

type Server struct {
	schema   graphql.Schema
	logger   *zap.Logger
	database database.Database
	packages packages.Provider
}

func New(lc fx.Lifecycle, config *config.Config, database database.Database, logger *zap.Logger, packages packages.Provider) (*Server, error) {
	logger = logger.Named("app-backend")
	s, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}
	srv := Server{
		schema:   s,
		logger:   logger,
		database: database,
		packages: packages,
	}

	r := gin.Default()
	r.Use(middleware.CORS(), middleware.Authentication(config, logger.Named("auth")))
	r.POST("/graphql", func(c *gin.Context) {
		var req request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			logger.Warn("failed to decode request", zap.Error(err))
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		res := graphql.Do(graphql.Params{
			Schema:         s,
			RequestString:  req.Query,
			VariableValues: req.Variables,
			OperationName:  req.Operation,
			Context:        srv.graphqlContext(c.Request.Context()),
		})
		c.JSON(http.StatusOK, res)
	})

	var listener net.Listener
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			listener, err = net.Listen("tcp", ":"+strconv.Itoa(config.AppServer.Port))
			if err != nil {
				return err
			}
			logger.Info("listening on port", zap.Int("port", config.AppServer.Port))

			go http.Serve(listener, r.Handler())
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
