package keywords

import (
	"context"
	"fmt"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
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
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Keyword, error) {
	createdCount, err := c.keywordsRepository.CountByChannelID(ctx, input.ChannelID)
	if err != nil {
		return entity.KeywordNil, err
	}

	if createdCount >= MaxPerChannel {
		return entity.KeywordNil, fmt.Errorf("you can have only %v keywords", MaxPerChannel)
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
		},
	)
	if err != nil {
		return entity.KeywordNil, err
	}

	c.logger.Audit(
		"Keywords create",
		audit.Fields{
			NewValue:      k,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
			OperationType: audit.OperationCreate,
			ObjectID:      lo.ToPtr(k.ID.String()),
		},
	)

	if err := c.keywordsCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", err)
	}

	return c.dbToModel(k), nil
}
