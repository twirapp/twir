package chat_translations

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("text", text)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot detect language: %s", string(body))
	}

	var resp langDetectResult
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
