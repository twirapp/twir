package mcp

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/server"
	"github.com/twirapp/twir/apps/api-gql/internal/services/alerts"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_files"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_moderation_settings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_overlays"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_secret"
	"github.com/twirapp/twir/apps/api-gql/internal/services/channels_storage"
	"github.com/twirapp/twir/apps/api-gql/internal/services/chat_wall"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	commandsrelations "github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	"github.com/twirapp/twir/apps/api-gql/internal/services/discord_integration"
	donatellointegration "github.com/twirapp/twir/apps/api-gql/internal/services/donatello_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/donatepay_integration"
	donatestreamintegration "github.com/twirapp/twir/apps/api-gql/internal/services/donatestream_integration"
	donationalertsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/donationalerts_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/events"
	faceitintegration "github.com/twirapp/twir/apps/api-gql/internal/services/faceit_integration"
	gamesvoteban "github.com/twirapp/twir/apps/api-gql/internal/services/games_voteban"
	"github.com/twirapp/twir/apps/api-gql/internal/services/giveaways"
	"github.com/twirapp/twir/apps/api-gql/internal/services/greetings"
	"github.com/twirapp/twir/apps/api-gql/internal/services/keywords"
	lastfmintegration "github.com/twirapp/twir/apps/api-gql/internal/services/lastfm_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/kappagen"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays_dudes"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	"github.com/twirapp/twir/apps/api-gql/internal/services/quotes"
	"github.com/twirapp/twir/apps/api-gql/internal/services/seventv_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/api-gql/internal/services/song_requests"
	"github.com/twirapp/twir/apps/api-gql/internal/services/spotify_integration"
	streamlabsintegration "github.com/twirapp/twir/apps/api-gql/internal/services/streamlabs_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	valorantintegration "github.com/twirapp/twir/apps/api-gql/internal/services/valorant_integration"
	"github.com/twirapp/twir/apps/api-gql/internal/services/variables"
	vkintegration "github.com/twirapp/twir/apps/api-gql/internal/services/vk_integration"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
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
	CommandGroups     *commands_groups.Service
	CommandsRelations *commandsrelations.Service
	Timers            *timers.Service
	Variables         *variables.Service
	Quotes            *quotes.Service
	Keywords          *keywords.Service
	Secrets           *channels_secret.Service
	Storage           *channels_storage.Service
	Files             *channels_files.Service
	Moderation        *channels_moderation_settings.Service
	ChatWall          *chat_wall.Service
	Games             *gamesvoteban.Service
	SongRequests      *song_requests.Service
	Events            *events.Service
	Giveaways         *giveaways.Service
	Greetings         *greetings.Service
	Alerts            *alerts.Service
	Twitch            *twitchservice.Service
	Gorm              *gorm.DB
	Discord           *discord_integration.Service
	Spotify           *spotify_integration.Service
	LastFM            *lastfmintegration.Service
	Valorant          *valorantintegration.Service
	Faceit            *faceitintegration.Service
	DonationAlerts    *donationalertsintegration.Service
	DonatePay         *donatepay_integration.Service
	DonateStream      *donatestreamintegration.Service
	Donatello         *donatellointegration.Service
	Streamlabs        *streamlabsintegration.Service
	VK                *vkintegration.Service
	SevenTV           *seventv_integration.Service
	CustomOverlays    *channels_overlays.Service
	TTS               *tts.Service
	Dudes             *overlays_dudes.Service
	Kappagen          *kappagen.Service
	BeRightBack       *be_right_back.Service
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
	h.addSystemTools(s, requestScope)
	h.addEngagementTools(s, requestScope)
	h.addIntegrationTools(s, requestScope)
	h.addOverlayTools(s, requestScope)

	return s
}

func parseID(value string) (uuid.UUID, error) {
	return uuid.Parse(value)
}
