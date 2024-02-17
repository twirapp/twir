package messagehandler

import (
	"context"
	"log/slog"
	"regexp"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/grpc/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *MessageHandler) handleCommand(ctx context.Context, msg handleMessage) error {
	if !strings.HasPrefix(msg.GetMessage().GetText(), "!") {
		return nil
	}

	emotes := make([]*parser.Message_Emote, 0, len(msg.GetMessage().GetFragments()))

	for _, f := range msg.GetMessage().GetFragments() {
		if f.GetType() != shared.FragmentType_EMOTE {
			continue
		}

		re := regexp.MustCompile(regexp.QuoteMeta(f.GetText()))
		var emotePositions []*parser.Message_EmotePosition

		for _, match := range re.FindAllStringSubmatchIndex(msg.GetMessage().GetText(), -1) {
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
				Id:        f.GetEmote().GetId(),
				Name:      f.GetText(),
				Positions: emotePositions,
			},
		)
	}

	mentions := make([]*parser.Message_Mention, 0, len(msg.GetMessage().GetFragments()))
	for _, m := range msg.GetMessage().GetFragments() {
		if m.Type != shared.FragmentType_MENTION {
			continue
		}

		mentions = append(
			mentions, &parser.Message_Mention{
				UserId:    m.GetMention().GetUserId(),
				UserName:  m.GetMention().GetUserName(),
				UserLogin: m.GetMention().GetUserLogin(),
			},
		)
	}

	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.GetChatterUserId(),
			Name:        msg.GetChatterUserLogin(),
			DisplayName: msg.GetChatterUserName(),
			Badges:      createUserBadges(msg.Badges),
			Color:       msg.GetColor(),
		},
		Channel: &parser.Channel{
			Id:   msg.GetBroadcasterUserId(),
			Name: msg.GetBroadcasterUserLogin(),
		},
		Message: &parser.Message{
			Id:       msg.GetMessageId(),
			Text:     msg.GetMessage().GetText(),
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

			r := r

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
