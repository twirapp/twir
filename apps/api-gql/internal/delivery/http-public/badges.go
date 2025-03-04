package http_public

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
)

type badgesOutput struct {
	Body []badgeWithUsers `json:"body"`
}

type badgeWithUsers struct {
	Name    string    `json:"name"`
	URL     string    `json:"url"`
	Users   []string  `json:"users"`
	FFZSlot int       `json:"ffzSlot"`
	ID      uuid.UUID `json:"id"`
}

// TODO: use some gin middleware for cache response

func (p *Public) HandleBadgesGet(ctx context.Context) (*badgesOutput, error) {
	entities, err := p.badgesWithUsersService.GetMany(
		ctx,
		badges_with_users.GetManyInput{Enabled: lo.ToPtr(true)},
	)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal server error")
	}

	result := make([]badgeWithUsers, 0, len(entities))
	for _, entity := range entities {
		result = append(
			result,
			badgeWithUsers{
				ID:      entity.ID,
				Name:    entity.Name,
				FFZSlot: entity.FFZSlot,
				URL:     entity.FileURL,
				Users:   entity.Users,
			},
		)
	}

	return &badgesOutput{Body: result}, nil
}
