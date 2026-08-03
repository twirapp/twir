package mcp

import (
	"context"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	"github.com/twirapp/twir/apps/api-gql/internal/services/giveaways"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	channelsgiveaways "github.com/twirapp/twir/libs/entities/channels_giveaways"
	model "github.com/twirapp/twir/libs/gomodels"
)

type manageRewardInput struct {
	Action                            string `json:"action" jsonschema:"create, update, or delete"`
	ID                                string `json:"id,omitempty"`
	Title                             string `json:"title,omitempty"`
	Prompt                            string `json:"prompt,omitempty"`
	Cost                              int    `json:"cost,omitempty"`
	Enabled                           bool   `json:"enabled,omitempty"`
	BackgroundColor                   string `json:"backgroundColor,omitempty"`
	UserInputRequired                 bool   `json:"userInputRequired,omitempty"`
	MaxPerStreamEnabled               bool   `json:"maxPerStreamEnabled,omitempty"`
	MaxPerStream                      int    `json:"maxPerStream,omitempty"`
	MaxPerUserPerStreamEnabled        bool   `json:"maxPerUserPerStreamEnabled,omitempty"`
	MaxPerUserPerStream               int    `json:"maxPerUserPerStream,omitempty"`
	GlobalCooldownEnabled             bool   `json:"globalCooldownEnabled,omitempty"`
	GlobalCooldownSeconds             int    `json:"globalCooldownSeconds,omitempty"`
	ShouldRedemptionsSkipRequestQueue bool   `json:"shouldRedemptionsSkipRequestQueue,omitempty"`
}

type createGiveawayInput struct {
	Type                 string  `json:"type" jsonschema:"KEYWORD or ONLINE_CHATTERS"`
	Keyword              *string `json:"keyword,omitempty"`
	MinWatchedTime       *int64  `json:"minWatchedTime,omitempty"`
	MinMessages          *int32  `json:"minMessages,omitempty"`
	MinUsedChannelPoints *int64  `json:"minUsedChannelPoints,omitempty"`
	MinFollowDuration    *int64  `json:"minFollowDuration,omitempty"`
	RequireSubscription  bool    `json:"requireSubscription,omitempty"`
}

type updateGreetingInput struct {
	ID           string  `json:"id"`
	UserID       *string `json:"userId,omitempty"`
	Enabled      *bool   `json:"enabled,omitempty"`
	Text         *string `json:"text,omitempty"`
	IsReply      *bool   `json:"isReply,omitempty"`
	WithShoutOut *bool   `json:"withShoutOut,omitempty"`
}

type manageAlertInput struct {
	Action      string   `json:"action" jsonschema:"create, update, or delete"`
	ID          string   `json:"id,omitempty"`
	Name        *string  `json:"name,omitempty"`
	AudioID     *string  `json:"audioId,omitempty"`
	AudioVolume *int     `json:"audioVolume,omitempty"`
	CommandIDs  []string `json:"commandIds,omitempty"`
	RewardIDs   []string `json:"rewardIds,omitempty"`
	GreetingIDs []string `json:"greetingIds,omitempty"`
	KeywordIDs  []string `json:"keywordIds,omitempty"`
}

func (h *Handler) addEngagementTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_events", Description: "List configured automation events and operations for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Events.GetAll(ctx, channelID)
			return nil, map[string]any{"events": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_twir_events", Description: "List Twir event types that can be subscribed to or used by channel automations."},
		func(_ context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			return nil, map[string]any{"eventTypes": entity.AllEventType}, nil
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_rewards", Description: "List Twitch custom rewards for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			result, err := h.deps.Twitch.GetRewardsByChannelID(ctx, channelID)
			return nil, result, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "manage_rewards", Description: "Create, update, or delete a Twitch custom reward."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input manageRewardInput) (*modelsdk.CallToolResult, any, error) {
			result, err := h.deps.Twitch.ManageReward(ctx, channelID, twitchservice.ManageRewardInput{
				Action: input.Action, ID: input.ID, Title: input.Title, Prompt: input.Prompt,
				Cost: input.Cost, Enabled: input.Enabled, BackgroundColor: input.BackgroundColor,
				UserInputRequired: input.UserInputRequired, MaxPerStreamEnabled: input.MaxPerStreamEnabled,
				MaxPerStream: input.MaxPerStream, MaxPerUserPerStreamEnabled: input.MaxPerUserPerStreamEnabled,
				MaxPerUserPerStream: input.MaxPerUserPerStream, GlobalCooldownEnabled: input.GlobalCooldownEnabled,
				GlobalCooldownSeconds:             input.GlobalCooldownSeconds,
				ShouldRedemptionsSkipRequestQueue: input.ShouldRedemptionsSkipRequestQueue,
			})
			return nil, result, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_giveaways", Description: "List giveaways for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Giveaways.GiveawaysGetMany(ctx, channelID)
			return nil, map[string]any{"giveaways": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_giveaway", Description: "Create a keyword or online-chatters giveaway."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input createGiveawayInput) (*modelsdk.CallToolResult, any, error) {
			item, err := h.deps.Giveaways.Create(ctx, giveaways.CreateInput{
				ChannelID: channelID, CreatedByUserID: requestScope.ActorID,
				Type: channelsgiveaways.GiveawayType(input.Type), Keyword: input.Keyword,
				MinWatchedTime: input.MinWatchedTime, MinMessages: input.MinMessages,
				MinUsedChannelPoints: input.MinUsedChannelPoints, MinFollowDuration: input.MinFollowDuration,
				RequireSubscription: input.RequireSubscription,
			})
			return nil, item, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_greetings", Description: "List channel greetings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Greetings.GetManyByChannelID(ctx, channelID)
			return nil, map[string]any{"greetings": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_greeting", Description: "Update a channel greeting by UUID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateGreetingInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, err
			}
			item, err := h.deps.Greetings.Update(ctx, id, greetings.UpdateInput{ChannelID: channelID, ActorID: requestScope.ActorID, UserID: input.UserID, Enabled: input.Enabled, Text: input.Text, IsReply: input.IsReply, WithShoutOut: input.WithShoutOut})
			return nil, item, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_notifications", Description: "List global and channel-owner notifications."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			owner := ownerID(requestScope)
			var items []model.Notifications
			err := h.deps.Gorm.WithContext(ctx).Where(`"userId" = ? OR "userId" IS NULL`, owner).Order(`"createdAt" DESC`).Find(&items).Error
			return nil, map[string]any{"notifications": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_notification", Description: "Get one global or channel-owner notification by ID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
			owner := ownerID(requestScope)
			var item model.Notifications
			err := h.deps.Gorm.WithContext(ctx).Where(`id = ? AND ("userId" = ? OR "userId" IS NULL)`, input.ID, owner).First(&item).Error
			return nil, item, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_alerts", Description: "List channel alerts and their command/reward/greeting/keyword bindings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Alerts.GetManyByChannelID(ctx, channelID)
			return nil, map[string]any{"alerts": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "manage_alerts", Description: "Create, update, or delete a channel alert."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input manageAlertInput) (*modelsdk.CallToolResult, any, error) {
			switch input.Action {
			case "create":
				if input.Name == nil {
					return nil, nil, fmt.Errorf("name is required")
				}
				volume := 0
				if input.AudioVolume != nil {
					volume = *input.AudioVolume
				}
				item, err := h.deps.Alerts.Create(ctx, alerts.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Name: *input.Name, AudioID: input.AudioID, AudioVolume: volume, CommandIDS: input.CommandIDs, RewardIDS: input.RewardIDs, GreetingsIDS: input.GreetingIDs, KeywordsIDS: input.KeywordIDs})
				return nil, item, err
			case "update":
				id, err := parseID(input.ID)
				if err != nil {
					return nil, nil, err
				}
				item, err := h.deps.Alerts.Update(ctx, id, alerts.UpdateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Name: input.Name, AudioID: input.AudioID, AudioVolume: input.AudioVolume, CommandIDS: input.CommandIDs, RewardIDS: input.RewardIDs, GreetingsIDS: input.GreetingIDs, KeywordsIDS: input.KeywordIDs})
				return nil, item, err
			case "delete":
				id, err := parseID(input.ID)
				if err != nil {
					return nil, nil, err
				}
				err = h.deps.Alerts.Delete(ctx, id, channelID, requestScope.ActorID)
				return nil, map[string]bool{"deleted": err == nil}, err
			default:
				return nil, nil, fmt.Errorf("unsupported action %q", input.Action)
			}
		})
}
