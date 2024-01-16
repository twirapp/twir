package roles

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/roles"
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
			RequiredWatchTime:         int32(settings.RequiredWatchTime),
			RequiredMessages:          settings.RequiredMessages,
			RequiredUserChannelPoints: int32(settings.RequiredUsedChannelPoints),
		},
		Users: lo.Map(
			entity.Users, func(u *model.ChannelRoleUser, _ int) *roles.Role_User {
				return &roles.Role_User{
					Id:     u.ID,
					UserId: u.UserID,
					RoleId: u.RoleID,
				}
			},
		),
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
		Group(`"id"`).
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
		Where(`"channelId" = ? and id = ?`, dashboardId, request.Id).
		First(&entity).
		Error; err != nil {
		return nil, err
	}

	entity.Name = request.Role.Name
	entity.Permissions = request.Role.Permissions
	entity.Users = lo.Map(
		request.Role.Users,
		func(u *roles.CreateRequest_User, _ int) *model.ChannelRoleUser {
			return &model.ChannelRoleUser{
				ID:     uuid.New().String(),
				UserID: u.UserId,
				RoleID: entity.ID,
			}
		},
	)

	if entity.Permissions == nil {
		entity.Permissions = pq.StringArray{}
	}

	txErr := c.Db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Where(
				`"roleId" = ?`,
				entity.ID,
			).Delete(&model.ChannelRoleUser{}).Error; err != nil {
				return err
			}

			err := tx.WithContext(ctx).Save(entity).Error
			return err
		},
	)
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
	if err := c.Db.
		WithContext(ctx).
		Where(
			`"id" = ? AND "channelId" = ? AND "type" = ?`,
			request.Id,
			dashboardId,
			model.ChannelRoleTypeCustom.String(),
		).
		Delete(&model.ChannelRole{}).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Roles) RolesCreate(
	ctx context.Context,
	request *roles.CreateRequest,
) (*roles.Role, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	settings := &model.ChannelRoleSettings{
		RequiredWatchTime:         int64(request.Settings.RequiredWatchTime),
		RequiredMessages:          request.Settings.RequiredMessages,
		RequiredUsedChannelPoints: int64(request.Settings.RequiredUserChannelPoints),
	}

	settingsBytes, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	entity := &model.ChannelRole{
		ID:          uuid.New().String(),
		ChannelID:   dashboardId,
		Name:        request.Name,
		Type:        model.ChannelRoleTypeCustom,
		Permissions: request.Permissions,
		Settings:    settingsBytes,
		Users: lo.Map(
			request.Users, func(u *roles.CreateRequest_User, _ int) *model.ChannelRoleUser {
				return &model.ChannelRoleUser{
					ID:     uuid.New().String(),
					UserID: u.UserId,
				}
			},
		),
	}

	if entity.Permissions == nil {
		entity.Permissions = pq.StringArray{}
	}

	if err := c.Db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity)
}
