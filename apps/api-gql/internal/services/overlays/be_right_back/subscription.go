package be_right_back

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/bus-core/api"
	usermodel "github.com/twirapp/twir/libs/repositories/users/model"
)

func createSettingsSubscriptionKey(channelID string) string {
	return fmt.Sprintf("overlays:be_right_back:settings:%s", channelID)
}

func createStartSubscriptionKey(channelID string) string {
	return fmt.Sprintf("overlays:be_right_back:start:%s", channelID)
}

func createStopSubscriptionKey(channelID string) string {
	return fmt.Sprintf("overlays:be_right_back:stop:%s", channelID)
}

func (s *Service) SettingsSubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan entity.BeRightBackOverlay, error) {
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

	chann := make(chan entity.BeRightBackOverlay, 1)

	// get initial initialSettings
	initialSettings, err := s.GetOrCreate(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get be right back overlay: %w", err)
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
				var newSettings entity.BeRightBackOverlay
				if err := json.Unmarshal(data, &newSettings); err != nil {
					panic(err)
				}

				chann <- newSettings
			}
		}
	}()

	return chann, nil
}

func (s *Service) StartSubscriptionSignalerByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan api.TriggerBrbStart, error) {
	user, err := s.usersRepository.GetByApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by api key: %w", err)
	}
	if user == usermodel.Nil {
		return nil, fmt.Errorf("user not found for provided api key")
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createStartSubscriptionKey(user.ID)})
	if err != nil {
		return nil, err
	}

	outputChan := make(chan api.TriggerBrbStart, 1)

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
				var msg api.TriggerBrbStart
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

func (s *Service) StopSubscriptionSignalerByApiKey(
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

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createStopSubscriptionKey(user.ID)})
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
				var msg api.TriggerBrbStop
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
