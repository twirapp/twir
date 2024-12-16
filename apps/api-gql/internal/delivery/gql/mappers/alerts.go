package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func AlertEntityTo(e entity.Alert) gqlmodel.ChannelAlert {
	return gqlmodel.ChannelAlert{
		ID:           e.ID,
		Name:         e.Name,
		AudioID:      e.AudioID,
		AudioVolume:  &e.AudioVolume,
		CommandIds:   e.CommandIDS,
		RewardIds:    e.RewardIDS,
		GreetingsIds: e.GreetingsIDS,
		KeywordsIds:  e.KeywordsIDS,
	}
}
