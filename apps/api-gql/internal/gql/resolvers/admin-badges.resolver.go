package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// BadgesDelete is the resolver for the badgesDelete field.
func (r *mutationResolver) BadgesDelete(ctx context.Context, id string) (bool, error) {
	badge := model.Badge{}
	if err := r.gorm.
		WithContext(ctx).
		Where("id = ?", id).
		First(&badge).
		Error; err != nil {
		return false, fmt.Errorf("cannot find badge: %w", err)
	}

	if err := r.gorm.
		WithContext(ctx).
		Delete(&badge).
		Error; err != nil {
		return false, fmt.Errorf("cannot delete badge: %w", err)
	}

	if err := r.minioClient.RemoveObject(
		ctx,
		r.config.S3Bucket,
		fmt.Sprintf("badges/%s", id),
		minio.RemoveObjectOptions{},
	); err != nil {
		fmt.Println("cannot delete file")
	}

	return true, nil
}

// BadgesUpdate is the resolver for the badgesUpdate field.
func (r *mutationResolver) BadgesUpdate(
	ctx context.Context,
	id string,
	opts gqlmodel.TwirBadgeUpdateOpts,
) (*gqlmodel.Badge, error) {
	entity := model.Badge{}
	if err := r.gorm.
		WithContext(ctx).
		Joins("Users").
		Where(
			"badges.id = ?",
			id,
		).
		First(&entity).Error; err != nil {
		return nil, fmt.Errorf("cannot find badge: %w", err)
	}

	if opts.File.IsSet() {
		file := opts.File.Value()
		_, err := r.minioClient.PutObject(
			ctx,
			r.config.S3Bucket,
			fmt.Sprintf("badges/%s", entity.ID),
			file.File,
			file.Size,
			minio.PutObjectOptions{
				ContentType: file.ContentType,
			},
		)
		if err != nil {
			return nil, fmt.Errorf("cannot upload badge file: %w", err)
		}
	}

	if opts.Name.IsSet() {
		entity.Name = *opts.Name.Value()
	}

	if opts.Enabled.IsSet() {
		entity.Enabled = *opts.Enabled.Value()
	}

	if err := r.gorm.WithContext(ctx).Save(&entity).Error; err != nil {
		return nil, err
	}

	usersIds := make([]string, 0, len(entity.Users))
	for _, user := range entity.Users {
		usersIds = append(usersIds, user.UserID)
	}

	return &gqlmodel.Badge{
		ID:        entity.ID.String(),
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.String(),
		FileURL:   r.computeBadgeUrl(entity.ID.String()),
		Enabled:   entity.Enabled,
		Users:     usersIds,
	}, nil
}

// BadgesCreate is the resolver for the badgesCreate field.
func (r *mutationResolver) BadgesCreate(
	ctx context.Context,
	name string,
	file graphql.Upload,
) (*gqlmodel.Badge, error) {
	entity := model.Badge{
		ID:        uuid.New(),
		Name:      name,
		Enabled:   true,
		CreatedAt: time.Now().UTC(),
	}

	if err := r.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	_, err := r.minioClient.PutObject(
		ctx,
		r.config.S3Bucket,
		fmt.Sprintf("badges/%s", entity.ID),
		file.File,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.ContentType,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot upload badge file: %w", err)
	}

	return &gqlmodel.Badge{
		ID:        entity.ID.String(),
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.String(),
		FileURL:   r.computeBadgeUrl(entity.ID.String()),
		Enabled:   entity.Enabled,
		Users:     []string{},
	}, nil
}

// BadgesAddUser is the resolver for the badgesAddUser field.
func (r *mutationResolver) BadgesAddUser(ctx context.Context, id string, userID string) (
	bool,
	error,
) {
	entity := model.BadgeUser{
		ID:        uuid.New(),
		BadgeID:   uuid.MustParse(id),
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
	if err := r.gorm.WithContext(ctx).Create(&entity).Error; err != nil {
		return false, err
	}

	return true, nil
}

// BadgesRemoveUser is the resolver for the badgesRemoveUser field.
func (r *mutationResolver) BadgesRemoveUser(ctx context.Context, id string, userID string) (
	bool,
	error,
) {
	entity := model.BadgeUser{}
	if err := r.gorm.WithContext(ctx).
		Where("badge_id = ? AND user_id = ?", id, userID).
		First(&entity).
		Error; err != nil {
		return false, err
	}

	if err := r.gorm.WithContext(ctx).Delete(&entity, "id = ?", entity.ID).Error; err != nil {
		return false, err
	}

	return true, nil
}

// TwirBadges is the resolver for the twirBadges field.
func (r *queryResolver) TwirBadges(ctx context.Context) ([]gqlmodel.Badge, error) {
	var entities []model.Badge
	if err := r.gorm.
		WithContext(ctx).
		Preload("Users").
		Order("name ASC").
		Find(&entities).
		Error; err != nil {
		return nil, fmt.Errorf("cannot get badges: %w", err)
	}

	respBadges := make([]gqlmodel.Badge, 0, len(entities))

	if len(entities) == 0 {
		return []gqlmodel.Badge{}, nil
	}

	for _, entity := range entities {
		badgeUsers := make([]string, 0, len(entity.Users))
		for _, user := range entity.Users {
			badgeUsers = append(badgeUsers, user.UserID)
		}

		respBadges = append(
			respBadges,
			gqlmodel.Badge{
				ID:        entity.ID.String(),
				Name:      entity.Name,
				CreatedAt: entity.CreatedAt.String(),
				FileURL:   r.computeBadgeUrl(entity.ID.String()),
				Enabled:   entity.Enabled,
				Users:     badgeUsers,
			},
		)
	}

	return respBadges, nil
}