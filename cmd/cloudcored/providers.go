package main

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
	"os"
	"os/signal"
)

// signaller accepts a global tomb that doesn't need to be provided via
// the fx framework and returns a fx.Provider function that takes control
// of the tomb and connects it to a signal handler. When the signal is
// received, then the tomb is killed.
//
// Why a global tomb rather than a tomb scoped to the fx app you may ask?
// Because we need the final shutdown step of the application to be waiting
// for the tomb to die, and this needs to happen outside the fx app.
func signaller(t *tomb.Tomb) func(logger *zap.Logger) (*tomb.Tomb, context.Context) {
	return func(logger *zap.Logger) (*tomb.Tomb, context.Context) {
		ctx, _ := signal.NotifyContext(t.Context(context.Background()), os.Interrupt)
		go func() {
			<-ctx.Done()
			logger.Info("received shutdown signal")
			t.Kill(ctx.Err())
		}()
		return t, ctx
	}
}

// shutdowner is a fx.Invoke-compatible function that triggers an fx shutdown
// when we see that the tomb is dying.
func shutdowner(s fx.Shutdowner, tomb *tomb.Tomb, logger *zap.Logger) {
	go func() {
		<-tomb.Dying()
		err := s.Shutdown(fx.ExitCode(0))
		if err != nil {
			logger.Error("shutting down", zap.Error(err))
		}
	}()
}
