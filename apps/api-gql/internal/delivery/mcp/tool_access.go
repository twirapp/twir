package mcp

import (
	"strings"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	oauthentity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type toolAccessScope string

const (
	toolAccessScopeRead  toolAccessScope = "read"
	toolAccessScopeWrite toolAccessScope = "write"
)

type toolAccessScopes map[toolAccessScope]struct{}

type toolAccessClassification string

const (
	toolAccessClassificationSensitive toolAccessClassification = "sensitive"
	toolAccessClassificationRead      toolAccessClassification = "read"
	toolAccessClassificationOther     toolAccessClassification = "other"
)

var sensitiveToolNames = map[string]struct{}{
	"get_secret": {},
}

func fullToolAccessScopes() toolAccessScopes {
	return toolAccessScopes{
		toolAccessScopeRead:  {},
		toolAccessScopeWrite: {},
	}
}

func toolAccessScopesFromOAuthScopes(scopes []oauthentity.Scope) (toolAccessScopes, bool) {
	accessScopes := make(toolAccessScopes, len(scopes))
	for _, scope := range scopes {
		switch scope {
		case oauthentity.ScopeRead:
			accessScopes[toolAccessScopeRead] = struct{}{}
		case oauthentity.ScopeWrite:
			accessScopes[toolAccessScopeWrite] = struct{}{}
		default:
			return nil, false
		}
	}
	_, hasRead := accessScopes[toolAccessScopeRead]
	return accessScopes, hasRead
}

func (scopes toolAccessScopes) allowsTool(name string) bool {
	if _, ok := scopes[toolAccessScopeWrite]; ok {
		return true
	}
	if _, ok := scopes[toolAccessScopeRead]; !ok {
		return false
	}
	return classifyToolAccess(name) == toolAccessClassificationRead
}

func classifyToolAccess(name string) toolAccessClassification {
	if _, ok := sensitiveToolNames[name]; ok {
		return toolAccessClassificationSensitive
	}
	if strings.HasPrefix(name, "list_") || strings.HasPrefix(name, "get_") {
		return toolAccessClassificationRead
	}
	return toolAccessClassificationOther
}

type toolRegistrar struct {
	server *modelsdk.Server
	scopes toolAccessScopes
}

func newToolRegistrar(server *modelsdk.Server, scopes toolAccessScopes) toolRegistrar {
	return toolRegistrar{server: server, scopes: scopes}
}

func addTool[In, Out any](registrar toolRegistrar, tool *modelsdk.Tool, handler modelsdk.ToolHandlerFor[In, Out]) {
	if registrar.scopes.allowsTool(tool.Name) {
		modelsdk.AddTool(registrar.server, tool, handler)
	}
}
