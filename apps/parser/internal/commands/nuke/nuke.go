package nuke

import (
	"context"
	"encoding/json"
	"fmt"
	"tsuwari/parser/internal/types"
	variables_cache "tsuwari/parser/internal/variablescache"

	redis "github.com/go-redis/redis/v9"
	"github.com/samber/lo"
	natsbots "github.com/satont/tsuwari/nats/bots"
	"github.com/tidwall/gjson"
	proto "google.golang.org/protobuf/proto"
)

type Message struct {
	ID           string `json:"messageId"`
	CanBeDeleted bool   `json:"canBeDeleted"`
	UserId       string `json:"userId"`
}

var Command = types.DefaultCommand{
	Command: types.Command{
		Name:        "nuke",
		Description: lo.ToPtr("Mass remove messages in chat by message content. Usage: <b>!nuke phrase</b>"),
		Permission:  "MODERATOR",
		Visible:     true,
		Module:      lo.ToPtr("CHANNEL"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		query := fmt.Sprintf("( (@channelId:{%v}) (@message:'%s') )", ctx.ChannelId, *ctx.Text)

		cmd := createSearchCmd(ctx.Services.Redis, query)
		result, err := cmd.Result()

		if err != nil {
			return nil
		}

		messages := []Message{}

		for _, m := range result {
			arr := fmt.Sprintf("%v", m)
			parsed := gjson.Parse(arr)
			if parsed.Type.String() != "JSON" {
				continue
			}
			v := Message{}
			el := parsed.Array()[0]

			err := json.Unmarshal([]byte(el.String()), &v)
			if err == nil {
				messages = append(messages, v)
			}
		}

		messages = lo.Filter(messages, func(m Message, _ int) bool {
			return m.CanBeDeleted
		})
		mappedMessages := lo.Map(messages, func(m Message, _ int) string {
			return m.ID
		})

		request := natsbots.DeleteMessagesRequest{
			ChannelId:   ctx.ChannelId,
			MessageIds:  mappedMessages,
			ChannelName: ctx.ChannelName,
		}

		marshaled, err := proto.Marshal(&request)

		if err == nil {
			ctx.Services.Nats.Publish("bots.deleteMessages", marshaled)
		}

		go func() {
			for _, msg := range messages {
				ctx.Services.Redis.Del(context.TODO(), fmt.Sprintf("messages:%v:%s", msg.UserId, msg.ID))
			}
		}()

		return nil
	},
}

func createSearchCmd(redisdb *redis.Client, query string) *redis.SliceCmd {
	c := context.TODO()
	n := redis.NewSliceCmd(c, "ft.search", "messages:index", query)
	redisdb.Process(c, n)
	return n
}
