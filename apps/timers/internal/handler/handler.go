package handler

import (
	"time"

	"github.com/satont/tsuwari/apps/timers/internal/types"

	"github.com/satont/tsuwari/libs/twitch"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/go-co-op/gocron"
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/tsuwari/libs/nats/bots"
	"github.com/satont/tsuwari/libs/nats/parser"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	twitch *twitch.Twitch
	nats   *nats.Conn
	db     *gorm.DB
	logger *zap.Logger
	store  types.Store
}

func New(
	twitch *twitch.Twitch,
	nats *nats.Conn,
	db *gorm.DB,
	logger *zap.Logger,
	store types.Store,
) *Handler {
	return &Handler{twitch: twitch, nats: nats, db: db, logger: logger, store: store}
}

func (c *Handler) Handle(j gocron.Job) {
	t := c.store[j.Tags()[0]]

	streamData := model.ChannelsStreams{}

	err := c.db.Where(`"userId" = ?`, t.Model.ChannelID).First(&streamData).Error
	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	if t.Model.MessageInterval > 0 &&
		t.Model.LastTriggerMessageNumber-int32(
			streamData.ParsedMessages,
		)+t.Model.MessageInterval > 0 {
		return
	}

	stream := model.ChannelsStreams{}

	err = c.db.Where(`"userId" = ?`, t.Model.ChannelID).First(&stream).Error

	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	var timerResponse *model.ChannelsTimersResponses
	for index, r := range *t.Model.Responses {
		if index == t.SendIndex {
			timerResponse = &r
			break
		}
	}

	if timerResponse == nil {
		return
	}

	requestBytes, protoError := proto.Marshal(&parser.ParseResponseRequest{
		Sender: &parser.Sender{
			Id:          "",
			Name:        "bot",
			DisplayName: "Bot",
			Badges:      []string{"BROADCASTER"},
		},
		Channel: &parser.Channel{Id: stream.UserId, Name: stream.UserLogin},
		Message: &parser.Message{Text: timerResponse.Text},
	})
	if protoError != nil {
		c.logger.Sugar().Error(err)
		return
	}

	response, natsError := c.nats.Request("parser.parseTextResponse", requestBytes, 5*time.Second)
	if natsError != nil {
		c.logger.Sugar().Error(natsError)
		return
	}
	responseData := parser.ParseResponseResponse{}

	err = proto.Unmarshal(response.Data, &responseData)

	if err != nil {
		c.logger.Sugar().Error(err)
		return
	}

	for i := 0; i < len(responseData.Responses); i++ {
		message := responseData.Responses[i]

		botsRequest := bots.SendMessage{
			ChannelId:   stream.UserId,
			ChannelName: &stream.UserLogin,
			Message:     message,
			IsAnnounce:  &timerResponse.IsAnnounce,
		}
		bytes, _ := proto.Marshal(&botsRequest)
		c.nats.Publish("bots.sendMessage", bytes)
	}

	nextIndex := t.SendIndex + 1

	if nextIndex < len(*t.Model.Responses) {
		t.SendIndex = nextIndex
	} else {
		t.SendIndex = 0
	}

	t.Model.LastTriggerMessageNumber = int32(streamData.ParsedMessages)

	err = c.db.
		Model(&model.ChannelsTimers{}).
		Where(`"id" = ?`, t.Model.ID).
		Updates(model.ChannelsTimers{LastTriggerMessageNumber: int32(streamData.ParsedMessages)}).
		Error

	if err != nil {
		c.logger.Sugar().Error(err)
	}
}
