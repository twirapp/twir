package impl

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/satont/tsuwari/libs/grpc/generated/api/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/api/commands"
	"github.com/satont/tsuwari/libs/grpc/generated/api/community"
	"github.com/satont/tsuwari/libs/grpc/generated/api/events"
	"github.com/satont/tsuwari/libs/grpc/generated/api/greetings"
	"github.com/satont/tsuwari/libs/grpc/generated/api/integrations_donate_stream"
	"github.com/satont/tsuwari/libs/grpc/generated/api/meta"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Api struct {
	redis *redis.Client
	db    *gorm.DB
}

func (a *Api) BotInfo(ctx context.Context, meta *meta.BaseRequestMeta) (*bots.BotInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) BotJoinPart(ctx context.Context, request *bots.BotJoinPartRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsGetAll(ctx context.Context, empty *emptypb.Empty) (*commands.CommandsGetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsGetById(ctx context.Context, request *commands.GetByIdRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsCreate(ctx context.Context, request *commands.CreateRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsDelete(ctx context.Context, request *commands.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsPut(ctx context.Context, request *commands.PutRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommandsPatch(ctx context.Context, request *commands.PatchRequest) (*commands.Command, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommunityGetUsers(ctx context.Context, request *community.GetUsersRequest) (*community.GetUsersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) CommunityResetStats(ctx context.Context, request *community.ResetStatsRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsGetAll(ctx context.Context, empty *emptypb.Empty) (*events.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsGetById(ctx context.Context, request *events.GetByIdRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsCreate(ctx context.Context, request *events.CreateRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsDelete(ctx context.Context, request *events.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsPut(ctx context.Context, request *events.PutRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) EventsPatch(ctx context.Context, request *events.PatchRequest) (*events.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsGetAll(ctx context.Context, empty *emptypb.Empty) (*greetings.GetAllResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsGetById(ctx context.Context, request *greetings.GetByIdRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsCreate(ctx context.Context, request *greetings.CreateRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsDelete(ctx context.Context, request *greetings.DeleteRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsPut(ctx context.Context, request *greetings.PutRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) GreetingsPatch(ctx context.Context, request *greetings.PatchRequest) (*greetings.Greeting, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) IntegrationsDonateStreamGet(ctx context.Context, empty *emptypb.Empty) (*integrations_donate_stream.DonateStreamResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *Api) IntegrationsDonateStreamPostSecret(ctx context.Context, request *integrations_donate_stream.DonateStreamPostSecretRequest) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

type Opts struct {
	Redis *redis.Client
	DB    *gorm.DB
}

func NewApi(opts Opts) *Api {
	return &Api{
		redis: opts.Redis,
		db:    opts.DB,
	}
}
