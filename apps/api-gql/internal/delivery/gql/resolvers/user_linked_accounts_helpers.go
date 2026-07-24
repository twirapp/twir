package resolvers

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
)

func (r *authenticatedUserResolver) linkedAccountsForChannel(ctx context.Context, channel channelentity.Channel) ([]gqlmodel.LinkedAccount, error) {
	if channel.IsNil() {
		return []gqlmodel.LinkedAccount{}, nil
	}
	if r.deps.ChannelPlatformBindingsService == nil {
		return nil, fmt.Errorf("channel platform binding service is not configured")
	}

	bindings, err := r.deps.ChannelPlatformBindingsService.List(ctx, channel.ID)
	if err != nil {
		return nil, fmt.Errorf("list linked platform bindings: %w", err)
	}

	accounts := make([]gqlmodel.LinkedAccount, 0, len(bindings))
	for _, binding := range bindings {
		account := gqlmodel.LinkedAccount{
			Platform:       binding.Binding.Platform.String(),
			PlatformUserID: binding.Profile.PlatformID,
			PlatformLogin:  binding.Profile.Login,
		}
		if binding.Profile.Avatar != "" {
			account.PlatformAvatar = &binding.Profile.Avatar
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
