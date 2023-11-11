package services

import (
	"context"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	config *config.ServerConfig
	signer *token.Signer

	rpc.UnimplementedAuthenticationServer
}

func (s *AuthService) Ping(_ context.Context, _ *rpc.PingRequest) (*rpc.PingResponse, error) {
	return &rpc.PingResponse{}, nil
}

func (s *AuthService) Authenticate(context.Context, *rpc.AuthenticateRequest) (*rpc.AuthenticateResponse, error) {
	tk, err := s.signer.NewToken()
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	return &rpc.AuthenticateResponse{
		Token: tk,
	}, nil
}

func NewAuthService(config *config.ServerConfig, signer *token.Signer) *AuthService {
	return &AuthService{
		signer: signer,
		config: config,
	}
}
