package pastebins

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Service  *pastebins.Service
	Sessions *auth.Auth
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "pastebin-get-user-list",
			Summary:     "Get authenticated user pastebins",
			Description: "Requires api-key header.",
			Method:      http.MethodGet,
			Tags:        []string{"Pastebin"},
			Path:        "/v1/pastebin",
			Security: []map[string][]string{
				{"api-key": {}},
			},
		}, func(ctx context.Context, input *getManyInput) (*getManyOutput, error) {
			user, err := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if user == nil || err != nil {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			data, err := opts.Service.GetUserMany(
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

			return &getManyOutput{
				Body: getManyOutputDto{
					Total: data.Total,
					Items: result,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "pastebin-get-by-id",
			Method:      http.MethodGet,
			Path:        "/v1/pastebin/{id}",
			Tags:        []string{"Pastebin"},
			Summary:     "Get pastebin by id",
		}, func(
			ctx context.Context,
			input *struct {
				ID string `path:"id" maxLength:"5" minLength:"1" pattern:"^[-_a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*getByIdOutput, error,
		) {
			bin, err := opts.Service.GetByID(ctx, input.ID)
			if err != nil {
				if errors.Is(err, pastebins.ErrNotFound) {
					return nil, huma.NewError(http.StatusNotFound, "Cannot get pastebin", err)
				}

				return nil, huma.NewError(http.StatusInternalServerError, "Cannot get pastebin", err)
			}

			return &getByIdOutput{
				Body: pasteBinOutputDto{
					ID:          bin.ID,
					CreatedAt:   bin.CreatedAt,
					Content:     bin.Content,
					ExpireAt:    bin.ExpireAt,
					OwnerUserID: bin.OwnerUserID,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "pastebin-create",
			Method:      http.MethodPost,
			Path:        "/v1/pastebin",
			Tags:        []string{"Pastebin"},
			Summary:     "Create pastebin",
		},
		func(
			ctx context.Context,
			input *struct {
				Body pasteBinCreateDto
			},
		) (
			*getByIdOutput, error,
		) {
			createInput := pastebins.CreateInput{
				Content:  input.Body.Content,
				ExpireAt: input.Body.ExpireAt,
			}

			user, err := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if err == nil && user != nil {
				createInput.OwnerUserID = &user.ID
			}

			bin, err := opts.Service.Create(ctx, createInput)
			if err != nil {
				return nil, huma.NewError(http.StatusInternalServerError, "Cannot create pastebin", err)
			}

			return &getByIdOutput{
				Body: pasteBinOutputDto{
					ID:          bin.ID,
					CreatedAt:   bin.CreatedAt,
					Content:     bin.Content,
					ExpireAt:    bin.ExpireAt,
					OwnerUserID: bin.OwnerUserID,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "pastebin-delete",
			Method:      http.MethodDelete,
			Path:        "/v1/pastebin/{id}",
			Tags:        []string{"Pastebin"},
			Summary:     "Delete pastebin",
			Security: []map[string][]string{
				{"api-key": {}},
			},
		}, func(
			ctx context.Context,
			input *struct {
				ID string `path:"id" maxLength:"5" minLength:"1" pattern:"^[-_a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*getByIdOutput, error,
		) {
			paste, err := opts.Service.GetByID(ctx, input.ID)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get pastebin", err)
			}

			user, err := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if user == nil || err != nil || paste.OwnerUserID == nil || *paste.OwnerUserID != user.ID {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			if err := opts.Service.Delete(ctx, input.ID); err != nil {
				return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete pastebin", err)
			}

			return nil, nil
		},
	)
}

type getByIdOutput struct {
	Body pasteBinOutputDto
}

type pasteBinOutputDto struct {
	ID          string     `json:"id" example:"KKMEa"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time"`
	Content     string     `json:"content" example:"Hello world"`
	ExpireAt    *time.Time `json:"expire_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time" nullable:"true"`
	OwnerUserID *string    `json:"owner_user_id" example:"1234567890" nullable:"true"`
}

type pasteBinCreateDto struct {
	Content  string     `json:"content" required:"true" minLength:"1" maxLength:"100000" example:"Hello world"`
	ExpireAt *time.Time `json:"expire_at" format:"date-time" nullable:"true" required:"false"`
}

type getManyInput struct {
	Page    uint `json:"page" query:"page" example:"1" default:"1" minimum:"1"`
	PerPage uint `json:"perPage" query:"perPage" example:"20" default:"20"`
}

type getManyOutput struct {
	Body getManyOutputDto
}

type getManyOutputDto struct {
	Total int                 `json:"total" example:"1"`
	Items []pasteBinOutputDto `json:"items"`
}
