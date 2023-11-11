package server

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/services"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

func Listener(config *config.ServerConfig, logger *zap.Logger) (net.Listener, error) {
	logger.Info("listening on port", zap.Int("port", config.Port))
	return net.Listen("tcp", ":"+strconv.Itoa(config.Port))
}

func New(
	lc fx.Lifecycle,
	logger *zap.Logger,
	listener net.Listener,
	auth *services.AuthService,
	agent *services.AgentManagerService,
) *grpc.Server {
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(loggingUnaryInterceptor(logger)),
		grpc.StreamInterceptor(loggingStreamInterceptor(logger)))
	rpc.RegisterAuthenticationServer(srv, auth)
	rpc.RegisterAgentManagerServer(srv, agent)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := srv.Serve(listener)
				if err != nil {
					logger.Fatal("failed to serve", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.GracefulStop()
			return nil
		},
	})
	return srv
}
