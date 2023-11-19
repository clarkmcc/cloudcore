package main

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcored/config"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/client"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/tasks"
	_ "github.com/clarkmcc/cloudcore/internal/tasks/registered"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
	"os"
	"os/signal"
)

var cmd = &cobra.Command{
	Use: "cloudcored",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			fx.Provide(func() *cobra.Command {
				return cmd
			}),
			fx.Provide(func() (*tomb.Tomb, context.Context) {
				ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
				return tomb.WithContext(ctx)
			}),
			fx.Provide(config.New),
			// Extra the logging config from the Agent-specific config
			fx.Provide(func(config *config.Config) *config.Logging {
				return &config.Logging
			}),
			fx.Provide(func(config *config.Config) *zap.Logger {
				return logger.New(config.Logging.Level, config.Logging.Debug)
			}),
			fx.Provide(agentdb.New),
			fx.Provide(client.New),
			fx.Provide(tasks.NewExecutor),
			fx.Invoke(func(e *tasks.Executor) {
				e.Initialize()
			}),
			// Register a hook that will notify the server when we shut down
			fx.Invoke(func(lc fx.Lifecycle, client *client.Client) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return client.Notify(context.Background(), &rpc.ClientNotification{
							Type: rpc.ClientNotification_AGENT_SHUTDOWN,
						})
					},
				})
			}),
		)
		err := app.Err()
		if err != nil {
			return err
		}
		app.Run()
		return nil
	},
}

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
