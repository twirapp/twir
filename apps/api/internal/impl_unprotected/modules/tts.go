package modules

import (
	"bytes"
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/satont/twir/libs/grpc/generated/api/modules_tts"
	"net/url"
	"strconv"
)

func (c *Modules) ModulesTTSSay(
	ctx context.Context,
	request *modules_tts.SayRequest,
) (*modules_tts.SayResponse, error) {
	reqUrl, err := url.Parse(fmt.Sprintf("http://%s/say", c.Config.TTSServiceUrl))
	if err != nil {
		return nil, err
	}

	query := reqUrl.Query()

	query.Set("voice", request.Voice)
	query.Set("pitch", strconv.Itoa(int(request.Pitch)))
	query.Set("volume", strconv.Itoa(int(request.Volume)))
	query.Set("rate", strconv.Itoa(int(request.Rate)))
	query.Set("text", request.Text)

	reqUrl.RawQuery = query.Encode()

	var b bytes.Buffer
	resp, err := req.SetContext(ctx).SetOutput(&b).Get(reqUrl.String())
	if err != nil {
		return nil, fmt.Errorf("cannot use say %w", err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot use say %s", resp.String())
	}

	return &modules_tts.SayResponse{
		File: b.Bytes(),
	}, nil
}
