package handle_message

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
)

type translateRequest struct {
	Text          string   `json:"text"`
	SrcLang       string   `json:"src"`
	DestLang      string   `json:"dest"`
	ExcludedWords []string `json:"excluded_words,omitempty,omitzero"`
}

type translateResult struct {
	SourceLanguage      string   `json:"source_language"`
	SourceText          string   `json:"source_text"`
	TranslatedText      []string `json:"translated_text"`
	DestinationLanguage string   `json:"destination_language"`
}

func (c *Service) translate(
	ctx context.Context,
	input translateRequest,
) (*translateResult, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = fmt.Sprint("http://language-processor:3012/translate")
	} else {
		reqUrl = "http://localhost:3012/translate"
	}

	resp := translateResult{}
	res, err := req.R().
		SetContext(ctx).
		SetBody(input).
		SetHeader("Content-Type", "application/json").
		SetSuccessResult(&resp).
		Post(reqUrl)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot translate: %s", res.String())
	}

	return &resp, nil
}
