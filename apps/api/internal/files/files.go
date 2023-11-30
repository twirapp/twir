package files

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/satont/twir/apps/api/internal/handlers"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Db     *gorm.DB
	Config cfg.Config
	Logger logger.Logger
}

type Files struct {
	db          *gorm.DB
	config      cfg.Config
	minioClient *minio.Client
	logger      logger.Logger
}

const bucketName = "twir"

func NewFiles(opts Opts) handlers.IHandler {
	client, err := minio.New(
		opts.Config.S3Host,
		&minio.Options{
			Creds:  credentials.NewStaticV4(opts.Config.S3AccessToken, opts.Config.S3SecretToken, ""),
			Region: opts.Config.S3Region,
			Secure: opts.Config.AppEnv == "production",
		},
	)
	if err != nil {
		opts.Logger.Error("cannot create minio host", slog.Any("err", err))
	}

	return &Files{
		db:          opts.Db,
		config:      opts.Config,
		minioClient: client,
	}
}

func (c *Files) Pattern() string {
	return "/files/"
}

func (c *Files) Handler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			channelId := r.URL.Query().Get("channel_id")
			fileId := r.URL.Query().Get("file_id")

			if channelId == "" || fileId == "" {
				http.Error(w, "Missed params", http.StatusBadRequest)
				return
			}

			entity := model.ChannelFile{}
			if err := c.db.
				WithContext(r.Context()).
				Where("id = ? and channel_id = ?", fileId, channelId).Find(&entity).Error; err != nil {
				c.logger.Error(
					"cannot get file",
					slog.Any("err", err),
					slog.String("channelId", channelId),
					slog.String("fileId", fileId),
				)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			if entity.ID == "" {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			reader, err := c.minioClient.GetObject(
				r.Context(),
				bucketName,
				fmt.Sprintf("channels/%s/%s", channelId, entity.ID),
				minio.GetObjectOptions{},
			)
			if err != nil {
				c.logger.Error(
					"cannot get file",
					slog.Any("err", err),
					slog.String("channelId", channelId),
					slog.String("fileId", fileId),
				)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", entity.MimeType)
			_, copyErr := io.Copy(w, reader)
			if copyErr != nil {
				http.Error(
					w,
					fmt.Sprintf("Error copying file to the http response %s", copyErr.Error()),
					http.StatusInternalServerError,
				)
				return
			}
		},
	)
}
