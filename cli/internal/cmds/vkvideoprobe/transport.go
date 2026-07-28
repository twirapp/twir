package vkvideoprobe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var placeholderPattern = regexp.MustCompile(`{{[^{}]+}}`)
var sensitiveFieldNamePattern = regexp.MustCompile(`(?i)^(?:[a-z0-9_-]*(?:token|secret|password|authorization|credential)[a-z0-9_-]*|api[_-]?key|cookie)$`)

var allowedPlaceholders = map[string]struct{}{
	"{{connection_token}}":   {},
	"{{channel}}":            {},
	"{{subscription_token}}": {},
}

type rawTransportSpec struct {
	URL          string            `json:"url"`
	Subprotocols []string          `json:"subprotocols"`
	Frames       []json.RawMessage `json:"frames"`
}

func ReadTransportSpec(path string) (TransportSpec, error) {
	file, err := os.Open(path)
	if err != nil {
		return TransportSpec{}, fmt.Errorf("open transport spec: %w", err)
	}
	defer file.Close()

	contents, err := io.ReadAll(io.LimitReader(file, maxTransportSpecBytes+1))
	if err != nil {
		return TransportSpec{}, fmt.Errorf("read transport spec: %w", err)
	}
	if len(contents) > maxTransportSpecBytes {
		return TransportSpec{}, fmt.Errorf("transport spec exceeded %d bytes", maxTransportSpecBytes)
	}

	return ParseTransportSpec(contents)
}

func ParseTransportSpec(contents []byte) (TransportSpec, error) {
	decoder := json.NewDecoder(bytes.NewReader(contents))
	decoder.DisallowUnknownFields()

	var raw rawTransportSpec
	if err := decoder.Decode(&raw); err != nil {
		return TransportSpec{}, fmt.Errorf("decode transport spec: %w", err)
	}
	if err := ensureEOF(decoder); err != nil {
		return TransportSpec{}, err
	}

	parsedURL, err := url.Parse(raw.URL)
	if err != nil || parsedURL.Scheme != "wss" || parsedURL.Host == "" || parsedURL.User != nil {
		return TransportSpec{}, fmt.Errorf("transport spec URL must be an explicit wss URL without user info")
	}
	for key := range parsedURL.Query() {
		if sensitiveFieldNamePattern.MatchString(key) {
			return TransportSpec{}, fmt.Errorf("transport spec URL must not contain sensitive query values")
		}
	}
	if len(raw.Frames) == 0 || len(raw.Frames) > maxTransportFrameCount {
		return TransportSpec{}, fmt.Errorf("transport spec must contain 1 to %d frames", maxTransportFrameCount)
	}

	for _, subprotocol := range raw.Subprotocols {
		if strings.TrimSpace(subprotocol) == "" || strings.ContainsAny(subprotocol, "\r\n") {
			return TransportSpec{}, fmt.Errorf("transport spec contains an invalid subprotocol")
		}
	}
	for _, frame := range raw.Frames {
		if !json.Valid(frame) {
			return TransportSpec{}, fmt.Errorf("transport spec frame is not valid JSON")
		}
		if err := validatePlaceholders(frame); err != nil {
			return TransportSpec{}, err
		}
		if err := validateNonSecretFrame(frame); err != nil {
			return TransportSpec{}, err
		}
	}

	frames := make([][]byte, len(raw.Frames))
	for index, frame := range raw.Frames {
		frames[index] = bytes.Clone(frame)
	}

	return TransportSpec{URL: raw.URL, Subprotocols: raw.Subprotocols, Frames: frames}, nil
}

func (spec TransportSpec) RenderFrames(preflight PreflightResult) [][]byte {
	replacer := strings.NewReplacer(
		"{{connection_token}}", preflight.ConnectionToken,
		"{{channel}}", preflight.ChatChannel,
		"{{subscription_token}}", preflight.SubscriptionToken,
	)

	frames := make([][]byte, len(spec.Frames))
	for index, frame := range spec.Frames {
		frames[index] = []byte(replacer.Replace(string(frame)))
	}

	return frames
}

func ensureEOF(decoder *json.Decoder) error {
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return fmt.Errorf("transport spec contains multiple JSON values")
		}
		return fmt.Errorf("read transport spec: %w", err)
	}

	return nil
}

func validatePlaceholders(frame []byte) error {
	for _, placeholder := range placeholderPattern.FindAllString(string(frame), -1) {
		if _, allowed := allowedPlaceholders[placeholder]; !allowed {
			return fmt.Errorf("transport spec contains an unsupported placeholder")
		}
	}

	return nil
}

func validateNonSecretFrame(frame []byte) error {
	var value any
	if err := json.Unmarshal(frame, &value); err != nil {
		return fmt.Errorf("decode transport spec frame: %w", err)
	}

	return validateFrameValue(value)
}

func validateFrameValue(value any) error {
	switch value := value.(type) {
	case map[string]any:
		for key, nestedValue := range value {
			if sensitiveFieldNamePattern.MatchString(key) {
				stringValue, ok := nestedValue.(string)
				if !ok || !containsAllowedPlaceholder(stringValue) {
					return fmt.Errorf("transport spec contains a static sensitive value")
				}
			}
			if err := validateFrameValue(nestedValue); err != nil {
				return err
			}
		}
	case []any:
		for _, nestedValue := range value {
			if err := validateFrameValue(nestedValue); err != nil {
				return err
			}
		}
	case string:
		if err := validateFrameURL(value); err != nil {
			return err
		}
	}

	return nil
}

func containsAllowedPlaceholder(value string) bool {
	for placeholder := range allowedPlaceholders {
		if strings.Contains(value, placeholder) {
			return true
		}
	}

	return false
}

func validateFrameURL(value string) error {
	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil
	}
	for key, values := range parsedURL.Query() {
		if sensitiveFieldNamePattern.MatchString(key) {
			for _, queryValue := range values {
				if !containsAllowedPlaceholder(queryValue) {
					return fmt.Errorf("transport spec contains a static sensitive URL value")
				}
			}
		}
	}

	return nil
}
