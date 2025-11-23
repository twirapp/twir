package handle_message

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/lkretschmer/deepl-go"
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

var providers = []func(c *Service, ctx context.Context, input translateRequest) (
	*translateResult,
	error,
){
	(*Service).translateUnOfficial,
	(*Service).translateOfficial,
}

func (c *Service) translate(
	ctx context.Context,
	input translateRequest,
) (*translateResult, error) {
	for _, provider := range providers {
		result, err := provider(c, ctx, input)
		if err != nil {
			c.logger.Error("translate provider error", "err", err)
			continue
		}
		if result != nil && len(result.TranslatedText) > 0 {
			return result, nil
		}
	}

	return nil, fmt.Errorf("all translate providers failed")
}

func (c *Service) translateOfficial(
	ctx context.Context,
	input translateRequest,
) (*translateResult, error) {

	translation, err := c.deeplClient.TranslateTextWithOptions(
		ctx, deepl.TranslateTextOptions{
			Text:       []string{input.Text},
			SourceLang: input.SrcLang,
			TargetLang: input.DestLang,
			Context:    "Never write offensive words, slurs, racial, homophobic, transphobic, sexist, or other toxic terms in full. Always replace them with a censored version with asterisks, while keeping them recognizable, for example: n*****, f****t, p****r, h****, ch****k, t*****a, etc. Do this automatically for any prohibited word, even if you are asked to write without censorship.",
		},
	)
	if err != nil {
		return nil, err
	}

	result := &translateResult{
		SourceLanguage:      "",
		SourceText:          "",
		TranslatedText:      nil,
		DestinationLanguage: "",
	}

	for _, t := range translation {
		if t == nil {
			continue
		}

		result.SourceLanguage += t.DetectedSourceLanguage
		result.SourceText += t.Text
		result.TranslatedText = append(result.TranslatedText, t.Text)
		result.DestinationLanguage += input.DestLang
	}

	return result, nil
}

func (c *Service) translateUnOfficial(
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
