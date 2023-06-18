package bot

import (
	"context"
	"github.com/kr/pretty"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Bot struct {
	*impl_deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	pretty.Println(ctx.Value("user"))

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
