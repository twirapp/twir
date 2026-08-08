package dota

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/olahol/melody"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/twirapp/twir/apps/websockets/types"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	twirlogger "github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/users"
)

type Dota struct {
	manager *melody.Melody

	logger  *slog.Logger
	counter prometheus.Gauge
	twirBus *buscore.Bus
}

func New(
	logger *slog.Logger,
	twirBus *buscore.Bus,
	channelsRepository channels.Repository,
	usersRepository users.Repository,
) *Dota {
	m := melody.New()
	dota := &Dota{
		manager: m,
		logger:  logger,
		counter: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name:        "websockets_connections_count",
				ConstLabels: prometheus.Labels{"overlay": "dota"},
			},
		),
		twirBus: twirBus,
	}

	dota.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckChannelByApiKey(session, channelsRepository, usersRepository)
			if err != nil {
				if !errors.Is(err, helpers.ErrUserNotFound) {
					logger.Error("cannot check user by api key", twirlogger.Error(err))
				}
				return
			}

			dota.counter.Inc()
			session.Write([]byte(`{"eventName":"connected to dota namespace"}`))

			dota.sendCurrentState(session)
		},
	)

	dota.manager.HandleDisconnect(
		func(session *melody.Session) {
			dota.counter.Dec()
		},
	)

	http.HandleFunc("/overlays/dota", dota.HandleRequest)

	return dota
}

func (c *Dota) HandleRequest(w http.ResponseWriter, r *http.Request) {
	_ = c.manager.HandleRequest(w, r)
}

func (c *Dota) sendCurrentState(session *melody.Session) {
	userIDValue, ok := session.Get("userId")
	if !ok {
		return
	}
	channelID, ok := userIDValue.(string)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	response, err := c.twirBus.Dota.GetData.Request(ctx, busdota.GetDataRequest{ChannelID: channelID})
	if err != nil {
		c.logger.Error("cannot fetch initial dota state", twirlogger.Error(err))
		return
	}
	if response == nil {
		return
	}

	data := response.Data
	message := busapi.DotaStateUpdateMessage{
		ChannelID:      channelID,
		InGame:         data.InGame,
		Mmr:            data.Mmr,
		SessionWins:    data.SessionWins,
		SessionLosses:  data.SessionLosses,
		WinProbability: data.WinProbability,
		HeroName:       data.HeroName,
		MatchID:        data.MatchID,
		TeamIsRadiant:  data.TeamIsRadiant,
		TeamKnown:      data.TeamKnown,
	}

	bytes, err := json.Marshal(&types.WebSocketMessage{
		EventName: "dotaStateUpdate",
		Data:      message,
		CreatedAt: time.Now().UTC().String(),
	})
	if err != nil {
		c.logger.Error("cannot marshal initial dota state", twirlogger.Error(err))
		return
	}

	_ = session.Write(bytes)
}

func (c *Dota) SendEvent(channelId, eventName string, data any) error {
	message := &types.WebSocketMessage{
		EventName: eventName,
		Data:      data,
		CreatedAt: time.Now().UTC().String(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		c.logger.Error("cannot process message", twirlogger.Error(err))
		return err
	}

	err = c.manager.BroadcastFilter(
		bytes,
		func(session *melody.Session) bool {
			socketUserId, ok := session.Get("userId")
			return ok && socketUserId.(string) == channelId
		},
	)

	if err != nil {
		c.logger.Error("cannot broadcast message", twirlogger.Error(err))
		return err
	}

	return nil
}
