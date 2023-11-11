package services

import (
	"github.com/clarkmcc/cloudcore/internal/config"
	"github.com/clarkmcc/cloudcore/internal/rpc"
)

type AgentManagerService struct {
	config *config.ServerConfig

	rpc.UnimplementedAgentManagerServer
}

func NewAgentManagerService(config *config.ServerConfig) *AgentManagerService {
	return &AgentManagerService{
		config: config,
	}
}
