package badges

import (
	"context"
	"log/slog"

	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/badges_unprotected"
	"github.com/twitchtv/twirp"
	google_protobuf "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
)

type Badges struct {
	*impl_deps.Deps
}

func (c *Badges) computeBadgeUrl(id string) string {
	if c.Config.AppEnv == "development" {
		return c.Config.S3PublicUrl + "/" + c.Config.S3Bucket + "/badges/" + id
	}

	return c.Config.S3Host + "/badges/" + id
}

func (c *Badges) GetBadgesWithUsers(
	ctx context.Context,
	_ *google_protobuf.Empty,
) (*badges_unprotected.GetBadgesResponse, error) {
	var entities []model.Badge
	if err := c.Db.
		WithContext(ctx).
		Debug().
		Preload("Users").
		Order("name DESC").
		Find(&entities).
		Error; err != nil {
		c.Logger.Error("cannot get badges", slog.Any("err", err))
		return nil, twirp.InternalError("cannot get badges")
	}

	resp := &badges_unprotected.GetBadgesResponse{
		Badges: make([]*badges_unprotected.Badge, 0, len(entities)),
		Users:  make(map[string]*structpb.ListValue),
	}

	if len(entities) == 0 {
		return resp, nil
	}

	for _, entity := range entities {
		resp.Badges = append(
			resp.Badges, &badges_unprotected.Badge{
				Id:        entity.ID.String(),
				Name:      entity.Name,
				CreatedAt: entity.CreatedAt.String(),
				FileUrl:   c.computeBadgeUrl(entity.ID.String()),
				Enabled:   entity.Enabled,
			},
		)

		badgeUsers := make([]*structpb.Value, 0, len(entity.Users))
		for _, user := range entity.Users {
			badgeUsers = append(badgeUsers, structpb.NewStringValue(user.UserID))
		}

		resp.Users[entity.ID.String()] = &structpb.ListValue{
			Values: badgeUsers,
		}
	}

	return resp, nil
}
