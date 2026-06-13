package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	chatmessagesrepo "github.com/twirapp/twir/libs/repositories/chat_messages"
)

func resolveSelectedDashboardAnalyticsIdentity(ctx context.Context, deps Deps) (string, string, error) {
	dashboardID, err := deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return "", "", err
	}

	currentPlatform, err := deps.Sessions.GetCurrentPlatform(ctx)
	if err != nil {
		return "", "", err
	}

	parsedDashboardID, err := uuid.Parse(dashboardID)
	if err != nil {
		return "", "", err
	}

	channel, err := deps.ChannelsRepository.GetByID(ctx, parsedDashboardID)
	if err != nil {
		return "", "", err
	}

	platformChannelID, err := resolvePlatformChannelID(currentPlatform, channel)
	if err != nil {
		return "", "", err
	}

	return currentPlatform, platformChannelID, nil
}

func resolvePlatformChannelID(currentPlatform string, channel channelsmodel.Channel) (string, error) {
	switch currentPlatform {
	case "kick":
		if channel.KickPlatformID == nil || *channel.KickPlatformID == "" {
			return "", fmt.Errorf("kick platform channel id not found")
		}
		return *channel.KickPlatformID, nil
	case "twitch":
		if channel.TwitchPlatformID == nil || *channel.TwitchPlatformID == "" {
			return "", fmt.Errorf("twitch platform channel id not found")
		}
		return *channel.TwitchPlatformID, nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", currentPlatform)
	}
}

func resolveSelectedDashboardChatMessageTargets(
	ctx context.Context,
	deps Deps,
	platformFilter []string,
) ([]chatmessagesrepo.PlatformChannelIdentity, error) {
	dashboardID, err := deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	parsedDashboardID, err := uuid.Parse(dashboardID)
	if err != nil {
		return nil, err
	}

	channel, err := deps.ChannelsRepository.GetByID(ctx, parsedDashboardID)
	if err != nil {
		return nil, err
	}

	allowedPlatforms := make(map[string]struct{}, len(platformFilter))
	for _, platform := range platformFilter {
		allowedPlatforms[platform] = struct{}{}
	}

	targets := make([]chatmessagesrepo.PlatformChannelIdentity, 0, 2)
	if channel.TwitchPlatformID != nil && *channel.TwitchPlatformID != "" {
		if len(allowedPlatforms) == 0 {
			targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{Platform: "twitch", PlatformChannelID: *channel.TwitchPlatformID})
		} else if _, ok := allowedPlatforms["twitch"]; ok {
			targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{Platform: "twitch", PlatformChannelID: *channel.TwitchPlatformID})
		}
	}

	if channel.KickPlatformID != nil && *channel.KickPlatformID != "" {
		if len(allowedPlatforms) == 0 {
			targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{Platform: "kick", PlatformChannelID: *channel.KickPlatformID})
		} else if _, ok := allowedPlatforms["kick"]; ok {
			targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{Platform: "kick", PlatformChannelID: *channel.KickPlatformID})
		}
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no chat message targets found for selected dashboard")
	}

	return targets, nil
}
