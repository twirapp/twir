package chat_translations

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type detectedLangRequest struct {
	Text   string `json:"text"`
	Method string `json:"method,omitempty"`
}

type detectedLangResponse struct {
	Language   string  `json:"language"`
	Confidence float64 `json:"confidence"`
	Method     string  `json:"method"`
}

func (c *Service) detectLanguage(ctx context.Context, text string) (*detectedLangResponse, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = "http://language-processor:3000/detect"
	} else {
		reqUrl = "http://localhost:3012/detect"
	}

	requestBody := detectedLangRequest{
		Text:   text,
		Method: "mediapipe",
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(requestBody); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqUrl, body)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("cannot detect language: %s", string(responseBody))
	}

	var resp detectedLangResponse
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
