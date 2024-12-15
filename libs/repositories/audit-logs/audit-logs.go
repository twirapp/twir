package audit_logs

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/audit-logs/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string, limit int) ([]model.AuditLog, error)
}
