package handle_message

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
)

type detectedLang struct {
	Language    string  `json:"language"`
	Probability float64 `json:"probability"`
}

type langDetectResult struct {
	Text              string         `json:"text"`
	CleanedText       string         `json:"cleaned_text"`
	DetectedLanguages []detectedLang `json:"detected_languages"`
}

func (c *Service) detectLanguage(ctx context.Context, text string) (*langDetectResult, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = fmt.Sprint("http://language-processor:3012/detect")
	} else {
		reqUrl = "http://localhost:3012/detect"
	}

	resp := langDetectResult{}
	res, err := req.R().SetContext(ctx).
		SetQueryParam("text", text).
		SetSuccessResult(&resp).
		Get(reqUrl)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot detect language: %s", res.String())
	}

	return &resp, nil
}
