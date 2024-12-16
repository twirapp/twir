package badges

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
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
	Enabled bool
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
		return "", fmt.Errorf("file extension is empty")
	}
	if !strings.HasPrefix(file.ContentType, "image/") {
		return "", fmt.Errorf("file is not an image")
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
		return entity.BadgeNil, err
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
		return nil, err
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
		return err
	}

	if err := c.minioClient.RemoveObject(
		ctx,
		c.config.S3Bucket,
		"badges/"+b.FileName,
		minio.RemoveObjectOptions{},
	); err != nil {
		return err
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
		return entity.BadgeNil, fmt.Errorf("cannot create badge in db: %w", err)
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
		return entity.BadgeNil, fmt.Errorf("cannot upload badge file: %w", err)
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
		return entity.BadgeNil, err
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
			return entity.BadgeNil, fmt.Errorf("cannot compute badge file name: %w", err)
		}

		if b.FileName != fileName {
			if err := c.minioClient.RemoveObject(
				ctx,
				c.config.S3Bucket,
				fmt.Sprintf("badges/%s", b.FileName),
				minio.RemoveObjectOptions{},
			); err != nil {
				return entity.BadgeNil, fmt.Errorf("cannot delete old badge file: %w", err)
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
			return entity.BadgeNil, fmt.Errorf("cannot upload badge file: %w", err)
		}

		updateInput.FileName = &fileName
	}

	newBadge, err := c.badgesRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.BadgeNil, fmt.Errorf("cannot update badge in db: %w", err)
	}

	return c.modelToEntity(newBadge), nil
}
