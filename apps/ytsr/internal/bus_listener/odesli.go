package bus_listener

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type odesliPlatform struct {
	Url string `json:"url"`
}
type odesliPlatforms map[string]*odesliPlatform

type odesliResponse struct {
	PageUrl         string          `json:"pageUrl"`
	LinksByPlatform odesliPlatforms `json:"linksByPlatform"`
}

type odesliErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
}

func (c *YtsrServer) searchOdesli(ctx context.Context, query string) (odesliResponse, error) {
	result := odesliResponse{}

	reqUrl := fmt.Sprintf(
		"https://api.song.link/v1-alpha.1/links?url=%s&key=%s",
		query,
		c.config.OdesliApiKey,
	)
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return result, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != 200 {
		odesliError := odesliErrorResponse{}
		if err = json.Unmarshal(body, &odesliError); err != nil {
			return result, err
		}
		return result, fmt.Errorf(`odesli error for input "%s": %s`, query, odesliError.Code)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if _, ok := result.LinksByPlatform["youtube"]; !ok {
		return result, fmt.Errorf(`odesli error for input "%s": %s`, query, "no youtube link")
	}

	return result, nil
}
