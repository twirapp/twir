package buslistener

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/dota/internal/match"
	"github.com/twirapp/twir/apps/dota/internal/stats"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"go.uber.org/fx"
)

const steamID64Base int64 = 76561197960265728

type StatsProvider interface {
	WinProbability(context.Context, int64) (float64, error)
	NotablePlayers(context.Context, int64, string) ([]string, error)
	LastGame(context.Context, int64) (*stats.LastGame, error)
}

type SnapshotProvider interface {
	GetSnapshot(context.Context, uuid.UUID) (match.Snapshot, error)
}

type getDataQueue interface {
	SubscribeGroup(
		string,
		buscore.QueueSubscribeCallback[busdota.GetDataRequest, busdota.GetDataResponse],
	) error
	Unsubscribe()
}

type Opts struct {
	fx.In

	Bus        *buscore.Bus
	State      *match.StateMachine
	Repository dotarepository.Repository
	Stats      StatsProvider
	Logger     *slog.Logger
	Lifecycle  fx.Lifecycle
}

type BusListener struct {
	state      SnapshotProvider
	repository dotarepository.Repository
	stats      StatsProvider
	logger     *slog.Logger
}

func New(opts Opts) *BusListener {
	return newBusListener(
		opts.State,
		opts.Repository,
		opts.Stats,
		opts.Logger,
		opts.Lifecycle,
		opts.Bus.Dota.GetData,
	)
}

func newBusListener(
	state SnapshotProvider,
	repository dotarepository.Repository,
	stats StatsProvider,
	logger *slog.Logger,
	lifecycle fx.Lifecycle,
	queue getDataQueue,
) *BusListener {
	listener := &BusListener{
		state:      state,
		repository: repository,
		stats:      stats,
		logger:     logger,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return queue.SubscribeGroup("dota", listener.GetData)
		},
		OnStop: func(context.Context) error {
			queue.Unsubscribe()
			return nil
		},
	})

	return listener
}

func (l *BusListener) GetData(
	ctx context.Context,
	req busdota.GetDataRequest,
) (busdota.GetDataResponse, error) {
	channelID, err := uuid.Parse(req.ChannelID)
	if err != nil {
		return busdota.GetDataResponse{}, fmt.Errorf("parse Dota channel ID: %w", err)
	}

	settings, err := l.repository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, dotarepository.ErrNotFound) {
			return busdota.GetDataResponse{}, nil
		}

		return busdota.GetDataResponse{}, fmt.Errorf("get Dota settings: %w", err)
	}

	linked := settings.SteamAccountID != nil && *settings.SteamAccountID != ""
	response := busdota.GetDataResponse{
		Enabled:       settings.Enabled,
		Linked:        linked,
		Mmr:           settings.Mmr,
		SessionWins:   settings.SessionWins,
		SessionLosses: settings.SessionLosses,
	}

	snapshot, err := l.state.GetSnapshot(ctx, channelID)
	if err != nil {
		return response, fmt.Errorf("get Dota match snapshot: %w", err)
	}

	response.InGame = snapshot.InGame
	response.HeroName = snapshot.HeroName
	response.MatchID = snapshot.MatchID
	response.TeamIsRadiant = snapshot.IsRadiant
	response.RadiantScore = snapshot.RadiantScore
	response.DireScore = snapshot.DireScore
	response.GameTime = snapshot.GameTime

	if !settings.Enabled || !linked {
		return response, nil
	}

	if snapshot.InGame {
		if snapshot.MatchID == 0 {
			return response, nil
		}

		winProbability, err := l.stats.WinProbability(ctx, snapshot.MatchID)
		if err != nil {
			l.logger.WarnContext(
				ctx,
				"dota bus listener: failed to fetch win probability",
				logger.Error(err),
				slog.Int64("match_id", snapshot.MatchID),
			)
		} else {
			response.WinProbability = winProbability
		}

		streamerAccountID := snapshot.SteamAccountID
		if streamerAccountID == "" {
			streamerAccountID = *settings.SteamAccountID
		}

		notablePlayers, err := l.stats.NotablePlayers(ctx, snapshot.MatchID, streamerAccountID)
		if err != nil {
			l.logger.WarnContext(
				ctx,
				"dota bus listener: failed to fetch notable players",
				logger.Error(err),
				slog.Int64("match_id", snapshot.MatchID),
			)
		} else {
			response.NotablePlayers = notablePlayers
		}

		return response, nil
	}

	accountID, err := dotaAccountID(*settings.SteamAccountID)
	if err != nil {
		l.logger.WarnContext(
			ctx,
			"dota bus listener: invalid Steam account ID",
			logger.Error(err),
			slog.String("steam_account_id", *settings.SteamAccountID),
		)
		return response, nil
	}

	lastGame, err := l.stats.LastGame(ctx, accountID)
	if err != nil {
		l.logger.WarnContext(
			ctx,
			"dota bus listener: failed to fetch last game",
			logger.Error(err),
			slog.Int64("account_id", accountID),
		)
		return response, nil
	}
	if lastGame == nil {
		return response, nil
	}

	response.LastGame = &busdota.LastGameInfo{
		HeroName:  lastGame.HeroName,
		Kills:     lastGame.Kills,
		Deaths:    lastGame.Deaths,
		Assists:   lastGame.Assists,
		Win:       lastGame.Win,
		DurationS: lastGame.DurationS,
	}

	return response, nil
}

func dotaAccountID(steamAccountID string) (int64, error) {
	accountID, err := strconv.ParseInt(steamAccountID, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse Steam account ID: %w", err)
	}
	if accountID < 0 {
		return 0, fmt.Errorf("Steam account ID must not be negative")
	}
	if accountID >= steamID64Base {
		return accountID - steamID64Base, nil
	}

	return accountID, nil
}
