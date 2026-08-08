package mcp

import (
	"context"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type getOverlayInput struct {
	Type string `json:"type" jsonschema:"custom, tts, dudes, kappagen, or be-right-back"`
	ID   string `json:"id,omitempty" jsonschema:"required for custom and dudes overlays"`
}

func (h *Handler) addOverlayTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_overlays", Description: "List custom, TTS, dudes, Kappagen, and be-right-back overlay settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			custom, err := h.deps.CustomOverlays.GetManyByChannelID(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			ttsSettings, err := h.deps.TTS.GetOrCreate(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			dudes, err := h.deps.Dudes.GetManyByChannelID(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			kappagenSettings, err := h.deps.Kappagen.GetOrCreate(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			brb, err := h.deps.BeRightBack.GetOrCreate(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			return nil, map[string]any{"custom": custom, "tts": ttsSettings, "dudes": dudes, "kappagen": kappagenSettings, "beRightBack": brb}, nil
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_overlay", Description: "Get one overlay or singleton overlay settings by type."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input getOverlayInput) (*modelsdk.CallToolResult, any, error) {
			switch input.Type {
			case "custom":
				id, err := parseID(input.ID)
				if err != nil {
					return nil, nil, err
				}
				item, err := h.deps.CustomOverlays.GetByID(ctx, id)
				if err != nil || item.ChannelID != channelID {
					return nil, nil, fmt.Errorf("overlay not found")
				}
				return nil, item, nil
			case "tts":
				item, err := h.deps.TTS.GetOrCreate(ctx, channelID)
				return nil, item, err
			case "dudes":
				id, err := parseID(input.ID)
				if err != nil {
					return nil, nil, err
				}
				items, err := h.deps.Dudes.GetManyByChannelID(ctx, channelID)
				if err != nil {
					return nil, nil, err
				}
				for _, item := range items {
					if item.ID == id {
						return nil, item, nil
					}
				}
				return nil, nil, fmt.Errorf("overlay not found")
			case "kappagen":
				item, err := h.deps.Kappagen.GetOrCreate(ctx, channelID)
				return nil, item, err
			case "be-right-back":
				item, err := h.deps.BeRightBack.GetOrCreate(ctx, channelID)
				return nil, item, err
			default:
				return nil, nil, fmt.Errorf("unsupported overlay type %q", input.Type)
			}
		})
}
