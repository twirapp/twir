package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
)

type UpdateInput struct {
	ChannelID string
	ActorID   string

	ID               string
	Text             *string
	Response         *string
	Enabled          *bool
	Cooldown         *int
	CooldownExpireAt *time.Time
	IsReply          *bool
	IsRegular        *bool
	Usages           *int
}

func (c *Service) Update(ctx context.Context, input UpdateInput) (entity.Keyword, error) {
	parsedID, err := uuid.Parse(input.ID)
	if err != nil {
		return entity.KeywordNil, err
	}

	keyword, err := c.keywordsRepository.GetByID(ctx, parsedID)
	if err != nil {
		return entity.KeywordNil, err
	}

	if keyword.ChannelID != input.ChannelID {
		return entity.KeywordNil, ErrKeywordNotFound
	}

	newKeyword, err := c.keywordsRepository.Update(
		ctx,
		parsedID,
		keywordsrepository.UpdateInput{
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
		"Keywords update",
		audit.Fields{
			OldValue:      keyword,
			NewValue:      newKeyword,
			ActorID:       &input.ActorID,
			ChannelID:     &input.ChannelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(keyword.ID.String()),
		},
	)

	if err := c.keywordsCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", err)
	}

	return entity.KeywordNil, nil
}
