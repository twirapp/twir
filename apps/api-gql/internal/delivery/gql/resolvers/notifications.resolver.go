package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
)

// TwitchProfile is the resolver for the twitchProfile field.
func (r *adminNotificationResolver) TwitchProfile(ctx context.Context, obj *gqlmodel.AdminNotification) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.UserID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.UserID)
}

// NotificationsCreate is the resolver for the notificationsCreate field.
func (r *mutationResolver) NotificationsCreate(ctx context.Context, text *string, editorJsJSON *string, userID *string) (*gqlmodel.AdminNotification, error) {
	entity := model.Notifications{
		ID:           uuid.NewString(),
		CreatedAt:    time.Now().UTC(),
		UserID:       null.StringFromPtr(userID),
		Message:      null.StringFromPtr(text),
		EditorJsJson: null.StringFromPtr(editorJsJSON),
	}

	if err := r.deps.Gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	userNotification := gqlmodel.UserNotification{
		ID:           entity.ID,
		Text:         entity.Message.Ptr(),
		EditorJsJSON: entity.EditorJsJson.Ptr(),
		UserID:       entity.UserID.Ptr(),
		CreatedAt:    entity.CreatedAt,
	}

	go func() {
		subKey := notificationsSubscriptionKey
		if userID != nil {
			subKey += "." + *userID
		}

		if err := r.deps.WsRouter.Publish(subKey, &userNotification); err != nil {
			r.deps.Logger.Error("failed to publish notification", slog.Any("err", err))
		}
	}()

	adminNotification := gqlmodel.AdminNotification{
		ID:           entity.ID,
		UserID:       entity.UserID.Ptr(),
		Text:         entity.Message.Ptr(),
		EditorJsJSON: entity.EditorJsJson.Ptr(),
		CreatedAt:    entity.CreatedAt,
	}

	if adminNotification.UserID != nil {
		twitchUser, err := r.deps.CachedTwitchClient.GetUserById(ctx, *adminNotification.UserID)
		if err != nil {
			return nil, err
		}

		adminNotification.TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
			Login:           twitchUser.Login,
			DisplayName:     twitchUser.DisplayName,
			ProfileImageURL: twitchUser.ProfileImageURL,
			Description:     twitchUser.Description,
		}
	}

	return &adminNotification, nil
}

// NotificationsUpdate is the resolver for the notificationsUpdate field.
func (r *mutationResolver) NotificationsUpdate(ctx context.Context, id string, opts gqlmodel.NotificationUpdateOpts) (*gqlmodel.AdminNotification, error) {
	entity := model.Notifications{}
	if err := r.deps.Gorm.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}

	if opts.Text.IsSet() {
		entity.Message = null.StringFromPtr(opts.Text.Value())
	}

	if opts.EditorJsJSON.IsSet() {
		entity.EditorJsJson = null.StringFromPtr(opts.EditorJsJSON.Value())
	}

	if err := r.deps.Gorm.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	notification := gqlmodel.AdminNotification{
		ID:           entity.ID,
		UserID:       entity.UserID.Ptr(),
		Text:         entity.Message.Ptr(),
		EditorJsJSON: entity.EditorJsJson.Ptr(),
		CreatedAt:    entity.CreatedAt,
	}

	return &notification, nil
}

// NotificationsDelete is the resolver for the notificationsDelete field.
func (r *mutationResolver) NotificationsDelete(ctx context.Context, id string) (bool, error) {
	if err := r.deps.Gorm.WithContext(ctx).Where(
		"id = ?",
		id,
	).Delete(&model.Notifications{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

// NotificationsByUser is the resolver for the notificationsByUser field.
func (r *queryResolver) NotificationsByUser(ctx context.Context) ([]gqlmodel.UserNotification, error) {
	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.Notifications
	if err := r.deps.Gorm.WithContext(ctx).Where(
		`"userId" = ? OR "userId" IS NULL`,
		user.ID,
	).Order(`"createdAt" DESC`).Find(&entities).Error; err != nil {
		return nil, err
	}

	notifications := make([]gqlmodel.UserNotification, len(entities))
	for i, entity := range entities {
		notifications[i] = gqlmodel.UserNotification{
			ID:           entity.ID,
			UserID:       entity.UserID.Ptr(),
			Text:         entity.Message.Ptr(),
			EditorJsJSON: entity.EditorJsJson.Ptr(),
			CreatedAt:    entity.CreatedAt,
		}
	}

	return notifications, nil
}

// NotificationsByAdmin is the resolver for the notificationsByAdmin field.
func (r *queryResolver) NotificationsByAdmin(ctx context.Context, opts gqlmodel.AdminNotificationsParams) (*gqlmodel.AdminNotificationsResponse, error) {
	query := r.deps.Gorm.WithContext(ctx)

	if opts.Type.IsSet() {
		switch *opts.Type.Value() {
		case gqlmodel.NotificationTypeGlobal:
			query = query.Where(`"userId" IS NULL`)
		case gqlmodel.NotificationTypeUser:
			query = query.Where(`"userId" IS NOT NULL`)
		}
	} else {
		query = query.Where(`"userId" IS NULL`)
	}

	var page int
	perPage := 20

	if opts.Page.IsSet() {
		page = *opts.Page.Value()
	}

	if opts.PerPage.IsSet() {
		perPage = *opts.PerPage.Value()
	}

	var total int64
	if err := query.Model(&model.Notifications{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var entities []model.Notifications
	if err := query.
		Limit(perPage).
		Offset(page * perPage).
		Order(`"createdAt" DESC`).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	notifications := make([]gqlmodel.AdminNotification, len(entities))
	for i, entity := range entities {
		notifications[i] = gqlmodel.AdminNotification{
			ID:           entity.ID,
			UserID:       entity.UserID.Ptr(),
			Text:         entity.Message.Ptr(),
			EditorJsJSON: entity.EditorJsJson.Ptr(),
			CreatedAt:    entity.CreatedAt,
		}
	}

	return &gqlmodel.AdminNotificationsResponse{
		Notifications: notifications,
		Total:         int(total),
	}, nil
}

// NewNotification is the resolver for the newNotification field.
func (r *subscriptionResolver) NewNotification(ctx context.Context) (<-chan *gqlmodel.UserNotification, error) {
	user, err := r.deps.Sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	channel := make(chan *gqlmodel.UserNotification, 1)

	go func() {
		sub, err := r.deps.WsRouter.Subscribe(
			[]string{
				notificationsSubscriptionKey + "." + user.ID,
				notificationsSubscriptionKey,
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var notification gqlmodel.UserNotification
				if err := json.Unmarshal(data, &notification); err != nil {
					panic(err)
				}

				channel <- &notification
			}
		}
	}()

	return channel, nil
}

// AdminNotification returns graph.AdminNotificationResolver implementation.
func (r *Resolver) AdminNotification() graph.AdminNotificationResolver {
	return &adminNotificationResolver{r}
}

type adminNotificationResolver struct{ *Resolver }
