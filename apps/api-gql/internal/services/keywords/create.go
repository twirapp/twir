package keywords

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/repositories/keywords"
)

type CreateInput struct {
	ChannelID string
	ActorID   string

	Text             string
	Response         string
	Enabled          bool
	Cooldown         int
	CooldownExpireAt *time.Time
	IsReply          bool
	IsRegular        bool
	Usages           int
	RolesIDs         []uuid.UUID
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Keyword, error) {
	plan, err := c.plansRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.KeywordNil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return entity.KeywordNil, fmt.Errorf("plan not found for channel")
	}

	createdCount, err := c.keywordsRepository.CountByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.KeywordNil, err
	}

	if createdCount >= plan.MaxKeywords {
		return entity.KeywordNil, fmt.Errorf("you can have only %v keywords", plan.MaxKeywords)
	}

	k, err := c.keywordsRepository.Create(
		ctx, keywords.CreateInput{
			ChannelID:        input.ChannelID,
			Text:             input.Text,
			Response:         input.Response,
			Enabled:          input.Enabled,
			Cooldown:         input.Cooldown,
			CooldownExpireAt: input.CooldownExpireAt,
			IsReply:          input.IsReply,
			IsRegular:        input.IsRegular,
			Usages:           input.Usages,
			RolesIDs:         input.RolesIDs,
		},
	)
	if err != nil {
		return entity.KeywordNil, err
	}

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(k.ID.String()),
			},
			NewValue: k,
		},
	)

	if err := c.keywordsCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", err)
	}

	return c.dbToModel(k), nil
}
