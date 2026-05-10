package bus_listener

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const expireShift = 15 * time.Minute

func isTokenExpired(expiresIn int, obtainmentTimestamp time.Time) bool {
	currentTime := time.Now().UTC()
	currentTokenLiveTime := currentTime.Sub(obtainmentTimestamp.UTC())

	return int64(currentTokenLiveTime.Seconds())+int64(expireShift.Seconds()) >= int64(expiresIn)
}

func decodeJsonResponse(resp *http.Response, target any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("unmarshal response body: %w", err)
	}

	return nil
}

func boolPtr(v bool) *bool {
	return &v
}
