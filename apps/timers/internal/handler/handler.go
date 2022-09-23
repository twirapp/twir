package handler

import (
	"context"
	"encoding/json"
	"log"
	"time"

	model "tsuwari/models"
	"tsuwari/timers/internal/types"
	"tsuwari/twitch"

	"github.com/go-co-op/gocron"
	"github.com/go-redis/redis/v9"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/nicklaw5/helix"
	"github.com/satont/tsuwari/nats/bots"
	"github.com/satont/tsuwari/nats/parser"
	"gorm.io/gorm"
)

type Stream struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	UserLogin    string    `json:"user_login"`
	UserName     string    `json:"user_name"`
	GameID       string    `json:"game_id"`
	GameName     string    `json:"game_name"`
	TagIDs       []string  `json:"tag_ids"`
	IsMature     bool      `json:"is_mature"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	ViewerCount  int       `json:"viewer_count"`
	StartedAt    time.Time `json:"started_at"`
	Language     string    `json:"language"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Messages int `json:"parsedMessages"`
}

type Handler struct {
	redis *redis.Client
	twitch *twitch.Twitch
	nats *nats.Conn
	db *gorm.DB
}

func New(redis *redis.Client, twitch *twitch.Twitch, nats *nats.Conn, db *gorm.DB) *Handler {
	return &Handler{redis: redis, twitch: twitch, nats: nats, db: db}
}

func (c *Handler) Handle(t *types.Timer, j gocron.Job) {
	streamString, err := c.redis.Get(context.TODO(), "streams:" + t.Model.ChannelID).Result()

	if err != nil {
		log.Fatal(err)
		return
	}

	streamData := Stream{}

	if err = json.Unmarshal([]byte(streamString), &streamData); err != nil {
		log.Fatal(err)
		return
	}

	if t.Model.MessageInterval > 0 && t.Model.LastTriggerMessageNumber - int32(streamData.Messages) + t.Model.MessageInterval > 0 {
		return
	}
	
	
	users, err := c.twitch.Client.GetUsers(&helix.UsersParams{
		IDs: []string{t.Model.ChannelID},
	})
	
	if err != nil || len(users.Data.Users) == 0 {
		return;
	}
	
	user := users.Data.Users[0]
	
	rawMessage := t.Model.Responses[t.SendIndex]

	requestBytes, err := proto.Marshal(&parser.ParseResponseRequest{
		Sender: &parser.Sender{Id: "", Name: "bot", DisplayName: "Bot", Badges: []string{"BROADCASTER"}},
		Channel: &parser.Channel{Id: user.ID, Name: user.Login},
		Message: &parser.Message{Text: rawMessage},
	})

	if err != nil {
		log.Fatal(err)
		return
	}
	
	response, err := c.nats.Request("parser.parseTextResponse", requestBytes, 5*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	responseData := parser.ParseResponseResponse{}

	err = proto.Unmarshal(response.Data, &responseData)

	if err != nil {
		log.Fatal(err)
		return
	}

	botsRequest := bots.SendMessage{
		ChannelId: user.ID,
		ChannelName: user.Login,
		Message: rawMessage,
	}
	bytes, _ := proto.Marshal(&botsRequest)
	c.nats.Publish("bots.sendMessage", bytes)

	nextIndex := t.SendIndex + 1

	if nextIndex+1 <= len(t.Model.Responses) {
		t.SendIndex = nextIndex
	} else {
		t.SendIndex = 0
	}

	t.Model.LastTriggerMessageNumber = int32(streamData.Messages)

	err = c.db.
		Model(&model.ChannelsTimers{}).
		Where(`"id" = ?`, t.Model.ID).
		Updates(model.ChannelsTimers{LastTriggerMessageNumber: int32(streamData.Messages)}).
		Error

	if err != nil {
		log.Fatal(err)
	}
}