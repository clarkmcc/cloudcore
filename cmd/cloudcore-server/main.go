package main

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/server"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/services"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/token"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	app := fx.New(
		fx.Provide(config.New),
		fx.Provide(token.NewSigner),
		fx.Provide(database.New),
		fx.Provide(services.NewAuthService),
		fx.Provide(services.NewAgentManagerService),
		fx.Provide(server.Listener),
		// Extra the logging config from the Agent-specific config
		fx.Provide(func(config *config.Config) *config.Logging {
			return &config.Logging
		}),
		fx.Provide(func(config *config.Config) *zap.Logger {
			return logger.New(config.Logging.Level, config.Logging.Debug)
		}),
		fx.Provide(server.New),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(db database.Database) error {
			return db.Migrate()
		}),
		fx.Invoke(func(_ *grpc.Server) {}),
	)
	err := app.Err()
	if err != nil {
		panic(err)
	}
	app.Run()
}
