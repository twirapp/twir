package commands_group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *CommandsGroup) getService(channelId string) ([]model.ChannelCommandGroup, error) {
	var groups []model.ChannelCommandGroup
	err := c.services.Gorm.
		Where(`"channelId" = ?`, channelId).
		Find(&groups).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return groups, nil
}

func (c *CommandsGroup) postService(
	channelId string,
	dto *groupDto,
) (*model.ChannelCommandGroup, error) {
	group := &model.ChannelCommandGroup{
		ID:        uuid.New().String(),
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := c.services.Gorm.Create(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}

func (c *CommandsGroup) deleteService(channelId, groupId string) error {
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).
		Delete(&model.ChannelCommandGroup{}).
		Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func (c *CommandsGroup) putService(
	channelId,
	groupId string,
	dto *groupDto,
) (*model.ChannelCommandGroup, error) {
	group := &model.ChannelCommandGroup{
		ID:        groupId,
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := c.services.Gorm.
		Model(&model.ChannelCommandGroup{}).
		Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).
		Updates(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}
