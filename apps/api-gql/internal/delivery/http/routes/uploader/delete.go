package uploader

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
	"github.com/twirapp/twir/libs/repositories/uploaded_files"
)

type deleteRequestDto struct {
	PublicId string `path:"publicId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
	Key      string `query:"key"`
}

type deleteByKeyRequestDto struct {
	Key      string `query:"key" required:"true"`
	PublicId string `query:"id" required:"true"`
}

var _ httpbase.Route[*deleteRequestDto, *httpbase.BaseOutputJson[deleteOutputDto]] = (*deleteRoute)(nil)
var _ httpbase.Route[*deleteByKeyRequestDto, *httpbase.BaseOutputJson[deleteOutputDto]] = (*deleteByKeyRoute)(nil)

type DeleteOpts struct {
	Service  *uploaderservice.Service
	Sessions *auth.Auth
}

type deleteRoute struct {
	service  *uploaderservice.Service
	sessions *auth.Auth
}

func newDelete(opts DeleteOpts) *deleteRoute {
	return &deleteRoute{service: opts.Service, sessions: opts.Sessions}
}

func (d *deleteRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "uploader-delete-file",
		Method:      http.MethodDelete,
		Path:        "/v1/uploader/files/{publicId}",
		Tags:        []string{"Uploader"},
		Summary:     "Delete an uploaded file",
	}
}

func (d *deleteRoute) Register(api huma.API) {
	huma.Register(api, d.GetMeta(), d.Handler)
}

func (d *deleteRoute) Handler(ctx context.Context, input *deleteRequestDto) (*httpbase.BaseOutputJson[deleteOutputDto], error) {
	entity, err := d.service.GetByPublicID(ctx, input.PublicId)
	if err != nil {
		if errors.Is(err, uploadedfiles.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found", err)
		}
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get uploaded file", err)
	}
	if entity.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found")
	}
	if err := d.authorize(ctx, entity, input.Key); err != nil {
		return nil, huma.NewError(http.StatusForbidden, err.Error(), err)
	}
	if err := d.service.Delete(ctx, entity); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete uploaded file", err)
	}
	return httpbase.CreateBaseOutputJson(deleteOutputDto{Success: true}), nil
}

func (d *deleteRoute) authorize(ctx context.Context, entity uploadedfile.Entity, key string) error {
	user, _ := d.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil && entity.UploadedByUserID != nil && *entity.UploadedByUserID == user.ID {
		return nil
	}
	if key != "" && key == entity.DeleteKey {
		return nil
	}
	if ids, err := d.sessions.GetLatestUploadedFilesIds(ctx); err == nil && containsID(ids, entity.PublicID) {
		return nil
	}
	return errUploadedFileForbidden
}

type DeleteByKeyOpts struct {
	Service *uploaderservice.Service
}

type deleteByKeyRoute struct {
	service *uploaderservice.Service
}

func newDeleteByKey(opts DeleteByKeyOpts) *deleteByKeyRoute {
	return &deleteByKeyRoute{service: opts.Service}
}

func (d *deleteByKeyRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "uploader-delete-file-by-key",
		Method:      http.MethodGet,
		Path:        "/v1/uploader/files/delete",
		Tags:        []string{"Uploader"},
		Summary:     "Delete an uploaded file by key",
	}
}

func (d *deleteByKeyRoute) Register(api huma.API) {
	huma.Register(api, d.GetMeta(), d.Handler)
}

func (d *deleteByKeyRoute) Handler(ctx context.Context, input *deleteByKeyRequestDto) (*httpbase.BaseOutputJson[deleteOutputDto], error) {
	entity, err := d.service.GetByPublicID(ctx, input.PublicId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Uploaded file not found", err)
	}
	if entity.IsNil() || entity.DeleteKey != input.Key {
		return nil, huma.NewError(http.StatusForbidden, "Invalid delete key")
	}
	if err := d.service.Delete(ctx, entity); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete uploaded file", err)
	}
	return httpbase.CreateBaseOutputJson(deleteOutputDto{Success: true}), nil
}
