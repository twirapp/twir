package files

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/files"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Files struct {
	*impl_deps.Deps
	s3Client *minio.Client
}

const bucketName = "twir"

func New(deps *impl_deps.Deps) *Files {
	var minioClient *minio.Client

	if deps.Config.S3Host != "" {
		client, err := minio.New(
			deps.Config.S3Host,
			&minio.Options{
				Creds:  credentials.NewStaticV4(deps.Config.S3AccessToken, deps.Config.S3SecretToken, ""),
				Region: deps.Config.S3Region,
				Secure: deps.Config.AppEnv == "production",
			},
		)
		if err != nil {
			deps.Logger.Error("cannot create minio host", slog.Any("err", err))
		}
		minioClient = client

		ctx := context.Background()
		err = minioClient.MakeBucket(
			ctx,
			deps.Config.S3Bucket,
			minio.MakeBucketOptions{Region: deps.Config.S3Region},
		)

		if err != nil {
			exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
			if errBucketExists != nil && !exists {
				deps.Logger.Error("Cannot create bucket", slog.Any("err", err))
			}
		} else {
			deps.Logger.Info("Successfully created bucket", slog.String("name", bucketName))
		}

		err = minioClient.SetBucketPolicy(
			ctx,
			bucketName,
			`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"AWS": ["*"]
				},
				"Action": ["s3:GetObject"],
				"Resource": [
					"arn:aws:s3:::`+bucketName+`/**"
				]
			}
		]
	}`,
		)
		if err != nil {
			deps.Logger.Error("cannot set policy", slog.Any("err", err))
		} else {
			deps.Logger.Info("Bucket policy was set")
		}
	}

	return &Files{
		Deps:     deps,
		s3Client: minioClient,
	}
}

// 10MB
const bytesRestriction = (1 << 20) * 10

var acceptedMimeTypes = []string{"audio", "image"}

func (c *Files) FilesUpload(ctx context.Context, req *files.UploadRequest) (
	*files.FileMeta,
	error,
) {
	if len(req.Content) > bytesRestriction {
		return nil, twirp.NewError(twirp.OutOfRange, "File cannot be bigger than 10mb")
	}

	if utf8.RuneCountInString(req.Name) > 100 {
		return nil, twirp.NewError(twirp.OutOfRange, "File name is too long")
	}

	if !lo.SomeBy(
		acceptedMimeTypes, func(t string) bool {
			return strings.HasPrefix(req.Mimetype, t)
		},
	) {
		return nil, twirp.NewError(twirp.OutOfRange, "Wrong file type")
	}

	dashboardId := ctx.Value("dashboardId").(string)

	type NResult struct {
		N int
	}
	var filesSize NResult
	if err := c.Db.WithContext(ctx).
		Model(&model.ChannelFile{}).
		Select("sum(size) as n").Scan(&filesSize).
		Error; err != nil {
		c.Logger.Error(
			"cannot count user files size",
			slog.String("channelId", dashboardId),
			slog.Any("err", err),
		)
		return nil, err
	}

	if filesSize.N+len(req.Content) > bytesRestriction*10 {
		return nil, twirp.NewError(twirp.OutOfRange, "Limit of storage reached")
	}

	fileEntity := model.ChannelFile{
		ID:        uuid.New().String(),
		ChannelID: dashboardId,
		MimeType:  req.Mimetype,
		Name:      req.Name,
		Size:      len(req.Content),
	}

	err := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Create(&fileEntity).Error; err != nil {
				return err
			}

			reader := bytes.NewReader(req.Content)

			_, err := c.s3Client.PutObject(
				ctx,
				bucketName,
				fmt.Sprintf("channels/%s/%s", dashboardId, fileEntity.ID),
				reader,
				int64(len(req.Content)),
				minio.PutObjectOptions{
					ContentType: req.Mimetype,
				},
			)

			return err
		},
	)

	if err != nil {
		c.Logger.Error(
			"cannot upload file",
			slog.String("channelId", dashboardId),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &files.FileMeta{
		Id:        fileEntity.ID,
		Mimetype:  fileEntity.MimeType,
		Name:      fileEntity.Name,
		ChannelId: fileEntity.ChannelID,
		Size:      int64(fileEntity.Size),
	}, nil
}

func (c *Files) FilesGetAll(ctx context.Context, _ *emptypb.Empty) (*files.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []model.ChannelFile
	if err := c.Db.WithContext(ctx).Where(
		"channel_id = ?",
		dashboardId,
	).Find(&entities).Error; err != nil {
		c.Logger.Error(
			"cannot get files",
			slog.String("channelId", dashboardId),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &files.GetAllResponse{
		Files: lo.Map(
			entities,
			func(item model.ChannelFile, _ int) *files.FileMeta {
				return &files.FileMeta{
					Id:        item.ID,
					Mimetype:  item.MimeType,
					Name:      item.Name,
					ChannelId: item.ChannelID,
					Size:      int64(item.Size),
				}
			},
		),
	}, nil
}

func (c *Files) FilesDelete(ctx context.Context, req *files.RemoveRequest) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	err := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where(
				"id = ? AND channel_id = ?",
				req.Id,
				dashboardId,
			).Delete(&model.ChannelFile{}).
				Error; err != nil {
				return err
			}

			if err := c.s3Client.RemoveObject(
				ctx,
				bucketName,
				fmt.Sprintf("channels/%s/%s", dashboardId, req.Id),
				minio.RemoveObjectOptions{},
			); err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		c.Logger.Error(
			"cannot remove file",
			slog.String("channelId", dashboardId),
			slog.String("fileId", req.Id),
			slog.Any("err", err),
		)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
