package mcp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_secret"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_storage"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	commandsrelations "github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/services/quotes"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
)

type contextKey struct{}

type scope struct {
	Channel channelentity.Channel
	ActorID string
}

type Deps struct {
	fx.In

	Channels          *channelservice.ChannelService
	Commands          *commands.Service
	CommandsRelations *commandsrelations.Service
	Timers            *timers.Service
	Variables         *variables.Service
	Quotes            *quotes.Service
	Keywords          *keywords.Service
	Secrets           *channels_secret.Service
	Storage           *channels_storage.Service
	Pastebins         *pastebins.Service
	ShortURLs         *shortenedurls.Service
}

type Handler struct {
	deps      Deps
	transport http.Handler
}

func New(deps Deps) *Handler {
	h := &Handler{deps: deps}
	h.transport = modelsdk.NewStreamableHTTPHandler(
		func(r *http.Request) *modelsdk.Server {
			requestScope, ok := r.Context().Value(contextKey{}).(scope)
			if !ok {
				return nil
			}

			return h.newServer(requestScope)
		},
		&modelsdk.StreamableHTTPOptions{Stateless: true},
	)

	return h
}

func Register(s *server.Server, h *Handler) {
	s.Any("/mcp", gin.WrapH(h))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiKey := r.Header.Get("Api-Key")
	if apiKey == "" {
		http.Error(w, "Api-Key header is required", http.StatusUnauthorized)
		return
	}

	channel, err := h.deps.Channels.GetChannelByApiKey(r.Context(), apiKey)
	if err != nil || channel.IsNil() {
		http.Error(w, "invalid channel API key", http.StatusUnauthorized)
		return
	}

	actorID := channel.ID.String()
	if len(channel.Bindings) > 0 {
		actorID = channel.Bindings[0].UserID.String()
	}

	ctx := context.WithValue(r.Context(), contextKey{}, scope{Channel: channel, ActorID: actorID})
	h.transport.ServeHTTP(w, r.WithContext(ctx))
}

func (h *Handler) newServer(requestScope scope) *modelsdk.Server {
	s := modelsdk.NewServer(
		&modelsdk.Implementation{Name: "twir", Version: "1.0.0"},
		&modelsdk.ServerOptions{Instructions: "Manage the Twir channel authorized by the Api-Key request header. All operations are restricted to that channel."},
	)

	h.addCommandTools(s, requestScope)
	h.addTimerTools(s, requestScope)
	h.addVariableTools(s, requestScope)
	h.addQuoteTools(s, requestScope)
	h.addKeywordTools(s, requestScope)
	h.addSecretTools(s, requestScope)
	h.addStorageTools(s, requestScope)
	h.addPastebinTools(s, requestScope)
	h.addShortURLTools(s, requestScope)

	return s
}

func parseID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
