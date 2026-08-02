package mcp

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
)

type uploadFileInput struct {
	Filename      string `json:"filename"`
	ContentType   string `json:"contentType" jsonschema:"image or audio MIME type"`
	ContentBase64 string `json:"contentBase64" jsonschema:"base64-encoded file bytes; maximum 10 MB"`
}

type manageQueueInput struct {
	Action   string   `json:"action" jsonschema:"delete, clear, or reorder"`
	VideoID  string   `json:"videoId,omitempty" jsonschema:"required for delete"`
	VideoIDs []string `json:"videoIds,omitempty" jsonschema:"complete ordered video ID list for reorder"`
}

type moderationInput struct {
	ID                              string   `json:"id,omitempty" jsonschema:"rule UUID; omit to create"`
	Name                            *string  `json:"name,omitempty"`
	Type                            string   `json:"type" jsonschema:"links, deny_list, symbols, long_message, caps, emotes, language, or one_man_spam"`
	Enabled                         bool     `json:"enabled"`
	BanTime                         int32    `json:"banTime,omitempty"`
	BanMessage                      string   `json:"banMessage,omitempty"`
	WarningMessage                  string   `json:"warningMessage,omitempty"`
	CheckClips                      bool     `json:"checkClips,omitempty"`
	TriggerLength                   int      `json:"triggerLength,omitempty"`
	MaxPercentage                   int      `json:"maxPercentage,omitempty"`
	DenyList                        []string `json:"denyList,omitempty"`
	DenyListRegexpEnabled           bool     `json:"denyListRegexpEnabled,omitempty"`
	DenyListWordBoundaryEnabled     bool     `json:"denyListWordBoundaryEnabled,omitempty"`
	DenyListSensitivityEnabled      bool     `json:"denyListSensitivityEnabled,omitempty"`
	DeniedChatLanguages             []string `json:"deniedChatLanguages,omitempty"`
	ExcludedRoles                   []string `json:"excludedRoles,omitempty"`
	MaxWarnings                     int      `json:"maxWarnings,omitempty"`
	OneManSpamMinimumStoredMessages int      `json:"oneManSpamMinimumStoredMessages,omitempty"`
	OneManSpamMessageMemorySeconds  int      `json:"oneManSpamMessageMemorySeconds,omitempty"`
	LanguageExcludedWords           []string `json:"languageExcludedWords,omitempty"`
}

func moderationServiceInput(channelID string, input moderationInput) channels_moderation_settings.CreateOrUpdateInput {
	return channels_moderation_settings.CreateOrUpdateInput{
		ChannelID: channelID, Name: input.Name, Type: entity.ModerationSettingsType(input.Type),
		Enabled: input.Enabled, BanTime: input.BanTime, BanMessage: input.BanMessage,
		WarningMessage: input.WarningMessage, CheckClips: input.CheckClips,
		TriggerLength: input.TriggerLength, MaxPercentage: input.MaxPercentage,
		DenyList: input.DenyList, DenyListRegexpEnabled: input.DenyListRegexpEnabled,
		DenyListWordBoundaryEnabled: input.DenyListWordBoundaryEnabled,
		DenyListSensitivityEnabled:  input.DenyListSensitivityEnabled,
		DeniedChatLanguages:         input.DeniedChatLanguages, ExcludedRoles: input.ExcludedRoles,
		MaxWarnings:                     input.MaxWarnings,
		OneManSpamMinimumStoredMessages: input.OneManSpamMinimumStoredMessages,
		OneManSpamMessageMemorySeconds:  input.OneManSpamMessageMemorySeconds,
		LanguageExcludedWords:           input.LanguageExcludedWords,
	}
}

func (h *Handler) addSystemTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_files", Description: "List uploaded image and audio files for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Files.GetMany(ctx, channelID)
			return nil, map[string]any{"files": items}, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "upload_file", Description: "Upload a base64-encoded image or audio file to channel storage."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input uploadFileInput) (*modelsdk.CallToolResult, any, error) {
			content, err := base64.StdEncoding.DecodeString(input.ContentBase64)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid base64 content: %w", err)
			}
			item, err := h.deps.Files.Upload(ctx, channelID, entity.Upload{File: bytes.NewReader(content), Filename: input.Filename, Size: int64(len(content)), ContentType: input.ContentType})
			return nil, item, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_games", Description: "List channel game modules and their settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			voteban, err := h.deps.Games.GetByChannelID(ctx, channelID)
			return nil, map[string]any{"games": []any{map[string]any{"name": "voteban", "settings": voteban}}}, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_song_requests", Description: "List the full song request queue for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.SongRequests.GetQueue(ctx, channelID)
			return nil, map[string]any{"queue": items}, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_current_song", Description: "Get the current song request playback state."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			item, err := h.deps.SongRequests.GetCurrentSong(ctx, channelID)
			return nil, map[string]any{"currentSong": item}, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "skip_song", Description: "Skip and remove the currently playing song."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			err := h.deps.SongRequests.Skip(ctx, channelID)
			return nil, map[string]bool{"skipped": err == nil}, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "manage_queue", Description: "Delete one song, clear the queue, or reorder the complete queue."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input manageQueueInput) (*modelsdk.CallToolResult, any, error) {
			var err error
			switch input.Action {
			case "delete":
				err = h.deps.SongRequests.DeleteFromQueue(ctx, channelID, input.VideoID)
			case "clear":
				err = h.deps.SongRequests.ClearQueue(ctx, channelID)
			case "reorder":
				err = h.deps.SongRequests.ReorderQueue(ctx, channelID, input.VideoIDs)
			default:
				err = fmt.Errorf("unsupported action %q", input.Action)
			}
			return nil, map[string]bool{"updated": err == nil}, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_moderation_settings", Description: "Get all moderation rules for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Moderation.GetByChannelID(ctx, channelID)
			return nil, map[string]any{"rules": items}, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "update_moderation", Description: "Create a moderation rule, or replace one when id is provided."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input moderationInput) (*modelsdk.CallToolResult, any, error) {
			serviceInput := moderationServiceInput(channelID, input)
			if input.ID == "" {
				item, err := h.deps.Moderation.Create(ctx, serviceInput)
				return nil, item, err
			}
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, err
			}
			existing, err := h.deps.Moderation.GetByID(ctx, id)
			if err != nil || existing.ChannelID != channelID {
				return nil, nil, fmt.Errorf("moderation rule not found")
			}
			item, err := h.deps.Moderation.Update(ctx, id, serviceInput)
			return nil, item, err
		})
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_mod_chat_wall", Description: "List moderation chat-wall entries and channel wall settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			walls, err := h.deps.ChatWall.GetChatWalls(ctx, channelID)
			if err != nil {
				return nil, nil, err
			}
			settings, settingsErr := h.deps.ChatWall.GetChannelSettings(ctx, channelID)
			if settingsErr != nil {
				settings = entity.ChatWallSettings{}
			}
			return nil, map[string]any{"walls": walls, "settings": settings}, nil
		})
}
