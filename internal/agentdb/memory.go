package agentdb

import (
	"context"
	"time"
)

var _ AgentDB = (*memoryDB)(nil)

type AuthToken struct {
	Token      string
	Expiration time.Time
	Duration   time.Duration
}

type memoryDB struct {
	authToken *AuthToken
	agentID   string
}

func (m *memoryDB) AuthToken(_ context.Context) (*AuthToken, error) {
	if m.authToken == nil {
		return nil, ErrAuthTokenNotFound
	}
	return m.authToken, nil
}

func (m *memoryDB) SaveAuthToken(_ context.Context, token *AuthToken) error {
	m.authToken = token
	return nil
}

func (m *memoryDB) AgentID(_ context.Context) (string, error) {
	return m.agentID, nil
}

func (m *memoryDB) SaveAgentID(_ context.Context, agentID string) error {
	m.agentID = agentID
	return nil
}

func newMemoryDB() (*memoryDB, error) {
	return &memoryDB{}, nil
}
