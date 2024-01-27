package messagehandler

import (
	"context"
	"log/slog"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/grpc/parser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *MessageHandler) handleCommand(ctx context.Context, msg handleMessage) error {
	if !strings.HasPrefix(msg.GetMessage().GetText(), "!") {
		return nil
	}

	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.GetChatterUserId(),
			Name:        msg.GetChatterUserLogin(),
			DisplayName: msg.GetChatterUserName(),
			Badges:      createUserBadges(msg.Badges),
		},
		Channel: &parser.Channel{
			Id:   msg.GetBroadcasterUserId(),
			Name: msg.GetBroadcasterUserLogin(),
		},
		Message: &parser.Message{
			Id:     msg.GetMessageId(),
			Text:   msg.GetMessage().GetText(),
			Emotes: []*parser.Message_Emote{},
		},
	}

	resp, err := c.parserGrpc.ProcessCommand(ctx, requestStruct)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil
		}

		c.logger.Error("cannot process command", slog.Any("err", err))
		return nil
	}

	if resp.GetKeepOrder() {
		for _, r := range resp.GetResponses() {
			if r == "" || r == " " {
				continue
			}

			err := c.twitchActions.SendMessage(
				ctx,
				twitchactions.SendMessageOpts{
					BroadcasterID:        msg.GetBroadcasterUserId(),
					SenderID:             msg.DbChannel.BotID,
					Message:              r,
					ReplyParentMessageID: lo.If(resp.GetIsReply(), msg.GetMessageId()).Else(""),
				},
			)
			if err != nil {
				c.logger.Error("cannot send message", slog.Any("err", err))
			}
		}
	} else {
		for _, r := range resp.GetResponses() {
			if r == "" || r == " " {
				continue
			}

			go func() {
				e := c.twitchActions.SendMessage(
					ctx,
					twitchactions.SendMessageOpts{
						BroadcasterID:        msg.GetBroadcasterUserId(),
						SenderID:             msg.DbChannel.BotID,
						Message:              r,
						ReplyParentMessageID: lo.If(resp.GetIsReply(), msg.GetMessageId()).Else(""),
					},
				)
				if e != nil {
					c.logger.Error("cannot send message", slog.Any("err", e))
				}
			}()
		}
	}

	return nil
}
