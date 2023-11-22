package services

import (
	"context"
	"github.com/clarkmcc/brpc"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"go.uber.org/zap"
)

type AgentManagerService struct {
	config   *config.Config
	logger   *zap.Logger
	signer   *token.Signer
	db       database.Database
	shutdown <-chan struct{}
	server   *brpc.Server[rpc.AgentClient]

	rpc.UnimplementedAgentManagerServer
}

func (s *AgentManagerService) UploadMetadata(ctx context.Context, req *rpc.UploadMetadataRequest) (*rpc.UploadMetadataResponse, error) {
	//agentID, err := s.signer.ValidateRequest(ctx)
	//if err != nil {
	//	return nil, status.Error(codes.Unauthenticated, err.Error())
	//}
	//id, err := s.db.UpdateMetadata(ctx, req.GetSystemMetadata())
	//if err != nil {
	//	return nil, status.Error(codes.Internal, err.Error())
	//}
	return &rpc.UploadMetadataResponse{}, nil
}

func (s *AgentManagerService) Notifications(srv rpc.AgentManager_NotificationsServer) error {
	return s.handleNotifications(srv, func(agentID string, c *rpc.ClientNotification) error {
		switch c.GetType() {
		case rpc.ClientNotification_HEARTBEAT:
			return s.db.Heartbeat(srv.Context(), agentID)
		case rpc.ClientNotification_AGENT_SHUTDOWN:
			// use background context in case client connection closes before
			// we're able to finish what we need to do.
			return s.db.AgentShutdown(context.Background(), agentID)
		case rpc.ClientNotification_AGENT_STARTUP:
			return s.db.AgentStartup(srv.Context(), agentID)
		default:
			return nil
		}
	})
}

func NewAgentManagerService(ctx context.Context, config *config.Config, logger *zap.Logger, signer *token.Signer, db database.Database, server *brpc.Server[rpc.AgentClient]) *AgentManagerService {
	return &AgentManagerService{
		server:   server,
		shutdown: ctx.Done(),
		logger:   logger.Named("agent-manager"),
		config:   config,
		signer:   signer,
		db:       db,
	}
}
