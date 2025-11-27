package channels_files

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/channels_files"
	"github.com/twirapp/twir/libs/repositories/channels_files/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	FilesRepo  channels_files.Repository
	TrmManager trm.Manager
	S3Client   *minio.Client
	Config     config.Config
}

type cachedFile struct {
	content []byte
	addedAt time.Time
}

func New(opts Opts) *Service {
	s := &Service{
		config:            opts.Config,
		filesRepo:         opts.FilesRepo,
		trmManager:        opts.TrmManager,
		s3Client:          opts.S3Client,
		filesContentCache: make(map[string]*cachedFile),
	}

	cacheTime := 10 * time.Minute
	contentClearTicker := time.NewTicker(cacheTime)
	contentClearTickerDone := make(chan struct{})

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					for {
						select {
						case <-contentClearTickerDone:
							contentClearTicker.Stop()
							return
						case <-contentClearTicker.C:
							for k, v := range s.filesContentCache {
								if v.addedAt.Add(cacheTime).Before(time.Now()) {
									delete(s.filesContentCache, k)
								}
							}
						}
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				contentClearTickerDone <- struct{}{}
				close(contentClearTickerDone)
				return nil
			},
		},
	)

	return s
}

type Service struct {
	config     config.Config
	filesRepo  channels_files.Repository
	trmManager trm.Manager
	s3Client   *minio.Client

	filesContentCache map[string]*cachedFile
}

// 10MB
const bytesRestriction = (1 << 20) * 10

var acceptedMimeTypes = []string{"audio", "image"}

func (c *Service) validateUpload(ctx context.Context, channelID string, file entity.Upload) error {
	if file.Size > bytesRestriction {
		return fmt.Errorf("file cannot be bigger than 10mb, got: %v", file.Size)
	}

	if utf8.RuneCountInString(file.Filename) > 100 {
		return fmt.Errorf("file name is too long")
	}

	var isCorrectMimeType bool
	for _, t := range acceptedMimeTypes {
		if strings.HasPrefix(file.ContentType, t) {
			isCorrectMimeType = true
			break
		}
	}

	if !isCorrectMimeType {
		return fmt.Errorf("unsupported file type")
	}

	totalUploaded, err := c.filesRepo.GetTotalChannelUploadedSizeBytes(ctx, channelID)
	if err != nil {
		return err
	}

	if totalUploaded+file.Size > bytesRestriction*10 {
		return fmt.Errorf("limit of storage reached")
	}

	return nil
}

func (c *Service) modelToEntity(m model.ChannelFile) entity.ChannelFile {
	return entity.ChannelFile{
		ID:        m.ID,
		ChannelID: m.ChannelID,
		MimeType:  m.MimeType,
		FileName:  m.FileName,
		Size:      m.Size,
	}
}

func (c *Service) buildS3Path(channelID string, fileID uuid.UUID) string {
	return fmt.Sprintf("channels/%s/%s", channelID, fileID)
}

func (c *Service) Upload(
	ctx context.Context,
	channelID string,
	file entity.Upload,
) (entity.ChannelFile, error) {
	if err := c.validateUpload(ctx, channelID, file); err != nil {
		return entity.ChannelFileNil, err
	}

	var createdFile entity.ChannelFile

	trErr := c.trmManager.Do(
		ctx,
		func(trCtx context.Context) error {
			newFile, err := c.filesRepo.Create(
				trCtx,
				channels_files.CreateInput{
					ChannelID: channelID,
					FileName:  file.Filename,
					MimeType:  file.ContentType,
					Size:      file.Size,
				},
			)
			if err != nil {
				return err
			}

			_, err = c.s3Client.PutObject(
				trCtx,
				c.config.S3Bucket,
				c.buildS3Path(channelID, newFile.ID),
				file.File,
				file.Size,
				minio.PutObjectOptions{
					ContentType: file.ContentType,
				},
			)

			if err != nil {
				return err
			}

			createdFile = c.modelToEntity(newFile)

			return nil
		},
	)
	if trErr != nil {
		return entity.ChannelFileNil, fmt.Errorf("cannot upload file: %w", trErr)
	}

	return createdFile, nil
}

func (c *Service) GetMany(ctx context.Context, channelID string) ([]entity.ChannelFile, error) {
	files, err := c.filesRepo.GetMany(ctx, channels_files.GetManyInput{ChannelID: channelID})
	if err != nil {
		return nil, err
	}

	entities := make([]entity.ChannelFile, len(files))
	for i, file := range files {
		entities[i] = c.modelToEntity(file)
	}

	return entities, nil
}

func (c *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.ChannelFile, error) {
	file, err := c.filesRepo.GetByID(ctx, id)
	if err != nil {
		return entity.ChannelFileNil, err
	}

	return c.modelToEntity(file), nil
}

func (c *Service) DeleteById(ctx context.Context, channelID string, id uuid.UUID) error {
	file, err := c.filesRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if file.ChannelID != channelID {
		return fmt.Errorf("file not found")
	}

	trErr := c.trmManager.Do(
		ctx, func(trCtx context.Context) error {
			err := c.filesRepo.DeleteByID(ctx, id)
			if err != nil {
				return err
			}

			err = c.s3Client.RemoveObject(
				ctx,
				c.config.S3Bucket,
				c.buildS3Path(channelID, id),
				minio.RemoveObjectOptions{},
			)

			return nil
		},
	)
	if trErr != nil {
		return fmt.Errorf("cannot delete file: %w", err)
	}

	return nil
}

func (c *Service) GetFileContent(ctx context.Context, channelID string, fileID uuid.UUID) (
	io.Reader,
	error,
) {
	if f, ok := c.filesContentCache[fileID.String()]; ok {
		c.filesContentCache[fileID.String()].addedAt = time.Now()

		return bytes.NewReader(f.content), nil
	}

	object, err := c.s3Client.GetObject(
		ctx,
		c.config.S3Bucket,
		c.buildS3Path(channelID, fileID),
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot get file: %w", err)
	}

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	}

	c.filesContentCache[fileID.String()] = &cachedFile{
		content: content,
		addedAt: time.Now(),
	}

	return bytes.NewReader(content), nil
}
