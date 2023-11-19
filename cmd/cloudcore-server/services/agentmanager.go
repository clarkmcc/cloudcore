package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/config"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

type AgentManagerService struct {
	config   *config.Config
	logger   *zap.Logger
	signer   *token.Signer
	db       database.Database
	shutdown <-chan struct{}

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
			return s.db.SetOffline(context.Background(), agentID)
		default:
			return nil
		}
	})
}

func NewAgentManagerService(ctx context.Context, config *config.Config, logger *zap.Logger, signer *token.Signer, db database.Database) *AgentManagerService {
	return &AgentManagerService{
		shutdown: ctx.Done(),
		logger:   logger.Named("agent-manager"),
		config:   config,
		signer:   signer,
		db:       db,
	}
}

func (s *AgentManagerService) handleNotifications(srv rpc.AgentManager_NotificationsServer, handler func(agentId string, notification *rpc.ClientNotification) error) error {
	agentID, err := s.signer.ValidateRequest(srv.Context())
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	logger := s.logger.Named("notifications")
	for {
		ns, errs := receive(srv.Recv)
		select {
		case <-s.shutdown:
			return fmt.Errorf("server shutdown")
		case <-srv.Context().Done():
			return srv.Context().Err()
		case err := <-errs:
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}
				return err
			}
		case n := <-ns:
			logger.Debug("received notification",
				zap.String("notification", n.GetType().String()),
				zap.String("agent", agentID))
			err = handler(agentID, n)
			if err != nil {
				logger.Warn("failed to handle notification", zap.Error(err))
				continue
			}
		}
	}
}

func receive[T any](recv func() (T, error)) (<-chan T, <-chan error) {
	ch := make(chan T)
	errs := make(chan error)
	go func() {
		defer close(ch)
		v, err := recv()
		if err != nil {
			errs <- err
			return
		}
		ch <- v
	}()
	return ch, errs
}
