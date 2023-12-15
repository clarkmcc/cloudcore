package server

import (
	"context"
	"fmt"
	"github.com/clarkmcc/brpc"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/services"
	"github.com/clarkmcc/cloudcore/internal/envtls"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"github.com/quic-go/quic-go"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"strconv"
)

func New(
	lc fx.Lifecycle,
	ctx context.Context,
	config *config.Config,
	logger *zap.Logger,
	signer *token.Signer,
	db database.Database,
) *brpc.Server[rpc.AgentClient] {
	srv := brpc.NewServer(brpc.ServerConfig[rpc.AgentClient]{
		// Create our gRPC server
		Server: grpc.NewServer(
			grpc.UnaryInterceptor(loggingUnaryInterceptor(logger)),
			grpc.StreamInterceptor(loggingStreamInterceptor(logger))),
		// Create a gRPC client service which allows for server->client RPCs
		ClientServiceBuilder: func(cc grpc.ClientConnInterface) rpc.AgentClient {
			return rpc.NewAgentClient(cc)
		},
	})
	rpc.RegisterAuthenticationServer(srv, services.NewAuthService(config, signer, db))
	rpc.RegisterAgentManagerServer(srv, services.NewAgentManagerService(ctx, config, logger, signer, db, srv))

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			cfg, err := envtls.TLSConfig()
			if err != nil {
				return fmt.Errorf("getting tls config: %w", err)
			}
			l, err := quic.ListenAddr(":"+strconv.Itoa(config.AgentServer.Port), cfg, nil)
			if err != nil {
				return err
			}
			go func() {
				err := srv.Serve(ctx, l)
				if err != nil {
					logger.Fatal("failed to serve", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Server.GracefulStop()
			return nil
		},
	})
	return srv
}
