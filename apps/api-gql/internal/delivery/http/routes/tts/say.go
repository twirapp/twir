package tts

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newSay),
)

type sayRequestDto struct {
	Voice  string `query:"voice" minLength:"1" maxLength:"100" example:"alan" doc:"Voice name to use for TTS"`
	Text   string `query:"text" minLength:"1" maxLength:"5000" example:"Hello world" doc:"Text to convert to speech"`
	Pitch  int    `query:"pitch" minimum:"0" maximum:"100" default:"50" example:"50" doc:"Voice pitch (0-100)"`
	Rate   int    `query:"rate" minimum:"0" maximum:"100" default:"50" example:"50" doc:"Speech rate (0-100)"`
	Volume int    `query:"volume" minimum:"0" maximum:"100" default:"50" example:"50" doc:"Volume level (0-100)"`
}

type sayResponseDto struct {
	ContentType string `header:"Content-Type"`
	Body        []byte
}

var _ httpbase.Route[*sayRequestDto, *sayResponseDto] = (*say)(nil)

type SayOpts struct {
	fx.In

	Config  config.Config
	Service *tts.Service
}

func newSay(opts SayOpts) *say {
	return &say{
		config:  opts.Config,
		service: opts.Service,
	}
}

type say struct {
	config  config.Config
	service *tts.Service
}

func (s *say) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "tts-say",
		Method:      http.MethodGet,
		Path:        "/v1/tts/say",
		Tags:        []string{"TTS"},
		Summary:     "Text-to-Speech Say",
		Description: "Convert text to speech using the TTS service. Returns an audio file.",
		Responses: map[string]*huma.Response{
			"200": {
				Description: "Successful TTS conversion",
				Content: map[string]*huma.MediaType{
					"audio/wav": {
						Schema: &huma.Schema{
							Type:        "string",
							Format:      "binary",
							Description: "File content",
						},
					},
					"audio/mpeg": {
						Schema: &huma.Schema{
							Type:        "string",
							Format:      "binary",
							Description: "File content",
						},
					},
					"audio/mp3": {
						Schema: &huma.Schema{
							Type:        "string",
							Format:      "binary",
							Description: "File content",
						},
					},
				},
			},
		},
	}
}

func (s *say) Handler(
	ctx context.Context,
	input *sayRequestDto,
) (*sayResponseDto, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("http://%s/say", s.config.TTSServiceUrl))
	if err != nil {
		return nil, huma.NewError(http.StatusBadRequest, "Invalid TTS service URL", err)
	}

	query := reqUrl.Query()
	query.Set("voice", input.Voice)
	query.Set("pitch", strconv.Itoa(input.Pitch))
	query.Set("volume", strconv.Itoa(input.Volume))
	query.Set("rate", strconv.Itoa(input.Rate))
	query.Set("text", input.Text)
	reqUrl.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl.String(), nil)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to create request", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to call TTS service", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to read TTS response", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, huma.NewError(
			http.StatusBadGateway,
			fmt.Sprintf("TTS service returned error: %s", resp.Status),
		)
	}

	return &sayResponseDto{
		ContentType: resp.Header.Get("Content-Type"),
		Body:        body,
	}, nil
}

func (s *say) Register(api huma.API) {
	huma.Register(api, s.GetMeta(), s.Handler)
}
