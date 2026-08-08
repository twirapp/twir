package badges_with_users

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/badges"
	badges_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-users"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"golang.org/x/sync/errgroup"
)

func New(
	badgesService *badges.Service,
	badgesUsersService *badges_users.Service,
	usersRepository usersrepository.Repository,
) *Service {
	return &Service{
		badgesService:      badgesService,
		badgesUsersService: badgesUsersService,
		usersRepository:    usersRepository,
	}
}

type Service struct {
	badgesService      *badges.Service
	badgesUsersService *badges_users.Service
	usersRepository    usersrepository.Repository
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

				userIDs := make([]uuid.UUID, 0, len(users))
				for _, user := range users {
					userIDs = append(userIDs, user.UserID)
				}

				if len(userIDs) > 0 {
					dbUsers, err := s.usersRepository.GetManyByIDS(
						wgCtx,
						usersrepository.GetManyInput{
							IDs:     userIDs,
							PerPage: len(userIDs),
						},
					)
					if err != nil {
						return err
					}

					for _, dbUser := range dbUsers {
						if dbUser.Platform != platformentity.PlatformTwitch {
							continue
						}

						badge.Users = append(badge.Users, dbUser.PlatformID)
					}
				}

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
