package chat_client

import (
	"sync"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
)

func (c *ChatClient) getChannels() ([]string, error) {
	twitchClient, err := twitch.NewBotClient(c.Model.ID, c.services.Cfg, c.services.TokensGrpc)
	if err != nil {
		return nil, err
	}

	var botChannels []model.Channels
	c.services.DB.
		Where(
			`"botId" = ? AND "isEnabled" = ? AND "isBanned" = ? AND "isTwitchBanned" = ?`,
			c.Model.ID,
			true,
			false,
			false,
		).Find(&botChannels)
	channelsChunks := lo.Chunk(botChannels, 100)

	var twitchUsers []helix.User
	var twitchUsersMU sync.Mutex

	wg := utils.NewGoroutinesGroup()

	for _, chunk := range channelsChunks {
		chunk := chunk

		wg.Go(
			func() {
				usersIds := lo.Map(
					chunk,
					func(item model.Channels, _ int) string {
						return item.ID
					},
				)

				twitchUsersReq, err := twitchClient.GetUsers(
					&helix.UsersParams{
						IDs: usersIds,
					},
				)
				if err != nil {
					panic(err)
				}
				twitchUsersMU.Lock()
				twitchUsers = append(twitchUsers, twitchUsersReq.Data.Users...)
				twitchUsersMU.Unlock()
			},
		)
	}

	wg.Wait()

	names := make([]string, len(twitchUsers))
	for i, u := range twitchUsers {
		if u.Login == "" {
			continue
		}
		names[i] = u.Login
	}
	return names, nil
}
