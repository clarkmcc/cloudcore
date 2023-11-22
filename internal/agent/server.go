package agent

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
	"gopkg.in/tomb.v2"
)

type Server struct {
	logger *zap.Logger
	tomb   *tomb.Tomb
	rpc.UnimplementedAgentServer
}

func (s *Server) Shutdown(ctx context.Context, req *rpc.ShutdownRequest) (*rpc.ShutdownResponse, error) {
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
