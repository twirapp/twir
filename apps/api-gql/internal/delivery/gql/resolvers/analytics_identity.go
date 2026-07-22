package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/channelbinding"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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

	channel, err := deps.ChannelService.GetChannelByID(ctx, parsedDashboardID)
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
	platform := platformentity.Platform(currentPlatform)
	if !platform.IsValid() {
		return "", fmt.Errorf("unsupported platform: %s", currentPlatform)
	}

	binding, found := channelbinding.Find(channel, platform)
	if !found || binding.PlatformChannelID == "" {
		return "", fmt.Errorf("%s platform channel id not found", currentPlatform)
	}

	return binding.PlatformChannelID, nil
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

	channel, err := deps.ChannelService.GetChannelByID(ctx, parsedDashboardID)
	if err != nil {
		return nil, err
	}

	allowedPlatforms := make(map[string]struct{}, len(platformFilter))
	for _, platform := range platformFilter {
		allowedPlatforms[platform] = struct{}{}
	}

	targets := make([]chatmessagesrepo.PlatformChannelIdentity, 0, len(channel.Bindings))
	for _, binding := range channel.Bindings {
		if binding.PlatformChannelID == "" {
			continue
		}
		if len(allowedPlatforms) > 0 {
			if _, ok := allowedPlatforms[binding.Platform.String()]; !ok {
				continue
			}
		}

		targets = append(targets, chatmessagesrepo.PlatformChannelIdentity{
			Platform:          binding.Platform.String(),
			PlatformChannelID: binding.PlatformChannelID,
		})
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no chat message targets found for selected dashboard")
	}

	return targets, nil
}
