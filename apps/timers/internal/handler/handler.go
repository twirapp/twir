package handler

import (
	"context"

	"github.com/satont/twir/apps/timers/internal/types"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"

	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	db         *gorm.DB
	logger     *zap.Logger
	store      types.Store
	parserGrpc parser.ParserClient
	botsGrpc   bots.BotsClient
}

func New(
	db *gorm.DB,
	logger *zap.Logger,
	store types.Store,
	parserGrpc parser.ParserClient,
	botsGrpc bots.BotsClient,
) *Handler {
	return &Handler{db: db, logger: logger, store: store, parserGrpc: parserGrpc, botsGrpc: botsGrpc}
}

func (c *Handler) Handle(j gocron.Job) {
	t := c.store[j.Tags()[0]]

	streamData := model.ChannelsStreams{}

	err := c.db.Where(`"userId" = ?`, t.Model.ChannelID).First(&streamData).Error
	if err != nil && err != gorm.ErrRecordNotFound {
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
	for index, r := range t.Model.Responses {
		if index == t.SendIndex {
			timerResponse = r
			break
		}
	}

	if timerResponse == nil {
		return
	}

	req, err := c.parserGrpc.ParseTextResponse(
		context.Background(), &parser.ParseTextRequestData{
			Sender: &parser.Sender{
				Id:          "",
				Name:        "bot",
				DisplayName: "Bot",
				Badges:      []string{"BROADCASTER"},
			},
			Channel: &parser.Channel{Id: stream.UserId, Name: stream.UserLogin},
			Message: &parser.Message{Text: timerResponse.Text},
		},
	)
	if err != nil {
		return
	}

	for i := 0; i < len(req.Responses); i++ {
		message := req.Responses[i]

		c.botsGrpc.SendMessage(
			context.Background(), &bots.SendMessageRequest{
				ChannelId:   stream.UserId,
				ChannelName: &stream.UserLogin,
				Message:     message,
				IsAnnounce:  &timerResponse.IsAnnounce,
			},
		)
	}

	nextIndex := t.SendIndex + 1

	if nextIndex < len(t.Model.Responses) {
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
