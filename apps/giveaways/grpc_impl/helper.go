package grpc_impl

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/giveaways/internal/types/services"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/giveaways"
)

func handleKeywordGiveawayMessage(
	ctx context.Context,
	services *services.Services,
	data *giveaways.HandleChatMessageRequest,
	giveaway *model.ChannelGiveaway,
) error {
	msg := data.Text
	/*
		Check message to contain giveaway keyword
	*/
	if !strings.Contains(msg, giveaway.Keyword.String) {
		return nil
	}

	/*
		Get info about user, that want to enter the giveaway
	*/
	user := &model.Users{}
	err := services.Gorm.WithContext(ctx).Where(`"id" = ?`).First(user).Error
	if err != nil {
		return err
	}

	/*
		Get info about user roles on this channel
	*/
	userRoles := &[]model.ChannelRoleUser{}
	err = services.Gorm.WithContext(ctx).Where(`"userId" = ? AND "channelId" = ?`, user.ID, data.Channel.Id).Find(userRoles).Error
	if err != nil {
		return err
	}

	/*
		Get user stats for this channel
	*/

	/*
		Check user roles for this channel
	*/
	isSub := false
	for _, role := range *userRoles {
		if role.Role.Type == model.ChannelRoleTypeSubscriber {
			isSub = true
		}
	}

	/*
		Add participant to giveaway
	*/
	participant := &model.ChannelGiveawayParticipant{}
	participant.GiveawayID = giveaway.ID
	participant.ID = uuid.New().String()
	participant.UserID = user.ID
	participant.IsSubscriber = isSub
	participant.SubscriberTier = 1
	// participant.UserFollowSince

	return nil
}

func handleRundomNumberGiveawayMessage(
	ctx context.Context,
	data *giveaways.HandleChatMessageRequest,
	giveaway *model.ChannelGiveaway,
) error {
	return nil
}
