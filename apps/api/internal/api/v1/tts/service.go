package tts

import (
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	cfg "github.com/satont/tsuwari/libs/config"
	"io"
	"net/url"
)

func handleGetSay(voice, pitch, volume, rate, text string) (io.ReadCloser, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)

	reqUrl, err := url.Parse(fmt.Sprintf("http://%s/say", config.TTSServiceUrl))
	if err != nil {
		return nil, err
	}

	query := reqUrl.Query()

	query.Set("voice", voice)
	query.Set("pitch", pitch)
	query.Set("volume", volume)
	query.Set("rate", rate)
	query.Set("text", text)

	reqUrl.RawQuery = query.Encode()

	response, err := req.Get(reqUrl.String())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return response.Body, nil
}
