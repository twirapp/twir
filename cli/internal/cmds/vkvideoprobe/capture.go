package vkvideoprobe

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/coder/websocket"
)

var ErrOutputExists = errors.New("probe output already exists")
var errOutputLimit = errors.New("probe output limit reached")

type captureRecord struct {
	Direction string `json:"direction"`
	Encoding  string `json:"encoding"`
	Frame     string `json:"frame"`
}

func CreateOutput(path string) (*os.File, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if errors.Is(err, os.ErrExist) {
		return nil, ErrOutputExists
	}
	if err != nil {
		return nil, fmt.Errorf("create probe output: %w", err)
	}

	return file, nil
}

func Capture(ctx context.Context, spec TransportSpec, preflight PreflightResult, output *os.File, options CaptureOptions) (err error) {
	connection, _, err := websocket.Dial(ctx, spec.URL, &websocket.DialOptions{Subprotocols: spec.Subprotocols})
	if err != nil {
		return fmt.Errorf("dial probe websocket: %w", err)
	}
	defer func() {
		err = errors.Join(err, connection.Close(websocket.StatusNormalClosure, "probe complete"))
	}()
	connection.SetReadLimit(int64(options.MaxFrameBytes))

	for _, frame := range spec.RenderFrames(preflight) {
		if err := connection.Write(ctx, websocket.MessageText, frame); err != nil {
			return fmt.Errorf("write probe websocket frame: %w", err)
		}
	}

	captureContext, cancel := context.WithTimeout(ctx, options.Duration)
	defer cancel()
	secrets := []string{preflight.accessToken, preflight.ConnectionToken, preflight.SubscriptionToken}
	writer := boundedWriter{file: output, remaining: options.MaxOutputBytes}

	for eventCount := 0; eventCount < options.MaxEvents; eventCount++ {
		messageType, frame, readErr := connection.Read(captureContext)
		if errors.Is(readErr, context.DeadlineExceeded) {
			return nil
		}
		if readErr != nil {
			return fmt.Errorf("read probe websocket frame: %w", readErr)
		}

		record := newCaptureRecord(messageType, frame, secrets)
		if err := writer.write(record); errors.Is(err, errOutputLimit) {
			return nil
		} else if err != nil {
			return err
		}
	}

	return nil
}

func newCaptureRecord(messageType websocket.MessageType, frame []byte, secrets []string) captureRecord {
	redacted := RedactFrame(frame, secrets)
	if messageType == websocket.MessageText {
		return captureRecord{Direction: "inbound", Encoding: "utf-8", Frame: string(redacted)}
	}

	return captureRecord{Direction: "inbound", Encoding: "base64", Frame: base64.StdEncoding.EncodeToString(redacted)}
}

type boundedWriter struct {
	file      *os.File
	remaining int
}

func (writer *boundedWriter) write(record captureRecord) error {
	line, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode probe frame: %w", err)
	}
	line = append(line, '\n')
	if len(line) > writer.remaining {
		return errOutputLimit
	}
	if _, err := writer.file.Write(line); err != nil {
		return fmt.Errorf("write probe frame: %w", err)
	}
	writer.remaining -= len(line)

	return nil
}
