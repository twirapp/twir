package roles_users

import (
	"github.com/gofiber/fiber/v2"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"
	"gorm.io/gorm"
	"net/http"
)

type roleUser struct {
	model.ChannelRoleUser
	UserName        string `json:"userName"`
	UserDisplayName string `json:"userDisplayName"`
	UserAvatar      string `json:"userAvatar"`
}

func getUsersService(roleId string) ([]*roleUser, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role := &model.ChannelRole{}
	if err := db.Preload("Users").Where("id = ?", roleId).First(role).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "Role not found")
	}

	response := make([]*roleUser, 0, len(role.Users))
	if len(role.Users) == 0 {
		return response, nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: lo.Map(role.Users, func(user *model.ChannelRoleUser, _ int) string {
			return user.UserID
		}),
	})

	if err != nil || twitchUsers.ErrorMessage != "" {
		logger.Error(err, twitchUsers.ErrorMessage)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	for _, user := range role.Users {
		twitchUser, ok := lo.Find(twitchUsers.Data.Users, func(twitchUser helix.User) bool {
			return twitchUser.ID == user.UserID
		})

		if !ok {
			continue
		}

		response = append(response, &roleUser{
			ChannelRoleUser: *user,
			UserName:        twitchUser.Login,
			UserDisplayName: twitchUser.DisplayName,
			UserAvatar:      twitchUser.ProfileImageURL,
		})
	}

	return response, nil
}

func updateUsersService(roleId string, userNames []string) error {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)
	db := do.MustInvoke[*gorm.DB](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	role := &model.ChannelRole{}
	if err := db.Preload("Users").Where("id = ?", roleId).First(role).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Role not found")
	}

	if len(userNames) == 0 {
		if err = db.Where(`"roleId" = ?`, role.ID).Delete(&model.ChannelRoleUser{}).Error; err != nil {
			logger.Error(err)
			return fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		return nil
	}

	twitchUsers, err := twitchClient.GetUsers(&helix.UsersParams{
		Logins: userNames,
	})

	if err != nil || twitchUsers.ErrorMessage != "" {
		logger.Error(err, twitchUsers.ErrorMessage)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	roleUsers := make([]*model.ChannelRoleUser, 0, len(twitchUsers.Data.Users))
	for _, twitchUser := range twitchUsers.Data.Users {
		roleUsers = append(roleUsers, &model.ChannelRoleUser{
			RoleID: role.ID,
			UserID: twitchUser.ID,
		})
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err = tx.Where(`"roleId" = ?`, role.ID).Delete(&model.ChannelRoleUser{}).Error; err != nil {
			return err
		}

		if err = tx.Create(roleUsers).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	return nil
}
