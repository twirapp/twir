package vkvideoprobe

import "time"

const (
	defaultAPIBaseURL       = "https://api.live.vkvideo.ru"
	maxHTTPResponseBytes    = 1 << 20
	maxTransportSpecBytes   = 64 << 10
	defaultDuration         = 30 * time.Second
	maxDuration             = 5 * time.Minute
	defaultMaxFrameBytes    = 256 << 10
	maxMaxFrameBytes        = 1 << 20
	defaultMaxEvents        = 100
	maxMaxEvents            = 1_000
	defaultMaxOutputBytes   = 1 << 20
	maxMaxOutputBytes       = 10 << 20
	maxTransportFrameCount  = 32
	preflightRequestTimeout = 15 * time.Second
)

type PreflightResult struct {
	ChannelURL        string
	StreamID          string
	ChatChannel       string
	ConnectionToken   string
	SubscriptionToken string
	accessToken       string
}

type TransportSpec struct {
	URL          string   `json:"url"`
	Subprotocols []string `json:"subprotocols,omitempty"`
	Frames       [][]byte `json:"frames"`
}

type CaptureOptions struct {
	Duration       time.Duration
	MaxFrameBytes  int
	MaxEvents      int
	MaxOutputBytes int
}
