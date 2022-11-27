package handlers

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/nats/watched"
	"github.com/satont/tsuwari/libs/twitch"
	uuid "github.com/satori/go.uuid"
)

func (c *Handlers) ProcessWatchedStreams(m *nats.Msg) {
	twitch := twitch.NewUserClient(twitch.UsersServiceOpts{
		Db:           c.db,
		ClientId:     c.cfg.TwitchClientId,
		ClientSecret: c.cfg.TwitchClientSecret,
	})

	requestData := watched.ParseRequest{}
	if err := proto.Unmarshal(m.Data, &requestData); err != nil {
		fmt.Println(err)
		return
	}
	api, err := twitch.CreateBot(requestData.BotId)
	if err != nil {
		fmt.Println(err)
		return
	}

	wg := sync.WaitGroup{}

	for _, channel := range requestData.ChannelsId {
		wg.Add(1)

		go func(channel string) {
			defer wg.Done()
			chatters := make([]helix.ChatChatter, 0)
			cursor := ""

			for {
				reqParams := &helix.GetChatChattersParams{
					BroadcasterID: channel,
					ModeratorID:   requestData.BotId,
					First:         "1000",
				}

				if cursor != "" {
					reqParams.After = cursor
				}

				req, err := api.GetChannelChatChatters(reqParams)
				if err != nil {
					break
				}

				chatters = append(chatters, req.Data.Chatters...)

				if req.Data.Pagination.Cursor == "" || len(req.Data.Chatters) == 0 {
					break
				}

				cursor = req.Data.Pagination.Cursor
			}

			chattersWg := sync.WaitGroup{}
			for _, chatter := range chatters {
				chattersWg.Add(1)

				go func(chatter helix.ChatChatter) {
					defer chattersWg.Done()
					user := model.Users{}
					if err := c.db.Where("id = ?", chatter.UserID).Preload("Stats").Find(&user).Error; err != nil {
						return
					}

					if user.ID == "" {
						c.db.Create(&model.Users{
							ID:     chatter.UserID,
							ApiKey: uuid.NewV4().String(),
						})
						c.db.Create(&model.UsersStats{
							ID:        uuid.NewV4().String(),
							UserID:    chatter.UserID,
							ChannelID: channel,
							Messages:  0,
							Watched:   0,
						})
					} else if user.Stats == nil {
						c.db.Create(&model.UsersStats{
							ID:        uuid.NewV4().String(),
							UserID:    chatter.UserID,
							ChannelID: channel,
							Messages:  0,
							Watched:   0,
						})
					} else {
						time := 5 * time.Minute
						c.db.Model(&model.UsersStats{}).
							Where("id = ?", user.Stats.ID).Select("*").
							Updates(map[string]any{
								"watched": user.Stats.Watched + time.Milliseconds(),
							})
					}
				}(chatter)
			}

			chattersWg.Wait()
		}(channel)
	}

	wg.Wait()
	m.Ack()
}
