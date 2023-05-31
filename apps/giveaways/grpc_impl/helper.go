package grpc_impl

import (
	"context"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/giveaways"
)

func handleKeywordGiveawayMessage(
	ctx context.Context,
	data *giveaways.HandleChatMessageRequest,
	giveaway *model.ChannelGiveaway,
) error {

}

func handleRundomNumberGiveawayMessage(
	ctx context.Context,
	data *giveaways.HandleChatMessageRequest,
	giveaway *model.ChannelGiveaway,
) error {

}
