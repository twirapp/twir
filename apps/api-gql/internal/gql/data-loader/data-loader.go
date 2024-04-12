package data_loader

import (
	"context"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/libs/cache/twitch"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	LoadersKey = ctxKey("dataloaders")
)

type DataLoader struct {
	cachedTwitchClient *twitch.CachedTwitchClient
	helixUserLoader    *dataloadgen.Loader[string, *helix.User]
}

type Opts struct {
	CachedTwitchClient *twitch.CachedTwitchClient
}

func New(opts Opts) *DataLoader {
	loader := &DataLoader{
		cachedTwitchClient: opts.CachedTwitchClient,
	}

	loader.helixUserLoader = dataloadgen.NewLoader(
		loader.getHelixUsersByIds,
		dataloadgen.WithWait(time.Millisecond),
	)

	return loader
}

// For returns the dataloader for a given context
func GetLoaderForRequest(ctx context.Context) *DataLoader {
	return ctx.Value(LoadersKey).(*DataLoader)
}
