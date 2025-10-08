package pastebins

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"go.uber.org/fx"
)

type getByIdRequestDto struct {
	ID string `path:"id" maxLength:"5" minLength:"1" pattern:"^[-_a-zA-Z0-9]+$" required:"true"`
}

var _ httpbase.Route[*getByIdRequestDto, *httpbase.BaseOutputJson[pasteBinOutputDto]] = (*getById)(nil)

type GetByIdOpts struct {
	fx.In

	Service *pastebins.Service
}

func newGetById(opts GetByIdOpts) *getById {
	return &getById{
		service: opts.Service,
	}
}

type getById struct {
	service *pastebins.Service
}

func (g *getById) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "pastebin-get-by-id",
		Method:      http.MethodGet,
		Path:        "/v1/pastebin/{id}",
		Tags:        []string{"Pastebin"},
		Summary:     "Get pastebin by id",
	}
}

func (g *getById) Handler(
	ctx context.Context,
	input *getByIdRequestDto,
) (*httpbase.BaseOutputJson[pasteBinOutputDto], error) {
	bin, err := g.service.GetByID(ctx, input.ID)
	if err != nil {
		if errors.Is(err, pastebins.ErrNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get pastebin", err)
		}

		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get pastebin", err)
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

func (g *getById) Register(api huma.API) {
	huma.Register(
		api,
		g.GetMeta(),
		g.Handler,
	)
}
