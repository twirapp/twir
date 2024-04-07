package notifications

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	messages_admin_notifications "github.com/twirapp/twir/libs/api/messages/admin_notifications"
	google_protobuf "google.golang.org/protobuf/types/known/emptypb"
)

type Notifications struct {
	*impl_deps.Deps
}

func convertModelToMessage(notification model.Notifications) *messages_admin_notifications.Notification {
	return &messages_admin_notifications.Notification{
		Id:        notification.ID,
		CreatedAt: notification.CreatedAt.UTC().String(),
		UserId:    notification.UserID.Ptr(),
		Message:   notification.Message,
	}
}

func (c *Notifications) NotificationsCreate(
	ctx context.Context,
	req *messages_admin_notifications.CreateNotificationRequest,
) (*messages_admin_notifications.Notification, error) {
	notification := model.Notifications{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UserID:    null.StringFromPtr(req.UserId),
		Message:   req.Message,
	}

	if err := c.Db.WithContext(ctx).Create(&notification).Error; err != nil {
		return nil, err
	}

	return convertModelToMessage(notification), nil
}

func (c *Notifications) NotificationsGetAll(
	ctx context.Context,
	_ *google_protobuf.Empty,
) (*messages_admin_notifications.GetNotificationsResponse, error) {
	var notifications []model.Notifications

	if err := c.Db.
		WithContext(ctx).
		Order(`"createdAt" DESC`).
		Find(&notifications).
		Error; err != nil {
		return nil, err
	}

	mappedNotifications := make([]*messages_admin_notifications.Notification, 0, len(notifications))
	for _, notification := range notifications {
		mappedNotifications = append(mappedNotifications, convertModelToMessage(notification))
	}

	return &messages_admin_notifications.GetNotificationsResponse{
		Notifications: mappedNotifications,
	}, nil
}
func (c *Notifications) NotificationsUpdate(
	ctx context.Context,
	req *messages_admin_notifications.UpdateNotificationRequest,
) (*messages_admin_notifications.Notification, error) {
	notification := model.Notifications{}
	if err := c.Db.WithContext(ctx).Where("id = ?", req.Id).First(&notification).Error; err != nil {
		return nil, err
	}

	notification.Message = req.Message

	if err := c.Db.WithContext(ctx).Save(&notification).Error; err != nil {
		return nil, err
	}

	return convertModelToMessage(notification), nil
}

func (c *Notifications) NotificationsDelete(
	ctx context.Context,
	req *messages_admin_notifications.DeleteNotificationRequest,
) (
	*google_protobuf.Empty,
	error,
) {
	if err := c.Db.WithContext(ctx).Where(
		"id = ?",
		req.Id,
	).Delete(&model.Notifications{}).Error; err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, nil
}