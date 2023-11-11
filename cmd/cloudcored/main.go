package main

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/client"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/tasks"
	_ "github.com/clarkmcc/cloudcore/internal/tasks/registered"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
	"gopkg.in/tomb.v2"
	"os"
	"os/signal"
)

var cmd = &cobra.Command{
	Use: "cloudcored",
	RunE: func(cmd *cobra.Command, args []string) error {
		app := fx.New(
			fx.Provide(func() (*tomb.Tomb, context.Context) {
				ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
				return tomb.WithContext(ctx)
			}),
			fx.Provide(config.NewAgentConfig),
			// Extra the logging config from the Agent-specific config
			fx.Provide(func(config *config.AgentConfig) *config.Logging {
				return &config.Logging
			}),
			fx.Provide(logger.New),
			fx.Provide(agentdb.New),
			fx.Provide(client.New),
			fx.Provide(tasks.NewExecutor),
			fx.Invoke(func(e *tasks.Executor) {
				e.Initialize()
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
