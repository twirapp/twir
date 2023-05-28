package roles

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type Role struct {
	*model.ChannelRole
	Settings *model.ChannelRoleSettings `json:"settings"`
}

func getRole(id string) *Role {
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

	settings := &model.ChannelRoleSettings{}
	err = json.Unmarshal(role.Settings, settings)

	r := &Role{
		ChannelRole: role,
		Settings:    settings,
	}

	return r
}

func getRolesService(channelId string) ([]Role, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	var channelsRoles []*model.ChannelRole
	err := db.
		Where(`"channelId" = ?`, channelId).
		Preload("Users").
		Find(&channelsRoles).Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	roles := make([]Role, len(channelsRoles))
	for i, role := range channelsRoles {
		settings := &model.ChannelRoleSettings{}
		err = json.Unmarshal(role.Settings, settings)
		if err != nil {
			logger.Error(err)
		}

		r := Role{
			ChannelRole: role,
			Settings:    settings,
		}
		roles[i] = r
	}

	return roles, nil
}

func updateRoleService(channelId, roleId string, dto *roleDto) (*Role, error) {
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

	spew.Dump(dto.Settings)
	settings, err := json.Marshal(dto.Settings)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role.Settings = settings

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

func createRoleService(channelId string, dto *roleDto) (*Role, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	settingsBytes, err := json.Marshal(dto.Settings)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role := &model.ChannelRole{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Name:        dto.Name,
		Type:        model.ChannelRoleTypeCustom,
		Permissions: dto.Permissions,
		Settings:    settingsBytes,
	}

	err = db.Create(&role).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(role.ID)
	return newRole, nil
}
