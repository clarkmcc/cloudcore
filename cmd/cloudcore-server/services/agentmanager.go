package services

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
)

type AgentManagerService struct {
	config *config.ServerConfig
	logger *zap.Logger

	rpc.UnimplementedAgentManagerServer
}

func (s *AgentManagerService) UploadMetadata(ctx context.Context, req *rpc.UploadMetadataRequest) (*rpc.UploadMetadataResponse, error) {
	id := req.GetSystemMetadata().GetIdentifiers().GetAgentIdentifier()
	s.logger.Info("received metadata", zap.String("id", id))
	return &rpc.UploadMetadataResponse{}, nil
}

func NewAgentManagerService(config *config.ServerConfig, logger *zap.Logger) *AgentManagerService {
	return &AgentManagerService{
		logger: logger.Named("agent-manager"),
		config: config,
	}
}
