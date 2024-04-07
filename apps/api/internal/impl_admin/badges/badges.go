package badges

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/admin_badges"
	"github.com/twirapp/twir/libs/api/messages/badges_unprotected"
	"github.com/twitchtv/twirp"
	google_protobuf "google.golang.org/protobuf/types/known/emptypb"
)

type Badges struct {
	*impl_deps.Deps
	s3Client *minio.Client
}

func NewBadges(deps *impl_deps.Deps) *Badges {
	instance := &Badges{
		Deps: deps,
	}

	if deps.Config.S3AccessToken != "" && deps.Config.S3SecretToken != "" {
		client, err := minio.New(
			deps.Config.S3Host,
			&minio.Options{
				Creds:  credentials.NewStaticV4(deps.Config.S3AccessToken, deps.Config.S3SecretToken, ""),
				Region: deps.Config.S3Region,
				Secure: deps.Config.AppEnv == "production",
			},
		)
		if err != nil {
			deps.Logger.Error("cannot create minio host", "err", err)
		}

		instance.s3Client = client
	} else {
		deps.Logger.Warn("S3 secrets not settuped, badges wont work")
	}

	return instance
}

func (c *Badges) computeBadgeUrl(id string) string {
	if c.Config.AppEnv == "development" {
		return c.Config.S3PublicUrl + "/" + c.Config.S3Bucket + "/badges/" + id
	}

	return c.Config.S3Host + "/badges/" + id
}

func (c *Badges) BadgesCreate(
	ctx context.Context,
	req *admin_badges.CreateBadgeRequest,
) (*badges_unprotected.Badge, error) {
	if c.s3Client == nil {
		return nil, twirp.InternalError("S3 not configured, badges wont work")
	}

	if len(req.FileBytes) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "file is required")
	}

	if len(req.Name) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "name is required")
	}

	if len(req.FileMimeType) == 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "file mime type is required")
	}

	badgeId := uuid.New()

	_, err := c.s3Client.PutObject(
		ctx,
		c.Config.S3Bucket,
		fmt.Sprintf("badges/%s", badgeId),
		bytes.NewReader(req.FileBytes),
		int64(len(req.FileBytes)),
		minio.PutObjectOptions{
			ContentType: req.FileMimeType,
		},
	)

	if err != nil {
		c.Logger.Error("cannot upload badge file", slog.Any("err", err))
		return nil, fmt.Errorf("cannot upload badge: %w", err)
	}

	entity := model.Badge{
		ID:        badgeId,
		Name:      req.Name,
		Enabled:   req.Enabled,
		CreatedAt: time.Now().UTC(),
	}
	if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	return &badges_unprotected.Badge{
		Id:        entity.ID.String(),
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.UTC().String(),
		FileUrl:   c.computeBadgeUrl(entity.ID.String()),
		Enabled:   entity.Enabled,
	}, nil
}

func (c *Badges) BadgesUpdate(
	ctx context.Context,
	req *admin_badges.UpdateBadgeRequest,
) (*badges_unprotected.Badge, error) {
	if c.s3Client == nil {
		return nil, twirp.InternalError("S3 not configured, badges wont work")
	}

	entity := model.Badge{}
	if err := c.Db.WithContext(ctx).Where("id = ?", req.Id).First(&entity).Error; err != nil {
		return nil, err
	}

	if len(req.FileBytes) > 0 {
		if req.FileMimeType == nil || len(*req.FileMimeType) == 0 {
			return nil, twirp.NewError(twirp.InvalidArgument, "file mime type is required")
		}
		_, err := c.s3Client.PutObject(
			ctx,
			c.Config.S3Bucket,
			fmt.Sprintf("badges/%s", entity.ID),
			bytes.NewReader(req.FileBytes),
			int64(len(req.FileBytes)),
			minio.PutObjectOptions{
				ContentType: *req.FileMimeType,
			},
		)
		if err != nil {
			c.Logger.Error("cannot upload badge file", slog.Any("err", err))
			return nil, fmt.Errorf("cannot delete badge: %w", err)
		}
	}

	if req.Name != nil {
		entity.Name = *req.Name
	}

	if req.Enabled != nil {
		entity.Enabled = *req.Enabled
	}

	return &badges_unprotected.Badge{
		Id:        entity.ID.String(),
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt.UTC().String(),
		FileUrl:   c.computeBadgeUrl(entity.ID.String()),
		Enabled:   entity.Enabled,
	}, nil
}

func (c *Badges) BadgesDelete(
	ctx context.Context,
	req *admin_badges.DeleteBadgeRequest,
) (*google_protobuf.Empty, error) {
	if c.s3Client == nil {
		return nil, twirp.InternalError("S3 not configured, badges wont work")
	}

	entity := model.Badge{}
	if err := c.Db.WithContext(ctx).Where("id = ?", req.Id).First(&entity).Error; err != nil {
		return nil, err
	}

	err := c.s3Client.RemoveObject(
		ctx,
		c.Config.S3Bucket,
		fmt.Sprintf("badges/%s", req.Id),
		minio.RemoveObjectOptions{},
	)
	if err != nil {
		c.Logger.Error("cannot delete badge file", slog.Any("err", err))
		return nil, fmt.Errorf("cannot delete badge: %w", err)
	}

	if err := c.Db.WithContext(ctx).Delete(&entity).Error; err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, nil
}

func (c *Badges) BadgeAddUser(
	ctx context.Context,
	req *admin_badges.AddUserRequest,
) (
	*google_protobuf.Empty,
	error,
) {
	entity := model.BadgeUser{
		ID:        uuid.New(),
		BadgeID:   uuid.MustParse(req.BadgeId),
		UserID:    req.UserId,
		CreatedAt: time.Now().UTC(),
	}
	if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, nil
}

func (c *Badges) BadgeDeleteUser(
	ctx context.Context,
	req *admin_badges.DeleteUserRequest,
) (
	*google_protobuf.Empty,
	error,
) {
	entity := model.BadgeUser{}
	if err := c.Db.WithContext(ctx).
		Debug().
		Where("badge_id = ? AND user_id = ?", req.BadgeId, req.UserId).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	if err := c.Db.WithContext(ctx).Delete(&entity, "id = ?", entity.ID).Error; err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, nil
}
