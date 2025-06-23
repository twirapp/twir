package channels_emotes_usages

import (
	"context"

	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsEmotesUsagesRepository channelsemotesusagesrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		channelsEmotesUsagesRepository: opts.ChannelsEmotesUsagesRepository,
	}
}

type Service struct {
	channelsEmotesUsagesRepository channelsemotesusagesrepository.Repository
}

func (c *Service) Count(ctx context.Context) (uint64, error) {
	return c.channelsEmotesUsagesRepository.Count(ctx, channelsemotesusagesrepository.CountInput{})
}
