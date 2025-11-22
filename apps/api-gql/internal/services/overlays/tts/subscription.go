package tts

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/bus-core/api"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func createSaySubscriptionKey(channelID string) string {
	return "tts:say:" + channelID
}

func createSkipSubscriptionKey(channelID string) string {
	return "tts:skip:" + channelID
}

func (s *Service) SettingsSubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan entity.TTSOverlay, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user == usermodel.Nil {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createSettingsSubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	chann := make(chan entity.TTSOverlay, 1)

	// get initial settings
	initialSettings, err := s.GetOrCreate(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tts overlay: %w", err)
	}

	chann <- initialSettings

	go func() {
		defer func() {
			wsRouterSub.Unsubscribe()
			close(chann)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-wsRouterSub.GetChannel():
				var newSettings entity.TTSOverlay
				if err := json.Unmarshal(data, &newSettings); err != nil {
					panic(err)
				}

				chann <- newSettings
			}
		}
	}()

	return chann, nil
}

func (s *Service) SaySubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan api.TriggerTtsSay, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user == usermodel.Nil {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createSaySubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	outputChan := make(chan api.TriggerTtsSay, 1)

	go func() {
		defer func() {
			wsRouterSub.Unsubscribe()
			close(outputChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-wsRouterSub.GetChannel():
				var msg api.TriggerTtsSay
				if err := json.Unmarshal(data, &msg); err != nil {
					panic(err)
				}

				if msg.ChannelId != user.ID {
					continue
				}

				outputChan <- msg
			}
		}
	}()

	return outputChan, nil
}

func (s *Service) SkipSubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan bool, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user == usermodel.Nil {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createSkipSubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	chann := make(chan bool, 1)

	go func() {
		defer func() {
			wsRouterSub.Unsubscribe()
			close(chann)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-wsRouterSub.GetChannel():
				var msg api.TriggerTtsSkip
				if err := json.Unmarshal(data, &msg); err != nil {
					panic(err)
				}

				if msg.ChannelId != user.ID {
					continue
				}

				chann <- true
			}
		}
	}()

	return chann, nil
}

