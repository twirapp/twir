package keywords

import (
	"net/http"
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsKeywords, error) {
	keywords := []model.ChannelsKeywords{}
	err := services.DB.Find(&keywords).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get keywords")
	}

	return keywords, nil
}

func handlePost(
	channelId string,
	dto *keywordDto,
	services types.Services,
) (*model.ChannelsKeywords, error) {
	existedKeyword := model.ChannelsKeywords{}
	err := services.DB.Where(`"text" = ?`, dto.Text).First(&existedKeyword).Error
	if err == nil {
		return nil, fiber.NewError(400, "keyword with that text already exists")
	}

	newKeyword := model.ChannelsKeywords{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		Text:      dto.Text,
		Response:  dto.Response,
		Enabled:   *dto.Enabled,
		Cooldown:  null.IntFrom(int64(dto.Cooldown)),
	}
	err = services.DB.Save(&newKeyword).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create keyword")
	}

	return &newKeyword, nil
}

func handleDelete(keywordId string, services types.Services) error {
	keyword := getById(services.DB, keywordId)
	if keyword == nil {
		return fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	err := services.DB.Delete(keyword).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete keyword")
	}

	return nil
}

func handleUpdate(
	keywordId string,
	dto *keywordDto,
	services types.Services,
) (*model.ChannelsKeywords, error) {
	currentKeyword := getById(services.DB, keywordId)
	if currentKeyword == nil {
		return nil, fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	newKeyword := model.ChannelsKeywords{
		ID:        currentKeyword.ID,
		ChannelID: currentKeyword.ChannelID,
		Text:      dto.Text,
		Response:  dto.Response,
		Enabled:   *dto.Enabled,
		Cooldown:  null.IntFrom(int64(dto.Cooldown)),
		IsReply:   *dto.IsReply,
	}

	err := services.DB.Model(currentKeyword).Select("*").Updates(newKeyword).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update keyword")
	}

	return &newKeyword, nil
}
