package badges_with_users

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	BadgesService      *badges.Service
	BadgesUsersService *badges_users.Service
}

func New(opts Opts) *Service {
	return &Service{
		badgesService:      opts.BadgesService,
		badgesUsersService: opts.BadgesUsersService,
	}
}

type Service struct {
	badgesService      *badges.Service
	badgesUsersService *badges_users.Service
}

type GetManyInput struct {
	Enabled bool
}

func (s *Service) GetMany(ctx context.Context, input GetManyInput) (
	[]entity.BadgeWithUsers,
	error,
) {
	badgesEntities, err := s.badgesService.GetMany(ctx, badges.GetManyInput{Enabled: input.Enabled})
	if err != nil {
		return nil, err
	}

	var mu sync.Mutex
	badgesWithUsers := make([]entity.BadgeWithUsers, 0, len(badgesEntities))

	wg, wgCtx := errgroup.WithContext(ctx)
	for _, b := range badgesEntities {
		b := b

		wg.Go(
			func() error {
				badge := entity.BadgeWithUsers{
					Badge: b,
					Users: nil,
				}

				users, err := s.badgesUsersService.GetMany(
					wgCtx,
					badges_users.GetManyInput{BadgeID: b.ID},
				)
				if err != nil {
					return err
				}

				userIds := make([]string, 0, len(users))
				for _, user := range users {
					userIds = append(userIds, user.UserID)
				}

				badge.Users = userIds

				mu.Lock()
				defer mu.Unlock()
				badgesWithUsers = append(badgesWithUsers, badge)
				return nil
			},
		)
	}

	if err := wg.Wait(); err != nil {
		return nil, err
	}

	slices.SortFunc(
		badgesWithUsers,
		func(i, j entity.BadgeWithUsers) int {
			return cmp.Compare(i.FFZSlot, j.FFZSlot)
		},
	)

	return badgesWithUsers, nil
}
