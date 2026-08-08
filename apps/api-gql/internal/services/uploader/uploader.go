package uploader

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/minio/minio-go/v7"
	minioclient "github.com/twirapp/twir/apps/api-gql/internal/minio"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	cfg "github.com/twirapp/twir/libs/config"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
	logger "github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/uploaded_files"
)

const (
	publicIDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	publicIDLength   = 8
	deleteIDLength   = 24
	cleanupLimit     = 500
	cleanupInterval  = 5 * time.Minute
)

var (
	ErrFileTooLarge        = errors.New("file is too large")
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrPublicIDUnavailable = errors.New("public ID unavailable")
	supportedMimeTypes     = map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
		"image/gif":  ".gif",
		"image/webp": ".webp",
		"image/avif": ".avif",
		"image/bmp":  ".bmp",
	}
)

type UploadInput struct {
	File      io.Reader
	Size      int64
	FileName  *string
	UserID    *string
	UserAgent *string
	UserIP    *netip.Addr
}

type Service struct {
	config         cfg.Config
	logger         *slog.Logger
	uploaderClient *minioclient.UploaderS3Client
	repository     uploadedfiles.Repository

	cleanupCtx    context.Context
	cleanupCancel context.CancelFunc
	cleanupWG     sync.WaitGroup
}

func New(
	config cfg.Config,
	logger *slog.Logger,
	uploaderClient *minioclient.UploaderS3Client,
	repository uploadedfiles.Repository,
	lc *lifecycle.Lifecycle,
) *Service {
	cleanupCtx, cleanupCancel := context.WithCancel(context.Background())
	service := &Service{
		config:         config,
		logger:         logger,
		uploaderClient: uploaderClient,
		repository:     repository,
		cleanupCtx:     cleanupCtx,
		cleanupCancel:  cleanupCancel,
	}

	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			service.cleanupWG.Go(func() {
				service.cleanupLoop()
			})
			return nil
		},
		OnStop: func(context.Context) error {
			service.cleanupCancel()
			service.cleanupWG.Wait()
			return nil
		},
	})

	return service
}

func (c *Service) Upload(ctx context.Context, input UploadInput) (uploadedfile.Entity, error) {
	if input.Size > c.config.UploaderMaxFileSizeBytes {
		return uploadedfile.Nil, ErrFileTooLarge
	}

	sample, err := io.ReadAll(io.LimitReader(input.File, 512))
	if err != nil {
		return uploadedfile.Nil, fmt.Errorf("read upload header: %w", err)
	}
	detectedMime := http.DetectContentType(sample)
	extension, ok := supportedMimeTypes[detectedMime]
	if !ok {
		return uploadedfile.Nil, fmt.Errorf("%s: %w", detectedMime, ErrUnsupportedFileType)
	}
	reader := io.MultiReader(bytes.NewReader(sample), input.File)

	now := time.Now()
	expiresAt := now.Add(c.config.UploaderAnonymousFileTTL)
	if input.UserID != nil {
		expiresAt = now.Add(c.config.UploaderAuthenticatedFileTTL)
	}

	var publicID string
	for attempts := 0; ; attempts++ {
		candidate, err := gonanoid.Generate(publicIDAlphabet, publicIDLength)
		if err != nil {
			return uploadedfile.Nil, fmt.Errorf("generate public ID: %w", err)
		}

		existing, err := c.repository.GetByPublicID(ctx, candidate)
		switch {
		case err == nil && !existing.IsNil():
			if attempts >= 4 {
				return uploadedfile.Nil, ErrPublicIDUnavailable
			}
			continue
		case err != nil && !errors.Is(err, uploadedfiles.ErrNotFound):
			return uploadedfile.Nil, fmt.Errorf("check public ID: %w", err)
		}

		publicID = candidate
		break
	}

	deleteKey, err := gonanoid.New(deleteIDLength)
	if err != nil {
		return uploadedfile.Nil, fmt.Errorf("generate delete key: %w", err)
	}
	s3Key := "uploads/" + publicID + extension
	if _, err := c.uploaderClient.PutObject(
		ctx,
		c.uploaderClient.Bucket,
		s3Key,
		reader,
		input.Size,
		minio.PutObjectOptions{ContentType: detectedMime},
	); err != nil {
		return uploadedfile.Nil, fmt.Errorf("upload file to storage: %w", err)
	}

	created, err := c.repository.Create(ctx, uploadedfiles.CreateInput{
		PublicID:         publicID,
		UploadedByUserID: input.UserID,
		FileName:         input.FileName,
		MimeType:         detectedMime,
		Extension:        extension,
		SizeBytes:        input.Size,
		S3Key:            s3Key,
		DeleteKey:        deleteKey,
		UserAgent:        input.UserAgent,
		UserIP:           input.UserIP,
		ExpiresAt:        expiresAt,
	})
	if err == nil {
		return created, nil
	}

	if cleanupErr := c.removeObject(ctx, s3Key); cleanupErr != nil {
		c.logger.WarnContext(ctx, "failed to clean up uploaded object", logger.Error(cleanupErr), slog.String("s3_key", s3Key))
	}
	return uploadedfile.Nil, fmt.Errorf("create uploaded file: %w", err)
}

func (c *Service) Delete(ctx context.Context, entity uploadedfile.Entity) error {
	if err := c.removeObject(ctx, entity.S3Key); err != nil {
		return fmt.Errorf("delete uploaded object: %w", err)
	}
	if err := c.repository.DeleteByID(ctx, entity.ID); err != nil {
		return fmt.Errorf("delete uploaded file: %w", err)
	}
	return nil
}

func (c *Service) GetByPublicID(ctx context.Context, publicID string) (uploadedfile.Entity, error) {
	return c.repository.GetByPublicID(ctx, publicID)
}

func (c *Service) GetManyByPublicIDs(ctx context.Context, ids []string) ([]uploadedfile.Entity, error) {
	return c.repository.GetManyByPublicIDs(ctx, ids)
}

func (c *Service) GetList(ctx context.Context, input uploadedfiles.GetListInput) (uploadedfiles.GetListOutput, error) {
	return c.repository.GetList(ctx, input)
}

func (c *Service) BuildPublicURL(entity uploadedfile.Entity) string {
	return strings.TrimRight(c.config.SiteBaseUrl, "/") + "/u/" + entity.PublicID
}

func (c *Service) GetObject(ctx context.Context, entity uploadedfile.Entity) (io.ReadCloser, error) {
	return c.uploaderClient.GetObject(ctx, c.uploaderClient.Bucket, entity.S3Key, minio.GetObjectOptions{})
}

func (c *Service) CleanupExpired(ctx context.Context) error {
	files, err := c.repository.GetExpired(ctx, cleanupLimit)
	if err != nil {
		return fmt.Errorf("get expired uploaded files: %w", err)
	}
	for _, file := range files {
		if err := c.removeObject(ctx, file.S3Key); err != nil {
			c.logger.ErrorContext(ctx, "failed to remove expired uploaded object", logger.Error(err), slog.String("s3_key", file.S3Key))
			continue
		}
		if err := c.repository.DeleteByID(ctx, file.ID); err != nil {
			c.logger.ErrorContext(ctx, "failed to delete expired uploaded file", logger.Error(err), slog.String("public_id", file.PublicID))
		}
	}
	return nil
}

func (c *Service) cleanupLoop() {
	if err := c.CleanupExpired(context.Background()); err != nil {
		c.logger.Error("failed to clean up expired uploaded files", logger.Error(err))
	}
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-c.cleanupCtx.Done():
			return
		case <-ticker.C:
			if err := c.CleanupExpired(context.Background()); err != nil {
				c.logger.Error("failed to clean up expired uploaded files", logger.Error(err))
			}
		}
	}
}

func (c *Service) removeObject(ctx context.Context, key string) error {
	err := c.uploaderClient.RemoveObject(ctx, c.uploaderClient.Bucket, key, minio.RemoveObjectOptions{})
	if err != nil && !isObjectNotFound(err) {
		return err
	}
	return nil
}

func isObjectNotFound(err error) bool {
	code := minio.ToErrorResponse(err).Code
	return code == "NoSuchKey" || code == "NoSuchObject" || code == "NotFound"
}
