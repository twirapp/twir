package uploader

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
)

type uploadFormData struct {
	File huma.FormFile `form:"file" contentType:"image/*" required:"true"`
}

type createInput struct {
	RawBody huma.MultipartFormFiles[uploadFormData]
}

type createOutput struct {
	Body struct {
		Data uploadedFileWithDeleteLinkDto `json:"data"`
	}
}

var _ httpbase.Route[*createInput, *createOutput] = (*create)(nil)

type CreateOpts struct {
	Config            config.Config
	Service           *uploaderservice.Service
	Sessions          *auth.Auth
	Logger            *slog.Logger
	Middlewares       *middlewares.Middlewares
	ClientInfoService *clientinfo.Service
}

type create struct {
	config            config.Config
	service           *uploaderservice.Service
	sessions          *auth.Auth
	logger            *slog.Logger
	middlewares       *middlewares.Middlewares
	clientInfoService *clientinfo.Service
}

func newCreate(opts CreateOpts) *create {
	return &create{
		config:            opts.Config,
		service:           opts.Service,
		sessions:          opts.Sessions,
		logger:            opts.Logger,
		middlewares:       opts.Middlewares,
		clientInfoService: opts.ClientInfoService,
	}
}

func (c *create) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "uploader-upload-file",
		Method:      http.MethodPost,
		Path:        "/v1/uploader/files",
		Tags:        []string{"Uploader"},
		Summary:     "Upload a file",
		Middlewares: huma.Middlewares{c.middlewares.RateLimit("uploader-upload-file", 10, time.Minute)},
	}
}

func (c *create) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}

func (c *create) Handler(ctx context.Context, input *createInput) (*createOutput, error) {
	file := input.RawBody.Data().File
	defer func() {
		if err := file.File.Close(); err != nil {
			c.logger.WarnContext(ctx, "Cannot close uploaded file", logger.Error(err))
		}
	}()

	content, err := io.ReadAll(io.LimitReader(file.File, c.config.UploaderMaxFileSizeBytes+1))
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot read file", err)
	}
	if int64(len(content)) > c.config.UploaderMaxFileSizeBytes {
		return nil, huma.NewError(http.StatusRequestEntityTooLarge, "File is too large", uploaderservice.ErrFileTooLarge)
	}

	var userID *string
	user, _ := c.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil {
		userID = &user.ID
	}

	clientInfo, err := c.clientInfoService.GetClientInfo(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Internal error on getting your information", err)
	}
	clientIP, err := netip.ParseAddr(clientInfo.IP)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Internal error on getting your information", err)
	}

	var fileName *string
	if file.Filename != "" {
		fileName = &file.Filename
	}
	entity, err := c.service.Upload(ctx, uploaderservice.UploadInput{
		File:      bytes.NewReader(content),
		Size:      int64(len(content)),
		FileName:  fileName,
		UserID:    userID,
		UserAgent: &clientInfo.UserAgent,
		UserIP:    &clientIP,
	})
	if err != nil {
		switch {
		case errors.Is(err, uploaderservice.ErrFileTooLarge):
			return nil, huma.NewError(http.StatusRequestEntityTooLarge, "File is too large", err)
		case errors.Is(err, uploaderservice.ErrUnsupportedFileType):
			return nil, huma.NewError(http.StatusUnsupportedMediaType, "Unsupported file type", err)
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot upload file", err)
		}
	}

	if err := c.sessions.AddLatestUploadedFileId(ctx, entity.PublicID); err != nil {
		c.logger.WarnContext(ctx, "Cannot save latest uploaded file id to session", logger.Error(err))
	}

	deleteURL := strings.TrimRight(c.config.SiteBaseUrl, "/") + "/api/v1/uploader/files/delete?key=" + url.QueryEscape(entity.DeleteKey) + "&id=" + url.QueryEscape(entity.PublicID)
	return &createOutput{Body: struct {
		Data uploadedFileWithDeleteLinkDto `json:"data"`
	}{Data: uploadedFileWithDeleteLinkDto{
		UploadedFileOutputDto: toOutput(c.service, entity),
		DeleteLink:            deleteURL,
	}}}, nil
}
