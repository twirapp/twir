package messagehandler

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/libs/types/types/services/twitch"
	"github.com/twirapp/twir/libs/grpc/parser"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *MessageHandler) handleCommand(ctx context.Context, msg handleMessage) error {
	if !strings.HasPrefix(msg.Message.Text, "!") {
		return nil
	}

	emotes := make([]*parser.Message_Emote, 0, len(msg.Message.Fragments))

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE {
			continue
		}

		re := regexp.MustCompile(regexp.QuoteMeta(f.Text))
		var emotePositions []*parser.Message_EmotePosition

		for _, match := range re.FindAllStringSubmatchIndex(msg.Message.Text, -1) {
			emotePositions = append(
				emotePositions,
				&parser.Message_EmotePosition{
					Start: int64(match[0]),
					End:   int64(match[1]),
				},
			)
		}

		emotes = append(
			emotes, &parser.Message_Emote{
				Id:        f.Emote.Id,
				Name:      f.Text,
				Positions: emotePositions,
			},
		)
	}

	mentions := make([]*parser.Message_Mention, 0, len(msg.Message.Fragments))
	for _, m := range msg.Message.Fragments {
		if m.Type != twitch.FragmentType_MENTION {
			continue
		}

		mentions = append(
			mentions, &parser.Message_Mention{
				UserId:    m.Mention.UserId,
				UserName:  m.Mention.UserName,
				UserLogin: m.Mention.UserLogin,
			},
		)
	}

	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.ChatterUserId,
			Name:        msg.ChatterUserLogin,
			DisplayName: msg.ChatterUserName,
			Badges:      createUserBadges(msg.Badges),
			Color:       msg.Color,
		},
		Channel: &parser.Channel{
			Id:   msg.BroadcasterUserId,
			Name: msg.BroadcasterUserLogin,
		},
		Message: &parser.Message{
			Id:       msg.MessageId,
			Text:     msg.Message.Text,
			Emotes:   emotes,
			Mentions: mentions,
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
					BroadcasterID:        msg.BroadcasterUserId,
					SenderID:             msg.DbChannel.BotID,
					Message:              r,
					ReplyParentMessageID: lo.If(resp.GetIsReply(), msg.MessageId).Else(""),
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

			r := r

			go func() {
				e := c.twitchActions.SendMessage(
					ctx,
					twitchactions.SendMessageOpts{
						BroadcasterID:        msg.BroadcasterUserId,
						SenderID:             msg.DbChannel.BotID,
						Message:              r,
						ReplyParentMessageID: lo.If(resp.GetIsReply(), msg.MessageId).Else(""),
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
