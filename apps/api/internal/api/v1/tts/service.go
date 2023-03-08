package tts

import (
	"fmt"
	"io"
	"net/url"

	"github.com/imroc/req/v3"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func handleGetSay(services *types.Services, voice, pitch, volume, rate, text string) (io.ReadCloser, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("http://%s/say", services.Config.TTSServiceUrl))
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
