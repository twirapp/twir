package impl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/deps"
	"github.com/satont/tsuwari/apps/api-twirp/internal/impl/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/commands"
	"github.com/satont/tsuwari/libs/grpc/generated/api/community"
	"github.com/satont/tsuwari/libs/grpc/generated/api/events"
	"github.com/satont/tsuwari/libs/grpc/generated/api/greetings"
	"github.com/satont/tsuwari/libs/grpc/generated/api/keywords"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"github.com/satont/tsuwari/libs/grpc/generated/api/modules_obs_websocket"
	"github.com/satont/tsuwari/libs/grpc/generated/api/modules_sr"
	"github.com/satont/tsuwari/libs/grpc/generated/api/modules_tts"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Api struct {
	*deps.Deps
	*integrations.Integrations
}

func (c *Api) KeywordsGetAll(ctx context.Context, empty *emptypb.Empty) (*keywords.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) KeywordsGetById(ctx context.Context, request *keywords.GetByIdRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) KeywordsCreate(ctx context.Context, request *keywords.CreateRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) KeywordsDelete(ctx context.Context, request *keywords.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) KeywordsPut(ctx context.Context, request *keywords.PutRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) KeywordsPatch(ctx context.Context, request *keywords.PatchRequest) (*keywords.Keyword, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesOBSWebsocketGet(ctx context.Context, empty *emptypb.Empty) (*modules_obs_websocket.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesOBSWebsocketPost(ctx context.Context, request *modules_obs_websocket.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesTTSGet(ctx context.Context, empty *emptypb.Empty) (*modules_tts.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesTTSPost(ctx context.Context, request *modules_tts.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesTTSGetInfo(ctx context.Context, empty *emptypb.Empty) (*modules_tts.GetInfoResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesSRGet(ctx context.Context, empty *emptypb.Empty) (*modules_sr.GetResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesSRGetSearch(ctx context.Context, request *modules_sr.GetSearchRequest) (*modules_sr.GetSearchResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) ModulesSRPost(ctx context.Context, request *modules_sr.PostRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	return &bots.BotInfo{
		IsMod:   false,
		BotId:   "123",
		BotName: "",
		Enabled: false,
	}, nil
}

func (c *Api) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsGetAll(ctx context.Context, empty *emptypb.Empty) (*commands.CommandsGetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsGetById(ctx context.Context, request *commands.GetByIdRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsCreate(ctx context.Context, request *commands.CreateRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsDelete(ctx context.Context, request *commands.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsPut(ctx context.Context, request *commands.PutRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommandsPatch(ctx context.Context, request *commands.PatchRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommunityGetUsers(ctx context.Context, request *community.GetUsersRequest) (*community.GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) CommunityResetStats(ctx context.Context, request *community.ResetStatsRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsGetAll(ctx context.Context, empty *emptypb.Empty) (*events.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsGetById(ctx context.Context, request *events.GetByIdRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsCreate(ctx context.Context, request *events.CreateRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsDelete(ctx context.Context, request *events.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsPut(ctx context.Context, request *events.PutRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) EventsPatch(ctx context.Context, request *events.PatchRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsGetAll(ctx context.Context, empty *emptypb.Empty) (*greetings.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsGetById(ctx context.Context, request *greetings.GetByIdRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsCreate(ctx context.Context, request *greetings.CreateRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsDelete(ctx context.Context, request *greetings.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsPut(ctx context.Context, request *greetings.PutRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Api) GreetingsPatch(ctx context.Context, request *greetings.PatchRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewApi(opts Opts) *Api {
	d := &deps.Deps{
		Redis: opts.Redis,
		Db:    opts.DB,
	}

	return &Api{
		Deps: d,
		Integrations: &integrations.Integrations{
			Deps: d,
		},
	}
}
