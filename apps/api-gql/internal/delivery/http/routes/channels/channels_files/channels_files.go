package channels_files

import (
	"context"
	"io"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_files"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api                  huma.API
	Config               config.Config
	ChannelsFilesService *channels_files.Service
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/channels/{channel_id}/files/content/{file_id}",
			Tags:        []string{"Files"},
			Summary:     "Get file content",
			Description: "Get file content by id",
			Responses: map[string]*huma.Response{
				"200": {
					Description: "File content",
					Content: map[string]*huma.MediaType{
						"application/octet-stream": {
							Schema: &huma.Schema{
								Type:        "string",
								Format:      "binary",
								Description: "File content",
							},
						},
					},
				},
			},
		},
		func(
			ctx context.Context, i *struct {
				ChannelID string    `path:"channel_id" maxLength:"36" minLength:"1" pattern:"^[0-9]+$" required:"true"`
				FileID    uuid.UUID `path:"file_id" maxLength:"36" minLength:"1" format:"uuid" required:"true"`
			},
		) (*huma.StreamResponse, error) {
			foundFile, err := opts.ChannelsFilesService.GetByID(ctx, i.FileID)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get file", err)
			}

			if foundFile.ChannelID != i.ChannelID {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get file")
			}

			return &huma.StreamResponse{
				Body: func(humaCtx huma.Context) {
					humaCtx.SetHeader("Content-Type", foundFile.MimeType)
					humaCtx.SetHeader("Content-Disposition", "attachment; filename="+foundFile.FileName)
					writer := humaCtx.BodyWriter()

					reader, err := opts.ChannelsFilesService.GetFileContent(ctx, i.ChannelID, i.FileID)
					if err != nil {
						_ = huma.WriteErr(
							opts.Api,
							humaCtx,
							http.StatusInternalServerError,
							"Cannot get file",
							err,
						)
						return
					}

					_, err = io.Copy(writer, reader)
					if err != nil {
						_ = huma.WriteErr(
							opts.Api,
							humaCtx,
							http.StatusInternalServerError,
							"Cannot get file",
							err,
						)
						return
					}
				},
			}, nil
		},
	)
}
