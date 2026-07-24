package dashboard_access

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	model "github.com/twirapp/twir/libs/gomodels"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

const canAccessDashboardPermission = "CAN_ACCESS_DASHBOARD"

type Opts struct {
	fx.In

	ChannelService *channelservice.ChannelService
	Gorm           *gorm.DB
}

func NewFx(opts Opts) *Service {
	return New(opts.ChannelService, &gormStore{gorm: opts.Gorm})
}

type Service struct {
	channels ChannelReader
	store    Store
}

type Subject struct {
	ID         string
	IsBotAdmin bool
}

type ChannelReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error)
}

type Store interface {
	GetLegacyChannel(context.Context, uuid.UUID) (model.Channels, error)
	GetRoles(context.Context, uuid.UUID, string) ([]model.ChannelRole, error)
	GetUserStat(context.Context, string, uuid.UUID) (model.UsersStats, error)
}

func New(channels ChannelReader, store Store) *Service {
	return &Service{channels: channels, store: store}
}

func (s *Service) CanAccess(
	ctx context.Context,
	subject Subject,
	channelID uuid.UUID,
	permission string,
) (bool, error) {
	if subject.IsBotAdmin {
		return true, nil
	}

	isOwner, err := s.IsOwner(ctx, subject.ID, channelID)
	if err != nil {
		return false, err
	}
	if isOwner {
		return true, nil
	}

	roles, err := s.store.GetRoles(ctx, channelID, subject.ID)
	if err != nil {
		return false, fmt.Errorf("get channel roles: %w", err)
	}
	userStat, err := s.store.GetUserStat(ctx, subject.ID, channelID)
	if err != nil {
		return false, fmt.Errorf("get user stat: %w", err)
	}

	return rolesAllowAccess(roles, subject.ID, userStat, permission), nil
}

func (s *Service) IsOwner(ctx context.Context, userID string, channelID uuid.UUID) (bool, error) {
	if s == nil || s.channels == nil || s.store == nil {
		return false, fmt.Errorf("dashboard access service is not configured")
	}

	channel, err := s.channels.GetChannelByID(ctx, channelID)
	if err != nil {
		return false, fmt.Errorf("get channel: %w", err)
	}
	if len(channel.Bindings) > 0 {
		for _, binding := range channel.Bindings {
			if binding.UserID.String() == userID {
				return true, nil
			}
		}

		return false, nil
	}

	legacyChannel, err := s.store.GetLegacyChannel(ctx, channelID)
	if err != nil {
		return false, fmt.Errorf("get legacy channel: %w", err)
	}

	return legacyChannel.IsOwner(userID), nil
}

func rolesAllowAccess(
	roles []model.ChannelRole,
	userID string,
	userStat model.UsersStats,
	permission string,
) bool {
	roleToStats := map[model.ChannelRoleEnum]bool{
		model.ChannelRoleTypeModerator:  userStat.IsMod,
		model.ChannelRoleTypeVip:        userStat.IsVip,
		model.ChannelRoleTypeSubscriber: userStat.IsSubscriber,
	}

	for i, role := range roles {
		if roleToStats[role.Type] {
			roles[i].Users = append(
				role.Users,
				&model.ChannelRoleUser{ID: "", UserID: userID, RoleID: role.ID},
			)
		}
	}

	for _, role := range roles {
		if len(role.Users) == 0 || len(role.Permissions) == 0 {
			continue
		}

		for _, roleUser := range role.Users {
			if roleUser.UserID != userID {
				continue
			}
			if permission == "" {
				return true
			}

			for _, rolePermission := range role.Permissions {
				if rolePermission == canAccessDashboardPermission || rolePermission == permission {
					return true
				}
			}
		}
	}

	return false
}

type gormStore struct {
	gorm *gorm.DB
}

func (s *gormStore) GetLegacyChannel(ctx context.Context, channelID uuid.UUID) (model.Channels, error) {
	var channel model.Channels
	if err := s.gorm.WithContext(ctx).Where("id = ?", channelID).First(&channel).Error; err != nil {
		return model.Channels{}, err
	}

	return channel, nil
}

func (s *gormStore) GetRoles(ctx context.Context, channelID uuid.UUID, userID string) ([]model.ChannelRole, error) {
	var channelRoles []model.ChannelRole
	if err := s.gorm.
		WithContext(ctx).
		Where(`"channelId" = ?::uuid`, channelID).
		Preload("Users", `user_id = ?`, userID).
		Find(&channelRoles).
		Error; err != nil {
		return nil, err
	}

	return channelRoles, nil
}

func (s *gormStore) GetUserStat(ctx context.Context, userID string, channelID uuid.UUID) (model.UsersStats, error) {
	var userStat model.UsersStats
	err := s.gorm.
		WithContext(ctx).
		Where(`user_id = ? AND channel_id = ?::uuid`, userID, channelID).
		First(&userStat).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return userStat, nil
	}
	if err != nil {
		return model.UsersStats{}, err
	}

	return userStat, nil
}
