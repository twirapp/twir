package bot

import (
	"context"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Bot struct {
	*deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	return &bots.BotInfo{
		IsMod:   false,
		BotId:   "123",
		BotName: "",
		Enabled: false,
	}, nil
}

func (c *Bot) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
