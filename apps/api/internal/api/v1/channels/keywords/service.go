package keywords

import (
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func (c *Keywords) getService(channelId string) ([]model.ChannelsKeywords, error) {
	keywords := []model.ChannelsKeywords{}
	err := c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&keywords).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get keywords")
	}

	return keywords, nil
}

func (c *Keywords) postService(
	channelId string,
	dto *keywordDto,
) (*model.ChannelsKeywords, error) {
	existedKeyword := model.ChannelsKeywords{}
	err := c.services.Gorm.Where(`"channelId" = ? AND "text" = ?`, channelId, dto.Text).
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
	err = c.services.Gorm.Save(&newKeyword).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create keyword")
	}

	return &newKeyword, nil
}

func (c *Keywords) deleteService(keywordId string) error {
	keyword := c.getById(keywordId)
	if keyword == nil {
		return fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	err := c.services.Gorm.Delete(keyword).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete keyword")
	}

	return nil
}

func (c *Keywords) putService(
	keywordId string,
	dto *keywordDto,
) (*model.ChannelsKeywords, error) {
	currentKeyword := c.getById(keywordId)
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

	err := c.services.Gorm.Model(currentKeyword).Select("*").Updates(newKeyword).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update keyword")
	}

	return &newKeyword, nil
}

func (c *Keywords) patchService(
	channelId,
	keywordId string,
	dto *keywordPatchDto,
) (*model.ChannelsKeywords, error) {
	keyword := model.ChannelsKeywords{}
	err := c.services.Gorm.Where(`"channelId" = ? AND "id" = ?`, channelId, keywordId).
		Find(&keyword).
		Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if keyword.ID == "" {
		return nil, fiber.NewError(http.StatusNotFound, "keyword not found")
	}

	keyword.Enabled = *dto.Enabled

	err = c.services.Gorm.Save(&keyword).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot update keyword")
	}

	return &keyword, nil
}
