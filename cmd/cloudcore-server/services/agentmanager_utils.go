package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

func (s *AgentManagerService) handleNotifications(srv rpc.AgentManager_NotificationsServer, handler func(agentId string, notification *rpc.ClientNotification) error) error {
	agentID, err := s.signer.ValidateRequest(srv.Context())
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}
	client, err := s.server.ClientFromContext(srv.Context())
	if err != nil {
		return status.Errorf(codes.Internal, "failed to get client from context: %v", err)
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
				if errors.Is(err, rpc.ErrAgentDeactivated) {
					requestAgentShutdown(srv.Context(), client, logger)
					return status.Error(codes.NotFound, err.Error())
				}
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

func requestAgentShutdown(ctx context.Context, client rpc.AgentClient, logger *zap.Logger) {
	_, err := client.Shutdown(ctx, &rpc.ShutdownRequest{})
	if err != nil {
		logger.Warn("failed to request agent shutdown", zap.Error(err))
		return
	}
}
