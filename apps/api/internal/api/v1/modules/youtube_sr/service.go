package youtube_sr

import (
	"encoding/json"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	"net/http"
	"strings"
	"sync"

	ytsr "github.com/SherlockYigit/youtube-go"
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

func handleSearch(query string, searchType string) ([]youtube.SearchResult, error) {
	if query == "" {
		return []youtube.SearchResult{}, nil
	}

	if searchType != "video" && searchType != "channel" {
		return nil, fiber.NewError(400, "type can be only video or channel")
	}

	search := ytsr.Search(query, ytsr.SearchOptions{
		Limit: 20,
		Type:  searchType,
	})

	result := make([]youtube.SearchResult, 0)

	if len(search) == 0 {
		return result, nil
	}

	for _, item := range search {
		var res youtube.SearchResult
		if searchType == "video" {
			res = youtube.SearchResult{
				ID:        item.Video.Id,
				Title:     item.Video.Title,
				ThumbNail: item.Video.Thumbnail.Url,
			}
		}
		if searchType == "channel" {
			res = youtube.SearchResult{
				ID:        item.Channel.Id,
				Title:     item.Channel.Name,
				ThumbNail: item.Channel.Icon.Url,
			}
		}
		if res.ID != "" {
			result = append(result, res)
		}
	}

	return result, nil
}

func handlePost(channelId string, dto *youtube.YoutubeSettings, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	var existedSettings *model.ChannelModulesSettings
	err := services.DB.Where(`"channelId" = ?`, channelId).First(&existedSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if len(dto.BlackList.Users) > 0 {
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
		logger.Error(err)
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
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	} else {
		err = services.DB.Model(existedSettings).Updates(map[string]interface{}{"settings": bytes}).Error
		if err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	}
}
