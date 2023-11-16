package services

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
)

type AgentManagerService struct {
	config *config.Config
	logger *zap.Logger
	db     database.Database

	rpc.UnimplementedAgentManagerServer
}

func (s *AgentManagerService) UploadMetadata(ctx context.Context, req *rpc.UploadMetadataRequest) (*rpc.UploadMetadataResponse, error) {
	//id, err := s.db.UpdateMetadata(ctx, req.GetSystemMetadata())
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	//s.logger.Info("received metadata", zap.String("id", id))
	return &rpc.UploadMetadataResponse{}, nil
}

func NewAgentManagerService(config *config.Config, logger *zap.Logger, db database.Database) *AgentManagerService {
	return &AgentManagerService{
		logger: logger.Named("agent-manager"),
		config: config,
		db:     db,
	}
}
