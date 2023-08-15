package handler

import (
	"context"
	"fmt"
	"github.com/satont/twir/apps/timers/internal/types"
	cfg "github.com/satont/twir/libs/config"
	"golang.org/x/exp/slog"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"

	"github.com/go-co-op/gocron"
	"gorm.io/gorm"
)

type Handler struct {
	db         *gorm.DB
	store      types.Store
	parserGrpc parser.ParserClient
	botsGrpc   bots.BotsClient
	config     *cfg.Config
}

func New(
	db *gorm.DB,
	store types.Store,
	parserGrpc parser.ParserClient,
	botsGrpc bots.BotsClient,
	config *cfg.Config,
) *Handler {
	return &Handler{db: db, store: store, parserGrpc: parserGrpc, botsGrpc: botsGrpc, config: config}
}

func (c *Handler) Handle(j gocron.Job) {
	t := c.store[j.Tags()[0]]

	stream := model.ChannelsStreams{}
	if err := c.db.Where(`"userId" = ?`, t.Model.ChannelID).Find(&stream).Error; err != nil {
		slog.Error(err.Error(), "userId", t.Model.ChannelID)
		return
	}

	if stream.ID == "" && c.config.AppEnv == "production" {
		return
	}

	fmt.Println("curr", t.SendIndex)
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
		context.Background(),
		&parser.ParseTextRequestData{
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

	for _, message := range req.Responses {
		_, err = c.botsGrpc.SendMessage(
			context.Background(),
			&bots.SendMessageRequest{
				ChannelId:   t.Model.ChannelID,
				ChannelName: &stream.UserLogin,
				Message:     message,
				IsAnnounce:  &timerResponse.IsAnnounce,
			},
		)
		if err != nil {
			slog.Error(err.Error(), "name", t.Model.Name, "userId", t.Model.ChannelID)
			return
		}
	}

	nextIndex := t.SendIndex + 1

	if nextIndex < len(t.Model.Responses) {
		t.SendIndex = nextIndex
	} else {
		t.SendIndex = 0
	}

	fmt.Println("setted to", t.SendIndex)

	t.Model.LastTriggerMessageNumber = int32(stream.ParsedMessages)

	err = c.db.
		Model(&model.ChannelsTimers{}).
		Where(`"id" = ?`, t.Model.ID).
		Updates(model.ChannelsTimers{LastTriggerMessageNumber: int32(stream.ParsedMessages)}).
		Error

	if err != nil {
		slog.Error(err.Error(), "name", t.Model.Name, "userId", t.Model.ChannelID)
	}
}
