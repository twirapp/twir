package mcp_oauth

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var ErrInvalidScope = errors.New("mcp oauth: invalid scope")

type ScopeGroup string

const (
	ScopeGroupCommands      ScopeGroup = "commands"
	ScopeGroupTimers        ScopeGroup = "timers"
	ScopeGroupFiles         ScopeGroup = "files"
	ScopeGroupGames         ScopeGroup = "games"
	ScopeGroupSongRequests  ScopeGroup = "song_requests"
	ScopeGroupModeration    ScopeGroup = "moderation"
	ScopeGroupOverlays      ScopeGroup = "overlays"
	ScopeGroupIntegrations  ScopeGroup = "integrations"
	ScopeGroupEvents        ScopeGroup = "events"
	ScopeGroupRewards       ScopeGroup = "rewards"
	ScopeGroupGiveaways     ScopeGroup = "giveaways"
	ScopeGroupGreetings     ScopeGroup = "greetings"
	ScopeGroupNotifications ScopeGroup = "notifications"
	ScopeGroupAlerts        ScopeGroup = "alerts"
	ScopeGroupSecrets       ScopeGroup = "secrets"
	ScopeGroupStorage       ScopeGroup = "storage"
	ScopeGroupPastes        ScopeGroup = "pastes"
	ScopeGroupShortURLs     ScopeGroup = "short_urls"
	ScopeGroupDashboard     ScopeGroup = "dashboard"
	ScopeGroupVariables     ScopeGroup = "variables"
	ScopeGroupQuotes        ScopeGroup = "quotes"
	ScopeGroupKeywords      ScopeGroup = "keywords"
)

type ScopeAction string

const (
	ScopeActionRead ScopeAction = "read"
	ScopeActionEdit ScopeAction = "edit"
)

type ScopeGroupEntry struct {
	Group       ScopeGroup
	Name        string
	Description string
}

var scopeGroupCatalog = [...]ScopeGroupEntry{
	{Group: ScopeGroupCommands, Name: "Commands", Description: "View and manage custom commands, groups, and role cooldowns"},
	{Group: ScopeGroupTimers, Name: "Timers", Description: "View and manage chat timers"},
	{Group: ScopeGroupFiles, Name: "Files", Description: "View and manage uploaded image and audio files"},
	{Group: ScopeGroupGames, Name: "Games", Description: "View and manage channel games and their settings"},
	{Group: ScopeGroupSongRequests, Name: "Song Requests", Description: "View and manage the song request queue and playback"},
	{Group: ScopeGroupModeration, Name: "Moderation", Description: "View and manage channel moderation rules and chat wall settings"},
	{Group: ScopeGroupOverlays, Name: "Overlays", Description: "View and manage custom and built-in overlay settings"},
	{Group: ScopeGroupIntegrations, Name: "Integrations", Description: "View and manage connected third-party integrations"},
	{Group: ScopeGroupEvents, Name: "Events", Description: "View and manage channel automation events and operations"},
	{Group: ScopeGroupRewards, Name: "Rewards", Description: "View and manage Twitch custom rewards"},
	{Group: ScopeGroupGiveaways, Name: "Giveaways", Description: "View and manage channel giveaways"},
	{Group: ScopeGroupGreetings, Name: "Greetings", Description: "View and manage channel greetings"},
	{Group: ScopeGroupNotifications, Name: "Notifications", Description: "View and manage channel notifications"},
	{Group: ScopeGroupAlerts, Name: "Alerts", Description: "View and manage channel alerts and their bindings"},
	{Group: ScopeGroupSecrets, Name: "Secrets", Description: "View and manage encrypted channel secrets"},
	{Group: ScopeGroupStorage, Name: "Storage", Description: "View and manage channel JSON storage entries"},
	{Group: ScopeGroupPastes, Name: "Pastes", Description: "View and manage owned pastes"},
	{Group: ScopeGroupShortURLs, Name: "Short URLs", Description: "View and manage owned short URLs and their settings"},
	{Group: ScopeGroupDashboard, Name: "Dashboard", Description: "View and manage channel dashboard settings and statistics"},
	{Group: ScopeGroupVariables, Name: "Variables", Description: "View and manage channel custom variables"},
	{Group: ScopeGroupQuotes, Name: "Quotes", Description: "View and manage channel quotes"},
	{Group: ScopeGroupKeywords, Name: "Keywords", Description: "View and manage keyword triggers"},
}

func AllScopeGroups() []ScopeGroupEntry {
	groups := make([]ScopeGroupEntry, len(scopeGroupCatalog))
	copy(groups, scopeGroupCatalog[:])
	return groups
}

func AllScopes() []Scope {
	scopes := make([]Scope, 0, len(scopeGroupCatalog)*2)
	for _, entry := range scopeGroupCatalog {
		scopes = append(scopes, scopeFor(entry.Group, ScopeActionRead), scopeFor(entry.Group, ScopeActionEdit))
	}
	return scopes
}

func ParseScopes(raw string) ([]Scope, error) {
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return nil, invalidScope(raw)
	}

	scopes := make([]Scope, len(tokens))
	for index, token := range tokens {
		scopes[index] = Scope(token)
	}
	return NormalizeScopes(scopes)
}

func NormalizeScopes(scopes []Scope) ([]Scope, error) {
	if len(scopes) == 0 {
		return nil, invalidScope("")
	}

	included := make(map[Scope]struct{}, len(scopes)*2)
	for _, scope := range scopes {
		switch scope {
		case ScopeRead:
			for _, entry := range scopeGroupCatalog {
				includeScope(included, entry.Group, ScopeActionRead)
			}
		case ScopeWrite:
			for _, entry := range scopeGroupCatalog {
				includeScope(included, entry.Group, ScopeActionEdit)
			}
		default:
			group, action, ok := parseCanonicalScope(scope)
			if !ok {
				return nil, invalidScope(string(scope))
			}
			includeScope(included, group, action)
		}
	}

	return canonicalScopes(included), nil
}

func ScopeSubset(requested, allowed []Scope) bool {
	normalizedRequested, requestedErr := NormalizeScopes(requested)
	if requestedErr != nil {
		return false
	}
	normalizedAllowed, allowedErr := NormalizeScopes(allowed)
	if allowedErr != nil {
		return false
	}

	allowedSet := make(map[Scope]struct{}, len(normalizedAllowed))
	for _, scope := range normalizedAllowed {
		allowedSet[scope] = struct{}{}
	}
	for _, scope := range normalizedRequested {
		if _, ok := allowedSet[scope]; !ok {
			return false
		}
	}
	return true
}

func HasScope(scopes []Scope, group ScopeGroup, action ScopeAction) bool {
	normalized, err := NormalizeScopes(scopes)
	if err != nil {
		return false
	}
	wanted := scopeFor(group, action)
	return slices.Contains(normalized, wanted)
}

func ScopeStrings(scopes []Scope) []string {
	strings := make([]string, len(scopes))
	for index, scope := range scopes {
		strings[index] = string(scope)
	}
	return strings
}

func scopeFor(group ScopeGroup, action ScopeAction) Scope {
	return Scope(string(group) + ":" + string(action))
}

func includeScope(included map[Scope]struct{}, group ScopeGroup, action ScopeAction) {
	included[scopeFor(group, ScopeActionRead)] = struct{}{}
	if action == ScopeActionEdit {
		included[scopeFor(group, ScopeActionEdit)] = struct{}{}
	}
}

func canonicalScopes(included map[Scope]struct{}) []Scope {
	scopes := make([]Scope, 0, len(included))
	for _, entry := range scopeGroupCatalog {
		for _, action := range [...]ScopeAction{ScopeActionRead, ScopeActionEdit} {
			scope := scopeFor(entry.Group, action)
			if _, ok := included[scope]; ok {
				scopes = append(scopes, scope)
			}
		}
	}
	return scopes
}

func parseCanonicalScope(scope Scope) (ScopeGroup, ScopeAction, bool) {
	parts := strings.Split(string(scope), ":")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	group := ScopeGroup(parts[0])
	action := ScopeAction(parts[1])
	if !isScopeGroup(group) || !isScopeAction(action) {
		return "", "", false
	}
	return group, action, true
}

func isScopeGroup(group ScopeGroup) bool {
	for _, entry := range scopeGroupCatalog {
		if entry.Group == group {
			return true
		}
	}
	return false
}

func isScopeAction(action ScopeAction) bool {
	return action == ScopeActionRead || action == ScopeActionEdit
}

func invalidScope(scope string) error {
	return fmt.Errorf("scope %q: %w", scope, ErrInvalidScope)
}
