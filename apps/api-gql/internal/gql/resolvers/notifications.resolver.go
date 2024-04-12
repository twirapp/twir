package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	twitchcahe "github.com/twirapp/twir/libs/cache/twitch"
)

// NotificationsCreate is the resolver for the notificationsCreate field.
func (r *mutationResolver) NotificationsCreate(ctx context.Context, text string, userID *string) (*gqlmodel.AdminNotification, error) {
	entity := model.Notifications{
		ID:        uuid.NewString(),
		CreatedAt: time.Now().UTC(),
		UserID:    null.StringFromPtr(userID),
		Message:   text,
	}

	if err := r.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	userNotification := gqlmodel.UserNotification{
		ID:        entity.ID,
		Text:      entity.Message,
		UserID:    entity.UserID.Ptr(),
		CreatedAt: entity.CreatedAt,
	}

	if userID == nil {
		for _, channel := range r.subscriptionsStore.NewNotificationsChannels {
			channel <- &userNotification
		}
	} else {
		if r.subscriptionsStore.NewNotificationsChannels[*userID] != nil {
			r.subscriptionsStore.NewNotificationsChannels[*userID] <- &userNotification
		}
	}

	adminNotification := gqlmodel.AdminNotification{
		ID:        entity.ID,
		UserID:    entity.UserID.Ptr(),
		Text:      entity.Message,
		CreatedAt: entity.CreatedAt,
	}

	if adminNotification.UserID != nil {
		twitchUser, err := r.cachedTwitchClient.GetUserById(ctx, *adminNotification.UserID)
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
	if err := r.gorm.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
		return nil, err
	}

	if opts.Text.IsSet() {
		entity.Message = *opts.Text.Value()
	}

	if err := r.gorm.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	notification := gqlmodel.AdminNotification{
		ID:        entity.ID,
		UserID:    entity.UserID.Ptr(),
		Text:      entity.Message,
		CreatedAt: entity.CreatedAt,
	}

	if notification.UserID != nil {
		twitchUser, err := r.cachedTwitchClient.GetUserById(ctx, *notification.UserID)
		if err != nil {
			return nil, err
		}

		notification.TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
			Login:           twitchUser.Login,
			DisplayName:     twitchUser.DisplayName,
			ProfileImageURL: twitchUser.ProfileImageURL,
			Description:     twitchUser.Description,
		}
	}

	return &notification, nil
}

// NotificationsDelete is the resolver for the notificationsDelete field.
func (r *mutationResolver) NotificationsDelete(ctx context.Context, id string) (bool, error) {
	if err := r.gorm.WithContext(ctx).Where(
		"id = ?",
		id,
	).Delete(&model.Notifications{}).Error; err != nil {
		return false, err
	}

	return true, nil
}

// NotificationsByUser is the resolver for the notificationsByUser field.
func (r *queryResolver) NotificationsByUser(ctx context.Context) ([]gqlmodel.UserNotification, error) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	var entities []model.Notifications
	if err := r.gorm.WithContext(ctx).Where(
		`"userId" = ? OR "userId" IS NULL`,
		user.ID,
	).Order(`"createdAt" DESC`).Find(&entities).Error; err != nil {
		return nil, err
	}

	notifications := make([]gqlmodel.UserNotification, len(entities))
	for i, entity := range entities {
		notifications[i] = gqlmodel.UserNotification{
			ID:        entity.ID,
			UserID:    entity.UserID.Ptr(),
			Text:      entity.Message,
			CreatedAt: entity.CreatedAt,
		}
	}

	return notifications, nil
}

// NotificationsByAdmin is the resolver for the notificationsByAdmin field.
func (r *queryResolver) NotificationsByAdmin(ctx context.Context, opts gqlmodel.AdminNotificationsParams) (*gqlmodel.AdminNotificationsResponse, error) {
	query := r.gorm.WithContext(ctx)

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
		Find(&entities).Error; err != nil {
		return nil, err
	}

	notifications := make([]gqlmodel.AdminNotification, len(entities))
	for i, entity := range entities {
		notifications[i] = gqlmodel.AdminNotification{
			ID:        entity.ID,
			UserID:    entity.UserID.Ptr(),
			Text:      entity.Message,
			CreatedAt: entity.CreatedAt,
		}
	}

	needTwitch := slices.Contains(
		GetPreloads(ctx),
		"notifications.twitchProfile",
	)

	if needTwitch && opts.Type.IsSet() && *opts.Type.Value() == gqlmodel.NotificationTypeUser {
		usersIdsForRequest := make([]string, len(notifications))
		for i, notification := range notifications {
			usersIdsForRequest[i] = *notification.UserID
		}

		twitchUsers, err := r.cachedTwitchClient.GetUsersByIds(ctx, usersIdsForRequest)
		if err != nil {
			return nil, err
		}

		for i, notification := range notifications {
			notificationTwitchUser, ok := lo.Find(
				twitchUsers, func(item twitchcahe.TwitchUser) bool {
					return item.ID == *notification.UserID
				},
			)
			if !ok {
				continue
			}

			notifications[i].TwitchProfile = &gqlmodel.TwirUserTwitchInfo{
				Login:           notificationTwitchUser.Login,
				DisplayName:     notificationTwitchUser.DisplayName,
				ProfileImageURL: notificationTwitchUser.ProfileImageURL,
				Description:     notificationTwitchUser.Description,
			}
		}
	}

	return &gqlmodel.AdminNotificationsResponse{
		Notifications: notifications,
		Total:         int(total),
	}, nil
}

// NewNotification is the resolver for the newNotification field.
func (r *subscriptionResolver) NewNotification(ctx context.Context) (<-chan *gqlmodel.UserNotification, error) {
	user, err := r.sessions.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	channel := make(chan *gqlmodel.UserNotification, 1)
	if r.subscriptionsStore.NewNotificationsChannels[user.ID] == nil {
		r.subscriptionsStore.NewNotificationsChannels[user.ID] = channel
	}

	go func() {
		defer close(channel)

		for {
			select {
			case <-ctx.Done():
				// close(r.subscriptionsStore.NewNotificationsChannels[user.ID])
				delete(r.subscriptionsStore.NewNotificationsChannels, user.ID)
				return
			case notification := <-r.subscriptionsStore.NewNotificationsChannels[user.ID]:
				channel <- notification
			}
		}
	}()

	return channel, nil
}
