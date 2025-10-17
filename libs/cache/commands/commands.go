package commands

import (
	"context"
	"time"

	kvotter "github.com/twirapp/kv/stores/otter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	commandswithgroupsandresponsesrepository "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
)

const KeyPrefix = "cache:twir:commands:channel:"

func New(
	repo commandswithgroupsandresponsesrepository.Repository,
	bus *buscore.Bus,
) *generic_cacher.GenericCacher[[]model.CommandWithGroupAndResponses] {
	return generic_cacher.New(
		generic_cacher.Opts[[]model.CommandWithGroupAndResponses]{
			KV:        kvotter.New(),
			KeyPrefix: KeyPrefix,
			LoadFn: func(ctx context.Context, key string) ([]model.CommandWithGroupAndResponses, error) {
				return repo.GetManyByChannelID(ctx, key)
			},
			Ttl:                24 * time.Hour,
			InvalidateSignaler: generic_cacher.NewBusCoreInvalidator(bus),
		},
	)
}
