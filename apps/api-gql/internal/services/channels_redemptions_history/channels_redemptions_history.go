package channels_redemptions_history

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/cache/twitch"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsRedemptionsHistory channelsredemptionshistory.Repository
	CachedTwitchClient         *twitch.CachedTwitchClient
}

func New(opts Opts) *Service {
	return &Service{
		channelsRedemptionsHistory: opts.ChannelsRedemptionsHistory,
		cachedTwitchClient:         opts.CachedTwitchClient,
	}
}

type Service struct {
	channelsRedemptionsHistory channelsredemptionshistory.Repository
	cachedTwitchClient         *twitch.CachedTwitchClient
}

type GetManyInput struct {
	ChannelID    string
	Page         int
	PerPage      int
	UserNameLike *string
	RewardsIDs   []string
}

type GetManyPayload struct {
	Items []entity.ChannelRedemptionHistoryItem
	Total uint64
}

func (c *Service) GetMany(
	ctx context.Context,
	input GetManyInput,
) (GetManyPayload, error) {
	page := input.Page
	perPage := input.PerPage
	if perPage == 0 || perPage > 1000 {
		perPage = 20
	}

	repoInput := channelsredemptionshistory.GetManyInput{
		ChannelID:  input.ChannelID,
		Page:       page,
		PerPage:    perPage,
		RewardsIDs: input.RewardsIDs,
	}

	if input.UserNameLike != nil {
		users, err := c.cachedTwitchClient.SearchChannels(ctx, *input.UserNameLike)
		if err != nil {
			return GetManyPayload{}, err
		}

		userIDs := make([]string, len(users))
		for i, user := range users {
			userIDs[i] = user.ID
		}

		repoInput.UserIDs = userIDs
	}

	items, err := c.channelsRedemptionsHistory.GetMany(ctx, repoInput)
	if err != nil {
		return GetManyPayload{}, err
	}

	convertedItems := make([]entity.ChannelRedemptionHistoryItem, len(items.Items))
	for i, item := range items.Items {
		convertedItems[i] = entity.ChannelRedemptionHistoryItem{
			ChannelID:    item.ChannelID,
			UserID:       item.UserID,
			RewardID:     item.RewardID,
			RewardPrompt: item.RewardPrompt,
			RewardTitle:  item.RewardTitle,
			RewardCost:   item.RewardCost,
			RedeemedAt:   item.CreatedAt,
		}
	}

	return GetManyPayload{
		Items: convertedItems,
		Total: items.Total,
	}, nil
}

type CountInput struct {
	ChannelID  *string
	RewardsIDs []string
}

func (c *Service) Count(
	ctx context.Context,
	input CountInput,
) (uint64, error) {
	return c.channelsRedemptionsHistory.Count(
		ctx,
		channelsredemptionshistory.CountInput{
			ChannelID:  input.ChannelID,
			RewardsIDs: input.RewardsIDs,
		},
	)
}
