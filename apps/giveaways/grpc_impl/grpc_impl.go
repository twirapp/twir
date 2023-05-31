package grpc_impl

import (
	"context"
	"errors"
	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
	"time"

	"github.com/satont/tsuwari/apps/giveaways/internal/types/services"
	"github.com/satont/tsuwari/libs/grpc/generated/giveaways"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GiveawaysGrpcServer struct {
	giveaways.UnimplementedGiveawaysServer

	services *services.Services
}

func NewServer() *GiveawaysGrpcServer {
	return &GiveawaysGrpcServer{}
}

func (server *GiveawaysGrpcServer) HandleChatMessage(
	ctx context.Context,
	data *giveaways.HandleChatMessageRequest,
) (*emptypb.Empty, error) {
	/*
		Firstly, we need to check that some giveaway is running for current channel
	*/
	giveaway := &model.ChannelGiveaway{}
	err := server.services.Gorm.WithContext(ctx).
		Where(`"channel_id" = ? AND "start_at" > ? AND "end_at < ?"`, data.Channel.Id, time.Now().String(), time.Now().String()).
		First(&giveaway).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	err = lo.If(giveaway.Type == model.ChannelGiveAwayTypeByKeyword, handleKeywordGiveawayMessage(ctx, data, giveaway)).
		Else(handleRundomNumberGiveawayMessage(ctx, data, giveaway))
	if err != nil {

	}

	return nil, nil
}

func (server *GiveawaysGrpcServer) SelectWinners(
	context.Context,
	*giveaways.ChooseWinnersRequest,
) (*giveaways.ChooseWinnersResponse, error) {
	return nil, nil
}

func (server *GiveawaysGrpcServer) ReSelectWinners(
	context.Context,
	*giveaways.ChooseWinnersRequest,
) (*giveaways.ChooseWinnersResponse, error) {
	return nil, nil
}
