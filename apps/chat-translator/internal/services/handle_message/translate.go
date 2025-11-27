package handle_message

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	googletranslate "cloud.google.com/go/translate"
	"github.com/lkretschmer/deepl-go"
	"golang.org/x/text/language"
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
	for _, p := range providers {
		result, err := p(c, ctx, input)
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

func (c *Service) translateDeeplOfficial(
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

func (c *Service) translateDeeplUnOfficial(
	ctx context.Context,
	input translateRequest,
) (*translateResult, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = "http://language-processor:3012/translate"
	} else {
		reqUrl = "http://localhost:3012/translate"
	}

	bodyBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	respBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot translate: %s", string(respBody))
	}

	var resp translateResult
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Service) translateGoogleOfficial(
	ctx context.Context,
	input translateRequest,
) (*translateResult, error) {
	sourceLang, err := language.Parse(input.SrcLang)
	if err != nil {
		return nil, err
	}
	targetLang, err := language.Parse(input.DestLang)
	if err != nil {
		return nil, err
	}

	result, err := c.googleTranslateClient.Translate(
		ctx, []string{input.Text}, targetLang, &googletranslate.Options{
			Source: sourceLang,
		},
	)

	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}

	return &translateResult{
		SourceLanguage:      result[0].Source.String(),
		SourceText:          input.Text,
		TranslatedText:      []string{result[0].Text},
		DestinationLanguage: targetLang.String(),
	}, nil
}
