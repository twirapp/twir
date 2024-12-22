package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/bots/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	KeywordsRepository keywords.Repository
	KeywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
}

func New(opts Opts) *Service {
	return &Service{
		keywordsRepository: opts.KeywordsRepository,
		keywordsCacher:     opts.KeywordsCacher,
	}
}

type Service struct {
	keywordsRepository keywords.Repository
	keywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
}

func (c *Service) mapToEntity(m model.Keyword) entity.Keyword {
	return entity.Keyword{
		ID:               m.ID,
		ChannelID:        m.ChannelID,
		Text:             m.Text,
		Response:         m.Response,
		Enabled:          m.Enabled,
		Cooldown:         m.Cooldown,
		CooldownExpireAt: m.CooldownExpireAt,
		IsReply:          m.IsReply,
		IsRegular:        m.IsRegular,
		Usages:           m.Usages,
	}
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.Keyword,
	error,
) {
	models, err := c.keywordsCacher.Get(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.Keyword, 0, len(models))
	for _, m := range models {
		entities = append(entities, c.mapToEntity(m))
	}

	return entities, nil
}

type UpdateInput struct {
	Text             *string
	Response         *string
	Enabled          *bool
	Cooldown         *int
	CooldownExpireAt *time.Time
	IsReply          *bool
	IsRegular        *bool
	Usages           *int
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, channelID string, input UpdateInput) (
	entity.Keyword,
	error,
) {
	updateInput := keywords.UpdateInput{
		Text:             input.Text,
		Response:         input.Response,
		Enabled:          input.Enabled,
		Cooldown:         input.Cooldown,
		CooldownExpireAt: input.CooldownExpireAt,
		IsReply:          input.IsReply,
		IsRegular:        input.IsRegular,
		Usages:           input.Usages,
	}

	m, err := c.keywordsRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.Keyword{}, err
	}

	if err := c.keywordsCacher.Invalidate(ctx, channelID); err != nil {
		return entity.Keyword{}, err
	}

	return c.mapToEntity(m), nil
}
