package users

import (
	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *adminUsers) postService(dto *ignoredUsersPostDto) error {
	for _, user := range dto.Users {
		newIgnoredUser := model.IgnoredUser{
			ID:          user.ID,
			Login:       null.StringFrom(user.Name),
			DisplayName: null.StringFrom(user.DisplayName),
		}
		err := c.services.Gorm.Save(&newIgnoredUser).Error
		if err != nil {
			c.services.Logger.Error(err)
			continue
		}
	}

	return nil
}
