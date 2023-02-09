package grpc_impl

import (
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/twitch"
	"github.com/valyala/fasttemplate"
	"time"
)

type Data struct {
	UserName        string `json:"userName,omitempty"`
	UserDisplayName string

	RaidViewers int64 `json:"raidViewers,omitempty"`

	ResubMonths  int64  `json:"resubMonths"`
	ResubStreak  int64  `json:"resubStreak"`
	ResubMessage string `json:"resubMessage"`
	SubLevel     string `json:"subLevel"`

	OldTitle    string `json:"oldTitle"`
	NewTitle    string `json:"newTitle"`
	OldCategory string `json:"oldCategory"`
	NewCategory string `json:"newCategory"`

	StreamTitle    string `json:"streamTitle"`
	StreamCategory string `json:"streamCategory"`

	RewardName  string  `json:"rewardName"`
	RewardCost  string  `json:"rewardCost"`
	RewardInput *string `json:"rewardInput"`

	CommandName string `json:"commandName"`

	TargetUserName        string `json:"targetUserName"`
	TargetUserDisplayName string `json:"targetUserDisplayName"`
}

func hydrateStringWithData(str string, data Data) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteString(m)

	return s, nil
}

func (c *EventsGrpcImplementation) processOperations(channelId string, operations []model.EventOperation, data Data) {
	streamerApiClient, err := twitch.NewUserClient(channelId, *c.services.Cfg, c.services.TokensGrpc)
	if err != nil {
		c.services.Logger.Sugar().Error(err)
		return
	}

	for _, operation := range operations {
		if operation.Delay.Valid {
			duration := time.Duration(operation.Delay.Int64) * time.Second
			time.Sleep(duration)
		}

		switch operation.Type {
		case "SEND_MESSAGE":
			if operation.Input.Valid {
				msg, err := hydrateStringWithData(operation.Input.String, data)
				if err == nil {
					c.services.BotsGrpc.SendMessage(context.Background(), &bots.SendMessageRequest{
						ChannelId:   channelId,
						ChannelName: nil,
						Message:     msg,
						IsAnnounce:  nil,
					})
				}
			}
		case "BAN", "UNBAN":
			if data.UserName == "" {
				continue
			}

			user, err := streamerApiClient.GetUsers(&helix.UsersParams{
				Logins: []string{data.UserName},
			})

			if err != nil || len(user.Data.Users) == 0 {
				if err != nil {
					c.services.Logger.Sugar().Error(err)
				}
				continue
			}

			if operation.Type == "BAN" {
				streamerApiClient.BanUser(&helix.BanUserParams{
					BroadcasterID: channelId,
					ModeratorId:   channelId,
					Body: helix.BanUserRequestBody{
						Duration: 0,
						Reason:   "banned from twirapp",
						UserId:   user.Data.Users[0].ID,
					},
				})
			} else {
				streamerApiClient.UnbanUser(&helix.UnbanUserParams{
					BroadcasterID: channelId,
					ModeratorID:   channelId,
					UserID:        user.Data.Users[0].ID,
				})
			}
		case "BAN_RANDOM":
			randomOnlineUser := &model.UsersOnline{}
			err := c.services.DB.Where(`"channelId" = ?`, channelId).Find(&randomOnlineUser).Error
			if err != nil {
				c.services.Logger.Sugar().Error(err)
				continue
			}

			if randomOnlineUser == nil || !randomOnlineUser.UserId.Valid {
				continue
			}

			streamerApiClient.BanUser(&helix.BanUserParams{
				BroadcasterID: channelId,
				ModeratorId:   channelId,
				Body: helix.BanUserRequestBody{
					Duration: 0,
					Reason:   "randomly banned from twirapp",
					UserId:   randomOnlineUser.UserId.String,
				},
			})
		case "VIP", "UNVIP":
			if data.UserName == "" {
				continue
			}

			user, err := streamerApiClient.GetUsers(&helix.UsersParams{
				Logins: []string{data.UserName},
			})

			if err != nil || len(user.Data.Users) == 0 {
				if err != nil {
					c.services.Logger.Sugar().Error(err)
				}
				continue
			}

			if operation.Type == "VIP" {
				streamerApiClient.AddChannelVip(&helix.AddChannelVipParams{
					BroadcasterID: channelId,
					UserID:        user.Data.Users[0].ID,
				})
			} else {
				streamerApiClient.RemoveChannelVip(&helix.RemoveChannelVipParams{
					BroadcasterID: channelId,
					UserID:        user.Data.Users[0].ID,
				})

			}
		case "ENABLE_SUBMODE", "DISABLE_SUBMODE":
			streamerApiClient.UpdateChatSettings(&helix.UpdateChatSettingsParams{
				BroadcasterID: channelId,
				ModeratorID:   channelId,
				SubscriberMode: lo.ToPtr(lo.
					If(operation.Type == "ENABLE_SUBMODE", true).
					Else(false)),
			})
		case "ENABLE_EMOTEONLY", "DISABLE_EMOTEONLY":
			streamerApiClient.UpdateChatSettings(&helix.UpdateChatSettingsParams{
				BroadcasterID: channelId,
				ModeratorID:   channelId,
				EmoteMode: lo.ToPtr(lo.
					If(operation.Type == "ENABLE_EMOTEONLY", true).
					Else(false)),
			})
		}
	}
}
