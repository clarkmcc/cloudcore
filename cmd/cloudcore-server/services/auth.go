package services

import (
	"context"
	"github.com/clarkmcc/cloudcore/cmd/cloudcore-server/database"
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
	"github.com/clarkmcc/cloudcore/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	config *config.ServerConfig
	signer *token.Signer
	db     database.Database

	rpc.UnimplementedAuthenticationServer
}

func (s *AuthService) Ping(_ context.Context, _ *rpc.PingRequest) (*rpc.PingResponse, error) {
	return &rpc.PingResponse{}, nil
}

func (s *AuthService) Authenticate(ctx context.Context, req *rpc.AuthenticateRequest) (*rpc.AuthenticateResponse, error) {
	switch req.Flow {
	case rpc.AuthenticateRequest_TOKEN:
		// when using a token, the agent must have already been authenticated
		// and will already exist in the database. We issue a token as long as
		// the current token is valid.
		err := s.signer.ValidateToken(req.Token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		tk, err := s.signer.NewToken()
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return &rpc.AuthenticateResponse{
			Token: tk,
		}, nil
	case rpc.AuthenticateRequest_PRE_SHARED_KEY:
		// todo: lookup pre-shared key
		agentID, err := s.db.AuthenticateAgent(ctx, req.GetPreSharedKey(), req.GetSystemMetadata())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		tk, err := s.signer.NewToken()
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		return &rpc.AuthenticateResponse{
			AgentId: agentID,
			Token:   tk,
		}, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid authentication flow")
	}
}

func NewAuthService(config *config.ServerConfig, signer *token.Signer, db database.Database) *AuthService {
	return &AuthService{
		signer: signer,
		config: config,
		db:     db,
	}
}
