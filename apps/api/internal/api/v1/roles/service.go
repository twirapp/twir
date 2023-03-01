package roles

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func getRole(id string) *model.ChannelRole {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	role := &model.ChannelRole{}
	err := db.
		Where(`"id" = ?`, id).
		Preload("Users").
		First(&role).Error
	if err != nil {
		logger.Error(err)
		return nil
	}

	return role
}

func getRolesService(channelId string) ([]*model.ChannelRole, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	channelsRoles := []*model.ChannelRole{}
	err := db.
		Where(`"channelId" = ?`, channelId).
		Preload("Users").
		Find(&channelsRoles).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return channelsRoles, nil
}

func updateRoleService(channelId, roleId string, dto *roleDto) (*model.ChannelRole, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	role := &model.ChannelRole{}
	err := db.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "Role not found")
	}

	role.Name = dto.Name
	role.Permissions = dto.Permissions

	err = db.Save(&role).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(roleId)

	return newRole, nil
}

func deleteRoleService(channelId, roleId string) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	role := &model.ChannelRole{}
	err := db.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Role not found")
	}

	if role.Type != model.ChannelRoleTypeCustom {
		return fiber.NewError(http.StatusForbidden, "System role can't be deleted")
	}

	err = db.Delete(&role).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func createRoleService(channelId string, dto *roleDto) (*model.ChannelRole, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	role := &model.ChannelRole{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Name:        dto.Name,
		Type:        model.ChannelRoleTypeCustom,
		Permissions: dto.Permissions,
	}

	err := db.Create(&role).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(role.ID)
	return newRole, nil
}
