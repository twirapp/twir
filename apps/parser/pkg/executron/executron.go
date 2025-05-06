package executron

import (
	"context"
	"fmt"
	"net/url"

	"github.com/imroc/req/v3"
	config "github.com/satont/twir/libs/config"
)

type request struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

type response struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func New(cfg config.Config) Executron {
	return Executron{
		apiUrl: cfg.ExecutronAddr,
	}
}

type Executron struct {
	apiUrl string
}

func (c *Executron) ExecuteUserCode(ctx context.Context, language, code string) (string, error) {
	u, _ := url.Parse(c.apiUrl)
	u.Path = "/run"

	var executeResponse response
	resp, err := req.R().
		SetContext(ctx).
		SetBodyJsonMarshal(
			request{
				Language: language,
				Code:     code,
			},
		).
		SetSuccessResult(&executeResponse).
		Post(u.String())
	if err != nil {
		return "", err
	}
	if !resp.IsSuccessState() {
		return "", fmt.Errorf("cannot execute code: %s", resp.String())
	}

	if executeResponse.Error != "" {
		return "", fmt.Errorf("cannot execute code: %s", executeResponse.Error)
	}

	return executeResponse.Result, nil
}
