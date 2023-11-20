package types

import "time"

type AgentEventType string

const (
	AgentEventType_AGENT_STARTUP  AgentEventType = "AGENT_STARTUP"
	AgentEventType_AGENT_SHUTDOWN AgentEventType = "AGENT_SHUTDOWN"
)

type AgentEventLog struct {
	ID        string         `db:"id"`
	CreatedAt time.Time      `db:"created_at"`
	ProjectID string         `db:"project_id"`
	AgentID   string         `db:"agent_id"`
	HostID    string         `db:"host_id"`
	Type      AgentEventType `db:"type"`
	Message   string         `db:"message"`
}
