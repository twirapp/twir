package obs_websocket_module

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	gojson "github.com/goccy/go-json"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/api"
	obsentity "github.com/twirapp/twir/libs/entities/obs"
	"github.com/twirapp/twir/libs/redis_keys"
	channelsmodulesobswebsocket "github.com/twirapp/twir/libs/repositories/channels_modules_obs_websocket"
	"github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	ObsWebsocketRepository channelsmodulesobswebsocket.Repository
	WsRouter               wsrouter.WsRouter
	UsersRepository        users.Repository
	Bus                    *buscore.Bus
	Logger                 *slog.Logger
	KV                     kv.KV
}

type Service struct {
	obsWebsocketRepository channelsmodulesobswebsocket.Repository
	wsRouter               wsrouter.WsRouter
	usersRepository        users.Repository
	bus                    *buscore.Bus
	kv                     kv.KV
}

func New(opts Opts) *Service {
	s := &Service{
		obsWebsocketRepository: opts.ObsWebsocketRepository,
		wsRouter:               opts.WsRouter,
		usersRepository:        opts.UsersRepository,
		bus:                    opts.Bus,
		kv:                     opts.KV,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s.bus.Api.TriggerObsCommand.SubscribeGroup(
					"api",
					func(ctx context.Context, data api.TriggerObsCommand) (struct{}, error) {
						cmd := mapBusCommandToEntity(data)
						return struct{}{}, s.wsRouter.Publish(
							createCommandsSubscriptionKey(data.ChannelId),
							cmd,
						)
					},
				)

				opts.Logger.Info("Subscribed to TriggerObsCommand events")

				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.bus.Api.TriggerObsCommand.Unsubscribe()

				opts.Logger.Info("Unsubscribed from TriggerObsCommand events")

				return nil
			},
		},
	)

	return s
}

// mapBusCommandToEntity maps bus API command to entity command
func mapBusCommandToEntity(cmd api.TriggerObsCommand) obsentity.ObsWebsocketCommand {
	var action obsentity.ObsWebsocketCommandAction
	switch cmd.Action {
	case api.ObsCommandActionSetScene:
		action = obsentity.ObsWebsocketCommandActionSetScene
	case api.ObsCommandActionToggleSource:
		action = obsentity.ObsWebsocketCommandActionToggleSource
	case api.ObsCommandActionToggleAudio:
		action = obsentity.ObsWebsocketCommandActionToggleAudio
	case api.ObsCommandActionSetVolume:
		action = obsentity.ObsWebsocketCommandActionSetVolume
	case api.ObsCommandActionIncreaseVolume:
		action = obsentity.ObsWebsocketCommandActionIncreaseVolume
	case api.ObsCommandActionDecreaseVolume:
		action = obsentity.ObsWebsocketCommandActionDecreaseVolume
	case api.ObsCommandActionEnableAudio:
		action = obsentity.ObsWebsocketCommandActionEnableAudio
	case api.ObsCommandActionDisableAudio:
		action = obsentity.ObsWebsocketCommandActionDisableAudio
	case api.ObsCommandActionStartStream:
		action = obsentity.ObsWebsocketCommandActionStartStream
	case api.ObsCommandActionStopStream:
		action = obsentity.ObsWebsocketCommandActionStopStream
	}

	return obsentity.ObsWebsocketCommand{
		Action:      action,
		Target:      cmd.Target,
		VolumeValue: cmd.VolumeValue,
		VolumeStep:  cmd.VolumeStep,
	}
}

func createSettingsSubscriptionKey(channelID string) string {
	return "obs-websocket:settings:" + channelID
}

func createCommandsSubscriptionKey(channelID string) string {
	return "obs-websocket:commands:" + channelID
}

func (s *Service) GetObsWebsocketData(
	ctx context.Context,
	channelID string,
) (*obsentity.ObsWebsocketData, error) {
	module, err := s.obsWebsocketRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get obs websocket data: %w", err)
	}

	if module.IsNil() {
		return nil, nil
	}

	// Ensure slices are not nil
	scenes := module.Scenes
	if scenes == nil {
		scenes = []string{}
	}
	sources := module.Sources
	if sources == nil {
		sources = []string{}
	}
	audioSources := module.AudioSources
	if audioSources == nil {
		audioSources = []string{}
	}

	return &obsentity.ObsWebsocketData{
		ServerPort:     module.ServerPort,
		ServerAddress:  module.ServerAddress,
		ServerPassword: module.ServerPassword,
		Sources:        sources,
		AudioSources:   audioSources,
		Scenes:         scenes,
	}, nil
}

type UpdateInput struct {
	ChannelID      string
	ServerPort     *int
	ServerAddress  *string
	ServerPassword *string
	Scenes         *[]string
	Sources        *[]string
	AudioSources   *[]string
}

func (s *Service) UpdateObsWebsocket(
	ctx context.Context,
	input UpdateInput,
) error {
	_, err := s.obsWebsocketRepository.Upsert(
		ctx,
		channelsmodulesobswebsocket.UpsertInput{
			ChannelID:      input.ChannelID,
			ServerPort:     input.ServerPort,
			ServerAddress:  input.ServerAddress,
			ServerPassword: input.ServerPassword,
			Scenes:         input.Scenes,
			Sources:        input.Sources,
			AudioSources:   input.AudioSources,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update obs websocket: %w", err)
	}

	// Publish update to subscribers
	data, err := s.GetObsWebsocketData(ctx, input.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get obs websocket data for publish: %w", err)
	}

	if err := s.wsRouter.Publish(createSettingsSubscriptionKey(input.ChannelID), data); err != nil {
		return fmt.Errorf("failed to publish obs websocket update: %w", err)
	}

	return nil
}

// SettingsSubscriptionSignalerByApiKey subscribes to obs websocket settings changes by API key
// This is used by the overlay to receive settings updates
func (s *Service) SettingsSubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan obsentity.ObsWebsocketData, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user.IsNil() {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createSettingsSubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	chann := make(chan obsentity.ObsWebsocketData, 1)

	// get initial settings
	initialSettings, err := s.GetObsWebsocketData(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get obs websocket data: %w", err)
	}

	chann <- *initialSettings

	go func() {
		defer func() {
			_ = wsRouterSub.Unsubscribe()
			close(chann)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-wsRouterSub.GetChannel():
				var newSettings obsentity.ObsWebsocketData
				if err := gojson.Unmarshal(data, &newSettings); err != nil {
					panic(err)
				}

				chann <- newSettings
			}
		}
	}()

	return chann, nil
}

type UpdateFromOverlayInput struct {
	Scenes       []string
	Sources      []string
	AudioSources []string
}

// UpdateFromOverlay updates OBS data from overlay and publishes changes to subscribers
func (s *Service) UpdateFromOverlay(
	ctx context.Context,
	apiKey string,
	input UpdateFromOverlayInput,
) error {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user.IsNil() {
		return fmt.Errorf("user not found for provided api key")
	}

	channelID := user.ID

	// Store data in database
	newValue, err := s.obsWebsocketRepository.Upsert(
		ctx,
		channelsmodulesobswebsocket.UpsertInput{
			ChannelID:    channelID,
			Scenes:       &input.Scenes,
			Sources:      &input.Sources,
			AudioSources: &input.AudioSources,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update sources in database: %w", err)
	}

	if err := s.wsRouter.Publish(createSettingsSubscriptionKey(channelID), newValue); err != nil {
		return fmt.Errorf("failed to publish obs websocket update: %w", err)
	}

	return nil
}

func (s *Service) SetConnectedState(ctx context.Context, apiKey string, connected bool) error {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user.IsNil() {
		return fmt.Errorf("user not found for provided api key")
	}

	key := redis_keys.ObsOverlayConnection(user.ID)

	return s.kv.Set(ctx, key, connected, kvoptions.WithExpire(5*time.Second))
}

func (s *Service) IsConnected(ctx context.Context, channelID string) (bool, error) {
	key := redis_keys.ObsOverlayConnection(channelID)
	return s.kv.Get(ctx, key).Bool()
}

// CommandsSubscriptionByApiKey subscribes to OBS commands by API key
func (s *Service) CommandsSubscriptionByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan obsentity.ObsWebsocketCommand, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user.IsNil() {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createCommandsSubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	chann := make(chan obsentity.ObsWebsocketCommand, 1)

	go func() {
		defer func() {
			_ = wsRouterSub.Unsubscribe()
			close(chann)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-wsRouterSub.GetChannel():
				var cmd obsentity.ObsWebsocketCommand
				if err := gojson.Unmarshal(data, &cmd); err != nil {
					continue
				}

				chann <- cmd
			}
		}
	}()

	return chann, nil
}
