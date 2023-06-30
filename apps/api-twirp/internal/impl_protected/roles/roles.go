package roles

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api-twirp/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/roles"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Roles struct {
	*impl_deps.Deps
}

func (c *Roles) convertEntity(entity *model.ChannelRole) (*roles.Role, error) {
	settings := &model.ChannelRoleSettings{}
	if err := json.Unmarshal(entity.Settings, settings); err != nil {
		return nil, err
	}

	return &roles.Role{
		Id:          entity.ID,
		ChannelId:   entity.ChannelID,
		Name:        entity.Name,
		Type:        entity.Type.String(),
		Permissions: entity.Permissions,
		Settings: &roles.Role_Settings{
			RequiredWatchTime:         settings.RequiredWatchTime,
			RequiredMessages:          settings.RequiredMessages,
			RequiredUserChannelPoints: settings.RequiredUsedChannelPoints,
		},
		Users: lo.Map(entity.Users, func(u *model.ChannelRoleUser, _ int) *roles.Role_User {
			return &roles.Role_User{
				Id:     u.ID,
				UserId: u.UserID,
				RoleId: u.RoleID,
			}
		}),
	}, nil
}

func (c *Roles) RolesGetAll(
	ctx context.Context,
	_ *emptypb.Empty,
) (*roles.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []*model.ChannelRole
	if err := c.Db.
		WithContext(ctx).
		Preload("Users").
		Where(`"channelId" = ?`, dashboardId).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}

	res := make([]*roles.Role, 0, len(entities))
	for _, entity := range entities {
		converted, err := c.convertEntity(entity)
		if err != nil {
			return nil, err
		}

		res = append(res, converted)
	}

	return &roles.GetAllResponse{
		Roles: res,
	}, nil
}

func (c *Roles) RolesUpdate(
	ctx context.Context,
	request *roles.UpdateRequest,
) (*roles.Role, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelRole{}
	if err := c.Db.
		WithContext(ctx).
		Preload("Users").
		Where(`"channelId" = ?`, dashboardId).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	entity.Name = request.Role.Name
	entity.Permissions = request.Role.Permissions
	entity.Users = lo.Map(request.Role.Users, func(u *roles.Role_User, _ int) *model.ChannelRoleUser {
		return &model.ChannelRoleUser{
			ID:     uuid.New().String(),
			UserID: u.UserId,
			RoleID: u.RoleId,
		}
	})

	txErr := c.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(`"roleId" = ?`, entity.ID).Delete(&model.ChannelRoleUser{}).Error; err != nil {
			return err
		}

		err := tx.WithContext(ctx).Save(entity).Error
		return err
	})
	if txErr != nil {
		return nil, txErr
	}

	return c.convertEntity(entity)
}

func (c *Roles) RolesDelete(
	ctx context.Context,
	request *roles.DeleteRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
}

func (c *Roles) RolesCreate(
	ctx context.Context,
	request *roles.CreateRequest,
) (*roles.Role, error) {
	dashboardId := ctx.Value("dashboardId").(string)
}
