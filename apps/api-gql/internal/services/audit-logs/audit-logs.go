package audit_logs

import (
	"context"

	auditlogs "github.com/satont/twir/libs/pubsub/audit-logs"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	auditlogsrepository "github.com/twirapp/twir/libs/repositories/audit_logs"
	"github.com/twirapp/twir/libs/repositories/audit_logs/model"
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
		TableName:     m.TableName,
		OperationType: entity.AuditOperationType(m.OperationType),
		OldValue:      m.OldValue.Ptr(),
		NewValue:      m.NewValue.Ptr(),
		ObjectID:      m.ObjectID.Ptr(),
		ChannelID:     m.ChannelID.Ptr(),
		UserID:        m.UserID.Ptr(),
		CreatedAt:     m.CreatedAt,
	}
}

type GetManyInput struct {
	ChannelID      *string
	ActorID        *string
	ObjectID       *string
	Limit          int
	Page           int
	Systems        []string
	OperationTypes []entity.AuditOperationType
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) (
	[]entity.AuditLog,
	error,
) {
	operationTypes := make([]model.AuditOperationType, 0, len(input.OperationTypes))
	for _, t := range input.OperationTypes {
		operationTypes = append(operationTypes, model.AuditOperationType(t))
	}

	logs, err := c.auditLogsRepository.GetMany(
		ctx,
		auditlogsrepository.GetManyInput{
			ChannelID:      input.ChannelID,
			ActorID:        input.ActorID,
			ObjectID:       input.ObjectID,
			Limit:          input.Limit,
			Page:           input.Page,
			Systems:        input.Systems,
			OperationTypes: operationTypes,
		},
	)
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
					TableName:     auditLog.Table,
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

type GetCountInput struct {
	ChannelID *string
}

func (c *Service) Count(ctx context.Context, input GetCountInput) (uint64, error) {
	return c.auditLogsRepository.Count(
		ctx,
		auditlogsrepository.GetCountInput{
			ChannelID: input.ChannelID,
		},
	)
}
