package agent

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

// Server is the gRPC server that runs on the agent itself and is accessible
// by the real gRPC server that the agent connects to.
type Server struct {
	logger *zap.Logger
	tomb   *tomb.Tomb

	rpc.UnimplementedAgentServer
}

func (s *Server) Shutdown(_ context.Context, _ *rpc.ShutdownRequest) (*rpc.ShutdownResponse, error) {
	s.logger.Info("server requested shutdown")
	go func() {
		s.tomb.Kill(rpc.ErrAgentDeactivated)
	}()
	return &rpc.ShutdownResponse{}, nil
}

func NewServer(logger *zap.Logger, tomb *tomb.Tomb) rpc.AgentServer {
	return &Server{
		tomb:   tomb,
		logger: logger,
	}
}
