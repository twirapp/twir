package pastebins

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"go.uber.org/fx"
)

type profileRequestDto struct {
	Page    uint `json:"page" query:"page" example:"1" default:"1" minimum:"1"`
	PerPage uint `json:"perPage" query:"perPage" example:"20" default:"20"`
}

type profileResponseDto struct {
	Total int                 `json:"total" example:"1"`
	Items []pasteBinOutputDto `json:"items"`
}

var _ httpbase.Route[*profileRequestDto, *httpbase.BaseOutputJson[profileResponseDto]] = (*profile)(nil)

type ProfileOpts struct {
	fx.In

	Service  *pastebins.Service
	Sessions *auth.Auth
}

func newProfile(opts ProfileOpts) *profile {
	return &profile{
		sessions: opts.Sessions,
		service:  opts.Service,
	}
}

type profile struct {
	sessions *auth.Auth
	service  *pastebins.Service
}

func (p *profile) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "pastebin-get-user-list",
		Summary:     "Get authenticated user pastebins",
		Description: "Requires api-key header.",
		Method:      http.MethodGet,
		Tags:        []string{"Pastebin"},
		Path:        "/v1/pastebin",
		Security: []map[string][]string{
			{"api-key": {}},
		},
	}
}

func (p *profile) Handler(
	ctx context.Context,
	input *profileRequestDto,
) (*httpbase.BaseOutputJson[profileResponseDto], error) {
	user, err := p.sessions.GetAuthenticatedUserModel(ctx)
	if user == nil || err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
	}

	data, err := p.service.GetUserMany(
		ctx, pastebins.GetManyInput{
			Page:        int(input.Page),
			PerPage:     int(input.PerPage),
			OwnerUserID: user.ID,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get pastebins", err)
	}

	result := make([]pasteBinOutputDto, 0, len(data.Items))
	for _, bin := range data.Items {
		result = append(
			result,
			pasteBinOutputDto{
				ID:          bin.ID,
				CreatedAt:   bin.CreatedAt,
				Content:     bin.Content,
				ExpireAt:    bin.ExpireAt,
				OwnerUserID: bin.OwnerUserID,
			},
		)
	}

	return httpbase.CreateBaseOutputJson(
		profileResponseDto{
			Total: data.Total,
			Items: result,
		},
	), nil
}

func (p *profile) Register(api huma.API) {
	huma.Register(api, p.GetMeta(), p.Handler)
}
