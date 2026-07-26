package webhook

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/entities/channel"
	channelplatformsmodel "github.com/twirapp/twir/libs/entities/channel_platform"
	"github.com/twirapp/twir/libs/entities/platform"
)

func TestStartLogsVKTransportRegistrationState(t *testing.T) {
	var logBuf bytes.Buffer
	manager := &Manager{
		logger:       slog.New(slog.NewTextHandler(&logBuf, nil)),
		channelsRepo: &bulkChannelsRepo{},
		transports:   newTransportRegistry(),
	}

	if err := manager.start(context.Background()); err != nil {
		t.Fatalf("start returned error: %v", err)
	}

	if !strings.Contains(logBuf.String(), "vk_registered=false") {
		t.Fatalf("logs = %s", logBuf.String())
	}
}

func TestSubscribeAllPlatformsLogsVKTransportSkipWhenMissing(t *testing.T) {
	vkBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformVKVideoLive,
		Enabled:  true,
	}
	repo := &bulkChannelsRepo{
		channelsByPlatform: map[platform.Platform][]channelsmodel.Channel{
			platform.PlatformVKVideoLive: {
				{
					ID: uuid.New(),
					Bindings: []channelplatformsmodel.ChannelPlatform{
						vkBinding,
					},
				},
			},
		},
	}
	var logBuf bytes.Buffer
	manager := &Manager{
		logger:       slog.New(slog.NewTextHandler(&logBuf, nil)),
		channelsRepo: repo,
		transports:   newTransportRegistry(),
	}

	if err := manager.subscribeAllPlatforms(context.Background()); err != nil {
		t.Fatalf("subscribeAllPlatforms returned error: %v", err)
	}

	logs := logBuf.String()
	for _, want := range []string{
		"webhook manager: VK EventSub subscribe summary",
		"transport_registered=false",
		"eligible_bindings=1",
		"routed_bindings=0",
		"skipped_bindings=1",
		"skip_reason=transport_not_registered",
	} {
		if !strings.Contains(logs, want) {
			t.Fatalf("logs = %s", logs)
		}
	}
}

func TestUnsubscribeAllPlatformsLogsVKTransportSkipWhenMissing(t *testing.T) {
	vkBinding := channelplatformsmodel.ChannelPlatform{
		ID:       uuid.New(),
		Platform: platform.PlatformVKVideoLive,
		Enabled:  true,
	}
	repo := &bulkChannelsRepo{
		channelsByPlatform: map[platform.Platform][]channelsmodel.Channel{
			platform.PlatformVKVideoLive: {
				{
					ID: uuid.New(),
					Bindings: []channelplatformsmodel.ChannelPlatform{
						vkBinding,
					},
				},
			},
		},
	}
	var logBuf bytes.Buffer
	manager := &Manager{
		logger:       slog.New(slog.NewTextHandler(&logBuf, nil)),
		channelsRepo: repo,
		transports:   newTransportRegistry(),
	}

	if err := manager.unsubscribeAllPlatforms(context.Background()); err != nil {
		t.Fatalf("unsubscribeAllPlatforms returned error: %v", err)
	}

	logs := logBuf.String()
	for _, want := range []string{
		"webhook manager: VK EventSub unsubscribe summary",
		"transport_registered=false",
		"eligible_bindings=1",
		"routed_bindings=0",
		"skipped_bindings=1",
		"skip_reason=transport_not_registered",
	} {
		if !strings.Contains(logs, want) {
			t.Fatalf("logs = %s", logs)
		}
	}
}
