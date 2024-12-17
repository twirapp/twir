package data_loader

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_groups"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands_responses"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type DataLoader struct {
	cachedTwitchClient       *twitch.CachedTwitchClient
	commandsGroupsService    *commands_groups.Service
	commandsResponsesService *commands_responses.Service

	helixUserByIdLoader         *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
	helixUserByNameLoader       *dataloadgen.Loader[string, *gqlmodel.TwirUserTwitchInfo]
	twitchCategoriesByIdLoader  *dataloadgen.Loader[string, *gqlmodel.TwitchCategory]
	commandsGroupsByIdLoader    *dataloadgen.Loader[uuid.UUID, *gqlmodel.CommandGroup]
	commandsResponsesByIDLoader *dataloadgen.Loader[uuid.UUID, []gqlmodel.CommandResponse]
}

type Opts struct {
	CachedTwitchClient       *twitch.CachedTwitchClient
	CommandsGroupsService    *commands_groups.Service
	CommandsResponsesService *commands_responses.Service
}

func New(opts Opts) *DataLoader {
	loader := &DataLoader{
		cachedTwitchClient:       opts.CachedTwitchClient,
		commandsGroupsService:    opts.CommandsGroupsService,
		commandsResponsesService: opts.CommandsResponsesService,
	}

	loader.helixUserByIdLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByIds,
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

	return loader
}

// GetLoaderForRequest returns the dataloader for a given context
func GetLoaderForRequest(ctx context.Context) *DataLoader {
	return ctx.Value(LoadersKey).(*DataLoader)
}
