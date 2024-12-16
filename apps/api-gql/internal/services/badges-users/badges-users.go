package badges_users

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	badges_users "github.com/twirapp/twir/libs/repositories/badges-users"
	"github.com/twirapp/twir/libs/repositories/badges-users/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	BadgesUsersRepository badges_users.Repository
	Redis                 *redis.Client
}

func New(opts Opts) *Service {
	return &Service{
		badgesUsersRepository: opts.BadgesUsersRepository,
		redis:                 opts.Redis,
	}
}

type Service struct {
	badgesUsersRepository badges_users.Repository
	redis                 *redis.Client
}

type GetManyInput struct {
	BadgeID uuid.UUID
}

func modelToEntity(b model.BadgeUser) entity.BadgeUser {
	return entity.BadgeUser{
		ID:        b.ID,
		BadgeID:   b.BadgeID,
		UserID:    b.UserID,
		CreatedAt: b.CreatedAt,
	}
}

const badgeUsersCacheKey = "cache:twir:badges_users:"

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.BadgeUser, error) {
	cachedUsers, _ := c.redis.Get(ctx, badgeUsersCacheKey+input.BadgeID.String()).Bytes()
	if len(cachedUsers) > 0 {
		var result []entity.BadgeUser
		if err := json.Unmarshal(cachedUsers, &result); err != nil {
			return nil, err
		}

		return result, nil
	}

	selectedBadgeUsers, err := c.badgesUsersRepository.GetMany(
		ctx,
		badges_users.GetManyInput{
			BadgeID: input.BadgeID,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.BadgeUser, 0, len(selectedBadgeUsers))
	for _, b := range selectedBadgeUsers {
		result = append(result, modelToEntity(b))
	}

	if len(result) > 0 {
		cacheBytes, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		if err := c.redis.Set(
			ctx,
			badgeUsersCacheKey+input.BadgeID.String(),
			cacheBytes,
			1*time.Minute,
		).Err(); err != nil {
			return nil, err
		}
	}

	return result, nil
}
