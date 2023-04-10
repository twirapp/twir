package youtube_sr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	ytsr "github.com/SherlockYigit/youtube-go"
	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	youtube "github.com/satont/tsuwari/libs/types/types/api/modules"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) (*youtube.YouTubeSettings, error) {
	settings := model.ChannelModulesSettings{}
	err := services.DB.Where(`"channelId" = ? and "type" = ?`, channelId, "youtube_song_requests").First(&settings).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "settings not found")
	}

	data := youtube.YouTubeSettings{}
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

	search, err := ytsr.Search(query, ytsr.SearchOptions{
		Limit: 20,
		Type:  searchType,
	})

	result := make([]youtube.SearchResult, 0)
	if err != nil {
		return result, nil
	}

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

func handlePost(channelId string, dto *youtube.YouTubeSettings, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		return nil
	}

	var existedSettings *model.ChannelModulesSettings
	err = services.DB.Where(`"channelId" = ? AND "type" = ?`, channelId, "youtube_song_requests").First(&existedSettings).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if len(dto.DenyList.Users) > 0 {
		twitchUsers := []helix.User{}
		twitchUsersChunks := lo.Chunk(dto.DenyList.Users, 100)
		mu := sync.Mutex{}
		wg := sync.WaitGroup{}
		wg.Add(len(twitchUsersChunks))

		for _, chunk := range twitchUsersChunks {
			go func(chunk []youtube.YouTubeDenySettingsUsers) {
				defer wg.Done()
				req, _ := twitchClient.GetUsers(&helix.UsersParams{
					Logins: lo.Map(
						chunk,
						func(item youtube.YouTubeDenySettingsUsers, _ int) string {
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
		for i, u := range dto.DenyList.Users {
			userInSlice, ok := lo.Find(twitchUsers, func(item helix.User) bool {
				return item.Login == strings.ToLower(u.UserName)
			})

			if !ok {
				errors = append(errors, fmt.Sprintf("user %s not found on twitch", u.UserName))
			} else {
				dto.DenyList.Users[i].UserName = userInSlice.Login
				dto.DenyList.Users[i].UserID = userInSlice.ID
			}
		}

		if len(errors) > 0 {
			return fiber.NewError(http.StatusNotFound, strings.Join(errors, ", "))
		}
	}

	bytes, err := json.Marshal(*dto)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if existedSettings.ID == "" {
		err = services.DB.Model(&model.ChannelModulesSettings{}).Create(map[string]interface{}{
			"id":        uuid.NewV4().String(),
			"type":      "youtube_song_requests",
			"settings":  bytes,
			"channelId": channelId,
		}).Error
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
