package bot

import (
	"context"
	"fmt"
	"github.com/kr/pretty"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl_deps"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type Bot struct {
	*impl_deps.Deps
}

func (c *Bot) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	dashboard := ctx.Value("dashboardId")
	pretty.Println(dashboard)

	return &bots.BotInfo{
		IsMod:   false,
		BotId:   "123",
		BotName: fmt.Sprintf("%v", time.Now().UnixMilli()),
		Enabled: false,
	}, nil
}

func (c *Bot) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}
