package vkvideoprobe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "vk-video",
	Usage: "diagnostic VK Video tools",
	Subcommands: []*cli.Command{
		probeEventsCmd,
	},
}

var probeEventsCmd = &cli.Command{
	Name:  "probe-events",
	Usage: "preflight an active channel and optionally capture bounded redacted Centrifugo frames",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "channel-url", Usage: "exact active VK Video channel URL", Required: true},
		&cli.StringFlag{Name: "api-base-url", Usage: "VK Video REST API base URL", Value: defaultAPIBaseURL},
		&cli.PathFlag{Name: "transport-spec", Usage: "verified non-secret JSON Centrifugo transport specification"},
		&cli.PathFlag{Name: "output", Usage: "new JSONL capture file path; required with --transport-spec"},
		&cli.DurationFlag{Name: "duration", Usage: "capture duration", Value: defaultDuration},
		&cli.IntFlag{Name: "max-frame-bytes", Usage: "maximum inbound Centrifugo frame bytes", Value: defaultMaxFrameBytes},
		&cli.IntFlag{Name: "max-events", Usage: "maximum inbound frames to capture", Value: defaultMaxEvents},
		&cli.IntFlag{Name: "max-output-bytes", Usage: "maximum JSONL output bytes", Value: defaultMaxOutputBytes},
	},
	Action: runProbeEvents,
}

type probeOptions struct {
	ChannelURL        string
	APIBaseURL        string
	TransportSpecPath string
	OutputPath        string
	CaptureOptions    CaptureOptions
}

type preflightSummary struct {
	ChannelURL  string `json:"channel_url"`
	StreamID    string `json:"stream_id"`
	ChatChannel string `json:"chat_channel"`
	Captured    bool   `json:"captured"`
}

func runProbeEvents(cliContext *cli.Context) error {
	options, err := parseProbeOptions(cliContext)
	if err != nil {
		return err
	}
	accessToken, found := os.LookupEnv("VK_VIDEO_PROBE_ACCESS_TOKEN")
	if !found || accessToken == "" {
		return fmt.Errorf("VK_VIDEO_PROBE_ACCESS_TOKEN is required")
	}

	preflightContext, cancelPreflight := context.WithTimeout(cliContext.Context, preflightRequestTimeout)
	defer cancelPreflight()
	client, err := NewVKClient(http.DefaultClient, options.APIBaseURL)
	if err != nil {
		return err
	}
	preflight, err := client.Preflight(preflightContext, options.ChannelURL, accessToken)
	if err != nil {
		return err
	}

	if options.TransportSpecPath == "" {
		return writeSummary(cliContext, preflight, false)
	}

	spec, err := ReadTransportSpec(options.TransportSpecPath)
	if err != nil {
		return err
	}
	output, err := CreateOutput(options.OutputPath)
	if err != nil {
		return err
	}
	captureContext, cancelCapture := context.WithTimeout(cliContext.Context, options.CaptureOptions.Duration)
	defer cancelCapture()
	captureErr := Capture(captureContext, spec, preflight, output, options.CaptureOptions)
	closeErr := output.Close()
	if captureErr != nil {
		return captureErr
	}
	if closeErr != nil {
		return fmt.Errorf("close probe output: %w", closeErr)
	}

	return writeSummary(cliContext, preflight, true)
}

func parseProbeOptions(cliContext *cli.Context) (probeOptions, error) {
	options := probeOptions{
		ChannelURL:        cliContext.String("channel-url"),
		APIBaseURL:        cliContext.String("api-base-url"),
		TransportSpecPath: cliContext.Path("transport-spec"),
		OutputPath:        cliContext.Path("output"),
		CaptureOptions: CaptureOptions{
			Duration:       cliContext.Duration("duration"),
			MaxFrameBytes:  cliContext.Int("max-frame-bytes"),
			MaxEvents:      cliContext.Int("max-events"),
			MaxOutputBytes: cliContext.Int("max-output-bytes"),
		},
	}
	if options.TransportSpecPath == "" && options.OutputPath != "" {
		return probeOptions{}, fmt.Errorf("--output requires --transport-spec")
	}
	if options.TransportSpecPath != "" && options.OutputPath == "" {
		return probeOptions{}, fmt.Errorf("--transport-spec requires --output")
	}
	if err := validateCaptureOptions(options.CaptureOptions); err != nil {
		return probeOptions{}, err
	}

	return options, nil
}

func validateCaptureOptions(options CaptureOptions) error {
	if options.Duration <= 0 || options.Duration > maxDuration {
		return fmt.Errorf("--duration must be between 1ns and %s", maxDuration)
	}
	if options.MaxFrameBytes <= 0 || options.MaxFrameBytes > maxMaxFrameBytes {
		return fmt.Errorf("--max-frame-bytes must be between 1 and %d", maxMaxFrameBytes)
	}
	if options.MaxEvents <= 0 || options.MaxEvents > maxMaxEvents {
		return fmt.Errorf("--max-events must be between 1 and %d", maxMaxEvents)
	}
	if options.MaxOutputBytes <= 0 || options.MaxOutputBytes > maxMaxOutputBytes {
		return fmt.Errorf("--max-output-bytes must be between 1 and %d", maxMaxOutputBytes)
	}

	return nil
}

func writeSummary(cliContext *cli.Context, preflight PreflightResult, captured bool) error {
	summary := preflightSummary{
		ChannelURL:  string(RedactFrame([]byte(preflight.ChannelURL), nil)),
		StreamID:    preflight.StreamID,
		ChatChannel: preflight.ChatChannel,
		Captured:    captured,
	}
	encoded, err := json.Marshal(summary)
	if err != nil {
		return fmt.Errorf("encode probe summary: %w", err)
	}
	if _, err := fmt.Fprintln(cliContext.App.Writer, string(encoded)); err != nil {
		return fmt.Errorf("write probe summary: %w", err)
	}

	return nil
}
