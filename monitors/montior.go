package monitors

import (
	"context"
	"time"
)

type Status string

const (
	StatusUp       Status = "up"
	StatusDown     Status = "down"
	StatusDegraded Status = "degraded"
	StatusUnknown  Status = "unknown"
	StatusWarning  Status = "warning"
)

func GetStatusIcon(status Status) string {
	switch status {
	case StatusUp:
		return "✅"
	case StatusDown:
		return "❌"
	case StatusDegraded:
		return "🟠"
	case StatusWarning:
		return "⚠️"
	case StatusUnknown:
		return "❓"
	default:
		return "🔷"
	}
}

type Result struct {
	MonitorID string        `json:"monitor_id"`
	Type      string        `json:"type"`
	Status    Status        `json:"status"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Message   string        `json:"message,omitempty"`
	Error     string        `json:"error,omitempty"`
	Success   bool          `json:"success"`
	CheckedAt time.Time     `json:"checked_at"`
	Details   any           `json:"details,omitempty"`
}

type Monitor interface {
	GetID() string
	GetContext() context.Context
	Start(ctx context.Context)
	Stop()
	IsRunning() bool
	SetRunning(bool)
	Check(ctx context.Context) Result
	GetHBInterval() time.Duration
	GetLastHB() time.Time
	SetLastHB(time.Time)
}
