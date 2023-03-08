package roles

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func getRole(id string, services *types.Services) *model.ChannelRole {
	role := &model.ChannelRole{}
	err := services.Gorm.
		Where(`"id" = ?`, id).
		Preload("Users").
		First(&role).Error
	if err != nil {
		services.Logger.Error(err)
		return nil
	}

	return role
}

func getRolesService(channelId string, services *types.Services) ([]*model.ChannelRole, error) {
	channelsRoles := []*model.ChannelRole{}
	err := services.Gorm.
		Where(`"channelId" = ?`, channelId).
		Preload("Users").
		Find(&channelsRoles).Error
	if err != nil {
		services.Logger.Error(err)
		return nil, err
	}

	return channelsRoles, nil
}

func updateRoleService(channelId, roleId string, dto *roleDto, services *types.Services) (*model.ChannelRole, error) {
	role := &model.ChannelRole{}
	err := services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "Role not found")
	}

	role.Name = dto.Name
	role.Permissions = dto.Permissions

	err = services.Gorm.Save(&role).Error
	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(roleId, services)

	return newRole, nil
}

func deleteRoleService(channelId, roleId string, services *types.Services) error {
	role := &model.ChannelRole{}
	err := services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Role not found")
	}

	if role.Type != model.ChannelRoleTypeCustom {
		return fiber.NewError(http.StatusForbidden, "System role can't be deleted")
	}

	err = services.Gorm.Delete(&role).Error
	if err != nil {
		services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func createRoleService(channelId string, dto *roleDto, services *types.Services) (*model.ChannelRole, error) {
	role := &model.ChannelRole{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Name:        dto.Name,
		Type:        model.ChannelRoleTypeCustom,
		Permissions: dto.Permissions,
	}

	err := services.Gorm.Create(&role).Error
	if err != nil {
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(role.ID, services)
	return newRole, nil
}
