package uploader

import (
	"context"
	"net/http"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	uploaderservice "github.com/twirapp/twir/apps/api-gql/internal/services/uploader"
	uploadedfiles "github.com/twirapp/twir/libs/repositories/uploaded_files"
)

type profileRequestDto struct {
	Page    int `query:"page" minimum:"0" default:"0"`
	PerPage int `query:"per_page" minimum:"1" default:"10" maximum:"50"`
}

var _ httpbase.Route[*profileRequestDto, *httpbase.BaseOutputJson[profileOutputDto]] = (*profile)(nil)

type ProfileOpts struct {
	Service  *uploaderservice.Service
	Sessions *auth.Auth
}

type profile struct {
	service  *uploaderservice.Service
	sessions *auth.Auth
}

func newProfile(opts ProfileOpts) *profile {
	return &profile{service: opts.Service, sessions: opts.Sessions}
}

func (p *profile) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "uploader-get-files",
		Method:      http.MethodGet,
		Path:        "/v1/uploader/files",
		Tags:        []string{"Uploader"},
		Summary:     "Get uploaded files",
	}
}

func (p *profile) Register(api huma.API) {
	huma.Register(api, p.GetMeta(), p.Handler)
}

func (p *profile) Handler(ctx context.Context, input *profileRequestDto) (*httpbase.BaseOutputJson[profileOutputDto], error) {
	files := make([]uploadedFileOutputDto, 0)
	total := 0
	user, userErr := p.sessions.GetAuthenticatedUserModel(ctx)
	if userErr == nil && user != nil {
		data, err := p.service.GetList(ctx, uploadedfiles.GetListInput{
			UserID: user.ID, Page: input.Page, PerPage: input.PerPage,
		})
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get uploaded files", err)
		}
		total = data.Total
		for _, file := range data.Items {
			files = append(files, toOutput(p.service, file))
		}
	}

	if ids, err := p.sessions.GetLatestUploadedFilesIds(ctx); err == nil {
		data, err := p.service.GetManyByPublicIDs(ctx, ids)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get uploaded files", err)
		}
		for _, file := range data {
			files = append(files, toOutput(p.service, file))
		}
	}

	seen := make(map[string]struct{}, len(files))
	unique := make([]uploadedFileOutputDto, 0, len(files))
	for _, file := range files {
		if _, ok := seen[file.Id]; ok {
			continue
		}
		seen[file.Id] = struct{}{}
		unique = append(unique, file)
	}
	slices.SortFunc(unique, func(left, right uploadedFileOutputDto) int {
		return right.CreatedAt.Compare(left.CreatedAt)
	})
	if len(unique) > input.PerPage {
		unique = unique[:input.PerPage]
	}
	if total == 0 {
		total = len(unique)
	}

	return httpbase.CreateBaseOutputJson(profileOutputDto{Files: unique, Total: total}), nil
}
