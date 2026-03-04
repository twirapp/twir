package badges

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/errors"
	"github.com/twirapp/twir/libs/repositories/badges"
	"github.com/twirapp/twir/libs/repositories/badges/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	BadgesRepository badges.Repository
	Config           config.Config
	MinioClient      *minio.Client
}

func New(opts Opts) *Service {
	return &Service{
		badgesRepository: opts.BadgesRepository,
		config:           opts.Config,
		minioClient:      opts.MinioClient,
	}
}

type Service struct {
	badgesRepository badges.Repository
	config           config.Config
	minioClient      *minio.Client
}

type GetManyInput struct {
	Enabled *bool
}

func (c *Service) modelToEntity(b model.Badge) entity.Badge {
	return entity.Badge{
		ID:        b.ID,
		Name:      b.Name,
		Enabled:   b.Enabled,
		CreatedAt: b.CreatedAt,
		FileName:  b.FileName,
		FFZSlot:   b.FFZSlot,
		FileURL:   c.computeBadgeUrl(b.FileName),
	}
}

func (c *Service) computeBadgeFileName(file entity.Upload, fileID uuid.UUID) (string, error) {
	fileExtension := filepath.Ext(file.Filename)
	if fileExtension == "" {
		return "", errors.NewBadRequestError("File must have a valid extension")
	}
	if !strings.HasPrefix(file.ContentType, "image/") {
		return "", errors.NewBadRequestError("Only image files are allowed")
	}

	fileExtension = strings.ToLower(fileExtension)
	fileName := fmt.Sprintf("%s%s", fileID, fileExtension)

	return fileName, nil
}

func (c *Service) computeBadgeUrl(fileName string) string {
	if c.config.AppEnv == "development" {
		return c.config.S3PublicUrl + "/" + c.config.S3Bucket + "/badges/" + fileName
	}

	return c.config.S3PublicUrl + "/badges/" + fileName
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.Badge, error) {
	b, err := c.badgesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.BadgeNil, errors.NewInternalError("Failed to fetch badge", err)
	}

	return c.modelToEntity(b), nil
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.Badge, error) {
	selectedBadges, err := c.badgesRepository.GetMany(
		ctx,
		badges.GetManyInput{
			Enabled: input.Enabled,
		},
	)
	if err != nil {
		return nil, errors.NewInternalError("Failed to fetch badges", err)
	}

	result := make([]entity.Badge, 0, len(selectedBadges))
	for _, b := range selectedBadges {
		result = append(result, c.modelToEntity(b))
	}

	return result, nil
}

func (c *Service) Delete(ctx context.Context, id uuid.UUID) error {
	b, err := c.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := c.badgesRepository.Delete(ctx, id); err != nil {
		return errors.NewInternalError("Failed to delete badge", err)
	}

	if err := c.minioClient.RemoveObject(
		ctx,
		c.config.S3Bucket,
		"badges/"+b.FileName,
		minio.RemoveObjectOptions{},
	); err != nil {
		return errors.NewInternalError("Failed to delete badge file", err)
	}

	return nil
}

type CreateInput struct {
	Name    string
	Enabled bool
	FfzSlot int
	File    entity.Upload
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Badge, error) {
	fileId := uuid.New()
	fileName, err := c.computeBadgeFileName(input.File, fileId)
	if err != nil {
		return entity.BadgeNil, err
	}

	newBadge, err := c.badgesRepository.Create(
		ctx, badges.CreateInput{
			ID:       fileId,
			Name:     input.Name,
			Enabled:  input.Enabled,
			FFZSlot:  input.FfzSlot,
			FileName: fileName,
		},
	)
	if err != nil {
		return entity.BadgeNil, errors.NewInternalError("Failed to create badge in database", err)
	}

	_, err = c.minioClient.PutObject(
		ctx,
		c.config.S3Bucket,
		fmt.Sprintf("badges/%s", fileName),
		input.File.File,
		input.File.Size,
		minio.PutObjectOptions{
			ContentType: input.File.ContentType,
		},
	)
	if err != nil {
		return entity.BadgeNil, errors.NewInternalError("Failed to upload badge file", err)
	}

	return c.modelToEntity(newBadge), nil
}

type UpdateInput struct {
	Name    *string
	Enabled *bool
	FfzSlot *int
	File    *entity.Upload
}

func (c *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (
	entity.Badge,
	error,
) {
	b, err := c.badgesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.BadgeNil, errors.NewInternalError("Failed to fetch badge", err)
	}

	updateInput := badges.UpdateInput{
		Name:     input.Name,
		Enabled:  input.Enabled,
		FFZSlot:  input.FfzSlot,
		FileName: nil,
	}

	if input.File != nil {
		file := *input.File
		fileName, err := c.computeBadgeFileName(file, b.ID)
		if err != nil {
			return entity.BadgeNil, err
		}

		if b.FileName != fileName {
			if err := c.minioClient.RemoveObject(
				ctx,
				c.config.S3Bucket,
				fmt.Sprintf("badges/%s", b.FileName),
				minio.RemoveObjectOptions{},
			); err != nil {
				return entity.BadgeNil, errors.NewInternalError("Failed to delete old badge file", err)
			}
		}

		_, err = c.minioClient.PutObject(
			ctx,
			c.config.S3Bucket,
			fmt.Sprintf("badges/%s", fileName),
			file.File,
			file.Size,
			minio.PutObjectOptions{
				ContentType: file.ContentType,
			},
		)
		if err != nil {
			return entity.BadgeNil, errors.NewInternalError("Failed to upload badge file", err)
		}

		updateInput.FileName = &fileName
	}

	newBadge, err := c.badgesRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.BadgeNil, errors.NewInternalError("Failed to update badge in database", err)
	}

	return c.modelToEntity(newBadge), nil
}
