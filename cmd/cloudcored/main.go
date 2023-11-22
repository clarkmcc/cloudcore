package main

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcored/config"
	"github.com/clarkmcc/cloudcore/internal/agent"
	"github.com/clarkmcc/cloudcore/internal/agentdb"
	"github.com/clarkmcc/cloudcore/internal/client"
	"github.com/clarkmcc/cloudcore/internal/events"
	"github.com/clarkmcc/cloudcore/internal/logger"
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
				tomb := tomb.Tomb{}
				ctx, _ := signal.NotifyContext(tomb.Context(context.Background()), os.Interrupt)
				return &tomb, ctx
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
			fx.Provide(agent.NewServer),
			fx.Provide(client.New),
			fx.Provide(tasks.NewExecutor),
			fx.Invoke(func(e *tasks.Executor) {
				e.Initialize()
			}),
			fx.Invoke(events.NewLifecycleNotifications),
			fx.Invoke(func(s fx.Shutdowner, tomb *tomb.Tomb) error {
				<-tomb.Dead()
				return s.Shutdown(fx.ExitCode(0))
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

func init() {
	cmd.PersistentFlags().String("psk", "", "Pre-shared key for authenticating with the server")
	cmd.PersistentFlags().Bool("insecure-skip-verify", false, "Whether to skip verifying the server's TLS certificate")
}

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
