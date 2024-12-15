package user_audit_log

import (
	"context"

	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	auditlogsrepository "github.com/twirapp/twir/libs/repositories/audit-logs"
	"github.com/twirapp/twir/libs/repositories/audit-logs/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	AuditLogsRepository auditlogsrepository.Repository
	LogsPubsub          auditlogs.PubSub
}

func New(opts Opts) *Service {
	return &Service{
		auditLogsRepository: opts.AuditLogsRepository,
		logsPubsub:          opts.LogsPubsub,
	}
}

type Service struct {
	auditLogsRepository auditlogsrepository.Repository
	logsPubsub          auditlogs.PubSub
}

func modelToEntity(m model.AuditLog) entity.AuditLog {
	return entity.AuditLog{
		ID:            m.ID,
		Table:         m.Table,
		OperationType: entity.AuditOperationType(m.OperationType),
		OldValue:      m.OldValue.Ptr(),
		NewValue:      m.NewValue.Ptr(),
		ObjectID:      m.ObjectID.Ptr(),
		ChannelID:     m.ChannelID.Ptr(),
		UserID:        m.UserID.Ptr(),
		CreatedAt:     m.CreatedAt,
	}
}

func (c *Service) GetMany(ctx context.Context, channelID string, limit int) (
	[]entity.AuditLog,
	error,
) {
	logs, err := c.auditLogsRepository.GetByChannelID(ctx, channelID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]entity.AuditLog, 0, len(logs))
	for _, e := range logs {
		result = append(result, modelToEntity(e))
	}

	return result, nil
}

func (c *Service) Subscribe(ctx context.Context, channelID string) (chan entity.AuditLog, error) {
	auditLogs, err := c.logsPubsub.Subscribe(ctx, channelID)
	if err != nil {
		return nil, err
	}

	channel := make(chan entity.AuditLog)

	go func() {
		defer func() {
			_ = auditLogs.Close()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case auditLog := <-auditLogs.Channel():
				channel <- entity.AuditLog{
					ID:            auditLog.ID,
					Table:         auditLog.Table,
					OperationType: entity.AuditOperationType(auditLog.OperationType),
					OldValue:      auditLog.OldValue.Ptr(),
					NewValue:      auditLog.NewValue.Ptr(),
					ObjectID:      auditLog.ObjectID.Ptr(),
					ChannelID:     auditLog.ChannelID.Ptr(),
					UserID:        auditLog.UserID.Ptr(),
					CreatedAt:     auditLog.CreatedAt,
				}
			}
		}
	}()

	return channel, nil
}
