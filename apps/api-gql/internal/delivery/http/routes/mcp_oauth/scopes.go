package mcp_oauth

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	entity "github.com/twirapp/twir/libs/entities/mcp_oauth"
)

type scopeCatalog struct{}

type scopeCatalogResponse struct {
	Scopes []scopeCatalogItem `json:"scopes"`
}

type scopeCatalogItem struct {
	Group       entity.ScopeGroup    `json:"group"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Actions     []entity.ScopeAction `json:"actions"`
}

var _ httpbase.Route[*emptyInput, *metadataOutput[scopeCatalogResponse]] = (*scopeCatalog)(nil)

func newScopeCatalog() *scopeCatalog {
	return &scopeCatalog{}
}

func (*scopeCatalog) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "mcp-oauth-scopes",
		Method:        http.MethodGet,
		Path:          "/oauth/scopes",
		Tags:          []string{"MCP OAuth"},
		Summary:       "OAuth scope catalog",
		DefaultStatus: http.StatusOK,
	}
}

func (*scopeCatalog) Handler(context.Context, *emptyInput) (*metadataOutput[scopeCatalogResponse], error) {
	groups := entity.AllScopeGroups()
	scopes := make([]scopeCatalogItem, 0, len(groups))
	for _, group := range groups {
		scopes = append(scopes, scopeCatalogItem{
			Group:       group.Group,
			Name:        group.Name,
			Description: group.Description,
			Actions:     []entity.ScopeAction{entity.ScopeActionRead, entity.ScopeActionEdit},
		})
	}

	return &metadataOutput[scopeCatalogResponse]{
		AccessControlAllowOrigin: "*",
		CacheControl:             "public, max-age=3600",
		Body: scopeCatalogResponse{
			Scopes: scopes,
		},
	}, nil
}

func (route *scopeCatalog) Register(api huma.API) {
	meta := route.GetMeta()
	meta.Responses = map[string]*huma.Response{
		"200": {
			Description: "OAuth scope catalog",
			Content:     jsonContent[scopeCatalogResponse](api),
		},
	}
	huma.Register(api, meta, route.Handler)
}
