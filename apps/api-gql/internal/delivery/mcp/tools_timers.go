package mcp

import (
	"context"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/libs/bus-core/bots"
)

type timerResponseInput struct {
	Text          string `json:"text"`
	IsAnnounce    bool   `json:"isAnnounce,omitempty"`
	Count         int    `json:"count,omitempty" jsonschema:"relative selection weight; defaults to 1"`
	AnnounceColor string `json:"announceColor,omitempty"`
}

type createTimerInput struct {
	Name            string               `json:"name"`
	Enabled         bool                 `json:"enabled"`
	OfflineEnabled  bool                 `json:"offlineEnabled"`
	OnlineEnabled   bool                 `json:"onlineEnabled"`
	TimeInterval    int                  `json:"timeInterval" jsonschema:"interval in seconds"`
	MessageInterval int                  `json:"messageInterval" jsonschema:"minimum chat messages between sends"`
	Responses       []timerResponseInput `json:"responses"`
}

type updateTimerInput struct {
	ID              string               `json:"id"`
	Name            *string              `json:"name,omitempty"`
	Enabled         *bool                `json:"enabled,omitempty"`
	OfflineEnabled  *bool                `json:"offlineEnabled,omitempty"`
	OnlineEnabled   *bool                `json:"onlineEnabled,omitempty"`
	TimeInterval    *int                 `json:"timeInterval,omitempty"`
	MessageInterval *int                 `json:"messageInterval,omitempty"`
	Responses       []timerResponseInput `json:"responses,omitempty"`
}

func timerResponses(input []timerResponseInput) []timers.CreateResponse {
	result := make([]timers.CreateResponse, 0, len(input))
	for _, response := range input {
		result = append(result, timers.CreateResponse{Text: response.Text, IsAnnounce: response.IsAnnounce, Count: response.Count, AnnounceColor: announceColor(response.AnnounceColor)})
	}
	return result
}

func announceColor(value string) bots.AnnounceColor {
	switch value {
	case "primary":
		return bots.AnnounceColorPrimary
	case "blue":
		return bots.AnnounceColorBlue
	case "green":
		return bots.AnnounceColorGreen
	case "orange":
		return bots.AnnounceColorOrange
	case "purple":
		return bots.AnnounceColorPurple
	default:
		return bots.AnnounceColorRandom
	}
}

func (h *Handler) addTimerTools(s *modelsdk.Server, requestScope scope) {
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_timers", Description: "List all chat timers for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Timers.GetAllByChannelID(ctx, requestScope.Channel.ID.String())
			return nil, map[string]any{"timers": items}, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_timer", Description: "Create a chat timer for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input createTimerInput) (*modelsdk.CallToolResult, any, error) {
			result, err := h.deps.Timers.Create(ctx, timers.CreateInput{ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID, Name: input.Name, Enabled: input.Enabled, OfflineEnabled: input.OfflineEnabled, OnlineEnabled: input.OnlineEnabled, TimeInterval: input.TimeInterval, MessageInterval: input.MessageInterval, Responses: timerResponses(input.Responses)})
			return nil, result, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_timer", Description: "Update a chat timer by UUID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateTimerInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid timer id: %w", err)
			}
			result, err := h.deps.Timers.Update(ctx, timers.UpdateInput{ID: id, ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID, Name: input.Name, Enabled: input.Enabled, OfflineEnabled: input.OfflineEnabled, OnlineEnabled: input.OnlineEnabled, TimeInterval: input.TimeInterval, MessageInterval: input.MessageInterval, Responses: timerResponses(input.Responses)})
			return nil, result, err
		})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_timer", Description: "Delete a chat timer by UUID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid timer id: %w", err)
			}
			err = h.deps.Timers.Delete(ctx, id, requestScope.Channel.ID.String(), requestScope.ActorID)
			return nil, map[string]bool{"deleted": err == nil}, err
		})
}
