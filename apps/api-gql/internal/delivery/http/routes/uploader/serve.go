package uploader

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	"github.com/twirapp/twir/libs/logger"
)

type serveRequestDto struct {
	PublicId string `path:"publicId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

type serveResponseDto struct {
	Status       int
	Location     string `header:"Location"`
	CacheControl string `header:"Cache-Control"`
}

var _ httpbase.Route[*serveRequestDto, *serveResponseDto] = (*serve)(nil)

type ServeOpts struct {
	Service *uploaderservice.Service
	Logger  *slog.Logger
}

type serve struct {
	service *uploaderservice.Service
	logger  *slog.Logger
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
		DefaultStatus: http.StatusFound,
	}
}

func (s *serve) Register(api huma.API) {
	huma.Register(api, s.GetMeta(), s.Handler)
}

func (s *serve) Handler(ctx context.Context, input *serveRequestDto) (*serveResponseDto, error) {
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

	return &serveResponseDto{
		Status:       http.StatusFound,
		Location:     s.service.BuildS3URL(entity),
		CacheControl: "public, max-age=86400",
	}, nil
}
