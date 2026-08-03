package mcp

import (
	"context"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/quotes"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
)

const variableScriptGuide = `SCRIPT variable guide:
- Put JavaScript source in evalValue, set type to SCRIPT and scriptLanguage to javascript. response is normally empty for SCRIPT variables.
- The source runs inside an async function, so top-level await is supported. Always return the final value; a script without return produces an empty result. Returned values are converted to strings, and command output is limited to 1000 characters.
- Supported APIs: fetch(url, options) with response.ok/status/statusText/headers/json()/text(); storage and twir.storage persistent async methods get/set/delete/has/keys/clear/push/pop/find/filter/splice/getProperty/setProperty/deleteProperty/hasProperty/merge; twir.secrets.get(name); twir.channel.id; lodash as _; setTimeout/clearTimeout/setInterval/clearInterval; URL/URLSearchParams; TextEncoder/TextDecoder; btoa/atob; crypto.randomUUID(); AbortController/AbortSignal; and standard JavaScript built-ins.
- The total execution budget is about 5 seconds. HTTP redirects, private/internal network targets, responses over 1 MB, Node/Bun APIs, require(), and arbitrary imports are unavailable.
- Twir built-in variables are expanded as raw text before JavaScript runs. Put them inside a quoted string, for example const sender = "$(sender.displayName)"; for a number use Number("$(user.messages)"). Use list_builtin_variables to discover current names and examples. $(command.param) contains command arguments.
- Handle fetch failures and non-2xx statuses when a friendly chat response is needed. After creating or updating a SCRIPT variable, call evaluate_variable; pass testAsUserName when the script contains sender/user-dependent built-in variables.
Example:
const sender = "$(sender.displayName)";
const response = await fetch("https://jsonplaceholder.typicode.com/posts/1");
if (!response.ok) return "API error: HTTP " + response.status;
const post = await response.json();
return sender + ": " + post.title;`

type createVariableInput struct {
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Type           string  `json:"type" jsonschema:"TEXT, NUMBER, or SCRIPT"`
	EvalValue      string  `json:"evalValue,omitempty" jsonschema:"JavaScript source when type is SCRIPT; top-level await is supported and the script must return its final value"`
	Response       string  `json:"response" jsonschema:"literal value or response template; normally empty for SCRIPT variables"`
	ScriptLanguage string  `json:"scriptLanguage,omitempty" jsonschema:"must be javascript for SCRIPT variables"`
}

type updateVariableInput struct {
	ID             string  `json:"id"`
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	Type           *string `json:"type,omitempty"`
	EvalValue      *string `json:"evalValue,omitempty"`
	Response       *string `json:"response,omitempty"`
	ScriptLanguage *string `json:"scriptLanguage,omitempty"`
}

type setVariableInput struct {
	ID             *string `json:"id,omitempty" jsonschema:"existing variable UUID; omit to create"`
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	Type           string  `json:"type" jsonschema:"TEXT, NUMBER, or SCRIPT"`
	EvalValue      string  `json:"evalValue,omitempty"`
	Response       string  `json:"response" jsonschema:"normally empty for SCRIPT variables"`
	ScriptLanguage string  `json:"scriptLanguage,omitempty" jsonschema:"must be javascript for SCRIPT variables"`
}

type evaluateVariableInput struct {
	ID             string  `json:"id"`
	TestAsUserName *string `json:"testAsUserName,omitempty" jsonschema:"optional Twitch login used as evaluation context"`
}

func (h *Handler) addVariableTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_builtin_variables", Description: "List Twir built-in variables available for templates and SCRIPT source. In JavaScript, wrap expansions in quotes because Twir substitutes their raw text before execution, for example: const sender = \"$(sender.displayName)\"."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Variables.GetBuiltIn(ctx)
		return nil, map[string]any{"variables": items, "scriptGuide": variableScriptGuide}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_variables", Description: "List custom variables for this channel."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Variables.GetAll(ctx, channelID)
		return nil, map[string]any{"variables": items}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_variable", Description: "Get a custom variable by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Variables.GetByID(ctx, input.ID)
		if err == nil && item.ChannelID != channelID {
			return nil, nil, variables.ErrNotFound
		}
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_variable", Description: "Create a channel custom variable. For SCRIPT variables, follow the server SCRIPT variable guide, always return a value, then verify it with evaluate_variable."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createVariableInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Variables.Create(ctx, variables.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Name: input.Name, Description: input.Description, Type: entity.CustomVarType(input.Type), EvalValue: input.EvalValue, Response: input.Response, ScriptLanguage: input.ScriptLanguage})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "set_variable", Description: "Create a custom variable, or replace one when id is provided. For SCRIPT variables, follow the server SCRIPT variable guide and verify with evaluate_variable."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input setVariableInput) (*modelsdk.CallToolResult, any, error) {
		if input.ID == nil {
			item, err := h.deps.Variables.Create(ctx, variables.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Name: input.Name, Description: input.Description, Type: entity.CustomVarType(input.Type), EvalValue: input.EvalValue, Response: input.Response, ScriptLanguage: input.ScriptLanguage})
			return nil, item, err
		}

		id, err := parseID(*input.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid variable id: %w", err)
		}
		variableType := entity.CustomVarType(input.Type)
		language := entity.CustomVarScriptLanguage(input.ScriptLanguage)
		item, err := h.deps.Variables.Update(ctx, variables.UpdateInput{ID: id, ChannelID: channelID, ActorID: requestScope.ActorID, Name: &input.Name, Description: input.Description, Type: &variableType, EvalValue: &input.EvalValue, Response: &input.Response, ScriptLanguage: &language})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_variable", Description: "Update a channel custom variable by UUID. For SCRIPT variables, follow the server SCRIPT variable guide and verify with evaluate_variable."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateVariableInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid variable id: %w", err)
		}
		var variableType *entity.CustomVarType
		if input.Type != nil {
			value := entity.CustomVarType(*input.Type)
			variableType = &value
		}
		var language *entity.CustomVarScriptLanguage
		if input.ScriptLanguage != nil {
			value := entity.CustomVarScriptLanguage(*input.ScriptLanguage)
			language = &value
		}
		item, err := h.deps.Variables.Update(ctx, variables.UpdateInput{ID: id, ChannelID: channelID, ActorID: requestScope.ActorID, Name: input.Name, Description: input.Description, Type: variableType, EvalValue: input.EvalValue, Response: input.Response, ScriptLanguage: language})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_variable", Description: "Delete a channel custom variable by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, fmt.Errorf("invalid variable id: %w", err)
		}
		err = h.deps.Variables.Delete(ctx, id, channelID, requestScope.ActorID)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "evaluate_variable", Description: "Evaluate a variable using the current channel context. For SCRIPT variables containing sender/user-dependent built-ins, pass testAsUserName so those placeholders are expanded before execution."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input evaluateVariableInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Variables.GetByID(ctx, input.ID)
		if err != nil {
			return nil, nil, err
		}
		if item.ChannelID != channelID {
			return nil, nil, variables.ErrNotFound
		}
		if item.Type != entity.CustomVarScript {
			return nil, map[string]string{"value": item.Response}, nil
		}
		value, err := h.deps.Variables.EvaluateScript(ctx, channelID, item.EvalValue, item.ScriptLanguage, input.TestAsUserName)
		return nil, map[string]string{"value": value}, err
	})
}

type createQuoteInput struct {
	Text        string  `json:"text"`
	CreatorName *string `json:"creatorName,omitempty"`
}
type updateQuoteInput struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func (h *Handler) addQuoteTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_quotes", Description: "List quotes for this channel."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Quotes.GetAllByChannelID(ctx, channelID)
		return nil, map[string]any{"quotes": items}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_quote", Description: "Create a quote for this channel."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createQuoteInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Quotes.Create(ctx, quotes.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Text: input.Text, CreatorName: input.CreatorName})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_quote", Description: "Update a quote by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateQuoteInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, err
		}
		item, err := h.deps.Quotes.Update(ctx, quotes.UpdateInput{ID: id, ChannelID: channelID, ActorID: requestScope.ActorID, Text: &input.Text})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_quote", Description: "Delete a quote by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, err
		}
		err = h.deps.Quotes.Delete(ctx, channelID, requestScope.ActorID, id)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
}

type createKeywordInput struct {
	Text      string `json:"text"`
	Response  string `json:"response"`
	Enabled   bool   `json:"enabled"`
	Cooldown  int    `json:"cooldown,omitempty"`
	IsReply   bool   `json:"isReply,omitempty"`
	IsRegular bool   `json:"isRegular,omitempty"`
}
type updateKeywordInput struct {
	ID        string  `json:"id"`
	Text      *string `json:"text,omitempty"`
	Response  *string `json:"response,omitempty"`
	Enabled   *bool   `json:"enabled,omitempty"`
	Cooldown  *int    `json:"cooldown,omitempty"`
	IsReply   *bool   `json:"isReply,omitempty"`
	IsRegular *bool   `json:"isRegular,omitempty"`
}

func (h *Handler) addKeywordTools(s *modelsdk.Server, requestScope scope) {
	channelID := requestScope.Channel.ID.String()
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_keywords", Description: "List keyword triggers for this channel."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
		items, err := h.deps.Keywords.GetAllByChannelID(ctx, channelID)
		return nil, map[string]any{"keywords": items}, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_keyword", Description: "Create a keyword trigger."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input createKeywordInput) (*modelsdk.CallToolResult, any, error) {
		item, err := h.deps.Keywords.Create(ctx, keywords.CreateInput{ChannelID: channelID, ActorID: requestScope.ActorID, Text: input.Text, Response: input.Response, Enabled: input.Enabled, Cooldown: input.Cooldown, IsReply: input.IsReply, IsRegular: input.IsRegular})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_keyword", Description: "Update a keyword trigger by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateKeywordInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, err
		}
		item, err := h.deps.Keywords.Update(ctx, keywords.UpdateInput{ID: id, ChannelID: channelID, ActorID: requestScope.ActorID, Text: input.Text, Response: input.Response, Enabled: input.Enabled, Cooldown: input.Cooldown, IsReply: input.IsReply, IsRegular: input.IsRegular})
		return nil, item, err
	})
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_keyword", Description: "Delete a keyword trigger by UUID."}, func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
		id, err := parseID(input.ID)
		if err != nil {
			return nil, nil, err
		}
		err = h.deps.Keywords.Delete(ctx, channelID, requestScope.ActorID, id)
		return nil, map[string]bool{"deleted": err == nil}, err
	})
}
