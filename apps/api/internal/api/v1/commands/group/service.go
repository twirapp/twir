package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func getGroupsService(channelId string, services *types.Services) ([]model.ChannelCommandGroup, error) {
	var groups []model.ChannelCommandGroup
	err := services.Gorm.
		Where(`"channelId" = ?`, channelId).
		Find(&groups).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return groups, nil
}

func createGroupService(channelId string, dto *groupDto, services *types.Services) (*model.ChannelCommandGroup, error) {
	group := &model.ChannelCommandGroup{
		ID:        uuid.New().String(),
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := services.Gorm.Create(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}

func deleteGroupService(channelId, groupId string, services *types.Services) error {
	err := services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).
		Delete(&model.ChannelCommandGroup{}).
		Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func updateGroupService(channelId, groupId string, dto *groupDto, services *types.Services) (*model.ChannelCommandGroup, error) {
	group := &model.ChannelCommandGroup{
		ID:        groupId,
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := services.Gorm.
		Model(&model.ChannelCommandGroup{}).
		Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).
		Updates(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}
