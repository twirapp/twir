package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
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
