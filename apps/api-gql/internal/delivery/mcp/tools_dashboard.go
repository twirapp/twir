package mcp

import (
	"context"
	"fmt"
	"strings"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type updateBotSettingsInput struct {
	Platform string `json:"platform" jsonschema:"platform binding to update: twitch, kick, vk_video_live, or youtube"`
	Enabled  bool   `json:"enabled"`
}

type listCommunityUsersInput struct {
	Search  string `json:"search,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"perPage,omitempty"`
}

type communityUser struct {
	UserID            string `json:"userId"`
	Platform          string `json:"platform"`
	PlatformID        string `json:"platformId"`
	Login             string `json:"login"`
	DisplayName       string `json:"displayName"`
	Avatar            string `json:"avatar,omitempty"`
	Messages          int    `json:"messages"`
	WatchedMs         int64  `json:"watchedMs"`
	UsedEmotes        int    `json:"usedEmotes"`
	UsedChannelPoints int64  `json:"usedChannelPoints"`
	IsMod             bool   `json:"isMod"`
	IsVIP             bool   `json:"isVip"`
	IsSubscriber      bool   `json:"isSubscriber"`
}

type communityUserInfoInput struct {
	UserID string `json:"userId" jsonschema:"internal Twir user UUID returned by list_community_users"`
}

func (h *Handler) addDashboardTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_stats", Description: "Get stream and channel dashboard statistics."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		stats, err := h.deps.Dashboard.GetDashboardStats(ctx, channelID)
		return nil, stats, err
	})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_bot_settings", Description: "Get bot status and enabled state for every connected channel platform."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		statuses, err := h.deps.Dashboard.GetBotStatuses(ctx, channelID)
		return nil, map[string]any{"platforms": statuses}, err
	})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "update_bot_settings", Description: "Enable or disable the bot on one connected channel platform."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateBotSettingsInput) (*modelsdk.CallToolResult, any, error) {
		platform := platformentity.Platform(strings.ToLower(input.Platform))
		binding, err := h.deps.ChannelPlatforms.SetEnabled(ctx, requestScope.Channel.ID, platform, input.Enabled)
		return nil, binding, err
	})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_community_users", Description: "List users who have channel statistics, with optional search and pagination."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input listCommunityUsersInput) (*modelsdk.CallToolResult, any, error) {
		if input.Page < 0 {
			input.Page = 0
		}
		if input.PerPage < 1 {
			input.PerPage = 20
		}
		if input.PerPage > 100 {
			input.PerPage = 100
		}

		query := h.deps.Gorm.WithContext(ctx).
			Table("users_stats").
			Select(`users_stats.user_id, users.platform, users.platform_id, users.login, users.display_name, users.avatar, users_stats.messages, users_stats.watched, users_stats.emotes, users_stats."usedChannelPoints" AS used_channel_points, users_stats.is_mod, users_stats.is_vip, users_stats.is_subscriber`).
			Joins("JOIN users ON users.id = users_stats.user_id").
			Where("users_stats.channel_id = ?::uuid", channelID)
		if search := strings.TrimSpace(input.Search); search != "" {
			pattern := "%" + search + "%"
			query = query.Where("users.login ILIKE ? OR users.display_name ILIKE ? OR users.platform_id ILIKE ?", pattern, pattern, pattern)
		}

		var total int64
		if err := query.Count(&total).Error; err != nil {
			return nil, nil, fmt.Errorf("count community users: %w", err)
		}

		items := make([]communityUser, 0, input.PerPage)
		if err := query.
			Order(`users_stats.watched DESC`).
			Limit(input.PerPage).
			Offset(input.Page * input.PerPage).
			Scan(&items).Error; err != nil {
			return nil, nil, fmt.Errorf("list community users: %w", err)
		}

		return nil, map[string]any{"users": items, "total": total}, nil
	})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_user_info", Description: "Get channel-specific statistics and follow state for a community user."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input communityUserInfoInput) (*modelsdk.CallToolResult, any, error) {
		info, err := h.deps.Users.GetChannelUserInfo(ctx, users.ChannelUserInfoInput{ChannelID: channelID, UserID: input.UserID})
		return nil, info, err
	})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_scheduled_vips", Description: "List VIP removals scheduled for this channel."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.ScheduledVIPs.GetScheduledVips(ctx, channelID)
		return nil, map[string]any{"scheduledVips": items}, err
	})
}
