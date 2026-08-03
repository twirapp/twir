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
	return isReadTool(name)
}

func isReadTool(name string) bool {
	return strings.HasPrefix(name, "list_") || strings.HasPrefix(name, "get_")
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
