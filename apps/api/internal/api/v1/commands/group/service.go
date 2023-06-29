package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/do"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"gorm.io/gorm"
)

func getGroupsService(channelId string, services types.Services) ([]model.ChannelCommandGroup, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)

	var groups []model.ChannelCommandGroup
	err := db.
		Where(`"channelId" = ?`, channelId).
		Find(&groups).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return groups, nil
}

func createGroupService(channelId string, dto *groupDto) (*model.ChannelCommandGroup, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)

	group := &model.ChannelCommandGroup{
		ID:        uuid.New().String(),
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := db.Create(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}

func deleteGroupService(channelId, groupId string) error {
	db := do.MustInvoke[*gorm.DB](di.Provider)

	err := db.Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).Delete(&model.ChannelCommandGroup{}).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func updateGroupService(channelId, groupId string, dto *groupDto) (*model.ChannelCommandGroup, error) {
	db := do.MustInvoke[*gorm.DB](di.Provider)

	group := &model.ChannelCommandGroup{
		ID:        groupId,
		ChannelID: channelId,
		Name:      dto.Name,
		Color:     dto.Color,
	}

	err := db.
		Model(&model.ChannelCommandGroup{}).
		Where(`"channelId" = ? AND "id" = ?`, channelId, groupId).
		Updates(group).Error
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return group, nil
}
