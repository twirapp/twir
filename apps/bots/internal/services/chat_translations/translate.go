package chat_translations

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"
	"unicode"

	googletranslate "cloud.google.com/go/translate"
	"github.com/lkretschmer/deepl-go"
	kvoptions "github.com/twirapp/kv/options"
	"golang.org/x/text/language"
)

var translationURLPattern = regexp.MustCompile(`(?i)(?:https?://|www\.)\S+`)

func prepareTextForTranslation(text string) string {
	text = translationURLPattern.ReplaceAllString(text, " ")
	return strings.Join(strings.Fields(text), " ")
}

func hasTranslatableText(text string) bool {
	letters := 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			letters++
		}
	}

	return letters >= 3
}

// translationsEquivalent reports whether a translation is a near-copy of the
// source, which happens when the source language was misdetected.
func translationsEquivalent(source, translated string) bool {
	src := lettersOnlyRunes(source)
	dst := lettersOnlyRunes(translated)

	if len(src) == 0 || len(dst) == 0 {
		return false
	}

	distance := levenshteinDistance(src, dst)
	similarity := 1 - float64(distance)/float64(max(len(src), len(dst)))

	return similarity >= 0.8
}

func lettersOnlyRunes(text string) []rune {
	runes := make([]rune, 0, len(text))
	for _, r := range strings.ToLower(text) {
		if unicode.IsLetter(r) {
			runes = append(runes, r)
		}
	}
	return runes
}

func levenshteinDistance(s1, s2 []rune) int {
	previous := make([]int, len(s2)+1)
	for j := range previous {
		previous[j] = j
	}

	for i := 1; i <= len(s1); i++ {
		current := make([]int, len(s2)+1)
		current[0] = i
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}
			current[j] = min(
				previous[j-1]+cost,
				previous[j]+1,
				current[j-1]+1,
			)
		}
		previous = current
	}

	return previous[len(s2)]
}

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
	start := time.Now()

	cacheKey := fmt.Sprintf(
		"chat-translator:v2:%s",
		base64.StdEncoding.EncodeToString(
			[]byte(fmt.Sprintf(
				"%s:%s:%s",
				input.SrcLang,
				input.DestLang,
				input.Text,
			)),
		),
	)

	var (
		r         *translateResult
		fromCache = false
	)

	if cachedValue, err := c.kv.Get(ctx, cacheKey).String(); err == nil && len(cachedValue) > 0 {
		r = &translateResult{
			SourceLanguage:      input.SrcLang,
			SourceText:          input.Text,
			TranslatedText:      []string{cachedValue},
			DestinationLanguage: input.DestLang,
		}
		fromCache = true
	} else {
		for _, p := range providers {
			result, err := p(c, ctx, input)
			if err != nil {
				c.logger.Error("translate provider error", "err", err)
				continue
			}
			if result != nil && len(result.TranslatedText) > 0 {
				r = result
				break
			}
		}
	}

	if r != nil && len(r.TranslatedText) > 0 {
		if !fromCache {
			err := c.kv.Set(
				ctx,
				cacheKey,
				strings.Join(r.TranslatedText, " "),
				kvoptions.WithExpire(24*7*time.Hour),
			)
			if err != nil {
				c.logger.Error("failed to cache translation", "err", err)
			}
		}

		c.logger.Info(
			"translated text",
			slog.String("src_lang", input.SrcLang),
			slog.String("dest_lang", input.DestLang),
			slog.String("original_text", input.Text),
			slog.String("translated_text", strings.Join(r.TranslatedText, " ")),
			slog.Bool("from_cache", fromCache),
			slog.Int64("duration_ms", time.Since(start).Milliseconds()),
		)

		return r, nil
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
