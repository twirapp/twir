package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_secret"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinksua "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents"
)

type secretOutput struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Value       *string `json:"value,omitempty"`
}
type createSecretInput struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Value       string  `json:"value"`
}
type updateSecretInput struct {
	ID          string  `json:"id"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Value       *string `json:"value,omitempty"`
}

func secretMetadata(id, name string, description *string) secretOutput {
	return secretOutput{ID: id, Name: name, Description: description}
}

func (h *Handler) ownedSecret(ctx context.Context, requestScope scope, id string) (secretOutput, error) {
	parsed, err := parseID(id)
	if err != nil {
		return secretOutput{}, err
	}
	item, err := h.deps.Secrets.GetByID(ctx, parsed)
	if err != nil {
		return secretOutput{}, err
	}
	if item.IsNil() || item.ChannelID != requestScope.Channel.ID {
		return secretOutput{}, channels_secret.ErrNotFound
	}
	return secretMetadata(item.ID.String(), item.Name, item.Description.Ptr()), nil
}

func (h *Handler) addSecretTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_secrets", Description: "List secret metadata. Values are never included in list results."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Secrets.GetAllByChannelID(ctx, channelID)
		if err != nil {
			return nil, nil, err
		}
		out := make([]secretOutput, 0, len(items))
		for _, item := range items {
			out = append(out, secretMetadata(item.ID.String(), item.Name, item.Description.Ptr()))
		}
		return nil, map[string]any{"secrets": out}, nil
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_secret", Description: "Get and decrypt one channel secret. Treat the returned value as sensitive."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		meta, err := h.ownedSecret(ctx, requestScope, input.ID)
		if err != nil {
			return nil, nil, err
		}
		id, _ := parseID(input.ID)
		value, err := h.deps.Secrets.GetDecryptedValue(ctx, id)
		meta.Value = &value
		return nil, meta, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_secret", Description: "Create an encrypted channel secret."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createSecretInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Secrets.Create(ctx, channels_secret.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Name: input.Name, Description: input.Description, Value: input.Value})
		return nil, secretMetadata(item.ID.String(), item.Name, item.Description.Ptr()), err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_secret", Description: "Update an encrypted channel secret."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateSecretInput) (*modelsdk.CallToolResult, any, error) {
		if _, err := h.ownedSecret(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		id, _ := parseID(input.ID)
		item, err := h.deps.Secrets.Update(ctx, id, channels_secret.UpdateInput{Name: input.Name, Description: input.Description, Value: input.Value})
		return nil, secretMetadata(item.ID.String(), item.Name, item.Description.Ptr()), err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_secret", Description: "Delete an encrypted channel secret."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		if _, err := h.ownedSecret(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		id, _ := parseID(input.ID)
		err := h.deps.Secrets.Delete(ctx, id)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
}

type storageKeyInput struct {
	Key string `json:"key"`
}
type setStorageInput struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value" jsonschema:"any valid JSON value"`
}

func (h *Handler) addStorageTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_storage_files", Description: "List channel JSON storage entries."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Storage.GetAllByChannelID(ctx, channelID)
		return nil, map[string]any{"files": items}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_storage_file", Description: "Get a channel JSON storage entry by key."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input storageKeyInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Storage.GetByKey(ctx, channelID, input.Key)
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "upload_storage_file", Description: "Create or replace a channel JSON storage entry."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input setStorageInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Storage.Set(ctx, channelID, input.Key, input.Value)
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_storage_file", Description: "Delete a channel JSON storage entry."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input storageKeyInput) (*modelsdk.CallToolResult, any, error) {
		err := h.deps.Storage.Delete(ctx, channelID, input.Key)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_storage_usage", Description: "Get used and maximum channel storage bytes."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		used, err := h.deps.Storage.GetTotalSizeByChannelID(ctx, channelID)
		return nil, map[string]int64{"usedBytes": used, "limitBytes": 30 * 1024 * 1024}, err
	})
}

type createPasteInput struct {
	Content  string     `json:"content"`
	ExpireAt *time.Time `json:"expireAt,omitempty"`
}
type updatePasteInput struct {
	ID       string     `json:"id"`
	Content  string     `json:"content"`
	ExpireAt *time.Time `json:"expireAt,omitempty"`
}
type listInput struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"perPage,omitempty"`
}

func ownerID(requestScope scope) string {
	if len(requestScope.Channel.Bindings) > 0 {
		return requestScope.Channel.Bindings[0].UserID.String()
	}
	return requestScope.Channel.ID.String()
}
func (h *Handler) ownedPaste(ctx context.Context, requestScope scope, id string) (any, error) {
	item, err := h.deps.Pastebins.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item.OwnerUserID == nil || *item.OwnerUserID != ownerID(requestScope) {
		return nil, errors.New("paste not found")
	}
	return item, nil
}

func (h *Handler) addPastebinTools(s *modelsdk.Server, requestScope scope) {
	owner := ownerID(requestScope)
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_pastes", Description: "List pastes owned by this channel owner."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input listInput) (*modelsdk.CallToolResult, any, error) {
		if input.Page < 1 {
			input.Page = 1
		}
		if input.PerPage < 1 {
			input.PerPage = 20
		}
		result, err := h.deps.Pastebins.GetUserMany(ctx, pastebins.GetManyInput{Page: input.Page, PerPage: input.PerPage, OwnerUserID: owner})
		return nil, result, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_paste", Description: "Get an owned paste including its content."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.ownedPaste(ctx, requestScope, input.ID)
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_paste_raw", Description: "Get raw content of an owned paste."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		raw, err := h.deps.Pastebins.GetByID(ctx, input.ID)
		if err != nil || raw.OwnerUserID == nil || *raw.OwnerUserID != owner {
			return nil, nil, errors.New("paste not found")
		}
		return nil, map[string]string{"content": raw.Content}, nil
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_paste", Description: "Create an owned paste."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createPasteInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Pastebins.Create(ctx, pastebins.CreateInput{Content: input.Content, ExpireAt: input.ExpireAt, OwnerUserID: &owner})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_paste", Description: "Update the content and expiration of an owned paste."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updatePasteInput) (*modelsdk.CallToolResult, any, error) {
		if _, err := h.ownedPaste(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		item, err := h.deps.Pastebins.Update(ctx, input.ID, input.Content, input.ExpireAt)
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_paste", Description: "Delete an owned paste."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		if _, err := h.ownedPaste(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		err := h.deps.Pastebins.Delete(ctx, input.ID)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
}

type createShortURLInput struct {
	URL   string `json:"url"`
	Alias string `json:"alias,omitempty"`
}
type updateShortURLInput struct {
	ID    string  `json:"id"`
	URL   *string `json:"url,omitempty"`
	Alias *string `json:"alias,omitempty"`
}
type statsInput struct {
	ID       string    `json:"id"`
	From     time.Time `json:"from"`
	To       time.Time `json:"to"`
	Interval string    `json:"interval,omitempty" jsonschema:"hour or day"`
}
type bannedUAInput struct {
	ID          string  `json:"id"`
	Pattern     string  `json:"pattern"`
	Description *string `json:"description,omitempty"`
}
type deleteBannedUAInput struct {
	ID                string `json:"id"`
	BannedUserAgentID string `json:"bannedUserAgentId"`
}

func (h *Handler) ownedShortURL(ctx context.Context, requestScope scope, id string) error {
	item, err := h.deps.ShortURLs.GetByShortID(ctx, nil, id)
	if err != nil {
		return err
	}
	owner := ownerID(requestScope)
	if item.IsNil() || item.CreatedByUserId == nil || *item.CreatedByUserId != owner {
		return errors.New("short URL not found")
	}
	return nil
}

func (h *Handler) addShortURLTools(s *modelsdk.Server, requestScope scope) {
	owner := ownerID(requestScope)
	listShortURLs := func(ctx context.Context, _ *modelsdk.CallToolRequest, input listInput) (*modelsdk.CallToolResult, any, error) {
		if input.PerPage < 1 {
			input.PerPage = 20
		}
		result, err := h.deps.ShortURLs.GetList(ctx, shortenedurls.GetListInput{Page: input.Page, PerPage: input.PerPage, OwnerUserID: &owner, SortBy: "created_at"})
		return nil, result, err
	}
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "shorten_url", Description: "Create a short URL owned by this channel owner."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createShortURLInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.ShortURLs.Create(ctx, shortenedurls.CreateInput{CreatedByUserID: &owner, ShortID: input.Alias, URL: input.URL})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_short_url", Description: "Get one owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		item, err := h.deps.ShortURLs.GetByShortID(ctx, nil, input.ID)
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_short_urls", Description: "List short URLs owned by this channel owner."}, listShortURLs)
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_short_links", Description: "List short links owned by this channel owner."}, listShortURLs)
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_short_url", Description: "Update an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateShortURLInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		item, err := h.deps.ShortURLs.Update(ctx, nil, input.ID, shortenedurls.UpdateInput{URL: input.URL, ShortID: input.Alias})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_short_url", Description: "Delete an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		err := h.deps.ShortURLs.Delete(ctx, nil, input.ID)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_short_url_stats", Description: "Get click statistics for an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input statsInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		if input.Interval == "" {
			input.Interval = "day"
		}
		points, err := h.deps.ShortURLs.GetStatistics(ctx, shortenedurls.GetStatisticsInput{ShortLinkID: input.ID, From: input.From, To: input.To, Interval: input.Interval})
		return nil, map[string]any{"points": points}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_banned_user_agents", Description: "List user-agent patterns banned for an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		items, err := h.deps.ShortURLs.GetLinkBannedUserAgents(ctx, input.ID)
		return nil, map[string]any{"patterns": items}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "ban_user_agent", Description: "Ban a user-agent pattern for an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input bannedUAInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		item, err := h.deps.ShortURLs.CreateLinkBannedUserAgent(ctx, shortlinksua.CreateInput{LinkID: input.ID, Pattern: input.Pattern, Description: input.Description})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "unban_user_agent", Description: "Remove a user-agent ban from an owned short URL."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input deleteBannedUAInput) (*modelsdk.CallToolResult, any, error) {
		if err := h.ownedShortURL(ctx, requestScope, input.ID); err != nil {
			return nil, nil, err
		}
		err := h.deps.ShortURLs.DeleteLinkBannedUserAgent(ctx, input.BannedUserAgentID, input.ID)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
}

var _ = fmt.Sprintf
