package keywords

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/bots/internal/entity"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/repositories/keywords"
	"github.com/twirapp/twir/libs/repositories/keywords/model"
	"github.com/twirapp/twir/libs/repositories/roles"
	rolesmodel "github.com/twirapp/twir/libs/repositories/roles/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	KeywordsRepository keywords.Repository
	KeywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
	RolesCache         *generic_cacher.GenericCacher[[]rolesmodel.Role]
	RolesRepository    roles.Repository
	Redis              *redis.Client
}

func New(opts Opts) *Service {
	return &Service{
		keywordsRepository: opts.KeywordsRepository,
		keywordsCacher:     opts.KeywordsCacher,
		rolesRepository:    opts.RolesRepository,
		rolesCache:         opts.RolesCache,
		redis:              opts.Redis,
	}
}

type Service struct {
	keywordsRepository keywords.Repository
	keywordsCacher     *generic_cacher.GenericCacher[[]model.Keyword]
	rolesRepository    roles.Repository
	rolesCache         *generic_cacher.GenericCacher[[]rolesmodel.Role]
	redis              *redis.Client
}

func (c *Service) mapToEntity(m model.Keyword) entity.Keyword {
	return entity.Keyword{
		ID:               m.ID,
		ChannelID:        m.ChannelID.String(),
		Text:             m.Text,
		Response:         m.Response,
		Enabled:          m.Enabled,
		Cooldown:         m.Cooldown,
		CooldownExpireAt: m.CooldownExpireAt,
		IsReply:          m.IsReply,
		IsRegular:        m.IsRegular,
		Usages:           m.Usages,
		RolesIDs:         m.RolesIDs,
		Platforms:        m.Platforms,
	}
}

func (c *Service) GetChannelRoles(ctx context.Context, channelID string) (
	[]rolesmodel.Role,
	error,
) {
	return c.rolesCache.Get(ctx, channelID)
}

func (c *Service) GetUserAccessibleRoles(
	ctx context.Context,
	channelID, userID string,
) ([]rolesmodel.Role, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return c.rolesRepository.GetUserAccessibleRoles(ctx, parsedChannelID, parsedUserID)
}

func (c *Service) GetManyByChannelID(ctx context.Context, channelID string) (
	[]entity.Keyword,
	error,
) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	models, err := c.keywordsCacher.Get(ctx, parsedChannelID.String())
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
	RolesIDs         *[]uuid.UUID
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, channelID string, input UpdateInput) (
	entity.Keyword,
	error,
) {
	currentKeyword, err := c.keywordsRepository.GetByID(ctx, id)
	if err != nil {
		return entity.Keyword{}, err
	}

	updateInput := keywords.UpdateInput{
		Text:             input.Text,
		Response:         input.Response,
		Enabled:          input.Enabled,
		Cooldown:         input.Cooldown,
		CooldownExpireAt: input.CooldownExpireAt,
		IsReply:          input.IsReply,
		IsRegular:        input.IsRegular,
		Usages:           input.Usages,
		RolesIDs:         input.RolesIDs,
		Platforms:        currentKeyword.Platforms,
	}

	m, err := c.keywordsRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.Keyword{}, err
	}

	err = c.keywordsCacher.SetValueFiltered(
		ctx,
		channelID,
		func(data []model.Keyword) []model.Keyword {
			for i, v := range data {
				if v.ID == id {
					data[i] = m
					break
				}
			}

			return data
		},
	)
	if err != nil {
		return entity.Keyword{}, err
	}

	return c.mapToEntity(m), nil
}
