package notifications

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/nicklaw5/helix/v2"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
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
	req *messages_admin_notifications.GetNotificationsRequest,
) (*messages_admin_notifications.GetNotificationsResponse, error) {
	var notifications []model.Notifications

	page := req.GetPage()
	perPage := req.GetPerPage()
	if perPage == 0 {
		perPage = 50
	}

	query := c.Db.WithContext(ctx).
		Order(`"createdAt" DESC`)

	if req.GetIsUser() {
		query = query.Where(`"userId" IS NOT NULL`)
	}

	if req.GetSearch() != "" {
		query = query.Where(`message ILIKE ?`, "%"+req.GetSearch()+"%")
	}

	if err := query.
		Limit(int(perPage)).
		Offset(int(page * perPage)).
		Find(&notifications).
		Error; err != nil {
		return nil, err
	}

	twitchClient, err := twitch.NewAppClientWithContext(ctx, c.Config, c.Grpc.Tokens)
	if err != nil {
		return nil, err
	}

	twitchUsers := make([]helix.User, 0, len(notifications))
	usersIdsForRequest := make([]string, 0, len(notifications))
	for _, notification := range notifications {
		if notification.UserID.Valid {
			usersIdsForRequest = append(usersIdsForRequest, notification.UserID.String)
		}
	}

	if len(usersIdsForRequest) > 0 {
		usersReq, err := twitchClient.GetUsers(&helix.UsersParams{IDs: usersIdsForRequest})
		if err != nil {
			return nil, err
		}
		if usersReq.ErrorMessage != "" {
			return nil, fmt.Errorf(usersReq.ErrorMessage)
		}

		twitchUsers = usersReq.Data.Users
	}

	mappedNotifications := make([]*messages_admin_notifications.Notification, 0, len(notifications))
	for _, notification := range notifications {
		convertedModel := convertModelToMessage(notification)
		for _, user := range twitchUsers {
			if notification.UserID.String == user.ID {
				convertedModel.UserName = &user.Login
				convertedModel.UserDisplayName = &user.DisplayName
				convertedModel.UserAvatar = &user.ProfileImageURL
			}
		}

		mappedNotifications = append(mappedNotifications, convertedModel)
	}

	var total int64
	if err := c.Db.WithContext(ctx).Model(&model.Notifications{}).Count(&total).Error; err != nil {
		return nil, err
	}

	return &messages_admin_notifications.GetNotificationsResponse{
		Notifications: mappedNotifications,
		Total:         int32(total),
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
