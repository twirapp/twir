package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/RediSearch/redisearch-go/v2/redisearch"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	twircfg "github.com/satont/twir/libs/config"
	twirlogger "github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	chat_messages_store "github.com/twirapp/twir/libs/bus-core/chat-messages-store"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

var (
	ignoredBadges = []string{"broadcaster", "moderator", "vip", "subscriber"}
)

func main() {
	cfg, err := twircfg.New()
	if err != nil {
		panic(err)
	}
	logger := twirlogger.New(twirlogger.Opts{
		Env:     cfg.AppEnv,
		Service: "chat-messages-store",
	})

	url, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		panic("Wrong redis url")
	}
	rdb := redis.NewClient(
		&redis.Options{
			Addr:     url.Addr,
			Password: url.Password,
			DB:       url.DB,
			Username: url.Username,
		},
	)
	defer rdb.Close()
	rdb.Conn()

	pool := &redigo.Pool{Dial: func() (redigo.Conn, error) {
		return redigo.Dial("tcp", url.Addr, redigo.DialPassword(url.Password))
	}}
	rsc := redisearch.NewClientFromPool(pool, "chat-messages-store:message:")

	sc := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTagField("can_be_deleted")).
		AddField(redisearch.NewTextFieldOptions("text", redisearch.TextFieldOptions{Weight: 5.0})).
		AddField(redisearch.NewTagField("channel_id"))

	if err := rsc.Drop(); err != nil {
		panic(err)
	}

	if err := rsc.CreateIndex(sc); err != nil {
		panic(err)
	}

	nc, err := nats.Connect(cfg.NatsUrl)
	if err != nil {
		panic(err)
	}
	defer nc.Close()

	bus := buscore.NewNatsBus(nc)
	err = bus.ChatMessages.SubscribeGroup(
		"chat-messages-store",
		func(ctx context.Context, msg twitch.TwitchChatMessage) struct{} {
			canBeDeleted := !lo.SomeBy(msg.Badges, func(b twitch.ChatMessageBadge) bool {
				return lo.Contains(ignoredBadges, b.SetId)
			})

			doc := redisearch.NewDocument("chat-messages-store:message:"+msg.MessageId, 1.0)
			doc.Set("can_be_deleted", strconv.FormatBool(canBeDeleted)).
				Set("text", msg.Message.Text).
				Set("channel_id", msg.BroadcasterUserId).
				Set("message_id", msg.MessageId).
				Set("user_id", msg.ChatterUserId).
				Set("user_login", msg.ChatterUserLogin).
				Set("created_at", time.Now())

			if err := rsc.Index([]redisearch.Document{doc}...); err != nil {
				logger.Error("Error in indexing message", slog.String("err", err.Error()))
				return struct{}{}
			}

			err = rdb.Expire(ctx, "chat-messages-store:message:"+msg.MessageId, 60*time.Minute).
				Err()
			if err != nil {
				logger.Error("Error in setting ttl", slog.String("err", err.Error()))
				return struct{}{}
			}

			return struct{}{}
		},
	)
	if err != nil {
		panic(err)
	}
	defer bus.ChatMessages.Unsubscribe()

	err = bus.ChatMessagesStore.GetChatMessagesByTextForDelete.SubscribeGroup(
		"chat-messages-store",
		func(ctx context.Context, req chat_messages_store.GetChatMessagesByTextRequest) chat_messages_store.GetChatMessagesByTextResponse {
			if len(req.ChannelID) == 0 || len(req.Text) == 0 {
				return chat_messages_store.GetChatMessagesByTextResponse{
					Messages: []chat_messages_store.StoredChatMessage{},
				}
			}

			offset := 0
			limit := 100
			var allDocs []redisearch.Document
			for {
				docs, total, err := rsc.Search(
					redisearch.NewAggregateQuery().
						SetQuery(redisearch.NewQuery(fmt.Sprintf("(@deleted:{false}) (@can_be_deleted:{true}) (@channel_id:{%s}) (@text:*%s*)", req.ChannelID, req.Text))).
						Query.Limit(offset, limit).
						SetReturnFields("message_id", "channel_id",
							"user_id",
							"user_login",
							"text",
							"can_be_deleted",
						),
				)
				if err != nil {
					logger.Error("Error in searching", slog.String("err", err.Error()))
					return chat_messages_store.GetChatMessagesByTextResponse{
						Messages: []chat_messages_store.StoredChatMessage{},
					}
				}

				allDocs = append(allDocs, docs...)
				if len(allDocs) >= total {
					break
				}

				offset += limit
			}
			var chatMessages []chat_messages_store.StoredChatMessage
			for _, doc := range allDocs {
				chatMessages = append(chatMessages, chat_messages_store.StoredChatMessage{
					MessageID:    doc.Properties["message_id"].(string),
					ChannelID:    doc.Properties["channel_id"].(string),
					UserID:       doc.Properties["user_id"].(string),
					UserLogin:    doc.Properties["user_login"].(string),
					Text:         doc.Properties["text"].(string),
					CanBeDeleted: doc.Properties["can_be_deleted"].(string) == "true",
				})
			}

			return chat_messages_store.GetChatMessagesByTextResponse{
				Messages: chatMessages,
			}
		},
	)
	if err != nil {
		panic(err)
	}
	defer bus.ChatMessagesStore.GetChatMessagesByTextForDelete.Unsubscribe()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	fmt.Printf("Received signal: %v\n", sig)

	fmt.Println("Graceful shutdown complete")
}
