package youtube_sr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) (*youtube.YoutubeSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.DB.Where(`"channelId" = ?`, channelId).First(&settings).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "settings not found")
	}

	data := youtube.YoutubeSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return &data, nil
}

func handlePost(channelId string, dto *youtube.YoutubeSettings, services types.Services) error {
	var existedSettings *model.ChannelModulesSettings
	err := services.DB.Where(`"channelId" = ?`, channelId).First(&existedSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if dto.BlackList != nil && dto.BlackList.Users != nil {
		twitchUsers := []helix.User{}
		twitchUsersChunks := lo.Chunk(dto.BlackList.Users, 100)
		mu := sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(twitchUsersChunks))

		for _, chunk := range twitchUsersChunks {
			go func(chunk []youtube.YoutubeBlacklistSettingsUsers) {
				defer wg.Done()
				req, _ := services.Twitch.Client.GetUsers(&helix.UsersParams{
					Logins: lo.Map(
						chunk,
						func(item youtube.YoutubeBlacklistSettingsUsers, _ int) string {
							return item.UserName
						},
					),
				})
				mu.Lock()
				twitchUsers = append(twitchUsers, req.Data.Users...)
				mu.Unlock()
			}(chunk)
		}

		wg.Wait()

		errors := []string{}
		for i, u := range dto.BlackList.Users {
			userInSlice, ok := lo.Find(twitchUsers, func(item helix.User) bool {
				return item.Login == strings.ToLower(u.UserName)
			})

			if !ok {
				errors = append(errors, fmt.Sprintf("user %s not found on twitch", u.UserName))
			} else {
				dto.BlackList.Users[i].UserName = userInSlice.Login
				dto.BlackList.Users[i].UserID = userInSlice.ID
			}
		}

		if len(errors) > 0 {
			return fiber.NewError(http.StatusNotFound, strings.Join(errors, ", "))
		}
	}

	bytes, err := json.Marshal(dto)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedSettings.ID == "" {
		newSettings := model.ChannelModulesSettings{
			ID:        uuid.NewV4().String(),
			Type:      "youtube_song_requests",
			Settings:  bytes,
			ChannelId: channelId,
		}

		err = services.DB.Create(&newSettings).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	} else {
		err = services.DB.Model(existedSettings).Updates(map[string]interface{}{"settings": bytes}).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	}
}

func handlePatch(
	channelId string,
	blackListType string,
	data any,
	services types.Services,
) error {
	settings := model.ChannelModulesSettings{}
	err := services.DB.Where(`"channelId" = ? AND type = ?`, channelId, "youtube_song_requests").
		First(&settings).
		Error
	if err != nil && err != gorm.ErrRecordNotFound {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if settings.ID == "" {
		bytes, err := json.Marshal(&youtube.YoutubeSettings{})
		if err != nil {
			services.Logger.Sugar().Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		newSettings := model.ChannelModulesSettings{
			ID:        uuid.NewV4().String(),
			Type:      "youtube_song_requests",
			ChannelId: channelId,
			Settings:  bytes,
		}
		err = services.DB.Save(&newSettings).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		settings = newSettings
	}

	originalSettings := youtube.YoutubeSettings{}
	err = json.Unmarshal(settings.Settings, &originalSettings)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if blackListType == "users" {
		user := data.(youtube.YoutubeBlacklistSettingsUsers)
		userReq, err := services.Twitch.Client.GetUsers(&helix.UsersParams{
			Logins: []string{strings.ToLower(user.UserName)},
		})

		if err != nil || len(userReq.Data.Users) == 0 {
			return fiber.NewError(http.StatusNotFound, "twitch user not found")
		}

		originalSettings.BlackList.Users = append(
			originalSettings.BlackList.Users,
			youtube.YoutubeBlacklistSettingsUsers{
				UserID:   userReq.Data.Users[0].ID,
				UserName: userReq.Data.Users[0].Login,
			},
		)
	}
	if blackListType == "channels" {
		originalSettings.BlackList.Channels = append(
			originalSettings.BlackList.Channels,
			data.(youtube.YoutubeBlacklistSettingsChannels),
		)
	}
	if blackListType == "songs" {
		originalSettings.BlackList.Songs = append(
			originalSettings.BlackList.Songs,
			data.(youtube.YoutubeBlacklistSettingsSongs),
		)
	}

	newSettingsBytes, err := json.Marshal(originalSettings)
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}
	settings.Settings = newSettingsBytes
	err = services.DB.Save(&settings).Error

	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
