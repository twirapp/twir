package modules

import (
	"context"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/satont/twir/libs/grpc/generated/api/modules_tts"
	"google.golang.org/protobuf/types/known/emptypb"
	"net/url"
	"strconv"
)

func (c *Modules) ModulesTTSSay(
	ctx context.Context,
	request *modules_tts.SayRequest,
) (*emptypb.Empty, error) {
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

	resp, err := req.SetContext(ctx).Get(reqUrl.String())
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot use say %s", resp.String())
	}

	return &emptypb.Empty{}, nil
}
