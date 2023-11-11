package main

import (
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/server"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/services"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/token"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	app := fx.New(
		fx.Provide(config.NewServerConfig),
		fx.Provide(token.NewSigner),
		fx.Provide(services.NewAuthService),
		fx.Provide(services.NewAgentManagerService),
		fx.Provide(server.Listener),
		// Extra the logging config from the Agent-specific config
		fx.Provide(func(config *config.ServerConfig) *config.Logging {
			return &config.Logging
		}),
		fx.Provide(logger.New),
		fx.Provide(server.New),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(_ *grpc.Server) {}),
	)
	err := app.Err()
	if err != nil {
		panic(err)
	}
	app.Run()
}
