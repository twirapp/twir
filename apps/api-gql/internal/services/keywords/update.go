package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/logger"
	keywordsrepository "github.com/twirapp/twir/libs/repositories/keywords"
)

type UpdateInput struct {
	ChannelID string
	ActorID   string

	ID               uuid.UUID
	Text             *string
	Response         *string
	Enabled          *bool
	Cooldown         *int
	CooldownExpireAt *time.Time
	IsReply          *bool
	IsRegular        *bool
	Usages           *int
	RolesIDs         []uuid.UUID
}

func (c *Service) Update(ctx context.Context, input UpdateInput) (entity.Keyword, error) {
	keyword, err := c.keywordsRepository.GetByID(ctx, input.ID)
	if err != nil {
		return entity.KeywordNil, err
	}

	if keyword.ChannelID != input.ChannelID {
		return entity.KeywordNil, ErrKeywordNotFound
	}

	newKeyword, err := c.keywordsRepository.Update(
		ctx,
		input.ID,
		keywordsrepository.UpdateInput{
			Text:             input.Text,
			Response:         input.Response,
			Enabled:          input.Enabled,
			Cooldown:         input.Cooldown,
			CooldownExpireAt: input.CooldownExpireAt,
			IsReply:          input.IsReply,
			IsRegular:        input.IsRegular,
			Usages:           input.Usages,
			RolesIDs:         &input.RolesIDs,
		},
	)
	if err != nil {
		return entity.KeywordNil, err
	}

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  lo.ToPtr(keyword.ID.String()),
			},
			NewValue: newKeyword,
			OldValue: keyword,
		},
	)

	if err = c.keywordsCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", logger.Error(err))
	}

	return c.dbToModel(newKeyword), nil
}
