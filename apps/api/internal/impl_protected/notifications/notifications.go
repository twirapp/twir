package notifications

import (
	"context"

	"github.com/satont/twir/apps/api/internal/helpers"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	"github.com/twirapp/twir/libs/api/messages/notifications"
	google_protobuf "google.golang.org/protobuf/types/known/emptypb"
)

type Notifications struct {
	*impl_deps.Deps
}

func (c *Notifications) NotificationsGetAll(
	ctx context.Context,
	_ *google_protobuf.Empty,
) (*notifications.NotificationsGetAllResponse, error) {
	user, err := helpers.GetUserModelFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	entities := []*notifications.Notification{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"userId" = ? OR "userId" IS NULL`, user.ID).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}

	convertedNotifications := make([]*notifications.Notification, 0, len(entities))
	for _, entity := range entities {
		convertedNotifications = append(
			convertedNotifications,
			&notifications.Notification{
				Id:      entity.Id,
				Message: entity.Message,
			},
		)
	}

	return &notifications.NotificationsGetAllResponse{
		Notifications: convertedNotifications,
	}, nil
}
