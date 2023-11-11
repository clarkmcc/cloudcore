package main

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/client"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

var cmd = &cobra.Command{
	Use: "cloudcored",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			fx.Provide(func() context.Context {
				ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
				return ctx
			}),
			fx.Provide(config.NewAgentConfig),
			// Extra the logging config from the Agent-specific config
			fx.Provide(func(config *config.AgentConfig) *config.Logging {
				return &config.Logging
			}),
			fx.Provide(logger.New),
			fx.Provide(agentdb.New),
			fx.Provide(client.New),
			fx.Invoke(func(ctx context.Context, client *client.Client, logger *zap.Logger) {
				err := client.Ping(ctx)
				if err != nil {
					logger.Error("failed to ping server", zap.Error(err))
					return
				}
				logger.Info("pinged server")
				err = client.UploadMetadata(ctx, &rpc.SystemMetadata{})
				if err != nil {
					logger.Error("failed to upload metadata", zap.Error(err))
					return
				}
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
