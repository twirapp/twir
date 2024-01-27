package messagehandler

import (
	"context"
	"time"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *MessageHandler) handleStoreMessage(ctx context.Context, msg handleMessage) error {
	badges := createUserBadges(msg.GetBadges())

	canBeDeleted := !lo.Some(
		badges,
		[]string{"BROADCASTER", "MODERATOR", "SUBSCRIBER", "VIP"},
	)

	entity := model.ChannelChatMessage{
		MessageId:    msg.GetMessageId(),
		ChannelId:    msg.GetBroadcasterUserId(),
		UserId:       msg.GetChatterUserId(),
		UserName:     msg.GetChatterUserLogin(),
		Text:         msg.GetMessage().GetText(),
		CanBeDeleted: canBeDeleted,
		CreatedAt:    time.Now().UTC(),
	}

	err := c.gorm.WithContext(ctx).Create(&entity).Error
	return err
}
