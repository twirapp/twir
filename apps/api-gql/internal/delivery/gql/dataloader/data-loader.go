package dataloader

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelsemotesusages "github.com/twirapp/twir/apps/api-gql/internal/services/channels_emotes_usages"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	twitchservice "github.com/twirapp/twir/apps/api-gql/internal/services/twitch"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/vikstrous/dataloadgen"
	"go.uber.org/fx"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type Opts struct {
	fx.In

	AuthService              *auth.Auth
	CachedTwitchClient       *twitch.CachedTwitchClient
	CommandsGroupsService    *commands_groups.Service
	CommandsResponsesService *commands_responses.Service
	TwitchService            *twitchservice.Service
	EmoteStatisticService    *channelsemotesusages.Service
}

type dataLoader struct {
	deps Opts

	helixUserByIdLoader         *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
	helixUserByNameLoader       *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
	twitchCategoriesByIdLoader  *dataloadgen.Loader[string, *gqlmodel.TwitchCategory]
	commandsGroupsByIdLoader    *dataloadgen.Loader[uuid.UUID, *gqlmodel.CommandGroup]
	commandsResponsesByIDLoader *dataloadgen.Loader[uuid.UUID, []gqlmodel.CommandResponse]
	emoteStatistic              *dataloadgen.Loader[EmoteRangeKey, []gqlmodel.EmoteStatisticUsage]
}

type LoaderFactory struct {
	deps Opts
}

func New(opts Opts) *LoaderFactory {
	return &LoaderFactory{
		deps: opts,
	}
}

func (c *LoaderFactory) Load() *dataLoader {
	loader := &dataLoader{
		deps: c.deps,
	}

	loader.helixUserByIdLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByIds,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.helixUserByNameLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByNames,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.twitchCategoriesByIdLoader = dataloadgen.NewLoader(
		loader.getTwitchCategoriesByIDs,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.commandsGroupsByIdLoader = dataloadgen.NewLoader(
		loader.getCommandsGroupsByIDs,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.commandsResponsesByIDLoader = dataloadgen.NewLoader(
		loader.getCommandsResponsesByIDs,
		dataloadgen.WithWait(time.Millisecond),
	)

	loader.emoteStatistic = dataloadgen.NewLoader(
		loader.getEmoteStatistic,
		dataloadgen.WithWait(time.Millisecond),
	)

	return loader
}

func (c *LoaderFactory) LoadMiddleware(g *gin.Context) {
	loaderForRequest := c.Load()

	g.Request = g.Request.WithContext(
		context.WithValue(g.Request.Context(), LoadersKey, loaderForRequest),
	)

	g.Next()
}

// GetLoaderForRequest returns the dataloader for a given context
func GetLoaderForRequest(ctx context.Context) *dataLoader {
	return ctx.Value(LoadersKey).(*dataLoader)
}
