package roles

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	model "github.com/satont/tsuwari/libs/gomodels"
	uuid "github.com/satori/go.uuid"
)

func (c *Roles) getById(id string) *model.ChannelRole {
	role := &model.ChannelRole{}
	err := c.services.Gorm.
		Where(`"id" = ?`, id).
		Preload("Users").
		First(&role).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil
	}

	return role
}

func (c *Roles) getService(channelId string) ([]*model.ChannelRole, error) {
	channelsRoles := []*model.ChannelRole{}
	err := c.services.Gorm.
		Where(`"channelId" = ?`, channelId).
		Preload("Users").
		Find(&channelsRoles).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	return channelsRoles, nil
}

func (c *Roles) putService(channelId, roleId string, dto *roleDto) (*model.ChannelRole, error) {
	role := &model.ChannelRole{}
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusNotFound, "Role not found")
	}

	role.Name = dto.Name
	role.Permissions = dto.Permissions

	err = c.services.Gorm.Save(&role).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := c.getById(roleId)

	return newRole, nil
}

func (c *Roles) deleteService(channelId, roleId string) error {
	role := &model.ChannelRole{}
	err := c.services.Gorm.
		Where(`"channelId" = ? AND "id" = ?`, channelId, roleId).
		First(&role).Error
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Role not found")
	}

	if role.Type != model.ChannelRoleTypeCustom {
		return fiber.NewError(http.StatusForbidden, "System role can't be deleted")
	}

	err = c.services.Gorm.Delete(&role).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}

func (c *Roles) postService(channelId string, dto *roleDto) (*model.ChannelRole, error) {
	role := &model.ChannelRole{
		ID:          uuid.NewV4().String(),
		ChannelID:   channelId,
		Name:        dto.Name,
		Type:        model.ChannelRoleTypeCustom,
		Permissions: dto.Permissions,
	}

	err := c.services.Gorm.Create(&role).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	newRole := c.getById(role.ID)
	return newRole, nil
}
