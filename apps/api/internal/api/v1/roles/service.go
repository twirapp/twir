package roles

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"net/http"
)

func getRole(id string) *model.ChannelRole {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)

	role := &model.ChannelRole{}
	err := db.
		Where(`"id" = ?`, id).
		Preload("Permissions").
		Preload("Permissions.Flag").
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
		Preload("Permissions").
		Preload("Permissions.Flag").
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
		logger.Error(err)
		return nil, fiber.NewError(http.StatusNotFound, "Role not found")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = db.
			Delete(&model.ChannelRolePermission{}, `"roleId" = ?`, role.ID).
			Error
		if err != nil {
			return err
		}

		for _, flag := range dto.Flags {
			dbFlag := &model.RoleFlag{}
			err = db.Where("flag = ?", flag).First(&dbFlag).Error
			if err != nil {
				return err
			}

			err = db.Create(&model.ChannelRolePermission{
				ID:     uuid.NewV4().String(),
				RoleID: role.ID,
				FlagID: dbFlag.ID,
			}).Error
			if err != nil {
				return err
			}
		}

		err = db.Model(&role).Update("name", dto.Name).Error
		if err != nil {
			return err
		}

		return nil
	})

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
		logger.Error(err)
		return fiber.NewError(http.StatusNotFound, "Role not found")
	}

	if role.System {
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
		ID:        uuid.NewV4().String(),
		ChannelID: channelId,
		Name:      dto.Name,
		System:    false,
		Type:      model.ChannelRoleTypeCustom,
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		err := db.Create(&role).Error
		if err != nil {
			return err
		}

		for _, flag := range dto.Flags {
			dbFlag := &model.RoleFlag{}
			err = db.Where("flag = ?", flag).First(&dbFlag).Error
			if err != nil {
				return err
			}

			err = db.Create(&model.ChannelRolePermission{
				ID:     uuid.NewV4().String(),
				RoleID: role.ID,
				FlagID: dbFlag.ID,
			}).Error
			if err != nil {
				logger.Error(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := getRole(role.ID)
	return newRole, nil
}
