package badges_with_users

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	BadgesService      *badges.Service
	BadgesUsersService *badges_users.Service
	UsersRepository    users.Repository
}

func New(opts Opts) *Service {
	return &Service{
		badgesService:      opts.BadgesService,
		badgesUsersService: opts.BadgesUsersService,
		usersRepository:    opts.UsersRepository,
	}
}

type Service struct {
	badgesService      *badges.Service
	badgesUsersService *badges_users.Service
	usersRepository    users.Repository
}

type GetManyInput struct {
	Enabled *bool
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

				userIds := make([]uuid.UUID, 0, len(users))
				for _, user := range users {
					userIds = append(userIds, user.UserID)
				}

				twitchIDs, err := s.resolveTwitchPlatformIDs(wgCtx, userIds)
				if err != nil {
					return err
				}

				badge.Users = twitchIDs

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

// resolveTwitchPlatformIDs maps internal user IDs to their Twitch platform IDs.
// Users without a Twitch account are skipped.
func (s *Service) resolveTwitchPlatformIDs(ctx context.Context, userIDs []uuid.UUID) (
	[]string,
	error,
) {
	if len(userIDs) == 0 {
		return nil, nil
	}

	twitchUsers, err := s.usersRepository.GetManyByIDS(
		ctx,
		users.GetManyInput{
			IDs:      userIDs,
			Platform: lo.ToPtr(platform.PlatformTwitch),
			PerPage:  len(userIDs),
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(twitchUsers))
	for _, user := range twitchUsers {
		result = append(result, user.PlatformID)
	}

	return result, nil
}
