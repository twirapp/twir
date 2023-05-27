package timers

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	lo_parallel "github.com/samber/lo/parallel"
	"github.com/satont/tsuwari/apps/scheduler/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
	"time"
)

type DefaultCommandsTimer struct {
	services *types.Services
}

func NewDefaultCommands(services *types.Services) *DefaultCommandsTimer {
	d := &DefaultCommandsTimer{}

	return d
}

var neededRoles = []model.ChannelRoleEnum{
	model.ChannelRoleTypeBroadcaster,
	model.ChannelRoleTypeModerator,
	model.ChannelRoleTypeVip,
	model.ChannelRoleTypeSubscriber,
}

func (c *DefaultCommandsTimer) Run(ctx context.Context) {
	timeTick := lo.If(c.services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)

	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				break
			case <-ticker.C:
				defaultCommands, err := c.services.Grpc.Parser.GetDefaultCommands(ctx, &emptypb.Empty{})
				if err != nil {
					zap.S().Error(err)
					continue
				}

				var channels []model.Channels
				err = c.services.Gorm.Find(&channels).Error
				if err != nil {
					zap.S().Error(err)
					continue
				}

				lo_parallel.ForEach(channels, func(item model.Channels, index int) {
					err = c.CreateCommandsAndRoles(defaultCommands, item.ID)
					if err != nil {
						zap.S().Error(err)
					}
				})
			}
		}
	}()
}

func (c *DefaultCommandsTimer) CreateCommandsAndRoles(
	defaultCommands *parser.GetDefaultCommandsResponse,
	channelId string,
) error {
	var currentRoles []model.ChannelRole
	var currentCommands []model.ChannelsCommands

	err := c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&currentCommands).Error
	if err != nil {
		return err
	}

	err = c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&currentRoles).Error
	if err != nil {
		return err
	}

	var roles []*model.ChannelRole

	err = c.services.Gorm.Transaction(func(tx *gorm.DB) error {
		for _, role := range neededRoles {
			existedRole, ok := lo.Find(currentRoles, func(item model.ChannelRole) bool {
				return item.Type == role
			})

			if ok {
				roles = append(roles, &existedRole)
				continue
			}

			newRole := &model.ChannelRole{
				ID:        uuid.New().String(),
				ChannelID: channelId,
				Name:      role.String(),
				Type:      role,
				Permissions: lo.
					If(
						role == model.ChannelRoleTypeBroadcaster,
						pq.StringArray{model.RolePermissionCanAccessDashboard.String()},
					).
					Else(pq.StringArray{}),
			}
			err = tx.Save(newRole).Error
			if err != nil {
				return err
			}
			roles = append(roles, newRole)
		}

		for _, command := range defaultCommands.List {
			_, ok := lo.Find(currentCommands, func(item model.ChannelsCommands) bool {
				return item.DefaultName.Valid && item.DefaultName.String == command.Name
			})

			if ok {
				continue
			}

			roleIds := pq.StringArray{}

			for _, roleName := range command.RolesNames {
				role, ok := lo.Find(roles, func(role *model.ChannelRole) bool {
					return role.Type.String() == roleName
				})
				if !ok {
					continue
				}

				roleIds = append(roleIds, role.ID)
			}

			newCommand := &model.ChannelsCommands{
				ID:                 uuid.New().String(),
				ChannelID:          channelId,
				Enabled:            true,
				Name:               command.Name,
				Description:        null.StringFrom(command.Description),
				Visible:            command.Visible,
				RolesIDS:           roleIds,
				Module:             command.Module,
				IsReply:            command.IsReply,
				KeepResponsesOrder: command.KeepResponsesOrder,
				Aliases:            command.Aliases,
				Default:            true,
				DefaultName:        null.StringFrom(command.Name),
				CooldownType:       "GLOBAL",
			}

			err = tx.Save(newCommand).Error
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
