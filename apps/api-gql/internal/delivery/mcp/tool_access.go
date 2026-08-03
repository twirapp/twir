package mcp

import (
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	oauthentity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type toolAccessScope = oauthentity.Scope

const (
	toolAccessScopeRead  toolAccessScope = "commands:read"
	toolAccessScopeWrite toolAccessScope = "integrations:edit"
)

type toolAccessScopes map[oauthentity.Scope]struct{}

type toolPermission struct {
	Group  oauthentity.ScopeGroup
	Action oauthentity.ScopeAction
}

var toolPermissions = map[string]toolPermission{
	"list_command_roles":  {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionRead},
	"list_command_groups": {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionRead},
	"list_commands":       {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionRead},
	"get_command":         {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionRead},
	"create_command":      {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionEdit},
	"update_command":      {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionEdit},
	"delete_command":      {Group: oauthentity.ScopeGroupCommands, Action: oauthentity.ScopeActionEdit},

	"list_timers":  {Group: oauthentity.ScopeGroupTimers, Action: oauthentity.ScopeActionRead},
	"create_timer": {Group: oauthentity.ScopeGroupTimers, Action: oauthentity.ScopeActionEdit},
	"update_timer": {Group: oauthentity.ScopeGroupTimers, Action: oauthentity.ScopeActionEdit},
	"delete_timer": {Group: oauthentity.ScopeGroupTimers, Action: oauthentity.ScopeActionEdit},

	"list_files":  {Group: oauthentity.ScopeGroupFiles, Action: oauthentity.ScopeActionRead},
	"upload_file": {Group: oauthentity.ScopeGroupFiles, Action: oauthentity.ScopeActionEdit},
	"list_games":  {Group: oauthentity.ScopeGroupGames, Action: oauthentity.ScopeActionRead},

	"list_song_requests": {Group: oauthentity.ScopeGroupSongRequests, Action: oauthentity.ScopeActionRead},
	"get_current_song":   {Group: oauthentity.ScopeGroupSongRequests, Action: oauthentity.ScopeActionRead},
	"skip_song":          {Group: oauthentity.ScopeGroupSongRequests, Action: oauthentity.ScopeActionEdit},
	"manage_queue":       {Group: oauthentity.ScopeGroupSongRequests, Action: oauthentity.ScopeActionEdit},

	"get_moderation_settings": {Group: oauthentity.ScopeGroupModeration, Action: oauthentity.ScopeActionRead},
	"list_mod_chat_wall":      {Group: oauthentity.ScopeGroupModeration, Action: oauthentity.ScopeActionRead},
	"update_moderation":       {Group: oauthentity.ScopeGroupModeration, Action: oauthentity.ScopeActionEdit},

	"list_overlays": {Group: oauthentity.ScopeGroupOverlays, Action: oauthentity.ScopeActionEdit},
	"get_overlay":   {Group: oauthentity.ScopeGroupOverlays, Action: oauthentity.ScopeActionEdit},

	"get_integration_status": {Group: oauthentity.ScopeGroupIntegrations, Action: oauthentity.ScopeActionRead},
	"toggle_integration":     {Group: oauthentity.ScopeGroupIntegrations, Action: oauthentity.ScopeActionEdit},

	"list_events":      {Group: oauthentity.ScopeGroupEvents, Action: oauthentity.ScopeActionRead},
	"list_twir_events": {Group: oauthentity.ScopeGroupEvents, Action: oauthentity.ScopeActionRead},
	"list_rewards":     {Group: oauthentity.ScopeGroupRewards, Action: oauthentity.ScopeActionRead},
	"manage_rewards":   {Group: oauthentity.ScopeGroupRewards, Action: oauthentity.ScopeActionEdit},

	"list_giveaways":  {Group: oauthentity.ScopeGroupGiveaways, Action: oauthentity.ScopeActionRead},
	"create_giveaway": {Group: oauthentity.ScopeGroupGiveaways, Action: oauthentity.ScopeActionEdit},
	"list_greetings":  {Group: oauthentity.ScopeGroupGreetings, Action: oauthentity.ScopeActionRead},
	"update_greeting": {Group: oauthentity.ScopeGroupGreetings, Action: oauthentity.ScopeActionEdit},

	"list_notifications": {Group: oauthentity.ScopeGroupNotifications, Action: oauthentity.ScopeActionRead},
	"get_notification":   {Group: oauthentity.ScopeGroupNotifications, Action: oauthentity.ScopeActionRead},
	"list_alerts":        {Group: oauthentity.ScopeGroupAlerts, Action: oauthentity.ScopeActionRead},
	"manage_alerts":      {Group: oauthentity.ScopeGroupAlerts, Action: oauthentity.ScopeActionEdit},

	"list_secrets":  {Group: oauthentity.ScopeGroupSecrets, Action: oauthentity.ScopeActionRead},
	"get_secret":    {Group: oauthentity.ScopeGroupSecrets, Action: oauthentity.ScopeActionEdit},
	"create_secret": {Group: oauthentity.ScopeGroupSecrets, Action: oauthentity.ScopeActionEdit},
	"update_secret": {Group: oauthentity.ScopeGroupSecrets, Action: oauthentity.ScopeActionEdit},
	"delete_secret": {Group: oauthentity.ScopeGroupSecrets, Action: oauthentity.ScopeActionEdit},

	"list_storage_files":  {Group: oauthentity.ScopeGroupStorage, Action: oauthentity.ScopeActionRead},
	"get_storage_file":    {Group: oauthentity.ScopeGroupStorage, Action: oauthentity.ScopeActionRead},
	"get_storage_usage":   {Group: oauthentity.ScopeGroupStorage, Action: oauthentity.ScopeActionRead},
	"upload_storage_file": {Group: oauthentity.ScopeGroupStorage, Action: oauthentity.ScopeActionEdit},
	"delete_storage_file": {Group: oauthentity.ScopeGroupStorage, Action: oauthentity.ScopeActionEdit},

	"list_pastes":   {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionRead},
	"get_paste":     {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionRead},
	"get_paste_raw": {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionRead},
	"create_paste":  {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionEdit},
	"update_paste":  {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionEdit},
	"delete_paste":  {Group: oauthentity.ScopeGroupPastes, Action: oauthentity.ScopeActionEdit},

	"get_short_url":           {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionRead},
	"list_short_urls":         {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionRead},
	"list_short_links":        {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionRead},
	"get_short_url_stats":     {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionRead},
	"list_banned_user_agents": {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionRead},
	"shorten_url":             {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionEdit},
	"update_short_url":        {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionEdit},
	"delete_short_url":        {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionEdit},
	"ban_user_agent":          {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionEdit},
	"unban_user_agent":        {Group: oauthentity.ScopeGroupShortURLs, Action: oauthentity.ScopeActionEdit},

	"get_stats":            {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionRead},
	"get_bot_settings":     {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionRead},
	"list_community_users": {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionRead},
	"get_user_info":        {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionRead},
	"list_scheduled_vips":  {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionRead},
	"update_bot_settings":  {Group: oauthentity.ScopeGroupDashboard, Action: oauthentity.ScopeActionEdit},

	"list_builtin_variables": {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionRead},
	"list_variables":         {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionRead},
	"get_variable":           {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionRead},
	"create_variable":        {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionEdit},
	"set_variable":           {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionEdit},
	"update_variable":        {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionEdit},
	"delete_variable":        {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionEdit},
	"evaluate_variable":      {Group: oauthentity.ScopeGroupVariables, Action: oauthentity.ScopeActionEdit},

	"list_quotes":    {Group: oauthentity.ScopeGroupQuotes, Action: oauthentity.ScopeActionRead},
	"create_quote":   {Group: oauthentity.ScopeGroupQuotes, Action: oauthentity.ScopeActionEdit},
	"update_quote":   {Group: oauthentity.ScopeGroupQuotes, Action: oauthentity.ScopeActionEdit},
	"delete_quote":   {Group: oauthentity.ScopeGroupQuotes, Action: oauthentity.ScopeActionEdit},
	"list_keywords":  {Group: oauthentity.ScopeGroupKeywords, Action: oauthentity.ScopeActionRead},
	"create_keyword": {Group: oauthentity.ScopeGroupKeywords, Action: oauthentity.ScopeActionEdit},
	"update_keyword": {Group: oauthentity.ScopeGroupKeywords, Action: oauthentity.ScopeActionEdit},
	"delete_keyword": {Group: oauthentity.ScopeGroupKeywords, Action: oauthentity.ScopeActionEdit},
}

func fullToolAccessScopes() toolAccessScopes {
	return toolAccessScopesFromNormalized(oauthentity.AllScopes())
}

func toolAccessScopesFromOAuthScopes(scopes []oauthentity.Scope) (toolAccessScopes, bool) {
	normalized, err := oauthentity.NormalizeScopes(scopes)
	if err != nil {
		return nil, false
	}
	return toolAccessScopesFromNormalized(normalized), true
}

func toolAccessScopesFromNormalized(scopes []oauthentity.Scope) toolAccessScopes {
	accessScopes := make(toolAccessScopes, len(scopes))
	for _, scope := range scopes {
		accessScopes[scope] = struct{}{}
	}
	return accessScopes
}

func (scopes toolAccessScopes) allowsTool(name string) bool {
	permission, ok := toolPermissions[name]
	if !ok {
		return false
	}

	granted := make([]oauthentity.Scope, 0, len(scopes))
	for scope := range scopes {
		granted = append(granted, scope)
	}
	return oauthentity.HasScope(granted, permission.Group, permission.Action)
}

type toolRegistrar struct {
	server *modelsdk.Server
	scopes toolAccessScopes
}

func newToolRegistrar(server *modelsdk.Server, scopes toolAccessScopes) toolRegistrar {
	return toolRegistrar{server: server, scopes: scopes}
}

func addTool[In, Out any](registrar toolRegistrar, tool *modelsdk.Tool, handler modelsdk.ToolHandlerFor[In, Out]) {
	if _, ok := toolPermissions[tool.Name]; !ok {
		panic(fmt.Sprintf("MCP tool %q is missing a static permission mapping", tool.Name))
	}
	if registrar.scopes.allowsTool(tool.Name) {
		modelsdk.AddTool(registrar.server, tool, handler)
	}
}
