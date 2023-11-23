package main

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/agent"
	"github.com/clarkmcc/cloudcore/internal/logger"
	"github.com/clarkmcc/cloudcore/internal/sysinfo"
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
			fx.Provide(literal(cmd)),
			fx.Provide(agent.NewConfig),
			fx.Provide(agent.NewDatabase),
			fx.Provide(agent.NewServer),
			fx.Provide(agent.NewClient),
			fx.Provide(tasks.NewExecutor),
			fx.Provide(fx.Annotate(
				sysinfo.NewSystemMetadataProvider,
				fx.As(new(agent.SystemMetadataProvider)))),
			fx.Invoke(agent.NewLifecycleNotifications),
			fx.Provide(func() (*tomb.Tomb, context.Context) {
				tomb := tomb.Tomb{}
				ctx, _ := signal.NotifyContext(tomb.Context(context.Background()), os.Interrupt)
				return &tomb, ctx
			}),
			fx.Decorate(func(config *agent.Config) *agent.Logging {
				return &config.Logging
			}),
			fx.Provide(func(config *agent.Config) *zap.Logger {
				return logger.New(config.Logging.Level, config.Logging.Debug)
			}),
			fx.Invoke(func(e *tasks.Executor) {
				e.Initialize()
			}),
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

// literal returns a fx provider function that returns the value
// passed to this function. It is a utility that avoids having
// to write a full anonymous inline function just to literal a
// type to fx.
func literal[T any](v T) func() T {
	return func() T {
		return v
	}
}
