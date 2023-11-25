package main

import (
	"context"
	"github.com/clarkmcc/brpc"
	appbackend "github.com/clarkmcc/cloudcore/app/backend"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/server"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"github.com/clarkmcc/cloudcore/pkg/packages"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	app := fx.New(
		fx.Provide(func() (context.Context, context.CancelFunc) {
			return signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
		}),
		fx.Provide(config.New),
		fx.Provide(componentConfigs),
		fx.Provide(token.NewSigner),
		fx.Provide(database.New),
		fx.Provide(appbackend.New),
		fx.Provide(server.New),
		fx.Provide(fx.Annotate(
			packages.NewGithubReleaseProvider,
			fx.As(new(packages.Provider)))),
		fx.Provide(func(config *config.Config) *zap.Logger {
			return logger.New(config.Logging.Level, config.Logging.Debug)
		}),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger}
		}),
		fx.Invoke(func(ctx context.Context, db database.Database) error {
			return db.Migrate(ctx)
		}),
		fx.Invoke(func(_ *brpc.Server[rpc.AgentClient]) {}),
		fx.Invoke(func(_ *appbackend.Server) {}),
	)
	err := app.Err()
	if err != nil {
		panic(err)
	}
	app.Run()
}
