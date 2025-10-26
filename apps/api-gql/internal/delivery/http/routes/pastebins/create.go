package pastebins

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"go.uber.org/fx"
)

type pasteBinCreateRequestDto struct {
	Body struct {
		Content  string     `json:"content" required:"true" minLength:"1" maxLength:"100000" example:"Hello world"`
		ExpireAt *time.Time `json:"expire_at" format:"date-time" nullable:"true" required:"false"`
	}
}

var _ httpbase.Route[*pasteBinCreateRequestDto, *httpbase.BaseOutputJson[pasteBinOutputDto]] = (*create)(nil)

type CreateOpts struct {
	fx.In

	Service  *pastebins.Service
	Sessions *auth.Auth
}

func newCreate(opts CreateOpts) *create {
	return &create{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type create struct {
	service  *pastebins.Service
	sessions *auth.Auth
}

func (c *create) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "pastebin-create",
		Method:      http.MethodPost,
		Path:        "/v1/pastebin",
		Tags:        []string{"Pastebin"},
		Summary:     "Create pastebin",
	}
}

func (c *create) Handler(
	ctx context.Context,
	input *pasteBinCreateRequestDto,
) (*httpbase.BaseOutputJson[pasteBinOutputDto], error) {
	createInput := pastebins.CreateInput{
		Content:  input.Body.Content,
		ExpireAt: input.Body.ExpireAt,
	}

	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err == nil && user != nil {
		createInput.OwnerUserID = &user.ID
	}

	bin, err := c.service.Create(ctx, createInput)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot create pastebin", err)
	}

	return httpbase.CreateBaseOutputJson(
		pasteBinOutputDto{
			ID:          bin.ID,
			CreatedAt:   bin.CreatedAt,
			Content:     bin.Content,
			ExpireAt:    bin.ExpireAt,
			OwnerUserID: bin.OwnerUserID,
		},
	), nil
}

func (c *create) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
