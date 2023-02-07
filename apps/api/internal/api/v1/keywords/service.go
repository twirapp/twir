package keywords

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsKeywords, error) {
	keywords := []model.ChannelsKeywords{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&keywords).Error
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
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	existedKeyword := model.ChannelsKeywords{}
	err := services.DB.Where(`"channelId" = ? AND "text" = ?`, channelId, dto.Text).
		First(&existedKeyword).
		Error
	if err == nil {
		return nil, fiber.NewError(400, "keyword with that text already exists")
	}

	newKeyword := model.ChannelsKeywords{
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		Text:      dto.Text,
		Response:  dto.Response,
		Enabled:   *dto.Enabled,
		Cooldown:  int(dto.Cooldown),
		IsRegular: *dto.IsRegular,
		IsReply:   *dto.IsReply,
		Usages:    *dto.Usages,
	}
	err = services.DB.Save(&newKeyword).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create keyword")
	}

	return &newKeyword, nil
}

func handleDelete(keywordId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	keyword := getById(services.DB, keywordId)
	if keyword == nil {
		return fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	err := services.DB.Delete(keyword).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete keyword")
	}

	return nil
}

func handleUpdate(
	keywordId string,
	dto *keywordDto,
	services types.Services,
) (*model.ChannelsKeywords, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

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
		Cooldown:  int(dto.Cooldown),
		IsReply:   *dto.IsReply,
		IsRegular: *dto.IsRegular,
		Usages:    *dto.Usages,
	}

	err := services.DB.Model(currentKeyword).Select("*").Updates(newKeyword).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update keyword")
	}

	return &newKeyword, nil
}

func handlePatch(
	channelId,
	keywordId string,
	dto *keywordPatchDto,
	services types.Services,
) (*model.ChannelsKeywords, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	keyword := model.ChannelsKeywords{}
	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, keywordId).
		Find(&keyword).
		Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if keyword.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	keyword.Enabled = *dto.Enabled

	err = services.DB.Save(&keyword).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update keyword")
	}

	return &keyword, nil
}
