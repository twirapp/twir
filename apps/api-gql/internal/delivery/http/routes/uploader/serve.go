package uploader

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
	"github.com/twirapp/twir/libs/logger"
)

type serveRequestDto struct {
	PublicId string `path:"publicId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

var _ httpbase.Route[*serveRequestDto, *huma.StreamResponse] = (*serve)(nil)

type ServeOpts struct {
	Service *uploaderservice.Service
	Logger  *slog.Logger
}

type serve struct {
	service *uploaderservice.Service
	logger  *slog.Logger
	api     huma.API
}

func newServe(opts ServeOpts) *serve {
	return &serve{service: opts.Service, logger: opts.Logger}
}

func (s *serve) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "uploader-serve-file",
		Method:        http.MethodGet,
		Path:          "/v1/u/{publicId}",
		Tags:          []string{"Uploader"},
		Summary:       "Serve an uploaded file",
		DefaultStatus: http.StatusOK,
	}
}

func (s *serve) Register(api huma.API) {
	s.api = api
	huma.Register(api, s.GetMeta(), s.Handler)
}

func (s *serve) Handler(ctx context.Context, input *serveRequestDto) (*huma.StreamResponse, error) {
	entity, err := s.service.GetByPublicID(ctx, input.PublicId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found", err)
	}
	if entity.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found")
	}
	if isExpired(entity) {
		if err := s.service.Delete(ctx, entity); err != nil {
			s.logger.WarnContext(ctx, "Cannot delete expired uploaded file", logger.Error(err))
		}
		return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found")
	}

	return &huma.StreamResponse{
		Body: func(humaCtx huma.Context) {
			humaCtx.SetHeader("Content-Type", entity.MimeType)
			humaCtx.SetHeader("Content-Length", strconv.FormatInt(entity.SizeBytes, 10))
			humaCtx.SetHeader("Cache-Control", "public, max-age=86400")
			humaCtx.SetHeader("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", safeFileName(entity)))

			reader, err := s.service.GetObject(ctx, entity)
			if err != nil {
				s.logger.WarnContext(ctx, "Cannot get uploaded object", logger.Error(err))
				_ = huma.WriteErr(s.api, humaCtx, http.StatusNotFound, "Uploaded file not found", err)
				return
			}
			defer func() {
				if closeErr := reader.Close(); closeErr != nil {
					s.logger.WarnContext(ctx, "Cannot close uploaded object", logger.Error(closeErr))
				}
			}()

			if _, err := io.Copy(humaCtx.BodyWriter(), reader); err != nil {
				s.logger.WarnContext(ctx, "Cannot stream uploaded object", logger.Error(err))
			}
		},
	}, nil
}

func safeFileName(entity uploadedfile.Entity) string {
	name := "image" + entity.Extension
	if entity.FileName != nil && *entity.FileName != "" {
		name = *entity.FileName
	}
	name = strings.NewReplacer("\"", "", "\\", "", "\r", "", "\n", "").Replace(name)
	if name == "" {
		return "image" + entity.Extension
	}
	return name
}
